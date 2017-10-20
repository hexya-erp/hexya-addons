// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"fmt"

	"log"

	"strings"

	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProductAttribute().DeclareModel()
	pool.ProductAttribute().SetDefaultOrder("Sequence", "Name")

	pool.ProductAttribute().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Required: true, Translate: true},
		"Values": models.One2ManyField{RelationModel: pool.ProductAttributeValue(), ReverseFK: "Attribute",
			JSON: "value_ids", NoCopy: false},
		"Sequence": models.IntegerField{Help: "Determine the display order"},
		"AttributeLines": models.One2ManyField{String: "Lines", RelationModel: pool.ProductAttributeLine(),
			ReverseFK: "Attribute", JSON: "attribute_line_ids"},
		"CreateVariant": models.BooleanField{Default: models.DefaultValue(true),
			Help: "Check this if you want to create multiple variants for this attribute."},
	})

	pool.ProductAttributeValue().DeclareModel()
	pool.ProductAttributeValue().SetDefaultOrder("Sequence")

	pool.ProductAttributeValue().AddFields(map[string]models.FieldDefinition{
		"Name":     models.CharField{String: "Value", Required: true, Translate: true},
		"Sequence": models.IntegerField{Help: "Determine the display order"},
		"Attribute": models.Many2OneField{RelationModel: pool.ProductAttribute(), OnDelete: models.Cascade,
			Required: true},
		"Products": models.Many2ManyField{String: "Variants", RelationModel: pool.ProductProduct(),
			JSON: "product_ids"},
		"PriceExtra": models.FloatField{String: "Attribute Price Extra",
			Compute: pool.ProductAttributeValue().Methods().ComputePriceExtra(),
			Inverse: pool.ProductAttributeValue().Methods().InversePriceExtra(),
			Default: models.DefaultValue(0), Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: "Price Extra: Extra price for the variant with this attribute value on sale price. eg. 200 price extra, 1000 + 200 = 1200."},
		"Prices": models.One2ManyField{String: "Attribute Prices", RelationModel: pool.ProductAttributePrice(),
			ReverseFK: "Value", JSON: "price_ids" /* readonly */},
	})

	pool.ProductAttributeValue().AddSQLConstraint("ValueCompanyUniq", "unique (name,attribute_id)", "This attribute value already exists !")

	pool.ProductAttributeValue().Methods().ComputePriceExtra().DeclareMethod(
		`ComputePriceExtra returns the price extra for this attribute for the product
		template passed as 'active_id' in the context. Returns 0 if there is not 'active_id'.`,
		func(rs pool.ProductAttributeValueSet) (*pool.ProductAttributeValueData, []models.FieldNamer) {
			var priceExtra float64
			if rs.Env().Context().HasKey("active_id") {
				productTmpl := pool.ProductTemplate().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("active_id")})
				price := rs.Prices().Search(pool.ProductAttributePrice().ProductTmpl().Equals(productTmpl))
				priceExtra = price.PriceExtra()
			}
			return &pool.ProductAttributeValueData{
				PriceExtra: priceExtra,
			}, []models.FieldNamer{pool.ProductAttributeValue().PriceExtra()}
		})

	pool.ProductAttributeValue().Methods().InversePriceExtra().DeclareMethod(
		`InversePriceExtra sets the price extra based on the product
		template passed as 'active_id'. Does nothing if there is not 'active_id'.`,
		func(rs pool.ProductAttributeValueSet, value float64) {
			if !rs.Env().Context().HasKey("active_id") {
				return
			}
			productTmpl := pool.ProductTemplate().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("active_id")})
			prices := pool.ProductAttributePrice().Search(rs.Env(),
				pool.ProductAttributePrice().Value().In(rs).And().ProductTmpl().Equals(productTmpl))
			if !prices.IsEmpty() {
				prices.SetPriceExtra(value)
				return
			}
			updated := pool.ProductAttributeValue().NewSet(rs.Env())
			for _, price := range prices.Records() {
				updated = updated.Union(price.Value())
			}
			for _, val := range rs.Subtract(updated).Records() {
				pool.ProductAttributePrice().Create(rs.Env(), &pool.ProductAttributePriceData{
					ProductTmpl: productTmpl,
					Value:       val,
					PriceExtra:  value,
				})
			}
		})

	pool.ProductAttributeValue().Methods().NameGet().Extend("",
		func(rs pool.ProductAttributeValueSet) string {
			if rs.Env().Context().HasKey("show_attribute") && !rs.Env().Context().GetBool("show_attribute") {
				return rs.Super().NameGet()
			}
			return fmt.Sprintf("%s: %s", rs.Attribute().Name(), rs.Name())
		})

	pool.ProductAttributeValue().Methods().Unlink().Extend("",
		func(rs pool.ProductAttributeValueSet) int64 {
			linkedProducts := pool.ProductProduct().NewSet(rs.Env()).WithContext("active_test", false).Search(
				pool.ProductProduct().AttributeValues().In(rs))
			if !linkedProducts.IsEmpty() {
				log.Panic(rs.T(`The operation cannot be completed:
You are trying to delete an attribute value with a reference on a product variant.`))
			}
			return rs.Super().Unlink()
		})

	pool.ProductAttributeValue().Methods().VariantName().DeclareMethod(
		`VariantName returns a comma separated list of this product's
		attributes values of the given variable attributes'`,
		func(rs pool.ProductAttributeValueSet, variableAttribute pool.ProductAttributeSet) string {
			var names []string
			for _, attrValue := range pool.ProductAttributeValue().NewSet(rs.Env()).Browse(rs.Ids()).OrderBy("Attribute.Name").Records() {
				if attrValue.Attribute().Intersect(variableAttribute).IsEmpty() {
					continue
				}
				names = append(names, attrValue.Name())
			}
			return strings.Join(names, ", ")
		})

	pool.ProductAttributePrice().DeclareModel()

	pool.ProductAttributePrice().AddFields(map[string]models.FieldDefinition{
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: pool.ProductTemplate(),
			OnDelete: models.Cascade, Required: true},
		"Value": models.Many2OneField{String: "Product Attribute Value", RelationModel: pool.ProductAttributeValue(),
			OnDelete: models.Cascade, Required: true},
		"PriceExtra": models.FloatField{String: "Price Extra", Digits: decimalPrecision.GetPrecision("Product Price")},
	})

	pool.ProductAttributeLine().DeclareModel()

	pool.ProductAttributeLine().AddFields(map[string]models.FieldDefinition{
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: pool.ProductTemplate(),
			OnDelete: models.Cascade, Required: true},
		"Attribute": models.Many2OneField{RelationModel: pool.ProductAttribute(),
			OnDelete: models.Restrict, Required: true,
			Constraint: pool.ProductAttributeLine().Methods().CheckValidAttribute()},
		"Values": models.Many2ManyField{String: "Attribute Values", RelationModel: pool.ProductAttributeValue(),
			JSON: "value_ids", Constraint: pool.ProductAttributeLine().Methods().CheckValidAttribute()},
	})

	pool.ProductAttributeLine().Methods().CheckValidAttribute().DeclareMethod(
		`CheckValidAttribute check that attributes values are valid for the given attributes.`,
		func(rs pool.ProductAttributeLineSet) {
			for _, line := range rs.Records() {
				if !line.Values().Subtract(line.Attribute().Values()).IsEmpty() {
					log.Panic(rs.T("Error ! You cannot use this attribute with the following value."))
				}
			}
		})

	pool.ProductAttributeLine().Methods().NameGet().Extend("",
		func(rs pool.ProductAttributeLineSet) string {
			return rs.Attribute().NameGet()
		})

	pool.ProductAttributeLine().Methods().SearchByName().Extend("",
		func(rs pool.ProductAttributeLineSet, name string, op operator.Operator, additionalCond pool.ProductAttributeLineCondition, limit int) pool.ProductAttributeLineSet {
			// TDE FIXME: currently overriding the domain; however as it includes a
			// search on a m2o and one on a m2m, probably this will quickly become
			// difficult to compute - check if performance optimization is required
			if name != "" && op.IsPositive() {
				additionalCond = pool.ProductAttributeLine().Attribute().AddOperator(op, name).
					Or().Values().AddOperator(op, name)
			}
			return rs.Super().SearchByName(name, op, additionalCond, limit)
		})

}
