// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import "github.com/hexya-erp/hexya/pool/h"

func init() {
	h.ReportAccountReportOverdue().DeclareTransientModel()
	h.ReportAccountReportOverdue().Methods().GetAccountMoveLines().DeclareMethod(
		`GetAccountMoveLines`,
		func(rs h.ReportAccountReportOverdueSet, args struct {
			PartnerIds interface{}
		}) {
			/*def _get_account_move_lines(self, partner_ids):
			  res = dict(map(lambda x:(x,[]), partner_ids))
			  self.env.cr.execute("SELECT m.name AS move_id, l.date, l.name, l.ref, l.date_maturity, l.partner_id, l.blocked, l.amount_currency, l.currency_id, "
			      "CASE WHEN at.type = 'receivable' "
			          "THEN SUM(l.debit) "
			          "ELSE SUM(l.credit * -1) "
			      "END AS debit, "
			      "CASE WHEN at.type = 'receivable' "
			          "THEN SUM(l.credit) "
			          "ELSE SUM(l.debit * -1) "
			      "END AS credit, "
			      "CASE WHEN l.date_maturity < %s "
			          "THEN SUM(l.debit - l.credit) "
			          "ELSE 0 "
			      "END AS mat "
			      "FROM account_move_line l "
			      "JOIN account_account_type at ON (l.user_type_id = at.id) "
			      "JOIN account_move m ON (l.move_id = m.id) "
			      "WHERE l.partner_id IN %s AND at.type IN ('receivable', 'payable') AND NOT l.reconciled GROUP BY l.date, l.name, l.ref, l.date_maturity, l.partner_id, at.type, l.blocked, l.amount_currency, l.currency_id, l.move_id, m.name", (((*/
		})
	h.ReportAccountReportOverdue().Methods().RenderHtml().DeclareMethod(
		`RenderHtml`,
		func(rs h.ReportAccountReportOverdueSet, args struct {
			Docids interface{}
			Data   interface{}
		}) {
			//@api.model
			/*def render_html(self, docids, data=None):
			  totals = {}
			  lines = self._get_account_move_lines(docids)
			  lines_to_display = {}
			  company_currency = self.env.user.company_id.currency_id
			  for partner_id in docids:
			      lines_to_display[partner_id] = {}
			      totals[partner_id] = {}
			      for line_tmp in lines[partner_id]:
			          line = line_tmp.copy()
			          currency = line['currency_id'] and self.env['res.currency'].browse(line['currency_id']) or company_currency
			          if currency not in lines_to_display[partner_id]:
			              lines_to_display[partner_id][currency] = []
			              totals[partner_id][currency] = dict((fn, 0.0) for fn in ['due', 'paid', 'mat', 'total'])
			          if line['debit'] and line['currency_id']:
			              line['debit'] = line['amount_currency']
			          if line['credit'] and line['currency_id']:
			              line['credit'] = line['amount_currency']
			          if line['mat'] and line['currency_id']:
			              line['mat'] = line['amount_currency']
			          lines_to_display[partner_id][currency].append(line)
			          if not line['blocked']:
			              totals[partner_id][currency]['due'] += line['debit']
			              totals[partner_id][currency]['paid'] += line['credit']
			              totals[partner_id][currency]['mat'] += line['mat']
			              totals[partner_id][currency]['total'] += line['debit'] - line['credit']
			  docargs = {
			      'doc_ids': docids,
			      'doc_model': 'res.partner',
			      'docs': self.env['res.partner'].browse(docids),
			      'time': time,
			      'Lines': lines_to_display,
			      'Totals': totals,
			      'Date': */
		})

}
