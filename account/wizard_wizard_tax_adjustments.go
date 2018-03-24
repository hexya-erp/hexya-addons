// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.TaxAdjustmentsWizard().DeclareTransientModel()

	h.TaxAdjustmentsWizard().AddFields(map[string]models.FieldDefinition{
		"Reason": models.CharField{String: "Justification", Required: true},
		"Journal": models.Many2OneField{RelationModel: h.AccountJournal(), Required: true,
			Default: func(env models.Environment) interface{} {
				return h.AccountJournal().Search(env, q.AccountJournal().Type().Equals("general")).Limit(1)
			}, Filter: q.AccountJournal().Type().Equals("general")},
		"Date": models.DateField{Required: true, Default: func(env models.Environment) interface{} {
			return dates.Today()
		}},
		"DebitAccount": models.Many2OneField{RelationModel: h.AccountAccount(), Required: true,
			Filter: q.AccountAccount().Deprecated().Equals(false)},
		"CreditAccount": models.Many2OneField{RelationModel: h.AccountAccount(), Required: true,
			Filter: q.AccountAccount().Deprecated().Equals(false)},
		"Amount": models.FloatField{ /*[currency_field 'company_currency_id']*/ Required: true},
		"CompanyCurrency": models.Many2OneField{RelationModel: h.Currency(), ReadOnly: true,
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company()
			}},
		"Tax": models.Many2OneField{String: "Adjustment Tax", RelationModel: h.AccountTax(),
			OnDelete: models.Restrict, Required: true,
			Filter: q.AccountTax().TypeTaxUse().Equals("none").And().TaxAdjustment().Equals(true)},
	})

	h.TaxAdjustmentsWizard().Methods().CreateMovePrivate().DeclareMethod(
		`CreateMovePrivate`,
		func(rs h.TaxAdjustmentsWizardSet) {
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

	h.TaxAdjustmentsWizard().Methods().CreateMove().DeclareMethod(
		`CreateMove`,
		func(rs h.TaxAdjustmentsWizardSet) {
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
