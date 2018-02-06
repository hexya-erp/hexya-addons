package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.Company().AddFields(map[string]models.FieldDefinition{
		"FiscalyearLastDay": models.IntegerField{Default: models.DefaultValue(31), Required: true},
		"FiscalyearLastMonth": models.SelectionField{Selection: types.Selection{
			"1":  "January",
			"2":  "February",
			"3":  "March",
			"4":  "April",
			"5":  "May",
			"6":  "June",
			"7":  "July",
			"8":  "August",
			"9":  "September",
			"10": "October",
			"11": "November",
			"12": "December",
		}, Default: models.DefaultValue("12"), Required: true},
		"PeriodLockDate": models.DateField{String: "Lock Date for Non-Advisers",
			Help: `Only users with the 'Adviser' role can edit accounts prior to and inclusive of this date.
Use it for period locking inside an open fiscal year, for example.`},
		"FiscalyearLockDate": models.DateField{String: "Lock Date",
			Help: `No users, including Advisers, can edit accounts prior to and inclusive of this date.
Use it for fiscal year locking for example.`},
		"TransferAccount": models.Many2OneField{String: "Inter-Banks Transfer Account",
			RelationModel: h.AccountAccount(), Filter: q.AccountAccount().Reconcile().Equals(true).
					And().Deprecated().Equals(false).
					And().UserTypeFilteredOn(
				q.AccountAccountType().HexyaExternalID().Equals("account.data_account_type_current_assets")),
			Help: "Intermediary account used when moving money from a liquidity account to another"},
		"ExpectsChartOfAccounts": models.BooleanField{String: "Expects a Chart of Accounts",
			Default: models.DefaultValue(true)},
		"ChartTemplate": models.Many2OneField{RelationModel: h.AccountChartTemplate(),
			Help: "The chart template for the company (if any)"},
		"BankAccountCodePrefix": models.CharField{String: "Prefix of the bank accounts"},
		"CashAccountCodePrefix": models.CharField{String: "Prefix of the cash accounts"},
		"AccountsCodeDigits":    models.IntegerField{String: "Number of digits in an account code"},
		"TaxCalculationRoundingMethod": models.SelectionField{String: "Tax Calculation Rounding Method",
			Selection: types.Selection{
				"round_per_line": "Round per Line",
				"round_globally": "Round Globally",
			}, Default: models.DefaultValue("round_per_line"),
			Help: `If you select 'Round per Line' : for each tax the tax amount will first be computed and
rounded for each PO/SO/invoice line and then these rounded amounts will be summed
leading to the total amount for that tax.
If you select 'Round Globally': for each tax the tax amount will be computed for
each PO/SO/invoice line then these amounts will be summed and eventually this
total tax amount will be rounded. If you sell with tax included you should
choose 'Round per line' because you certainly want the sum of your tax-included
line subtotals to be equal to the total amount with taxes.`},
		"CurrencyExchangeJournal": models.Many2OneField{String: "Exchange Gain or Loss Journal",
			RelationModel: h.AccountJournal(), Filter: q.AccountJournal().Type().Equals("general")},
		"IncomeCurrencyExchangeAccount": models.Many2OneField{String: "Gain Exchange Rate Account",
			RelationModel: h.AccountAccount(), Related: "CurrencyExchangeJournal.DefaultCreditAccount",
			Filter: q.AccountAccount().InternalType().Equals("other").
				And().Deprecated().Equals(false).
				And().Company().EqualsEval("id")},
		"ExpenseCurrencyExchangeAccount": models.Many2OneField{String: "Loss Exchange Rate Account",
			RelationModel: h.AccountAccount(), Related: "CurrencyExchangeJournal.DefaultDebitAccount",
			Filter: q.AccountAccount().InternalType().Equals("other").
				And().Deprecated().Equals(false).
				And().Company().EqualsEval("id")},
		"AngloSaxonAccounting": models.BooleanField{String: "Use anglo-saxon accounting"},
		"PropertyStockAccountInputCateg": models.Many2OneField{String: "Input Account for Stock Valuation",
			RelationModel: h.AccountAccount()},
		"PropertyStockAccountOutputCateg": models.Many2OneField{String: "Output Account for Stock Valuation",
			RelationModel: h.AccountAccount()},
		"PropertyStockValuationAccount": models.Many2OneField{String: "Account Template for Stock Valuation",
			RelationModel: h.AccountAccount()},
		"BankJournals": models.One2ManyField{RelationModel: h.AccountJournal(), ReverseFK: "Company",
			JSON: "bank_journal_ids", Filter: q.AccountJournal().Type().Equals("bank")},
		"OverdueMsg": models.TextField{String: "Overdue Payments Message", Translate: true,
			Default: models.DefaultValue(`Dear Sir/Madam,

Our records indicate that some payments on your account are still due. Please find details below.
If the amount has already been paid, please disregard this notice. Otherwise, please forward us the total amount stated below.
If you have any queries regarding your account, Please contact us.

Thank you in advance for your cooperation.
Best Regards,`)},
	})

	h.Company().Methods().ComputeFiscalyearDates().DeclareMethod(
		`ComputeFiscalyearDates`,
		func(rs h.CompanySet, date dates.Date) (dates.Date, dates.Date) {
			//@api.multi
			/*def compute_fiscalyear_dates(self, date):
			  """ Computes the start and end dates of the fiscalyear where the given 'date' belongs to
			      @param date: a datetime object
			      @returns: a dictionary with date_from and date_to
			  """
			  self = self[0]
			  last_month = self.fiscalyear_last_month
			  last_day = self.fiscalyear_last_day
			  if (date.month < last_month or (date.month == last_month and date.day <= last_day)):
			      date = date.replace(month=last_month, day=last_day)
			  else:
			      if last_month == 2 and last_day == 29 and (date.year + 1) % 4 != 0:
			          date = date.replace(month=last_month, day=28, year=date.year + 1)
			      else:
			          date = date.replace(month=last_month, day=last_day, year=date.year + 1)
			  date_to = date
			  date_from = date + timedelta(days=1)
			  if date_from.month == 2 and date_from.day == 29:
			      date_from = date_from.replace(day=28, year=date_from.year - 1)
			  else:
			      date_from = date_from.replace(year=date_from.year - 1)
			  return {'date_from': date_from, 'date_to': date_to}

			*/
			return dates.Today(), dates.Today()
		})

	h.Company().Methods().GetNewAccountCode().DeclareMethod(
		`GetNewAccountCode`,
		func(rs h.CompanySet, currentCode, oldPrefix, newPrefix string, digits int) string {
			/*def get_new_account_code(self, current_code, old_prefix, new_prefix, digits):
			  return new_prefix + current_code.replace(old_prefix, '', 1).lstrip('0').rjust(digits-len(new_prefix), '0')

			*/
			return ""
		})

	h.Company().Methods().ReflectCodePrefixChange().DeclareMethod(
		`ReflectCodePrefixChange`,
		func(rs h.CompanySet, oldCode, newCode string, digits int) {
			/*def reflect_code_prefix_change(self, old_code, new_code, digits):
			  accounts = self.env['account.account'].search([('code', 'like', old_code), ('internal_type', '=', 'liquidity'),
			      ('company_id', '=', self.id)], order='code asc')
			  for account in accounts:
			      if account.code.startswith(old_code):
			          account.write({'code': self.get_new_account_code(account.code, old_code, new_code, digits)})

			*/
		})

	h.Company().Methods().ReflectCodeDigitsChange().DeclareMethod(
		`ReflectCodeDigitsChange`,
		func(rs h.CompanySet, digits int) {
			/*def reflect_code_digits_change(self, digits):
			  accounts = self.env['account.account'].search([('company_id', '=', self.id)], order='code asc')
			  for account in accounts:
			      account.write({'code': account.code.rstrip('0').ljust(digits, '0')})

			*/
		})

	h.Company().Methods().ValidateFiscalyearLock().DeclareMethod(
		`ValidateFiscalyearLock`,
		func(rs h.CompanySet, values *h.CompanyData, fieldsToReset ...models.FieldNamer) {
			//@api.multi
			/*def _validate_fiscalyear_lock(self, values):
			  if values.get('fiscalyear_lock_date'):
			      nb_draft_entries = self.env['account.move'].search([
			          ('company_id', 'in', [c.id for c in self]),
			          ('state', '=', 'draft'),
			          ('date', '<=', values['fiscalyear_lock_date'])])
			      if nb_draft_entries:
			          raise ValidationError(_('There are still unposted entries in the period you want to lock. You should either post or delete them.'))

			*/
		})

	h.Company().Methods().Write().Extend("",
		func(rs h.CompanySet, data *h.CompanyData, fieldsToReset ...models.FieldNamer) bool {
			//@api.multi
			/*def write(self, values):
			  #restrict the closing of FY if there are still unposted entries
			  self._validate_fiscalyear_lock(values)

			  # Reflect the change on accounts
			  for company in self:
			      digits = values.get('accounts_code_digits') or company.accounts_code_digits
			      if values.get('bank_account_code_prefix') or values.get('accounts_code_digits'):
			          new_bank_code = values.get('bank_account_code_prefix') or company.bank_account_code_prefix
			          company.reflect_code_prefix_change(company.bank_account_code_prefix, new_bank_code, digits)
			      if values.get('cash_account_code_prefix') or values.get('accounts_code_digits'):
			          new_cash_code = values.get('cash_account_code_prefix') or company.cash_account_code_prefix
			          company.reflect_code_prefix_change(company.cash_account_code_prefix, new_cash_code, digits)
			      if values.get('accounts_code_digits'):
			          company.reflect_code_digits_change(digits)
			  return super(ResCompany, self).write(values)
			*/
			return rs.Super().Write(data, fieldsToReset...)
		})

}
