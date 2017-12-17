// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProductProduct().AddFields(map[string]models.FieldDefinition{
		"SalesCount": models.IntegerField{String: "# Sales", Compute: pool.ProductProduct().Methods().ComputeSalesCount(),
			GoType: new(int)},
	})

	pool.ProductProduct().Methods().ComputeSalesCount().DeclareMethod(
		`ComputeSalesCount returns the number of sales for this product`,
		func(rs pool.ProductProductSet) (*pool.ProductProductData, []models.FieldNamer) {
			cond := pool.SaleReport().State().In([]string{"sale", "done"}).And().Product().In(rs)
			return &pool.ProductProductData{
				SalesCount: pool.SaleReport().NewSet(rs.Env()).Search(cond).SearchCount(),
			}, []models.FieldNamer{pool.ProductProduct().SalesCount()}
		})
}
