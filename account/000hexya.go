// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	// Import dependencies
	_ "github.com/hexya-erp/hexya-addons/analytic"
	"github.com/hexya-erp/hexya-base/base"
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
}
