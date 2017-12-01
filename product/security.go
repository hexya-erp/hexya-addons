// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool"
)

var (
	// GroupSalePriceList enables the "Sales Pricelists" feature
	GroupSalePriceList *security.Group
	// GroupPricelistItem enables the "Manage Pricelist Items" feature
	GroupPricelistItem *security.Group
	// GroupProductPricelist enables the "Pricelists On Product" feature
	GroupProductPricelist *security.Group
	// GroupUom enables the "Manage Multiple Units of Measure" feature
	GroupUom *security.Group
	// GroupStockPackaging enables the "Manage Product Packaging" feature
	GroupStockPackaging *security.Group
	// GroupMRPProperties enables the "Manage Properties of Product" feature
	GroupMRPProperties *security.Group
	// GroupProductVariant enables the "Manage Product Variants" feature
	GroupProductVariant *security.Group
)

func init() {
	pool.ProductUomCateg().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductUom().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductCategory().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductTemplate().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductPackaging().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductSupplierinfo().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductPricelist().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductPricelistItem().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductPricelist().Methods().Load().AllowGroup(base.GroupPartnerManager)
	pool.ProductProduct().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductPriceHistory().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductAttribute().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductAttributeValue().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductAttributePrice().Methods().Load().AllowGroup(base.GroupUser)
	pool.ProductAttributeLine().Methods().Load().AllowGroup(base.GroupUser)
}
