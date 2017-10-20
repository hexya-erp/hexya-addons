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
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProductTemplate().DeclareModel()
	pool.ProductTemplate().SetDefaultOrder("Name")

	pool.ProductTemplate().AddFields(map[string]models.FieldDefinition{
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
		"Categ": models.Many2OneField{String: "Internal Category", RelationModel: pool.ProductCategory(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				if env.Context().HasKey("categ_id") {
					return pool.ProductCategory().Browse(env, []int64{env.Context().GetInteger("categ_id")})
				}
				if env.Context().HasKey("default_categ_id") {
					return pool.ProductCategory().Browse(env, []int64{env.Context().GetInteger("default_categ_id")})
				}
				category := pool.ProductCategory().Search(env, pool.ProductCategory().HexyaExternalID().Equals("product_product_category_all"))
				if category.Type() != "normal" {
					return pool.ProductCategory().NewSet(env)
				}
				return category
			}, Filter: pool.ProductCategory().Type().Equals("normal"), Required: true,
			Help: "Select category for the current product"},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(),
			Compute: pool.ProductTemplate().Methods().ComputeCurrency()},
		"Price": models.FloatField{Compute: pool.ProductTemplate().Methods().ComputeTemplatePrice(),
			Inverse: pool.ProductTemplate().Methods().InverseTemplatePrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price")},
		"ListPrice": models.FloatField{String: "Sale Price", Default: models.DefaultValue(1.0),
			Digits: decimalPrecision.GetPrecision("Product Price"),
			Help:   "Base price to compute the customer price. Sometimes called the catalog price."},
		"LstPrice": models.FloatField{String: "Public Price", Related: "ListPrice",
			Digits: decimalPrecision.GetPrecision("Product Price")},
		"StandardPrice": models.FloatField{String: "Cost",
			Compute: pool.ProductTemplate().Methods().ComputeStandardPrice(),
			Depends: []string{"ProductVariants", "ProductVariants.StandardPrice"},
			Inverse: pool.ProductTemplate().Methods().InverseStandardPrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Help:    "Cost of the product, in the default unit of measure of the product."},
		"Volume": models.FloatField{Compute: pool.ProductTemplate().Methods().ComputeVolume(),
			Depends: []string{"ProductVariants", "ProductVariants.Volume"},
			Inverse: pool.ProductTemplate().Methods().InverseVolume(), Help: "The volume in m3.", Stored: true},
		"Weight": models.FloatField{Compute: pool.ProductTemplate().Methods().ComputeWeight(),
			Depends: []string{"ProductVariants", "ProductVariants.Weight"},
			Inverse: pool.ProductTemplate().Methods().InverseWeight(),
			Digits:  decimalPrecision.GetPrecision("Stock Weight"), Stored: true,
			Help: "The weight of the contents in Kg, not including any packaging, etc."},
		"Warranty": models.FloatField{},
		"SaleOk": models.BooleanField{String: "Can be Sold", Default: models.DefaultValue(true),
			Help: "Specify if the product can be selected in a sales order line."},
		"PurchaseOk": models.BooleanField{String: "Can be Purchased", Default: models.DefaultValue(true)},
		"Pricelist": models.Many2OneField{String: "Pricelist", RelationModel: pool.ProductPricelist(),
			Stored: false, Help: "Technical field. Used for searching on pricelists, not stored in database."},
		"Uom": models.Many2OneField{String: "Unit of Measure", RelationModel: pool.ProductUom(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.ProductUom().NewSet(env).SearchAll().Limit(1).OrderBy("ID")
			}, Required: true, Help: "Default Unit of Measure used for all stock operation.",
			Constraint: pool.ProductTemplate().Methods().CheckUom(),
			OnChange:   pool.ProductTemplate().Methods().OnchangeUom()},
		"UomPo": models.Many2OneField{String: "Purchase Unit of Measure", RelationModel: pool.ProductUom(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.ProductUom().NewSet(env).SearchAll().Limit(1).OrderBy("ID")
			}, Required: true, Constraint: pool.ProductTemplate().Methods().CheckUom(),
			Help: "Default Unit of Measure used for purchase orders. It must be in the same category than the default unit of measure."},
		"Company": models.Many2OneField{String: "Company", RelationModel: pool.Company(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.ProductUom().NewSet(env).SearchAll().Limit(1).OrderBy("ID")
			}, Index: true},
		"Packagings": models.One2ManyField{String: "Logistical Units", RelationModel: pool.ProductPackaging(),
			ReverseFK: "ProductTmpl", JSON: "packaging_ids",
			Help: `Gives the different ways to package the same product. This has no impact on
the picking order and is mainly used if you use the EDI module.`},
		"Sellers": models.One2ManyField{String: "Vendors", RelationModel: pool.ProductSupplierinfo(),
			ReverseFK: "ProductTmpl", JSON: "seller_ids"},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If unchecked, it will allow you to hide the product without removing it."},
		"Color": models.IntegerField{String: "Color Index"},
		"AttributeLines": models.One2ManyField{String: "Product Attributes",
			RelationModel: pool.ProductAttributeLine(), ReverseFK: "ProductTmpl", JSON: "attribute_line_ids"},
		"ProductVariants": models.One2ManyField{String: "Products", RelationModel: pool.ProductProduct(),
			ReverseFK: "ProductTmpl", JSON: "product_variant_ids", Required: true},
		"ProductVariant": models.Many2OneField{String: "Product", RelationModel: pool.ProductProduct(),
			Compute: pool.ProductTemplate().Methods().ComputeProductVariant(),
			Depends: []string{"ProductVariants"}},
		"ProductVariantCount": models.IntegerField{String: "# Product Variants",
			Compute: pool.ProductTemplate().Methods().ComputeProductVariantCount(),
			Depends: []string{"ProductVariants"}, GoType: new(int)},
		"Barcode": models.CharField{},
		"DefaultCode": models.CharField{String: "Internal Reference",
			Compute: pool.ProductTemplate().Methods().ComputeDefaultCode(),
			Depends: []string{"ProductVariants", "ProductVariants.DefaultCode"},
			Inverse: pool.ProductTemplate().Methods().InverseDefaultCode(), Stored: true},
		"Items": models.One2ManyField{String: "Pricelist Items", RelationModel: pool.ProductPricelistItem(),
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

	pool.ProductTemplate().Fields().StandardPrice().RevokeAccess(security.GroupEveryone, security.All)
	pool.ProductTemplate().Fields().StandardPrice().GrantAccess(base.GroupUser, security.All)

	pool.ProductTemplate().Methods().ComputeProductVariant().DeclareMethod(
		`ComputeProductVariant returns the first variant of this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			return &pool.ProductTemplateData{
				ProductVariant: rs.ProductVariants().Records()[0],
			}, []models.FieldNamer{pool.ProductTemplate().ProductVariant()}
		})

	pool.ProductTemplate().Methods().ComputeCurrency().DeclareMethod(
		`ComputeCurrency computes the currency of this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			mainCompany := pool.Company().NewSet(rs.Env()).Sudo().Search(
				pool.Company().HexyaExternalID().Equals("base_main_company"))
			if mainCompany.IsEmpty() {
				mainCompany = pool.Company().NewSet(rs.Env()).Sudo().SearchAll().Limit(1).OrderBy("ID")
			}
			currency := mainCompany.Currency()
			if !rs.Company().Sudo().Currency().IsEmpty() {
				currency = rs.Company().Sudo().Currency()
			}
			return &pool.ProductTemplateData{
				Currency: currency,
			}, []models.FieldNamer{pool.ProductTemplate().Currency()}
		})

	pool.ProductTemplate().Methods().ComputeTemplatePrice().DeclareMethod(
		`ComputeTemplatePrice returns the price of this template depending on the context:

		- 'partner' => int64 (id of the partner)
		- 'pricelist' => int64 (id of the price list)
		- 'quantity' => float64`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			if !rs.Env().Context().HasKey("pricelist") {
				return new(pool.ProductTemplateData), []models.FieldNamer{pool.ProductProduct().Price()}
			}
			priceListID := rs.Env().Context().GetInteger("pricelist")
			priceList := pool.ProductPricelist().Browse(rs.Env(), []int64{priceListID})
			if priceList.IsEmpty() {
				return new(pool.ProductTemplateData), []models.FieldNamer{pool.ProductProduct().Price()}
			}
			partnerID := rs.Env().Context().GetInteger("partner")
			partner := pool.Partner().Browse(rs.Env(), []int64{partnerID})
			quantity := rs.Env().Context().GetFloat("quantity")
			if quantity == 0 {
				quantity = 1
			}
			return &pool.ProductTemplateData{
				Price: priceList.GetProductPrice(rs.ProductVariant(), quantity, partner, dates.Today(), pool.ProductUom().NewSet(rs.Env())),
			}, []models.FieldNamer{pool.ProductProduct().Price()}
		})

	pool.ProductTemplate().Methods().InverseTemplatePrice().DeclareMethod(
		`InverseTemplatePrice sets the template's price`,
		func(rs pool.ProductTemplateSet, price float64) {
			if rs.Env().Context().HasKey("uom") {
				uom := pool.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
				value := uom.ComputePrice(price, rs.Uom())
				rs.SetListPrice(value)
				return
			}
			rs.SetListPrice(price)
		})

	pool.ProductTemplate().Methods().ComputeStandardPrice().DeclareMethod(
		`ComputeStandardPrice returns the standard price for this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			if rs.ProductVariants().Len() == 1 {
				return &pool.ProductTemplateData{
					StandardPrice: rs.ProductVariant().StandardPrice(),
				}, []models.FieldNamer{pool.ProductTemplate().StandardPrice()}
			}
			return new(pool.ProductTemplateData), []models.FieldNamer{pool.ProductTemplate().StandardPrice()}
		})

	pool.ProductTemplate().Methods().InverseStandardPrice().DeclareMethod(
		`InverseStandardPrice sets this template's standard price`,
		func(rs pool.ProductTemplateSet, price float64) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetStandardPrice(price)
			}
		})

	pool.ProductTemplate().Methods().ComputeVolume().DeclareMethod(
		`ComputeVolume compute the volume of this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			if rs.ProductVariants().Len() == 1 {
				return &pool.ProductTemplateData{
					Volume: rs.ProductVariant().Volume(),
				}, []models.FieldNamer{pool.ProductTemplate().Volume()}
			}
			return new(pool.ProductTemplateData), []models.FieldNamer{pool.ProductTemplate().Volume()}
		})

	pool.ProductTemplate().Methods().InverseVolume().DeclareMethod(
		`InverseVolume sets this template's volume`,
		func(rs pool.ProductTemplateSet, volume float64) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetVolume(volume)
			}
		})

	pool.ProductTemplate().Methods().ComputeWeight().DeclareMethod(
		`ComputeWeight compute the weight of this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			if rs.ProductVariants().Len() == 1 {
				return &pool.ProductTemplateData{
					Weight: rs.ProductVariant().Weight(),
				}, []models.FieldNamer{pool.ProductTemplate().Weight()}
			}
			return new(pool.ProductTemplateData), []models.FieldNamer{pool.ProductTemplate().Weight()}
		})

	pool.ProductTemplate().Methods().InverseWeight().DeclareMethod(
		`InverseWeightsets this template's weight`,
		func(rs pool.ProductTemplateSet, weight float64) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetWeight(weight)
			}
		})

	pool.ProductTemplate().Methods().ComputeProductVariantCount().DeclareMethod(
		`ComputeProductVariantCount returns the number of variants for this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			return &pool.ProductTemplateData{
				ProductVariantCount: rs.ProductVariants().Len(),
			}, []models.FieldNamer{pool.ProductTemplate().ProductVariantCount()}
		})

	pool.ProductTemplate().Methods().ComputeDefaultCode().DeclareMethod(
		`ComputeDefaultCode returns the default code for this template`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			if rs.ProductVariants().Len() == 1 {
				return &pool.ProductTemplateData{
					DefaultCode: rs.ProductVariant().DefaultCode(),
				}, []models.FieldNamer{pool.ProductTemplate().DefaultCode()}
			}
			return new(pool.ProductTemplateData), []models.FieldNamer{pool.ProductTemplate().DefaultCode()}
		})

	pool.ProductTemplate().Methods().InverseDefaultCode().DeclareMethod(
		`InverseDefaultCode sets the default code of this template`,
		func(rs pool.ProductTemplateSet, code string) {
			if rs.ProductVariants().Len() == 1 {
				rs.ProductVariant().SetDefaultCode(code)
			}
		})

	pool.ProductTemplate().Methods().CheckUom().DeclareMethod(
		`CheckUom checks that this template's uom is of the same category as the purchase uom`,
		func(rs pool.ProductTemplateSet) {
			if !rs.Uom().IsEmpty() && !rs.UomPo().IsEmpty() && !rs.Uom().Category().Equals(rs.UomPo().Category()) {
				log.Panic(rs.T("Error: The default Unit of Measure and the purchase Unit of Measure must be in the same category."))
			}
		})

	pool.ProductTemplate().Methods().OnchangeUom().DeclareMethod(
		`OnchangeUom updates UomPo when uom is changed`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateData, []models.FieldNamer) {
			if !rs.Uom().IsEmpty() {
				return &pool.ProductTemplateData{
					UomPo: rs.Uom(),
				}, []models.FieldNamer{pool.ProductTemplate().UomPo()}
			}
			return new(pool.ProductTemplateData), []models.FieldNamer{}
		})

	pool.ProductTemplate().Methods().Create().Extend("",
		func(rs pool.ProductTemplateSet, data *pool.ProductTemplateData) pool.ProductTemplateSet {
			// tools.image_resize_images(vals)
			template := rs.Super().Create(data)
			if !rs.Env().Context().HasKey("create_product_product") {
				template.WithContext("create_from_tmpl", true).CreateVariants()
			}
			// This is needed to set given values to first variant after creation
			relatedVals := &pool.ProductTemplateData{
				Barcode:       data.Barcode,
				DefaultCode:   data.DefaultCode,
				StandardPrice: data.StandardPrice,
				Volume:        data.Volume,
				Weight:        data.Weight,
			}
			template.Write(relatedVals)
			return template
		})

	pool.ProductTemplate().Methods().Write().Extend("",
		func(rs pool.ProductTemplateSet, vals *pool.ProductTemplateData, fieldsToUnset ...models.FieldNamer) bool {
			// tools.image_resize_images(vals)
			res := rs.Super().Write(vals, fieldsToUnset...)
			if _, exists := vals.Get(pool.ProductTemplate().AttributeLines(), fieldsToUnset...); exists || vals.Active {
				rs.CreateVariants()
			}
			if active, exists := vals.Get(pool.ProductTemplate().AttributeLines(), fieldsToUnset...); exists && !active.(bool) {
				rs.WithContext("active_test", false).ProductVariants().SetActive(vals.Active)
			}
			return res
		})

	pool.ProductTemplate().Methods().Copy().Extend("",
		func(rs pool.ProductTemplateSet, overrides *pool.ProductTemplateData, fieldsToUnset ...models.FieldNamer) pool.ProductTemplateSet {
			rs.EnsureOne()
			if _, exists := overrides.Get(pool.ProductTemplate().Name(), fieldsToUnset...); !exists {
				overrides.Name = rs.T("%s (Copy)", rs.Name())
			}
			return rs.Super().Copy(overrides, fieldsToUnset...)
		})

	pool.ProductTemplate().Methods().NameGet().Extend("",
		func(rs pool.ProductTemplateSet) string {
			return pool.ProductProduct().NewSet(rs.Env()).NameFormat(rs.Name(), rs.DefaultCode())
		})

	pool.ProductTemplate().Methods().SearchByName().Extend("",
		func(rs pool.ProductTemplateSet, name string, op operator.Operator, additionalCond pool.ProductTemplateCondition, limit int) pool.ProductTemplateSet {
			// Only use the product.product heuristics if there is a search term and the domain
			// does not specify a match on `product.template` IDs.
			if name == "" {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			for _, term := range additionalCond.Fields() {
				if pool.ProductTemplate().JSONizeFieldName(term) == "id" {
					return rs.Super().SearchByName(name, op, additionalCond, limit)
				}
			}

			templates := pool.ProductTemplate().NewSet(rs.Env())
			if limit == 0 {
				limit = 100
			}
			for templates.Len() > limit {
				var prodCond pool.ProductProductCondition
				if !templates.IsEmpty() {
					prodCond = pool.ProductProduct().ProductTmpl().In(templates)
				}
				products := pool.ProductProduct().NewSet(rs.Env()).SearchByName(name, op,
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

	pool.ProductTemplate().Methods().PriceCompute().DeclareMethod(
		`PriceCompute returns the price field defined by priceType in the given uom and currency
		for the given company.`,
		func(rs pool.ProductTemplateSet, priceType models.FieldNamer, uom pool.ProductUomSet, currency pool.CurrencySet, company pool.CompanySet) {
			rs.EnsureOne()
			template := rs
			if priceType == pool.ProductTemplate().StandardPrice() {
				// StandardPrice field can only be seen by users in base.group_user
				// Thus, in order to compute the sale price from the cost for users not in this group
				// We fetch the standard price as the superuser
				if company.IsEmpty() {
					company = pool.User().NewSet(rs.Env()).CurrentUser().Company()
					if rs.Env().Context().HasKey("force_company") {
						company = pool.Company().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("force_company")})
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

	pool.ProductTemplate().Methods().CreateVariants().DeclareMethod(
		`CreateVariants`,
		func(rs pool.ProductTemplateSet) {
			//@api.multi
			for _, tmpl := range rs.WithContext("active_test", false).Records() {
				// adding an attribute with only one value should not recreate product
				// write this attribute on every product to make sure we don't lose them
				variantAlone := pool.ProductAttributeValue().NewSet(rs.Env())
				for _, attrLine := range tmpl.AttributeLines().Records() {
					if attrLine.Values().Len() != 1 {
						continue
					}
					variantAlone = variantAlone.Union(attrLine.Values())
				}
				for _, value := range variantAlone.Records() {
					for _, prod := range tmpl.ProductVariants().Records() {
						valuesAttributes := pool.ProductAttribute().NewSet(rs.Env())
						for _, val := range prod.AttributeValues().Records() {
							valuesAttributes = valuesAttributes.Union(val.Attribute())
						}
						if value.Attribute().Intersect(valuesAttributes).IsEmpty() {
							prod.SetAttributeValues(prod.AttributeValues().Union(value))
						}
					}
				}

				// list of values combination
				var existingVariants []pool.ProductAttributeValueSet
				for _, prod := range tmpl.ProductVariants().Records() {
					prodVariant := pool.ProductAttributeValue().NewSet(rs.Env())
					for _, attrVal := range prod.AttributeValues().Records() {
						if attrVal.Attribute().CreateVariant() {
							prodVariant = prodVariant.Union(attrVal)
						}
					}
					existingVariants = append(existingVariants, prodVariant)
				}
				var matrixValues []pool.ProductAttributeValueSet
				for _, attrLine := range tmpl.AttributeLines().Records() {
					if !attrLine.Attribute().CreateVariant() {
						continue
					}
					matrixValues = append(matrixValues, attrLine.Values())
				}
				var variantMatrix []pool.ProductAttributeValueSet
				if len(matrixValues) > 0 {
					variantMatrix = variantMatrix[0].CartesianProduct(variantMatrix[1:]...)
				}

				var toCreateVariants []pool.ProductAttributeValueSet
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
				variantsToActivate := pool.ProductProduct().NewSet(rs.Env())
				variantsToUnlink := pool.ProductProduct().NewSet(rs.Env())
				for _, product := range tmpl.ProductVariants().Records() {
					tcAttrs := pool.ProductAttributeValue().NewSet(rs.Env())
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
					pool.ProductProduct().Create(rs.Env(), &pool.ProductProductData{
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
