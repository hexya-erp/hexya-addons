// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProcurementOrder().AddFields(map[string]models.FieldDefinition{
		"SaleLine": models.Many2OneField{String: "Sale Order Line", RelationModel: pool.SaleOrderLine()},
	})

}
