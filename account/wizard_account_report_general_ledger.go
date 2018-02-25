// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountReportGeneralLedger().DeclareTransientModel()
	h.AccountReportGeneralLedger().InheritModel(h.AccountCommonAccountReport())

	h.AccountReportGeneralLedger().AddFields(map[string]models.FieldDefinition{
		"InitialBalance": models.BooleanField{String: "Include Initial Balances", Help: `If you selected date, this field allow you to add a row to display the amount of debit/credit/balance that precedes the filter you\'ve set."/*[ this field allow you to add a row to display the amount of debit/credit/balance that precedes the filter you\'ve set.']`},
		"Sortby": models.SelectionField{String: "Sort by", Selection: types.Selection{
			"sort_date":            "Date",
			"sort_journal_partner": "Journal & Partner",
		}, Required: true, Default: models.DefaultValue("sort_date")},
	})
	h.AccountReportGeneralLedger().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs h.AccountCommonAccountReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  data = self.pre_print_report(data)
			  data['form'].update(self.read(['initial_balance', 'sortby'])[0])
			  if data['form'].get('initial_balance') and not data['form'].get('date_from'):
			      raise UserError(_("You must define a Start Date"))
			  records = self.env[data['model']].browse(data.get('ids', []))
			  return self.env['report'].with_context(landscape=True).get_action(records, 'account.report_generalledger', data=data)
			*/
		})

}
