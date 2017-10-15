package decimalPrecision

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool"
)

func init() {
	pool.DecimalPrecision().Methods().AllowAllToGroup(base.GroupSystem)
	pool.DecimalPrecision().Methods().Load().AllowGroup(security.GroupEveryone)
}
