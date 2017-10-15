package product

import (
	"fmt"
	"log"
	"strings"

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

	pool.ProductCategory().DeclareModel()
	pool.ProductCategory().SetDefaultOrder("Parent.Name")

	pool.ProductCategory().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Index: true, Required: true, Translate: true},
		"Parent": models.Many2OneField{String: "Parent Category", RelationModel: pool.ProductCategory(), Index: true,
			OnDelete: models.Cascade, Constraint: pool.ProductCategory().Methods().CheckCategoryRecursion()},
		"Children": models.One2ManyField{String: "Child Categories", RelationModel: pool.ProductCategory(),
			ReverseFK: "Parent", JSON: "child_id"},
		"Type": models.SelectionField{String: "Category Type", Selection: types.Selection{"view": "View", "normal": "Normal"},
			Default: models.DefaultValue("normal"), Help: "A category of the view type is a virtual category that can be used as the parent of another category to create a hierarchical structure."},
		"Products": models.One2ManyField{RelationModel: pool.ProductTemplate(), ReverseFK: "Categ"},
		"ProductCount": models.IntegerField{String: "# Products", Compute: pool.ProductCategory().Methods().ComputeProductCount(),
			Help:    "The number of products under this category (Does not consider the children categories)",
			Depends: []string{"Products"}, GoType: new(int)},
	})

	pool.ProductCategory().Methods().ComputeProductCount().DeclareMethod(
		`ComputeProductCount returns the number of products within this category (not considering children categories)`,
		func(rs pool.ProductCategorySet) (*pool.ProductCategoryData, []models.FieldNamer) {
			return &pool.ProductCategoryData{
				ProductCount: pool.ProductTemplate().Search(rs.Env(), pool.ProductTemplate().Categ().Equals(rs)).SearchCount(),
			}, []models.FieldNamer{pool.ProductCategory().ProductCount()}
		})

	pool.ProductCategory().Methods().CheckCategoryRecursion().DeclareMethod(
		`CheckCategoryRecursion panics if there is a recursion in the category tree.`,
		func(rs pool.ProductCategorySet) {
			if !rs.CheckRecursion() {
				log.Panic(rs.T("Error ! You cannot create recursive categories."))
			}
		})

	pool.ProductCategory().Methods().NameGet().Extend("",
		func(rs pool.ProductCategorySet) string {
			var names []string
			for current := rs; !current.IsEmpty(); current = current.Parent() {
				names = append([]string{current.Name()}, names...)
			}
			return strings.Join(names, " / ")
		})

	pool.ProductCategory().Methods().SearchByName().Extend("",
		func(rs pool.ProductCategorySet, name string, op operator.Operator, additionalCond pool.ProductCategoryCondition, limit int) pool.ProductCategorySet {
			if name == "" {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			// Be sure name_search is symetric to name_get
			categoryNames := strings.Split(name, " / ")
			child := categoryNames[len(categoryNames)-1]
			cond := pool.ProductCategory().Name().AddOperator(op, child)
			var categories pool.ProductCategorySet
			if len(categoryNames) > 1 {
				parents := rs.SearchByName(strings.Join(categoryNames[:len(categoryNames)-1], " / "), operator.IContains, additionalCond, limit)
				if op.IsNegative() {
					categories = pool.ProductCategory().Search(rs.Env(), pool.ProductCategory().ID().NotIn(parents.Ids()))
					cond = cond.Or().Parent().In(categories)
				} else {
					cond = cond.And().Parent().In(parents)
				}
				for i := 1; i < len(categoryNames); i++ {
					if op.IsNegative() {
						cond = cond.AndCond(pool.ProductCategory().Name().AddOperator(op, strings.Join(categoryNames[len(categoryNames)-1-i:], " / ")))
					} else {
						cond = cond.OrCond(pool.ProductCategory().Name().AddOperator(op, strings.Join(categoryNames[len(categoryNames)-1-i:], " / ")))
					}
				}
			}
			return pool.ProductCategory().Search(rs.Env(), cond.AndCond(additionalCond))
		})

	pool.ProductPriceHistory().DeclareModel()
	pool.ProductPriceHistory().SetDefaultOrder("Datetime DESC")

	pool.ProductPriceHistory().AddFields(map[string]models.FieldDefinition{
		"Company": models.Many2OneField{RelationModel: pool.Company(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				if env.Context().HasKey("force_company") {

					return pool.Company().Browse(env, []int64{env.Context().GetInteger("force_company")})
				}
				currentUser := pool.User().NewSet(env).CurrentUser()
				return currentUser.Company()
			}, Required: true},
		"Product": models.Many2OneField{RelationModel: pool.ProductProduct(), JSON: "product_id",
			OnDelete: models.Cascade, Required: true},
		"Datetime": models.DateTimeField{String: "Date", Default: func(env models.Environment, vals models.FieldMap) interface{} {
			return dates.Now()
		}},
		"Cost": models.FloatField{String: "Cost", Digits: decimalPrecision.GetPrecision("Product Price")},
	})

	pool.ProductProduct().DeclareModel()
	pool.ProductProduct().SetDefaultOrder("DefaultCode", "Name", "ID")

	pool.ProductProduct().AddFields(map[string]models.FieldDefinition{
		"Price": models.FloatField{Compute: pool.ProductProduct().Methods().ComputeProductPrice(),
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Inverse: pool.ProductProduct().Methods().InverseProductPrice()},
		"PriceExtra": models.FloatField{String: "Variant Price Extra",
			Compute: pool.ProductProduct().Methods().ComputeProductPriceExtra(),
			Depends: []string{"AttributeValues", "AttributeValues.Prices", "AttributeValues.Prices.PriceExtra", "AttributeValues.Prices.PriceExtra.ProductTmpl"},
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Help:    "This is the sum of the extra price of all attributes"},
		"LstPrice": models.FloatField{String: "Sale Price",
			Compute: pool.ProductProduct().Methods().ComputeProductLstPrice(),
			Depends: []string{"ListPrice", "PriceExtra"},
			Digits:  decimalPrecision.GetPrecision("Product Price"),
			Inverse: pool.ProductProduct().Methods().InverseProductLstPrice(),
			Help:    "The sale price is managed from the product template. Click on the 'Variant Prices' button to set the extra attribute prices."},
		"DefaultCode": models.CharField{String: "Internal Reference", Index: true},
		"Code": models.CharField{String: "Internal Reference",
			Compute: pool.ProductProduct().Methods().ComputeProductCode(), Depends: []string{""}},
		"PartnerRef": models.CharField{String: "Customer Ref",
			Compute: pool.ProductProduct().Methods().ComputePartnerRef(), Depends: []string{""}},
		"Active": models.BooleanField{String: "Active",
			Default: models.DefaultValue(true),
			Help:    "If unchecked, it will allow you to hide the product without removing it."},
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: pool.ProductTemplate(),
			Index: true, OnDelete: models.Cascade, Required: true, Embed: true},
		"Barcode": models.CharField{String: "Barcode", NoCopy: true, Unique: true,
			Help: "International Article Number used for product identification."},
		"AttributeValues": models.Many2ManyField{String: "Attributes", RelationModel: pool.ProductAttributeValue(),
			JSON: "attribute_value_ids" /*, OnDelete: models.Restrict*/},
		"ImageVariant": models.BinaryField{String: "Variant Image",
			Help: "This field holds the image used as image for the product variant, limited to 1024x1024px."},
		"Image": models.BinaryField{String: "Big-sized image",
			Compute: pool.ProductProduct().Methods().ComputeImages(),
			Depends: []string{"ImageVariant", "ProductTmpl", "ProductTmpl.Image"},
			Inverse: pool.ProductProduct().Methods().InverseImageValue(),
			Help: `Image of the product variant (Big-sized image of product template if false). It is automatically
resized as a 1024x1024px image, with aspect ratio preserved.`},
		"ImageSmall": models.BinaryField{String: "Small-sized image",
			Compute: pool.ProductProduct().Methods().ComputeImages(),
			Depends: []string{"ImageVariant", "ProductTmpl", "ProductTmpl.Image"},
			Inverse: pool.ProductProduct().Methods().InverseImageValue(),
			Help:    "Image of the product variant (Small-sized image of product template if false)."},
		"ImageMedium": models.BinaryField{String: "Medium-sized image",
			Compute: pool.ProductProduct().Methods().ComputeImages(),
			Depends: []string{"ImageVariant", "ProductTmpl", "ProductTmpl.Image"},
			Inverse: pool.ProductProduct().Methods().InverseImageValue(),
			Help:    "Image of the product variant (Medium-sized image of product template if false)."},
		"StandardPrice": models.FloatField{String: "Cost", /*, CompanyDependent : true*/
			Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: `Cost of the product template used for standard stock valuation in accounting and used as a
base price on purchase orders. Expressed in the default unit of measure of the product.`},
		"Volume": models.FloatField{Help: "The volume in m3."},
		"Weight": models.FloatField{Digits: decimalPrecision.GetPrecision("Stock Weight"),
			Help: "The weight of the contents in Kg, not including any packaging, etc."},
		"PricelistItems": models.Many2ManyField{RelationModel: pool.ProductPricelistItem(),
			JSON: "pricelist_item_ids", Compute: pool.ProductProduct().Methods().GetPricelistItems()},
	})

	pool.ProductProduct().Fields().StandardPrice().RevokeAccess(security.GroupEveryone, security.All).GrantAccess(base.GroupUser, security.All)

	pool.ProductProduct().Methods().ComputeProductPrice().DeclareMethod(
		`ComputeProductPrice computes the price of this product based on the given context keys:
		- 'partner' => int64 (id of the partner)
		- 'pricelist' => int64 (id of the price list)
		- 'quantity' => float64`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			if !rs.Env().Context().HasKey("pricelist") {
				return new(pool.ProductProductData), []models.FieldNamer{pool.ProductProduct().Price()}
			}
			priceListID := rs.Env().Context().GetInteger("pricelist")
			priceList := pool.ProductPricelist().Browse(rs.Env(), []int64{priceListID})
			quantity := rs.Env().Context().GetFloat("quantity")
			partnerID := rs.Env().Context().GetInteger("partner")
			partner := pool.Partner().Browse(rs.Env(), []int64{partnerID})
			if quantity == 0 {
				quantity = 1
			}
			if priceList.IsEmpty() {
				return new(pool.ProductProductData), []models.FieldNamer{pool.ProductProduct().Price()}
			}
			return &pool.ProductProductData{
				Price: priceList.GetProductPrice(rs, quantity, partner, dates.Today(), pool.ProductUom().NewSet(rs.Env())),
			}, []models.FieldNamer{pool.ProductProduct().Price()}
		})

	pool.ProductProduct().Methods().InverseProductPrice().DeclareMethod(
		`InverseProductPrice updates ListPrice from the given Price`,
		func(rs pool.ProductProductSet, price float64) {
			if rs.Env().Context().HasKey("uom") {
				price = pool.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")}).ComputePrice(price, rs.Uom())
			}
			price -= rs.PriceExtra()
			rs.SetListPrice(price)
		})

	pool.ProductProduct().Methods().InverseProductLstPrice().DeclareMethod(
		`InverseProductLstPrice updates ListPrice from the given LstPrice`,
		func(rs pool.ProductProductSet, price float64) {
			if rs.Env().Context().HasKey("uom") {
				price = pool.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")}).ComputePrice(price, rs.Uom())
			}
			price -= rs.PriceExtra()
			rs.SetListPrice(price)
		})

	pool.ProductProduct().Methods().ComputeProductPriceExtra().DeclareMethod(
		`ComputeProductPriceExtra computes the price extra of this product by suming the extras of each attribute`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			var priceExtra float64
			for _, attributeValue := range rs.AttributeValues().Records() {
				for _, attributePrice := range attributeValue.Prices().Records() {
					if attributePrice.ProductTmpl().Equals(rs.ProductTmpl()) {
						priceExtra += attributePrice.PriceExtra()
					}
				}
			}
			return &pool.ProductProductData{
				PriceExtra: priceExtra,
			}, []models.FieldNamer{pool.ProductProduct().PriceExtra()}
		})

	pool.ProductProduct().Methods().ComputeProductLstPrice().DeclareMethod(
		`ComputeProductLstPrice computes the LstPrice from the ListPrice and the extras`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			listPrice := rs.ListPrice()
			if rs.Env().Context().HasKey("uom") {
				toUoM := pool.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
				listPrice = rs.Uom().ComputePrice(listPrice, toUoM)
			}
			return &pool.ProductProductData{
				LstPrice: listPrice + rs.PriceExtra(),
			}, []models.FieldNamer{pool.ProductProduct().LstPrice()}
		})

	pool.ProductProduct().Methods().ComputeProductCode().DeclareMethod(
		`ComputeProductCode computes the product code based on the context:
- 'partner_id' => int64 (id of the considered partner)`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			var code string
			for _, supplierInfo := range rs.Sellers().Records() {
				if supplierInfo.Name().ID() == rs.Env().Context().GetInteger("partner_id") {
					code = supplierInfo.ProductCode()
					if code != "" {
						break
					}
				}
			}
			if code != "" {
				code = rs.DefaultCode()
			}
			return &pool.ProductProductData{
				Code: code,
			}, []models.FieldNamer{pool.ProductProduct().Code()}
		})

	pool.ProductProduct().Methods().ComputePartnerRef().DeclareMethod(
		`ComputePartnerRef computes the product's reference (i.e. "[code] description") based on the context:
- 'partner_id' => int64 (id of the considered partner)`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			var code, productName string
			for _, supplierInfo := range rs.Sellers().Records() {
				if supplierInfo.Name().ID() == rs.Env().Context().GetInteger("partner_id") {
					code = supplierInfo.ProductCode()
					if code != "" {
						break
					}
					productName = supplierInfo.ProductName()
					if productName != "" {
						break
					}
				}
			}
			if code != "" {
				code = rs.DefaultCode()
			}
			if productName != "" {
				productName = rs.Name()
			}
			partnerRef := fmt.Sprintf("[%s]%s", code, productName)
			if code == "" {
				partnerRef = productName
			}
			return &pool.ProductProductData{
				PartnerRef: partnerRef,
			}, []models.FieldNamer{pool.ProductProduct().PartnerRef()}
		})

	pool.ProductProduct().Methods().ComputeImages().DeclareMethod(
		`ComputeImages computes the images in different sizes.`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			// TODO implement image resizing
			//@api.depends('image_variant','product_tmpl_id.image')
			/*def _compute_images(self):
			  if self._context.get('bin_size'):
			      self.image_medium = self.image_variant
			      self.image_small = self.image_variant
			      self.image = self.image_variant
			  else:
			      resized_images = tools.image_get_resized_images(self.image_variant, return_big=True, avoid_resize_medium=True)
			      self.image_medium = resized_images['image_medium']
			      self.image_small = resized_images['image_small']
			      self.image = resized_images['image']
			  if not self.image_medium:
			      self.image_medium = self.product_tmpl_id.image_medium
			  if not self.image_small:
			      self.image_small = self.product_tmpl_id.image_small
			  if not self.image:
			      self.image = self.product_tmpl_id.image

			*/
			return &pool.ProductProductData{}, []models.FieldNamer{}
		})

	pool.ProductProduct().Methods().InverseImageValue().DeclareMethod(
		`InverseImageValue`,
		func(rs pool.ProductProductSet, args struct {
			Value interface{}
		}) {
			//@api.one
			/*def _set_image_value(self, value):
			  image = tools.image_resize_image_big(value)
			  if self.product_tmpl_id.image:
			      self.image_variant = image
			  else:
			      self.product_tmpl_id.image = image

			*/
		})
	pool.ProductProduct().Methods().GetPricelistItems().DeclareMethod(
		`GetPricelistItems`,
		func(rs pool.ProductProductSet) {
			//@api.one
			/*def _get_pricelist_items(self):
			  self.pricelist_item_ids = self.env['product.pricelist.item'].search([
			      '|',
			      ('product_id', '=', self.id),
			      ('product_tmpl_id', '=', self.product_tmpl_id.id)]).ids

			*/
		})
	pool.ProductProduct().Methods().CheckAttributeValueIds().DeclareMethod(
		`CheckAttributeValueIds`,
		func(rs pool.ProductProductSet) {
			//@api.constrains('attribute_value_ids')
			/*def _check_attribute_value_ids(self):
			  for product in self:
			      attributes = self.env['product.attribute']
			      for value in product.attribute_value_ids:
			          if value.attribute_id in attributes:
			              raise ValidationError(_('Error! It is not allowed to choose more than one value for a given attribute.'))
			          attributes |= value.attribute_id
			  return True

			*/
		})
	pool.ProductProduct().Methods().OnchangeUom().DeclareMethod(
		`OnchangeUom`,
		func(rs pool.ProductProductSet) {
			//@api.onchange('uom_id','uom_po_id')
			/*def _onchange_uom(self):
			  if self.uom_id and self.uom_po_id and self.uom_id.category_id != self.uom_po_id.category_id:
			      self.uom_po_id = self.uom_id

			*/
		})
	pool.ProductProduct().Methods().Create().Extend("",
		func(rs pool.ProductProductSet, data *pool.ProductProductData) pool.ProductProductSet {
			//@api.model
			/*def create(self, vals):
			  product = super(ProductProduct, self.with_context(create_product_product=True)).create(vals)
			  # When a unique variant is created from tmpl then the standard price is set by _set_standard_price
			  if not (self.env.context.get('create_from_tmpl') and len(product.product_tmpl_id.product_variant_ids) == 1):
			      product._set_standard_price(vals.get('standard_price') or 0.0)
			  return product

			*/
		})

	pool.ProductProduct().Methods().Write().Extend("",
		func(rs pool.ProductProductSet, data *pool.ProductProductData, fieldsToUnset ...models.FieldNamer) bool {
			//@api.multi
			/*def write(self, values):
			  ''' Store the standard price change in order to be able to retrieve the cost of a product for a given date'''
			  res = super(ProductProduct, self).write(values)
			  if 'standard_price' in values:
			      self._set_standard_price(values['standard_price'])
			  return res

			*/
		})

	pool.ProductProduct().Methods().Unlink().Extend("",
		func(rs pool.ProductProductSet) int64 {
			//@api.multi
			/*def unlink(self):
			  unlink_products = self.env['product.product']
			  unlink_templates = self.env['product.template']
			  for product in self:
			      # Check if product still exists, in case it has been unlinked by unlinking its template
			      if not product.exists():
			          continue
			      # Check if the product is last product of this template
			      other_products = self.search([('product_tmpl_id', '=', product.product_tmpl_id.id), ('id', '!=', product.id)])
			      if not other_products:
			          unlink_templates |= product.product_tmpl_id
			      unlink_products |= product
			  res = super(ProductProduct, unlink_products).unlink()
			  # delete templates after calling super, as deleting template could lead to deleting
			  # products due to ondelete='cascade'
			  unlink_templates.unlink()
			  return res

			*/
		})

	pool.ProductProduct().Methods().Copy().Extend("",
		func(rs pool.ProductProductSet, overrides *pool.ProductProductData, fieldsToUnset ...models.FieldNamer) pool.ProductProductSet {
			//@api.multi
			/*def copy(self, default=None):
			  # TDE FIXME: clean context / variant brol
			  if default is None:
			      default = {}
			  if self._context.get('variant'):
			      # if we copy a variant or create one, we keep the same template
			      default['product_tmpl_id'] = self.product_tmpl_id.id
			  elif 'name' not in default:
			      default['name'] = self.name

			  return super(ProductProduct, self).copy(default=default)

			*/
		})

	pool.ProductProduct().Methods().Search().Extend("",
		func(rs pool.ProductProductSet, cond pool.ProductProductCondition) pool.ProductProductSet {
			//@api.model
			/*def search(self, args, offset=0, limit=None, order=None, count=False):
			  # TDE FIXME: strange
			  if self._context.get('search_default_categ_id'):
			      args.append((('categ_id', 'child_of', self._context['search_default_categ_id'])))
			  return super(ProductProduct, self).search(args, offset=offset, limit=limit, order=order, count=count)

			*/
		})

	pool.ProductProduct().Methods().NameGet().Extend("",
		func(rs pool.ProductProductSet) string {
			//@api.multi
			/*def _name_get(d):
			      name = d.get('name', '')
			      code = self._context.get('display_default_code', True) and d.get('default_code', False) or False
			      if code:
			          name = '[%s] %s' % (code,name)
			      return (d['id'], name)

			  partner_id = self._context.get('partner_id')
			  if partner_id:
			      partner_ids = [partner_id, self.env['res.partner'].browse(partner_id).commercial_partner_id.id]
			  else:
			      partner_ids = []

			  # all user don't have access to seller and partner
			  # check access and use superuser
			  self.check_access_rights("read")
			  self.check_access_rule("read")

			  result = []
			  for product in self.sudo():
			      # display only the attributes with multiple possible values on the template
			      variable_attributes = product.attribute_line_ids.filtered(lambda l: len(l.value_ids) > 1).mapped('attribute_id')
			      variant = product.attribute_value_ids._variant_name(variable_attributes)

			      name = variant and "%s (%s)" % (product.name, variant) or product.name
			      sellers = []
			      if partner_ids:
			          sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and (x.product_id == product)]
			          if not sellers:
			              sellers = [x for x in product.seller_ids if (x.name.id in partner_ids) and not x.product_id]
			      if sellers:
			          for s in sellers:
			              seller_variant = s.product_name and (
			                  variant and "%s (%s)" % (s.product_name, variant) or s.product_name
			                  ) or False
			              mydict = {
			                        'id': product.id,
			                        'name': seller_variant or name,
			                        'default_code': s.product_code or product.default_code,
			                        }
			              temp = _name_get(mydict)
			              if temp not in result:
			                  result.append(temp)
			      else:
			          mydict = {
			                    'id': product.id,
			                    'name': name,
			                    'default_code': product.default_code,
			                    }
			          result.append(_name_get(mydict))
			  return result

			*/
		})

	pool.ProductProduct().Methods().SearchByName().Extend("",
		func(rs pool.ProductProductSet, name string, op operator.Operator, additionalCond pool.ProductProductCondition, limit int) pool.ProductProductSet {
			//@api.model
			/*def name_search(self, name, args=None, operator='ilike', limit=100):
			  if not args:
			      args = []
			  if name:
			      # Be sure name_search is symetric to name_get
			      category_names = name.split(' / ')
			      parents = list(category_names)
			      child = parents.pop()
			      domain = [('name', operator, child)]
			      if parents:
			          names_ids = self.name_search(' / '.join(parents), args=args, operator='ilike', limit=limit)
			          category_ids = [name_id[0] for name_id in names_ids]
			          if operator in expression.NEGATIVE_TERM_OPERATORS:
			              categories = self.search([('id', 'not in', category_ids)])
			              domain = expression.OR([[('parent_id', 'in', categories.ids)], domain])
			          else:
			              domain = expression.AND([[('parent_id', 'in', category_ids)], domain])
			          for i in range(1, len(category_names)):
			              domain = [[('name', operator, ' / '.join(category_names[-1 - i:]))], domain]
			              if operator in expression.NEGATIVE_TERM_OPERATORS:
			                  domain = expression.AND(domain)
			              else:
			                  domain = expression.OR(domain)
			      categories = self.search(expression.AND([domain, args]), limit=limit)
			  else:
			      categories = self.search(args, limit=limit)
			  return categories.name_get()


			*/
		})

	pool.ProductProduct().Methods().OpenProductTemplate().DeclareMethod(
		`OpenProductTemplate`,
		func(rs pool.ProductProductSet) {
			//@api.multi
			/*def open_product_template(self):
			  """ Utility method used to add an "Open Template" button in product views """
			  self.ensure_one()
			  return {'type': 'ir.actions.act_window',
			          'res_model': 'product.template',
			          'view_mode': 'form',
			          'res_id': self.product_tmpl_id.id,
			          'target': 'new'}

			*/
		})

	pool.ProductProduct().Methods().SelectSeller().DeclareMethod(
		`SelectSeller`,
		func(rs pool.ProductProductSet, args struct {
			PartnerId interface{}
			Quantity  interface{}
			Date      interface{}
			UomId     interface{}
		}) {
			//@api.multi
			/*
				def _select_seller(self, partner_id=False, quantity=0.0, date=None, uom_id=False):
					self.ensure_one()
					if date is None:
					  	date = fields.Date.today()
					res = self.env['product.supplierinfo']
					for seller in self.seller_ids:
						# Set quantity in UoM of seller
						quantity_uom_seller = quantity
						if quantity_uom_seller and uom_id and uom_id != seller.product_uom:
							quantity_uom_seller = uom_id._compute_quantity(quantity_uom_seller, seller.product_uom)

						if seller.date_start and seller.date_start > date:
							continue
						if seller.date_end and seller.date_end < date:
							continue
						if partner_id and seller.name not in [partner_id, partner_id.parent_id]:
							continue
						if quantity_uom_seller < seller.min_qty:
							continue
						if seller.product_id and seller.product_id != self:
							continue

						res |= seller
						break
					return res*/
		})

	pool.ProductProduct().Methods().PriceCompute().DeclareMethod(
		`PriceCompute`,
		func(rs pool.ProductProductSet, args struct {
			PriceType interface{}
			Uom       interface{}
			Currency  interface{}
			Company   interface{}
		}) {
			//@api.multi
			/*def price_compute(self, price_type, uom=False, currency=False, company=False):
			  # TDE FIXME: delegate to template or not ? fields are reencoded here ...
			  # compatibility about context keys used a bit everywhere in the code
			  if not uom and self._context.get('uom'):
			      uom = self.env['product.uom'].browse(self._context['uom'])
			  if not currency and self._context.get('currency'):
			      currency = self.env['res.currency'].browse(self._context['currency'])

			  products = self
			  if price_type == 'standard_price':
			      # standard_price field can only be seen by users in base.group_user
			      # Thus, in order to compute the sale price from the cost for users not in this group
			      # We fetch the standard price as the superuser
			      products = self.with_context(force_company=company and company.id or self._context.get('force_company', self.env.user.company_id.id)).sudo()

			  prices = dict.fromkeys(self.ids, 0.0)
			  for product in products:
			      prices[product.id] = product[price_type] or 0.0
			      if price_type == 'list_price':
			          prices[product.id] += product.price_extra

			      if uom:
			          prices[product.id] = product.uom_id._compute_price(prices[product.id], uom)

			      # Convert from current user company currency to asked one
			      # This is right cause a field cannot be in more than one currency
			      if currency:
			          prices[product.id] = product.currency_id.compute(prices[product.id], currency)

			  return prices

			*/
		})

	pool.ProductProduct().Methods().DefineStandardPrice().DeclareMethod(
		`DefineStandardPrice stores the standard price change in order to be able to retrieve the cost of a product for
		a given date`,
		func(rs pool.ProductProductSet, args struct {
			Value interface{}
		}) {
			//@api.multi
			/*def _set_standard_price(self, value):
			  ''' Store the standard price change in order to be able to retrieve the cost of a product for a given date'''
			  PriceHistory = self.env['product.price.history']
			  for product in self:
			      PriceHistory.create({
			          'product_id': product.id,
			          'cost': value,
			          'company_id': self._context.get('force_company', self.env.user.company_id.id),
			      })

			*/
		})

	pool.ProductProduct().Methods().GetHistoryPrice().DeclareMethod(
		`GetHistoryPrice`,
		func(rs pool.ProductProductSet, args struct {
			CompanyId interface{}
			Date      interface{}
		}) {
			//@api.multi
			/*def get_history_price(self, company_id, date=None):
			  history = self.env['product.price.history'].search([
			      ('company_id', '=', company_id),
			      ('product_id', 'in', self.ids),
			      ('datetime', '<=', date or */
		})

	pool.ProductProduct().Methods().NeedProcurement().DeclareMethod(
		`NeedProcurement`,
		func(rs pool.ProductProductSet) bool {
			// When sale/product is installed alone, there is no need to create procurements. Only
			// sale_stock and sale_service need procurements
			return false
		})

	pool.ProductPackaging().DeclareModel()
	pool.ProductPackaging().SetDefaultOrder("Sequence")

	pool.ProductPackaging().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Packaging Type", Required: true},
		"Sequence": models.IntegerField{Default: models.DefaultValue(1),
			Help: "The first in the sequence is the default one."},
		"ProductTmpl": models.Many2OneField{String: "Product", RelationModel: pool.ProductTemplate()},
		"Qty": models.FloatField{String: "Quantity per Package",
			Help: "The total number of products you can have per pallet or box."},
	})

	pool.ProductSupplierinfo().DeclareModel()
	pool.ProductSupplierinfo().SetDefaultOrder("Sequence", "MinQty DESC", "Price")

	pool.ProductSupplierinfo().AddFields(map[string]models.FieldDefinition{
		"Name": models.Many2OneField{String: "Vendor", RelationModel: pool.Partner(),
			Filter: pool.Partner().Supplier().Equals(true), OnDelete: models.Cascade, Required: true,
			Help: "Vendor of this product"},
		"ProductName": models.CharField{String: "Vendor Product Name",
			Help: `This vendor's product name will be used when printing a request for quotation.
Keep empty to use the internal one.`},
		"ProductCode": models.CharField{String: "Vendor Product Code",
			Help: `This vendor's product code will be used when printing a request for quotation.
Keep empty to use the internal one.`},
		"Sequence": models.IntegerField{Default: models.DefaultValue(1),
			Help: "Assigns the priority to the list of product vendor."},
		"ProductUom": models.Many2OneField{String: "Vendor Unit of Measure", RelationModel: pool.ProductUom(), /* readonly=true */
			Related: "ProductTmpl.UomPo", Help: "This comes from the product form."},
		"MinQty": models.FloatField{String: "Minimal Quantity", Default: models.DefaultValue(0), Required: true,
			Help: `The minimal quantity to purchase from this vendor, expressed in the vendor Product Unit of Measure if any,
or in the default unit of measure of the product otherwise.`},
		"Price": models.FloatField{Default: models.DefaultValue(0), Digits: decimalPrecision.GetPrecision("Product Price"),
			Required: true, Help: "The price to purchase a product"},
		"Company": models.Many2OneField{RelationModel: pool.Company(), Default: func(env models.Environment, vals models.FieldMap) interface{} {
			return pool.User().NewSet(env).CurrentUser().Company()
		}, Index: true},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(), Default: func(env models.Environment, vals models.FieldMap) interface{} {
			return pool.User().NewSet(env).CurrentUser().Company().Currency()
		}, Required: true},
		"DateStart": models.DateField{String: "Start Date", Help: "Start date for this vendor price"},
		"DateEnd":   models.DateField{String: "End Date", Help: "End date for this vendor price"},
		"Product": models.Many2OneField{String: "Product Variant", RelationModel: pool.ProductProduct(),
			Help: "When this field is filled in, the vendor data will only apply to the variant."},
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: pool.ProductTemplate(),
			Index: true, OnDelete: models.Cascade},
		"Delay": models.IntegerField{String: "Delivery Lead Time", Default: models.DefaultValue(1), Required: true,
			Help: `Lead time in days between the confirmation of the purchase order and the receipt of the
products in your warehouse. Used by the scheduler for automatic computation of the purchase order planning.`},
	})

}
