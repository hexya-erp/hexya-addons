package decimalPrecision

import (
	// decimalPrecision depends on base module
	_ "github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/server"
)

const MODULE_NAME string = "decimalPrecision"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}
