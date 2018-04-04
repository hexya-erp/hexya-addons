// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"time"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.CRMTeam().AddFields(map[string]models.FieldDefinition{
		"UseQuotations": models.BooleanField{String: "Quotations", Default: models.DefaultValue(true),
			OnChange: h.CRMTeam().Methods().OnchangeUseQuotation(),
			Help:     "Check this box to manage quotations in this sales team."},
		"UseInvoices": models.BooleanField{String: "Invoices",
			Help: "Check this box to manage invoices in this sales team."},
		"Invoiced": models.FloatField{String: "Invoiced This Month",
			Compute: h.CRMTeam().Methods().ComputeInvoiced(),
			Help: `Invoice revenue for the current month. This is the amount the sales
team has invoiced this month. It is used to compute the progression ratio
of the current and target revenue on the kanban view.`},
		"InvoicedTarget": models.FloatField{String: "Invoice Target",
			Help: `Target of invoice revenue for the current month. This is the amount the sales
team estimates to be able to invoice this month.`},
		"SalesToInvoiceAmount": models.FloatField{String: "Amount of sales to invoice",
			Compute: h.CRMTeam().Methods().ComputeSalesToInvoiceAmount()},
		"Currency": models.Many2OneField{RelationModel: h.Currency(), Related: "Company.Currency", ReadOnly: true,
			Required: true},
	})

	h.CRMTeam().Methods().ComputeSalesToInvoiceAmount().DeclareMethod(
		`ComputeSalesToInvoiceAmount computes the total amount of sale orders that have not yet been invoiced`,
		func(rs h.CRMTeamSet) *h.CRMTeamData {
			amounts := h.SaleOrder().Search(rs.Env(),
				q.SaleOrder().Team().Equals(rs).
					And().InvoiceStatus().Equals("to invoice")).
				GroupBy(h.SaleOrder().Team()).
				Aggregates(h.SaleOrder().Team(), h.SaleOrder().AmountTotal())
			if len(amounts) == 0 {
				return &h.CRMTeamData{}
			}
			amount, _ := amounts[0].Values.Get("AmountTotal", h.SaleOrder().Underlying())
			return &h.CRMTeamData{
				SalesToInvoiceAmount: amount.(float64),
			}
		})

	h.CRMTeam().Methods().ComputeInvoiced().DeclareMethod(
		`ComputeInvoiced returns the total amount invoiced by this sale team this month.`,
		func(rs h.CRMTeamSet) *h.CRMTeamData {
			firstDayOfMonth := dates.Date{Time: time.Date(dates.Today().Year(), dates.Today().Month(), 1,
				0, 0, 0, 0, time.UTC)}
			invoices := h.AccountInvoice().Search(rs.Env(),
				q.AccountInvoice().State().In([]string{"open", "paid"}).
					And().Team().Equals(rs).
					And().Date().LowerOrEqual(dates.Today()).
					And().Date().GreaterOrEqual(firstDayOfMonth).
					And().Type().In([]string{"out_invoice", "out_refund"})).
				GroupBy(h.AccountInvoice().Team()).
				Aggregates(h.AccountInvoice().Team(), h.AccountInvoice().AmountUntaxedSigned())
			if len(invoices) == 0 {
				return &h.CRMTeamData{}
			}
			amount, _ := invoices[0].Values.Get("AmountUntaxedSigned", h.AccountInvoice().Underlying())
			return &h.CRMTeamData{
				Invoiced: amount.(float64),
			}
		})

	h.CRMTeam().Methods().UpdateInvoicedTarget().DeclareMethod(
		`UpdateInvoicedTarget updates the invoice target with the given value`,
		func(rs h.CRMTeamSet, value float64) bool {
			return rs.Write(&h.CRMTeamData{
				InvoicedTarget: nbutils.Round(value, 1),
			}, h.CRMTeam().InvoicedTarget())
		})

	h.CRMTeam().Methods().OnchangeUseQuotation().DeclareMethod(
		`OnchangeUseQuotation makes sure we use invoices if we use quotations.`,
		func(rs h.CRMTeamSet) (*h.CRMTeamData, []models.FieldNamer) {
			//@api.onchange('use_quotations')
			/*def _onchange_use_quotation(self):
			  if self.use_quotations:
			      self.use_invoices = True
			*/
			var useInvoices bool
			if rs.UseQuotations() {
				useInvoices = true
			}
			return &h.CRMTeamData{
				UseInvoices: useInvoices,
			}, []models.FieldNamer{h.CRMTeam().UseInvoices()}
		})

}
