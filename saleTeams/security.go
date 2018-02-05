// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool/h"
)

var (
	GroupSaleSalesman         *security.Group
	GroupSaleManager          *security.Group
	GroupSaleSalesmanAllLeads *security.Group
)

func init() {

	h.CRMTeam().Methods().Load().AllowGroup(base.GroupUser)
	h.CRMTeam().Methods().Load().AllowGroup(GroupSaleSalesman)
	h.CRMTeam().Methods().AllowAllToGroup(GroupSaleManager)
}
