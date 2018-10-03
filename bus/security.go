//filevalid
package bus

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {
	h.BusPresence().Methods().AllowAllToGroup(base.GroupUser)
	h.BusPresence().Methods().AllowAllToGroup(base.GroupPortal)
}
