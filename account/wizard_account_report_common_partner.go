// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountCommonPartnerReport().DeclareMixinModel()
	h.AccountCommonPartnerReport().InheritModel(h.AccountCommonReport())

	h.AccountCommonPartnerReport().AddFields(map[string]models.FieldDefinition{
		"ResultSelection": models.SelectionField{ /*result_selection = fields.Selection([('customer', 'Receivable Accounts'), ('supplier', 'Payable Accounts'), ('customer_supplier', 'Receivable and Payable Accounts')*/ },
	})
	h.AccountCommonPartnerReport().Methods().PrePrintReport().DeclareMethod(
		`PrePrintReport`,
		func(rs h.AccountCommonPartnerReportSet, args struct {
			Data interface{}
		}) {
			/*def pre_print_report(self, data):
			  data['form'].update(self.read(['result_selection'])[0])
			  return data
			*/
		})

}
