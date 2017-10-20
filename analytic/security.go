// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package analytic

import "github.com/hexya-erp/hexya/pool"

func init() {
	pool.AccountAnalyticAccount().Methods().AllowAllToGroup(GroupAnalyticAccounting)
	pool.AccountAnalyticLine().Methods().AllowAllToGroup(GroupAnalyticAccounting)
	pool.AccountAnalyticTag().Methods().AllowAllToGroup(GroupAnalyticAccounting)
}
