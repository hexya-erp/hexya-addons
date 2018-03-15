// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"fmt"

	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.ProductTemplate().AddFields(map[string]models.FieldDefinition{
		"TrackService": models.SelectionField{String: "Track Service", Selection: types.Selection{
			"manual": "Manually set quantities on order",
		}, Help: `Manually set quantities on order: Invoice based on the manually entered quantity without creating an
analytic account.
Timesheets on contract: Invoice based on the tracked hours on the related timesheet.
Create a task and track hours: Create a task on the sale order validation and track the work hours.`,
			Default: models.DefaultValue("manual")},
		"SaleLineWarn": models.SelectionField{Selection: base.WarningMessage, String: "Sales Order Line",
			/*Help: base.WarningHelp,*/ Required: true, Default: models.DefaultValue("no-message")},
		"SaleLineWarnMsg": models.TextField{String: "Message for Sales Order Line')"},
		"ExpensePolicy": models.SelectionField{String: "Re-Invoice Expenses", Selection: types.Selection{
			"no":          "No",
			"cost":        "At cost",
			"sales_price": "At sale price",
		}, Default: models.DefaultValue("no")},
		"SalesCount": models.IntegerField{String: "# Sales",
			Compute: h.ProductTemplate().Methods().ComputeSalesCount(), GoType: new(int)},
		"InvoicePolicy": models.SelectionField{String: "Invoicing Policy", Selection: types.Selection{
			"order":    "Ordered quantities",
			"delivery": "Delivered quantities",
		}, Help: `Ordered Quantity: Invoice based on the quantity the customer ordered.
Delivered Quantity: Invoiced based on the quantity the vendor delivered (time or deliveries).`,
			Default: models.DefaultValue("order")},
	})

	h.ProductTemplate().Methods().ComputeSalesCount().DeclareMethod(
		`ComputeSalesCount returns the number of sales for this product template.`,
		func(rs h.ProductTemplateSet) *h.ProductTemplateData {
			var count int
			for _, product := range rs.ProductVariants().Records() {
				count += product.SalesCount()
			}
			return &h.ProductTemplateData{
				SalesCount: count,
			}
		})

	h.ProductTemplate().Methods().ActionViewSales().DeclareMethod(
		`ActionViewSales`,
		func(rs h.ProductTemplateSet) *actions.Action {
			rs.EnsureOne()
			action := actions.Registry.MustGetById("sale_action_product_sale_list")
			products := rs.WithContext("active_test", false).ProductVariants()
			returnedAction := *action
			returnedAction.Context = types.NewContext(map[string]interface{}{
				"default_product_id": products.Ids()[0],
			})
			returnedAction.Domain = fmt.Sprintf("[('state', 'in', ['sale', 'done']), ('product_id.product_tmpl_id', '=', %d)]", rs.ID())
			return &returnedAction
		})

}
