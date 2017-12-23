// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountCommonJournalReport().DeclareMixinModel()
	pool.AccountCommonJournalReport().InheritModel(pool.AccountCommonReport())

	pool.AccountCommonJournalReport().AddFields(map[string]models.FieldDefinition{
		"AmountCurrency": models.BooleanField{String: "With Currency" /*['With Currency']*/, Help: "Print Report with the currency column if the currency differs from the company currency."},
	})
	pool.AccountCommonJournalReport().Methods().PrePrintReport().DeclareMethod(
		`PrePrintReport`,
		func(rs pool.AccountCommonJournalReportSet, args struct {
			Data interface{}
		}) {
			//@api.multi
			/*def pre_print_report(self, data):
			  data['form'].update({'amount_currency': self.amount_currency})
			  return data
			*/
		})

}
