// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountAnalyticLine().AddFields(map[string]models.FieldDefinition{
		"ProductUom": models.Many2OneField{String: "Unit of Measure", RelationModel: pool.ProductUom(),
			OnChange: pool.AccountAnalyticLine().Methods().OnChangeUnitAmount()},
		"Product": models.Many2OneField{RelationModel: pool.ProductProduct(),
			OnChange: pool.AccountAnalyticLine().Methods().OnChangeUnitAmount()},
		"GeneralAccount": models.Many2OneField{String: "Financial Account", RelationModel: pool.AccountAccount(),
			OnDelete: models.Restrict /* readonly=true */, Related: "Move.Account",
			Filter: pool.AccountAccount().Deprecated().Equals(false)},
		"Move": models.Many2OneField{String: "Move Line", RelationModel: pool.AccountMoveLine(),
			JSON: "move_id", OnDelete: models.Cascade, Index: true},
		"Code": models.CharField{String: "Code", Size: 8},
		"Ref":  models.CharField{},
		"CompanyCurrency": models.Many2OneField{RelationModel: pool.Currency(),
			Related: "Company.Currency" /* readonly=true */, Help: "Utility field to express amount currency"},
		"AmountCurrency": models.FloatField{Related: "Move.AmountCurrency",
			Help: "The amount expressed in the related account currency if not equal to the company one." /* readonly=True */},
		"AnalyticAmountCurrency": models.FloatField{String: "Amount Currency",
			Compute: pool.AccountAnalyticLine().Methods().GetAnalyticAmountCurrency(), /*[ readonly True]*/
			Help:    "The amount expressed in the related account currency if not equal to the company one."},
	})

	pool.AccountAnalyticLine().Fields().Currency().
		SetString("Account Currency").
		SetRelated("Move.Currency").
		SetOnchange(pool.AccountAnalyticLine().Methods().OnChangeUnitAmount()).
		SetHelp("The related account currency if not equal to the company one.")

	pool.AccountAnalyticLine().Fields().Partner().SetRelated("Account.Partner")

	pool.AccountAnalyticLine().Fields().UnitAmount().SetOnchange(pool.AccountAnalyticLine().Methods().OnChangeUnitAmount())

	pool.AccountAnalyticLine().Methods().GetAnalyticAmountCurrency().DeclareMethod(
		`GetAnalyticAmountCurrency`,
		func(rs pool.AccountAnalyticLineSet) (*pool.AccountAnalyticAccountData, []models.FieldNamer) {
			/*def _get_analytic_amount_currency(self):
			  for line in self:
			      line.analytic_amount_currency = abs(line.amount_currency) * copysign(1, line.amount)

			*/
			return new(pool.AccountAnalyticAccountData), []models.FieldNamer{}
		})

	pool.AccountAnalyticLine().Methods().OnChangeUnitAmount().DeclareMethod(
		`OnChangeUnitAmount`,
		func(rs pool.AccountAnalyticLineSet) (*pool.AccountAnalyticAccountData, []models.FieldNamer) {
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
			return new(pool.AccountAnalyticAccountData), []models.FieldNamer{}
		})
}
