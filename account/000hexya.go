// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	// Import dependencies
	_ "github.com/hexya-erp/hexya-addons/analytic"
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya-base/web/controllers"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/server"
)

const MODULE_NAME string = "account"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})

	GroupAccountInvoice = security.Registry.NewGroup("account_group_account_invoice", "Billing", base.GroupUser)
	GroupAccountUser = security.Registry.NewGroup("account_group_account_user", "Accountant", GroupAccountInvoice)
	GroupAccountManager = security.Registry.NewGroup("account_group_account_manager", "Adviser", GroupAccountUser)

	controllers.BackendCSS = append(controllers.BackendCSS,
		"/static/account/src/css/account_bank_and_cash.css",
		"/static/account/src/css/account.css")
	controllers.BackendLess = append(controllers.BackendLess,
		"/static/account/src/less/account_reconciliation.less",
		"/static/account/src/less/account_journal_dashboard.less")
	controllers.BackendJS = append(controllers.BackendJS,
		"/static/account/src/js/account_reconciliation_widgets.js",
		"/static/account/src/js/move_line_quickadd.js",
		"/static/account/src/js/account_payment_widget.js",
		"/static/account/src/js/account_journal_dashboard_widget.js")
}
