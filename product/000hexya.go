// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	// product module dependencies
	_ "github.com/hexya-erp/hexya-addons/decimalPrecision"
	_ "github.com/hexya-erp/hexya-addons/webKanban"
	_ "github.com/hexya-erp/hexya-base/web"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/server"
)

const MODULE_NAME string = "product"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})

	GroupSalePriceList = security.Registry.NewGroup("product_group_sale_pricelist", "Sales Pricelists")
	GroupPricelistItem = security.Registry.NewGroup("product_group_pricelist_item", "Manage Pricelist Items")
	GroupProductPricelist = security.Registry.NewGroup("product_group_product_pricelist", "Pricelists On Product")
	GroupUom = security.Registry.NewGroup("product_group_uom", "Manage Multiple Units of Measure")
	GroupStockPackaging = security.Registry.NewGroup("product_group_stock_packaging", "Manage Product Packaging")
	GroupMRPProperties = security.Registry.NewGroup("product_group_mrp_properties", "Manage Properties of Product")
	GroupProductVariant = security.Registry.NewGroup("product_group_product_variant", "Manage Product Variants")
}
