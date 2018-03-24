// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.AccountInvoiceReport().DeclareModel()
	h.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency().DeclareMethod(
		`ComputeAmountsInUserCurrency`,
		func(rs h.AccountInvoiceReportSet) *h.AccountInvoiceReportData {
			//@api.depends('currency_id','date','price_total','price_average','residual')
			/*def _compute_amounts_in_user_currency(self):
			  """Compute the amounts in the currency of the user
			  """
			  context = dict(self._context or {})
			  user_currency_id = self.env.user.company_id.currency_id
			  currency_rate_id = self.env['res.currency.rate'].search([
			      ('rate', '=', 1),
			      '|', ('company_id', '=', self.env.user.company_id.id), ('company_id', '=', False)], limit=1)
			  base_currency_id = currency_rate_id.currency_id
			  ctx = context.copy()
			  for record in self:
			      ctx['date'] = record.date
			      record.user_currency_price_total = base_currency_id.with_context(ctx).compute(record.price_total, user_currency_id)
			      record.user_currency_price_average = base_currency_id.with_context(ctx).compute(record.price_average, user_currency_id)
			      record.user_currency_residual = base_currency_id.with_context(ctx).compute(record.residual, user_currency_id)
			*/
			return &h.AccountInvoiceReportData{}
		})
	h.AccountInvoiceReport().AddFields(map[string]models.FieldDefinition{
		"Date":       models.DateField{String: "Date", ReadOnly: true},
		"Product":    models.Many2OneField{RelationModel: h.ProductProduct(), ReadOnly: true},
		"ProductQty": models.FloatField{String: "Product Quantity", ReadOnly: true},
		"UomName":    models.CharField{String: "Reference Unit of Measure", ReadOnly: true},
		"PaymentTerm": models.Many2OneField{String: "Payment Terms", RelationModel: h.AccountPaymentTerm(),
			ReadOnly: true},
		"FiscalPosition": models.Many2OneField{RelationModel: h.AccountFiscalPosition(), ReadOnly: true},
		"Currency":       models.Many2OneField{RelationModel: h.Currency(), ReadOnly: true},
		"Categ": models.Many2OneField{String: "Product Category", RelationModel: h.ProductCategory(),
			ReadOnly: true},
		"Journal": models.Many2OneField{RelationModel: h.AccountJournal(), ReadOnly: true},
		"Partner": models.Many2OneField{RelationModel: h.Partner(), ReadOnly: true},
		"CommercialPartner": models.Many2OneField{String: "Partner Company", RelationModel: h.Partner(),
			Help: "Commercial Entity", ReadOnly: true},
		"Company":    models.Many2OneField{RelationModel: h.Company(), ReadOnly: true},
		"User":       models.Many2OneField{String: "Salesperson", RelationModel: h.User(), ReadOnly: true},
		"PriceTotal": models.FloatField{String: "Total Without Tax", ReadOnly: true},
		"UserCurrencyPriceTotal": models.FloatField{String: "Total Without Tax",
			Compute: h.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency(),
			Digits:  nbutils.Digits{0, 0}},
		"PriceAverage": models.FloatField{String: "Average Price", ReadOnly: true, GroupOperator: "avg"},
		"UserCurrencyPriceAverage": models.FloatField{String: "Average Price",
			Compute: h.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency(),
			Digits:  nbutils.Digits{0, 0}},
		"CurrencyRate": models.FloatField{ReadOnly: true, GroupOperator: "avg"},
		"Nbr":          models.IntegerField{String: "# of Lines", ReadOnly: true},
		"Type": models.SelectionField{Selection: types.Selection{
			"out_invoice": "Customer Invoice",
			"in_invoice":  "Vendor Bill",
			"out_refund":  "Customer Refund",
			"in_refund":   "Vendor Refund",
		}, ReadOnly: true},
		"State": models.SelectionField{Selection: types.Selection{
			"draft":     "Draft",
			"proforma":  "Pro-forma",
			"proforma2": "Pro-forma",
			"open":      "Open",
			"paid":      "Done",
			"cancel":    "Cancelled",
		}, ReadOnly: true},
		"DateDue": models.DateField{String: "Due Date", ReadOnly: true},
		"Account": models.Many2OneField{RelationModel: h.AccountAccount(), ReadOnly: true,
			Filter: q.AccountAccount().Deprecated().Equals(false)},
		"AccountLine": models.Many2OneField{RelationModel: h.AccountAccount(), ReadOnly: true,
			Filter: q.AccountAccount().Deprecated().Equals(false)},
		"PartnerBank": models.Many2OneField{String: "Bank Account", RelationModel: h.BankAccount(),
			ReadOnly: true},
		"Residual": models.FloatField{String: "Total Residual", ReadOnly: true},
		"UserCurrencyResidual": models.FloatField{String: "Total Residual",
			Compute: h.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency(),
			Digits:  nbutils.Digits{0, 0}},
		"Country": models.Many2OneField{String: "Country of the Partner Company", RelationModel: h.Country(),
			ReadOnly: true},
		"Weight": models.FloatField{String: "Gross Weight", ReadOnly: true},
		"Volume": models.FloatField{ReadOnly: true},
	})
	h.AccountInvoiceReport().Methods().Select().DeclareMethod(
		`Select`,
		func(rs h.AccountInvoiceReportSet) string {
			/*def _select(self):
			  select_str = """
			      SELECT sub.id, sub.date, sub.product_id, sub.partner_id, sub.country_id, sub.account_analytic_id,
			          sub.payment_term_id, sub.uom_name, sub.currency_id, sub.journal_id,
			          sub.fiscal_position_id, sub.user_id, sub.company_id, sub.nbr, sub.type, sub.state,
			          sub.weight, sub.volume,
			          sub.categ_id, sub.date_due, sub.account_id, sub.account_line_id, sub.partner_bank_id,
			          sub.product_qty, sub.price_total as price_total, sub.price_average as price_average,
			          COALESCE(cr.rate, 1) as currency_rate, sub.residual as residual, sub.commercial_partner_id as commercial_partner_id
			  """
			  return select_str

			*/
			return ""
		})

	h.AccountInvoiceReport().Methods().SubSelect().DeclareMethod(
		`SubSelect`,
		func(rs h.AccountInvoiceReportSet) string {
			/*def _sub_select(self):
			  select_str = """
			          SELECT ail.id AS id,
			              ai.date_invoice AS date,
			              ail.product_id, ai.partner_id, ai.payment_term_id, ail.account_analytic_id,
			              u2.name AS uom_name,
			              ai.currency_id, ai.journal_id, ai.fiscal_position_id, ai.user_id, ai.company_id,
			              1 AS nbr,
			              ai.type, ai.state, pt.categ_id, ai.date_due, ai.account_id, ail.account_id AS account_line_id,
			              ai.partner_bank_id,
			              SUM ((invoice_type.sign * ail.quantity) / u.factor * u2.factor) AS product_qty,
			              SUM(ail.price_subtotal_signed) AS price_total,
			              SUM(ABS(ail.price_subtotal_signed)) / CASE
			                      WHEN SUM(ail.quantity / u.factor * u2.factor) <> 0::numeric
			                         THEN SUM(ail.quantity / u.factor * u2.factor)
			                         ELSE 1::numeric
			                      END AS price_average,
			              ai.residual_company_signed / (SELECT count(*) FROM account_invoice_line l where invoice_id = ai.id) *
			              count(*) * invoice_type.sign AS residual,
			              ai.commercial_partner_id as commercial_partner_id,
			              partner.country_id,
			              SUM(pr.weight * (invoice_type.sign*ail.quantity) / u.factor * u2.factor) AS weight,
			              SUM(pr.volume * (invoice_type.sign*ail.quantity) / u.factor * u2.factor) AS volume
			  """
			  return select_str

			*/
			return ""
		})

	h.AccountInvoiceReport().Methods().From().DeclareMethod(
		`From`,
		func(rs h.AccountInvoiceReportSet) string {
			/*def _from(self):
			  from_str = """
			          FROM account_invoice_line ail
			          JOIN account_invoice ai ON ai.id = ail.invoice_id
			          JOIN res_partner partner ON ai.commercial_partner_id = partner.id
			          LEFT JOIN product_product pr ON pr.id = ail.product_id
			          left JOIN product_template pt ON pt.id = pr.product_tmpl_id
			          LEFT JOIN product_uom u ON u.id = ail.uom_id
			          LEFT JOIN product_uom u2 ON u2.id = pt.uom_id
			          JOIN (
			              -- Temporary table to decide if the qty should be added or retrieved (Invoice vs Refund)
			              SELECT id,(CASE
			                   WHEN ai.type::text = ANY (ARRAY['out_refund'::character varying::text, 'in_invoice'::character varying::text])
			                      THEN -1
			                      ELSE 1
			                  END) AS sign
			              FROM account_invoice ai
			          ) AS invoice_type ON invoice_type.id = ai.id
			  """
			  return from_str

			*/
			return ""
		})

	h.AccountInvoiceReport().Methods().GroupByClause().DeclareMethod(
		`GroupBy`,
		func(rs h.AccountInvoiceReportSet) string {
			/*def _group_by(self):
			  group_by_str = """
			          GROUP BY ail.id, ail.product_id, ail.account_analytic_id, ai.date_invoice, ai.id,
			              ai.partner_id, ai.payment_term_id, u2.name, u2.id, ai.currency_id, ai.journal_id,
			              ai.fiscal_position_id, ai.user_id, ai.company_id, ai.type, invoice_type.sign, ai.state, pt.categ_id,
			              ai.date_due, ai.account_id, ail.account_id, ai.partner_bank_id, ai.residual_company_signed,
			              ai.amount_total_company_signed, ai.commercial_partner_id, partner.country_id
			  """
			  return group_by_str

			*/
			return ""
		})

	h.AccountInvoiceReport().Methods().Init().DeclareMethod(
		`Init`,
		func(rs h.AccountInvoiceReportSet) {
			//@api.model_cr
			/*def init(self):
			  # self._table = account_invoice_report
			  tools.drop_view_if_exists(self.env.cr, self._table)
			  self.env.cr.execute("""CREATE or REPLACE VIEW %s as (
			      WITH currency_rate AS (%s)
			      %s
			      FROM (
			          %s %s %s
			      ) AS sub
			      LEFT JOIN currency_rate cr ON
			          (cr.currency_id = sub.currency_id AND
			           cr.company_id = sub.company_id AND
			           cr.date_start <= COALESCE(sub.date, NOW()) AND
			           (cr.date_end IS NULL OR cr.date_end > COALESCE(sub.date, NOW())))
			  )""" % (
			              self._table, self.env['res.currency']._select_companies_rates(),
			              self._select(), self._sub_select(), self._from(), self._group_by()))
			*/
		})

}
