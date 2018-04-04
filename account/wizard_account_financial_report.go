// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountingReport().DeclareTransientModel()
	h.AccountingReport().InheritModel(h.AccountCommonReport())
	h.AccountingReport().Methods().GetAccountReport().DeclareMethod(
		`GetAccountReport`,
		func(rs h.AccountingReportSet) {
			//@api.model
			/*def _get_account_report(self):
			    reports = []
			    if self._context.get('active_id'):
			        menu = self.env['ir.ui.menu'].browse(self._context.get('active_id')).name
			        reports = self.env['account.financial.report'].search([('name', 'ilike', menu)])
			    return reports and reports[0] or False

			enable_filter = */
		})
	h.AccountingReport().AddFields(map[string]models.FieldDefinition{
		"EnableFilter": models.BooleanField{String: "EnableFilter" /*[string 'Enable Comparison']*/},
		"AccountReport": models.Many2OneField{String: "Account Reports", RelationModel: h.AccountFinancialReport(), JSON: "account_report_id" /*['account.financial.report']*/, Required: true, Default: func(env models.Environment) interface{} {
			/*_get_account_report(self):
			    reports = []
			    if self._context.get('active_id'):
			        menu = self.env['ir.ui.menu'].browse(self._context.get('active_id')).name
			        reports = self.env['account.financial.report'].search([('name', 'ilike', menu)])
			    return reports and reports[0] or False

			enable_filter = */
			return 0
		}},
		"LabelFilter": models.CharField{String: "LabelFilter" /*[string 'Column Label']*/, Help: "This label will be displayed on report to show the balance computed for the given comparison filter."},
		"FilterCmp": models.SelectionField{String: "Filter by", Selection: types.Selection{
			"filter_no":   "No Filters",
			"filter_date": "Date",
		}, /*[]*/ Required: true, Default: models.DefaultValue("filter_no")},
		"DateFromCmp": models.DateField{String: "DateFromCmp" /*[string 'Start Date']*/},
		"DateToCmp":   models.DateField{String: "DateToCmp" /*[string 'End Date']*/},
		"DebitCredit": models.BooleanField{String: "DebitCredit" /*[string 'Display Debit/Credit Columns']*/, Help: "This option allows you to get more details about the way your balances are computed. Because it is space consuming, we do not allow to use it while doing a comparison." /*[ we do not allow to use it while doing a comparison."]*/},
	})
	h.AccountingReport().Methods().BuildComparisonContext().DeclareMethod(
		`BuildComparisonContext`,
		func(rs h.AccountingReportSet, args struct {
			Data interface{}
		}) {
			/*def _build_comparison_context(self, data):
			  result = {}
			  result['journal_ids'] = 'journal_ids' in data['form'] and data['form']['journal_ids'] or False
			  result['state'] = 'target_move' in data['form'] and data['form']['target_move'] or ''
			  if data['form']['filter_cmp'] == 'filter_date':
			      result['date_from'] = data['form']['date_from_cmp']
			      result['date_to'] = data['form']['date_to_cmp']
			      result['strict_range'] = True
			  return result

			*/
		})
	h.AccountingReport().Methods().CheckReport().DeclareMethod(
		`CheckReport`,
		func(rs h.AccountingReportSet) {
			//@api.multi
			/*def check_report(self):
			  res = super(AccountingReport, self).check_report()
			  data = {}
			  data['form'] = self.read(['account_report_id', 'date_from_cmp', 'date_to_cmp', 'journal_ids', 'filter_cmp', 'target_move'])[0]
			  for field in ['account_report_id']:
			      if isinstance(data['form'][field], tuple):
			          data['form'][field] = data['form'][field][0]
			  comparison_context = self._build_comparison_context(data)
			  res['data']['form']['comparison_context'] = comparison_context
			  return res

			*/
		})
	h.AccountingReport().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs h.AccountingReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  data['form'].update(self.read(['date_from_cmp', 'debit_credit', 'date_to_cmp', 'filter_cmp', 'account_report_id', 'enable_filter', 'label_filter', 'target_move'])[0])
			  return self.env['report'].get_action(self, 'account.report_financial', data=data)
			*/
		})

}
