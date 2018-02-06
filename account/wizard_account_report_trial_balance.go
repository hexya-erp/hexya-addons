// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import "github.com/hexya-erp/hexya/pool/h"

func init() {

	h.AccountBalanceReport().DeclareTransientModel()
	h.AccountBalanceReport().InheritModel(h.AccountCommonAccountReport())

	h.AccountBalanceReport().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs h.AccountCommonAccountReportSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  data = self.pre_print_report(data)
			  records = self.env[data['model']].browse(data.get('ids', []))
			  return self.env['report'].get_action(records, 'account.report_trialbalance', data=data)
			*/
		})

}
