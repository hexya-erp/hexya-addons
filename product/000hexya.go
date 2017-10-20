// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	// product module dependencies
	_ "github.com/hexya-erp/hexya-addons/decimalPrecision"
	_ "github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/server"
)

const MODULE_NAME string = "product"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}
