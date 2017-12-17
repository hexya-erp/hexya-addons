// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"time"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.CRMTeam().AddFields(map[string]models.FieldDefinition{
		"UseQuotations": models.BooleanField{String: "Quotations", Default: models.DefaultValue(true),
			OnChange: pool.CRMTeam().Methods().OnchangeUseQuotation(),
			Help:     "Check this box to manage quotations in this sales team."},
		"UseInvoices": models.BooleanField{String: "Invoices",
			Help: "Check this box to manage invoices in this sales team."},
		"Invoiced": models.FloatField{String: "Invoiced This Month",
			Compute: pool.CRMTeam().Methods().ComputeInvoiced(), /*[ readonly True]*/
			Help: `Invoice revenue for the current month. This is the amount the sales
team has invoiced this month. It is used to compute the progression ratio
of the current and target revenue on the kanban view.`},
		"InvoicedTarget": models.FloatField{String: "Invoice Target",
			Help: `Target of invoice revenue for the current month. This is the amount the sales
team estimates to be able to invoice this month.`},
		"SalesToInvoiceAmount": models.FloatField{String: "Amount of sales to invoice",
			Compute: pool.CRMTeam().Methods().ComputeSalesToInvoiceAmount() /*[ readonly True]*/},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(), Related: "Company.Currency", /* readonly=true */
			Required: true},
	})

	pool.CRMTeam().Methods().ComputeSalesToInvoiceAmount().DeclareMethod(
		`ComputeSalesToInvoiceAmount computes the total amount of sale orders that have not yet been invoiced`,
		func(rs pool.CRMTeamSet) (*pool.CRMTeamData, []models.FieldNamer) {
			amounts := pool.SaleOrder().Search(rs.Env(),
				pool.SaleOrder().Team().Equals(rs).
					And().InvoiceStatus().Equals("to invoice")).
				GroupBy(pool.SaleOrder().Team()).
				Aggregates(pool.SaleOrder().Team(), pool.SaleOrder().AmountTotal())
			if len(amounts) == 0 {
				return &pool.CRMTeamData{}, []models.FieldNamer{pool.CRMTeam().SalesToInvoiceAmount()}
			}
			amount, _ := amounts[0].Values.Get("AmountTotal", pool.SaleOrder().Underlying())
			return &pool.CRMTeamData{
				SalesToInvoiceAmount: amount.(float64),
			}, []models.FieldNamer{pool.CRMTeam().SalesToInvoiceAmount()}
		})

	pool.CRMTeam().Methods().ComputeInvoiced().DeclareMethod(
		`ComputeInvoiced returns the total amount invoiced by this sale team this month.`,
		func(rs pool.CRMTeamSet) (*pool.CRMTeamData, []models.FieldNamer) {
			firstDayOfMonth := dates.Date{Time: time.Date(dates.Today().Year(), dates.Today().Month(), 1,
				0, 0, 0, 0, time.UTC)}
			invoices := pool.AccountInvoice().Search(rs.Env(),
				pool.AccountInvoice().State().In([]string{"open", "paid"}).
					And().Team().Equals(rs).
					And().Date().LowerOrEqual(dates.Today()).
					And().Date().GreaterOrEqual(firstDayOfMonth).
					And().Type().In([]string{"out_invoice", "out_refund"})).
				GroupBy(pool.AccountInvoice().Team()).
				Aggregates(pool.AccountInvoice().Team(), pool.AccountInvoice().AmountUntaxedSigned())
			if len(invoices) == 0 {
				return &pool.CRMTeamData{}, []models.FieldNamer{pool.CRMTeam().Invoiced()}
			}
			amount, _ := invoices[0].Values.Get("AmountUntaxedSigned", pool.AccountInvoice().Underlying())
			return &pool.CRMTeamData{
				Invoiced: amount.(float64),
			}, []models.FieldNamer{pool.CRMTeam().Invoiced()}
		})

	pool.CRMTeam().Methods().UpdateInvoicedTarget().DeclareMethod(
		`UpdateInvoicedTarget updates the invoice target with the given value`,
		func(rs pool.CRMTeamSet, value float64) bool {
			return rs.Write(&pool.CRMTeamData{
				InvoicedTarget: nbutils.Round(value, 1),
			}, pool.CRMTeam().InvoicedTarget())
		})

	pool.CRMTeam().Methods().OnchangeUseQuotation().DeclareMethod(
		`OnchangeUseQuotation makes sure we use invoices if we use quotations.`,
		func(rs pool.CRMTeamSet) (*pool.CRMTeamData, []models.FieldNamer) {
			//@api.onchange('use_quotations')
			/*def _onchange_use_quotation(self):
			  if self.use_quotations:
			      self.use_invoices = True
			*/
			var useInvoices bool
			if rs.UseQuotations() {
				useInvoices = true
			}
			return &pool.CRMTeamData{
				UseInvoices: useInvoices,
			}, []models.FieldNamer{pool.CRMTeam().UseInvoices()}
		})

}
