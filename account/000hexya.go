// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	// Import dependencies
	_ "github.com/hexya-erp/hexya-addons/analytic"
	"github.com/hexya-erp/hexya/hexya/server"
)

const MODULE_NAME string = "account"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}
