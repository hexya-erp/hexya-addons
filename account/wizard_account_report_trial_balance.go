// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountCommonAccountReport().AddFields(map[string]models.FieldDefinition{
		"Journals": models.Many2ManyField{String: "Journals", RelationModel: pool.AccountJournal(), JSON: "journal_ids" /*['account.journal']*/ /*['account_balance_report_journal_rel']*/ /*[ 'account_id']*/ /*[ 'journal_id']*/ /*[ required True]*/ /*[ default []]*/},
	})
	pool.AccountCommonAccountReport().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs pool.AccountCommonAccountReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  data = self.pre_print_report(data)
			  records = self.env[data['model']].browse(data.get('ids', []))
			  return self.env['report'].get_action(records, 'account.report_trialbalance', data=data)
			*/
		})

}
