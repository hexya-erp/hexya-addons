// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountAnalyticLine().AddFields(map[string]models.FieldDefinition{
		"Amount":                 models.FloatField{String: "Amount" /*[currency_field 'company_currency_id']*/},
		"ProductUom":             models.Many2OneField{String: "Unit of Measure", RelationModel: pool.ProductUom(), JSON: "product_uom_id" /*['product.uom']*/},
		"Product":                models.Many2OneField{String: "Product", RelationModel: pool.ProductProduct(), JSON: "product_id" /*['product.product']*/},
		"GeneralAccount":         models.Many2OneField{String: "Financial Account", RelationModel: pool.AccountAccount(), JSON: "general_account_id" /*['account.account']*/, OnDelete: models.Restrict /* readonly=true */, Related: "Move.Account", Stored: true /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/},
		"Move":                   models.Many2OneField{String: "Move Line", RelationModel: pool.AccountMoveLine(), JSON: "move_id" /*['account.move.line']*/, OnDelete: models.Cascade, Index: true},
		"Code":                   models.CharField{String: "Code" /*[size 8]*/},
		"Ref":                    models.CharField{String: "Ref" /*[string 'Ref.']*/},
		"CompanyCurrency":        models.Many2OneField{String: "CompanyCurrencyId", RelationModel: pool.Currency(), JSON: "company_currency_id" /*['res.currency']*/, Related: "Company.Currency" /* readonly=true */, Help: "Utility field to express amount currency"},
		"Currency":               models.Many2OneField{String: "Account Currency", RelationModel: pool.Currency(), JSON: "currency_id" /*['res.currency']*/, Related: "Move.Currency", Stored: true, Help: "The related account currency if not equal to the company one." /* readonly=true */},
		"AmountCurrency":         models.FloatField{String: "AmountCurrency", Related: "Move", Stored: true, Help: "The amount expressed in the related account currency if not equal to the company one." /*[ readonly True]*/},
		"AnalyticAmountCurrency": models.FloatField{String: "AnalyticAmountCurrency" /*[string 'Amount Currency']*/, Compute: pool.AccountAnalyticLine().Methods().GetAnalyticAmountCurrency(), Help: "The amount expressed in the related account currency if not equal to the company one." /*[ readonly True]*/},
		"Partner":                models.Many2OneField{String: "Partner", RelationModel: pool.Partner(), JSON: "partner_id" /*['res.partner']*/, Related: "Account.Partner", Stored: true /* readonly=true */},
	})
	pool.AccountAnalyticLine().Methods().GetAnalyticAmountCurrency().DeclareMethod(
		`GetAnalyticAmountCurrency`,
		func(rs pool.AccountAnalyticLineSet) {
			/*def _get_analytic_amount_currency(self):
			  for line in self:
			      line.analytic_amount_currency = abs(line.amount_currency) * copysign(1, line.amount)

			*/
		})
	pool.AccountAnalyticLine().Methods().OnChangeUnitAmount().DeclareMethod(
		`OnChangeUnitAmount`,
		func(rs pool.AccountAnalyticLineSet) {
			//@api.onchange('product_id','product_uom_id','unit_amount','currency_id')
			/*def on_change_unit_amount(self):
			  if not self.product_id:
			      return {}

			  result = 0.0
			  prod_accounts = self.product_id.product_tmpl_id._get_product_accounts()
			  unit = self.product_uom_id
			  account = prod_accounts['expense']
			  if not unit or self.product_id.uom_po_id.category_id.id != unit.category_id.id:
			      unit = self.product_id.uom_po_id

			  # Compute based on pricetype
			  amount_unit = self.product_id.price_compute('standard_price', uom=unit)[self.product_id.id]
			  amount = amount_unit * self.unit_amount or 0.0
			  result = round(amount, self.currency_id.decimal_places) * -1
			  self.amount = result
			  self.general_account_id = account
			  self.product_uom_id = unit

			*/
		})
	pool.AccountAnalyticLine().Methods().ViewHeaderGet().DeclareMethod(
		`ViewHeaderGet`,
		func(rs pool.AccountAnalyticLineSet, args struct {
			ViewId   interface{}
			ViewType interface{}
		}) {
			//@api.model
			/*def view_header_get(self, view_id, view_type):
			  context = (self._context or {})
			  header = False
			  if context.get('account_id', False):
			      analytic_account = self.env['account.analytic.account'].search([('id', '=', context['account_id'])], limit=1)
			      header = _('Entries: ') + (analytic_account.name or '')
			  return header
			*/
		})

}
