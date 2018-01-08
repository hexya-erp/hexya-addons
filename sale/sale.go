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
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.SaleOrder().DeclareModel()
	pool.SaleOrder().SetDefaultOrder("DateOrder DESC", "ID DESC")

	pool.SaleOrder().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Order Reference", Required: true, NoCopy: true, /*[ readonly True]*/
			/*[ states {'draft': [('readonly']*/ /*[ False)]}]*/ Index: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.SaleOrder().NewSet(env).T("New")
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
		}, /*[ readonly True]*/ NoCopy: true, Index: true, /*[ track_visibility 'onchange']*/
			Default: models.DefaultValue("draft")},
		"DateOrder": models.DateTimeField{String: "Order Date", Required: true, Index: true, /*[ readonly True]*/
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			NoCopy: true, Default: func(models.Environment, models.FieldMap) interface{} {
				return dates.Now()
			}},
		"ValidityDate": models.DateField{String: "Expiration Date" /*[ readonly True]*/, NoCopy: true,
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: `Manually set the expiration date of your quotation (offer), or it will set the date automatically
based on the template if online quotation is installed.`},
		"ConfirmationDate": models.DateTimeField{ /*[ readonly True]*/ Index: true,
			Help: "Date on which the sale order is confirmed."},
		"User": models.Many2OneField{String: "Salesperson", RelationModel: pool.User(), Index: true, /*[ track_visibility 'onchange']*/
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser()
			}},
		"Partner": models.Many2OneField{String: "Customer", RelationModel: pool.Partner(), /* readonly=true */
			/*[ states {'draft': [('readonly'] [ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			OnChange: pool.SaleOrder().Methods().OnchangePartner(),
			Required: true, Index: true /*[ track_visibility 'always']*/},
		"PartnerInvoice": models.Many2OneField{String: "Invoice Address", RelationModel: pool.Partner(), /* readonly=true */
			Required: true,
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: "Invoice address for current sales order."},
		"PartnerShipping": models.Many2OneField{String: "Delivery Address", RelationModel: pool.Partner(), /* readonly=true */
			Required: true,
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			OnChange: pool.SaleOrder().Methods().OnchangePartnerShipping(),
			Help:     "Delivery address for current sales order."},
		"Pricelist": models.Many2OneField{RelationModel: pool.ProductPricelist(), Required: true, /* readonly=true */
			/*[ states {'draft': [('readonly']*/ /*[ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: "Pricelist for current sales order."},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(),
			Related: "Pricelist.Currency" /* readonly=true */, Required: true},
		"Project": models.Many2OneField{String: "Analytic Account", RelationModel: pool.AccountAnalyticAccount(),
			/* readonly=true */
			/*[ states {'draft': [('readonly'] [ False)]]*/
			/*[ 'sent': [('readonly'] [ False)]}]*/
			Help: "The analytic account related to a sales order.", NoCopy: true},
		"RelatedProject": models.Many2OneField{String: "Analytic Account", RelationModel: pool.AccountAnalyticAccount(),
			Related: "Project", Help: "The analytic account related to a sales order."},
		"OrderLine": models.One2ManyField{String: "Order Lines", RelationModel: pool.SaleOrderLine(),
			ReverseFK: "Order",
			/*[ states {'cancel': [('readonly'] [ True)]]*/
			/*[ 'done': [('readonly'] [ True)]}]*/
			NoCopy: false},
		"InvoiceCount": models.IntegerField{String: "# of Invoices",
			Compute: pool.SaleOrder().Methods().GetInvoiced(),
			Depends: []string{"state", "OrderLine.InvoiceStatus"}, GoType: new(int) /*[ readonly True]*/},
		"Invoices": models.Many2ManyField{String: "Invoices", RelationModel: pool.AccountInvoice(),
			JSON: "invoice_ids", Compute: pool.SaleOrder().Methods().GetInvoiced(),
			Depends: []string{"state", "OrderLine.InvoiceStatus"} /*[ readonly True]*/, NoCopy: true},
		"InvoiceStatus": models.SelectionField{Selection: types.Selection{
			"upselling":  "Upselling Opportunity",
			"invoiced":   "Fully Invoiced",
			"to invoice": "To Invoice",
			"no":         "Nothing to Invoice",
		}, Compute: pool.SaleOrder().Methods().GetInvoiced(),
			Depends: []string{"state", "OrderLine.InvoiceStatus"}, Stored: true /*readonly=true*/},
		"Note": models.TextField{String: "Terms and conditions",
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser().Company().SaleNote()
			}},
		"AmountUntaxed": models.FloatField{String: "Untaxed Amount", Stored: true, /*[ readonly True]*/
			Compute: pool.SaleOrder().Methods().AmountAll(), /*[ track_visibility 'always']*/
			Depends: []string{"OrderLine.PriceTotal"}},
		"AmountTax": models.FloatField{String: "Taxes", Stored: true, /*[ readonly True]*/
			Compute: pool.SaleOrder().Methods().AmountAll(), /*[ track_visibility 'always']*/
			Depends: []string{"OrderLine.PriceTotal"}},
		"AmountTotal": models.FloatField{String: "Total", Stored: true, /*[ readonly True]*/
			Compute: pool.SaleOrder().Methods().AmountAll(), /*[ track_visibility 'always']*/
			Depends: []string{"OrderLine.PriceTotal"}},
		"PaymentTerm": models.Many2OneField{String: "Payment Terms", RelationModel: pool.AccountPaymentTerm()},
		"FiscalPosition": models.Many2OneField{RelationModel: pool.AccountFiscalPosition(),
			OnChange: pool.SaleOrder().Methods().ComputeTax()},
		"Company": models.Many2OneField{RelationModel: pool.Company(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.Company().NewSet(env).CompanyDefaultGet()
			}},
		"Team": models.Many2OneField{String: "Sales Team", RelationModel: pool.CRMTeam(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.CRMTeam().NewSet(env).GetDefaultTeam(pool.User().NewSet(env))
			}},
		"ProcurementGroup": models.Many2OneField{RelationModel: pool.ProcurementGroup(), NoCopy: true},
		"Product":          models.Many2OneField{RelationModel: pool.ProductProduct(), Related: "OrderLine.Product"},
	})

	pool.SaleOrder().Methods().AmountAll().DeclareMethod(
		`AmountAll computes all the amounts of this sale order by summing its sale order lines.`,
		func(rs pool.SaleOrderSet) (*pool.SaleOrderData, []models.FieldNamer) {
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
			return &pool.SaleOrderData{
					AmountUntaxed: rs.Pricelist().Currency().Round(amountUntaxed),
					AmountTax:     rs.Pricelist().Currency().Round(amountTaxed),
					AmountTotal:   amountTaxed + amountUntaxed,
				}, []models.FieldNamer{
					pool.SaleOrder().AmountUntaxed(),
					pool.SaleOrder().AmountTax(),
					pool.SaleOrder().AmountTotal()}
		})

	pool.SaleOrder().Methods().GetInvoiced().DeclareMethod(
		`GetInvoiced computes the invoice status of a SO. Possible statuses:

			  - no: if the SO is not in status 'sale' or 'done', we consider that there is nothing to
			    invoice. This is also hte default value if the conditions of no other status is met.
			  - to invoice: if any SO line is 'to invoice', the whole SO is 'to invoice'
			  - invoiced: if all SO lines are invoiced, the SO is invoiced.
			  - upselling: if all SO lines are invoiced or upselling, the status is upselling.

			  The invoice_ids are obtained thanks to the invoice lines of the SO lines, and we also search
			  for possible refunds created directly from existing invoices. This is necessary since such a
			  refund is not directly linked to the SO.`,
		func(rs pool.SaleOrderSet) (*pool.SaleOrderData, []models.FieldNamer) {
			invoices := pool.AccountInvoice().NewSet(rs.Env())
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
			invoices = invoices.Union(pool.AccountInvoice().Search(rs.Env(),
				pool.AccountInvoice().Origin().Like(rs.Name()).
					And().Name().In(origins).
					And().Type().In([]string{"out_invoice", "out_refund"})))

			refunds := pool.AccountInvoice().NewSet(rs.Env())
			for _, inv := range invoices.Records() {
				refunds = refunds.Union(pool.AccountInvoice().Search(rs.Env(),
					pool.AccountInvoice().Type().Equals("out_refund").
						And().Origin().Equals(inv.Number()).
						And().Origin().IsNotNull().
						And().Journal().Equals(inv.Journal())))
			}
			// Ignore the status of the deposit product
			depositProduct := pool.SaleAdvancePaymentInv().NewSet(rs.Env()).DefaultProduct()
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
			case lineInvoiceStatus["to_invoice"]:
				invoiceStatus = "to_invoice"
			case len(lineInvoiceStatus) == 1 && lineInvoiceStatus["invoiced"]:
				invoiceStatus = "invoiced"
			case len(lineInvoiceStatus) <= 2 && (lineInvoiceStatus["invoiced"] || lineInvoiceStatus["upselling"]):
				invoiceStatus = "upselling"
			default:
				invoiceStatus = "no"
			}
			return &pool.SaleOrderData{
				InvoiceCount:  invoices.Union(refunds).Len(),
				Invoices:      invoices.Union(refunds),
				InvoiceStatus: invoiceStatus,
			}, []models.FieldNamer{}
		})

	pool.SaleOrder().Methods().ComputeTax().DeclareMethod(
		`ComputeTax triggers the recompute of the taxes if the fiscal position is changed on the SO.`,
		func(rs pool.SaleOrderSet) (*pool.SaleOrderData, []models.FieldNamer) {
			//@api.onchange('fiscal_position_id')
			/*def _compute_tax_id(self):
			  """
			  Trigger the recompute of the taxes if the fiscal position is changed on the SO.
			  """
			  for order in self:
			      order.order_line._compute_tax_id()
			*/
			// TODO : need to implement onchange on relation fields first
			return &pool.SaleOrderData{}, []models.FieldNamer{}
		})

	pool.SaleOrder().Methods().GetCustomerLead().DeclareMethod(
		`GetCustomerLead returns the delay to deliver the given product template`,
		func(rs pool.SaleOrderSet, productTmpl pool.ProductTemplateSet) int {
			return 0
		})

	pool.SaleOrder().Methods().ButtonDummy().DeclareMethod(
		`ButtonDummy is a dummy function to force reload of the form on client side.`,
		func(rs pool.SaleOrderSet) bool {
			return true
		})

	pool.SaleOrder().Methods().Unlink().Extend("",
		func(rs pool.SaleOrderSet) int64 {
			for _, order := range rs.Records() {
				if order.State() != "draft" && order.State() != "cancel" {
					panic(rs.T("You can not delete a sent quotation or a sales order! Try to cancel it before."))
				}
			}
			return rs.Super().Unlink()
		})

	//pool.SaleOrder().Methods().TrackSubtype().DeclareMethod(
	//	`TrackSubtype`,
	//	func(rs pool.SaleOrderSet, initvalues interface{}) {
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

	pool.SaleOrder().Methods().OnchangePartnerShipping().DeclareMethod(
		`OnchangePartnerShipping triggers the change of fiscal position when the shipping address is modified.`,
		func(rs pool.SaleOrderSet) (*pool.SaleOrderData, []models.FieldNamer) {
			return &pool.SaleOrderData{
				FiscalPosition: pool.AccountFiscalPosition().NewSet(rs.Env()).GetFiscalPosition(rs.Partner(), rs.PartnerShipping()),
			}, []models.FieldNamer{pool.SaleOrder().FiscalPosition()}
		})

	pool.SaleOrder().Methods().OnchangePartner().DeclareMethod(
		`OnchangePartner updates the following fields when the partner is changed:
		- Pricelist
		- Payment term
		- Invoice address
		- Delivery address
		`,
		func(rs pool.SaleOrderSet) (*pool.SaleOrderData, []models.FieldNamer) {
			if rs.Partner().IsEmpty() {
				return &pool.SaleOrderData{
						PartnerInvoice:  pool.Partner().NewSet(rs.Env()),
						PartnerShipping: pool.Partner().NewSet(rs.Env()),
						PaymentTerm:     pool.AccountPaymentTerm().NewSet(rs.Env()),
						FiscalPosition:  pool.AccountFiscalPosition().NewSet(rs.Env()),
					}, []models.FieldNamer{
						pool.SaleOrder().PartnerInvoice(),
						pool.SaleOrder().PartnerShipping(),
						pool.SaleOrder().PaymentTerm(),
						pool.SaleOrder().FiscalPosition(),
					}
			}
			addr := rs.Partner().AddressGet([]string{"delivery", "invoice"})
			values := &pool.SaleOrderData{
				Pricelist:       rs.Partner().PropertyProductPricelist(),
				PaymentTerm:     rs.Partner().PropertyPaymentTerm(),
				PartnerInvoice:  addr["invoice"],
				PartnerShipping: addr["delivery"],
			}
			fields := []models.FieldNamer{
				pool.SaleOrder().PartnerInvoice(),
				pool.SaleOrder().PartnerShipping(),
				pool.SaleOrder().PaymentTerm(),
				pool.SaleOrder().Pricelist(),
			}
			if pool.User().NewSet(rs.Env()).CurrentUser().Company().SaleNote() != "" {
				values.Note = pool.User().NewSet(rs.Env()).WithContext("lang", rs.Partner().Lang()).
					CurrentUser().Company().SaleNote()
				fields = append(fields, pool.SaleOrder().Note())
			}
			if !rs.Partner().User().IsEmpty() {
				values.User = rs.Partner().User()
				fields = append(fields, pool.SaleOrder().User())
			}
			if !rs.Partner().Team().IsEmpty() {
				values.Team = rs.Partner().Team()
				fields = append(fields, pool.SaleOrder().Team())
			}
			return values, fields
		})

	//pool.SaleOrder().Methods().OnchangePartnerWarning().DeclareMethod(
	//	`OnchangePartnerWarning`,
	//	func(rs pool.SaleOrderSet) {
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

	pool.SaleOrder().Methods().Create().Extend("",
		func(rs pool.SaleOrderSet, data *pool.SaleOrderData) pool.SaleOrderSet {
			if data.Name == "" || data.Name == rs.T("New") {
				seq := pool.Sequence().NewSet(rs.Env())
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

	pool.SaleOrder().Methods().PrepareInvoice().DeclareMethod(
		`PrepareInvoice prepares the data to create the new invoice for a sales order. This method may be
			  overridden to implement custom invoice generation (making sure to call super() to establish
			  a clean extension chain).`,
		func(rs pool.SaleOrderSet) *pool.AccountInvoiceData {
			rs.EnsureOne()
			journal := pool.AccountInvoice().NewSet(rs.Env()).DefaultJournal()
			if journal.IsEmpty() {
				panic(rs.T("Please define an accounting sale journal for this company."))
			}
			fPos := rs.PartnerInvoice().PropertyAccountPosition()
			if !rs.FiscalPosition().IsEmpty() {
				fPos = rs.FiscalPosition()
			}
			invoiceVals := &pool.AccountInvoiceData{
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

	pool.SaleOrder().Methods().PrintQuotation().DeclareMethod(
		`PrintQuotation returns the action to print the quotation report`,
		func(rs pool.SaleOrderSet) *actions.Action {
			//@api.multi
			/*def print_quotation(self):
			  self.filtered(lambda s: s.state == 'draft').write({'state': 'sent'})
			  return self.env['report'].get_action(self, 'sale.report_saleorder')

			*/
			// TODO Implement reports first
			rs.Search(pool.SaleOrder().State().Equals("draft")).SetState("sent")
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	pool.SaleOrder().Methods().ActionViewInvoice().DeclareMethod(
		`ActionViewInvoice returns an action to view the invoice(s) related to this order.
		If there is a single invoice, then it will be opened in a form view, otherwise in list view.`,
		func(rs pool.SaleOrderSet) *actions.Action {
			invoices := pool.AccountInvoice().NewSet(rs.Env())
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
					Type: views.VIEW_TYPE_FORM,
				}}
				action.ResID = invoices.ID()
			default:
				action = &actions.Action{Type: actions.ActionCloseWindow}
			}
			return action
		})

	pool.SaleOrder().Methods().ActionInvoiceCreate().DeclareMethod(
		`ActionInvoiceCreate creates the invoice associated to the SO.

		- If grouped is true, invoices are grouped by sale orders.
		If False, invoices are grouped by (partner_invoice_id, currency)
        - If final is true, refunds will be generated if necessary.`,
		func(rs pool.SaleOrderSet, grouped, final bool) pool.AccountInvoiceSet {
			type keyStruct struct {
				OrderID    int64
				PartnerID  int64
				CurrencyID int64
			}
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			invoices := make(map[keyStruct]pool.AccountInvoiceSet)
			references := make(map[int64]pool.SaleOrderSet)
			for _, order := range rs.Records() {
				groupKey := keyStruct{PartnerID: order.PartnerInvoice().ID(), CurrencyID: order.Currency().ID()}
				if grouped {
					groupKey = keyStruct{OrderID: order.ID()}
				}
				for _, line := range pool.SaleOrderLine().Search(rs.Env(),
					pool.SaleOrderLine().ID().In(order.OrderLine().Ids())).OrderBy("QtyToInvoice").Records() {
					if nbutils.IsZero(line.QtyToInvoice(), precision) {
						continue
					}
					if _, exists := invoices[groupKey]; !exists {
						invData := order.PrepareInvoice()
						invoice := pool.AccountInvoice().Create(rs.Env(), invData)
						references[invoice.ID()] = order
						invoices[groupKey] = invoice
					} else {
						vals := pool.AccountInvoiceData{}
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
				panic(rs.T("There is no invoicable line."))
			}

			res := pool.AccountInvoice().NewSet(rs.Env())
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

	pool.SaleOrder().Methods().ActionDraft().DeclareMethod(
		`ActionDraft sets this sale order back to the draft state.`,
		func(rs pool.SaleOrderSet) bool {
			orders := pool.SaleOrder().NewSet(rs.Env())
			for _, order := range rs.Records() {
				if order.State() != "cancel" && order.State() != "sent" {
					continue
				}
				orders = orders.Union(order)
			}
			orders.Write(&pool.SaleOrderData{
				State:            "draft",
				ProcurementGroup: pool.ProcurementGroup().NewSet(rs.Env()),
			})
			for _, order := range orders.Records() {
				for _, line := range order.OrderLine().Records() {
					for _, proc := range line.Procurements().Records() {
						proc.SetSaleLine(pool.SaleOrderLine().NewSet(rs.Env()))
					}
				}
			}
			return true
		})

	pool.SaleOrder().Methods().ActionCancel().DeclareMethod(
		`ActionCancel cancels this sale order.`,
		func(rs pool.SaleOrderSet) bool {
			rs.SetState("cancel")
			return true
		})

	pool.SaleOrder().Methods().ActionQuotationSend().DeclareMethod(
		`ActionQuotationSend opens a window to compose an email,
		with the edi sale template message loaded by default`,
		func(rs pool.SaleOrderSet) *actions.Action {
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
			rs.Search(pool.SaleOrder().State().Equals("draft")).SetState("sent")
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	pool.SaleOrder().Methods().ForceQuotationSend().DeclareMethod(
		`ForceQuotationSend`,
		func(rs pool.SaleOrderSet) bool {
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

	pool.SaleOrder().Methods().ActionDone().DeclareMethod(
		`ActionDone sets the state of this sale order to done`,
		func(rs pool.SaleOrderSet) bool {
			rs.SetState("done")
			return true
		})

	pool.SaleOrder().Methods().PrepareProcurementGroup().DeclareMethod(
		`PrepareProcurementGroup returns the data that will be used to create the
		procurement group of this sale order`,
		func(rs pool.SaleOrderSet) *pool.ProcurementGroupData {
			return &pool.ProcurementGroupData{
				Name: rs.Name(),
			}
		})

	pool.SaleOrder().Methods().ActionConfirm().DeclareMethod(
		`ActionConfirm confirms this quotation into a sale order`,
		func(rs pool.SaleOrderSet) bool {
			for _, order := range rs.Records() {
				order.Write(&pool.SaleOrderData{
					State:            "sale",
					ConfirmationDate: dates.Now(),
				})
				if rs.Env().Context().HasKey("send_email") {
					rs.ForceQuotationSend()
				}
				order.OrderLine().ActionProcurementCreate()
			}
			autoDone := pool.ConfigParameter().Search(rs.Env(), pool.ConfigParameter().Key().Equals("sale.auto_done_setting"))
			if autoDone.Value() != "" {
				rs.ActionDone()
			}
			return true
		})

	pool.SaleOrder().Methods().CreateAnalyticAccount().DeclareMethod(
		`CreateAnalyticAccount creates the analytic account (project) for this sale order.`,
		func(rs pool.SaleOrderSet, prefix string) {
			for _, order := range rs.Records() {
				name := order.Name()
				if prefix != "" {
					name = fmt.Sprintf("%s: %s", prefix, order.Name())
				}
				analyticAccount := pool.AccountAnalyticAccount().Create(rs.Env(), &pool.AccountAnalyticAccountData{
					Name:    name,
					Code:    order.ClientOrderRef(),
					Company: order.Company(),
					Partner: order.Partner(),
				})
				order.SetProject(analyticAccount)
			}
		})

	pool.SaleOrder().Methods().OrderLinesLayouted().DeclareMethod(
		`OrderLinesLayouted returns this order lines classified by sale_layout_category and separated in
        pages according to the category pagebreaks. Used to render the report.`,
		func(rs pool.SaleOrderSet) {
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

	pool.SaleOrder().Methods().GetTaxAmountByGroup().DeclareMethod(
		`GetTaxAmountByGroup`,
		func(rs pool.SaleOrderSet) []accounttypes.TaxGroup {
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
				taxGroup := pool.AccountTaxGroup().Browse(rs.Env(), []int64{id})
				res[i] = accounttypes.TaxGroup{GroupName: taxGroup.Name(), TaxAmount: amount}
				i++
			}
			sort.Slice(res, func(i, j int) bool {
				return res[i].Sequence < res[j].Sequence
			})
			return res
		})

	pool.SaleOrderLine().DeclareModel()
	pool.SaleOrderLine().SetDefaultOrder("Order", "LayoutCategory", "Sequence", "ID")

	pool.SaleOrderLine().AddFields(map[string]models.FieldDefinition{
		"Order": models.Many2OneField{String: "Order Reference", RelationModel: pool.SaleOrder(),
			Required: true, OnDelete: models.Cascade, Index: true, NoCopy: true},
		"Name":     models.TextField{String: "Description", Required: true},
		"Sequence": models.IntegerField{String: "Sequence", Default: models.DefaultValue(10)},
		"InvoiceLines": models.Many2ManyField{String: "Invoice Lines",
			RelationModel: pool.AccountInvoiceLine(), NoCopy: true},
		"InvoiceStatus": models.SelectionField{Selection: types.Selection{
			"upselling":  "Upselling Opportunity",
			"invoiced":   "Fully Invoiced",
			"to invoice": "To Invoice",
			"no":         "Nothing to Invoice",
		},
			Compute: pool.SaleOrderLine().Methods().ComputeInvoiceStatus(), Stored: true, /* readonly=true */
			Depends: []string{"State", "ProductUom", "QtyDelivered", "QtyToInvoice", "QtyInvoiced"},
			Default: models.DefaultValue("no"),
		},
		"PriceUnit": models.FloatField{String: "Unit Price", Required: true,
			Digits:   decimalPrecision.GetPrecision("Product Price"),
			OnChange: pool.SaleOrderLine().Methods().OnchangeDiscount()},
		"PriceSubtotal": models.FloatField{String: "Subtotal",
			Compute: pool.SaleOrderLine().Methods().ComputeAmount() /*[ readonly True]*/, Stored: true,
			Depends: []string{"ProductUomQty", "Discount", "PriceUnit", "Tax"}},
		"PriceTax": models.FloatField{String: "Taxes",
			Compute: pool.SaleOrderLine().Methods().ComputeAmount() /*[ readonly True]*/, Stored: true,
			Depends: []string{"ProductUomQty", "Discount", "PriceUnit", "Tax"}},
		"PriceTotal": models.FloatField{String: "Total",
			Compute: pool.SaleOrderLine().Methods().ComputeAmount() /*[ readonly True]*/, Stored: true,
			Depends: []string{"ProductUomQty", "Discount", "PriceUnit", "Tax"}},
		"PriceReduce": models.FloatField{String: "Price Reduce",
			Compute: pool.SaleOrderLine().Methods().GetPriceReduce() /*[ readonly True]*/, Stored: true,
			Depends: []string{"PriceUnit", "Discount"}},
		"Tax": models.Many2ManyField{String: "Taxes",
			RelationModel: pool.AccountTax(), JSON: "tax_id",
			OnChange: pool.SaleOrderLine().Methods().OnchangeDiscount(),
			Filter:   pool.AccountTax().Active().Equals(true).Or().Active().Equals(false)},
		"PriceReduceTaxInc": models.FloatField{Compute: pool.SaleOrderLine().Methods().GetPriceReduceTax(), /*[ readonly True]*/
			Stored: true, Depends: []string{"PriceTotal", "ProductUomQty"}},
		"PriceReduceTaxExcl": models.FloatField{Compute: pool.SaleOrderLine().Methods().GetPriceReduceNotax(), /*[ readonly True]*/
			Stored: true, Depends: []string{"PriceSubtotal", "ProductUomQty"}},
		"Discount": models.FloatField{String: "Discount (%)",
			Digits: decimalPrecision.GetPrecision("Discount")},
		"Product": models.Many2OneField{String: "Product", RelationModel: pool.ProductProduct(),
			OnChange: pool.SaleOrderLine().Methods().ProductChange(),
			Filter:   pool.ProductProduct().SaleOk().Equals(true), OnDelete: models.Restrict, Required: true},
		"ProductUomQty": models.FloatField{String: "Quantity",
			Digits: decimalPrecision.GetPrecision("Product Unit of Measure"), Required: true,
			OnChange: pool.SaleOrderLine().Methods().ProductUomChange(),
			Default:  models.DefaultValue(1.0)},
		"ProductUom": models.Many2OneField{String: "Unit of Measure", RelationModel: pool.ProductUom(),
			OnChange: pool.SaleOrderLine().Methods().ProductUomChange(),
			Required: true},
		"QtyDeliveredUpdateable": models.BooleanField{String: "Can Edit Delivered",
			Compute: pool.SaleOrderLine().Methods().ComputeQtyDeliveredUpdateable(), /*[ readonly True]*/
			Depends: []string{"Product.InvoicePolicy", "Order.State"},
			Default: models.DefaultValue(true)},
		"QtyDelivered": models.FloatField{String: "Delivered", NoCopy: true,
			Digits: decimalPrecision.GetPrecision("Product Unit of Measure")},
		"QtyToInvoice": models.FloatField{String: "To Invoice",
			Compute: pool.SaleOrderLine().Methods().GetToInvoiceQty(), Stored: true, /*[ readonly True]*/
			Depends: []string{"QtyInvoiced", "QtyDelivered", "ProductUomQty", "Order.State"},
			Digits:  decimalPrecision.GetPrecision("Product Unit of Measure")},
		"QtyInvoiced": models.FloatField{String: "Invoiced", Compute: pool.SaleOrderLine().Methods().GetInvoiceQty(),
			Depends: []string{"InvoiceLines.Invoice.State", "InvoiceLines.Quantity"},
			Stored:  true /*[ readonly True]*/, Digits: decimalPrecision.GetPrecision("Product Unit of Measure")},
		"Salesman": models.Many2OneField{String: "Salesperson", RelationModel: pool.User(), Related: "Order.User" /* readonly=true */},
		"Currency": models.Many2OneField{String: "Currency", RelationModel: pool.Currency(),
			Related: "Order.Currency" /* readonly=true */},
		"Company": models.Many2OneField{String: "Company", RelationModel: pool.Company(), Related: "Order.Company" /* readonly=true */},
		"OrderPartner": models.Many2OneField{String: "Customer", RelationModel: pool.Partner(),
			Related: "Order.Partner"},
		"AnalyticTags": models.Many2ManyField{String: "Analytic Tags", RelationModel: pool.AccountAnalyticTag(),
			JSON: "analytic_tag_ids"},
		"State": models.SelectionField{String: "Order Status", Selection: types.Selection{
			"draft":  "Quotation",
			"sent":   "Quotation Sent",
			"sale":   "Sale Order",
			"done":   "Done",
			"cancel": "Cancelled",
		},
			Related: "Order.State" /*[ readonly True]*/, NoCopy: true, Default: models.DefaultValue("draft")},
		"CustomerLead": models.IntegerField{String: "Delivery Lead Time", Required: true, GoType: new(int),
			Help: "Number of days between the order confirmation and the shipping of the products to the customer"},
		"Procurements": models.One2ManyField{String: "ProcurementIds", RelationModel: pool.ProcurementOrder(),
			ReverseFK: "SaleLine", JSON: "procurement_ids"},
		"LayoutCategory":         models.Many2OneField{String: "Section", RelationModel: pool.SaleLayoutCategory()},
		"LayoutCategorySequence": models.IntegerField{String: "Layout Sequence"},
	})

	pool.SaleOrderLine().Methods().ComputeInvoiceStatus().DeclareMethod(
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
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			invoiceStatus := "no"
			for _, line := range rs.Records() {
				switch {
				case line.State() != "sale" && line.State() != "done":
					invoiceStatus = "no"
				case !nbutils.IsZero(line.QtyToInvoice(), precision):
					invoiceStatus = "to_invoice"
				case line.State() == "sale" && line.Product().InvoicePolicy() == "order" &&
					nbutils.Compare(line.QtyDelivered(), line.ProductUomQty(), precision) > 0:
					invoiceStatus = "upselling"
				case nbutils.Compare(line.QtyInvoiced(), line.ProductUomQty(), precision) >= 0:
					invoiceStatus = "invoiced"
				default:
					invoiceStatus = "no"
				}
			}
			return &pool.SaleOrderLineData{
				InvoiceStatus: invoiceStatus,
			}, []models.FieldNamer{pool.SaleOrderLine().InvoiceStatus()}
		})

	pool.SaleOrderLine().Methods().ComputeAmount().DeclareMethod(
		`ComputeAmount computes the amounts of the SO line.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			price := rs.PriceUnit() * (1 - rs.Discount()/100)
			_, totalExcluded, totalIncluded, _ := rs.Tax().ComputeAll(price, rs.Order().Currency(), rs.ProductUomQty(), rs.Product(), rs.Order().PartnerShipping())
			return &pool.SaleOrderLineData{
				PriceTax:      totalIncluded - totalExcluded,
				PriceTotal:    totalIncluded,
				PriceSubtotal: totalExcluded,
			}, []models.FieldNamer{pool.SaleOrderLine().PriceTax(), pool.SaleOrderLine().PriceTotal(), pool.SaleOrderLine().PriceSubtotal()}

		})

	pool.SaleOrderLine().Methods().ComputeQtyDeliveredUpdateable().DeclareMethod(
		`ComputeQtyDeliveredUpdateable checks if the delivered quantity can be updated`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			qtyDeliveredUpdateable := rs.Order().State() == "sale" && rs.Product().TrackService() == "manual" && rs.Product().ExpensePolicy() == "no"
			return &pool.SaleOrderLineData{
				QtyDeliveredUpdateable: qtyDeliveredUpdateable,
			}, []models.FieldNamer{pool.SaleOrderLine().QtyDeliveredUpdateable()}
		})

	pool.SaleOrderLine().Methods().GetToInvoiceQty().DeclareMethod(
		`GetToInvoiceQty compute the quantity to invoice. If the invoice policy is order,
		the quantity to invoice is calculated from the ordered quantity. Otherwise, the quantity
		delivered is used.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			if rs.Order().State() != "sale" && rs.Order().State() != "done" {
				return &pool.SaleOrderLineData{}, []models.FieldNamer{pool.SaleOrderLine().QtyToInvoice()}
			}
			qtyToInvoice := rs.QtyDelivered() - rs.QtyInvoiced()
			if rs.Product().InvoicePolicy() == "order" {
				qtyToInvoice = rs.ProductUomQty() - rs.QtyInvoiced()
			}
			return &pool.SaleOrderLineData{
				QtyToInvoice: qtyToInvoice,
			}, []models.FieldNamer{pool.SaleOrderLine().QtyToInvoice()}
		})

	pool.SaleOrderLine().Methods().GetInvoiceQty().DeclareMethod(
		`GetInvoiceQty computes the quantity invoiced. If case of a refund, the quantity invoiced is decreased.
		Note that this is the case only if the refund is generated from the SO and that is intentional: if
        a refund made would automatically decrease the invoiced quantity, then there is a risk of reinvoicing
        it automatically, which may not be wanted at all. That's why the refund has to be created from the SO`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
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
			return &pool.SaleOrderLineData{
				QtyInvoiced: qtyInvoiced,
			}, []models.FieldNamer{pool.SaleOrderLine().QtyInvoiced()}
		})

	pool.SaleOrderLine().Methods().GetPriceReduce().DeclareMethod(
		`GetPriceReduce computes the unit price with discount.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			return &pool.SaleOrderLineData{
				PriceReduce: rs.PriceUnit() * (1 - rs.Discount()/100),
			}, []models.FieldNamer{pool.SaleOrderLine().PriceReduce()}
		})

	pool.SaleOrderLine().Methods().GetPriceReduceTax().DeclareMethod(
		`GetPriceReduceTax computes the total price with tax and discount.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			var price float64
			if rs.ProductUomQty() != 0 {
				price = rs.PriceTotal() / rs.ProductUomQty()
			}
			return &pool.SaleOrderLineData{
				PriceReduceTaxInc: price,
			}, []models.FieldNamer{pool.SaleOrderLine().PriceReduceTaxInc()}
		})

	pool.SaleOrderLine().Methods().GetPriceReduceNotax().DeclareMethod(
		`GetPriceReduceNotax  computes the total price with discount but without taxes.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			var price float64
			if rs.ProductUomQty() != 0 {
				price = rs.PriceSubtotal() / rs.ProductUomQty()
			}
			return &pool.SaleOrderLineData{
				PriceReduceTaxExcl: price,
			}, []models.FieldNamer{pool.SaleOrderLine().PriceReduceTaxExcl()}

		})

	pool.SaleOrderLine().Methods().ComputeTax().DeclareMethod(
		`ComputeTax computes the taxes applicable for this sale order line.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			rs.EnsureOne()
			fPos := rs.Order().Partner().PropertyAccountPosition()
			if !rs.Order().FiscalPosition().IsEmpty() {
				fPos = rs.Order().FiscalPosition()
			}
			taxes := rs.Product().Taxes()
			if !rs.Company().IsEmpty() {
				taxes = taxes.Search(pool.AccountTax().Company().Equals(rs.Company()))
			}
			if !fPos.IsEmpty() {
				taxes = fPos.MapTax(taxes, rs.Product(), rs.Order().PartnerShipping())
			}
			return &pool.SaleOrderLineData{
				Tax: taxes,
			}, []models.FieldNamer{pool.SaleOrderLine().Tax()}
		})

	pool.SaleOrderLine().Methods().PrepareOrderLineProcurement().DeclareMethod(
		`PrepareOrderLineProcurement returns the data to create the procurement of this sale order line.`,
		func(rs pool.SaleOrderLineSet, group pool.ProcurementGroupSet) *pool.ProcurementOrderData {
			rs.EnsureOne()
			return &pool.ProcurementOrderData{
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

	pool.SaleOrderLine().Methods().ActionProcurementCreate().DeclareMethod(
		`ActionProcurementCreate creates procurements based on quantity ordered. If the quantity is increased, new
			  procurements are created. If the quantity is decreased, no automated action is taken.`,
		func(rs pool.SaleOrderLineSet) pool.ProcurementOrderSet {
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			newProcs := pool.ProcurementOrder().NewSet(rs.Env())
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
				newProc := pool.ProcurementOrder().NewSet(rs.Env()).WithContext("procurement_autorun_defer", true).Create(vals)
				//new_proc.message_post_with_view('mail.message_origin_link',
				//    values={'self': new_proc, 'origin': line.order_id},
				//    subtype_id=self.env.ref('mail.mt_note').id)
				newProcs = newProcs.Union(newProc)
			}
			newProcs.Run(false)
			return newProcs
		})

	pool.SaleOrderLine().Methods().PrepareAddMissingFields().DeclareMethod(
		`PrepareAddMissingFields deduces missing required fields from the onchange`,
		func(rs pool.SaleOrderLineSet, values *pool.SaleOrderLineData) *pool.SaleOrderLineData {
			if !values.Order.IsEmpty() && !values.Product.IsEmpty() {
				// line := pool.SaleOrderLine().New(rs.Env(), values)
				// data, _ = line.ProductChange()
				// values = values.Update(data)
				// TODO Implement New and Update
			}
			return values
		})

	pool.SaleOrderLine().Methods().Create().Extend("",
		func(rs pool.SaleOrderLineSet, data *pool.SaleOrderLineData) pool.SaleOrderLineSet {
			data = rs.PrepareAddMissingFields(data)
			line := rs.Super().Create(data)
			if line.Order().State() == "sale" {
				line.ActionProcurementCreate()
				//msg = _("Extra line with %s ") % (line.product_id.display_name,)
				//line.order_id.message_post(body=msg)
			}
			return line
		})

	pool.SaleOrderLine().Methods().Write().Extend("",
		func(rs pool.SaleOrderLineSet, data *pool.SaleOrderLineData, fieldsToReset ...models.FieldNamer) bool {
			lines := pool.SaleOrderLine().NewSet(rs.Env())
			changedLines := pool.SaleOrderLine().NewSet(rs.Env())
			if _, ok := data.Get(pool.SaleOrderLine().ProductUomQty()); ok {
				lines = rs.Search(pool.SaleOrderLine().State().Equals("sale").And().ProductUomQty().Lower(data.ProductUomQty))
				changedLines = rs.Search(pool.SaleOrderLine().State().Equals("sale").And().ProductUomQty().NotEquals(data.ProductUomQty))
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

	pool.SaleOrderLine().Methods().PrepareInvoiceLine().DeclareMethod(
		`PrepareInvoiceLine prepares the data to create the new invoice line for a sales order line.`,
		func(rs pool.SaleOrderLineSet, qty float64) *pool.AccountInvoiceLineData {
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
			return &pool.AccountInvoiceLineData{
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

	pool.SaleOrderLine().Methods().InvoiceLineCreate().DeclareMethod(
		`InvoiceLineCreate creates an invoice line. The quantity to invoice can be positive (invoice) or negative
			  (refund).`,
		func(rs pool.SaleOrderLineSet, invoice pool.AccountInvoiceSet, qty float64) {
			precision := decimalPrecision.GetPrecision("Product Unit of Measure").ToPrecision()
			for _, line := range rs.Records() {
				if nbutils.IsZero(qty, precision) {
					continue
				}
				vals := line.PrepareInvoiceLine(qty)
				vals.Invoice = invoice
				vals.SaleLines = line
				pool.AccountInvoiceLine().Create(rs.Env(), vals)
			}
		})

	pool.SaleOrderLine().Methods().GetDisplayPrice().DeclareMethod(
		`GetDisplayPrice returns the price to display for this order line.`,
		func(rs pool.SaleOrderLineSet, product pool.ProductProductSet) float64 {
			if rs.Order().Pricelist().DiscountPolicy() == "with_discount" {
				return product.WithContext("pricelist", rs.Order().Pricelist().ID()).Price()
			}
			qty := float64(1)
			if rs.ProductUomQty() != 0 {
				qty = rs.ProductUomQty()
			}
			finalPrice, rule := rs.Order().Pricelist().ComputePriceRule(product, qty, rs.Order().Partner(), dates.Date{}, pool.ProductUom().NewSet(rs.Env()))
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

	pool.SaleOrderLine().Methods().ProductChange().DeclareMethod(
		`ProductChange updates data when product is changed in the user interface.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			if rs.Product().IsEmpty() {
				return &pool.SaleOrderLineData{}, []models.FieldNamer{}
			}
			qty := rs.ProductUomQty()
			data, fields := rs.ComputeTax()
			if rs.ProductUom().IsEmpty() || rs.Product().Uom() != rs.ProductUom() {
				data.ProductUom = rs.Product().Uom()
				data.ProductUomQty = 1
				fields = append(fields, pool.SaleOrderLine().ProductUom(), pool.SaleOrderLine().ProductUomQty())
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
			fields = append(fields, pool.SaleOrderLine().Name())
			if !rs.Order().Pricelist().IsEmpty() && !rs.Order().Partner().IsEmpty() {
				data.PriceUnit = pool.AccountTax().NewSet(rs.Env()).FixTaxIncludedPrice(rs.GetDisplayPrice(product),
					product.Taxes(), rs.Tax())
				fields = append(fields, pool.SaleOrderLine().PriceUnit())
			}
			return rs.UpdateOnchangeDiscount(data, fields)
			// TODO Add messages and domains when implemented
		})

	pool.SaleOrderLine().Methods().ProductUomChange().DeclareMethod(
		`ProductUomChange updates data when quantity or unit of measure is changed in the user interface.`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			if rs.ProductUom().IsEmpty() || rs.Product().IsEmpty() {
				return &pool.SaleOrderLineData{}, []models.FieldNamer{pool.SaleOrderLine().PriceUnit()}
			}
			if rs.Order().Pricelist().IsEmpty() || !rs.Order().Partner().IsEmpty() {
				return &pool.SaleOrderLineData{}, []models.FieldNamer{}
			}
			product := rs.Product().
				WithContext("lang", rs.Order().Partner().Lang()).
				WithContext("partner", rs.Order().Partner()).
				WithContext("quantity", rs.ProductUomQty()).
				WithContext("date", rs.Order().DateOrder()).
				WithContext("pricelist", rs.Order().Pricelist()).
				WithContext("uom", rs.ProductUom()).
				WithContext("fiscal_position", rs.Env().Context().GetInteger("fiscal_position"))
			data := &pool.SaleOrderLineData{
				PriceUnit: pool.AccountTax().NewSet(rs.Env()).FixTaxIncludedPrice(rs.GetDisplayPrice(product),
					product.Taxes(), rs.Tax()),
			}
			fields := []models.FieldNamer{pool.SaleOrderLine().PriceUnit()}
			return rs.UpdateOnchangeDiscount(data, fields)
		})

	pool.SaleOrderLine().Methods().Unlink().Extend("",
		func(rs pool.SaleOrderLineSet) int64 {
			if !rs.Search(pool.SaleOrderLine().State().In([]string{"sale", "done"})).IsEmpty() {
				panic(rs.T("You can not remove a sale order line.\nDiscard changes and try setting the quantity to 0."))
			}
			return rs.Super().Unlink()
		})

	pool.SaleOrderLine().Methods().GetDeliveredQty().DeclareMethod(
		`GetDeliveredQty is intended to be overridden in saleStock and saleMRP`,
		func(rs pool.SaleOrderLineSet) float64 {
			return 0
		})

	pool.SaleOrderLine().Methods().GetRealPriceCurrency().DeclareMethod(
		`GetRealPriceCurrency retrieve the price before applying the pricelist`,
		func(rs pool.SaleOrderLineSet, product pool.ProductProductSet, rule pool.ProductPricelistItemSet, qty float64,
			uom pool.ProductUomSet, pricelist pool.ProductPricelistSet) (float64, pool.CurrencySet) {
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
			currency := pool.Currency().NewSet(rs.Env())
			productCurrency := pool.Currency().NewSet(rs.Env())
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
				productCurrency = pool.User().NewSet(rs.Env()).CurrentUser().Company().Currency()
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
				productUom = pool.ProductUom().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("uom")})
			}
			uomFactor := float64(1)
			if !uom.IsEmpty() && !uom.Equals(productUom) {
				uomFactor = uom.ComputePrice(1, product.Uom())
			}

			return product.Get(fieldName).(float64) * uomFactor * curFactor, currency
		})

	pool.SaleOrderLine().Methods().OnchangeDiscount().DeclareMethod(
		`OnchangeDiscount`,
		func(rs pool.SaleOrderLineSet) (*pool.SaleOrderLineData, []models.FieldNamer) {
			return rs.UpdateOnchangeDiscount(&pool.SaleOrderLineData{}, []models.FieldNamer{})
		})

	pool.SaleOrderLine().Methods().UpdateOnchangeDiscount().DeclareMethod(
		`UpdateOnchangeDiscount updates the given data and fields with the discount business logic.`,
		func(rs pool.SaleOrderLineSet, data *pool.SaleOrderLineData, fields []models.FieldNamer) (*pool.SaleOrderLineData, []models.FieldNamer) {
			data.Discount = 0
			fields = append(fields, pool.SaleOrderLine().Discount())
			if rs.Product().IsEmpty() || rs.ProductUom().IsEmpty() || rs.Order().Partner().IsEmpty() ||
				rs.Order().Pricelist().IsEmpty() || rs.Order().Pricelist().DiscountPolicy() != "without_discount" ||
				!pool.User().NewSet(rs.Env()).CurrentUser().HasGroup("sale_group_discount_per_so_line") {
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
