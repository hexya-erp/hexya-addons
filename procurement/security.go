package procurement

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProcurementOrder().Methods().AllowAllToGroup(base.GroupUser)
	pool.ProcurementGroup().Methods().AllowAllToGroup(base.GroupUser)
	pool.ProcurementRule().Methods().AllowAllToGroup(base.GroupUser)

}
