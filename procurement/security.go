// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package procurement

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.ProcurementOrder().Methods().AllowAllToGroup(base.GroupUser)
	h.ProcurementGroup().Methods().AllowAllToGroup(base.GroupUser)
	h.ProcurementRule().Methods().AllowAllToGroup(base.GroupUser)

}
