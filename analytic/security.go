// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package analytic

import "github.com/hexya-erp/hexya/pool/h"

func init() {
	h.AccountAnalyticAccount().Methods().AllowAllToGroup(GroupAnalyticAccounting)
	h.AccountAnalyticLine().Methods().AllowAllToGroup(GroupAnalyticAccounting)
	h.AccountAnalyticTag().Methods().AllowAllToGroup(GroupAnalyticAccounting)
}
