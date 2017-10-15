package product

import (
	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProductPricelist().DeclareModel()
	pool.ProductPricelist().SetDefaultOrder("Sequence ASC", "ID DESC")

	pool.ProductPricelist().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Pricelist Name", Required: true, Translate: true},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If unchecked, it will allow you to hide the pricelist without removing it."},
		"Items": models.One2ManyField{String: "Pricelist Items", RelationModel: pool.ProductPricelistItem(),
			ReverseFK: "Pricelist", JSON: "item_ids", NoCopy: false,
			Default: func(models.Environment, models.FieldMap) interface{} {
				/*_get_default_item_ids(self):
				  ProductPricelistItem = self.env['product.pricelist.item']
				  vals = ProductPricelistItem.default_get(ProductPricelistItem._*/
				return 0
			}},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(),
			Default: func(models.Environment, models.FieldMap) interface{} {
				/*_get_default_currency_id(self):
				  return self.env.user.company_id.currency_id.id

				*/
				return 0
			}, Required: true},
		"Company":       models.Many2OneField{RelationModel: pool.Company()},
		"Sequence":      models.IntegerField{Default: models.DefaultValue(16)},
		"CountryGroups": models.Many2ManyField{RelationModel: pool.CountryGroup(), JSON: "country_group_ids"},
	})

	pool.ProductPricelist().Methods().NameGet().Extend("",
		func(rs pool.ProductPricelistSet) string {
			//@api.multi
			/*def name_get(self):
			  return [(pricelist.id, '%s (%s)' % (pricelist.name, pricelist.currency_id.name)) for pricelist in self]

			*/
		})

	pool.ProductPricelist().Methods().SearchByName().Extend("",
		func(rs pool.ProductPricelistSet, name string, op operator.Operator, additionalCondition pool.ProductPricelistCondition, limit int) pool.ProductPricelistSet {
			//@api.model
			/*def name_search(self, name, args=None, operator='ilike', limit=100):
			  if name and operator == '=' and not args:
			      # search on the name of the pricelist and its currency, opposite of name_get(),
			      # Used by the magic context filter in the product search view.
			      query_args = {'name': name, 'limit': limit, 'lang': self._context.get('lang') or 'en_US'}
			      query = """SELECT p.id
			                 FROM ((
			                          SELECT pr.id, pr.name
			                          FROM product_pricelist pr JOIN
			                               res_currency cur ON
			                                   (pr.currency_id = cur.id)
			                          WHERE pr.name || ' (' || cur.name || ')' = %(name)s
			                      )
			                      UNION (
			                          SELECT tr.res_id as id, tr.value as name
			                          FROM ir_translation tr JOIN
			                               product_pricelist pr ON (
			                                  pr.id = tr.res_id AND
			                                  tr.type = 'model' AND
			                                  tr.name = 'product.pricelist,name' AND
			                                  tr.lang = %(lang)s
			                               ) JOIN
			                               res_currency cur ON
			                                   (pr.currency_id = cur.id)
			                          WHERE tr.value || ' (' || cur.name || ')' = %(name)s
			                      )
			                  ) p
			                 ORDER BY p.name"""
			      if limit:
			          query += " LIMIT %(limit)s"
			      self._cr.execute(query, query_args)
			      ids = [r[0] for r in self._cr.fetchall()]
			      # regular search() to apply ACLs - may limit results below limit in some cases
			      pricelists = self.search([('id', 'in', ids)], limit=limit)
			      if pricelists:
			          return pricelists.name_get()
			  return super(Pricelist, self).name_search(name, args, operator=operator, limit=limit)

			*/
		})

	pool.ProductPricelist().Methods().ComputePriceRule().DeclareMethod(
		`ComputePriceRule is the low-level method computing the price of the given product according to this
		price list. Price depends on quantity, partner and date, and is given for the uom.`,
		func(rs pool.ProductPricelistSet, product pool.ProductProductSet, quantity float64, partner pool.PartnerSet,
			Date dates.Date, uom pool.ProductUomSet) (float64, pool.ProductPricelistItemSet) {

			rs.EnsureOne()
			//@api.multi
			/*def _compute_price_rule(self, products_qty_partner, date=False, uom_id=False):
						  """ Low-level method - Mono pricelist, multi products
						  Returns: dict{product_id: (price, suitable_rule) for the given pricelist}

						  If date in context: Date of the pricelist (%Y-%m-%d)

						      :param products_qty_partner: list of typles products, quantity, partner
						      :param datetime date: validity date
						      :param ID uom_id: intermediate unit of measure
						  """
						  self.ensure_one()
						  if not date:
						      date = self._context.get('date') or fields.Date.today()
			        if not uom_id and self._context.get('uom'):
			            uom_id = self._context['uom']
			        if uom_id:
			            # rebrowse with uom if given
			            products = [item[0].with_context(uom=uom_id) for item in products_qty_partner]
			            products_qty_partner = [(products[index], data_struct[1], data_struct[2]) for index, data_struct in enumerate(products_qty_partner)]
			        else:
			            products = [item[0] for item in products_qty_partner]

			        if not products:
			            return {}

			        categ_ids = {}
			        for p in products:
			            categ = p.categ_id
			            while categ:
			                categ_ids[categ.id] = True
			                categ = categ.parent_id
			        categ_ids = categ_ids.keys()

			        is_product_template = products[0]._name == "product.template"
			        if is_product_template:
			            prod_tmpl_ids = [tmpl.id for tmpl in products]
			            # all variants of all products
			            prod_ids = [p.id for p in
			                        list(chain.from_iterable([t.product_variant_ids for t in products]))]
			        else:
			            prod_ids = [product.id for product in products]
			            prod_tmpl_ids = [product.product_tmpl_id.id for product in products]

			        # Load all rules
			        self._cr.execute(
			            'SELECT item.id '
			            'FROM product_pricelist_item AS item '
			            'LEFT JOIN product_category AS categ '
			            'ON item.categ_id = categ.id '
			            'WHERE (item.product_tmpl_id IS NULL OR item.product_tmpl_id = any(%s))'
			            'AND (item.product_id IS NULL OR item.product_id = any(%s))'
			            'AND (item.categ_id IS NULL OR item.categ_id = any(%s)) '
			            'AND (item.pricelist_id = %s) '
			            'AND (item.date_start IS NULL OR item.date_start<=%s) '
			            'AND (item.date_end IS NULL OR item.date_end>=%s)'
			            'ORDER BY item.applied_on, item.min_quantity desc, categ.parent_left desc',
			            (prod_tmpl_ids, prod_ids, categ_ids, self.id, date, date))

			        item_ids = [x[0] for x in self._cr.fetchall()]
			        items = self.env['product.pricelist.item'].browse(item_ids)
			        results = {}
			        for product, qty, partner in products_qty_partner:
			            results[product.id] = 0.0
			            suitable_rule = False

			            # Final unit price is computed according to `qty` in the `qty_uom_id` UoM.
			            # An intermediary unit price may be computed according to a different UoM, in
			            # which case the price_uom_id contains that UoM.
			            # The final price will be converted to match `qty_uom_id`.
			            qty_uom_id = self._context.get('uom') or product.uom_id.id
			            price_uom_id = product.uom_id.id
			            qty_in_product_uom = qty
			            if qty_uom_id != product.uom_id.id:
			                try:
			                    qty_in_product_uom = self.env['product.uom'].browse([self._context['uom']])._compute_quantity(qty, product.uom_id)
			                except UserError:
			                    # Ignored - incompatible UoM in context, use default product UoM
			                    pass

			            # if Public user try to access standard price from website sale, need to call price_compute.
			            # TDE SURPRISE: product can actually be a template
			            price = product.price_compute('list_price')[product.id]

			            price_uom = self.env['product.uom'].browse([qty_uom_id])
			            for rule in items:
			                if rule.min_quantity and qty_in_product_uom < rule.min_quantity:
			                    continue
			                if is_product_template:
			                    if rule.product_tmpl_id and product.id != rule.product_tmpl_id.id:
			                        continue
			                    if rule.product_id and not (product.product_variant_count == 1 and product.product_variant_id.id == rule.product_id.id):
			                        # product rule acceptable on template if has only one variant
			                        continue
			                else:
			                    if rule.product_tmpl_id and product.product_tmpl_id.id != rule.product_tmpl_id.id:
			                        continue
			                    if rule.product_id and product.id != rule.product_id.id:
			                        continue

			                if rule.categ_id:
			                    cat = product.categ_id
			                    while cat:
			                        if cat.id == rule.categ_id.id:
			                            break
			                        cat = cat.parent_id
			                    if not cat:
			                        continue

			                if rule.base == 'pricelist' and rule.base_pricelist_id:
			                    price_tmp = rule.base_pricelist_id._compute_price_rule([(product, qty, partner)])[product.id][0]  # TDE: 0 = price, 1 = rule
			                    price = rule.base_pricelist_id.currency_id.compute(price_tmp, self.currency_id, round=False)
			                else:
			                    # if base option is public price take sale price else cost price of product
			                    # price_compute returns the price in the context UoM, i.e. qty_uom_id
			                    price = product.price_compute(rule.base)[product.id]

			                convert_to_price_uom = (lambda price: product.uom_id._compute_price(price, price_uom))

			                if price is not False:
			                    if rule.compute_price == 'fixed':
			                        price = convert_to_price_uom(rule.fixed_price)
			                    elif rule.compute_price == 'percentage':
			                        price = (price - (price * (rule.percent_price / 100))) or 0.0
			                    else:
			                        # complete formula
			                        price_limit = price
			                        price = (price - (price * (rule.price_discount / 100))) or 0.0
			                        if rule.price_round:
			                            price = tools.float_round(price, precision_rounding=rule.price_round)

			                        if rule.price_surcharge:
			                            price_surcharge = convert_to_price_uom(rule.price_surcharge)
			                            price += price_surcharge

			                        if rule.price_min_margin:
			                            price_min_margin = convert_to_price_uom(rule.price_min_margin)
			                            price = max(price, price_limit + price_min_margin)

			                        if rule.price_max_margin:
			                            price_max_margin = convert_to_price_uom(rule.price_max_margin)
			                            price = min(price, price_limit + price_max_margin)
			                    suitable_rule = rule
			                break
			            # Final price conversion into pricelist currency
			            if suitable_rule and suitable_rule.compute_price != 'fixed' and suitable_rule.base != 'pricelist':
			                price = product.currency_id.compute(price, self.currency_id, round=False)

			            results[product.id] = (price, suitable_rule and suitable_rule.id or False)

			        return results*/
		})

	pool.ProductPricelist().Methods().GetProductPrice().DeclareMethod(
		`GetProductPrice returns the price of the given product in the given quantity for the given partner, at
		the given date and in the given UoM according to this price list.`,
		func(rs pool.ProductPricelistSet, product pool.ProductProductSet, quantity float64, partner pool.PartnerSet,
			date dates.Date, uom pool.ProductUomSet) float64 {

			rs.EnsureOne()
			price, _ := rs.ComputePriceRule(product, quantity, partner, date, uom)
			return price
		})

	pool.ProductPricelist().Methods().GetProductPriceRule().DeclareMethod(
		`GetProductPriceRule returns the applicable price list rule for the given product in the given quantity
		for the given partner, at the given date and in the given UoM according to this price list.`,
		func(rs pool.ProductPricelistSet, product pool.ProductProductSet, quantity float64, partner pool.PartnerSet,
			date dates.Date, uom pool.ProductUomSet) pool.ProductPricelistItemSet {

			rs.EnsureOne()
			_, rule := rs.ComputePriceRule(product, quantity, partner, date, uom)
			return rule
		})

	pool.ProductPricelist().Methods().GetPartnerPricelist().DeclareMethod(
		`GetPartnerPricelist rtrieve the applicable pricelist for the given partner in the given company.`,
		func(rs pool.ProductPricelistSet, partner pool.PartnerSet, company pool.CompanySet) pool.ProductPricelistSet {
			/*def _get_partner_pricelist(self, partner_id, company_id=None):
			  """ Retrieve the applicable pricelist for a given partner in a given company.

			      :param company_id: if passed, used for looking up properties,
			       instead of current user's company
			  """
			  Partner = self.env['res.partner']
			  Property = self.env['ir.property'].with_context(force_company=company_id or self.env.user.company_id.id)

			  p = Partner.browse(partner_id)
			  pl = Property.get('property_product_pricelist', Partner._name, '%s,%s' % (Partner._name, p.id))
			  if pl:
			      pl = pl[0].id

			  if not pl:
			      if p.country_id.code:
			          pls = self.env['product.pricelist'].search([('country_group_ids.country_ids.code', '=', p.country_id.code)], limit=1)
			          pl = pls and pls[0].id

			  if not pl:
			      # search pl where no country
			      pls = self.env['product.pricelist'].search([('country_group_ids', '=', False)], limit=1)
			      pl = pls and pls[0].id

			  if not pl:
			      prop = Property.get('property_product_pricelist', 'res.partner')
			      pl = prop and prop[0].id

			  if not pl:
			      pls = self.env['product.pricelist'].search([], limit=1)
			      pl = pls and pls[0].id

			  return pl


			*/
		})

	pool.CountryGroup().AddFields(map[string]models.FieldDefinition{
		"Pricelists": models.Many2ManyField{String: "Pricelists", RelationModel: pool.ProductPricelist(),
			JSON: "pricelist_ids"},
	})

	pool.ProductPricelistItem().DeclareModel()
	pool.ProductPricelistItem().SetDefaultOrder("AppliedOn", "MinQuantity DESC", "Categ DESC", "ID")

	pool.ProductPricelistItem().AddFields(map[string]models.FieldDefinition{
		"ProductTmpl": models.Many2OneField{String: "Product Template", RelationModel: pool.ProductTemplate(),
			OnDelete: models.Cascade,
			Help:     "Specify a template if this rule only applies to one product template. Keep empty otherwise."},
		"Product": models.Many2OneField{RelationModel: pool.ProductProduct(), OnDelete: models.Cascade,
			Help: "Specify a product if this rule only applies to one product. Keep empty otherwise."},
		"Categ": models.Many2OneField{String: "Product Category", RelationModel: pool.ProductCategory(),
			OnDelete: models.Cascade,
			Help: `Specify a product category if this rule only applies to products belonging to this category or 
its children categories. Keep empty otherwise.`},
		"MinQuantity": models.IntegerField{Default: models.DefaultValue(1),
			Help: `For the rule to apply, bought/sold quantity must be greater
than or equal to the minimum quantity specified in this field.
Expressed in the default unit of measure of the product.`},
		"AppliedOn": models.SelectionField{String: "Apply On", Selection: types.Selection{
			"3_global":           "Global",
			"2_product_category": "Product Category",
			"1_product":          "Product",
			"0_product_variant":  "Product Variant",
		}, Default: models.DefaultValue("3_global"), Required: true,
			Help: "Pricelist Item applicable on selected option"},
		"Sequence": models.IntegerField{Default: models.DefaultValue(5), Required: true,
			Help: `Gives the order in which the pricelist items will be checked. The evaluation gives highest priority
to lowest sequence and stops as soon as a matching item is found.`},
		"Base": models.SelectionField{String: "Based on", Selection: types.Selection{
			"list_price":     "Public Price",
			"standard_price": "Cost",
			"pricelist":      "Other Pricelist",
		}, Default: models.DefaultValue("list_price"), Required: true,
			Help: `Base price for computation.
- Public Price: The base price will be the Sale/public Price.
- Cost Price : The base price will be the cost price.
- Other Pricelist : Computation of the base price based on another Pricelist.`},
		"BasePricelist": models.Many2OneField{String: "Other Pricelist", RelationModel: pool.ProductPricelist()},
		"Pricelist": models.Many2OneField{RelationModel: pool.ProductPricelist(), Index: true,
			OnDelete: models.Cascade},
		"PriceSurcharge": models.FloatField{Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: "Specify the fixed amount to add or substract(if negative) to the amount calculated with the discount."},
		"PriceDiscount": models.FloatField{Default: models.DefaultValue(0),
			Digits: nbutils.Digits{Precision: 16, Scale: 2}},
		"PriceRound": models.FloatField{Digits: decimalPrecision.GetPrecision("Product Price"),
			Help: `Sets the price so that it is a multiple of this value.
Rounding is applied after the discount and before the surcharge.
To have prices that end in 9.99, set rounding 10, surcharge -0.01`},
		"PriceMinMargin": models.FloatField{String: "Min. Price Margin",
			Digits: decimalPrecision.GetPrecision("Product Price"),
			Help:   "Specify the minimum amount of margin over the base price."},
		"PriceMaxMargin": models.FloatField{String: "Max. Price Margin",
			Digits: decimalPrecision.GetPrecision("Product Price"),
			Help:   "Specify the maximum amount of margin over the base price."},
		"Company": models.Many2OneField{RelationModel: pool.Company(), /* readonly=true */
			Related: "Pricelist.Company", Stored: true},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(), /* readonly=true */
			Related: "Pricelist.Currency", Stored: true},
		"DateStart": models.DateField{String: "Start Date", Help: "Starting date for the pricelist item validation"},
		"DateEnd":   models.DateField{String: "End Date", Help: "Ending valid for the pricelist item validation"},
		"ComputePrice": models.SelectionField{Selection: types.Selection{
			"fixed":      "Fix Price",
			"percentage": "Percentage (discount)",
			"formula":    "Formula",
		}, Index: true, Default: models.DefaultValue("fixed")},
		"FixedPrice":   models.FloatField{String: "Fixed Price", Digits: decimalPrecision.GetPrecision("Product Price")},
		"PercentPrice": models.FloatField{String: "Percentage Price"},
		"Name": models.CharField{Compute: pool.ProductPricelistItem().Methods().GetPricelistItemNamePrice(),
			Help: "Explicit rule name for this pricelist line."},
		"Price": models.CharField{Compute: pool.ProductPricelistItem().Methods().GetPricelistItemNamePrice(),
			Help: "Explicit rule name for this pricelist line."},
	})

	pool.ProductPricelistItem().Methods().CheckRecursion().DeclareMethod(
		`CheckRecursion`,
		func(rs pool.ProductPricelistItemSet) {
			//@api.constrains('base_pricelist_id','pricelist_id','base')
			/*def _check_recursion(self):
			  if any(item.base == 'pricelist' and item.pricelist_id and item.pricelist_id == item.base_pricelist_id for item in self):
			      raise ValidationError(_('Error! You cannot assign the Main Pricelist as Other Pricelist in PriceList Item!'))
			  return True

			*/
		})

	pool.ProductPricelistItem().Methods().CheckMargin().DeclareMethod(
		`CheckMargin`,
		func(rs pool.ProductPricelistItemSet) {
			//@api.constrains('price_min_margin','price_max_margin')
			/*def _check_margin(self):
			  if any(item.price_min_margin > item.price_max_margin for item in self):
			      raise ValidationError(_('Error! The minimum margin should be lower than the maximum margin.'))
			  return True

			*/
		})

	pool.ProductPricelistItem().Methods().GetPricelistItemNamePrice().DeclareMethod(
		`GetPricelistItemNamePrice`,
		func(rs pool.ProductPricelistItemSet) (*pool.ProductPricelistItemSet, models.FieldNamer) {
			/*def _get_pricelist_item_name_price(self):
			  if self.categ_id:
			      self.name = _("Category: %s") % (self.categ_id.name)
			  elif self.product_tmpl_id:
			      self.name = self.product_tmpl_id.name
			  elif self.product_id:
			      self.name = self.product_id.display_name.replace('[%s]' % self.product_id.code, '')
			  else:
			      self.name = _("All Products")

			  if self.compute_price == 'fixed':
			      self.price = ("%s %s") % (self.fixed_price, self.pricelist_id.currency_id.name)
			  elif self.compute_price == 'percentage':
			      self.price = _("%s %% discount") % (self.percent_price)
			  else:
			      self.price = _("%s %% discount and %s surcharge") % (abs(self.price_discount), self.price_surcharge)

			*/
		})

	pool.ProductPricelistItem().Methods().OnchangeAppliedOn().DeclareMethod(
		`OnchangeAppliedOn`,
		func(rs pool.ProductPricelistItemSet) (*pool.ProductPricelistItemSet, models.FieldNamer) {
			//@api.onchange('applied_on')
			/*def _onchange_applied_on(self):
			  if self.applied_on != '0_product_variant':
			      self.product_id = False
			  if self.applied_on != '1_product':
			      self.product_tmpl_id = False
			  if self.applied_on != '2_product_category':
			      self.categ_id = False

			*/
		})

	pool.ProductPricelistItem().Methods().OnchangeComputePrice().DeclareMethod(
		`OnchangeComputePrice`,
		func(rs pool.ProductPricelistItemSet) (*pool.ProductPricelistItemSet, models.FieldNamer) {
			//@api.onchange('compute_price')
			/*def _onchange_compute_price(self):
			  if self.compute_price != 'fixed':
			      self.fixed_price = 0.0
			  if self.compute_price != 'percentage':
			      self.percent_price = 0.0
			  if self.compute_price != 'formula':
			      self.update({
			          'price_discount': 0.0,
			          'price_surcharge': 0.0,
			          'price_round': 0.0,
			          'price_min_margin': 0.0,
			          'price_max_margin': 0.0,
			      })

			*/
		})

}
