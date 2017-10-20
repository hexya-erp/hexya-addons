// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.Partner().AddFields(map[string]models.FieldDefinition{
		"PropertyProductPricelist": models.Many2OneField{String: "Sale Pricelist", RelationModel: pool.ProductPricelist(),
			Compute: pool.Partner().Methods().ComputeProductPricelist(),
			Depends: []string{"Country"},
			Inverse: pool.Partner().Methods().InverseProductPricelist(),
			Help:    "This pricelist will be used instead of the default one for sales to the current partner"},
		"ProductPricelist": models.Many2OneField{String: "Stored Pricelist", RelationModel: pool.ProductPricelist()},
	})

	pool.Partner().Methods().ComputeProductPricelist().DeclareMethod(
		`ComputeProductPricelist returns the price list applicable for this partner`,
		func(rs pool.PartnerSet) (*pool.PartnerData, []models.FieldNamer) {
			if rs.ID() == 0 {
				// We are processing an Onchange
				return new(pool.PartnerData), []models.FieldNamer{}
			}
			company := pool.User().NewSet(rs.Env()).CurrentUser().Company()
			return &pool.PartnerData{
				PropertyProductPricelist: pool.ProductPricelist().NewSet(rs.Env()).GetPartnerPricelist(rs, company),
			}, []models.FieldNamer{pool.Partner().PropertyProductPricelist()}
		})

	pool.Partner().Methods().InverseProductPricelist().DeclareMethod(
		`InverseProductPricelist sets the price list for this partner to the given list`,
		func(rs pool.PartnerSet, priceList pool.ProductPricelistSet) {
			var defaultForCountry pool.ProductPricelistSet
			if !rs.Country().IsEmpty() {
				defaultForCountry = pool.ProductPricelist().Search(rs.Env(),
					pool.ProductPricelist().CountryGroupsFilteredOn(
						pool.CountryGroup().CountriesFilteredOn(
							pool.Country().Code().Equals(rs.Country().Code())))).Limit(1)
			} else {
				defaultForCountry = pool.ProductPricelist().Search(rs.Env(),
					pool.ProductPricelist().CountryGroups().IsNull()).Limit(1)
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

	pool.Partner().Methods().CommercialFields().Extend(
		`CommercialFields`,
		func(rs pool.PartnerSet) []models.FieldNamer {
			return append(rs.Super().CommercialFields(), pool.Partner().PropertyProductPricelist())
		})
}
