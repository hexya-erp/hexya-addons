// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.ProductCategory().AddFields(map[string]models.FieldDefinition{
		"PropertyAccountIncomeCateg": models.Many2OneField{String: "Income Account",
			RelationModel: h.AccountAccount(), /*, CompanyDependent : true*/
			Filter:        q.AccountAccount().Deprecated().Equals(false),
			Help:          "This account will be used for invoices to value sales."},
		"PropertyAccountExpenseCateg": models.Many2OneField{String: "Expense Account",
			RelationModel: h.AccountAccount(), /*, CompanyDependent : true*/
			Filter:        q.AccountAccount().Deprecated().Equals(false),
			Help:          "This account will be used for invoices to value expenses."},
	})

	h.ProductTemplate().AddFields(map[string]models.FieldDefinition{
		"Taxes": models.Many2ManyField{String: "Customer Taxes", RelationModel: h.AccountTax(),
			JSON: "taxes_id", Filter: q.AccountTax().TypeTaxUse().Equals("sale")},
		"SupplierTaxes": models.Many2ManyField{String: "Vendor Taxes", RelationModel: h.AccountTax(),
			JSON: "supplier_taxes_id", Filter: q.AccountTax().TypeTaxUse().Equals("purchase")},
		"PropertyAccountIncome": models.Many2OneField{String: "Income Account", RelationModel: h.AccountAccount(),
			/*, CompanyDependent : true*/ Filter: q.AccountAccount().Deprecated().Equals(false),
			Help: `This account will be used for invoices instead of the default one
to value sales for the current product.`},
		"PropertyAccountExpense": models.Many2OneField{String: "Expense Account", RelationModel: h.AccountAccount(),
			/*, CompanyDependent : true*/ Filter: q.AccountAccount().Deprecated().Equals(false),
			Help: `This account will be used for invoices instead of the default one
to value expenses for the current product.`},
	})

	h.ProductTemplate().Methods().Write().Extend("",
		func(rs h.ProductTemplateSet, data *h.ProductTemplateData, fieldsToReset ...models.FieldNamer) bool {
			//@api.multi
			/*def write(self, vals):
			  #TODO: really? i don't see the reason we'd need that constraint..
			  check = self.ids and 'uom_po_id' in vals
			  if check:
			      self._cr.execute("SELECT id, uom_po_id FROM product_template WHERE id IN %s", [tuple(self.ids)])
			      uoms = dict(self._cr.fetchall())
			  res = super(ProductTemplate, self).write(vals)
			  if check:
			      self._cr.execute("SELECT id, uom_po_id FROM product_template WHERE id IN %s", [tuple(self.ids)])
			      if dict(self._cr.fetchall()) != uoms:
			          products = self.env['product.product'].search([('product_tmpl_id', 'in', self.ids)])
			          if self.env['account.move.line'].search_count([('product_id', 'in', products.ids)]):
			              raise UserError(_('You can not change the unit of measure of a product that has been already used in an account journal item. If you need to change the unit of measure, you may deactivate this product.'))
			  return res

			*/
			return rs.Super().Write(data, fieldsToReset...)
		})

	h.ProductTemplate().Methods().GetProductDirectAccounts().DeclareMethod(
		`GetProductDirectAccounts`,
		func(rs h.ProductTemplateSet) (h.AccountAccountSet, h.AccountAccountSet) {
			//@api.multi
			/*def _get_product_accounts(self):
			  return {
			      'income': self.property_account_income_id or self.categ_id.property_account_income_categ_id,
			      'expense': self.property_account_expense_id or self.categ_id.property_account_expense_categ_id
			  }

			*/
			return h.AccountAccount().NewSet(rs.Env()), h.AccountAccount().NewSet(rs.Env())
		})

	h.ProductTemplate().Methods().GetAssetAccounts().DeclareMethod(
		`GetAssetAccounts`,
		func(rs h.ProductTemplateSet) (h.AccountAccountSet, h.AccountAccountSet) {
			//@api.multi
			/*def _get_asset_accounts(self):
			  res = {}
			  res['stock_input'] = False
			  res['stock_output'] = False
			  return res

			*/
			return h.AccountAccount().NewSet(rs.Env()), h.AccountAccount().NewSet(rs.Env())
		})

	h.ProductTemplate().Methods().GetProductAccounts().DeclareMethod(
		`GetProductAccounts`,
		func(rs h.ProductTemplateSet, fiscalPos h.AccountFiscalPositionSet) (h.AccountAccountSet, h.AccountAccountSet) {
			//@api.multi
			/*def get_product_accounts(self, fiscal_pos=None):
			  accounts = self._get_product_accounts()
			  if not fiscal_pos:
			      fiscal_pos = self.env['account.fiscal.position']
			  return fiscal_pos.map_accounts(accounts)
			*/
			return h.AccountAccount().NewSet(rs.Env()), h.AccountAccount().NewSet(rs.Env())
		})

}
