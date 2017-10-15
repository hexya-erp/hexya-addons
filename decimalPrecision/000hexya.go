package decimalPrecision

import (
	// decimalPrecision depends on base module
	_ "github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/server"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool"
)

const MODULE_NAME string = "decimalPrecision"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}

// GetPrecision returns the precision for the given application
func GetPrecision(app string) nbutils.Digits {
	var res int8
	models.ExecuteInNewEnvironment(security.SuperUserID, func(env models.Environment) {
		res = pool.DecimalPrecision().NewSet(env).PrecisionGet(app)
	})
	return nbutils.Digits{Precision: 16, Scale: res}
}
