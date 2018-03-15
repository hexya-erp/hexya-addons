// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.Partner().AddFields(map[string]models.FieldDefinition{
		"SaleOrderCount": models.IntegerField{String: "# of Sales Order",
			Compute: h.Partner().Methods().ComputeSaleOrderCount(), GoType: new(int)},
		"SaleOrders": models.One2ManyField{String: "Sales Order", RelationModel: h.SaleOrder(),
			ReverseFK: "Partner", JSON: "sale_order_ids"},
		"SaleWarn": models.SelectionField{Selection: base.WarningMessage, String: "Sales Order",
			Default: models.DefaultValue("no-message") /* Help: base.WarningHelp */, Required: true},
		"SaleWarnMsg": models.TextField{String: "Message for Sales Order"},
	})

	h.Partner().Methods().ComputeSaleOrderCount().DeclareMethod(
		`ComputeSaleOrderCount`,
		func(rs h.PartnerSet) *h.PartnerData {
			count := h.SaleOrder().Search(rs.Env(), q.SaleOrder().Partner().ChildOf(rs)).SearchCount()
			return &h.PartnerData{
				SaleOrderCount: count,
			}
		})

}
