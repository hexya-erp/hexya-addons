// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool/h"
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
	h.ProductUomCateg().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductUom().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductCategory().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductTemplate().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductPackaging().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductSupplierinfo().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductPricelist().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductPricelistItem().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductPricelist().Methods().Load().AllowGroup(base.GroupPartnerManager)
	h.ProductProduct().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductPriceHistory().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductAttribute().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductAttributeValue().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductAttributePrice().Methods().Load().AllowGroup(base.GroupUser)
	h.ProductAttributeLine().Methods().Load().AllowGroup(base.GroupUser)
}
