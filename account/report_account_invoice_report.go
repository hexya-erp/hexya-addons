// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountInvoiceReport().DeclareModel()
	pool.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency().DeclareMethod(
		`ComputeAmountsInUserCurrency`,
		func(rs pool.AccountInvoiceReportSet) {
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

			date = */
		})
	pool.AccountInvoiceReport().AddFields(map[string]models.FieldDefinition{
		"Date":                     models.DateField{String: "Date" /*[readonly True]*/},
		"Product":                  models.Many2OneField{String: "Product", RelationModel: pool.ProductProduct(), JSON: "product_id" /*['product.product']*/ /* readonly=true */},
		"ProductQty":               models.FloatField{String: "ProductQty" /*[string 'Product Quantity']*/ /*[ readonly True]*/},
		"UomName":                  models.CharField{String: "UomName" /*[string 'Reference Unit of Measure']*/ /*[ readonly True]*/},
		"PaymentTerm":              models.Many2OneField{String: "Payment Terms", RelationModel: pool.AccountPaymentTerm(), JSON: "payment_term_id" /*['account.payment.term']*/ /*[ oldname 'payment_term']*/ /* readonly=true */},
		"FiscalPosition":           models.Many2OneField{String: "Fiscal Position", RelationModel: pool.AccountFiscalPosition(), JSON: "fiscal_position_id" /*['account.fiscal.position']*/ /*[oldname 'fiscal_position']*/ /* readonly=true */},
		"Currency":                 models.Many2OneField{String: "Currency", RelationModel: pool.Currency(), JSON: "currency_id" /*['res.currency']*/ /* readonly=true */},
		"Categ":                    models.Many2OneField{String: "Product Category", RelationModel: pool.ProductCategory(), JSON: "categ_id" /*['product.category']*/ /* readonly=true */},
		"Journal":                  models.Many2OneField{String: "Journal", RelationModel: pool.AccountJournal(), JSON: "journal_id" /*['account.journal']*/ /* readonly=true */},
		"Partner":                  models.Many2OneField{String: "Partner", RelationModel: pool.Partner(), JSON: "partner_id" /*['res.partner']*/ /* readonly=true */},
		"CommercialPartner":        models.Many2OneField{String: "Partner Company", RelationModel: pool.Partner(), JSON: "commercial_partner_id" /*['res.partner']*/, Help: "Commercial Entity"},
		"Company":                  models.Many2OneField{String: "Company", RelationModel: pool.Company(), JSON: "company_id" /*['res.company']*/ /* readonly=true */},
		"User":                     models.Many2OneField{String: "Salesperson", RelationModel: pool.User(), JSON: "user_id" /*['res.users']*/ /* readonly=true */},
		"PriceTotal":               models.FloatField{String: "PriceTotal" /*[string 'Total Without Tax']*/ /*[ readonly True]*/},
		"UserCurrencyPriceTotal":   models.FloatField{String: "UserCurrencyPriceTotal" /*[string "Total Without Tax"]*/, Compute: pool.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency(), Digits: nbutils.Digits{0, 0}},
		"PriceAverage":             models.FloatField{String: "PriceAverage" /*[string 'Average Price']*/ /*[ readonly True]*/ /*[ group_operator "avg"]*/},
		"UserCurrencyPriceAverage": models.FloatField{String: "UserCurrencyPriceAverage" /*[string "Average Price"]*/, Compute: pool.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency(), Digits: nbutils.Digits{0, 0}},
		"CurrencyRate":             models.FloatField{String: "CurrencyRate" /*[string 'Currency Rate']*/ /*[ readonly True]*/ /*[ group_operator "avg"]*/},
		"Nbr":                      models.IntegerField{String: "Nbr" /*[string '# of Lines']*/ /*[ readonly True)  # TDE FIXME master: rename into nbr_lines type   fields.Selection([ ('out_invoice']*/ /*[ 'Customer Invoice']*/ /*[ ('in_invoice']*/ /*[ 'Vendor Bill']*/ /*[ ('out_refund']*/ /*[ 'Customer Refund']*/ /*[ ('in_refund']*/ /*[ 'Vendor Refund']*/ /*[ ]]*/ /*[ readonly True]*/},
		"Type": models.SelectionField{String: "Type", Selection: types.Selection{
			"out_invoice": "Customer Invoice",
			"in_invoice":  "Vendor Bill",
			"out_refund":  "Customer Refund",
			"in_refund":   "Vendor Refund",
			/*[ ('out_invoice', 'Customer Invoice'  ('in_invoice', 'Vendor Bill'  ('out_refund', 'Customer Refund'  ('in_refund', 'Vendor Refund'  ]*/} /*[]*/ /*[readonly True]*/},
		"State":                models.SelectionField{ /*state = fields.Selection([ ('draft', 'Draft'), ('proforma', 'Pro-forma'), ('proforma2', 'Pro-forma'), ('open', 'Open'), ('paid', 'Done'), ('cancel', 'Cancelled')*/ },
		"DateDue":              models.DateField{String: "DateDue" /*[string 'Due Date']*/ /*[ readonly True]*/},
		"Account":              models.Many2OneField{String: "Account", RelationModel: pool.AccountAccount(), JSON: "account_id" /*['account.account']*/ /* readonly=true */ /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/},
		"AccountLine":          models.Many2OneField{String: "Account Line", RelationModel: pool.AccountAccount(), JSON: "account_line_id" /*['account.account']*/ /* readonly=true */ /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/},
		"PartnerBank":          models.Many2OneField{String: "Bank Account", RelationModel: pool.BankAccount(), JSON: "partner_bank_id" /*['res.partner.bank']*/ /* readonly=true */},
		"Residual":             models.FloatField{String: "Residual" /*[string 'Total Residual']*/ /*[ readonly True]*/},
		"UserCurrencyResidual": models.FloatField{String: "UserCurrencyResidual" /*[string "Total Residual"]*/, Compute: pool.AccountInvoiceReport().Methods().ComputeAmountsInUserCurrency(), Digits: nbutils.Digits{0, 0}},
		"Country":              models.Many2OneField{String: "Country of the Partner Company", RelationModel: pool.Country(), JSON: "country_id" /*['res.country']*/},
		"Weight":               models.FloatField{String: "Weight" /*[string 'Gross Weight']*/ /*[ readonly True]*/},
		"Volume":               models.FloatField{String: "Volume" /*[string 'Volume']*/ /*[ readonly True]*/},
	})
	pool.AccountInvoiceReport().Methods().Select().DeclareMethod(
		`Select`,
		func(rs pool.AccountInvoiceReportSet) {
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
		})
	pool.AccountInvoiceReport().Methods().SubSelect().DeclareMethod(
		`SubSelect`,
		func(rs pool.AccountInvoiceReportSet) {
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
		})
	pool.AccountInvoiceReport().Methods().From().DeclareMethod(
		`From`,
		func(rs pool.AccountInvoiceReportSet) {
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
		})
	pool.AccountInvoiceReport().Methods().GroupBy().DeclareMethod(
		`GroupBy`,
		func(rs pool.AccountInvoiceReportSet) {
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
		})
	pool.AccountInvoiceReport().Methods().Init().DeclareMethod(
		`Init`,
		func(rs pool.AccountInvoiceReportSet) {
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
