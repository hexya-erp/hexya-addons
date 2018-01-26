// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.TaxAdjustmentsWizard().DeclareTransientModel()
	pool.TaxAdjustmentsWizard().Methods().GetDefaultJournal().DeclareMethod(
		`GetDefaultJournal`,
		func(rs pool.TaxAdjustmentsWizardSet) {
			//@api.multi
			/*def _get_default_journal(self):
			    return self.env['account.journal'].search([('type', '=', 'general')], limit=1).id

			reason = */
		})
	pool.TaxAdjustmentsWizard().AddFields(map[string]models.FieldDefinition{
		"Reason": models.CharField{String: "Reason" /*[string 'Justification']*/, Required: true},
		"Journal": models.Many2OneField{String: "Journal", RelationModel: pool.AccountJournal(), JSON: "journal_id" /*['account.journal']*/, Required: true, Default: func(env models.Environment) interface{} {
			/*_get_default_journal(self):
			    return self.env['account.journal'].search([('type', '=', 'general')], limit=1).id

			reason = */
			return 0
		} /*, Filter: [('type'*/ /*[ ' ']*/ /*[ 'general')]]*/},
		"Date":          models.DateField{String: "Date" /*[required True]*/ /*[ default fields.Date.context_today]*/},
		"DebitAccount":  models.Many2OneField{String: "Debit account", RelationModel: pool.AccountAccount(), JSON: "debit_account_id" /*['account.account']*/, Required: true /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/},
		"CreditAccount": models.Many2OneField{String: "Credit account", RelationModel: pool.AccountAccount(), JSON: "credit_account_id" /*['account.account']*/, Required: true /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/},
		"Amount":        models.FloatField{String: "Amount" /*[currency_field 'company_currency_id']*/, Required: true},
		"CompanyCurrency": models.Many2OneField{String: "CompanyCurrencyId", RelationModel: pool.Currency(), JSON: "company_currency_id" /*['res.currency']*/ /* readonly=true */, Default: func(env models.Environment) interface{} {
			/*lambda self: self.env.user.company_id.currency_id*/
			return 0
		}},
		"Tax": models.Many2OneField{String: "Adjustment Tax", RelationModel: pool.AccountTax(), JSON: "tax_id" /*['account.tax']*/, OnDelete: models.Restrict /*, Filter: [('type_tax_use'*/ /*[ ' ']*/ /*[ 'none']*/ /*[ ('tax_adjustment']*/ /*[ ' ']*/ /*[ True)]]*/, Required: true},
	})
	pool.TaxAdjustmentsWizard().Methods().CreateMovePrivate().DeclareMethod(
		`CreateMovePrivate`,
		func(rs pool.TaxAdjustmentsWizardSet) {
			//@api.multi
			/*def _create_move(self):
			  debit_vals = {
			      'name': self.reason,
			      'debit': self.amount,
			      'credit': 0.0,
			      'account_id': self.debit_account_id.id,
			      'tax_line_id': self.tax_id.id,
			  }
			  credit_vals = {
			      'name': self.reason,
			      'debit': 0.0,
			      'credit': self.amount,
			      'account_id': self.credit_account_id.id,
			      'tax_line_id': self.tax_id.id,
			  }
			  vals = {
			      'journal_id': self.journal_id.id,
			      'date': self.date,
			      'state': 'draft',
			      'line_ids': [(0, 0, debit_vals), (0, 0, credit_vals)]
			  }
			  move = self.env['account.move'].create(vals)
			  move.post()
			  return move.id

			*/
		})
	pool.TaxAdjustmentsWizard().Methods().CreateMove().DeclareMethod(
		`CreateMove`,
		func(rs pool.TaxAdjustmentsWizardSet) {
			//@api.multi
			/*def create_move(self):
			  #create the adjustment move
			  move_id = self._create_move()
			  #return an action showing the created move
			  action = self.env.ref(self.env.context.get('action', 'account.action_move_line_form'))
			  result = action.read()[0]
			  result['views'] = [(False, 'form')]
			  result['res_id'] = move_id
			  return result
			*/
		})

}
