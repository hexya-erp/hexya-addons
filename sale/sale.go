// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/hexya-erp/hexya-addons/account/accounttypes"
	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/hexya/views"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.SaleOrder().DeclareModel()
	h.SaleOrder().SetDefaultOrder("DateOrder DESC", "ID DESC")

	h.SaleOrder().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Order Reference", Required: true, NoCopy: true, /*[ readonly True]*/
			/*[ states {'draft': [('readonly']*/ /*[ False)]}]*/ Index: true,
			Default: func(env models.Environment) interface{} {
				return h.SaleOrder().NewSet(env).T("New")
			}},
		"Origin": models.CharField{String: "Source Document",
			Help: "Reference of the document that generated this sales order request."},
		"ClientOrderRef": models.CharField{String: "Customer Reference", NoCopy: true},
		"State": models.SelectionField{String: "Status", Selection: types.Selection{
			"draft":  "Quotation",
			"sent":   "Quotation Sent",
			"sale":   "Sales Order",
			"done":   "Locked",
			"cancel": "Cancelled",
		}, ReadOnly: true, NoCopy: true, Index: true, /*[ track_visibility 'onchange']*/
			Default: models.DefaultValue("draft")},
		"DateOrder": models.DateTimeField{String: "Order Date", Required: true, Index: true, /*[ readonly True]*/
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			NoCopy: true, Default: func(env models.Environment) interface{} {
				return dates.Now()
			}},
		"ValidityDate": models.DateField{String: "Expiration Date" /*[ readonly True]*/, NoCopy: true,
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: `Manually set the expiration date of your quotation (offer), or it will set the date automatically
based on the template if online quotation is installed.`},
		"ConfirmationDate": models.DateTimeField{Index: true, ReadOnly: true,
			Help: "Date on which the sale order is confirmed."},
		"User": models.Many2OneField{String: "Salesperson", RelationModel: h.User(), Index: true, /*[ track_visibility 'onchange']*/
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser()
			}},
		"Partner": models.Many2OneField{String: "Customer", RelationModel: h.Partner(), /* readonly=true */
			/*[ states {'draft': [('readonly'] [ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			OnChange: h.SaleOrder().Methods().OnchangePartner(),
			Required: true, Index: true /*[ track_visibility 'always']*/},
		"PartnerInvoice": models.Many2OneField{String: "Invoice Address", RelationModel: h.Partner(), /* readonly=true */
			Required: true,
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: "Invoice address for current sales order."},
		"PartnerShipping": models.Many2OneField{String: "Delivery Address", RelationModel: h.Partner(), /* readonly=true */
			Required: true,
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			OnChange: h.SaleOrder().Methods().OnchangePartnerShipping(),
			Help:     "Delivery address for current sales order."},
		"Pricelist": models.Many2OneField{RelationModel: h.ProductPricelist(), Required: true, /* readonly=true */
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: "Pricelist for current sales order."},
		"Currency": models.Many2OneField{RelationModel: h.Currency(),
			Related: "Pricelist.Currency", ReadOnly: true, Required: true},
		"Project": models.Many2OneField{String: "Analytic Account", RelationModel: h.AccountAnalyticAccount(),
			/* readonly=true */
			/*[ states {'draft': [('readonly'] [ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: "The analytic account related to a sales order.", NoCopy: true},
		"RelatedProject": models.Many2OneField{String: "Analytic Account", RelationModel: h.AccountAnalyticAccount(),
			Related: "Project", Help: "The analytic account related to a sales order."},
		"OrderLine": models.One2ManyField{String: "Order Lines", RelationModel: h.SaleOrderLine(),
			ReverseFK: "Order",
			/*[ states {'cancel': [('readonly'] [ True)]]*/
			/*[ 'done': [('readonly'] [ True)]}]*/
			NoCopy: false},
		"InvoiceCount": models.IntegerField{String: "# of Invoices",
			Compute: h.SaleOrder().Methods().GetInvoiced(),
			Depends: []string{"state", "OrderLine.InvoiceStatus"}, GoType: new(int)},
		"Invoices": models.Many2ManyField{String: "Invoices", RelationModel: h.AccountInvoice(),
			JSON: "invoice_ids", Compute: h.SaleOrder().Methods().GetInvoiced(),
			Depends: []string{"state", "OrderLine.InvoiceStatus"}, NoCopy: true},
		"InvoiceStatus": models.SelectionField{Selection: types.Selection{
			"upselling":  "Upselling Opportunity",
			"invoiced":   "Fully Invoiced",
			"to invoice": "To Invoice",
			"no":         "Nothing to Invoice",
		}, Compute: h.SaleOrder().Methods().GetInvoiced(),
			Depends: []string{"state", "OrderLine.InvoiceStatus"}, Stored: true},
		"Note": models.TextField{String: "Terms and conditions",
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company().SaleNote()
			}},
		"AmountUntaxed": models.FloatField{String: "Untaxed Amount", Stored: true,
			Compute: h.SaleOrder().Methods().AmountAll(), /*[ track_visibility 'always']*/
			Depends: []string{"OrderLine.PriceTotal"}},
		"AmountTax": models.FloatField{String: "Taxes", Stored: true,
			Compute: h.SaleOrder().Methods().AmountAll(), /*[ track_visibility 'always']*/
			Depends: []string{"OrderLine.PriceTotal"}},
		"AmountTotal": models.FloatField{String: "Total", Stored: true,
			Compute: h.SaleOrder().Methods().AmountAll(), /*[ track_visibility 'always']*/
			Depends: []string{"OrderLine.PriceTotal"}},
		"PaymentTerm": models.Many2OneField{String: "Payment Terms", RelationModel: h.AccountPaymentTerm()},
		"FiscalPosition": models.Many2OneField{RelationModel: h.AccountFiscalPosition(),
			OnChange: h.SaleOrder().Methods().ComputeTax()},
		"Company": models.Many2OneField{RelationModel: h.Company(),
			Default: func(env models.Environment) interface{} {
				return h.Company().NewSet(env).CompanyDefaultGet()
			}},
		"Team": models.Many2OneField{String: "Sales Team", RelationModel: h.CRMTeam(),
			Default: func(env models.Environment) interface{} {
				return h.CRMTeam().NewSet(env).GetDefaultTeam(h.User().NewSet(env))
			}},
		"ProcurementGroup": models.Many2OneField{RelationModel: h.ProcurementGroup(), NoCopy: true},
		"Product":          models.Many2OneField{RelationModel: h.ProductProduct(), Related: "OrderLine.Product"},
	})

	h.SaleOrder().Methods().AmountAll().DeclareMethod(
		`AmountAll computes all the amounts of this sale order by summing its sale order lines.`,
		func(rs h.SaleOrderSet) *h.SaleOrderData {
			var amountUntaxed, amountTaxed float64
			for _, line := range rs.OrderLine().Records() {
				amountUntaxed += line.PriceSubtotal()
				if rs.Company().TaxCalculationRoundingMethod() == "round_globally" {
					price := line.PriceUnit() * (1 - line.Discount()/100)
					_, _, _, taxes := line.Tax().ComputeAll(price, line.Order().Currency(), line.ProductUomQty(),
						line.Product(), rs.PartnerShipping())
					for _, t := range taxes {
						amountTaxed += t.Amount
					}
				} else {
					amountTaxed += line.PriceTax()
				}
			}
			return &h.SaleOrderData{
				AmountUntaxed: rs.Pricelist().Currency().Round(amountUntaxed),
				AmountTax:     rs.Pricelist().Currency().Round(amountTaxed),
				AmountTotal:   amountTaxed + amountUntaxed,
			}
		})

	h.SaleOrder().Methods().GetInvoiced().DeclareMethod(
		`GetInvoiced computes the invoice status of a SO. Possible statuses:

			  - no: if the SO is not in status 'sale' or 'done', we consider that there is nothing to
			    invoice. This is also hte default value if the conditions of no other status is met.
			  - to invoice: if any SO line is 'to invoice', the whole SO is 'to invoice'
			  - invoiced: if all SO lines are invoiced, the SO is invoiced.
			  - upselling: if all SO lines are invoiced or upselling, the status is upselling.

			  The invoice_ids are obtained thanks to the invoice lines of the SO lines, and we also search
			  for possible refunds created directly from existing invoices. This is necessary since such a
			  refund is not directly linked to the SO.`,
		func(rs h.SaleOrderSet) *h.SaleOrderData {
			invoices := h.AccountInvoice().NewSet(rs.Env())
			for _, line := range rs.OrderLine().Records() {
				for _, invLine := range line.InvoiceLines().Records() {
					if invLine.Invoice().Type() == "out_invoice" || invLine.Invoice().Type() == "out_refund" {
						invoices = invoices.Union(invLine.Invoice())
					}
				}
			}
			// Search for invoices which have been 'cancelled' (filter_refund = 'modify' in
			// 'account.invoice.refund')
			// use like as origin may contains multiple references (e.g. 'SO01, SO02')
			origins := strings.Split(rs.Origin(), ",")
			for i, o := range origins {
				origins[i] = strings.TrimSpace(o)
			}
			invoices = invoices.Union(h.AccountInvoice().Search(rs.Env(),
				q.AccountInvoice().Origin().Like(rs.Name()).
					And().Name().In(origins).
					And().Type().In([]string{"out_invoice", "out_refund"})))

			refunds := h.AccountInvoice().NewSet(rs.Env())
			for _, inv := range invoices.Records() {
				refunds = refunds.Union(h.AccountInvoice().Search(rs.Env(),
					q.AccountInvoice().Type().Equals("out_refund").
						And().Origin().Equals(inv.Number()).
						And().Origin().IsNotNull().
						And().Journal().Equals(inv.Journal())))
			}
			// Ignore the status of the deposit product
			depositProduct := h.SaleAdvancePaymentInv().NewSet(rs.Env()).DefaultProduct()
			lineInvoiceStatus := make(map[string]bool)
			for _, l := range rs.OrderLine().Records() {
				if l.Product().Equals(depositProduct) {
					continue
				}
				lineInvoiceStatus[l.InvoiceStatus()] = true
			}
			var invoiceStatus string
			switch {
			case rs.State() != "sale" && rs.State() != "done":
				invoiceStatus = "no"
			case lineInvoiceStatus["to invoice"]:
				invoiceStatus = "to invoice"
			case len(lineInvoiceStatus) == 1 && lineInvoiceStatus["invoiced"]:
				invoiceStatus = "invoiced"
			case len(lineInvoiceStatus) <= 2 && (lineInvoiceStatus["invoiced"] || lineInvoiceStatus["upselling"]):
				invoiceStatus = "upselling"
			default:
				invoiceStatus = "no"
			}
			return &h.SaleOrderData{
				InvoiceCount:  invoices.Union(refunds).Len(),
				Invoices:      invoices.Union(refunds),
				InvoiceStatus: invoiceStatus,
			}
		})

	h.SaleOrder().Methods().ComputeTax().DeclareMethod(
		`ComputeTax triggers the recompute of the taxes if the fiscal position is changed on the SO.`,
		func(rs h.SaleOrderSet) (*h.SaleOrderData, []models.FieldNamer) {
			//@api.onchange('fiscal_position_id')
			/*def _compute_tax_id(self):
			  """
			  Trigger the recompute of the taxes if the fiscal position is changed on the SO.
			  """
			  for order in self:
			      order.order_line._compute_tax_id()
			*/
			// TODO : need to implement onchange on relation fields first
			return &h.SaleOrderData{}, []models.FieldNamer{}
		})

	h.SaleOrder().Methods().GetCustomerLead().DeclareMethod(
		`GetCustomerLead returns the delay to deliver the given product template`,
		func(rs h.SaleOrderSet, productTmpl h.ProductTemplateSet) int {
			return 0
		})

	h.SaleOrder().Methods().ButtonDummy().DeclareMethod(
		`ButtonDummy is a dummy function to force reload of the form on client side.`,
		func(rs h.SaleOrderSet) bool {
			return true
		})

	h.SaleOrder().Methods().Unlink().Extend("",
		func(rs h.SaleOrderSet) int64 {
			for _, order := range rs.Records() {
				if order.State() != "draft" && order.State() != "cancel" {
					panic(rs.T("You can not delete a sent quotation or a sales order! Try to cancel it before."))
				}
			}
			return rs.Super().Unlink()
		})

	//h.SaleOrder().Methods().TrackSubtype().DeclareMethod(
	//	`TrackSubtype`,
	//	func(rs h.SaleOrderSet, initvalues interface{}) {
	//		//@api.multi
	//		/*def _track_subtype(self, init_values):
	//		  self.ensure_one()
	//		  if 'state' in init_values and self.state == 'sale':
	//		      return 'sale.mt_order_confirmed'
	//		  elif 'state' in init_values and self.state == 'sent':
	//		      return 'sale.mt_order_sent'
	//		  return super(SaleOrder, self)._track_subtype(init_values)
	//
	//		*/
	//	})

	h.SaleOrder().Methods().OnchangePartnerShipping().DeclareMethod(
		`OnchangePartnerShipping triggers the change of fiscal position when the shipping address is modified.`,
		func(rs h.SaleOrderSet) (*h.SaleOrderData, []models.FieldNamer) {
			return &h.SaleOrderData{
				FiscalPosition: h.AccountFiscalPosition().NewSet(rs.Env()).GetFiscalPosition(rs.Partner(), rs.PartnerShipping()),
			}, []models.FieldNamer{h.SaleOrder().FiscalPosition()}
		})

	h.SaleOrder().Methods().OnchangePartner().DeclareMethod(
		`OnchangePartner updates the following fields when the partner is changed:
		- Pricelist
		- Payment term
		- Invoice address
		- Delivery address
		`,
		func(rs h.SaleOrderSet) (*h.SaleOrderData, []models.FieldNamer) {
			if rs.Partner().IsEmpty() {
				return &h.SaleOrderData{
						PartnerInvoice:  h.Partner().NewSet(rs.Env()),
						PartnerShipping: h.Partner().NewSet(rs.Env()),
						PaymentTerm:     h.AccountPaymentTerm().NewSet(rs.Env()),
						FiscalPosition:  h.AccountFiscalPosition().NewSet(rs.Env()),
					}, []models.FieldNamer{
						h.SaleOrder().PartnerInvoice(),
						h.SaleOrder().PartnerShipping(),
						h.SaleOrder().PaymentTerm(),
						h.SaleOrder().FiscalPosition(),
					}
			}
			addr := rs.Partner().AddressGet([]string{"delivery", "invoice"})
			values := &h.SaleOrderData{
				Pricelist:       rs.Partner().PropertyProductPricelist(),
				PaymentTerm:     rs.Partner().PropertyPaymentTerm(),
				PartnerInvoice:  addr["invoice"],
				PartnerShipping: addr["delivery"],
			}
			fields := []models.FieldNamer{
				h.SaleOrder().PartnerInvoice(),
				h.SaleOrder().PartnerShipping(),
				h.SaleOrder().PaymentTerm(),
				h.SaleOrder().Pricelist(),
			}
			if h.User().NewSet(rs.Env()).CurrentUser().Company().SaleNote() != "" {
				values.Note = h.User().NewSet(rs.Env()).WithContext("lang", rs.Partner().Lang()).
					CurrentUser().Company().SaleNote()
				fields = append(fields, h.SaleOrder().Note())
			}
			if !rs.Partner().User().IsEmpty() {
				values.User = rs.Partner().User()
				fields = append(fields, h.SaleOrder().User())
			}
			if !rs.Partner().Team().IsEmpty() {
				values.Team = rs.Partner().Team()
				fields = append(fields, h.SaleOrder().Team())
			}
			return values, fields
		})

	//h.SaleOrder().Methods().OnchangePartnerWarning().DeclareMethod(
	//	`OnchangePartnerWarning`,
	//	func(rs h.SaleOrderSet) {
	//		//@api.onchange('partner_id')
	//		/*def onchange_partner_id_warning(self):
	//		  if not self.partner_id:
	//		      return
	//		  warning = {}
	//		  title = False
	//		  message = False
	//		  partner = self.partner_id
	//
	//		  # If partner has no warning, check its company
	//		  if partner.sale_warn == 'no-message' and partner.parent_id:
	//		      partner = partner.parent_id
	//
	//		  if partner.sale_warn != 'no-message':
	//		      # Block if partner only has warning but parent company is blocked
	//		      if partner.sale_warn != 'block' and partner.parent_id and partner.parent_id.sale_warn == 'block':
	//		          partner = partner.parent_id
	//		      title = ("Warning for %s") % partner.name
	//		      message = partner.sale_warn_msg
	//		      warning = {
	//		              'title': title,
	//		              'message': message,
	//		      }
	//		      if partner.sale_warn == 'block':
	//		          self.update({'partner_id': False, 'partner_invoice_id': False, 'partner_shipping_id': False, 'pricelist_id': False})
	//		          return {'warning': warning}
	//
	//		  if warning:
	//		      return {'warning': warning}
	//
	//		*/
	// TODO : need to implement onchange warnings first
	//	})

	h.SaleOrder().Methods().Create().Extend("",
		func(rs h.SaleOrderSet, data *h.SaleOrderData) h.SaleOrderSet {
			if data.Name == "" || data.Name == rs.T("New") {
				seq := h.Sequence().NewSet(rs.Env())
				if !data.Company.IsEmpty() {
					seq = seq.WithContext("force_company", data.Company.ID())
				}
				data.Name = rs.T("New")
				seqValue := seq.NextByCode("sale.order")
				if seqValue != "" {
					data.Name = seqValue
				}
			}
			// Makes sure PartnerInvoice, PartnerShipping and Pricelist are defined
			addr := data.Partner.AddressGet([]string{"delivery", "invoice"})
			if data.PartnerInvoice.IsEmpty() {
				data.PartnerInvoice = addr["invoice"]
			}
			if data.PartnerShipping.IsEmpty() {
				data.PartnerShipping = addr["delivery"]
			}
			if data.Pricelist.IsEmpty() {
				data.Pricelist = data.Partner.PropertyProductPricelist()
			}
			return rs.Super().Create(data)
		})

	h.SaleOrder().Methods().PrepareInvoice().DeclareMethod(
		`PrepareInvoice prepares the data to create the new invoice for a sales order. This method may be
			  overridden to implement custom invoice generation (making sure to call super() to establish
			  a clean extension chain).`,
		func(rs h.SaleOrderSet) *h.AccountInvoiceData {
			rs.EnsureOne()
			journal := h.AccountInvoice().NewSet(rs.Env()).DefaultJournal()
			if journal.IsEmpty() {
				panic(rs.T("Please define an accounting sale journal for this company."))
			}
			fPos := rs.PartnerInvoice().PropertyAccountPosition()
			if !rs.FiscalPosition().IsEmpty() {
				fPos = rs.FiscalPosition()
			}
			invoiceVals := &h.AccountInvoiceData{
				Name:            rs.ClientOrderRef(),
				Origin:          rs.Name(),
				Type:            "out_invoice",
				Account:         rs.PartnerInvoice().PropertyAccountReceivable(),
				Partner:         rs.PartnerInvoice(),
				PartnerShipping: rs.PartnerShipping(),
				Journal:         journal,
				Currency:        rs.Pricelist().Currency(),
				Comment:         rs.Note(),
				PaymentTerm:     rs.PaymentTerm(),
				FiscalPosition:  fPos,
				Company:         rs.Company(),
				User:            rs.User(),
				Team:            rs.Team(),
			}
			return invoiceVals
		})

	h.SaleOrder().Methods().PrintQuotation().DeclareMethod(
		`PrintQuotation returns the action to print the quotation report`,
		func(rs h.SaleOrderSet) *actions.Action {
			//@api.multi
			/*def print_quotation(self):
			  self.filtered(lambda s: s.state == 'draft').write({'state': 'sent'})
			  return self.env['report'].get_action(self, 'sale.report_saleorder')

			*/
			// TODO Implement reports first
			rs.Search(q.SaleOrder().State().Equals("draft")).SetState("sent")
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	h.SaleOrder().Methods().ActionViewInvoice().DeclareMethod(
		`ActionViewInvoice returns an action to view the invoice(s) related to this order.
		If there is a single invoice, then it will be opened in a form view, otherwise in list view.`,
		func(rs h.SaleOrderSet) *actions.Action {
			invoices := h.AccountInvoice().NewSet(rs.Env())
			for _, order := range rs.Records() {
				invoices = invoices.Union(order.Invoices())
			}
			action := actions.Registry.MustGetById("account_account_invoice_tree1")
			switch {
			case invoices.Len() > 1:
				idsStr := make([]string, invoices.Len())
				for i, inv := range invoices.Records() {
					idsStr[i] = strconv.FormatInt(inv.ID(), 10)
				}
				action.Domain = fmt.Sprintf("[('id', 'in', (%s))]", strings.Join(idsStr, ","))
			case invoices.Len() == 1:
				action.Views = []views.ViewTuple{{
					ID:   "account_invoice_form",
					Type: views.ViewTypeForm,
				}}
				action.ResID = invoices.ID()
			default:
				action = &actions.Action{Type: actions.ActionCloseWindow}
			}
			return action
		})

	h.SaleOrder().Methods().ActionInvoiceCreate().DeclareMethod(
		`ActionInvoiceCreate creates the invoice associated to the SO.

		- If grouped is true, invoices are grouped by sale orders.
		If False, invoices are grouped by (partner_invoice_id, currency)
        - If final is true, refunds will be generated if necessary.`,
		func(rs h.SaleOrderSet, grouped, final bool) h.AccountInvoiceSet {
			type keyStruct struct {
				OrderID    int64
				PartnerID  int64
				CurrencyID int64
			}
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			invoices := make(map[keyStruct]h.AccountInvoiceSet)
			references := make(map[int64]h.SaleOrderSet)
			for _, order := range rs.Records() {
				groupKey := keyStruct{PartnerID: order.PartnerInvoice().ID(), CurrencyID: order.Currency().ID()}
				if grouped {
					groupKey = keyStruct{OrderID: order.ID()}
				}
				lines := order.OrderLine().Records()
				sort.Slice(lines, func(i, j int) bool {
					return lines[i].QtyToInvoice() < lines[j].QtyToInvoice()
				})
				for _, line := range lines {
					if nbutils.IsZero(line.QtyToInvoice(), precision) {
						continue
					}
					if _, exists := invoices[groupKey]; !exists {
						invData := order.PrepareInvoice()
						invoice := h.AccountInvoice().Create(rs.Env(), invData)
						references[invoice.ID()] = order
						invoices[groupKey] = invoice
					} else {
						vals := h.AccountInvoiceData{}
						origins := strings.Split(invoices[groupKey].Origin(), ", ")
						var inOrigins bool
						for _, o := range origins {
							if o == order.Name() {
								inOrigins = true
								break
							}
						}
						if !inOrigins {
							vals.Origin = invoices[groupKey].Origin() + ", " + order.Name()
						}
						names := strings.Split(invoices[groupKey].Name(), ", ")
						var inNames bool
						for _, n := range names {
							if n == order.ClientOrderRef() {
								inNames = true
								break
							}
						}
						if !inNames && order.ClientOrderRef() != "" {
							vals.Name = invoices[groupKey].Name() + ", " + order.ClientOrderRef()
						}
					}
					if line.QtyToInvoice() > 0 || (line.QtyToInvoice() < 0 && final) {
						line.InvoiceLineCreate(invoices[groupKey], line.QtyToInvoice())
					}
				}
				if ref, exists := references[invoices[groupKey].ID()]; exists {
					if order.Intersect(ref).IsEmpty() {
						references[invoices[groupKey].ID()] = ref.Union(order)
					}
				}
			}
			if len(invoices) == 0 {
				panic(rs.T("There is no invoicable line"))
			}

			res := h.AccountInvoice().NewSet(rs.Env())
			for _, invoice := range invoices {
				if invoice.InvoiceLines().IsEmpty() {
					panic(rs.T("There is no invoicable line."))
				}
				// If invoice is negative, do a refund invoice instead
				if invoice.AmountUntaxed() < 0 {
					invoice.SetType("out_refund")
					for _, line := range invoice.InvoiceLines().Records() {
						line.SetQuantity(-line.Quantity())
					}
				}
				// Use additional field helper function (for account extensions)
				for _, line := range invoice.InvoiceLines().Records() {
					line.DefineAdditionalFields(invoice)
				}
				// Necessary to force computation of taxes. In account_invoice, they are triggered
				// by onchanges, which are not triggered when doing a create.
				invoice.ComputeTaxes()
				//invoice.message_post_with_view('mail.message_origin_link',
				//    values={'self': invoice, 'origin': references[invoice]},
				//    subtype_id=self.env.ref('mail.mt_note').id)
				res = res.Union(invoice)
			}
			return res
		})

	h.SaleOrder().Methods().ActionDraft().DeclareMethod(
		`ActionDraft sets this sale order back to the draft state.`,
		func(rs h.SaleOrderSet) bool {
			orders := h.SaleOrder().NewSet(rs.Env())
			for _, order := range rs.Records() {
				if order.State() != "cancel" && order.State() != "sent" {
					continue
				}
				orders = orders.Union(order)
			}
			orders.Write(&h.SaleOrderData{
				State:            "draft",
				ProcurementGroup: h.ProcurementGroup().NewSet(rs.Env()),
			})
			for _, order := range orders.Records() {
				for _, line := range order.OrderLine().Records() {
					for _, proc := range line.Procurements().Records() {
						proc.SetSaleLine(h.SaleOrderLine().NewSet(rs.Env()))
					}
				}
			}
			return true
		})

	h.SaleOrder().Methods().ActionCancel().DeclareMethod(
		`ActionCancel cancels this sale order.`,
		func(rs h.SaleOrderSet) bool {
			rs.SetState("cancel")
			return true
		})

	h.SaleOrder().Methods().ActionQuotationSend().DeclareMethod(
		`ActionQuotationSend opens a window to compose an email,
		with the edi sale template message loaded by default`,
		func(rs h.SaleOrderSet) *actions.Action {
			rs.EnsureOne()
			//		//@api.multi
			//		/*def action_quotation_send(self):
			//		  '''
			//		  This function opens a window to compose an email, with the edi sale template message loaded by default
			//		  '''
			//		  self.ensure_one()
			//		  ir_model_data = self.env['ir.model.data']
			//		  try:
			//		      template_id = ir_model_data.get_object_reference('sale', 'email_template_edi_sale')[1]
			//		  except ValueError:
			//		      template_id = False
			//		  try:
			//		      compose_form_id = ir_model_data.get_object_reference('mail', 'email_compose_message_wizard_form')[1]
			//		  except ValueError:
			//		      compose_form_id = False
			//		  ctx = dict()
			//		  ctx.update({
			//		      'default_model': 'sale.order',
			//		      'default_res_id': self.ids[0],
			//		      'default_use_template': bool(template_id),
			//		      'default_template_id': template_id,
			//		      'default_composition_mode': 'comment',
			//		      'mark_so_as_sent': True,
			//		      'custom_layout': "sale.mail_template_data_notification_email_sale_order"
			//		  })
			//		  return {
			//		      'type': 'ir.actions.act_window',
			//		      'view_type': 'form',
			//		      'view_mode': 'form',
			//		      'res_model': 'mail.compose.message',
			//		      'views': [(compose_form_id, 'form')],
			//		      'view_id': compose_form_id,
			//		      'target': 'new',
			//		      'context': ctx,
			//		  }
			//
			//		*/
			// FIXME: Next line for demo only
			rs.Search(q.SaleOrder().State().Equals("draft")).SetState("sent")
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	h.SaleOrder().Methods().ForceQuotationSend().DeclareMethod(
		`ForceQuotationSend`,
		func(rs h.SaleOrderSet) bool {
			//@api.multi
			/*def force_quotation_send(self):
			  for order in self:
			      email_act = order.action_quotation_send()
			      if email_act and email_act.get('context'):
			          email_ctx = email_act['context']
			          email_ctx.update(default_email_from=order.company_id.email)
			          order.with_context(email_ctx).message_post_with_template(email_ctx.get('default_template_id'))
			  return True

			*/
			return true
		})

	h.SaleOrder().Methods().ActionDone().DeclareMethod(
		`ActionDone sets the state of this sale order to done`,
		func(rs h.SaleOrderSet) bool {
			rs.SetState("done")
			return true
		})

	h.SaleOrder().Methods().PrepareProcurementGroup().DeclareMethod(
		`PrepareProcurementGroup returns the data that will be used to create the
		procurement group of this sale order`,
		func(rs h.SaleOrderSet) *h.ProcurementGroupData {
			return &h.ProcurementGroupData{
				Name: rs.Name(),
			}
		})

	h.SaleOrder().Methods().ActionConfirm().DeclareMethod(
		`ActionConfirm confirms this quotation into a sale order`,
		func(rs h.SaleOrderSet) bool {
			for _, order := range rs.Records() {
				order.Write(&h.SaleOrderData{
					State:            "sale",
					ConfirmationDate: dates.Now(),
				})
				if rs.Env().Context().HasKey("send_email") {
					rs.ForceQuotationSend()
				}
				order.OrderLine().ActionProcurementCreate()
			}
			autoDone := h.ConfigParameter().Search(rs.Env(), q.ConfigParameter().Key().Equals("sale.auto_done_setting"))
			if autoDone.Value() != "" {
				rs.ActionDone()
			}
			return true
		})

	h.SaleOrder().Methods().CreateAnalyticAccount().DeclareMethod(
		`CreateAnalyticAccount creates the analytic account (project) for this sale order.`,
		func(rs h.SaleOrderSet, prefix string) {
			for _, order := range rs.Records() {
				name := order.Name()
				if prefix != "" {
					name = fmt.Sprintf("%s: %s", prefix, order.Name())
				}
				analyticAccount := h.AccountAnalyticAccount().Create(rs.Env(), &h.AccountAnalyticAccountData{
					Name:    name,
					Code:    order.ClientOrderRef(),
					Company: order.Company(),
					Partner: order.Partner(),
				})
				order.SetProject(analyticAccount)
			}
		})

	h.SaleOrder().Methods().OrderLinesLayouted().DeclareMethod(
		`OrderLinesLayouted returns this order lines classified by sale_layout_category and separated in
        pages according to the category pagebreaks. Used to render the report.`,
		func(rs h.SaleOrderSet) {
			/*
			   @api.multi
			   def order_lines_layouted(self):
			       self.ensure_one()
			       report_pages = [[]]
			       for category, lines in groupby(self.order_line, lambda l: l.layout_category_id):
			           # If last added category induced a pagebreak, this one will be on a new page
			           if report_pages[-1] and report_pages[-1][-1]['pagebreak']:
			               report_pages.append([])
			           # Append category to current report page
			           report_pages[-1].append({
			               'name': category and category.name or 'Uncategorized',
			               'subtotal': category and category.subtotal,
			               'pagebreak': category and category.pagebreak,
			               'lines': list(lines)
			           })

			       return report_pages

			*/
		})

	h.SaleOrder().Methods().GetTaxAmountByGroup().DeclareMethod(
		`GetTaxAmountByGroup`,
		func(rs h.SaleOrderSet) []accounttypes.TaxGroup {
			rs.EnsureOne()
			currency := rs.Company().Currency()
			if !rs.Currency().IsEmpty() {
				currency = rs.Currency()
			}
			groups := make(map[int64]float64)
			for _, line := range rs.OrderLine().Records() {
				var baseTax float64
				for _, tax := range line.Tax().Records() {
					priceReduce := line.PriceUnit() * (1 - line.Discount()/100)
					_, _, _, taxes := tax.ComputeAll(priceReduce+baseTax, currency, line.ProductUomQty(), line.Product(), rs.PartnerShipping())
					for _, t := range taxes {
						groups[tax.TaxGroup().ID()] += t.Amount
					}
					if tax.IncludeBaseAmount() {
						_, _, _, taxesIncl := tax.ComputeAll(priceReduce+baseTax, currency, 1, line.Product(), rs.PartnerShipping())
						baseTax += taxesIncl[0].Amount
					}
				}
			}
			res := make([]accounttypes.TaxGroup, len(groups))
			var i int
			for id, amount := range groups {
				taxGroup := h.AccountTaxGroup().Browse(rs.Env(), []int64{id})
				res[i] = accounttypes.TaxGroup{GroupName: taxGroup.Name(), TaxAmount: amount}
				i++
			}
			sort.Slice(res, func(i, j int) bool {
				return res[i].Sequence < res[j].Sequence
			})
			return res
		})

	h.SaleOrderLine().DeclareModel()
	h.SaleOrderLine().SetDefaultOrder("Order", "LayoutCategory", "Sequence", "ID")

	h.SaleOrderLine().AddFields(map[string]models.FieldDefinition{
		"Order": models.Many2OneField{String: "Order Reference", RelationModel: h.SaleOrder(),
			Required: true, OnDelete: models.Cascade, Index: true, NoCopy: true},
		"Name":     models.TextField{String: "Description", Required: true},
		"Sequence": models.IntegerField{String: "Sequence", Default: models.DefaultValue(10)},
		"InvoiceLines": models.Many2ManyField{String: "Invoice Lines",
			RelationModel: h.AccountInvoiceLine(), NoCopy: true},
		"InvoiceStatus": models.SelectionField{Selection: types.Selection{
			"upselling":  "Upselling Opportunity",
			"invoiced":   "Fully Invoiced",
			"to invoice": "To Invoice",
			"no":         "Nothing to Invoice",
		},
			Compute: h.SaleOrderLine().Methods().ComputeInvoiceStatus(), Stored: true,
			Depends: []string{"State", "ProductUom", "QtyDelivered", "QtyToInvoice", "QtyInvoiced"},
			Default: models.DefaultValue("no"),
		},
		"PriceUnit": models.FloatField{String: "Unit Price", Required: true,
			Digits:   decimalPrecision.GetPrecision("Product Price"),
			OnChange: h.SaleOrderLine().Methods().OnchangeDiscount()},
		"PriceSubtotal": models.FloatField{String: "Subtotal",
			Compute: h.SaleOrderLine().Methods().ComputeAmount(), Stored: true,
			Depends: []string{"ProductUomQty", "Discount", "PriceUnit", "Tax"}},
		"PriceTax": models.FloatField{String: "Taxes",
			Compute: h.SaleOrderLine().Methods().ComputeAmount(), Stored: true,
			Depends: []string{"ProductUomQty", "Discount", "PriceUnit", "Tax"}},
		"PriceTotal": models.FloatField{String: "Total",
			Compute: h.SaleOrderLine().Methods().ComputeAmount(), Stored: true,
			Depends: []string{"ProductUomQty", "Discount", "PriceUnit", "Tax"}},
		"PriceReduce": models.FloatField{String: "Price Reduce",
			Compute: h.SaleOrderLine().Methods().GetPriceReduce(), Stored: true,
			Depends: []string{"PriceUnit", "Discount"}},
		"Tax": models.Many2ManyField{String: "Taxes",
			RelationModel: h.AccountTax(), JSON: "tax_id",
			OnChange: h.SaleOrderLine().Methods().OnchangeDiscount(),
			Filter:   q.AccountTax().Active().Equals(true).Or().Active().Equals(false)},
		"PriceReduceTaxInc": models.FloatField{Compute: h.SaleOrderLine().Methods().GetPriceReduceTax(),
			Stored: true, Depends: []string{"PriceTotal", "ProductUomQty"}},
		"PriceReduceTaxExcl": models.FloatField{Compute: h.SaleOrderLine().Methods().GetPriceReduceNotax(),
			Stored: true, Depends: []string{"PriceSubtotal", "ProductUomQty"}},
		"Discount": models.FloatField{String: "Discount (%)",
			Digits: decimalPrecision.GetPrecision("Discount")},
		"Product": models.Many2OneField{String: "Product", RelationModel: h.ProductProduct(),
			OnChange: h.SaleOrderLine().Methods().ProductChange(),
			Filter:   q.ProductProduct().SaleOk().Equals(true), OnDelete: models.Restrict, Required: true},
		"ProductUomQty": models.FloatField{String: "Quantity",
			Digits: decimalPrecision.GetPrecision("Product Unit of Measure"), Required: true,
			OnChange: h.SaleOrderLine().Methods().ProductUomChange(),
			Default:  models.DefaultValue(1.0)},
		"ProductUom": models.Many2OneField{String: "Unit of Measure", RelationModel: h.ProductUom(),
			OnChange: h.SaleOrderLine().Methods().ProductUomChange(),
			Required: true},
		"QtyDeliveredUpdateable": models.BooleanField{String: "Can Edit Delivered",
			Compute: h.SaleOrderLine().Methods().ComputeQtyDeliveredUpdateable(),
			Depends: []string{"Product.InvoicePolicy", "Order.State"},
			Default: models.DefaultValue(true)},
		"QtyDelivered": models.FloatField{String: "Delivered", NoCopy: true,
			Digits: decimalPrecision.GetPrecision("Product Unit of Measure")},
		"QtyToInvoice": models.FloatField{String: "To Invoice",
			Compute: h.SaleOrderLine().Methods().GetToInvoiceQty(), Stored: true,
			Depends: []string{"QtyInvoiced", "QtyDelivered", "ProductUomQty", "Order.State"},
			Digits:  decimalPrecision.GetPrecision("Product Unit of Measure")},
		"QtyInvoiced": models.FloatField{String: "Invoiced", Compute: h.SaleOrderLine().Methods().GetInvoiceQty(),
			Depends: []string{"InvoiceLines.Invoice.State", "InvoiceLines.Quantity"},
			Stored:  true, Digits: decimalPrecision.GetPrecision("Product Unit of Measure")},
		"Salesman": models.Many2OneField{String: "Salesperson", RelationModel: h.User(), Related: "Order.User",
			ReadOnly: true},
		"Currency": models.Many2OneField{String: "Currency", RelationModel: h.Currency(),
			Related: "Order.Currency", ReadOnly: true},
		"Company": models.Many2OneField{String: "Company", RelationModel: h.Company(), Related: "Order.Company",
			ReadOnly: true},
		"OrderPartner": models.Many2OneField{String: "Customer", RelationModel: h.Partner(),
			Related: "Order.Partner"},
		"AnalyticTags": models.Many2ManyField{String: "Analytic Tags", RelationModel: h.AccountAnalyticTag(),
			JSON: "analytic_tag_ids"},
		"State": models.SelectionField{String: "Order Status", Selection: types.Selection{
			"draft":  "Quotation",
			"sent":   "Quotation Sent",
			"sale":   "Sale Order",
			"done":   "Done",
			"cancel": "Cancelled",
		},
			Related: "Order.State", ReadOnly: true, NoCopy: true, Default: models.DefaultValue("draft")},
		"CustomerLead": models.IntegerField{String: "Delivery Lead Time", Required: true, GoType: new(int),
			Help: "Number of days between the order confirmation and the shipping of the products to the customer"},
		"Procurements": models.One2ManyField{String: "ProcurementIds", RelationModel: h.ProcurementOrder(),
			ReverseFK: "SaleLine", JSON: "procurement_ids"},
		"LayoutCategory":         models.Many2OneField{String: "Section", RelationModel: h.SaleLayoutCategory()},
		"LayoutCategorySequence": models.IntegerField{String: "Layout Sequence"},
	})

	h.SaleOrderLine().Methods().ComputeInvoiceStatus().DeclareMethod(
		`ComputeInvoiceStatus compute the invoice status of a SO line. Possible statuses:

			  - no: if the SO is not in status 'sale' or 'done', we consider that there is nothing to
			    invoice. This is also hte default value if the conditions of no other status is met.

			  - to invoice: we refer to the quantity to invoice of the line. Refer to method
			    'GetToInvoiceQty()' for more information on how this quantity is calculated.

			  - upselling: this is possible only for a product invoiced on ordered quantities for which
			    we delivered more than expected. The could arise if, for example, a project took more
			    time than expected but we decided not to invoice the extra cost to the client. This
			    occurs only in state 'sale', so that when a SO is set to done, the upselling opportunity
			    is removed from the list.

			  - invoiced: the quantity invoiced is larger or equal to the quantity ordered.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			invoiceStatus := "no"
			for _, line := range rs.Records() {
				switch {
				case line.State() != "sale" && line.State() != "done":
					invoiceStatus = "no"
				case !nbutils.IsZero(line.QtyToInvoice(), precision):
					invoiceStatus = "to invoice"
				case line.State() == "sale" && line.Product().InvoicePolicy() == "order" &&
					nbutils.Compare(line.QtyDelivered(), line.ProductUomQty(), precision) > 0:
					invoiceStatus = "upselling"
				case nbutils.Compare(line.QtyInvoiced(), line.ProductUomQty(), precision) >= 0:
					invoiceStatus = "invoiced"
				default:
					invoiceStatus = "no"
				}
			}
			return &h.SaleOrderLineData{
				InvoiceStatus: invoiceStatus,
			}
		})

	h.SaleOrderLine().Methods().ComputeAmount().DeclareMethod(
		`ComputeAmount computes the amounts of the SO line.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			price := rs.PriceUnit() * (1 - rs.Discount()/100)
			_, totalExcluded, totalIncluded, _ := rs.Tax().ComputeAll(price, rs.Order().Currency(), rs.ProductUomQty(), rs.Product(), rs.Order().PartnerShipping())
			return &h.SaleOrderLineData{
				PriceTax:      totalIncluded - totalExcluded,
				PriceTotal:    totalIncluded,
				PriceSubtotal: totalExcluded,
			}
		})

	h.SaleOrderLine().Methods().ComputeQtyDeliveredUpdateable().DeclareMethod(
		`ComputeQtyDeliveredUpdateable checks if the delivered quantity can be updated`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			qtyDeliveredUpdateable := rs.Order().State() == "sale" && rs.Product().TrackService() == "manual" && rs.Product().ExpensePolicy() == "no"
			return &h.SaleOrderLineData{
				QtyDeliveredUpdateable: qtyDeliveredUpdateable,
			}
		})

	h.SaleOrderLine().Methods().GetToInvoiceQty().DeclareMethod(
		`GetToInvoiceQty compute the quantity to invoice. If the invoice policy is order,
		the quantity to invoice is calculated from the ordered quantity. Otherwise, the quantity
		delivered is used.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			if rs.Order().State() != "sale" && rs.Order().State() != "done" {
				return &h.SaleOrderLineData{}
			}
			qtyToInvoice := rs.QtyDelivered() - rs.QtyInvoiced()
			if rs.Product().InvoicePolicy() == "order" {
				qtyToInvoice = rs.ProductUomQty() - rs.QtyInvoiced()
			}
			return &h.SaleOrderLineData{
				QtyToInvoice: qtyToInvoice,
			}
		})

	h.SaleOrderLine().Methods().GetInvoiceQty().DeclareMethod(
		`GetInvoiceQty computes the quantity invoiced. If case of a refund, the quantity invoiced is decreased.
		Note that this is the case only if the refund is generated from the SO and that is intentional: if
        a refund made would automatically decrease the invoiced quantity, then there is a risk of reinvoicing
        it automatically, which may not be wanted at all. That's why the refund has to be created from the SO`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			var qtyInvoiced float64
			for _, invoiceLine := range rs.InvoiceLines().Records() {
				if invoiceLine.Invoice().State() == "cancel" {
					continue
				}
				switch invoiceLine.Invoice().Type() {
				case "out_invoice":
					qtyInvoiced += invoiceLine.Uom().ComputeQuantity(invoiceLine.Quantity(), rs.ProductUom(), true)
				case "out_refund":
					qtyInvoiced -= invoiceLine.Uom().ComputeQuantity(invoiceLine.Quantity(), rs.ProductUom(), true)
				}
			}
			return &h.SaleOrderLineData{
				QtyInvoiced: qtyInvoiced,
			}
		})

	h.SaleOrderLine().Methods().GetPriceReduce().DeclareMethod(
		`GetPriceReduce computes the unit price with discount.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			return &h.SaleOrderLineData{
				PriceReduce: rs.PriceUnit() * (1 - rs.Discount()/100),
			}
		})

	h.SaleOrderLine().Methods().GetPriceReduceTax().DeclareMethod(
		`GetPriceReduceTax computes the total price with tax and discount.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			var price float64
			if rs.ProductUomQty() != 0 {
				price = rs.PriceTotal() / rs.ProductUomQty()
			}
			return &h.SaleOrderLineData{
				PriceReduceTaxInc: price,
			}
		})

	h.SaleOrderLine().Methods().GetPriceReduceNotax().DeclareMethod(
		`GetPriceReduceNotax  computes the total price with discount but without taxes.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			var price float64
			if rs.ProductUomQty() != 0 {
				price = rs.PriceSubtotal() / rs.ProductUomQty()
			}
			return &h.SaleOrderLineData{
				PriceReduceTaxExcl: price,
			}

		})

	h.SaleOrderLine().Methods().ComputeTax().DeclareMethod(
		`ComputeTax computes the taxes applicable for this sale order line.`,
		func(rs h.SaleOrderLineSet) *h.SaleOrderLineData {
			rs.EnsureOne()
			fPos := rs.Order().Partner().PropertyAccountPosition()
			if !rs.Order().FiscalPosition().IsEmpty() {
				fPos = rs.Order().FiscalPosition()
			}
			taxes := rs.Product().Taxes()
			if !rs.Company().IsEmpty() {
				taxes = taxes.Search(q.AccountTax().Company().Equals(rs.Company()))
			}
			if !fPos.IsEmpty() {
				taxes = fPos.MapTax(taxes, rs.Product(), rs.Order().PartnerShipping())
			}
			return &h.SaleOrderLineData{
				Tax: taxes,
			}
		})

	h.SaleOrderLine().Methods().PrepareOrderLineProcurement().DeclareMethod(
		`PrepareOrderLineProcurement returns the data to create the procurement of this sale order line.`,
		func(rs h.SaleOrderLineSet, group h.ProcurementGroupSet) *h.ProcurementOrderData {
			rs.EnsureOne()
			return &h.ProcurementOrderData{
				Name:        rs.Name(),
				Origin:      rs.Order().Name(),
				DatePlanned: rs.Order().DateOrder().AddDate(0, 0, rs.CustomerLead()),
				Product:     rs.Product(),
				ProductQty:  rs.ProductUomQty(),
				ProductUom:  rs.ProductUom(),
				Company:     rs.Company(),
				Group:       group,
				SaleLine:    rs,
			}

		})

	h.SaleOrderLine().Methods().ActionProcurementCreate().DeclareMethod(
		`ActionProcurementCreate creates procurements based on quantity ordered. If the quantity is increased, new
			  procurements are created. If the quantity is decreased, no automated action is taken.`,
		func(rs h.SaleOrderLineSet) h.ProcurementOrderSet {
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			newProcs := h.ProcurementOrder().NewSet(rs.Env())
			for _, line := range rs.Records() {
				if line.State() != "sale" || !line.Product().NeedProcurement() {
					continue
				}
				var qty float64
				for _, proc := range line.Procurements().Records() {
					qty += proc.ProductQty()
				}
				if nbutils.Compare(qty, line.ProductUomQty(), precision) >= 0 {
					continue
				}
				vals := line.PrepareOrderLineProcurement(line.Order().ProcurementGroup())
				vals.ProductQty = line.ProductUomQty() - qty
				newProc := h.ProcurementOrder().NewSet(rs.Env()).WithContext("procurement_autorun_defer", true).Create(vals)
				//new_proc.message_post_with_view('mail.message_origin_link',
				//    values={'self': new_proc, 'origin': line.order_id},
				//    subtype_id=self.env.ref('mail.mt_note').id)
				newProcs = newProcs.Union(newProc)
			}
			newProcs.Run(false)
			return newProcs
		})

	h.SaleOrderLine().Methods().PrepareAddMissingFields().DeclareMethod(
		`PrepareAddMissingFields deduces missing required fields from the onchange`,
		func(rs h.SaleOrderLineSet, values *h.SaleOrderLineData) *h.SaleOrderLineData {
			if !values.Order.IsEmpty() && !values.Product.IsEmpty() {
				// line := h.SaleOrderLine().New(rs.Env(), values)
				// data, _ = line.ProductChange()
				// values = values.Update(data)
				// TODO Implement New and Update
			}
			return values
		})

	h.SaleOrderLine().Methods().Create().Extend("",
		func(rs h.SaleOrderLineSet, data *h.SaleOrderLineData) h.SaleOrderLineSet {
			data = rs.PrepareAddMissingFields(data)
			line := rs.Super().Create(data)
			if line.Order().State() == "sale" {
				line.ActionProcurementCreate()
				//msg = _("Extra line with %s ") % (line.product_id.display_name,)
				//line.order_id.message_post(body=msg)
			}
			return line
		})

	h.SaleOrderLine().Methods().Write().Extend("",
		func(rs h.SaleOrderLineSet, data *h.SaleOrderLineData, fieldsToReset ...models.FieldNamer) bool {
			lines := h.SaleOrderLine().NewSet(rs.Env())
			changedLines := h.SaleOrderLine().NewSet(rs.Env())
			if _, ok := data.Get(h.SaleOrderLine().ProductUomQty()); ok {
				lines = rs.Search(q.SaleOrderLine().State().Equals("sale").And().ProductUomQty().Lower(data.ProductUomQty))
				changedLines = rs.Search(q.SaleOrderLine().State().Equals("sale").And().ProductUomQty().NotEquals(data.ProductUomQty))
				if !changedLines.IsEmpty() {
					/*
						orders = self.mapped('order_id')
						for order in orders:
							order_lines = changed_lines.filtered(lambda x: x.order_id == order)
							msg = ""
							if any([values['product_uom_qty'] < x.product_uom_qty for x in order_lines]):
								msg += "<b>" + _('The ordered quantity has been decreased. Do not forget to take it into account on your invoices and delivery orders.') + '</b>'
							msg += "<ul>"
							for line in order_lines:
								msg += "<li> %s:" % (line.product_id.display_name,)
								msg += "<br/>" + _("Ordered Quantity") + ": %s -> %s <br/>" % (line.product_uom_qty, float(values['product_uom_qty']),)
								if line.product_id.type in ('consu', 'product'):
									msg += _("Delivered Quantity") + ": %s <br/>" % (line.qty_delivered,)
								msg += _("Invoiced Quantity") + ": %s <br/>" % (line.qty_invoiced,)
							msg += "</ul>"
							order.message_post(body=msg)
					*/
				}
			}
			res := rs.Super().Write(data, fieldsToReset...)
			if !lines.IsEmpty() {
				lines.ActionProcurementCreate()
			}
			return res
		})

	h.SaleOrderLine().Methods().PrepareInvoiceLine().DeclareMethod(
		`PrepareInvoiceLine prepares the data to create the new invoice line for a sales order line.`,
		func(rs h.SaleOrderLineSet, qty float64) *h.AccountInvoiceLineData {
			rs.EnsureOne()
			account := rs.Product().Categ().PropertyAccountIncomeCateg()
			if !rs.Product().PropertyAccountIncome().IsEmpty() {
				account = rs.Product().PropertyAccountIncome()
			}
			if account.IsEmpty() {
				panic(rs.T("Please define income account for this product: '%s' (id:%d) - or for its category: '%s'.",
					rs.Product().Name(), rs.Product().ID(), rs.Product().Categ().Name()))
			}
			fPos := rs.Order().Partner().PropertyAccountPosition()
			if !rs.Order().FiscalPosition().IsEmpty() {
				fPos = rs.Order().FiscalPosition()
			}
			if !fPos.IsEmpty() {
				account = fPos.MapAccount(account)
			}
			return &h.AccountInvoiceLineData{
				Name:             rs.Name(),
				Sequence:         rs.Sequence(),
				Origin:           rs.Order().Name(),
				Account:          account,
				PriceUnit:        rs.PriceUnit(),
				Quantity:         qty,
				Discount:         rs.Discount(),
				Uom:              rs.ProductUom(),
				Product:          rs.Product(),
				LayoutCategory:   rs.LayoutCategory(),
				InvoiceLineTaxes: rs.Tax(),
				AccountAnalytic:  rs.Order().Project(),
				AnalyticTags:     rs.AnalyticTags(),
			}
		})

	h.SaleOrderLine().Methods().InvoiceLineCreate().DeclareMethod(
		`InvoiceLineCreate creates an invoice line. The quantity to invoice can be positive (invoice) or negative
			  (refund).`,
		func(rs h.SaleOrderLineSet, invoice h.AccountInvoiceSet, qty float64) {
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			for _, line := range rs.Records() {
				if nbutils.IsZero(qty, precision) {
					continue
				}
				vals := line.PrepareInvoiceLine(qty)
				vals.Invoice = invoice
				vals.SaleLines = line
				h.AccountInvoiceLine().Create(rs.Env(), vals)
			}
		})

	h.SaleOrderLine().Methods().GetDisplayPrice().DeclareMethod(
		`GetDisplayPrice returns the price to display for this order line.`,
		func(rs h.SaleOrderLineSet, product h.ProductProductSet) float64 {
			if rs.Order().Pricelist().DiscountPolicy() == "with_discount" {
				return product.WithContext("pricelist", rs.Order().Pricelist().ID()).Price()
			}
			qty := float64(1)
			if rs.ProductUomQty() != 0 {
				qty = rs.ProductUomQty()
			}
			finalPrice, rule := rs.Order().Pricelist().ComputePriceRule(product, qty, rs.Order().Partner(), dates.Date{}, h.ProductUom().NewSet(rs.Env()))
			basePrice, currency := rs.WithContext("partner_id", rs.Order().Partner().ID()).
				WithContext("date", rs.Order().DateOrder()).
				GetRealPriceCurrency(rs.Product(), rule, rs.ProductUomQty(), rs.ProductUom(), rs.Order().Pricelist())
			if !currency.Equals(rs.Order().Pricelist().Currency()) {
				basePrice = currency.WithContext("partner_id", rs.Order().Partner().ID()).
					WithContext("date", rs.Order().DateOrder()).
					Compute(basePrice, rs.Order().Pricelist().Currency(), true)
			}
			// negative discounts (= surcharge) are included in the display price
			return math.Max(basePrice, finalPrice)
		})

	h.SaleOrderLine().Methods().ProductChange().DeclareMethod(
		`ProductChange updates data when product is changed in the user interface.`,
		func(rs h.SaleOrderLineSet) (*h.SaleOrderLineData, []models.FieldNamer) {
			if rs.Product().IsEmpty() {
				return &h.SaleOrderLineData{}, []models.FieldNamer{}
			}
			qty := rs.ProductUomQty()
			data := rs.ComputeTax()
			fields := []models.FieldNamer{h.SaleOrderLine().Tax()}
			if rs.ProductUom().IsEmpty() || rs.Product().Uom() != rs.ProductUom() {
				data.ProductUom = rs.Product().Uom()
				data.ProductUomQty = 1
				fields = append(fields, h.SaleOrderLine().ProductUom(), h.SaleOrderLine().ProductUomQty())
				qty = 1
			}
			product := rs.Product().
				WithContext("lang", rs.Order().Partner().Lang()).
				WithContext("partner", rs.Order().Partner()).
				WithContext("quantity", qty).
				WithContext("date", rs.Order().DateOrder()).
				WithContext("pricelist", rs.Order().Pricelist()).
				WithContext("uom", rs.ProductUom())

			name := product.NameGet()
			if product.DescriptionSale() != "" {
				name += "\n" + product.DescriptionSale()
			}
			data.Name = name
			fields = append(fields, h.SaleOrderLine().Name())
			if !rs.Order().Pricelist().IsEmpty() && !rs.Order().Partner().IsEmpty() {
				data.PriceUnit = h.AccountTax().NewSet(rs.Env()).FixTaxIncludedPrice(rs.GetDisplayPrice(product),
					product.Taxes(), rs.Tax())
				fields = append(fields, h.SaleOrderLine().PriceUnit())
			}
			return rs.UpdateOnchangeDiscount(data, fields)
			// TODO Add messages and domains when implemented
		})

	h.SaleOrderLine().Methods().ProductUomChange().DeclareMethod(
		`ProductUomChange updates data when quantity or unit of measure is changed in the user interface.`,
		func(rs h.SaleOrderLineSet) (*h.SaleOrderLineData, []models.FieldNamer) {
			if rs.ProductUom().IsEmpty() || rs.Product().IsEmpty() {
				return &h.SaleOrderLineData{}, []models.FieldNamer{h.SaleOrderLine().PriceUnit()}
			}
			if rs.Order().Pricelist().IsEmpty() || !rs.Order().Partner().IsEmpty() {
				return &h.SaleOrderLineData{}, []models.FieldNamer{}
			}
			product := rs.Product().
				WithContext("lang", rs.Order().Partner().Lang()).
				WithContext("partner", rs.Order().Partner()).
				WithContext("quantity", rs.ProductUomQty()).
				WithContext("date", rs.Order().DateOrder()).
				WithContext("pricelist", rs.Order().Pricelist()).
				WithContext("uom", rs.ProductUom()).
				WithContext("fiscal_position", rs.Env().Context().GetInteger("fiscal_position"))
			data := &h.SaleOrderLineData{
				PriceUnit: h.AccountTax().NewSet(rs.Env()).FixTaxIncludedPrice(rs.GetDisplayPrice(product),
					product.Taxes(), rs.Tax()),
			}
			fields := []models.FieldNamer{h.SaleOrderLine().PriceUnit()}
			return rs.UpdateOnchangeDiscount(data, fields)
		})

	h.SaleOrderLine().Methods().Unlink().Extend("",
		func(rs h.SaleOrderLineSet) int64 {
			if !rs.Search(q.SaleOrderLine().State().In([]string{"sale", "done"})).IsEmpty() {
				panic(rs.T("You can not remove a sale order line.\nDiscard changes and try setting the quantity to 0."))
			}
			return rs.Super().Unlink()
		})

	h.SaleOrderLine().Methods().GetDeliveredQty().DeclareMethod(
		`GetDeliveredQty is intended to be overridden in saleStock and saleMRP`,
		func(rs h.SaleOrderLineSet) float64 {
			return 0
		})

	h.SaleOrderLine().Methods().GetRealPriceCurrency().DeclareMethod(
		`GetRealPriceCurrency retrieve the price before applying the pricelist`,
		func(rs h.SaleOrderLineSet, product h.ProductProductSet, rule h.ProductPricelistItemSet, qty float64,
			uom h.ProductUomSet, pricelist h.ProductPricelistSet) (float64, h.CurrencySet) {
			/*def _get_real_price_currency(self, product, rule_id, qty, uom, pricelist_id):
			  """Retrieve the price before applying the pricelist
			      :param obj product: object of current product record
			      :parem float qty: total quentity of product
			      :param tuple price_and_rule: tuple(price, suitable_rule) coming from pricelist computation
			      :param obj uom: unit of measure of current order line
			      :param integer pricelist_id: pricelist id of sale order"""
			  PricelistItem = self.env['product.pricelist.item']
			  field_name = 'lst_price'
			  currency_id = None
			  product_currency = None
			  if rule_id:
			      pricelist_item = PricelistItem.browse(rule_id)
			      if pricelist_item.pricelist_id.discount_policy == 'without_discount':
			          while pricelist_item.base == 'pricelist' and pricelist_item.base_pricelist_id and pricelist_item.base_pricelist_id.discount_policy == 'without_discount':
			              price, rule_id = pricelist_item.base_pricelist_id.with_context(uom=uom.id).get_product_price_rule(product, qty, self.order_id.partner_id)
			              pricelist_item = PricelistItem.browse(rule_id)

			      if pricelist_item.base == 'standard_price':
			          field_name = 'standard_price'
			      if pricelist_item.base == 'pricelist' and pricelist_item.base_pricelist_id:
			          field_name = 'price'
			          product = product.with_context(pricelist=pricelist_item.base_pricelist_id.id)
			          product_currency = pricelist_item.base_pricelist_id.currency_id
			      currency_id = pricelist_item.pricelist_id.currency_id

			  product_currency = product_currency or(product.company_id and product.company_id.currency_id) or self.env.user.company_id.currency_id
			  if not currency_id:
			      currency_id = product_currency
			      cur_factor = 1.0
			  else:
			      if currency_id.id == product_currency.id:
			          cur_factor = 1.0
			      else:
			          cur_factor = currency_id._get_conversion_rate(product_currency, currency_id)

			  product_uom = self.env.context.get('uom') or product.uom_id.id
			  if uom and uom.id != product_uom:
			      # the unit price is in a different uom
			      uom_factor = uom._compute_price(1.0, product.uom_id)
			  else:
			      uom_factor = 1.0

			  return product[field_name] * uom_factor * cur_factor, currency_id.id

			*/
			fieldName := "lstPrice"
			currency := h.Currency().NewSet(rs.Env())
			productCurrency := h.Currency().NewSet(rs.Env())
			if !rule.IsEmpty() {
				pricelistItem := rule
				if pricelistItem.Pricelist().DiscountPolicy() == "without_discount" {
					for pricelistItem.Base() == "pricelist" && !pricelistItem.BasePricelist().IsEmpty() &&
						pricelistItem.BasePricelist().DiscountPolicy() == "without_discount" {
						pricelistItem = pricelistItem.BasePricelist().GetProductPriceRule(product, qty, rs.Order().Partner(), dates.Date{}, uom)
					}
				}
				if pricelistItem.Base() == "standard_price" {
					fieldName = "standardPrice"
				}
				if pricelistItem.Base() == "pricelist" && !pricelistItem.BasePricelist().IsEmpty() {
					fieldName = "Price"
					product = product.WithContext("pricelist", pricelistItem.BasePricelist().ID())
					productCurrency = pricelistItem.BasePricelist().Currency()
				}
				currency = pricelistItem.Pricelist().Currency()
			}

			switch {
			case !productCurrency.IsEmpty():
				break
			case !product.Company().IsEmpty():
				productCurrency = product.Company().Currency()
			default:
				productCurrency = h.User().NewSet(rs.Env()).CurrentUser().Company().Currency()
			}

			curFactor := float64(1)
			if currency.IsEmpty() {
				currency = productCurrency
			}
			if !currency.Equals(productCurrency) {
				curFactor = productCurrency.GetConversionRateTo(currency)
			}

			productUom := product.Uom()
			if rs.Env().Context().HasKey("uom") {
				productUom = h.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
			}
			uomFactor := float64(1)
			if !uom.IsEmpty() && !uom.Equals(productUom) {
				uomFactor = uom.ComputePrice(1, product.Uom())
			}

			return product.Get(fieldName).(float64) * uomFactor * curFactor, currency
		})

	h.SaleOrderLine().Methods().OnchangeDiscount().DeclareMethod(
		`OnchangeDiscount`,
		func(rs h.SaleOrderLineSet) (*h.SaleOrderLineData, []models.FieldNamer) {
			return rs.UpdateOnchangeDiscount(&h.SaleOrderLineData{}, []models.FieldNamer{})
		})

	h.SaleOrderLine().Methods().UpdateOnchangeDiscount().DeclareMethod(
		`UpdateOnchangeDiscount updates the given data and fields with the discount business logic.`,
		func(rs h.SaleOrderLineSet, data *h.SaleOrderLineData, fields []models.FieldNamer) (*h.SaleOrderLineData, []models.FieldNamer) {
			data.Discount = 0
			fields = append(fields, h.SaleOrderLine().Discount())
			if rs.Product().IsEmpty() || rs.ProductUom().IsEmpty() || rs.Order().Partner().IsEmpty() ||
				rs.Order().Pricelist().IsEmpty() || rs.Order().Pricelist().DiscountPolicy() != "without_discount" ||
				!h.User().NewSet(rs.Env()).CurrentUser().HasGroup("sale_group_discount_per_so_line") {
				return data, fields
			}
			qty := float64(1)
			if rs.ProductUomQty() != 0 {
				qty = rs.ProductUomQty()
			}
			price, rule := rs.Order().Pricelist().WithContext("partner_id", rs.Order().Partner()).
				ComputePriceRule(rs.Product(), qty, rs.Order().Partner(), rs.Order().DateOrder().ToDate(), rs.ProductUom())
			newListPrice, currency := rs.WithContext("partner_id", rs.Order().Partner()).
				WithContext("date", rs.Order().DateOrder().ToDate()).
				GetRealPriceCurrency(rs.Product(), rule, rs.ProductUomQty(), rs.ProductUom(), rs.Order().Pricelist())
			if newListPrice == 0 {
				return data, fields
			}
			if !rs.Order().Pricelist().Currency().Equals(currency) {
				// we need new_list_price in the same currency as price, which is in the SO's pricelist's currency
				newListPrice = currency.WithContext("partner_id", rs.Order().Partner()).
					WithContext("date", rs.Order().DateOrder()).
					Compute(newListPrice, rs.Order().Pricelist().Currency(), true)
			}
			discount := (newListPrice - price) / newListPrice * 100
			if discount > 0 {
				data.Discount = discount
			}
			return data, fields
		})

}
