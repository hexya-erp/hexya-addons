// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProductPriceList(t *testing.T) {
	Convey("Testing Product Pricelist", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			//productPricelist := h.ProductPricelist().NewSet(env)
			partner4 := h.Partner().NewSet(env).GetRecord("base_res_partner_4")
			//computerSC234 := h.ProductProduct().NewSet(env).GetRecord("product_product_product_3")
			ipadRetinaDisplay := h.ProductProduct().NewSet(env).GetRecord("product_product_product_4")
			//customComputerKit := h.ProductProduct().NewSet(env).GetRecord("product_product_product_5")
			ipadMini := h.ProductProduct().NewSet(env).GetRecord("product_product_product_6")
			appleInEarHeadphones := h.ProductProduct().NewSet(env).GetRecord("product_product_product_7")
			lapTopE5023 := h.ProductProduct().NewSet(env).GetRecord("product_product_delivery_01")
			lapTopS3450 := h.ProductProduct().NewSet(env).GetRecord("product_product_product_25")
			category5 := h.ProductCategory().NewSet(env).GetRecord("product_product_category_5")
			uomUnit := h.ProductUom().NewSet(env).GetRecord("product_product_uom_unit")
			list0 := h.ProductPricelist().NewSet(env).GetRecord("product_list0")

			ipadRetinaDisplay.Write(&h.ProductProductData{
				Uom:      uomUnit,
				Category: category5,
			})
			customerPricelist := h.ProductPricelist().Create(env, &h.ProductPricelistData{
				Name: "Customer pricelist",
				Items: h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
					ComputePrice:  "formula",
					Base:          "pricelist",
					BasePricelist: list0,
				}).Union(h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
					AppliedOn:     "1_product",
					Sequence:      1,
					Product:       ipadRetinaDisplay,
					ComputePrice:  "formula",
					Base:          "ListPrice",
					PriceDiscount: 10,
				})).Union(h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
					AppliedOn:      "1_product",
					Sequence:       4,
					Product:        lapTopE5023,
					ComputePrice:   "formula",
					Base:           "ListPrice",
					PriceSurcharge: 1,
				})).Union(h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
					AppliedOn:     "2_product_category",
					Sequence:      1,
					MinQuantity:   2,
					ComputePrice:  "formula",
					Base:          "ListPrice",
					Category:      category5,
					PriceDiscount: 5,
				})).Union(h.ProductPricelistItem().Create(env, &h.ProductPricelistItemData{
					AppliedOn:     "0_product_variant",
					DateStart:     dates.ParseDate("2011-12-27"),
					DateEnd:       dates.ParseDate("2011-12-31"),
					Sequence:      1,
					ComputePrice:  "formula",
					Base:          "ListPrice",
					PriceDiscount: 30,
				})),
			})
			Convey("Test calculation of product price based on pricelist", func() {
				context := types.NewContext().
					WithKey("pricelist", customerPricelist.ID()).
					WithKey("quantity", 1)

				ipadRetinaDisplay = ipadRetinaDisplay.WithNewContext(context)
				So(ipadRetinaDisplay.Price(), ShouldAlmostEqual, ipadRetinaDisplay.LstPrice()*0.9, 0.01)
				So(ipadRetinaDisplay.Price(), ShouldAlmostEqual, 675, 0.01)

				lapTopE5023 = lapTopE5023.WithNewContext(context)
				So(lapTopE5023.Price(), ShouldAlmostEqual, lapTopE5023.LstPrice()+1, 0.01)
				So(lapTopE5023.Price(), ShouldAlmostEqual, 71, 0.01)

				appleHeadPhones := appleInEarHeadphones.WithNewContext(context)
				So(appleHeadPhones.Price(), ShouldAlmostEqual, appleHeadPhones.LstPrice(), 0.01)
				So(appleHeadPhones.Price(), ShouldAlmostEqual, 79, 0.01)

				context = context.WithKey("quantity", 5)
				lapTopS3450 = lapTopS3450.WithNewContext(context)
				So(lapTopS3450.Price(), ShouldAlmostEqual, lapTopS3450.LstPrice()*0.95, 0.01)
				So(lapTopS3450.Price(), ShouldAlmostEqual, 2802.5, 0.01)

				context = context.WithKey("quantity", 1)
				ipadMini = ipadMini.WithNewContext(context)
				So(ipadMini.Price(), ShouldAlmostEqual, ipadMini.LstPrice(), 0.01)
				So(ipadMini.Price(), ShouldAlmostEqual, 320, 0.01)

				context = context.
					WithKey("quantity", 1).
					WithKey("date", dates.ParseDate("2011-12-31"))
				ipadMini = ipadMini.WithNewContext(context)
				So(ipadMini.Price(), ShouldAlmostEqual, ipadMini.LstPrice()*0.7, 0.01)
				So(ipadMini.Price(), ShouldAlmostEqual, 224, 0.01)

				context = context.
					WithKey("quantity", 1).
					WithKey("date", dates.Date{}).
					WithKey("partner_id", partner4.ID())
				ipadMini = ipadMini.WithNewContext(context)
				partner := partner4.WithNewContext(context)
				So(ipadMini.SelectSeller(partner, 1, dates.Date{}, h.ProductUom().NewSet(env)).Price(), ShouldAlmostEqual, 790, 0.01)

				context = context.
					WithKey("quantity", 3)
				ipadMini = ipadMini.WithNewContext(context)
				partner = partner4.WithNewContext(context)
				So(ipadMini.SelectSeller(partner, 3, dates.Date{}, h.ProductUom().NewSet(env)).Price(), ShouldAlmostEqual, 785, 0.01)

			})
		}), ShouldBeNil)
	})
}
