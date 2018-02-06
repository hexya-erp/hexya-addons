// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

func init() {

	//h.AccountConfigSettings().DeclareTransientModel()
	//h.AccountConfigSettings().Methods().GetCurrency().DeclareMethod(
	//	`GetCurrency`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.depends('company_id')
	//		/*def _get_currency_id(self):
	//		  self.currency_id = self.company_id.currency_id
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().InverseCurrency().DeclareMethod(
	//	`InverseCurrency`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.one
	//		/*def _set_currency_id(self):
	//		    if self.currency_id != self.company_id.currency_id:
	//		        self.company_id.currency_id = self.currency_id
	//
	//		company_id = */
	//	})
	//h.AccountConfigSettings().AddFields(map[string]models.FieldDefinition{
	//	"Company": models.Many2OneField{String: "Company", RelationModel: h.Company(), JSON: "company_id" /*['res.company']*/, Required: true, Default: func(env models.Environment) interface{} {
	//		/*lambda self: self.env.user.company_id*/
	//		return 0
	//	}},
	//	"HasDefaultCompany":      models.BooleanField{String: "HasDefaultCompany" /*[readonly True]*/, Default: models.DefaultValue( /*lambda self: self._default_has_default_company()*/ 0)},
	//	"ExpectsChartOfAccounts": models.BooleanField{String: "ExpectsChartOfAccounts" /*[related 'company_id.expects_chart_of_accounts']*/ /*[ string 'This company has its own chart of accounts']*/, Help: "Check this box if this company is a legal entity."},
	//	"Currency":               models.Many2OneField{String: "Default company currency", RelationModel: h.Currency(), JSON: "currency_id" /*['res.currency']*/, Compute: h.AccountConfigSettings().Methods().GetCurrency(), Inverse: h.AccountConfigSettings().Methods().InverseCurrency(), Required: true, Help: "Main currency of the company."},
	//	"CompanyFooter":          models.TextField{String: "CompanyFooter" /*[related 'company_id.rml_footer']*/ /*[ string 'Bank accounts footer preview']*/ /*[ readonly True]*/, Help: "Bank accounts as printed in the footer of each printed document"},
	//	"HasChartOfAccounts":     models.BooleanField{String: "HasChartOfAccounts" /*[string 'Company has a chart of accounts']*/},
	//	"ChartTemplate":          models.Many2OneField{String: "Template", RelationModel: h.AccountChartTemplate(), JSON: "chart_template_id" /*['account.chart.template']*/ /*, Filter: "[('visible'*/ /*[' ']*/ /*[ True)]"]*/},
	//	"UseAngloSaxon":          models.BooleanField{String: "UseAngloSaxon" /*[string 'Use Anglo-Saxon Accounting *']*/ /*[ related 'company_id.anglo_saxon_accounting']*/},
	//	"CodeDigits":             models.IntegerField{String: "CodeDigits" /*[string '# of Digits *']*/ /*[ related 'company_id.accounts_code_digits']*/, Help: "No. of digits to use for account code"},
	//	"TaxCalculationRoundingMethod": models.SelectionField{String: "Tax calculation rounding method *", Selection: types.Selection{
	//		"round_per_line": "Round calculation of taxes per line",
	//		"round_globally": "Round globally calculation of taxes ",
	//		/*[ ('round_per_line', 'Round calculation of taxes per line'  ('round_globally', 'Round globally calculation of taxes '  ]*/}, /*[]*/ /*[related 'company_id.tax_calculation_rounding_method']*/ Help: "If you select 'Round per line' : for each tax" /*[ the tax amount will first be computed and rounded for each PO/SO/invoice line and then these rounded amounts will be summed]*/ /*[ leading to the total amount for that tax. If you select 'Round globally': for each tax]*/ /*[ the tax amount will be computed for each PO/SO/invoice line]*/ /*[ then these amounts will be summed and eventually this total tax amount will be rounded. If you sell with tax included]*/ /*[ you should choose 'Round per line' because you certainly want the sum of your tax-included line subtotals to be equal to the total amount with taxes."""]*/},
	//	"SaleTax":                 models.Many2OneField{String: "Default sale tax", RelationModel: h.AccountTaxTemplate(), JSON: "sale_tax_id" /*['account.tax.template']*/ /*[ oldname "sale_tax"]*/},
	//	"PurchaseTax":             models.Many2OneField{String: "Default purchase tax", RelationModel: h.AccountTaxTemplate(), JSON: "purchase_tax_id" /*['account.tax.template']*/ /*[ oldname "purchase_tax"]*/},
	//	"SaleTaxRate":             models.FloatField{String: "SaleTaxRate" /*[string 'Sales tax (%)']*/},
	//	"PurchaseTaxRate":         models.FloatField{String: "PurchaseTaxRate" /*[string 'Purchase tax (%)']*/},
	//	"BankAccountCodePrefix":   models.CharField{String: "BankAccountCodePrefix" /*[string 'Bank Accounts Prefix *']*/, Help: "Define the code prefix for the bank accounts', oldname='bank_account_code_char" /*[ oldname 'bank_account_code_char']*/},
	//	"CashAccountCodePrefix":   models.CharField{String: "CashAccountCodePrefix" /*[string 'Cash Accounts Prefix *']*/, Help: "Define the code prefix for the cash accounts"},
	//	"TemplateTransferAccount": models.Many2OneField{String: "TemplateTransferAccountId", RelationModel: h.AccountAccountTemplate(), JSON: "template_transfer_account_id" /*['account.account.template']*/, Help: "Intermediary account used when moving money from a liquidity account to another"},
	//	"TransferAccount":         models.Many2OneField{String: "TransferAccountId", RelationModel: h.AccountAccount(), JSON: "transfer_account_id" /*['account.account']*/, Related: "Company.TransferAccount" /*, Filter: lambda self: [('reconcile'*/ /*[ ' ']*/ /*[ True]*/ /*[ ('user_type_id.id']*/ /*[ ' ']*/ /*[ self.env.ref('account.data_account_type_current_assets').id)]]*/, Help: "Intermediary account used when moving money from a liquidity account to another"},
	//	"CompleteTaxSet":          models.BooleanField{String: "CompleteTaxSet" /*[string 'Complete set of taxes']*/, Help: "This boolean helps you to choose if you want to propose to the user to encode#~#~# the sales and purchase rates or use the usual m2o fields. This last choice assumes that#~#~# the set of tax defined for the chosen template is complete"},
	//	"FiscalyearLastDay":       models.IntegerField{String: "FiscalyearLastDay" /*[related 'company_id.fiscalyear_last_day']*/, Default: models.DefaultValue(31)},
	//	"FiscalyearLastMonth": models.SelectionField{String: "FiscalyearLastMonth", Selection: types.Selection{
	//		"1":  "January",
	//		"2":  "February",
	//		"3":  "March",
	//		"4":  "April",
	//		"5":  "May",
	//		"6":  "June",
	//		"7":  "July",
	//		"8":  "August",
	//		"9":  "September",
	//		"10": "October",
	//		"11": "November",
	//		"12": "December",
	//	}, /*[]*/ /*[related 'company_id.fiscalyear_last_month']*/ Default: models.DefaultValue("12")},
	//	"PeriodLockDate":                      models.DateField{String: "PeriodLockDate" /*[string "Lock Date for Non-Advisers"]*/ /*[ related 'company_id.period_lock_date']*/, Help: "Only users with the 'Adviser' role can edit accounts prior to and inclusive of this date. Use it for period locking inside an open fiscal year, for example." /*[ for example."]*/},
	//	"FiscalyearLockDate":                  models.DateField{String: "FiscalyearLockDate" /*[string "Lock Date"]*/ /*[ related 'company_id.fiscalyear_lock_date']*/, Help: "No users, including Advisers, can edit accounts prior to and inclusive of this date. Use it for fiscal year locking for example." /*[ including Advisers]*/ /*[ can edit accounts prior to and inclusive of this date. Use it for fiscal year locking for example."]*/},
	//	"ModuleAccountAccountant":             models.BooleanField{String: "ModuleAccountAccountant" /*[string 'Full accounting features: journals]*/ /*[ legal statements]*/ /*[ chart of accounts]*/ /*[ etc.']*/, Help: "If you do not check this box, you will be able to do invoicing & payments,#~#~# but not accounting (Journal Items, Chart of  Accounts, ...)" /*[ you will be able to do invoicing & payments]*/ /*[ but not accounting (Journal Items]*/ /*[ Chart of  Accounts]*/ /*[ ...)"""]*/},
	//	"ModuleAccountReports":                models.BooleanField{String: "Get dynamic accounting reports" /*["Get dynamic accounting reports"]*/},
	//	"GroupMultiCurrency":                  models.BooleanField{String: "GroupMultiCurrency" /*[string 'Allow multi currencies']*/ /*[ implied_group 'base.group_multi_currency']*/, Help: "Allows to work in a multi currency environment"},
	//	"GroupAnalyticAccounting":             models.BooleanField{String: "GroupAnalyticAccounting" /*[string 'Analytic accounting']*/ /*[ implied_group 'analytic.group_analytic_accounting']*/, Help: "Allows you to use the analytic accounting."},
	//	"GroupWarningAccount":                 models.SelectionField{ /*group_warning_account = fields.Selection([ (0, 'All the partners can be used in invoices'), (1, 'An informative or blocking warning can be set on a partner')*/ },
	//	"CurrencyExchangeJournal":             models.Many2OneField{String: "Rate Difference Journal", RelationModel: h.AccountJournal(), JSON: "currency_exchange_journal_id" /*['account.journal']*/, Related: "Company.CurrencyExchangeJournal" /*[]*/},
	//	"ModuleAccountAsset":                  models.BooleanField{String: "ModuleAccountAsset" /*[string 'Assets management']*/, Help: "Asset management: This allows you to manage the assets owned by a company or a person. '#~#~# 'It keeps track of the depreciation occurred on those assets, and creates account move for those depreciation lines.\n\n'#~#~# '-This installs the module account_asset." /*[ and creates account move for those depreciation lines.\n\n' '-This installs the module account_asset.']*/},
	//	"ModuleAccountDeferredRevenue":        models.BooleanField{String: "ModuleAccountDeferredRevenue" /*[string "Revenue Recognition"]*/, Help: "This allows you to manage the revenue recognition on selling products. '#~#~# 'It keeps track of the installments occurred on those revenue recognitions, '#~#~# 'and creates account moves for those installment lines\n'#~#~# '-This installs the module account_deferred_revenue." /*[ ' 'and creates account moves for those installment lines\n' '-This installs the module account_deferred_revenue.']*/},
	//	"ModuleAccountBudget":                 models.BooleanField{String: "ModuleAccountBudget" /*[string 'Budget management']*/, Help: "This allows accountants to manage analytic and crossovered budgets. '#~#~# 'Once the master budgets and the budgets are defined, '#~#~# 'the project managers can set the planned amount on each analytic account.\n'#~#~# '-This installs the module account_budget." /*[ ' 'the project managers can set the planned amount on each analytic account.\n' '-This installs the module account_budget.']*/},
	//	"ModuleAccountTaxCashBasis":           models.BooleanField{String: "ModuleAccountTaxCashBasis" /*[string "Allow Tax Cash Basis"]*/, Help: "Generate tax cash basis entrie when reconciliating entries"},
	//	"GroupProformaInvoices":               models.BooleanField{String: "GroupProformaInvoices" /*[string 'Allow pro-forma invoices']*/ /*[ implied_group 'account.group_proforma_invoices']*/, Help: "Allows you to put invoices in pro-forma state."},
	//	"ModuleAccountReportsFollowup":        models.BooleanField{String: "Enable payment followup management" /*["Enable payment followup management"]*/, Help: "This allows to automate letters for unpaid invoices, with multi-level recalls.\n'#~#~# '-This installs the module account_reports_followup." /*[ with multi-level recalls.\n' '-This installs the module account_reports_followup.']*/},
	//	"DefaultSaleTax":                      models.Many2OneField{String: "Default Sale Tax", RelationModel: h.AccountTax(), JSON: "default_sale_tax_id" /*['account.tax']*/, Help: "This sale tax will be assigned by default on new products." /*[ oldname "default_sale_tax"]*/},
	//	"DefaultPurchaseTax":                  models.Many2OneField{String: "Default Purchase Tax", RelationModel: h.AccountTax(), JSON: "default_purchase_tax_id" /*['account.tax']*/, Help: "This purchase tax will be assigned by default on new products." /*[ oldname "default_purchase_tax"]*/},
	//	"ModuleL10nUsCheckPrinting":           models.BooleanField{String: "Allow check printing and deposits" /*["Allow check printing and deposits"]*/},
	//	"ModuleAccountBatchDeposit":           models.BooleanField{String: "ModuleAccountBatchDeposit" /*[string 'Use batch deposit']*/, Help: "This allows you to group received checks before you deposit them to the bank.\n'#~#~# '-This installs the module account_batch_deposit."},
	//	"ModuleAccountSepa":                   models.BooleanField{String: "ModuleAccountSepa" /*[string 'Use SEPA payments']*/, Help: "If you check this box, you will be able to register your payment using SEPA.\n'#~#~# '-This installs the module account_sepa." /*[ you will be able to register your payment using SEPA.\n' '-This installs the module account_sepa.']*/},
	//	"ModuleAccountPlaid":                  models.BooleanField{String: "ModuleAccountPlaid" /*[string "Plaid Connector"]*/, Help: "Get your bank statements from you bank and import them through plaid.com.\n'#~#~# '-This installs the module account_plaid."},
	//	"ModuleAccountYodlee":                 models.BooleanField{String: "Bank Interface - Sync your bank feeds automatically" /*["Bank Interface - Sync your bank feeds automatically"]*/, Help: "Get your bank statements from your bank and import them through yodlee.com.\n'#~#~# '-This installs the module account_yodlee."},
	//	"ModuleAccountBankStatementImportQif": models.BooleanField{String: "Import .qif files" /*["Import .qif files"]*/, Help: "Get your bank statements from your bank and import them in Hexya in the .QIF format.\n'#~#~# '-This installs the module account_bank_statement_import_qif."},
	//	"ModuleAccountBankStatementImportOfx": models.BooleanField{String: "Import in .ofx format" /*["Import in .ofx format"]*/, Help: "Get your bank statements from your bank and import them in Hexya in the .OFX format.\n'#~#~# '-This installs the module account_bank_statement_import_ofx."},
	//	"ModuleAccountBankStatementImportCsv": models.BooleanField{String: "Import in .csv format" /*["Import in .csv format"]*/, Help: "Get your bank statements from your bank and import them in Hexya in the .CSV format.\n'#~#~# '-This installs the module account_bank_statement_import_csv."},
	//	"OverdueMsg":                          models.TextField{String: "OverdueMsg" /*[related 'company_id.overdue_msg']*/ /*[ string 'Overdue Payments Message *']*/},
	//})
	//h.AccountConfigSettings().Methods().DefaultHasDefaultCompany().DeclareMethod(
	//	`DefaultHasDefaultCompany`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.model
	//		/*def _default_has_default_company(self):
	//		  count = self.env['res.company'].search_count([])
	//		  return bool(count == 1)
	//
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OnchangeCompanyId().DeclareMethod(
	//	`OnchangeCompanyId`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.onchange('company_id')
	//		/*def onchange_company_id(self):
	//		  # update related fields
	//		  self.currency_id = False
	//		  if self.company_id:
	//		      company = self.company_id
	//		      self.chart_template_id = company.chart_template_id
	//		      self.has_chart_of_accounts = len(company.chart_template_id) > 0 or False
	//		      self.expects_chart_of_accounts = company.expects_chart_of_accounts
	//		      self.currency_id = company.currency_id
	//		      self.transfer_account_id = company.transfer_account_id
	//		      self.company_footer = company.rml_footer
	//		      self.tax_calculation_rounding_method = company.tax_calculation_rounding_method
	//		      self.bank_account_code_prefix = company.bank_account_code_prefix
	//		      self.cash_account_code_prefix = company.cash_account_code_prefix
	//		      self.code_digits = company.accounts_code_digits
	//
	//		      # update taxes
	//		      ir_values = self.env['ir.values']
	//		      taxes_id = ir_values.get_default('product.template', 'taxes_id', company_id = self.company_id.id)
	//		      supplier_taxes_id = ir_values.get_default('product.template', 'supplier_taxes_id', company_id = self.company_id.id)
	//		      self.default_sale_tax_id = isinstance(taxes_id, list) and len(taxes_id) > 0 and taxes_id[0] or taxes_id
	//		      self.default_purchase_tax_id = isinstance(supplier_taxes_id, list) and len(supplier_taxes_id) > 0 and supplier_taxes_id[0] or supplier_taxes_id
	//		  return {}
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OnchangeChartTemplateId().DeclareMethod(
	//	`OnchangeChartTemplateId`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.onchange('chart_template_id')
	//		/*def onchange_chart_template_id(self):
	//		  tax_templ_obj = self.env['account.tax.template']
	//		  self.complete_tax_set = self.sale_tax_id = self.purchase_tax_id = False
	//		  self.sale_tax_rate = self.purchase_tax_rate = 15
	//		  if self.chart_template_id and not self.has_chart_of_accounts:
	//		      # update complete_tax_set, sale_tax_id and purchase_tax_id
	//		      self.complete_tax_set = self.chart_template_id.complete_tax_set
	//		      if self.chart_template_id.complete_tax_set:
	//		          ir_values_obj = self.env['ir.values']
	//		          # default tax is given by the lowest sequence. For same sequence we will take the latest created as it will be the case for tax created while isntalling the generic chart of account
	//		          sale_tax = tax_templ_obj.search(
	//		              [('chart_template_id', 'parent_of', self.chart_template_id.id), ('type_tax_use', '=', 'sale')], limit=1,
	//		              order="sequence, id desc")
	//		          purchase_tax = tax_templ_obj.search(
	//		              [('chart_template_id', 'parent_of', self.chart_template_id.id), ('type_tax_use', '=', 'purchase')], limit=1,
	//		              order="sequence, id desc")
	//		          self.sale_tax_id = sale_tax
	//		          self.purchase_tax_id = purchase_tax
	//		      if self.chart_template_id.code_digits:
	//		          self.code_digits = self.chart_template_id.code_digits
	//		      if self.chart_template_id.transfer_account_id:
	//		          self.template_transfer_account_id = self.chart_template_id.transfer_account_id.id
	//		      if self.chart_template_id.bank_account_code_prefix:
	//		          self.bank_account_code_prefix = self.chart_template_id.bank_account_code_prefix
	//		      if self.chart_template_id.cash_account_code_prefix:
	//		          self.cash_account_code_prefix = self.chart_template_id.cash_account_code_prefix
	//		  return {}
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OnchangeTaxRate().DeclareMethod(
	//	`OnchangeTaxRate`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.onchange('sale_tax_rate')
	//		/*def onchange_tax_rate(self):
	//		  self.purchase_tax_rate = self.sale_tax_rate or False
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().DeclareGroupMultiCurrency().DeclareMethod(
	//	`DeclareGroupMultiCurrency`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_group_multi_currency(self):
	//		  ir_model = self.env['ir.model.data']
	//		  group_user = ir_model.get_object('base', 'group_user')
	//		  group_product = ir_model.get_object('product', 'group_sale_pricelist')
	//		  if self.group_multi_currency:
	//		      group_user.write({'implied_ids': [(4, group_product.id)]})
	//		  return True
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OpenBankAccounts().DeclareMethod(
	//	`OpenBankAccounts`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.multi
	//		/*def open_bank_accounts(self):
	//		  action_rec = self.env['ir.model.data'].xmlid_to_object('account.action_account_bank_journal_form')
	//		  return action_rec.read([])[0]
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().DeclareTransferAccount().DeclareMethod(
	//	`DeclareTransferAccount`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_transfer_account(self):
	//		  if self.transfer_account_id and self.transfer_account_id != self.company_id.transfer_account_id:
	//		      self.company_id.write({'transfer_account_id': self.transfer_account_id.id})
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().DeclareProductTaxes().DeclareMethod(
	//	`DeclareProductTaxes`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_product_taxes(self):
	//		  """ Set the product taxes if they have changed """
	//		  ir_values_obj = self.env['ir.values']
	//		  if self.default_sale_tax_id:
	//		      ir_values_obj.sudo().set_default('product.template', "taxes_id", [self.default_sale_tax_id.id], for_all_users=True, company_id=self.company_id.id)
	//		  if self.default_purchase_tax_id:
	//		      ir_values_obj.sudo().set_default('product.template', "supplier_taxes_id", [self.default_purchase_tax_id.id], for_all_users=True, company_id=self.company_id.id)
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().DeclareChartOfAccounts().DeclareMethod(
	//	`DeclareChartOfAccounts`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.multi
	//		/*def set_chart_of_accounts(self):
	//		  """ install a chart of accounts for the given company (if required) """
	//		  if self.chart_template_id and not self.has_chart_of_accounts and self.expects_chart_of_accounts:
	//		      if self.company_id.chart_template_id and self.chart_template_id != self.company_id.chart_template_id:
	//		          raise UserError(_('You can not change a company chart of account once it has been installed'))
	//		      wizard = self.env['wizard.multi.charts.accounts'].create({
	//		          'company_id': self.company_id.id,
	//		          'chart_template_id': self.chart_template_id.id,
	//		          'transfer_account_id': self.template_transfer_account_id.id,
	//		          'code_digits': self.code_digits or 6,
	//		          'sale_tax_id': self.sale_tax_id.id,
	//		          'purchase_tax_id': self.purchase_tax_id.id,
	//		          'sale_tax_rate': self.sale_tax_rate,
	//		          'purchase_tax_rate': self.purchase_tax_rate,
	//		          'complete_tax_set': self.complete_tax_set,
	//		          'currency_id': self.currency_id.id,
	//		          'bank_account_code_prefix': self.bank_account_code_prefix or self.chart_template_id.bank_account_code_prefix,
	//		          'cash_account_code_prefix': self.cash_account_code_prefix or self.chart_template_id.cash_account_code_prefix,
	//		      })
	//		      wizard.execute()
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OnchangeAnalyticAccounting().DeclareMethod(
	//	`OnchangeAnalyticAccounting`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.onchange('group_analytic_accounting')
	//		/*def onchange_analytic_accounting(self):
	//		  if self.group_analytic_accounting:
	//		      self.module_account_accountant = True
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OnchangeModuleAccountBudget().DeclareMethod(
	//	`OnchangeModuleAccountBudget`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.onchange('module_account_budget')
	//		/*def onchange_module_account_budget(self):
	//		  if self.module_account_budget:
	//		      self.group_analytic_accounting = True
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().OpenCompany().DeclareMethod(
	//	`OpenCompany`,
	//	func(rs h.AccountConfigSettingsSet) {
	//		//@api.multi
	//		/*def open_company(self):
	//		  return {
	//		      'type': 'ir.actions.act_window',
	//		      'name': 'My Company',
	//		      'view_type': 'form',
	//		      'view_mode': 'form',
	//		      'res_model': 'res.company',
	//		      'res_id': self.env.user.company_id.id,
	//		      'target': 'current',
	//		  }
	//
	//		*/
	//	})
	//h.AccountConfigSettings().Methods().Create().DeclareMethod(
	//	`Create`,
	//	func(rs h.AccountConfigSettingsSet, args struct {
	//		Vals interface{}
	//	}) {
	//		//@api.model
	//		/*def create(self, vals):
	//		  """
	//		  Avoid to rewrite the `accounts_code_digits` on the company if the value is the same. As all the values are
	//		  passed on the res.config creation, the related fields are rewriten on each res_config creation. Rewriting
	//		  the `account_code_digits` on the company will trigger a write on all the account.account to complete the
	//		  code the missing characters to complete the desired number of digit, leading to a sql_constraint.
	//		  """
	//		  if ('company_id' in vals and 'code_digits' in vals):
	//		      if self.env['res.company'].browse(vals.get('company_id')).accounts_code_digits == vals.get('code_digits'):
	//		          vals.pop('code_digits')
	//		  res = super(AccountConfigSettings, self).create(vals)
	//		  return res
	//		*/
	//	})

}
