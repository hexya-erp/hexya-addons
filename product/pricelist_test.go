// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool/h"
	. "github.com/smartystreets/goconvey/convey"
)

type priceListTestData struct {
	dataCard        h.ProductProductSet
	usbAdapter      h.ProductProductSet
	uomTon          h.ProductUomSet
	uomUnit         h.ProductUomSet
	uomDozen        h.ProductUomSet
	uomKgm          h.ProductUomSet
	publicPriceList h.ProductPricelistSet
	salePriceList   h.ProductPricelistSet
}

func getTestPriceListData(env models.Environment) *priceListTestData {
	pltd := &priceListTestData{
		dataCard:        h.ProductProduct().NewSet(env).GetRecord("product_product_delivery_02"),
		usbAdapter:      h.ProductProduct().NewSet(env).GetRecord("product_product_delivery_01"),
		uomTon:          h.ProductUom().NewSet(env).GetRecord("product_product_uom_ton"),
		uomUnit:         h.ProductUom().NewSet(env).GetRecord("product_product_uom_unit"),
		uomDozen:        h.ProductUom().NewSet(env).GetRecord("product_product_uom_dozen"),
		uomKgm:          h.ProductUom().NewSet(env).GetRecord("product_product_uom_kgm"),
		publicPriceList: h.ProductPricelist().NewSet(env).GetRecord("product_list0"),
	}
	pltd.salePriceList = h.ProductPricelist().Create(env, &h.ProductPricelistData{
		Name: "Sale pricelist",
		Items: h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
			ComputePrice:  "formula",
			Base:          "ListPrice",
			PriceDiscount: 10,
			Product:       pltd.usbAdapter,
			AppliedOn:     "0_product_variant",
		}).Union(h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
			ComputePrice:   "formula",
			Base:           "ListPrice",
			PriceSurcharge: -0.5,
			Product:        pltd.dataCard,
			AppliedOn:      "0_product_variant",
		})),
	})
	return pltd
}

func TestPriceList(t *testing.T) {
	Convey("Testing Price lists", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			Convey("Test Discount", func() {
				pltd := getTestPriceListData(env)
				publicContext := types.NewContext().WithKey("pricelist", pltd.publicPriceList.ID())
				pricelistContext := types.NewContext().WithKey("pricelist", pltd.salePriceList.ID())

				usbAdapterWithoutPriceList := pltd.usbAdapter.WithNewContext(publicContext)
				usbAdapterWithPriceList := pltd.usbAdapter.WithNewContext(pricelistContext)
				So(usbAdapterWithPriceList.Price(), ShouldEqual, 63)
				So(usbAdapterWithoutPriceList.Price(), ShouldEqual, 70)
				So(usbAdapterWithPriceList.Price(), ShouldEqual, usbAdapterWithoutPriceList.Price()*0.9)

				dataCardWithoutPriceList := pltd.dataCard.WithNewContext(publicContext)
				dataCardWithPriceList := pltd.dataCard.WithNewContext(pricelistContext)
				So(dataCardWithPriceList.Price(), ShouldEqual, 39.5)
				So(dataCardWithoutPriceList.Price(), ShouldEqual, 40)
				So(dataCardWithPriceList.Price(), ShouldEqual, dataCardWithoutPriceList.Price()-0.5)

				// Make sure that changing the unit of measure does not break the unit price (after converting)
				unitContext := types.NewContext().WithKey("pricelist", pltd.salePriceList.ID()).WithKey("uom", pltd.uomUnit.ID())
				dozenContext := types.NewContext().WithKey("pricelist", pltd.salePriceList.ID()).WithKey("uom", pltd.uomDozen.ID())
				usbAdapterUnit := pltd.usbAdapter.WithNewContext(unitContext)
				usbAdapterDozen := pltd.usbAdapter.WithNewContext(dozenContext)
				So(usbAdapterUnit.Price()*12, ShouldAlmostEqual, usbAdapterDozen.Price(), .000000001)
				dataCardUnit := pltd.dataCard.WithNewContext(unitContext)
				dataCardDozen := pltd.dataCard.WithNewContext(dozenContext)
				So(dataCardUnit.Price()*12, ShouldAlmostEqual, dataCardDozen.Price(), .000000001)
			})
		}), ShouldBeNil)
	})
}
