// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

func init() {

	//h.SaleConfigSettings().AddFields(map[string]models.FieldDefinition{
	//	"CompanyId": models.Many2OneField{String: "Company", RelationModel: h.Company() /*['res.company']*/, Required: true, Default: func(env models.Environment) interface{} {
	//		/*lambda self: self.env.user.company_id*/
	//		return 0
	//	}},
	//	"SaleNote":                 models.TextField{String: "Default Terms and Conditions *" /*[related 'company_id.sale_note']*/},
	//	"GroupProductVariant":      models.SelectionField{ /*group_product_variant = fields.Selection([ (0, "No variants on products"), (1, 'Products can have several attributes, defining variants (Example: size, color,...)')*/ },
	//	"GroupSalePricelist":       models.BooleanField{String: "Use pricelists to adapt your price per customers", String: "Use pricelists to adapt your price per customers" /*["Use pricelists to adapt your price per customers"]*/ /*[ implied_group 'product.group_sale_pricelist']*/, Help: "Allows to manage different prices based on rules per category of customers.#~#~# Example: 10% for retailers, promotion of 5 EUR on this product, etc." /*[ promotion of 5 EUR on this product]*/ /*[ etc."""]*/},
	//	"GroupPricelistItem":       models.BooleanField{String: "Show pricelists to customers", String: "Show pricelists to customers" /*["Show pricelists to customers"]*/ /*[ implied_group 'product.group_pricelist_item']*/},
	//	"GroupProductPricelist":    models.BooleanField{String: "Show pricelists On Products", String: "Show pricelists On Products" /*["Show pricelists On Products"]*/ /*[ implied_group 'product.group_product_pricelist']*/},
	//	"GroupUom":                 models.SelectionField{ /*group_uom = fields.Selection([ (0, 'Products have only one unit of measure (easier)'), (1, 'Some products may be sold/purchased in different units of measure (advanced)')*/ },
	//	"GroupDiscountPerSoLine":   models.SelectionField{ /*group_discount_per_so_line = fields.Selection([ (0, 'No discount on sales order lines, global discount only'), (1, 'Allow discounts on sales order lines')*/ },
	//	"GroupDisplayIncoterm":     models.SelectionField{ /*group_display_incoterm = fields.Selection([ (0, 'No incoterm on reports'), (1, 'Show incoterms on sales orders and invoices')*/ },
	//	"ModuleSaleMargin":         models.SelectionField{ /*module_sale_margin = fields.Selection([ (0, 'Salespeople do not need to view margins when quoting'), (1, 'Display margins on quotations and sales orders')*/ },
	//	"GroupSaleLayout":          models.SelectionField{ /*group_sale_layout = fields.Selection([ (0, 'Do not personalize sales orders and invoice reports'), (1, 'Personalize the sales orders and invoice report with categories, subtotals and page-breaks')*/ },
	//	"GroupWarningSale":         models.SelectionField{ /*group_warning_sale = fields.Selection([ (0, 'All the products and the customers can be used in sales orders'), (1, 'An informative or blocking warning can be set on a product or a customer')*/ },
	//	"ModuleWebsiteQuote":       models.SelectionField{ /*module_website_quote = fields.Selection([ (0, 'Print quotes or send by email'), (1, 'Send quotations your customer can approve & pay online (advanced)')*/ },
	//	"GroupSaleDeliveryAddress": models.SelectionField{ /*group_sale_delivery_address = fields.Selection([ (0, "Invoicing and shipping addresses are always the same (Example: services companies)"), (1, 'Display 3 fields on sales orders: customer, invoice address, delivery address')*/ },
	//	"SalePricelistSetting":     models.SelectionField{ /*sale_pricelist_setting = fields.Selection([ ('fixed', 'A single sale price per product'), ('percentage', 'Specific prices per customer segment, currency, etc.'), ('formula', 'Advanced pricing based on formulas (discounts, margins, rounding)')*/ },
	//	"GroupShowPriceSubtotal":   models.BooleanField{String: "Show subtotal", String: "Show subtotal" /*["Show subtotal"]*/ /*[ implied_group 'sale.group_show_price_subtotal']*/ /*[ group 'base.group_portal]*/ /*[base.group_user]*/ /*[base.group_public']*/},
	//	"GroupShowPriceTotal":      models.BooleanField{String: "Show total", String: "Show total" /*["Show total"]*/ /*[ implied_group 'sale.group_show_price_total']*/ /*[ group 'base.group_portal]*/ /*[base.group_user]*/ /*[base.group_public']*/},
	//	"SaleShowTax": models.SelectionField{String: "Tax Display", Selection: types.Selection{
	//		"subtotal": "Show line subtotals without taxes (B2B)",
	//		"total":    "Show line subtotals with taxes included (B2C)",
	//	}, /*[]*/ /*["Tax Display"]*/ Default: models.DefaultValue("subtotal"), Required: true},
	//	"DefaultInvoicePolicy":     models.SelectionField{ /*default_invoice_policy = fields.Selection([ ('order', 'Invoice ordered quantities'), ('delivery', 'Invoice delivered quantities')*/ },
	//	"DepositProductIdSetting":  models.Many2OneField{String: "Deposit Product", RelationModel: h.ProductProduct() /*[ 'product.product']*/ /*['Deposit Product']*/ /*, Filter: "[('type'*/ /*[ ' ']*/ /*[ 'service')]"]*/, Help: "Default product used for payment advances"},
	//	"AutoDoneSetting":          models.SelectionField{ /*auto_done_setting = fields.Selection([ (0, "Allow to edit sales order from the 'Sales Order' menu (not from the Quotation menu)"), (1, "Never allow to modify a confirmed sales order")*/ },
	//	"ModuleSaleContract":       models.BooleanField{String: "Manage subscriptions and recurring invoicing" /*["Manage subscriptions and recurring invoicing"]*/},
	//	"ModuleWebsiteSaleDigital": models.BooleanField{String: "Sell digital products - provide downloadable content on your customer portal" /*["Sell digital products - provide downloadable content on your customer portal"]*/},
	//})
	//h.SaleConfigSettings().Methods().SetSaleDefaults().DeclareMethod(
	//	`SetSaleDefaults`,
	//	func(rs h.SaleConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_sale_defaults(self):
	//		  return self.env['ir.values'].sudo().set_default(
	//		      'sale.config.settings', 'sale_pricelist_setting', self.sale_pricelist_setting)
	//
	//		*/
	//	})
	//h.SaleConfigSettings().Methods().SetDepositProductIdDefaults().DeclareMethod(
	//	`SetDepositProductIdDefaults`,
	//	func(rs h.SaleConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_deposit_product_id_defaults(self):
	//		  return self.env['ir.values'].sudo().set_default(
	//		      'sale.config.settings', 'deposit_product_id_setting', self.deposit_product_id_setting.id)
	//
	//		*/
	//	})
	//h.SaleConfigSettings().Methods().SetAutoDoneDefaults().DeclareMethod(
	//	`SetAutoDoneDefaults`,
	//	func(rs h.SaleConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_auto_done_defaults(self):
	//		  return self.env['ir.values'].sudo().set_default(
	//		      'sale.config.settings', 'auto_done_setting', self.auto_done_setting)
	//
	//		*/
	//	})
	//h.SaleConfigSettings().Methods().OnchangeSalePrice().DeclareMethod(
	//	`OnchangeSalePrice`,
	//	func(rs h.SaleConfigSettingsSet) {
	//		//@api.onchange('sale_pricelist_setting')
	//		/*def _onchange_sale_price(self):
	//		  if self.sale_pricelist_setting == 'percentage':
	//		      self.update({
	//		          'group_product_pricelist': True,
	//		          'group_sale_pricelist': True,
	//		          'group_pricelist_item': False,
	//		      })
	//		  elif self.sale_pricelist_setting == 'formula':
	//		      self.update({
	//		          'group_product_pricelist': False,
	//		          'group_sale_pricelist': True,
	//		          'group_pricelist_item': True,
	//		      })
	//		  else:
	//		      self.update({
	//		          'group_product_pricelist': False,
	//		          'group_sale_pricelist': False,
	//		          'group_pricelist_item': False,
	//		      })
	//
	//		*/
	//	})
	//h.SaleConfigSettings().Methods().SetSaleTaxDefaults().DeclareMethod(
	//	`SetSaleTaxDefaults`,
	//	func(rs h.SaleConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_sale_tax_defaults(self):
	//		  return self.env['ir.values'].sudo().set_default(
	//		      'sale.config.settings', 'sale_show_tax', self.sale_show_tax)
	//
	//		*/
	//	})
	//h.SaleConfigSettings().Methods().OnchangeSaleTax().DeclareMethod(
	//	`OnchangeSaleTax`,
	//	func(rs h.SaleConfigSettingsSet) {
	//		//@api.onchange('sale_show_tax')
	//		/*def _onchange_sale_tax(self):
	//		  if self.sale_show_tax == "subtotal":
	//		      self.update({
	//		          'group_show_price_total': False,
	//		          'group_show_price_subtotal': True,
	//		      })
	//		  else:
	//		      self.update({
	//		          'group_show_price_total': True,
	//		          'group_show_price_subtotal': False,
	//		      })
	//		*/
	//	})

}
