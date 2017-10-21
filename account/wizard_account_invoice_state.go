// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountInvoiceConfirm().DeclareTransientModel()
	pool.AccountInvoiceConfirm().Methods().InvoiceConfirm().DeclareMethod(
		`InvoiceConfirm`,
		func(rs pool.AccountInvoiceConfirmSet) {
			//@api.multi
			/*def invoice_confirm(self):
			  context = dict(self._context or {})
			  active_ids = context.get('active_ids', []) or []

			  for record in self.env['account.invoice'].browse(active_ids):
			      if record.state not in ('draft', 'proforma', 'proforma2'):
			          raise UserError(_("Selected invoice(s) cannot be confirmed as they are not in 'Draft' or 'Pro-Forma' state."))
			      record.action_invoice_open()
			  return {'type': 'ir.actions.act_window_close'}


			*/
		})

	pool.AccountInvoiceCancel().DeclareTransientModel()
	pool.AccountInvoiceCancel().Methods().InvoiceCancel().DeclareMethod(
		`InvoiceCancel`,
		func(rs pool.AccountInvoiceCancelSet) {
			//@api.multi
			/*def invoice_cancel(self):
			  context = dict(self._context or {})
			  active_ids = context.get('active_ids', []) or []

			  for record in self.env['account.invoice'].browse(active_ids):
			      if record.state in ('cancel', 'paid'):
			          raise UserError(_("Selected invoice(s) cannot be cancelled as they are already in 'Cancelled' or 'Done' state."))
			      record.action_invoice_cancel()
			  return {'type': 'ir.actions.act_window_close'}
			*/
		})

}
