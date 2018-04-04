// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountMoveLineReconcile().DeclareTransientModel()
	h.AccountMoveLineReconcile().AddFields(map[string]models.FieldDefinition{
		"TransNbr": models.IntegerField{String: "# of Transaction", ReadOnly: true},
		"Credit": models.FloatField{String: "Credit amount", ReadOnly: true,
			Digits: nbutils.Digits{0, 0}},
		"Debit": models.FloatField{String: "Debit amount", ReadOnly: true,
			Digits: nbutils.Digits{0, 0}},
		"Writeoff": models.FloatField{String: "Write-off amount", ReadOnly: true,
			Digits: nbutils.Digits{0, 0}},
		"Company": models.Many2OneField{String: "Company", RelationModel: h.Company(),
			Required: true, Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company()
			}},
	})
	h.AccountMoveLineReconcile().Methods().DefaultGet().DeclareMethod(
		`DefaultGet`,
		func(rs h.AccountMoveLineReconcileSet, args struct {
			Fields interface{}
		}) {
			//@api.model
			/*def default_get(self, fields):
			  res = super(AccountMoveLineReconcile, self).default_get(fields)
			  data = self.trans_rec_get()
			  if 'trans_nbr' in fields:
			      res.update({'trans_nbr': data['trans_nbr']})
			  if 'credit' in fields:
			      res.update({'credit': data['credit']})
			  if 'debit' in fields:
			      res.update({'debit': data['debit']})
			  if 'writeoff' in fields:
			      res.update({'writeoff': data['writeoff']})
			  return res

			*/
		})
	h.AccountMoveLineReconcile().Methods().TransRecGet().DeclareMethod(
		`TransRecGet`,
		func(rs h.AccountMoveLineReconcileSet) {
			//@api.multi
			/*def trans_rec_get(self):
			  context = self._context or {}
			  credit = debit = 0
			  lines = self.env['account.move.line'].browse(context.get('active_ids', []))
			  for line in lines:
			      if not line.full_reconcile_id:
			          credit += line.credit
			          debit += line.debit
			  precision = self.env.user.company_id.currency_id.decimal_places
			  writeoff = float_round(debit - credit, precision_digits=precision)
			  credit = float_round(credit, precision_digits=precision)
			  debit = float_round(debit, precision_digits=precision)
			  return {'trans_nbr': len(lines), 'credit': credit, 'debit': debit, 'writeoff': writeoff}

			*/
		})
	h.AccountMoveLineReconcile().Methods().TransRecAddendumWriteoff().DeclareMethod(
		`TransRecAddendumWriteoff`,
		func(rs h.AccountMoveLineReconcileSet) {
			//@api.multi
			/*def trans_rec_addendum_writeoff(self):
			  return self.env['account.move.line.reconcile.writeoff'].trans_rec_addendum()

			*/
		})
	h.AccountMoveLineReconcile().Methods().TransRecReconcilePartialReconcile().DeclareMethod(
		`TransRecReconcilePartialReconcile`,
		func(rs h.AccountMoveLineReconcileSet) {
			//@api.multi
			/*def trans_rec_reconcile_partial_reconcile(self):
			  return self.env['account.move.line.reconcile.writeoff'].trans_rec_reconcile_partial()

			*/
		})
	h.AccountMoveLineReconcile().Methods().TransRecReconcileFull().DeclareMethod(
		`TransRecReconcileFull`,
		func(rs h.AccountMoveLineReconcileSet) {
			//@api.multi
			/*def trans_rec_reconcile_full(self):
			  move_lines = self.env['account.move.line'].browse(self._context.get('active_ids', []))
			  currency = False
			  for aml in move_lines:
			      if not currency and aml.currency_id.id:
			          currency = aml.currency_id.id
			      elif aml.currency_id:
			          if aml.currency_id.id == currency:
			              continue
			          raise UserError(_('Operation not allowed. You can only reconcile entries that share the same secondary currency or that don\'t have one. Edit your journal items or make another selection before proceeding any further.'))
			  #Don't consider entrires that are already reconciled
			  move_lines_filtered = move_lines.filtered(lambda aml: not aml.reconciled)
			  #Because we are making a full reconcilition in batch, we need to consider use cases as defined in the test test_manual_reconcile_wizard_opw678153
			  #So we force the reconciliation in company currency only at first
			  move_lines_filtered.with_context(skip_full_reconcile_check='amount_currency_excluded', manual_full_reconcile_currency=currency).reconcile()

			  #then in second pass the amounts in secondary currency, only if some lines are still not fully reconciled
			  move_lines_filtered = move_lines.filtered(lambda aml: not aml.reconciled)
			  if move_lines_filtered:
			      move_lines_filtered.with_context(skip_full_reconcile_check='amount_currency_only', manual_full_reconcile_currency=currency).reconcile()
			  move_lines.compute_full_after_batch_reconcile()
			  return {'type': 'ir.actions.act_window_close'}


			*/
		})

	h.AccountMoveLineReconcileWriteoff().DeclareTransientModel()
	h.AccountMoveLineReconcileWriteoff().AddFields(map[string]models.FieldDefinition{
		"Journal":     models.Many2OneField{String: "Write-Off Journal", RelationModel: h.AccountJournal(), JSON: "journal_id" /*['account.journal']*/, Required: true},
		"WriteoffAcc": models.Many2OneField{String: "Write-Off account", RelationModel: h.AccountAccount(), JSON: "writeoff_acc_id" /*['account.account']*/, Required: true /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/},
		"DateP":       models.DateField{String: "DateP" /*[string 'Date']*/ /*[ default fields.Date.context_today]*/},
		"Comment":     models.CharField{String: "Comment", Required: true /*[ default 'Write-off']*/},
		"Analytic":    models.Many2OneField{String: "Analytic Account", RelationModel: h.AccountAnalyticAccount(), JSON: "analytic_id" /*['account.analytic.account']*/},
	})
	h.AccountMoveLineReconcileWriteoff().Methods().TransRecAddendum().DeclareMethod(
		`TransRecAddendum`,
		func(rs h.AccountMoveLineReconcileWriteoffSet) {
			//@api.multi
			/*def trans_rec_addendum(self):
			  view = self.env.ref('account.account_move_line_reconcile_writeoff')
			  return {
			      'name': _('Reconcile Writeoff'),
			      'context': self._context,
			      'view_type': 'form',
			      'view_mode': 'form',
			      'res_model': 'account.move.line.reconcile.writeoff',
			      'views': [(view.id, 'form')],
			      'type': 'ir.actions.act_window',
			      'target': 'new',
			  }

			*/
		})
	h.AccountMoveLineReconcileWriteoff().Methods().TransRecReconcilePartial().DeclareMethod(
		`TransRecReconcilePartial`,
		func(rs h.AccountMoveLineReconcileWriteoffSet) {
			//@api.multi
			/*def trans_rec_reconcile_partial(self):
			  context = self._context or {}
			  self.env['account.move.line'].browse(context.get('active_ids', [])).reconcile()
			  return {'type': 'ir.actions.act_window_close'}

			*/
		})
	h.AccountMoveLineReconcileWriteoff().Methods().TransRecReconcile().DeclareMethod(
		`TransRecReconcile`,
		func(rs h.AccountMoveLineReconcileWriteoffSet) {
			//@api.multi
			/*def trans_rec_reconcile(self):
			  context = dict(self._context or {})
			  context['date_p'] = self.date_p
			  context['comment'] = self.comment
			  if self.analytic_id:
			      context['analytic_id'] = self.analytic_id.id
			  move_lines = self.env['account.move.line'].browse(self._context.get('active_ids', []))
			  currency = False
			  for aml in move_lines:
			      if not currency and aml.currency_id.id:
			          currency = aml.currency_id.id
			      elif aml.currency_id:
			          if aml.currency_id.id == currency:
			              continue
			          raise UserError(_('Operation not allowed. You can only reconcile entries that share the same secondary currency or that don\'t have one. Edit your journal items or make another selection before proceeding any further.'))
			  #Don't consider entrires that are already reconciled
			  move_lines_filtered = move_lines.filtered(lambda aml: not aml.reconciled)
			  #Because we are making a full reconcilition in batch, we need to consider use cases as defined in the test test_manual_reconcile_wizard_opw678153
			  #So we force the reconciliation in company currency only at first,
			  context['skip_full_reconcile_check'] = 'amount_currency_excluded'
			  context['manual_full_reconcile_currency'] = currency
			  writeoff = move_lines_filtered.with_context(context).reconcile(self.writeoff_acc_id, self.journal_id)
			  #then in second pass the amounts in secondary currency, only if some lines are still not fully reconciled
			  move_lines_filtered = move_lines.filtered(lambda aml: not aml.reconciled)
			  if move_lines_filtered:
			      move_lines_filtered.with_context(skip_full_reconcile_check='amount_currency_only', manual_full_reconcile_currency=currency).reconcile()
			  if not isinstance(writeoff, bool):
			      move_lines += writeoff
			  move_lines.compute_full_after_batch_reconcile()
			  return {'type': 'ir.actions.act_window_close'}
			*/
		})

}
