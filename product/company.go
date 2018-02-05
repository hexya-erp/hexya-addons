// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.Company().AddFields(map[string]models.FieldDefinition{
		"DefaultPriceList": models.Many2OneField{RelationModel: h.ProductPricelist(),
			Help: "Default Price list for partners of this company"},
	})

	h.Company().Methods().Create().Extend("",
		func(rs h.CompanySet, vals *h.CompanyData) h.CompanySet {
			newCompany := rs.Super().Create(vals)
			priceList := h.ProductPricelist().Search(rs.Env(),
				q.ProductPricelist().Currency().Equals(newCompany.Currency()).And().Company().IsNull()).Limit(1)
			if priceList.IsEmpty() {
				priceList = h.ProductPricelist().Create(rs.Env(), &h.ProductPricelistData{
					Name:     newCompany.Name(),
					Currency: newCompany.Currency(),
				})
			}
			newCompany.SetDefaultPriceList(priceList)
			return newCompany
		})

	h.Company().Methods().Write().Extend("",
		func(rs h.CompanySet, vals *h.CompanyData, fieldsToUnset ...models.FieldNamer) bool {
			// When we modify the currency of the company, we reflect the change on the list0 pricelist, if
			// that pricelist is not used by another company. Otherwise, we create a new pricelist for the
			// given currency.
			currency := vals.Currency
			mainPricelist := h.ProductPricelist().Search(rs.Env(), q.ProductPricelist().HexyaExternalID().Equals("product_list0"))
			if currency.IsEmpty() || mainPricelist.IsEmpty() {
				return rs.Super().Write(vals, fieldsToUnset...)
			}
			nbCompanies := h.Company().NewSet(rs.Env()).SearchAll().SearchCount()
			for _, company := range rs.Records() {
				if mainPricelist.Company().Equals(company) || (mainPricelist.Company().IsEmpty() && nbCompanies == 1) {
					mainPricelist.SetCurrency(currency)
				} else {
					priceList := h.ProductPricelist().Create(rs.Env(), &h.ProductPricelistData{
						Name:     company.Name(),
						Currency: currency,
					})
					company.SetDefaultPriceList(priceList)
				}
			}
			return rs.Super().Write(vals, fieldsToUnset...)
		})

}
