// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"log"

	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.ProductTemplate().DeclareModel()
	h.ProductTemplate().SetDefaultOrder("Name")

	h.ProductTemplate().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Index: true, Required: true, Translate: true},
		"Sequence": models.IntegerField{Default: models.DefaultValue(1),
			Help: "Gives the sequence order when displaying a product list"},
		"Description": models.TextField{Translate: true,
			Help: "A precise description of the Product, used only for internal information purposes."},
		"DescriptionPurchase": models.TextField{String: "Purchase Description", Translate: true,
			Help: `A description of the Product that you want to communicate to your vendors.
This description will be copied to every Purchase Order, Receipt and Vendor Bill/Refund.`},
		"DescriptionSale": models.TextField{String: "Sale Description", Translate: true,
			Help: `A description of the Product that you want to communicate to your customers.
This description will be copied to every Sale Order, Delivery Order and Customer Invoice/Refund`},
		"Type": models.SelectionField{String: "Product Type", Selection: types.Selection{
			"consu":   "Consumable",
			"service": "Service",
		}, Default: models.DefaultValue("consu"), Required: true,
			Help: `A stockable product is a product for which you manage stock. The "Inventory" app has to be installed.
- A consumable product on the other hand is a product for which stock is not managed.
- A service is a non-material product you provide.
- A digital content is a non-material product you sell online.
	The files attached to the products are the one that are sold on
	the e-commerce such as e-books, music, pictures,...
	The "Digital Product" module has to be installed.`},
		"Rental": models.BooleanField{String: "Can be Rent"},
		"Categ": models.Many2OneField{String: "Internal Category", RelationModel: h.ProductCategory(),
			Default: func(env models.Environment) interface{} {
				if env.Context().HasKey("categ_id") {
					return h.ProductCategory().Browse(env, []int64{env.Context().GetInteger("categ_id")})
				}
				if env.Context().HasKey("default_categ_id") {
					return h.ProductCategory().Browse(env, []int64{env.Context().GetInteger("default_categ_id")})
				}
				category := h.ProductCategory().Search(env, q.ProductCategory().HexyaExternalID().Equals("product_product_category_all"))
				if category.Type() != "normal" {
					return h.ProductCategory().NewSet(env)
				}
				return category
			}, Filter: q.ProductCategory().Type().Equals("normal"), Required: true,
			Help: "Select category for the current product"},
		"Currency": models.Many2OneField{RelationModel: h.Currency(),
			Compute: h.ProductTemplate().Methods().ComputeCurrency()},
		"Price": models.FloatField{Compute: h.ProductTemplate().Methods().ComputeTemplatePrice(),
			Inverse: h.ProductTemplate().Methods().InverseTemplatePrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price")},
		"ListPrice": models.FloatField{String: "Sale Price", Default: models.DefaultValue(1.0),
			Digits: decimalPrecision.GetPrecision("Product Price"),
			Help:   "Base price to compute the customer price. Sometimes called the catalog price."},
		"LstPrice": models.FloatField{String: "Public Price", Related: "ListPrice",
			Digits: decimalPrecision.GetPrecision("Product Price")},
		"StandardPrice": models.FloatField{String: "Cost",
			Compute: h.ProductTemplate().Methods().ComputeStandardPrice(),
			Depends: []string{"ProductVariants", "ProductVariants.StandardPrice"},
			Inverse: h.ProductTemplate().Methods().InverseStandardPrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Help:    "Cost of the product, in the default unit of measure of the product."},
		"Volume": models.FloatField{Compute: h.ProductTemplate().Methods().ComputeVolume(),
			Depends: []string{"ProductVariants", "ProductVariants.Volume"},
			Inverse: h.ProductTemplate().Methods().InverseVolume(), Help: "The volume in m3.", Stored: true},
		"Weight": models.FloatField{Compute: h.ProductTemplate().Methods().ComputeWeight(),
			Depends: []string{"ProductVariants", "ProductVariants.Weight"},
			Inverse: h.ProductTemplate().Methods().InverseWeight(),
			Digits:  decimalPrecision.GetPrecision("Stock Weight"), Stored: true,
			Help: "The weight of the contents in Kg, not including any packaging, etc."},
		"Warranty": models.FloatField{},
		"SaleOk": models.BooleanField{String: "Can be Sold", Default: models.DefaultValue(true),
			Help: "Specify if the product can be selected in a sales order line."},
		"PurchaseOk": models.BooleanField{String: "Can be Purchased", Default: models.DefaultValue(true)},
		"Pricelist": models.Many2OneField{String: "Pricelist", RelationModel: h.ProductPricelist(),
			Stored: false, Help: "Technical field. Used for searching on pricelists, not stored in database."},
		"Uom": models.Many2OneField{String: "Unit of Measure", RelationModel: h.ProductUom(),
			Default: func(env models.Environment) interface{} {
				return h.ProductUom().NewSet(env).SearchAll().Limit(1).OrderBy("ID")
			}, Required: true, Help: "Default Unit of Measure used for all stock operation.",
			Constraint: h.ProductTemplate().Methods().CheckUom(),
			OnChange:   h.ProductTemplate().Methods().OnchangeUom()},
		"UomPo": models.Many2OneField{String: "Purchase Unit of Measure", RelationModel: h.ProductUom(),
			Default: func(env models.Environment) interface{} {
				return h.ProductUom().NewSet(env).SearchAll().Limit(1).OrderBy("ID")
			}, Required: true, Constraint: h.ProductTemplate().Methods().CheckUom(),
			Help: "Default Unit of Measure used for purchase orders. It must be in the same category than the default unit of measure."},
		"Company": models.Many2OneField{String: "Company", RelationModel: h.Company(),
			Default: func(env models.Environment) interface{} {
				return h.ProductUom().NewSet(env).SearchAll().Limit(1).OrderBy("ID")
			}, Index: true},
		"Packagings": models.One2ManyField{String: "Logistical Units", RelationModel: h.ProductPackaging(),
			ReverseFK: "ProductTmpl", JSON: "packaging_ids",
			Help: `Gives the different ways to package the same product. This has no impact on
the picking order and is mainly used if you use the EDI module.`},
		"Sellers": models.One2ManyField{String: "Vendors", RelationModel: h.ProductSupplierinfo(),
			ReverseFK: "ProductTmpl", JSON: "seller_ids"},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If unchecked, it will allow you to hide the product without removing it."},
		"Color": models.IntegerField{String: "Color Index"},
		"AttributeLines": models.One2ManyField{String: "Product Attributes",
			RelationModel: h.ProductAttributeLine(), ReverseFK: "ProductTmpl", JSON: "attribute_line_ids"},
		"ProductVariants": models.One2ManyField{String: "Products", RelationModel: h.ProductProduct(),
			ReverseFK: "ProductTmpl", JSON: "product_variant_ids", Required: true},
		"ProductVariant": models.Many2OneField{String: "Product", RelationModel: h.ProductProduct(),
			Compute: h.ProductTemplate().Methods().ComputeProductVariant(),
			Depends: []string{"ProductVariants"}},
		"ProductVariantCount": models.IntegerField{String: "# Product Variants",
			Compute: h.ProductTemplate().Methods().ComputeProductVariantCount(),
			Depends: []string{"ProductVariants"}, GoType: new(int)},
		"Barcode": models.CharField{},
		"DefaultCode": models.CharField{String: "Internal Reference",
			Compute: h.ProductTemplate().Methods().ComputeDefaultCode(),
			Depends: []string{"ProductVariants", "ProductVariants.DefaultCode"},
			Inverse: h.ProductTemplate().Methods().InverseDefaultCode(), Stored: true},
		"Items": models.One2ManyField{String: "Pricelist Items", RelationModel: h.ProductPricelistItem(),
			ReverseFK: "ProductTmpl", JSON: "item_ids"},
		"Image": models.BinaryField{
			Help: "This field holds the image used as image for the product, limited to 1024x1024px."},
		"ImageMedium": models.BinaryField{String: "Medium-sized image",
			Help: `Medium-sized image of the product. It is automatically
resized as a 128x128px image, with aspect ratio preserved,
only when the image exceeds one of those sizes.
Use this field in form views or some kanban views.`},
		"ImageSmall": models.BinaryField{String: "Small-sized image",
			Help: `Small-sized image of the product. It is automatically
resized as a 64x64px image, with aspect ratio preserved.
Use this field anywhere a small image is required.`},
	})

	h.ProductTemplate().Fields().StandardPrice().RevokeAccess(security.GroupEveryone, security.All)
	h.ProductTemplate().Fields().StandardPrice().GrantAccess(base.GroupUser, security.All)

	h.ProductTemplate().Methods().ComputeProductVariant().DeclareMethod(
		`ComputeProductVariant returns the first variant of this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			return &h.ProductTemplateData{
				ProductVariant: rs.ProductVariants().Records()[0],
			}
		})

	h.ProductTemplate().Methods().ComputeCurrency().DeclareMethod(
		`ComputeCurrency computes the currency of this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			mainCompany := h.Company().NewSet(rs.Env()).Sudo().Search(
				q.Company().HexyaExternalID().Equals("base_main_company"))
			if mainCompany.IsEmpty() {
				mainCompany = h.Company().NewSet(rs.Env()).Sudo().SearchAll().Limit(1).OrderBy("ID")
			}
			currency := mainCompany.Currency()
			if !rs.Company().Sudo().Currency().IsEmpty() {
				currency = rs.Company().Sudo().Currency()
			}
			return &h.ProductTemplateData{
				Currency: currency,
			}
		})

	h.ProductTemplate().Methods().ComputeTemplatePrice().DeclareMethod(
		`ComputeTemplatePrice returns the price of this template depending on the context:

		- 'partner' => int64 (id of the partner)
		- 'pricelist' => int64 (id of the price list)
		- 'quantity' => float64`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			if !rs.Env().Context().HasKey("pricelist") {
				return new(h.ProductTemplateData)
			}
			priceListID := rs.Env().Context().GetInteger("pricelist")
			priceList := h.ProductPricelist().Browse(rs.Env(), []int64{priceListID})
			if priceList.IsEmpty() {
				return new(h.ProductTemplateData)
			}
			partnerID := rs.Env().Context().GetInteger("partner")
			partner := h.Partner().Browse(rs.Env(), []int64{partnerID})
			quantity := rs.Env().Context().GetFloat("quantity")
			if quantity == 0 {
				quantity = 1
			}
			return &h.ProductTemplateData{
				Price: priceList.GetProductPrice(rs.ProductVariant(), quantity, partner, dates.Today(), h.ProductUom().NewSet(rs.Env())),
			}
		})

	h.ProductTemplate().Methods().InverseTemplatePrice().DeclareMethod(
		`InverseTemplatePrice sets the template's price`,
		func(rs h.ProductTemplateSet, price float64) {
			if rs.Env().Context().HasKey("uom") {
				uom := h.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
				value := uom.ComputePrice(price, rs.Uom())
				rs.SetListPrice(value)
				return
			}
			rs.SetListPrice(price)
		})

	h.ProductTemplate().Methods().ComputeStandardPrice().DeclareMethod(
		`ComputeStandardPrice returns the standard price for this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			if rs.ProductVariants().Len() == 1 {
				return &h.ProductTemplateData{
					StandardPrice: rs.ProductVariant().StandardPrice(),
				}
			}
			return new(h.ProductTemplateData)
		})

	h.ProductTemplate().Methods().InverseStandardPrice().DeclareMethod(
		`InverseStandardPrice sets this template's standard price`,
		func(rs h.ProductTemplateSet, price float64) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetStandardPrice(price)
			}
		})

	h.ProductTemplate().Methods().ComputeVolume().DeclareMethod(
		`ComputeVolume compute the volume of this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			if rs.ProductVariants().Len() == 1 {
				return &h.ProductTemplateData{
					Volume: rs.ProductVariant().Volume(),
				}
			}
			return new(h.ProductTemplateData)
		})

	h.ProductTemplate().Methods().InverseVolume().DeclareMethod(
		`InverseVolume sets this template's volume`,
		func(rs h.ProductTemplateSet, volume float64) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetVolume(volume)
			}
		})

	h.ProductTemplate().Methods().ComputeWeight().DeclareMethod(
		`ComputeWeight compute the weight of this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			if rs.ProductVariants().Len() == 1 {
				return &h.ProductTemplateData{
					Weight: rs.ProductVariant().Weight(),
				}
			}
			return new(h.ProductTemplateData)
		})

	h.ProductTemplate().Methods().InverseWeight().DeclareMethod(
		`InverseWeightsets this template's weight`,
		func(rs h.ProductTemplateSet, weight float64) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetWeight(weight)
			}
		})

	h.ProductTemplate().Methods().ComputeProductVariantCount().DeclareMethod(
		`ComputeProductVariantCount returns the number of variants for this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			return &h.ProductTemplateData{
				ProductVariantCount: rs.ProductVariants().Len(),
			}
		})

	h.ProductTemplate().Methods().ComputeDefaultCode().DeclareMethod(
		`ComputeDefaultCode returns the default code for this template`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			if rs.ProductVariants().Len() == 1 {
				return &h.ProductTemplateData{
					DefaultCode: rs.ProductVariant().DefaultCode(),
				}
			}
			return new(h.ProductTemplateData)
		})

	h.ProductTemplate().Methods().InverseDefaultCode().DeclareMethod(
		`InverseDefaultCode sets the default code of this template`,
		func(rs h.ProductTemplateSet, code string) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetDefaultCode(code)
			}
		})

	h.ProductTemplate().Methods().CheckUom().DeclareMethod(
		`CheckUom checks that this template's uom is of the same category as the purchase uom`,
		func(rs h.ProductTemplateSet) {
			if !rs.Uom().IsEmpty() && !rs.UomPo().IsEmpty() && !rs.Uom().Category().Equals(rs.UomPo().Category()) {
				log.Panic(rs.T("Error: The default Unit of Measure and the purchase Unit of Measure must be in the same category."))
			}
		})

	h.ProductTemplate().Methods().OnchangeUom().DeclareMethod(
		`OnchangeUom updates UomPo when uom is changed`,
		func(rs h.ProductTemplateSet) (*h.ProductTemplateData, []models.FieldNamer) {
			if !rs.Uom().IsEmpty() {
				return &h.ProductTemplateData{
					UomPo: rs.Uom(),
				}, []models.FieldNamer{h.ProductTemplate().UomPo()}
			}
			return new(h.ProductTemplateData), []models.FieldNamer{}
		})

	h.ProductTemplate().Methods().Create().Extend("",
		func(rs h.ProductTemplateSet, data *h.ProductTemplateData) h.ProductTemplateSet {
			// tools.image_resize_images(vals)
			template := rs.Super().Create(data)
			if !rs.Env().Context().HasKey("create_product_product") {
				template.WithContext("create_from_tmpl", true).CreateVariants()
			}
			// This is needed to set given values to first variant after creation
			relatedVals := &h.ProductTemplateData{
				Barcode:       data.Barcode,
				DefaultCode:   data.DefaultCode,
				StandardPrice: data.StandardPrice,
				Volume:        data.Volume,
				Weight:        data.Weight,
			}
			template.Write(relatedVals)
			return template
		})

	h.ProductTemplate().Methods().Write().Extend("",
		func(rs h.ProductTemplateSet, vals *h.ProductTemplateData, fieldsToUnset ...models.FieldNamer) bool {
			// tools.image_resize_images(vals)
			res := rs.Super().Write(vals, fieldsToUnset...)
			if _, exists := vals.Get(h.ProductTemplate().AttributeLines(), fieldsToUnset...); exists || vals.Active {
				rs.CreateVariants()
			}
			if active, exists := vals.Get(h.ProductTemplate().Active(), fieldsToUnset...); exists && !active.(bool) {
				rs.WithContext("active_test", false).ProductVariants().SetActive(vals.Active)
			}
			return res
		})

	h.ProductTemplate().Methods().Copy().Extend("",
		func(rs h.ProductTemplateSet, overrides *h.ProductTemplateData, fieldsToUnset ...models.FieldNamer) h.ProductTemplateSet {
			rs.EnsureOne()
			if _, exists := overrides.Get(h.ProductTemplate().Name(), fieldsToUnset...); !exists {
				overrides.Name = rs.T("%s (Copy)", rs.Name())
			}
			return rs.Super().Copy(overrides, fieldsToUnset...)
		})

	h.ProductTemplate().Methods().NameGet().Extend("",
		func(rs h.ProductTemplateSet) string {
			return h.ProductProduct().NewSet(rs.Env()).NameFormat(rs.Name(), rs.DefaultCode())
		})

	h.ProductTemplate().Methods().SearchByName().Extend("",
		func(rs h.ProductTemplateSet, name string, op operator.Operator, additionalCond q.ProductTemplateCondition, limit int) h.ProductTemplateSet {
			// Only use the product.product heuristics if there is a search term and the domain
			// does not specify a match on `product.template` IDs.
			if name == "" {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			if additionalCond.HasField(h.ProductTemplate().Fields().ID()) {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}

			templates := h.ProductTemplate().NewSet(rs.Env())
			if limit == 0 {
				limit = 100
			}
			for templates.Len() > limit {
				var prodCond q.ProductProductCondition
				if !templates.IsEmpty() {
					prodCond = q.ProductProduct().ProductTmpl().In(templates)
				}
				products := h.ProductProduct().NewSet(rs.Env()).SearchByName(name, op,
					prodCond.And().ProductTmplFilteredOn(additionalCond), limit)
				for _, prod := range products.Records() {
					templates = templates.Union(prod.ProductTmpl())
				}
				if products.IsEmpty() {
					break
				}
			}
			return templates
		})

	h.ProductTemplate().Methods().PriceCompute().DeclareMethod(
		`PriceCompute returns the price field defined by priceType in the given uom and currency
		for the given company.`,
		func(rs h.ProductTemplateSet, priceType models.FieldNamer, uom h.ProductUomSet, currency h.CurrencySet, company h.CompanySet) {
			rs.EnsureOne()
			template := rs
			if priceType == h.ProductTemplate().StandardPrice() {
				// StandardPrice field can only be seen by users in base.group_user
				// Thus, in order to compute the sale price from the cost for users not in this group
				// We fetch the standard price as the superuser
				if company.IsEmpty() {
					company = h.User().NewSet(rs.Env()).CurrentUser().Company()
					if rs.Env().Context().HasKey("force_company") {
						company = h.Company().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("force_company")})
					}
				}
				template = rs.WithContext("force_company", company.ID()).Sudo()
			}
			price := template.Get(priceType.String()).(float64)
			if !uom.IsEmpty() {
				price = template.Uom().ComputePrice(price, uom)
			}
			// Convert from current user company currency to asked one
			// This is right cause a field cannot be in more than one currency
			if !currency.IsEmpty() {
				price = template.Currency().Compute(price, currency, true)
			}
		})

	h.ProductTemplate().Methods().CreateVariants().DeclareMethod(
		`CreateVariants`,
		func(rs h.ProductTemplateSet) {
			for _, tmpl := range rs.WithContext("active_test", false).Records() {
				// adding an attribute with only one value should not recreate product
				// write this attribute on every product to make sure we don't lose them
				variantAlone := h.ProductAttributeValue().NewSet(rs.Env())
				for _, attrLine := range tmpl.AttributeLines().Records() {
					if attrLine.Values().Len() == 1 {
						variantAlone = variantAlone.Union(attrLine.Values())
					}
				}
				for _, value := range variantAlone.Records() {
					for _, prod := range tmpl.ProductVariants().Records() {
						valuesAttributes := h.ProductAttribute().NewSet(rs.Env())
						for _, val := range prod.AttributeValues().Records() {
							valuesAttributes = valuesAttributes.Union(val.Attribute())
						}
						if value.Attribute().Intersect(valuesAttributes).IsEmpty() {
							prod.SetAttributeValues(prod.AttributeValues().Union(value))
						}
					}
				}

				// list of values combination
				var existingVariants []h.ProductAttributeValueSet
				for _, prod := range tmpl.ProductVariants().Records() {
					prodVariant := h.ProductAttributeValue().NewSet(rs.Env())
					for _, attrVal := range prod.AttributeValues().Records() {
						if attrVal.Attribute().CreateVariant() {
							prodVariant = prodVariant.Union(attrVal)
						}
					}
					existingVariants = append(existingVariants, prodVariant)
				}
				var matrixValues []h.ProductAttributeValueSet
				for _, attrLine := range tmpl.AttributeLines().Records() {
					if !attrLine.Attribute().CreateVariant() {
						continue
					}
					matrixValues = append(matrixValues, attrLine.Values())
				}
				var variantMatrix []h.ProductAttributeValueSet
				if len(matrixValues) > 0 {
					variantMatrix = matrixValues[0].CartesianProduct(matrixValues[1:]...)
				} else {
					variantMatrix = []h.ProductAttributeValueSet{h.ProductAttributeValue().NewSet(rs.Env())}
				}

				var toCreateVariants []h.ProductAttributeValueSet
				for _, mVariant := range variantMatrix {
					var exists bool
					for _, eVariant := range existingVariants {
						if mVariant.Equals(eVariant) {
							exists = true
							break
						}
					}
					if !exists {
						toCreateVariants = append(toCreateVariants, mVariant)
					}
				}

				// check product
				variantsToActivate := h.ProductProduct().NewSet(rs.Env())
				variantsToUnlink := h.ProductProduct().NewSet(rs.Env())
				for _, product := range tmpl.ProductVariants().Records() {
					tcAttrs := h.ProductAttributeValue().NewSet(rs.Env())
					for _, attrVal := range product.AttributeValues().Records() {
						if !attrVal.Attribute().CreateVariant() {
							continue
						}
						tcAttrs = tcAttrs.Union(attrVal)
					}
					var inMatrix bool
					for _, mVariant := range variantMatrix {
						if tcAttrs.Equals(mVariant) {

							inMatrix = true
							break
						}
					}
					switch {
					case inMatrix && !product.Active():
						variantsToActivate = variantsToActivate.Union(product)
					case !inMatrix:
						variantsToUnlink = variantsToUnlink.Union(product)
					}
				}
				if !variantsToActivate.IsEmpty() {
					variantsToActivate.SetActive(true)
				}

				// create new product
				for _, variants := range toCreateVariants {
					h.ProductProduct().Create(rs.Env(), &h.ProductProductData{
						ProductTmpl:     tmpl,
						AttributeValues: variants,
					})
				}

				// inactive product
				if !variantsToUnlink.IsEmpty() {
					variantsToUnlink.SetActive(false)
				}

			}
		})

}
