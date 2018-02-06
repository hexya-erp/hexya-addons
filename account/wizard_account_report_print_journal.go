// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountPrintJournal().DeclareTransientModel()
	h.AccountPrintJournal().InheritModel(h.AccountCommonJournalReport())

	h.AccountPrintJournal().AddFields(map[string]models.FieldDefinition{
		"SortSelection": models.SelectionField{String: "Entries Sorted by", Selection: types.Selection{
			"date":      "Date",
			"move_name": "Journal Entry Number",
			/*[('date', 'Date'  ('move_name', 'Journal Entry Number' ]*/}, /*[]*/ /*['Entries Sorted by']*/ Required: true, Default: models.DefaultValue("move_name")},
		"Journals": models.Many2ManyField{String: "Journals", RelationModel: h.AccountJournal(), JSON: "journal_ids" /*['account.journal']*/ /*[ required True]*/ /*[ default lambda self: self.env['account.journal'].search([('type']*/ /*[ 'in']*/ /*[ ['sale']*/ /*[ 'purchase'])]]*/},
	})
	h.AccountPrintJournal().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs h.AccountCommonJournalReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  data = self.pre_print_report(data)
			  data['form'].update({'sort_selection': self.sort_selection})
			  return self.env['report'].with_context(landscape=True).get_action(self, 'account.report_journal', data=data)
			*/
		})

}
