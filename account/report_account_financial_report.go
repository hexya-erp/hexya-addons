// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountFinancialReport().DeclareModel()
	h.AccountFinancialReport().Methods().GetLevel().DeclareMethod(
		`GetLevel`,
		func(rs h.AccountFinancialReportSet) *h.AccountFinancialReportData {
			//@api.depends('parent_id','parent_id.level')
			/*def _get_level(self):
			  '''Returns a dictionary with key=the ID of a record and value = the level of this
			     record in the tree structure.'''
			  for report in self:
			      level = 0
			      if report.parent_id:
			          level = report.parent_id.level + 1
			      report.level = level

			*/
			return &h.AccountFinancialReportData{}
		})
	h.AccountFinancialReport().Methods().GetChildrenByOrder().DeclareMethod(
		`GetChildrenByOrder`,
		func(rs h.AccountFinancialReportSet) {
			/*def _get_children_by_order(self):
			    '''returns a recordset of all the children computed recursively, and sorted by sequence. Ready for the printing'''
			    res = self
			    children = self.search([('parent_id', 'in', self.ids)], order='sequence ASC')
			    if children:
			        for child in children:
			            res += child._get_children_by_order()
			    return res

			name = */
		})
	h.AccountFinancialReport().AddFields(map[string]models.FieldDefinition{
		"Name":      models.CharField{String: "Report Name" /*['Report Name']*/, Required: true, Translate: true},
		"Parent":    models.Many2OneField{String: "Parent", RelationModel: h.AccountFinancialReport(), JSON: "parent_id" /*['account.financial.report']*/ /*['Parent']*/},
		"Childrens": models.One2ManyField{String: "Account Report", RelationModel: h.AccountFinancialReport(), ReverseFK: "Parent", JSON: "children_ids" /*['account.financial.report']*/ /*[ 'parent_id']*/ /*['Account Report']*/},
		"Sequence":  models.IntegerField{String: "Sequence')" /*['Sequence']*/},
		"Level":     models.IntegerField{String: "Level", Compute: h.AccountFinancialReport().Methods().GetLevel() /*[ string 'Level']*/ /*[ store True]*/},
		"Type": models.SelectionField{String: "Type", Selection: types.Selection{
			"sum":            "View",
			"accounts":       "Accounts",
			"account_type":   "Account Type",
			"account_report": "Report Value",
			/*[ ('sum', 'View'  ('accounts', 'Accounts'  ('account_type', 'Account Type'  ('account_report', 'Report Value'  ]*/}, /*[]*/ /*['Type']*/ Default: models.DefaultValue("sum")},
		"Accounts":      models.Many2ManyField{String: "account_account_financial_report", RelationModel: h.AccountAccount(), JSON: "account_ids" /*['account.account']*/ /*['account_account_financial_report']*/ /*[ 'report_line_id']*/ /*[ 'account_id']*/ /*[ 'Accounts']*/},
		"AccountReport": models.Many2OneField{String: "Report Value", RelationModel: h.AccountFinancialReport(), JSON: "account_report_id" /*['account.financial.report']*/ /*['Report Value']*/},
		"AccountTypes":  models.Many2ManyField{String: "account_account_financial_report_type", RelationModel: h.AccountAccountType(), JSON: "account_type_ids" /*['account.account.type']*/ /*['account_account_financial_report_type']*/ /*[ 'report_id']*/ /*[ 'account_type_id']*/ /*[ 'Account Types']*/},
		"Sign": models.SelectionField{String: "Sign on Reports", Selection: types.Selection{
			"-1": "Reverse balance sign",
			"1":  "Preserve balance sign",
		}, /*[]*/ /*['Sign on Reports']*/ Required: true, Default: models.DefaultValue("1"), Help: "For accounts that are typically more debited than credited and that you would like to print as negative amounts in your reports" /*[ you should reverse the sign of the balance; e.g.: Expense account. The same applies for accounts that are typically more credited than debited and that you would like to print as positive amounts in your reports; e.g.: Income account.']*/},
		"DisplayDetail": models.SelectionField{ /*display_detail = fields.Selection([ ('no_detail', 'No detail'), ('detail_flat', 'Display children flat'), ('detail_with_hierarchy', 'Display children with hierarchy')*/ },
		"StyleOverwrite": models.SelectionField{String: "Financial Report Style", Selection: types.Selection{
			"0": "Automatic formatting",
			"1": "Main Title 1 (bold underlined)",
			"2": "Title 2 (bold)",
			"3": "Title 3 (bold smaller)",
			"4": "Normal Text",
			"5": "Italic Text (smaller)",
			"6": "Smallest Text",
			/*[ (0, 'Automatic formatting'  (1, 'Main Title 1 (bold, underlined)'  (2, 'Title 2 (bold)'  (3, 'Title 3 (bold, smaller)'  (4, 'Normal Text'  (5, 'Italic Text (smaller)'  (6, 'Smallest Text'  ]*/}, /*[]*/ /*['Financial Report Style']*/ Default: models.DefaultValue("0"), Help: "You can set up here the format you want this record to be displayed. If you leave the automatic formatting" /*[ it will be computed based on the financial reports hierarchy (auto-computed field 'level')."]*/},
	})

}
