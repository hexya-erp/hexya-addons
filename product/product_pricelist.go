// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.ProductPricelist().DeclareModel()
	h.ProductPricelist().SetDefaultOrder("Sequence ASC", "ID DESC")

	h.ProductPricelist().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Pricelist Name", Required: true, Translate: true},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If unchecked, it will allow you to hide the pricelist without removing it."},
		"Items": models.One2ManyField{String: "Pricelist Items", RelationModel: h.ProductPricelistItem(),
			ReverseFK: "Pricelist", JSON: "item_ids", NoCopy: false,
			Default: func(env models.Environment) interface{} {
				listItems := h.ProductPricelistItem().NewSet(env)
				values, _ := listItems.DataStruct(listItems.DefaultGet())
				values.ComputePrice = "formula"
				return listItems.Create(values)
			}},
		"Currency": models.Many2OneField{RelationModel: h.Currency(),
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company().Currency()
			}, Required: true},
		"Company":       models.Many2OneField{RelationModel: h.Company()},
		"Sequence":      models.IntegerField{Default: models.DefaultValue(16)},
		"CountryGroups": models.Many2ManyField{RelationModel: h.CountryGroup(), JSON: "country_group_ids"},
	})

	h.ProductPricelist().Methods().NameGet().Extend("",
		func(rs h.ProductPricelistSet) string {
			return fmt.Sprintf("%s (%s)", rs.Name(), rs.Currency().Name())
		})

	h.ProductPricelist().Methods().SearchByName().Extend("",
		func(rs h.ProductPricelistSet, name string, op operator.Operator, additionalCondition q.ProductPricelistCondition, limit int) h.ProductPricelistSet {
			return rs.Super().SearchByName(name, op, additionalCondition, limit)
		})

	h.ProductPricelist().Methods().ComputePriceRule().DeclareMethod(
		`ComputePriceRule is the low-level method computing the price of the given product according to this
		price list. Price depends on quantity, partner and date, and is given for the uom.

		If date or uom are not given, this function will try to read them from the context 'date' and 'uom' keys`,
		func(rs h.ProductPricelistSet, product h.ProductProductSet, quantity float64, partner h.PartnerSet,
			date dates.Date, uom h.ProductUomSet) (float64, h.ProductPricelistItemSet) {

			rs.EnsureOne()
			if date.IsZero() {
				date = dates.Today()
				if rs.Env().Context().HasKey("date") {
					date = rs.Env().Context().GetDate("date")
				}
			}
			if uom.IsEmpty() && rs.Env().Context().HasKey("uom") {
				uom = h.ProductUom().NewSet(rs.Env()).Browse([]int64{rs.Env().Context().GetInteger("uom")})
			}
			if !uom.IsEmpty() {
				product = product.WithContext("uom", uom.ID())
			}
			if product.IsEmpty() {
				return 0, h.ProductPricelistItem().NewSet(rs.Env())
			}

			categs := h.ProductCategory().NewSet(rs.Env())
			for categ := product.Categ(); !categ.IsEmpty(); categ = categ.Parent() {
				categs = categs.Union(categ)
			}

			prodTmpl := product.ProductTmpl()

			// Load all rules
			tmplCond := q.ProductPricelistItem().ProductTmpl().IsNull().Or().ProductTmpl().Equals(prodTmpl)
			prodCond := q.ProductPricelistItem().Product().IsNull().Or().Product().Equals(product)
			categCond := q.ProductPricelistItem().Categ().IsNull().Or().Categ().In(categs)
			dateStartCond := q.ProductPricelistItem().DateStart().IsNull().Or().DateStart().LowerOrEqual(date)
			dateEndCond := q.ProductPricelistItem().DateEnd().IsNull().Or().DateEnd().LowerOrEqual(date)

			items := h.ProductPricelistItem().Search(rs.Env(),
				q.ProductPricelistItem().Pricelist().Equals(rs).
					AndCond(tmplCond).
					AndCond(prodCond).
					AndCond(categCond).
					AndCond(dateStartCond).
					AndCond(dateEndCond)).OrderBy("AppliedOn", "MinQuantity DESC", "Categ.Name")

			var price float64
			suitableRule := h.ProductPricelistItem().NewSet(rs.Env())
			// Final unit price is computed according to `qty` in the `qty_uom_id` UoM.
			// An intermediary unit price may be computed according to a different UoM, in
			// which case the price_uom_id contains that UoM.
			// The final price will be converted to match `qty_uom_id`.
			qtyUom := product.Uom()
			if rs.Env().Context().HasKey("uom") {
				qtyUom = h.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
			}
			priceUom := product.Uom()
			qtyInProductUom := quantity
			if !qtyUom.Equals(product.Uom()) {
				if qtyUom.Category().Equals(product.Uom().Category()) {
					qtyInProductUom = qtyUom.ComputeQuantity(quantity, product.Uom(), true)
				}
			}
			price = product.PriceCompute(h.ProductProduct().ListPrice(),
				h.ProductUom().NewSet(rs.Env()), h.Currency().NewSet(rs.Env()), h.Company().NewSet(rs.Env()))

			for _, rule := range items.Records() {
				if rule.MinQuantity() != 0 && qtyInProductUom < rule.MinQuantity() {
					continue
				}
				if !rule.ProductTmpl().IsEmpty() && !product.ProductTmpl().Equals(rule.ProductTmpl()) {
					continue
				}
				if !rule.Product().IsEmpty() && !product.Equals(rule.Product()) {
					continue
				}
				if !rule.Categ().IsEmpty() {
					cat := product.Categ()
					for ; !cat.IsEmpty(); cat = cat.Parent() {
						if cat.Equals(rule.Categ()) {
							break
						}
					}
					if cat.IsEmpty() {
						continue
					}
				}
				if rule.Base() == "pricelist" && !rule.BasePricelist().IsEmpty() {
					priceTmp, _ := rule.BasePricelist().ComputePriceRule(product, quantity, partner, dates.Date{},
						h.ProductUom().NewSet(rs.Env()))
					price = rule.BasePricelist().Currency().Compute(priceTmp, rs.Currency(), false)
				} else {
					// if base option is public price take sale price else cost price of product
					// price_compute returns the price in the context UoM, i.e. qty_uom_id
					price = product.PriceCompute(models.FieldName(rule.Base()), h.ProductUom().NewSet(rs.Env()),
						h.Currency().NewSet(rs.Env()), h.Company().NewSet(rs.Env()))
				}
				convertToPriceUom := func(p float64) float64 {
					return product.Uom().ComputePrice(p, priceUom)
				}

				if price == 0 {
					break
				}
				switch rule.ComputePrice() {
				case "fixed":
					price = convertToPriceUom(rule.FixedPrice())
				case "percentage":
					price = price - (price * (rule.PercentPrice() / 100))
				case "formula":
					priceLimit := price
					price = price - (price * (rule.PriceDiscount() / 100))
					if rule.PriceRound() != 0 {
						price = nbutils.Round(price, rule.PriceRound())
					}
					if rule.PriceSurcharge() != 0 {
						priceSurcharge := convertToPriceUom(rule.PriceSurcharge())
						price += priceSurcharge
					}
					if rule.PriceMinMargin() != 0 {
						priceMinMargin := convertToPriceUom(rule.PriceMinMargin())
						price = math.Max(price, priceLimit+priceMinMargin)
					}
					if rule.PriceMaxMargin() != 0 {
						priceMaxMargin := convertToPriceUom(rule.PriceMaxMargin())
						price = math.Min(price, priceLimit+priceMaxMargin)
					}
				}
				suitableRule = rule
				break
			}
			// Final price conversion into pricelist currency
			if !suitableRule.IsEmpty() && suitableRule.ComputePrice() != "fixed" && suitableRule.Base() != "pricelist" {
				price = product.Currency().Compute(price, rs.Currency(), false)
			}
			return price, suitableRule
		})

	h.ProductPricelist().Methods().GetProductPrice().DeclareMethod(
		`GetProductPrice returns the price of the given product in the given quantity for the given partner, at
		the given date and in the given UoM according to this price list.`,
		func(rs h.ProductPricelistSet, product h.ProductProductSet, quantity float64, partner h.PartnerSet,
			date dates.Date, uom h.ProductUomSet) float64 {

			rs.EnsureOne()
			price, _ := rs.ComputePriceRule(product, quantity, partner, date, uom)
			return price
		})

	h.ProductPricelist().Methods().GetProductPriceRule().DeclareMethod(
		`GetProductPriceRule returns the applicable price list rule for the given product in the given quantity
		for the given partner, at the given date and in the given UoM according to this price list.`,
		func(rs h.ProductPricelistSet, product h.ProductProductSet, quantity float64, partner h.PartnerSet,
			date dates.Date, uom h.ProductUomSet) h.ProductPricelistItemSet {

			rs.EnsureOne()
			_, rule := rs.ComputePriceRule(product, quantity, partner, date, uom)
			return rule
		})

	h.ProductPricelist().Methods().GetPartnerPricelist().DeclareMethod(
		`GetPartnerPricelist retrieve the applicable pricelist for the given partner in the given company.`,
		func(rs h.ProductPricelistSet, partner h.PartnerSet, company h.CompanySet) h.ProductPricelistSet {
			if company.IsEmpty() {
				company = h.User().NewSet(rs.Env()).CurrentUser().Company()
			}
			pl := partner.ProductPricelist()
			if pl.IsEmpty() {
				if !partner.Country().IsEmpty() {
					pl = h.ProductPricelist().Search(rs.Env(),
						q.ProductPricelist().CountryGroupsFilteredOn(
							q.CountryGroup().CountriesFilteredOn(
								q.Country().Code().Equals(partner.Country().Code())))).Limit(1)
				}
			}
			if pl.IsEmpty() {
				pl = h.ProductPricelist().Search(rs.Env(),
					q.ProductPricelist().CountryGroups().IsNull()).Limit(1)
			}
			if pl.IsEmpty() {
				pl = company.DefaultPriceList()
			}
			if pl.IsEmpty() {
				pl = h.ProductPricelist().NewSet(rs.Env()).SearchAll().Limit(1)
			}
			return pl
		})

	h.CountryGroup().AddFields(map[string]models.FieldDefinition{
		"Pricelists": models.Many2ManyField{String: "Pricelists", RelationModel: h.ProductPricelist(),
			JSON: "pricelist_ids"},
	})

	h.ProductPricelistItem().DeclareModel()
	h.ProductPricelistItem().SetDefaultOrder("AppliedOn", "MinQuantity DESC", "Categ DESC", "ID")

	h.ProductPricelistItem().AddFields(map[string]models.FieldDefinition{
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: h.ProductTemplate(),
			OnDelete: models.Cascade,
			Help:     "Specify a template if this rule only applies to one product template. Keep empty otherwise."},
		"Product": models.Many2OneField{RelationModel: h.ProductProduct(), OnDelete: models.Cascade,
			Help: "Specify a product if this rule only applies to one product. Keep empty otherwise."},
		"Categ": models.Many2OneField{String: "Product Category", RelationModel: h.ProductCategory(),
			OnDelete: models.Cascade,
			Help: `Specify a product category if this rule only applies to products belonging to this category or 
its children categories. Keep empty otherwise.`},
		"MinQuantity": models.FloatField{Default: models.DefaultValue(1),
			Help: `For the rule to apply, bought/sold quantity must be greater
than or equal to the minimum quantity specified in this field.
Expressed in the default unit of measure of the product.`},
		"AppliedOn": models.SelectionField{String: "Apply On", Selection: types.Selection{
			"3_global":           "Global",
			"2_product_category": "Product Category",
			"1_product":          "Product",
			"0_product_variant":  "Product Variant",
		}, Default: models.DefaultValue("3_global"), Required: true,
			Help:     "Pricelist Item applicable on selected option",
			OnChange: h.ProductPricelistItem().Methods().OnchangeAppliedOn()},
		"Sequence": models.IntegerField{Default: models.DefaultValue(5), Required: true,
			Help: `Gives the order in which the pricelist items will be checked. The evaluation gives highest priority
to lowest sequence and stops as soon as a matching item is found.`},
		"Base": models.SelectionField{String: "Based on", Selection: types.Selection{
			"ListPrice":     "Public Price",
			"StandardPrice": "Cost",
			"pricelist":     "Other Pricelist",
		}, Default: models.DefaultValue("list_price"), Required: true,
			Help: `Base price for computation.
- Public Price: The base price will be the Sale/public Price.
- Cost Price : The base price will be the cost price.
- Other Pricelist : Computation of the base price based on another Pricelist.`,
			Constraint: h.ProductPricelistItem().Methods().CheckOtherList()},
		"BasePricelist": models.Many2OneField{String: "Other Pricelist", RelationModel: h.ProductPricelist(),
			Constraint: h.ProductPricelistItem().Methods().CheckOtherList()},
		"Pricelist": models.Many2OneField{RelationModel: h.ProductPricelist(), Index: true,
			OnDelete: models.Cascade, Constraint: h.ProductPricelistItem().Methods().CheckOtherList()},
		"PriceSurcharge": models.FloatField{Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: "Specify the fixed amount to add or substract(if negative) to the amount calculated with the discount."},
		"PriceDiscount": models.FloatField{Default: models.DefaultValue(0),
			Digits: nbutils.Digits{Precision: 16, Scale: 2}},
		"PriceRound": models.FloatField{Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: `Sets the price so that it is a multiple of this value.
Rounding is applied after the discount and before the surcharge.
To have prices that end in 9.99, set rounding 10, surcharge -0.01`},
		"PriceMinMargin": models.FloatField{String: "Min. Price Margin",
			Digits:     decimalPrecision.GetPrecision("Product Price"),
			Help:       "Specify the minimum amount of margin over the base price.",
			Constraint: h.ProductPricelistItem().Methods().CheckMargin()},
		"PriceMaxMargin": models.FloatField{String: "Max. Price Margin",
			Digits:     decimalPrecision.GetPrecision("Product Price"),
			Help:       "Specify the maximum amount of margin over the base price.",
			Constraint: h.ProductPricelistItem().Methods().CheckMargin()},
		"Company": models.Many2OneField{RelationModel: h.Company(), ReadOnly: true,
			Related: "Pricelist.Company"},
		"Currency": models.Many2OneField{RelationModel: h.Currency(), ReadOnly: true,
			Related: "Pricelist.Currency"},
		"DateStart": models.DateField{String: "Start Date", Help: "Starting date for the pricelist item validation"},
		"DateEnd":   models.DateField{String: "End Date", Help: "Ending valid for the pricelist item validation"},
		"ComputePrice": models.SelectionField{Selection: types.Selection{
			"fixed":      "Fix Price",
			"percentage": "Percentage (discount)",
			"formula":    "Formula",
		},
			Index: true, Default: models.DefaultValue("fixed"),
			OnChange: h.ProductPricelistItem().Methods().OnchangeComputePrice()},
		"FixedPrice":   models.FloatField{String: "Fixed Price", Digits: decimalPrecision.GetPrecision("Product Price")},
		"PercentPrice": models.FloatField{String: "Percentage Price"},
		"Name": models.CharField{Compute: h.ProductPricelistItem().Methods().GetPricelistItemNamePrice(),
			Help: "Explicit rule name for this pricelist line."},
		"Price": models.CharField{Compute: h.ProductPricelistItem().Methods().GetPricelistItemNamePrice(),
			Help: "Explicit rule name for this pricelist line."},
	})

	h.ProductPricelistItem().Methods().CheckOtherList().DeclareMethod(
		`CheckOtherList panics if the other list used in a rule is the same as the base list`,
		func(rs h.ProductPricelistItemSet) {
			for _, item := range rs.Records() {
				if item.Base() == "pricelist" && !item.Pricelist().IsEmpty() && item.Pricelist().Equals(item.BasePricelist()) {
					log.Panic(rs.T("Error! You cannot assign the Main Pricelist as Other Pricelist in PriceList Item!"))
				}
			}
		})

	h.ProductPricelistItem().Methods().CheckMargin().DeclareMethod(
		`CheckMargin checks that the max margin is greater or equal to the min margin`,
		func(rs h.ProductPricelistItemSet) {
			for _, item := range rs.Records() {
				if item.PriceMinMargin() > item.PriceMaxMargin() {
					log.Panic(rs.T("Error! The minimum margin should be lower than the maximum margin."))
				}
			}
		})

	h.ProductPricelistItem().Methods().GetPricelistItemNamePrice().DeclareMethod(
		`GetPricelistItemNamePrice computes the name and the price fields of this line`,
		func(rs h.ProductPricelistItemSet) *h.ProductPricelistItemData {
			var name, price string
			switch {
			case !rs.Categ().IsEmpty():
				name = rs.T("Category: %s", rs.Categ().Name())
			case !rs.ProductTmpl().IsEmpty():
				name = rs.ProductTmpl().Name()
			case !rs.Product().IsEmpty():
				name = strings.Replace(rs.Product().DisplayName(),
					fmt.Sprintf("[%s]", rs.Product().DefaultCode()), "", 1)
			default:
				name = rs.T("All Products")
			}
			switch {
			case rs.ComputePrice() == "fixed":
				price = fmt.Sprintf("%v %v", rs.FixedPrice(), rs.Pricelist().Currency().Name())
			case rs.ComputePrice() == "percentage":
				price = rs.T("%v %% discount", rs.PercentPrice())
			default:
				price = rs.T("%v %% discount and %v surcharge", math.Abs(rs.PriceDiscount()), rs.PriceSurcharge())
			}
			return &h.ProductPricelistItemData{
				Price: price,
				Name:  name,
			}
		})

	h.ProductPricelistItem().Methods().OnchangeAppliedOn().DeclareMethod(
		`OnchangeAppliedOn updates values when the AppliedOn is changed`,
		func(rs h.ProductPricelistItemSet) (*h.ProductPricelistItemData, []models.FieldNamer) {
			var fieldsToReset []models.FieldNamer
			if rs.AppliedOn() != "0_product_variant" {
				fieldsToReset = append(fieldsToReset, h.ProductPricelistItem().Product())
			}
			if rs.AppliedOn() != "1_product" {
				fieldsToReset = append(fieldsToReset, h.ProductPricelistItem().ProductTmpl())
			}
			if rs.AppliedOn() != "2_product_category" {
				fieldsToReset = append(fieldsToReset, h.ProductPricelistItem().Categ())
			}
			return new(h.ProductPricelistItemData), fieldsToReset
		})

	h.ProductPricelistItem().Methods().OnchangeComputePrice().DeclareMethod(
		`OnchangeComputePrice updates values when the ComputePrice field is changed`,
		func(rs h.ProductPricelistItemSet) (*h.ProductPricelistItemData, []models.FieldNamer) {
			var fieldsToReset []models.FieldNamer
			if rs.ComputePrice() != "fixed" {
				fieldsToReset = append(fieldsToReset, h.ProductPricelistItem().FixedPrice())
			}
			if rs.ComputePrice() != "percentage" {
				fieldsToReset = append(fieldsToReset, h.ProductPricelistItem().PercentPrice())
			}
			if rs.ComputePrice() != "formula" {
				fieldsToReset = append(fieldsToReset,
					h.ProductPricelistItem().PriceDiscount(),
					h.ProductPricelistItem().PriceSurcharge(),
					h.ProductPricelistItem().PriceRound(),
					h.ProductPricelistItem().PriceMinMargin(),
					h.ProductPricelistItem().PriceMaxMargin(),
				)
			}
			return new(h.ProductPricelistItemData), fieldsToReset
		})

}
