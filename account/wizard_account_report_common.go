// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountCommonReport().DeclareTransientModel()
	pool.AccountCommonReport().AddFields(map[string]models.FieldDefinition{
		"Company": models.Many2OneField{String: "Company", RelationModel: pool.Company(), JSON: "company_id" /*['res.company']*/ /* readonly=true */, Default: func(models.Environment, models.FieldMap) interface{} {
			/*lambda self: self.env.user.company_id*/
			return 0
		}},
		"Journals": models.Many2ManyField{String: "Journals", RelationModel: pool.AccountJournal(), JSON: "journal_ids" /*['account.journal']*/ /*[ required True]*/ /*[ default lambda self: self.env['account.journal'].search([]]*/},
		"DateFrom": models.DateField{String: "DateFrom" /*[string 'Start Date']*/},
		"DateTo":   models.DateField{String: "DateTo" /*[string 'End Date']*/},
		"TargetMove": models.SelectionField{String: "Target Moves", Selection: types.Selection{
			"posted": "All Posted Entries",
			"all":    "All Entries",
			/*[('posted', 'All Posted Entries'  ('all', 'All Entries'  ]*/}, /*[]*/ Required: true, Default: models.DefaultValue("posted")},
	})
	pool.AccountCommonReport().Methods().BuildContexts().DeclareMethod(
		`BuildContexts`,
		func(rs pool.AccountCommonReportSet, args struct {
			Data interface{}
		}) {
			/*def _build_contexts(self, data):
			  result = {}
			  result['journal_ids'] = 'journal_ids' in data['form'] and data['form']['journal_ids'] or False
			  result['state'] = 'target_move' in data['form'] and data['form']['target_move'] or ''
			  result['date_from'] = data['form']['date_from'] or False
			  result['date_to'] = data['form']['date_to'] or False
			  result['strict_range'] = True if result['date_from'] else False
			  return result

			*/
		})
	pool.AccountCommonReport().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs pool.AccountCommonReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  raise (_('Error!'), _('Not implemented.'))

			*/
		})
	pool.AccountCommonReport().Methods().CheckReport().DeclareMethod(
		`CheckReport`,
		func(rs pool.AccountCommonReportSet) {
			//@api.multi
			/*def check_report(self):
			  self.ensure_one()
			  data = {}
			  data['ids'] = self.env.context.get('active_ids', [])
			  data['model'] = self.env.context.get('active_model', 'ir.ui.menu')
			  data['form'] = self.read(['date_from', 'date_to', 'journal_ids', 'target_move'])[0]
			  used_context = self._build_contexts(data)
			  data['form']['used_context'] = dict(used_context, lang=self.env.context.get('lang', 'en_US'))
			  return self._print_report(data)
			*/
		})

}
