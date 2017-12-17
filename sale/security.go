// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya-addons/account"
	"github.com/hexya-erp/hexya-addons/saleTeams"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.SaleOrder().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.SaleOrderLine().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.SaleOrderLine().Methods().AllowAllToGroup(account.GroupAccountUser)
	pool.AccountInvoiceTax().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.AccountInvoice().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.AccountInvoice().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.AccountInvoiceLine().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.AccountPaymentTerm().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountAnalyticTag().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountAnalyticAccount().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.SaleOrder().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.SaleOrder().Methods().AllowAllToGroup(account.GroupAccountUser)
	pool.SaleReport().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.SaleReport().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	//pool.IrProperty().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.AccountJournal().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.Partner().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.Partner().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductTemplate().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.ProductProduct().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountTax().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.Attachment().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.Attachment().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductUom().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.ProductPricelist().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountAccount().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.ProductUomCateg().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductUom().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductSupplierinfo().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.ProductSupplierinfo().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductPricelist().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.Partner().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.AccountMoveLine().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.SaleOrder().Methods().AllowAllToGroup(account.GroupAccountInvoice)
	pool.SaleOrderLine().Methods().AllowAllToGroup(account.GroupAccountInvoice)
	pool.SaleLayoutCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.SaleLayoutCategory().Methods().AllowAllToGroup(account.GroupAccountManager)
	pool.SaleLayoutCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	pool.SaleLayoutCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesmanAllLeads)
	pool.SaleLayoutCategory().Methods().Load().AllowGroup(account.GroupAccountInvoice)
	pool.ProductPricelistItem().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductPriceHistory().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductTemplate().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductProduct().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductAttribute().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductAttributeValue().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductAttributePrice().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.ProductAttributeLine().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	pool.AccountTax().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountJournal().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountInvoiceTax().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountTaxGroup().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	pool.AccountAccount().Methods().Load().AllowGroup(saleTeams.GroupSaleManager)

}
