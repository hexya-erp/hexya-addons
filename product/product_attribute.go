package product

import (
	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/models"
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
			Help: "Price Extra: Extra price for the variant with this attribute value on sale price. eg. 200 price extra, 1000 + 200 = 1200." /*[ 1000 + 200   1200."]*/},
		"Prices": models.One2ManyField{String: "Attribute Prices", RelationModel: pool.ProductAttributePrice(),
			ReverseFK: "Value", JSON: "price_ids" /* readonly */},
	})

	pool.ProductAttributeValue().AddSQLConstraint("ValueCompanyUniq", "unique (name,attribute_id)", "This attribute value already exists !")

	pool.ProductAttributeValue().Methods().ComputePriceExtra().DeclareMethod(
		`ComputePriceExtra`,
		func(rs pool.ProductAttributeValueSet) (*pool.ProductAttributeValueData, []models.FieldNamer) {
			//@api.one
			/*def _compute_price_extra(self):
			  if self._context.get('active_id'):
			      price = self.price_ids.filtered(lambda price: price.product_tmpl_id.id == self._context['active_id'])
			      self.price_extra = price.price_extra
			  else:
			      self.price_extra = 0.0

			*/
		})

	pool.ProductAttributeValue().Methods().InversePriceExtra().DeclareMethod(
		`InversePriceExtra`,
		func(rs pool.ProductAttributeValueSet, value float64) {
			/*def _set_price_extra(self):
			  if not self._context.get('active_id'):
			      return

			  AttributePrice = self.env['product.attribute.price']
			  prices = AttributePrice.search([('value_id', 'in', self.ids), ('product_tmpl_id', '=', self._context['active_id'])])
			  updated = prices.mapped('value_id')
			  if prices:
			      prices.write({'price_extra': self.price_extra})
			  else:
			      for value in self - updated:
			          AttributePrice.create({
			              'product_tmpl_id': self._context['active_id'],
			              'value_id': value.id,
			              'price_extra': self.price_extra,
			          })

			*/
		})

	pool.ProductAttributeValue().Methods().NameGet().Extend("",
		func(rs pool.ProductAttributeValueSet) string {
			//@api.multi
			/*def name_get(self):
			  if not self._context.get('show_attribute', True):  # TDE FIXME: not used
			      return super(ProductAttributevalue, self).name_get()
			  return [(value.id, "%s: %s" % (value.attribute_id.name, value.name)) for value in self]

			*/
		})

	pool.ProductAttributeValue().Methods().Unlink().Extend("",
		func(rs pool.ProductAttributeValueSet) int64 {
			//@api.multi
			/*def unlink(self):
			  linked_products = self.env['product.product'].with_context(active_test=False).search([('attribute_value_ids', 'in', self.ids)])
			  if linked_products:
			      raise UserError(_('The operation cannot be completed:\nYou are trying to delete an attribute value with a reference on a product variant.'))
			  return super(ProductAttributevalue, self).unlink()

			*/
		})

	pool.ProductAttributeValue().Methods().VariantName().DeclareMethod(
		`VariantName`,
		func(rs pool.ProductAttributeValueSet, variableAttribute pool.ProductAttributeSet) string {
			//@api.multi
			/*def _variant_name(self, variable_attributes):
			  return ", ".join([v.name for v in self.sorted(key=lambda r: r.attribute_id.name) if v.attribute_id in variable_attributes])


			*/
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
			OnDelete: models.Restrict, Required: true},
		"Values": models.Many2ManyField{String: "Attribute Values", RelationModel: pool.ProductAttributeValue(),
			JSON: "value_ids"},
	})

	pool.ProductAttributeLine().Methods().CheckValidAttribute().DeclareMethod(
		`CheckValidAttribute`,
		func(rs pool.ProductAttributeLineSet) {
			//@api.constrains('value_ids','attribute_id')
			/*def _check_valid_attribute(self):
			  if any(line.value_ids > line.attribute_id.value_ids for line in self):
			      raise ValidationError(_('Error ! You cannot use this attribute with the following value.'))
			  return True

			*/
		})

	pool.ProductAttributeLine().Methods().NameGet().Extend("",
		func(rs pool.ProductAttributeLineSet) string {
			return rs.Attribute().NameGet()
		})

	pool.ProductAttributeLine().Methods().NameSearch().DeclareMethod(
		`NameSearch`,
		func(rs pool.ProductAttributeLineSet, args struct {
			Name     interface{}
			Args     interface{}
			Operator interface{}
			Limit    interface{}
		}) {
			//@api.model
			/*def name_search(self, name='', args=None, operator='ilike', limit=100):
			  # TDE FIXME: currently overriding the domain; however as it includes a
			  # search on a m2o and one on a m2m, probably this will quickly become
			  # difficult to compute - check if performance optimization is required
			  if name and operator in ('=', 'ilike', '=ilike', 'like', '=like'):
			      new_args = ['|', ('attribute_id', operator, name), ('value_ids', operator, name)]
			  else:
			      new_args = args
			  return super(ProductAttributeLine, self).name_search(name=name, args=new_args, operator=operator, limit=limit)
			*/
		})

}
