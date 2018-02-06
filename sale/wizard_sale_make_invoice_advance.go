// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"fmt"
	"strconv"

	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.SaleAdvancePaymentInv().DeclareTransientModel()

	h.SaleAdvancePaymentInv().AddFields(map[string]models.FieldDefinition{
		"AdvancePaymentMethod": models.SelectionField{String: "What do you want to invoice?",
			Selection: types.Selection{
				"delivered":  "Invoiceable lines",
				"all":        "Invoiceable lines (deduct down payments)",
				"percentage": "Down payment (percentage)",
				"fixed":      "Down payment (fixed amount)"},
			OnChange: h.SaleAdvancePaymentInv().Methods().OnchangeAdvancePaymentMethod(),
			Default: func(env models.Environment) interface{} {
				if len(env.Context().GetIntegerSlice("active_ids")) == 1 {
					order := h.SaleOrder().Browse(env, env.Context().GetIntegerSlice("active_ids"))
					if order.InvoiceCount() > 0 {
						return "all"
					}
					for _, line := range order.OrderLine().Records() {
						if line.Product().InvoicePolicy() != "order" {
							return "delivered"
						}
					}
					return "all"
				}
				return "delivered"
			}},
		"Product": models.Many2OneField{String: "Down Payment Product", RelationModel: h.ProductProduct(),
			Filter: q.ProductProduct().Type().Equals("service"),
			Default: func(env models.Environment) interface{} {
				return h.SaleAdvancePaymentInv().NewSet(env).DefaultProduct()
			}},
		"Count": models.IntegerField{String: "# of Orders",
			Default: func(env models.Environment) interface{} {
				return len(env.Context().GetIntegerSlice("active_ids"))
			}},
		"Amount": models.FloatField{String: "Down Payment Amount", Digits: decimalPrecision.GetPrecision("Account"),
			Help: "The amount to be invoiced in advance, taxes excluded."},
		"DepositAccount": models.Many2OneField{String: "Income Account", RelationModel: h.AccountAccount(),
			Filter: q.AccountAccount().Deprecated().Equals(false), Help: "Account used for deposits",
			Default: func(env models.Environment) interface{} {
				return h.SaleAdvancePaymentInv().NewSet(env).DefaultProduct().PropertyAccountIncome()
			}},
		"DepositTaxes": models.Many2ManyField{String: "Customer Taxes", RelationModel: h.AccountTax(),
			JSON: "deposit_taxes_id", Help: "Taxes used for deposits",
			Default: func(env models.Environment) interface{} {
				return h.SaleAdvancePaymentInv().NewSet(env).DefaultProduct().Taxes()
			}},
	})

	h.SaleAdvancePaymentInv().Methods().DefaultProduct().DeclareMethod(
		`DefaultProduct returns the default deposit product`,
		func(rs h.SaleAdvancePaymentInvSet) h.ProductProductSet {
			conf := h.ConfigParameter().NewSet(rs.Env()).GetParam("deposit_product_id_setting", "")
			accountID, err := strconv.ParseInt(conf, 10, 64)
			if err != nil {
				return h.ProductProduct().NewSet(rs.Env())
			}
			return h.ProductProduct().Browse(rs.Env(), []int64{accountID})
		})

	h.SaleAdvancePaymentInv().Methods().OnchangeAdvancePaymentMethod().DeclareMethod(
		`OnchangeAdvancePaymentMethod sets the amount to 0 when percentage is selected.`,
		func(rs h.SaleAdvancePaymentInvSet) (*h.SaleAdvancePaymentInvData, []models.FieldNamer) {
			var fieldsToReset []models.FieldNamer
			if rs.AdvancePaymentMethod() == "percentage" {
				fieldsToReset = append(fieldsToReset, h.SaleAdvancePaymentInv().Amount())
			}
			return &h.SaleAdvancePaymentInvData{}, fieldsToReset
		})

	h.SaleAdvancePaymentInv().Methods().CreateInvoice().DeclareMethod(
		`CreateInvoice creates a deposit invoice for the given order and order line.`,
		func(rs h.SaleAdvancePaymentInvSet, order h.SaleOrderSet, soLine h.SaleOrderLineSet) h.AccountInvoiceSet {
			account := h.AccountAccount().NewSet(rs.Env())
			if !rs.Product().IsEmpty() {
				account = rs.Product().PropertyAccountIncome()
			}
			if account.IsEmpty() {
				//inc_acc = ir_property_obj.get('property_account_income_categ_id', 'product.category')
				//account_id = order.fiscal_position_id.map_account(inc_acc).id if inc_acc else False
			}
			if account.IsEmpty() {
				panic(rs.T("There is no income account defined for this product: '%s'."+
					" You may have to install a chart of account from Accounting app, settings menu.",
					rs.Product().Name()))
			}
			if rs.Amount() <= 0 {
				panic(rs.T("The value of the down payment amount must be positive."))
			}
			var (
				amount float64
				name   string
			)
			if rs.AdvancePaymentMethod() == "percentage" {
				amount = order.AmountUntaxed() * rs.Amount() / 100
				name = rs.T("Down payment of %s%%", rs.Amount())
			} else {
				amount = rs.Amount()
				name = rs.T("Down Payment")
			}
			taxes := rs.Product().Taxes()
			if !order.Company().IsEmpty() {
				taxes = taxes.Search(q.AccountTax().Company().Equals(order.Company()))
			}
			if !order.FiscalPosition().IsEmpty() && !taxes.IsEmpty() {
				taxes = order.FiscalPosition().MapTax(taxes, h.ProductProduct().NewSet(rs.Env()),
					h.Partner().NewSet(rs.Env()))
			}
			nameInv := order.Name()
			if order.ClientOrderRef() != "" {
				nameInv = order.ClientOrderRef()
			}
			fPos := order.Partner().PropertyAccountPosition()
			if !order.FiscalPosition().IsEmpty() {
				fPos = order.FiscalPosition()
			}
			invoiceLines := h.AccountInvoiceLine().Create(rs.Env(),
				&h.AccountInvoiceLineData{
					Name:             name,
					Origin:           order.Name(),
					Account:          account,
					PriceUnit:        amount,
					Quantity:         1,
					Discount:         0,
					Uom:              rs.Product().Uom(),
					Product:          rs.Product(),
					SaleLines:        soLine,
					InvoiceLineTaxes: taxes,
					AccountAnalytic:  order.Project(),
				})
			invoice := h.AccountInvoice().Create(rs.Env(),
				&h.AccountInvoiceData{
					Name:            nameInv,
					Origin:          order.Name(),
					Type:            "out_invoice",
					Reference:       "",
					Account:         order.Partner().PropertyAccountReceivable(),
					Partner:         order.PartnerInvoice(),
					PartnerShipping: order.PartnerShipping(),
					InvoiceLines:    invoiceLines,
					Currency:        order.Pricelist().Currency(),
					PaymentTerm:     order.PaymentTerm(),
					FiscalPosition:  fPos,
					Team:            order.Team(),
					User:            order.User(),
					Comment:         order.Note(),
				})
			invoice.ComputeTaxes()
			//invoice.message_post_with_view('mail.message_origin_link',
			//            values={'self': invoice, 'origin': order},
			//            subtype_id=self.env.ref('mail.mt_note').id)
			return invoice
		})

	h.SaleAdvancePaymentInv().Methods().CreateInvoices().DeclareMethod(
		`CreateInvoices is the main method called from the wizard to create the invoices.`,
		func(rs h.SaleAdvancePaymentInvSet) *actions.Action {
			rs.EnsureOne()
			saleOrders := h.SaleOrder().Browse(rs.Env(), rs.Env().Context().GetIntegerSlice("active_ids"))
			switch rs.AdvancePaymentMethod() {
			case "delivered":
				saleOrders.ActionInvoiceCreate(false, false)
			case "all":
				saleOrders.ActionInvoiceCreate(false, true)
			default:
				// Create deposit product if necessary
				if rs.Product().IsEmpty() {
					depositProduct := h.ProductProduct().Create(rs.Env(), rs.PrepareDepositProduct())
					rs.SetProduct(depositProduct)
					h.ConfigParameter().NewSet(rs.Env()).SetParam("deposit_product_id_setting",
						fmt.Sprintf("%d", depositProduct.ID()))
				}

				for _, order := range saleOrders.Records() {
					amount := rs.Amount()
					if rs.AdvancePaymentMethod() == "percentage" {
						amount = order.AmountUntaxed() * rs.Amount() / 100
					}
					if rs.Product().InvoicePolicy() != "order" {
						panic(rs.T(`The product used to invoice a down payment should have an invoice policy set
to 'Ordered quantities'. Please update your deposit product to be able to create a deposit invoice.`))
					}
					if rs.Product().Type() != "service" {
						panic(rs.T(`The product used to invoice a down payment should be of type 'Service'.
Please use another product or update this product.`))
					}
					taxes := rs.Product().Taxes()
					if !order.Company().IsEmpty() {
						taxes = taxes.Search(q.AccountTax().Company().Equals(order.Company()))
					}
					if !order.FiscalPosition().IsEmpty() && !taxes.IsEmpty() {
						taxes = order.FiscalPosition().MapTax(taxes, h.ProductProduct().NewSet(rs.Env()),
							h.Partner().NewSet(rs.Env()))
					}
					SOLine := h.SaleOrderLine().Create(rs.Env(),
						&h.SaleOrderLineData{
							Name:          rs.T("Advance: %v", dates.Today()),
							PriceUnit:     amount,
							ProductUomQty: 0,
							Order:         order,
							Discount:      0,
							ProductUom:    rs.Product().Uom(),
							Product:       rs.Product(),
							Tax:           taxes,
						})
					rs.CreateInvoice(order, SOLine)
				}
			}
			if rs.Env().Context().GetBool("open_invoice") {
				return saleOrders.ActionViewInvoice()
			}
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	h.SaleAdvancePaymentInv().Methods().PrepareDepositProduct().DeclareMethod(
		`PrepareDepositProduct returns the data used to create the deposit product.`,
		func(rs h.SaleAdvancePaymentInvSet) *h.ProductProductData {
			return &h.ProductProductData{
				Name:                  "Down payment",
				Type:                  "service",
				InvoicePolicy:         "order",
				PropertyAccountIncome: rs.DepositAccount(),
				Taxes: rs.DepositTaxes(),
			}
		})

}
