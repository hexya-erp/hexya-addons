// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.Partner().AddFields(map[string]models.FieldDefinition{
		"SaleOrderCount": models.IntegerField{String: "# of Sales Order",
			Compute: pool.Partner().Methods().ComputeSaleOrderCount(), GoType: new(int)},
		"SaleOrders": models.One2ManyField{String: "Sales Order", RelationModel: pool.SaleOrder(),
			ReverseFK: "Partner", JSON: "sale_order_ids"},
		"SaleWarn": models.SelectionField{Selection: base.WarningMessage, String: "Sales Order",
			Default: models.DefaultValue("no-message") /* Help: base.WarningHelp */, Required: true},
		"SaleWarnMsg": models.TextField{String: "Message for Sales Order"},
	})

	pool.Partner().Methods().ComputeSaleOrderCount().DeclareMethod(
		`ComputeSaleOrderCount`,
		func(rs pool.PartnerSet) (*pool.PartnerData, []models.FieldNamer) {
			count := pool.SaleOrder().Search(rs.Env(), pool.SaleOrder().Partner().ChildOf(rs)).SearchCount()
			return &pool.PartnerData{
				SaleOrderCount: count,
			}, []models.FieldNamer{pool.Partner().SaleOrderCount()}
		})

}
