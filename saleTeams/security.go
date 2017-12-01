// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool"
)

var (
	GroupSaleSalesman         *security.Group
	GroupSaleManager          *security.Group
	GroupSaleSalesmanAllLeads *security.Group
)

func init() {

	pool.CRMTeam().Methods().Load().AllowGroup(base.GroupUser)
	pool.CRMTeam().Methods().Load().AllowGroup(GroupSaleSalesman)
	pool.CRMTeam().Methods().AllowAllToGroup(GroupSaleManager)
}
