// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountMoveReversal().DeclareTransientModel()
	h.AccountMoveReversal().AddFields(map[string]models.FieldDefinition{
		"Date":    models.DateField{String: "Date" /*[string 'Reversal date']*/ /*[ default fields.Date.context_today]*/ /*[ required True]*/},
		"Journal": models.Many2OneField{String: "Use Specific Journal", RelationModel: h.AccountJournal(), JSON: "journal_id" /*['account.journal']*/, Help: "If empty, uses the journal of the journal entry to be reversed." /*[ uses the journal of the journal entry to be reversed.']*/},
	})
	h.AccountMoveReversal().Methods().ReverseMoves().DeclareMethod(
		`ReverseMoves`,
		func(rs h.AccountMoveReversalSet) {
			//@api.multi
			/*def reverse_moves(self):
			  ac_move_ids = self._context.get('active_ids', False)
			  res = self.env['account.move'].browse(ac_move_ids).reverse_moves(self.date, self.journal_id or False)
			  if res:
			      return {
			          'name': _('Reverse Moves'),
			          'type': 'ir.actions.act_window',
			          'view_type': 'form',
			          'view_mode': 'tree,form',
			          'res_model': 'account.move',
			          'domain': [('id', 'in', res)],
			      }
			  return {'type': 'ir.actions.act_window_close'}
			*/
		})

}
