// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.ProductProduct().Fields().SalesCount().SetCompute(h.ProductProduct().Methods().ComputeSalesCount())

	h.ProductProduct().Methods().ComputeSalesCount().DeclareMethod(
		`ComputeSalesCount returns the number of sales for this product`,
		func(rs h.ProductProductSet) (*h.ProductProductData, []models.FieldNamer) {
			cond := q.SaleReport().State().In([]string{"sale", "done"}).And().Product().In(rs)
			return &h.ProductProductData{
				SalesCount: h.SaleReport().NewSet(rs.Env()).Search(cond).SearchCount(),
			}, []models.FieldNamer{h.ProductProduct().SalesCount()}
		})
}
