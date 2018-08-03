// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool/h"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUom(t *testing.T) {
	Convey("Testing UoMs", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			uomGram := h.ProductUom().NewSet(env).GetRecord("product_product_uom_gram")
			uomKgm := h.ProductUom().NewSet(env).GetRecord("product_product_uom_kgm")
			uomTon := h.ProductUom().NewSet(env).GetRecord("product_product_uom_ton")
			uomUnit := h.ProductUom().NewSet(env).GetRecord("product_product_uom_unit")
			uomDozen := h.ProductUom().NewSet(env).GetRecord("product_product_uom_dozen")
			categUnit := h.ProductUomCategory().NewSet(env).GetRecord("product_product_uom_categ_unit")
			Convey("Conversions", func() {
				qty := uomGram.ComputeQuantity(1020000, uomTon, true)
				So(qty, ShouldEqual, 1.02)

				price := uomGram.ComputePrice(2, uomTon)
				So(price, ShouldEqual, 2000000.0)

				qty = uomDozen.ComputeQuantity(1, uomUnit, true)
				So(qty, ShouldEqual, 12.0)

				uomGram.SetRounding(1)
				qty = uomGram.ComputeQuantity(1234, uomKgm, true)
				So(qty, ShouldEqual, 1.234)
			})
			Convey("Testing Roundings", func() {
				productUom := h.ProductUom().Create(env, &h.ProductUomData{
					Name:      "Score",
					FactorInv: 20,
					UomType:   "bigger",
					Rounding:  1.0,
					Category:  categUnit,
				})
				qty := uomUnit.ComputeQuantity(2, productUom, true)
				// Unlike Odoo, we do not want to go into rounding issues with epsilons.
				So(qty, ShouldEqual, 0)
			})
		}), ShouldBeNil)
	})
}
