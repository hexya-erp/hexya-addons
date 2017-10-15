package product

import (
	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/pool"
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
