// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya-addons/account"
	"github.com/hexya-erp/hexya-addons/saleTeams"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/pool/h"
)

var (
	// GroupSaleLayout enables layouts in sale reports
	GroupSaleLayout *security.Group
	// GroupDeliveryInvoiceAddress enables different delivery and invoice addresses
	GroupDeliveryInvoiceAddress *security.Group
	// GroupShowPriceSubtotal shows line subtotals without taxes (B2B)
	GroupShowPriceSubtotal *security.Group
	// GroupShowPriceTotal shows line subtotals with taxes (B2C)
	GroupShowPriceTotal *security.Group
)

func init() {

	h.SaleOrder().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.SaleOrderLine().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.SaleOrderLine().Methods().AllowAllToGroup(account.GroupAccountUser)
	h.AccountInvoiceTax().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.AccountInvoice().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.AccountInvoice().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.AccountInvoiceLine().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.AccountPaymentTerm().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountAnalyticTag().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountAnalyticAccount().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.SaleOrder().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.SaleOrder().Methods().AllowAllToGroup(account.GroupAccountUser)
	h.SaleReport().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.SaleReport().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	//h.IrProperty().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.AccountJournal().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.Partner().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.Partner().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductTemplate().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.ProductProduct().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountTax().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.Attachment().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.Attachment().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductUom().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.ProductPricelist().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountAccount().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.ProductUomCateg().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductUom().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductSupplierinfo().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.ProductSupplierinfo().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductPricelist().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.Partner().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.AccountMoveLine().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.SaleOrder().Methods().AllowAllToGroup(account.GroupAccountInvoice)
	h.SaleOrderLine().Methods().AllowAllToGroup(account.GroupAccountInvoice)
	h.SaleLayoutCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.SaleLayoutCategory().Methods().AllowAllToGroup(account.GroupAccountManager)
	h.SaleLayoutCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesman)
	h.SaleLayoutCategory().Methods().AllowAllToGroup(saleTeams.GroupSaleSalesmanAllLeads)
	h.SaleLayoutCategory().Methods().Load().AllowGroup(account.GroupAccountInvoice)
	h.ProductPricelistItem().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductPriceHistory().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductTemplate().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductProduct().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductAttribute().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductAttributeValue().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductAttributePrice().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.ProductAttributeLine().Methods().AllowAllToGroup(saleTeams.GroupSaleManager)
	h.AccountTax().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountJournal().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountInvoiceTax().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountTaxGroup().Methods().Load().AllowGroup(saleTeams.GroupSaleSalesman)
	h.AccountAccount().Methods().Load().AllowGroup(saleTeams.GroupSaleManager)

}
