// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/tests"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func TestMain(m *testing.M) {
	tests.RunTests(m, "product")
}

type productTestData struct {
	partner1    h.PartnerSet
	uomUnit     h.ProductUomSet
	uomDozen    h.ProductUomSet
	uomDunit    h.ProductUomSet
	uomWeight   h.ProductUomSet
	product0    h.ProductProductSet
	product1    h.ProductProductSet
	product2    h.ProductProductSet
	product3    h.ProductProductSet
	product4    h.ProductProductSet
	product5    h.ProductProductSet
	product6    h.ProductProductSet
	product7    h.ProductProductSet
	product71   h.ProductProductSet
	product72   h.ProductProductSet
	product8    h.ProductProductSet
	product9    h.ProductProductSet
	product10   h.ProductProductSet
	template7   h.ProductTemplateSet
	prodAtt1    h.ProductAttributeSet
	prodAttr1V1 h.ProductAttributeValueSet
	prodAttr1V2 h.ProductAttributeValueSet
}

func getProductTestData(env models.Environment) *productTestData {
	var ptd productTestData
	ptd.partner1 = h.Partner().Create(env, &h.PartnerData{
		Name:  "Julia Agrolait",
		Email: "julia@agrolait.example.com",
	})
	ptd.uomUnit = h.ProductUom().Search(env, q.ProductUom().HexyaExternalID().Equals("product_product_uom_unit"))
	ptd.uomDozen = h.ProductUom().Search(env, q.ProductUom().HexyaExternalID().Equals("product_product_uom_dozen"))
	ptd.uomDunit = h.ProductUom().Create(env, &h.ProductUomData{
		Name:      "DeciUnit",
		Category:  ptd.uomUnit.Category(),
		FactorInv: 0.1,
		Factor:    10,
		UomType:   "smaller",
		Rounding:  0.001,
	})
	ptd.uomWeight = h.ProductUom().Search(env, q.ProductUom().HexyaExternalID().Equals("product_product_uom_kgm"))
	ptd.product0 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Work",
		Type:  "service",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	ptd.product1 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Courage",
		Type:  "consu",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomDunit,
	})
	ptd.product2 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Wood",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	ptd.product3 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Stone",
		Uom:   ptd.uomDozen,
		UomPo: ptd.uomDozen,
	})
	ptd.product4 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Stick",
		Uom:   ptd.uomDozen,
		UomPo: ptd.uomDozen,
	})
	ptd.product5 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Stone Tools",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	ptd.product6 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Door",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	ptd.prodAtt1 = h.ProductAttribute().Create(env, &h.ProductAttributeData{
		Name: "Color",
	})
	ptd.prodAttr1V1 = h.ProductAttributeValue().Create(env, &h.ProductAttributeValueData{
		Name:      "Red",
		Attribute: ptd.prodAtt1,
	})
	ptd.prodAttr1V2 = h.ProductAttributeValue().Create(env, &h.ProductAttributeValueData{
		Name:      "Blue",
		Attribute: ptd.prodAtt1,
	})
	ptd.template7 = h.ProductTemplate().Create(env, &h.ProductTemplateData{
		Name:  "Sofa",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
		AttributeLines: h.ProductAttributeLine().Create(env, &h.ProductAttributeLineData{
			Attribute: ptd.prodAtt1,
		}),
	})
	ptd.product7 = h.ProductProduct().Create(env, &h.ProductProductData{
		ProductTmpl: ptd.template7,
	})
	ptd.product71 = h.ProductProduct().Create(env, &h.ProductProductData{
		ProductTmpl:     ptd.template7,
		AttributeValues: ptd.prodAttr1V1,
	})
	ptd.product72 = h.ProductProduct().Create(env, &h.ProductProductData{
		ProductTmpl:     ptd.template7,
		AttributeValues: ptd.prodAttr1V2,
	})
	ptd.product8 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "House",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	ptd.product9 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Paper",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	ptd.product10 = h.ProductProduct().Create(env, &h.ProductProductData{
		Name:  "Stone",
		Uom:   ptd.uomUnit,
		UomPo: ptd.uomUnit,
	})
	return &ptd
}
