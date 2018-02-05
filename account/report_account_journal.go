// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import "github.com/hexya-erp/hexya/pool/h"

func init() {
	h.ReportAccountReportJournal().DeclareTransientModel()
	h.ReportAccountReportJournal().Methods().Lines().DeclareMethod(
		`Lines`,
		func(rs h.ReportAccountReportJournalSet, args struct {
			TargetMove    interface{}
			JournalIds    interface{}
			SortSelection interface{}
			Data          interface{}
		}) {
			/*def lines(self, target_move, journal_ids, sort_selection, data):
			  if isinstance(journal_ids, int):
			      journal_ids = [journal_ids]

			  move_state = ['draft', 'posted']
			  if target_move == 'posted':
			      move_state = ['posted']

			  query_get_clause = self._get_query_get_clause(data)
			  params = [tuple(move_state), tuple(journal_ids)] + query_get_clause[2]
			  query = 'SELECT "account_move_line".id FROM ' + query_get_clause[0] + ', account_move am, account_account acc WHERE "account_move_line".account_id = acc.id AND "account_move_line".move_id=am.id AND am.state IN %s AND "account_move_line".journal_id IN %s AND ' + query_get_clause[1] + ' ORDER BY '
			  if sort_selection == 'date':
			      query += '"account_move_line".date'
			  else:
			      query += 'am.name'
			  query += ', "account_move_line".move_id, acc.code'
			  self.env.cr.execute(query, tuple(params))
			  ids = map(lambda x: x[0], self.env.cr.fetchall())
			  return self.env['account.move.line'].browse(ids)

			*/
		})
	h.ReportAccountReportJournal().Methods().SumDebit().DeclareMethod(
		`SumDebit`,
		func(rs h.ReportAccountReportJournalSet, args struct {
			Data      interface{}
			JournalId interface{}
		}) {
			/*def _sum_debit(self, data, journal_id):
			  move_state = ['draft', 'posted']
			  if data['form'].get('target_move', 'all') == 'posted':
			      move_state = ['posted']

			  query_get_clause = self._get_query_get_clause(data)
			  params = [tuple(move_state), tuple(journal_id.ids)] + query_get_clause[2]
			  self.env.cr.execute('SELECT SUM(debit) FROM ' + query_get_clause[0] + ', account_move am '
			                  'WHERE "account_move_line".move_id=am.id AND am.state IN %s AND "account_move_line".journal_id IN %s AND ' + query_get_clause[1] + ' ',
			                  tuple(params))
			  return self.env.cr.fetchone()[0] or 0.0

			*/
		})
	h.ReportAccountReportJournal().Methods().SumCredit().DeclareMethod(
		`SumCredit`,
		func(rs h.ReportAccountReportJournalSet, args struct {
			Data      interface{}
			JournalId interface{}
		}) {
			/*def _sum_credit(self, data, journal_id):
			  move_state = ['draft', 'posted']
			  if data['form'].get('target_move', 'all') == 'posted':
			      move_state = ['posted']

			  query_get_clause = self._get_query_get_clause(data)
			  params = [tuple(move_state), tuple(journal_id.ids)] + query_get_clause[2]
			  self.env.cr.execute('SELECT SUM(credit) FROM ' + query_get_clause[0] + ', account_move am '
			                  'WHERE "account_move_line".move_id=am.id AND am.state IN %s AND "account_move_line".journal_id IN %s AND ' + query_get_clause[1] + ' ',
			                  tuple(params))
			  return self.env.cr.fetchone()[0] or 0.0

			*/
		})
	h.ReportAccountReportJournal().Methods().GetTaxes().DeclareMethod(
		`GetTaxes`,
		func(rs h.ReportAccountReportJournalSet, args struct {
			Data      interface{}
			JournalId interface{}
		}) {
			/*def _get_taxes(self, data, journal_id):
			  move_state = ['draft', 'posted']
			  if data['form'].get('target_move', 'all') == 'posted':
			      move_state = ['posted']

			  query_get_clause = self._get_query_get_clause(data)
			  params = [tuple(move_state), tuple(journal_id.ids)] + query_get_clause[2]
			  query = """
			      SELECT rel.account_tax_id, SUM("account_move_line".balance) AS base_amount
			      FROM account_move_line_account_tax_rel rel, """ + query_get_clause[0] + """
			      LEFT JOIN account_move am ON "account_move_line".move_id = am.id
			      WHERE "account_move_line".id = rel.account_move_line_id
			          AND am.state IN %s
			          AND "account_move_line".journal_id IN %s
			          AND """ + query_get_clause[1] + """
			     GROUP BY rel.account_tax_id"""
			  self.env.cr.execute(query, tuple(params))
			  ids = []
			  base_amounts = {}
			  for row in self.env.cr.fetchall():
			      ids.append(row[0])
			      base_amounts[row[0]] = row[1]


			  res = {}
			  for tax in self.env['account.tax'].browse(ids):
			      self.env.cr.execute('SELECT sum(debit - credit) FROM ' + query_get_clause[0] + ', account_move am '
			          'WHERE "account_move_line".move_id=am.id AND am.state IN %s AND "account_move_line".journal_id IN %s AND ' + query_get_clause[1] + ' AND tax_line_id = %s',
			          tuple(params + [tax.id]))
			      res[tax] = {
			          'base_amount': base_amounts[tax.id],
			          'tax_amount': self.env.cr.fetchone()[0] or 0.0,
			      }
			      if journal_id.type == 'sale':
			          #sales operation are credits
			          res[tax]['base_amount'] = res[tax]['base_amount'] * -1
			          res[tax]['tax_amount'] = res[tax]['tax_amount'] * -1
			  return res

			*/
		})
	h.ReportAccountReportJournal().Methods().GetQueryGetClause().DeclareMethod(
		`GetQueryGetClause`,
		func(rs h.ReportAccountReportJournalSet, args struct {
			Data interface{}
		}) {
			/*def _get_query_get_clause(self, data):
			  return self.env['account.move.line'].with_context(data['form'].get('used_context', {}))._query_get()

			*/
		})
	h.ReportAccountReportJournal().Methods().RenderHtml().DeclareMethod(
		`RenderHtml`,
		func(rs h.ReportAccountReportJournalSet, args struct {
			Docids interface{}
			Data   interface{}
		}) {
			//@api.model
			/*def render_html(self, docids, data=None):
			  if not data.get('form'):
			      raise UserError(_("Form content is missing, this report cannot be printed."))

			  target_move = data['form'].get('target_move', 'all')
			  sort_selection = data['form'].get('sort_selection', 'date')

			  res = {}
			  for journal in data['form']['journal_ids']:
			      res[journal] = self.with_context(data['form'].get('used_context', {})).lines(target_move, journal, sort_selection, data)
			  docargs = {
			      'doc_ids': data['form']['journal_ids'],
			      'doc_model': self.env['account.journal'],
			      'data': data,
			      'docs': self.env['account.journal'].browse(data['form']['journal_ids']),
			      'time': time,
			      'lines': res,
			      'sum_credit': self._sum_credit,
			      'sum_debit': self._sum_debit,
			      'get_taxes': self._get_taxes,
			  }
			  return self.env['report'].render('account.report_journal', docargs)
			*/
		})

}
