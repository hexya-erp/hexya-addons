// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.Partner().AddFields(map[string]models.FieldDefinition{
		"PropertyProductPricelist": models.Many2OneField{String: "Sale Pricelist", RelationModel: h.ProductPricelist(),
			Compute: h.Partner().Methods().ComputeProductPricelist(),
			Depends: []string{"Country"},
			Inverse: h.Partner().Methods().InverseProductPricelist(),
			Help:    "This pricelist will be used instead of the default one for sales to the current partner"},
		"ProductPricelist": models.Many2OneField{String: "Stored Pricelist", RelationModel: h.ProductPricelist()},
	})

	h.Partner().Methods().ComputeProductPricelist().DeclareMethod(
		`ComputeProductPricelist returns the price list applicable for this partner`,
		func(rs h.PartnerSet) *h.PartnerData {
			if rs.ID() == 0 {
				// We are processing an Onchange
				return new(h.PartnerData)
			}
			company := h.User().NewSet(rs.Env()).CurrentUser().Company()
			return &h.PartnerData{
				PropertyProductPricelist: h.ProductPricelist().NewSet(rs.Env()).GetPartnerPricelist(rs, company),
			}
		})

	h.Partner().Methods().InverseProductPricelist().DeclareMethod(
		`InverseProductPricelist sets the price list for this partner to the given list`,
		func(rs h.PartnerSet, priceList h.ProductPricelistSet) {
			var defaultForCountry h.ProductPricelistSet
			if !rs.Country().IsEmpty() {
				defaultForCountry = h.ProductPricelist().Search(rs.Env(),
					q.ProductPricelist().CountryGroupsFilteredOn(
						q.CountryGroup().CountriesFilteredOn(
							q.Country().Code().Equals(rs.Country().Code())))).Limit(1)
			} else {
				defaultForCountry = h.ProductPricelist().Search(rs.Env(),
					q.ProductPricelist().CountryGroups().IsNull()).Limit(1)
			}
			actual := rs.PropertyProductPricelist()
			if !priceList.IsEmpty() || (!actual.IsEmpty() && !defaultForCountry.Equals(actual)) {
				if priceList.IsEmpty() {
					rs.SetProductPricelist(defaultForCountry)
					return
				}
				rs.SetProductPricelist(priceList)
			}
		})

	h.Partner().Methods().CommercialFields().Extend(
		`CommercialFields`,
		func(rs h.PartnerSet) []models.FieldNamer {
			return append(rs.Super().CommercialFields(), h.Partner().PropertyProductPricelist())
		})
}
