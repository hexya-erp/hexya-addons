// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.ProductProduct().AddFields(map[string]models.FieldDefinition{
		"SalesCount": models.IntegerField{String: "# Sales", Compute: h.ProductProduct().Methods().ComputeSalesCount(),
			GoType: new(int)},
	})

	h.ProductProduct().Methods().ComputeSalesCount().DeclareMethod(
		`ComputeSalesCount returns the number of sales for this product`,
		func(rs h.ProductProductSet) *h.ProductProductData {
			cond := q.SaleReport().State().In([]string{"sale", "done"}).And().Product().In(rs)
			return &h.ProductProductData{
				SalesCount: h.SaleReport().NewSet(rs.Env()).Search(cond).SearchCount(),
			}
		})
}
