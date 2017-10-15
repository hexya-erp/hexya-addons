package product

import (
	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
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
			Default: func(models.Environment, models.FieldMap) interface{} {
				/*_get_default_category_id(self):
				  if self._context.get('categ_id') or self._context.get('default_categ_id'):
				      return self._context.get('categ_id') or self._context.get('default_categ_id')
				  category = self.env.ref('product.product_category_all', raise_if_not_found=False)
				  return category and category.type == 'normal' and category.id or False

				*/
				return 0
			}, Filter: pool.ProductCategory().Type().Equals("normal"), Required: true,
			Help: "Select category for the current product"},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(),
			Compute: pool.ProductTemplate().Methods().ComputeCurrencyId()},
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
			Inverse: pool.ProductTemplate().Methods().InverseStandardPrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Help:    "Cost of the product, in the default unit of measure of the product."},
		"Volume": models.FloatField{Compute: pool.ProductTemplate().Methods().ComputeVolume(),
			Inverse: pool.ProductTemplate().Methods().InverseVolume(), Help: "The volume in m3.", Stored: true},
		"Weight": models.FloatField{Compute: pool.ProductTemplate().Methods().ComputeWeight(),
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
			Default: func(models.Environment, models.FieldMap) interface{} {
				/*_get_default_uom_id(self):
				  return self.env["product.uom"].search([], limit=1, order='id').id */
				return 0
			}, Required: true, Help: "Default Unit of Measure used for all stock operation."},
		"UomPo": models.Many2OneField{String: "Purchase Unit of Measure", RelationModel: pool.ProductUom(),
			Default: func(models.Environment, models.FieldMap) interface{} {
				/*_get_default_uom_id(self):
				  return self.env["product.uom"].search([], limit=1, order='id').id*/
				return 0
			}, Required: true, Help: "Default Unit of Measure used for purchase orders. It must be in the same category than the default unit of measure."},
		"Company": models.Many2OneField{String: "Company", RelationModel: pool.Company(),
			Default: func(models.Environment, models.FieldMap) interface{} {
				/*lambda self: self.env['res.company']._company_default_get('product.template'*/
				return 0
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
			Compute: pool.ProductTemplate().Methods().ComputeProductVariantId()},
		"ProductVariantCount": models.IntegerField{String: "# Product Variants",
			Compute: pool.ProductTemplate().Methods().ComputeProductVariantCount()},
		"Barcode": models.CharField{},
		"DefaultCode": models.CharField{String: "Internal Reference",
			Compute: pool.ProductTemplate().Methods().ComputeDefaultCode(),
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

	pool.ProductTemplate().Methods().ComputeProductVariantId().DeclareMethod(
		`ComputeProductVariantId`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.depends('product_variant_ids')
			/*def _compute_product_variant_id(self):
			  for p in self:
			      p.product_variant_id = p.product_variant_ids[:1].id

			*/
		})

	pool.ProductTemplate().Methods().ComputeCurrencyId().DeclareMethod(
		`ComputeCurrencyId`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.multi
			/*def _compute_currency_id(self):
			  try:
			      main_company = self.sudo().env.ref('base.main_company')
			  except ValueError:
			      main_company = self.env['res.company'].sudo().search([], limit=1, order="id")
			  for template in self:
			      template.currency_id = template.company_id.sudo().currency_id.id or main_company.currency_id.id

			*/
		})

	pool.ProductTemplate().Methods().ComputeTemplatePrice().DeclareMethod(
		`ComputeTemplatePrice`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.multi
			/*def _compute_template_price(self):
			  prices = {}
			  pricelist_id_or_name = self._context.get('pricelist')
			  if pricelist_id_or_name:
			      pricelist = None
			      partner = self._context.get('partner')
			      quantity = self._context.get('quantity', 1.0)

			      # Support context pricelists specified as display_name or ID for compatibility
			      if isinstance(pricelist_id_or_name, basestring):
			          pricelist_data = self.env['product.pricelist'].name_search(pricelist_id_or_name, operator='=', limit=1)
			          if pricelist_data:
			              pricelist = self.env['product.pricelist'].browse(pricelist_data[0][0])
			      elif isinstance(pricelist_id_or_name, (int, long)):
			          pricelist = self.env['product.pricelist'].browse(pricelist_id_or_name)

			      if pricelist:
			          quantities = [quantity] * len(self)
			          partners = [partner] * len(self)
			          prices = pricelist.get_products_price(self, quantities, partners)

			  for template in self:
			      template.price = prices.get(template.id, 0.0)

			*/
		})

	pool.ProductTemplate().Methods().InverseTemplatePrice().DeclareMethod(
		`InverseTemplatePrice`,
		func(rs pool.ProductTemplateSet, price float64) {
			//@api.multi
			/*def _set_template_price(self):
			  if self._context.get('uom'):
			      for template in self:
			          value = self.env['product.uom'].browse(self._context['uom'])._compute_price(template.price, template.uom_id)
			          template.write({'list_price': value})
			  else:
			      self.write({'list_price': self.price})

			*/
		})

	pool.ProductTemplate().Methods().ComputeStandardPrice().DeclareMethod(
		`ComputeStandardPrice`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.depends('product_variant_ids','product_variant_ids.standard_price')
			/*def _compute_standard_price(self):
			  unique_variants = self.filtered(lambda template: len(template.product_variant_ids) == 1)
			  for template in unique_variants:
			      template.standard_price = template.product_variant_ids.standard_price
			  for template in (self - unique_variants):
			      template.standard_price = 0.0

			*/
		})

	pool.ProductTemplate().Methods().InverseStandardPrice().DeclareMethod(
		`InverseStandardPrice`,
		func(rs pool.ProductTemplateSet, price float64) {
			//@api.one
			/*def _set_standard_price(self):
			  if len(self.product_variant_ids) == 1:
			      self.product_variant_ids.standard_price = self.standard_price

			*/
		})

	pool.ProductTemplate().Methods().ComputeVolume().DeclareMethod(
		`ComputeVolume`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.depends('product_variant_ids','product_variant_ids.volume')
			/*def _compute_volume(self):
			  unique_variants = self.filtered(lambda template: len(template.product_variant_ids) == 1)
			  for template in unique_variants:
			      template.volume = template.product_variant_ids.volume
			  for template in (self - unique_variants):
			      template.volume = 0.0

			*/
		})

	pool.ProductTemplate().Methods().InverseVolume().DeclareMethod(
		`InverseVolume`,
		func(rs pool.ProductTemplateSet, volume float64) {
			//@api.one
			/*def _set_volume(self):
			  if len(self.product_variant_ids) == 1:
			      self.product_variant_ids.volume = self.volume

			*/
		})

	pool.ProductTemplate().Methods().ComputeWeight().DeclareMethod(
		`ComputeWeight`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.depends('product_variant_ids','product_variant_ids.weight')
			/*def _compute_weight(self):
			  unique_variants = self.filtered(lambda template: len(template.product_variant_ids) == 1)
			  for template in unique_variants:
			      template.weight = template.product_variant_ids.weight
			  for template in (self - unique_variants):
			      template.weight = 0.0

			*/
		})
	pool.ProductTemplate().Methods().InverseWeight().DeclareMethod(
		`InverseWeight`,
		func(rs pool.ProductTemplateSet, weight float64) {
			//@api.one
			/*def _set_weight(self):
			  if len(self.product_variant_ids) == 1:
			      self.product_variant_ids.weight = self.weight

			*/
		})
	pool.ProductTemplate().Methods().ComputeProductVariantCount().DeclareMethod(
		`ComputeProductVariantCount`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.depends('product_variant_ids.product_tmpl_id')
			/*def _compute_product_variant_count(self):
			  self.product_variant_count = len(self.product_variant_ids)

			*/
		})

	pool.ProductTemplate().Methods().ComputeDefaultCode().DeclareMethod(
		`ComputeDefaultCode`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.depends('product_variant_ids','product_variant_ids.default_code')
			/*def _compute_default_code(self):
			  unique_variants = self.filtered(lambda template: len(template.product_variant_ids) == 1)
			  for template in unique_variants:
			      template.default_code = template.product_variant_ids.default_code
			  for template in (self - unique_variants):
			      template.default_code = ''

			*/
		})

	pool.ProductTemplate().Methods().InverseDefaultCode().DeclareMethod(
		`InverseDefaultCode`,
		func(rs pool.ProductTemplateSet, code string) {
			//@api.one
			/*def _set_default_code(self):
			  if len(self.product_variant_ids) == 1:
			      self.product_variant_ids.default_code = self.default_code

			*/
		})

	pool.ProductTemplate().Methods().CheckUom().DeclareMethod(
		`CheckUom`,
		func(rs pool.ProductTemplateSet) {
			//@api.constrains('uom_id','uom_po_id')
			/*def _check_uom(self):
			  if any(template.uom_id and template.uom_po_id and template.uom_id.category_id != template.uom_po_id.category_id for template in self):
			      raise ValidationError(_('Error: The default Unit of Measure and the purchase Unit of Measure must be in the same category.'))
			  return True

			*/
		})

	pool.ProductTemplate().Methods().OnchangeUomId().DeclareMethod(
		`OnchangeUomId`,
		func(rs pool.ProductTemplateSet) (*pool.ProductTemplateSet, []models.FieldNamer) {
			//@api.onchange('uom_id')
			/*def _onchange_uom_id(self):
			  if self.uom_id:
			      self.uom_po_id = self.uom_id.id

			*/
		})

	pool.ProductTemplate().Methods().Create().Extend("",
		func(rs pool.ProductTemplateSet, data *pool.ProductTemplateData) pool.ProductTemplateSet {
			//@api.model
			/*def create(self, vals):
			  ''' Store the initial standard price in order to be able to retrieve the cost of a product template for a given date'''
			  # TDE FIXME: context brol
			  tools.image_resize_images(vals)
			  template = super(ProductTemplate, self).create(vals)
			  if "create_product_product" not in self._context:
			      template.with_context(create_from_tmpl=True).create_variant_ids()

			  # This is needed to set given values to first variant after creation
			  related_vals = {}
			  if vals.get('barcode'):
			      related_vals['barcode'] = vals['barcode']
			  if vals.get('default_code'):
			      related_vals['default_code'] = vals['default_code']
			  if vals.get('standard_price'):
			      related_vals['standard_price'] = vals['standard_price']
			  if vals.get('volume'):
			      related_vals['volume'] = vals['volume']
			  if vals.get('weight'):
			      related_vals['weight'] = vals['weight']
			  if related_vals:
			      template.write(related_vals)
			  return template

			*/
		})

	pool.ProductTemplate().Methods().Write().Extend("",
		func(rs pool.ProductTemplateSet, vals *pool.ProductTemplateData, fieldsToUnset ...models.FieldNamer) bool {
			//@api.multi
			/*def write(self, vals):
			  tools.image_resize_images(vals)
			  res = super(ProductTemplate, self).write(vals)
			  if 'attribute_line_ids' in vals or vals.get('active'):
			      self.create_variant_ids()
			  if 'active' in vals and not vals.get('active'):
			      self.with_context(active_test=False).mapped('product_variant_ids').write({'active': vals.get('active')})
			  return res

			*/
		})

	pool.ProductTemplate().Methods().Copy().Extend("",
		func(rs pool.ProductTemplateSet, overrides *pool.ProductTemplateData, fieldsToUnset ...models.FieldNamer) pool.ProductTemplateSet {
			//@api.multi
			/*def copy(self, default=None):
			  # TDE FIXME: should probably be copy_data
			  self.ensure_one()
			  if default is None:
			      default = {}
			  if 'name' not in default:
			      default['name'] = _("%s (copy)") % self.name
			  return super(ProductTemplate, self).copy(default=default)

			*/
		})

	pool.ProductTemplate().Methods().NameGet().Extend("",
		func(rs pool.ProductTemplateSet) string {
			//@api.multi
			/*def name_get(self):
			  return [(template.id, '%s%s' % (template.default_code and '[%s] ' % template.default_code or '', template.name))
			          for template in self]

			*/
		})

	pool.ProductTemplate().Methods().SearchByName().Extend("",
		func(rs pool.ProductTemplateSet, name string, op operator.Operator, additionalCond pool.ProductTemplateCondition, limit int) pool.ProductTemplateSet {
			//@api.model
			/*def name_search(self, name='', args=None, operator='ilike', limit=100):
			  # Only use the product.product heuristics if there is a search term and the domain
			  # does not specify a match on `product.template` IDs.
			  if not name or any(term[0] == 'id' for term in (args or [])):
			      return super(ProductTemplate, self).name_search(name=name, args=args, operator=operator, limit=limit)

			  Product = self.env['product.product']
			  templates = self.browse([])
			  while True:
			      domain = templates and [('product_tmpl_id', 'not in', templates.ids)] or []
			      args = args if args is not None else []
			      products_ns = Product.name_search(name, args+domain, operator=operator)
			      products = Product.browse([x[0] for x in products_ns])
			      templates |= products.mapped('product_tmpl_id')
			      if (not products) or (limit and (len(templates) > limit)):
			          break

			  # re-apply product.template order + name_get
			  return super(ProductTemplate, self).name_search(
			      '', args=[('id', 'in', list(set(templates.ids)))],
			      operator='ilike', limit=limit)

			*/
		})

	pool.ProductTemplate().Methods().PriceCompute().DeclareMethod(
		`PriceCompute`,
		func(rs pool.ProductTemplateSet, priceType string, uom pool.ProductUomSet, currency pool.CurrencySet, company pool.CompanySet) {
			//@api.multi
			/*def price_compute(self, price_type, uom=False, currency=False, company=False):
			      # TDE FIXME: delegate to template or not ? fields are reencoded here ...
			      # compatibility about context keys used a bit everywhere in the code
			      if not uom and self._context.get('uom'):
			          uom = self.env['product.uom'].browse(self._context['uom'])
			      if not currency and self._context.get('currency'):
			          currency = self.env['res.currency'].browse(self._context['currency'])

			      templates = self
			      if price_type == 'standard_price':
			          # standard_price field can only be seen by users in base.group_user
			          # Thus, in order to compute the sale price from the cost for users not in this group
			          # We fetch the standard price as the superuser
			          templates = self.with_context(force_company=company and company.id or self._context.get('force_company', self.env.user.company_id.id)).sudo()

			      prices = dict.fromkeys(self.ids, 0.0)
			      for template in templates:
			          prices[template.id] = template[price_type] or 0.0

			          if uom:
			              prices[template.id] = template.uom_id._compute_price(prices[template.id], uom)

			          # Convert from current user company currency to asked one
			          # This is right cause a field cannot be in more than one currency
			          if currency:
			              prices[template.id] = template.currency_id.compute(prices[template.id], currency)

			      return prices

			  # compatibility to remove after v10 - DEPRECATED
			*/
		})

	pool.ProductTemplate().Methods().CreateVariantIds().DeclareMethod(
		`CreateVariantIds`,
		func(rs pool.ProductTemplateSet) {
			//@api.multi
			/*def create_variant_ids(self):
			  Product = self.env["product.product"]
			  for tmpl_id in self.with_context(active_test=False):
			      # adding an attribute with only one value should not recreate product
			      # write this attribute on every product to make sure we don't lose them
			      variant_alone = tmpl_id.attribute_line_ids.filtered(lambda line: len(line.value_ids) == 1).mapped('value_ids')
			      for value_id in variant_alone:
			          updated_products = tmpl_id.product_variant_ids.filtered(lambda product: value_id.attribute_id not in product.mapped('attribute_value_ids.attribute_id'))
			          updated_products.write({'attribute_value_ids': [(4, value_id.id)]})

			      # list of values combination
			      existing_variants = [set(variant.attribute_value_ids.filtered(lambda r: r.attribute_id.create_variant).ids) for variant in tmpl_id.product_variant_ids]
			      variant_matrix = itertools.product(*(line.value_ids for line in tmpl_id.attribute_line_ids if line.value_ids and line.value_ids[0].attribute_id.create_variant))
			      variant_matrix = map(lambda record_list: reduce(lambda x, y: x+y, record_list, self.env['product.attribute.value']), variant_matrix)
			      to_create_variants = filter(lambda rec_set: set(rec_set.ids) not in existing_variants, variant_matrix)

			      # check product
			      variants_to_activate = self.env['product.product']
			      variants_to_unlink = self.env['product.product']
			      for product_id in tmpl_id.product_variant_ids:
			          if not product_id.active and product_id.attribute_value_ids.filtered(lambda r: r.attribute_id.create_variant) in variant_matrix:
			              variants_to_activate |= product_id
			          elif product_id.attribute_value_ids.filtered(lambda r: r.attribute_id.create_variant) not in variant_matrix:
			              variants_to_unlink |= product_id
			      if variants_to_activate:
			          variants_to_activate.write({'active': True})

			      # create new product
			      for variant_ids in to_create_variants:
			          new_variant = Product.create({
			              'product_tmpl_id': tmpl_id.id,
			              'attribute_value_ids': [(6, 0, variant_ids.ids)]
			          })

			      # unlink or inactive product
			      for variant in variants_to_unlink:
			          try:
			              with self._cr.savepoint(), tools.mute_logger('odoo.sql_db'):
			                  variant.unlink()
			          # We catch all kind of exception to be sure that the operation doesn't fail.
			          except (psycopg2.Error, except_orm):
			              variant.write({'active': False})
			              pass
			  return True
			*/
		})

}
