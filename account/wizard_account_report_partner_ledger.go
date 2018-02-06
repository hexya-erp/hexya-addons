// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountReportPartnerLedger().DeclareTransientModel()
	h.AccountReportPartnerLedger().InheritModel(h.AccountCommonPartnerReport())

	h.AccountCommonPartnerReport().AddFields(map[string]models.FieldDefinition{
		"AmountCurrency": models.BooleanField{String: "With Currency" /*["With Currency"]*/, Help: "It adds the currency column on report if the currency differs from the company currency."},
		"Reconciled":     models.BooleanField{String: "Reconciled Entries')" /*['Reconciled Entries']*/},
	})
	h.AccountCommonPartnerReport().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs h.AccountCommonPartnerReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  data = self.pre_print_report(data)
			  data['form'].update({'reconciled': self.reconciled, 'amount_currency': self.amount_currency})
			  return self.env['report'].get_action(self, 'account.report_partnerledger', data=data)
			*/
		})

}
