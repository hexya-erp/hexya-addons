// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"math"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.SaleOrderLine().Methods().ComputeAnalytic().DeclareMethod(
		`ComputeAnalytic updates analytic lines linked with this SaleOrderLine`,
		func(rs pool.SaleOrderLineSet, cond pool.AccountAnalyticLineCondition) bool {
			lines := make(map[int64]float64)
			forceSOLines := rs.Env().Context().GetIntegerSlice("force_so_lines")
			if cond.IsEmpty() {
				if rs.IsEmpty() && len(forceSOLines) == 0 {
					return true
				}
				cond = pool.AccountAnalyticLine().SoLine().In(rs).And().Amount().LowerOrEqual(0)
			}
			data := pool.AccountAnalyticLine().Search(rs.Env(), cond).
				GroupBy(pool.AccountAnalyticLine().ProductUom(), pool.AccountAnalyticLine().SoLine()).
				Aggregates(pool.AccountAnalyticLine().ProductUom(), pool.AccountAnalyticLine().SoLine(),
					pool.AccountAnalyticLine().UnitAmount())
			for _, d := range data {
				pUom, _ := d.Values.Get("ProductUom", pool.AccountAnalyticLine().Underlying())
				soLineID, _ := d.Values.Get("SOLine", pool.AccountAnalyticLine().Underlying())
				unitAmount, _ := d.Values.Get("UnitAmount", pool.AccountAnalyticLine().Underlying())
				if pUom.(models.RecordSet).IsEmpty() {
					continue
				}
				line := pool.SaleOrderLine().Browse(rs.Env(), []int64{soLineID.(int64)})
				uom := pool.ProductUom().Browse(rs.Env(), []int64{pUom.(int64)})
				qty := unitAmount.(float64)
				if line.ProductUom().Category().Equals(uom.Category()) {
					qty = uom.ComputeQuantity(unitAmount.(float64), line.ProductUom(), true)
				}
				lines[line.ID()] += qty
			}
			for l, q := range lines {
				pool.SaleOrderLine().Browse(rs.Env(), []int64{l}).SetQtyDelivered(q)
			}
			return true
		})

	pool.AccountAnalyticLine().AddFields(map[string]models.FieldDefinition{
		"SoLine": models.Many2OneField{String: "Sale Order Line", RelationModel: pool.SaleOrderLine()},
	})

	pool.AccountAnalyticLine().Methods().GetInvoicePrice().DeclareMethod(
		`GetInvoicePrice returns the unit price to set on invoice`,
		func(rs pool.AccountAnalyticLineSet, order pool.SaleOrderSet) float64 {
			if rs.Product().ExpensePolicy() == "sales_price" {
				return rs.Product().
					WithContext("partner", order.Partner().ID()).
					WithContext("date_order", order.DateOrder()).
					WithContext("pricelist", order.Pricelist().ID()).
					WithContext("uom", rs.ProductUom().ID()).Price()
			}
			if rs.UnitAmount() == 0 {
				return 0
			}
			// Prevent unnecessary currency conversion that could be impacted by exchange rate
			// fluctuations
			if !rs.Currency().IsEmpty() && rs.AmountCurrency() != 0 && rs.Currency().Equals(order.Currency()) {
				return math.Abs(rs.AmountCurrency() / rs.UnitAmount())
			}
			priceUnit := math.Abs(rs.Amount() / rs.UnitAmount())
			currency := rs.Company().Currency()
			if !currency.IsEmpty() && !currency.Equals(order.Currency()) {
				priceUnit = currency.Compute(priceUnit, order.Currency(), true)
			}
			return priceUnit
		})

	pool.AccountAnalyticLine().Methods().GetSaleOrderLineVals().DeclareMethod(
		`GetSaleOrderLineVals returns the data to create a sale order line from this account analytic line on
		the given order for the given price.`,
		func(rs pool.AccountAnalyticLineSet, order pool.SaleOrderSet, price float64) *pool.SaleOrderLineData {
			lastSOLine := pool.SaleOrderLine().Search(rs.Env(), pool.SaleOrderLine().Order().Equals(order)).
				OrderBy("Sequence DESC").Limit(1)
			lastSequence := int64(100)
			if !lastSOLine.IsEmpty() {
				lastSequence = lastSOLine.Sequence() + 1
			}
			fPos := order.Partner().PropertyAccountPosition()
			if !order.FiscalPosition().IsEmpty() {
				fPos = order.FiscalPosition()
			}
			taxes := fPos.MapTax(rs.Product().Taxes(), rs.Product(), order.Partner())

			return &pool.SaleOrderLineData{
				Order:         order,
				Name:          rs.Name(),
				Sequence:      lastSequence,
				PriceUnit:     price,
				Tax:           taxes,
				Discount:      0,
				Product:       rs.Product(),
				ProductUom:    rs.ProductUom(),
				ProductUomQty: 0,
				QtyDelivered:  rs.UnitAmount(),
			}

		})

	pool.AccountAnalyticLine().Methods().GetSaleOrderLine().DeclareMethod(
		`GetSaleOrderLine adds the sale order line data to the given vals.
		Returned data is a modified copy of vals.`,
		func(rs pool.AccountAnalyticLineSet, vals *pool.AccountAnalyticLineData) *pool.AccountAnalyticLineData {
			result := *vals
			SOLine := result.SoLine
			if SOLine.IsEmpty() {
				SOLine = rs.SoLine()
			}
			if !SOLine.IsEmpty() || rs.Account().IsEmpty() || rs.Product().IsEmpty() || rs.Product().ExpensePolicy() == "no" {
				return &result
			}
			orderInSale := pool.SaleOrder().Search(rs.Env(),
				pool.SaleOrder().Project().Equals(rs.Account()).
					And().State().Equals("sale")).Limit(1)
			order := orderInSale
			if order.IsEmpty() {
				order = pool.SaleOrder().Search(rs.Env(), pool.SaleOrder().Project().Equals(rs.Account())).Limit(1)
			}
			if order.IsEmpty() {
				return &result
			}
			price := rs.GetInvoicePrice(order)
			SOLines := pool.SaleOrderLine().Search(rs.Env(),
				pool.SaleOrderLine().Order().Equals(order).
					And().PriceUnit().Equals(price).
					And().Product().Equals(rs.Product()))
			if !SOLines.IsEmpty() {
				result.SoLine = SOLines.Records()[0]
				return &result
			}
			if order.State() != "sale" {
				panic(rs.T("The Sale Order %s linked to the Analytic Account must be validated before registering expenses.", order.Name()))
			}
			orderLineVals := rs.GetSaleOrderLineVals(order, price)
			NewSOLine := pool.SaleOrderLine().Create(rs.Env(), orderLineVals)
			data, ftr := NewSOLine.ComputeTax()
			NewSOLine.Write(data, ftr...)
			result.SoLine = NewSOLine

			return &result
		})

	pool.AccountAnalyticLine().Methods().Write().Extend("",
		func(rs pool.AccountAnalyticLineSet, data *pool.AccountAnalyticLineData, fieldsToReset ...models.FieldNamer) bool {
			if rs.Env().Context().GetBool("create") {
				return rs.Super().Write(data, fieldsToReset...)
			}
			res := rs.Super().Write(data, fieldsToReset...)
			for _, line := range rs.Records() {
				vals := line.Sudo().GetSaleOrderLine(data)
				rs.Super().Write(vals, pool.AccountAnalyticLine().SoLine())
			}
			SOLines := pool.SaleOrderLine().NewSet(rs.Env())
			for _, line := range rs.Records() {
				SOLines = SOLines.Union(line.SoLine())
			}
			SOLines.ComputeAnalytic(pool.AccountAnalyticLineCondition{})
			return res
		})

	pool.AccountAnalyticLine().Methods().Create().Extend("",
		func(rs pool.AccountAnalyticLineSet, data *pool.AccountAnalyticLineData) pool.AccountAnalyticLineSet {
			line := rs.Super().Create(data)
			vals := line.Sudo().GetSaleOrderLine(data)
			line.WithContext("create", true).Write(vals, pool.AccountAnalyticLine().SoLine())
			SOLines := pool.SaleOrderLine().NewSet(rs.Env())
			for _, l := range rs.Records() {
				SOLines = SOLines.Union(l.SoLine())
			}
			SOLines.ComputeAnalytic(pool.AccountAnalyticLineCondition{})
			return line
		})

	pool.AccountAnalyticLine().Methods().Unlink().Extend("",
		func(rs pool.AccountAnalyticLineSet) int64 {
			SOLines := pool.SaleOrderLine().NewSet(rs.Env())
			for _, line := range rs.Records() {
				SOLines = SOLines.Union(line.SoLine())
			}
			res := rs.Super().Unlink()
			SOLines.WithContext("force_so_lines", SOLines.Ids()).ComputeAnalytic(pool.AccountAnalyticLineCondition{})
			return res
		})

}
