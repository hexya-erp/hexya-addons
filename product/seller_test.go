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

func TestSeller(t *testing.T) {
	Convey("Testing sellers", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			productService := h.ProductProduct().NewSet(env).GetRecord("product_product_product_2")
			productService.SetDefaultCode("DEFCODE")
			asusTec := h.Partner().NewSet(env).GetRecord("base_res_partner_1")
			campToCamp := h.Partner().NewSet(env).GetRecord("base_res_partner_12")
			Convey("Product codes should match context", func() {
				productService.SetSellers(h.ProductSupplierinfo().Create(env, &h.ProductSupplierinfoData{
					Name:        asusTec,
					ProductCode: "ASUCODE",
				}).Union(h.ProductSupplierinfo().Create(env, &h.ProductSupplierinfoData{
					Name:        campToCamp,
					ProductCode: "C2CCODE",
				})))
				defaultCode := productService.Code()
				So(defaultCode, ShouldEqual, "DEFCODE")
				contextCode := productService.WithContext("partner_id", campToCamp.ID()).Code()
				So(contextCode, ShouldEqual, "C2CCODE")
			})
		}), ShouldBeNil)
	})
}
