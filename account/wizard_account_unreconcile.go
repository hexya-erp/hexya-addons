// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountUnreconcile().DeclareTransientModel()
	pool.AccountUnreconcile().Methods().TransUnrec().DeclareMethod(
		`TransUnrec`,
		func(rs pool.AccountUnreconcileSet) {
			//@api.multi
			/*def trans_unrec(self):
			  context = dict(self._context or {})
			  if context.get('active_ids', False):
			      self.env['account.move.line'].browse(context.get('active_ids')).remove_move_reconcile()
			  return {'type': 'ir.actions.act_window_close'}
			*/
		})

}
