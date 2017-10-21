package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountFiscalPosition().DeclareModel()
	pool.AccountFiscalPosition().AddFields(map[string]models.FieldDefinition{
		"Sequence":     models.IntegerField{String: "Sequence" /*[]*/},
		"Name":         models.CharField{String: "Name" /*[string 'Fiscal Position']*/, Required: true},
		"Active":       models.BooleanField{String: "Active", Default: models.DefaultValue(true), Help: "By unchecking the active field, you may hide a fiscal position without deleting it." /*[ you may hide a fiscal position without deleting it."]*/},
		"Company":      models.Many2OneField{String: "Company", RelationModel: pool.Company(), JSON: "company_id" /*['res.company']*/},
		"Accounts":     models.One2ManyField{String: "AccountIds", RelationModel: pool.AccountFiscalPositionAccount(), ReverseFK: "Position", JSON: "account_ids" /*['account.fiscal.position.account']*/ /*[ 'position_id']*/ /*[string 'Account Mapping']*/, NoCopy: false},
		"Taxs":         models.One2ManyField{String: "TaxIds", RelationModel: pool.AccountFiscalPositionTax(), ReverseFK: "Position", JSON: "tax_ids" /*['account.fiscal.position.tax']*/ /*[ 'position_id']*/ /*[string 'Tax Mapping']*/, NoCopy: false},
		"Note":         models.TextField{String: "Notes" /*['Notes']*/, Translate: true, Help: "Legal mentions that have to be printed on the invoices."},
		"AutoApply":    models.BooleanField{String: "AutoApply" /*[string 'Detect Automatically']*/, Help: "Apply automatically this fiscal position."},
		"VatRequired":  models.BooleanField{String: "VatRequired" /*[string 'VAT required']*/, Help: "Apply only if partner has a VAT number."},
		"Country":      models.Many2OneField{String: "Country", RelationModel: pool.Country(), JSON: "country_id" /*['res.country']*/, Help: "Apply only if delivery or invoicing country match."},
		"CountryGroup": models.Many2OneField{String: "Country Group", RelationModel: pool.CountryGroup(), JSON: "country_group_id" /*['res.country.group']*/, Help: "Apply only if delivery or invocing country match the group."},
		"States":       models.Many2ManyField{String: "Federal States", RelationModel: pool.CountryState(), JSON: "state_ids" /*['res.country.state']*/},
		"ZipFrom":      models.IntegerField{String: "ZipFrom" /*[string 'Zip Range From']*/, Default: models.DefaultValue(0)},
		"ZipTo":        models.IntegerField{String: "ZipTo" /*[string 'Zip Range To']*/, Default: models.DefaultValue(0)},
		"StatesCount":  models.IntegerField{String: "StatesCount", Compute: pool.AccountFiscalPosition().Methods().ComputeStatesCount()},
	})
	pool.AccountFiscalPosition().Methods().ComputeStatesCount().DeclareMethod(
		`ComputeStatesCount`,
		func(rs pool.AccountFiscalPositionSet) {
			//@api.one
			/*def _compute_states_count(self):
			  self.states_count = len(self.country_id.state_ids)

			*/
		})
	pool.AccountFiscalPosition().Methods().CheckZip().DeclareMethod(
		`CheckZip`,
		func(rs pool.AccountFiscalPositionSet) {
			//@api.constrains('zip_from','zip_to')
			/*def _check_zip(self):
			  if self.zip_from > self.zip_to:
			      raise ValidationError(_('Invalid "Zip Range", please configure it properly.'))
			  return True

			*/
		})
	pool.AccountFiscalPosition().Methods().MapTax().DeclareMethod(
		`MapTax`,
		func(rs pool.AccountFiscalPositionSet, args struct {
			Taxes   interface{}
			Product interface{}
			Partner interface{}
		}) {
			//@api.model#noqa
			/*def map_tax(self, taxes, product=None, partner=None):
			  result = self.env['account.tax'].browse()
			  for tax in taxes:
			      tax_count = 0
			      for t in self.tax_ids:
			          if t.tax_src_id == tax:
			              tax_count += 1
			              if t.tax_dest_id:
			                  result |= t.tax_dest_id
			      if not tax_count:
			          result |= tax
			  return result

			*/
		})
	pool.AccountFiscalPosition().Methods().MapAccount().DeclareMethod(
		`MapAccount`,
		func(rs pool.AccountFiscalPositionSet, args struct {
			Account interface{}
		}) {
			//@api.model
			/*def map_account(self, account):
			  for pos in self.account_ids:
			      if pos.account_src_id == account:
			          return pos.account_dest_id
			  return account

			*/
		})
	pool.AccountFiscalPosition().Methods().MapAccounts().DeclareMethod(
		`MapAccounts`,
		func(rs pool.AccountFiscalPositionSet, args struct {
			Accounts interface{}
		}) {
			//@api.model
			/*def map_accounts(self, accounts):
			  """ Receive a dictionary having accounts in values and try to replace those accounts accordingly to the fiscal position.
			  """
			  ref_dict = {}
			  for line in self.account_ids:
			      ref_dict[line.account_src_id] = line.account_dest_id
			  for key, acc in accounts.items():
			      if acc in ref_dict:
			          accounts[key] = ref_dict[acc]
			  return accounts

			*/
		})
	pool.AccountFiscalPosition().Methods().OnchangeCountryId().DeclareMethod(
		`OnchangeCountryId`,
		func(rs pool.AccountFiscalPositionSet) {
			//@api.onchange('country_id')
			/*def _onchange_country_id(self):
			  if self.country_id:
			      self.zip_from = self.zip_to = self.country_group_id = False
			      self.state_ids = [(5,)]
			      self.states_count = len(self.country_id.state_ids)

			*/
		})
	pool.AccountFiscalPosition().Methods().OnchangeCountryGroupId().DeclareMethod(
		`OnchangeCountryGroupId`,
		func(rs pool.AccountFiscalPositionSet) {
			//@api.onchange('country_group_id')
			/*def _onchange_country_group_id(self):
			  if self.country_group_id:
			      self.zip_from = self.zip_to = self.country_id = False
			      self.state_ids = [(5,)]

			*/
		})
	pool.AccountFiscalPosition().Methods().GetFposByRegion().DeclareMethod(
		`GetFposByRegion`,
		func(rs pool.AccountFiscalPositionSet, args struct {
			CountryId   interface{}
			StateId     interface{}
			Zipcode     interface{}
			VatRequired interface{}
		}) {
			//@api.model
			/*def _get_fpos_by_region(self, country_id=False, state_id=False, zipcode=False, vat_required=False):
			  if not country_id:
			      return False
			  base_domain = [('auto_apply', '=', True), ('vat_required', '=', vat_required)]
			  if self.env.context.get('force_company'):
			      base_domain.append(('company_id', '=', self.env.context.get('force_company')))
			  null_state_dom = state_domain = [('state_ids', '=', False)]
			  null_zip_dom = zip_domain = [('zip_from', '=', 0), ('zip_to', '=', 0)]
			  null_country_dom = [('country_id', '=', False), ('country_group_id', '=', False)]

			  if zipcode and zipcode.isdigit():
			      zipcode = int(zipcode)
			      zip_domain = [('zip_from', '<=', zipcode), ('zip_to', '>=', zipcode)]
			  else:
			      zipcode = 0

			  if state_id:
			      state_domain = [('state_ids', '=', state_id)]

			  domain_country = base_domain + [('country_id', '=', country_id)]
			  domain_group = base_domain + [('country_group_id.country_ids', '=', country_id)]

			  # Build domain to search records with exact matching criteria
			  fpos = self.search(domain_country + state_domain + zip_domain, limit=1)
			  # return records that fit the most the criteria, and fallback on less specific fiscal positions if any can be found
			  if not fpos and state_id:
			      fpos = self.search(domain_country + null_state_dom + zip_domain, limit=1)
			  if not fpos and zipcode:
			      fpos = self.search(domain_country + state_domain + null_zip_dom, limit=1)
			  if not fpos and state_id and zipcode:
			      fpos = self.search(domain_country + null_state_dom + null_zip_dom, limit=1)

			  # fallback: country group with no state/zip range
			  if not fpos:
			      fpos = self.search(domain_group + null_state_dom + null_zip_dom, limit=1)

			  if not fpos:
			      # Fallback on catchall (no country, no group)
			      fpos = self.search(base_domain + null_country_dom, limit=1)
			  return fpos or False

			*/
		})
	pool.AccountFiscalPosition().Methods().GetFiscalPosition().DeclareMethod(
		`GetFiscalPosition`,
		func(rs pool.AccountFiscalPositionSet, args struct {
			PartnerId  interface{}
			DeliveryId interface{}
		}) {
			//@api.model
			/*def get_fiscal_position(self, partner_id, delivery_id=None):
			  if not partner_id:
			      return False
			  # This can be easily overriden to apply more complex fiscal rules
			  PartnerObj = self.env['res.partner']
			  partner = PartnerObj.browse(partner_id)

			  # if no delivery use invoicing
			  if delivery_id:
			      delivery = PartnerObj.browse(delivery_id)
			  else:
			      delivery = partner

			  # partner manually set fiscal position always win
			  if delivery.property_account_position_id or partner.property_account_position_id:
			      return delivery.property_account_position_id.id or partner.property_account_position_id.id

			  # First search only matching VAT positions
			  vat_required = bool(partner.vat)
			  fp = self._get_fpos_by_region(delivery.country_id.id, delivery.state_id.id, delivery.zip, vat_required)

			  # Then if VAT required found no match, try positions that do not require it
			  if not fp and vat_required:
			      fp = self._get_fpos_by_region(delivery.country_id.id, delivery.state_id.id, delivery.zip, False)

			  return fp.id if fp else False


			*/
		})

	pool.AccountFiscalPositionTax().DeclareModel()
	pool.AccountFiscalPositionTax().AddFields(map[string]models.FieldDefinition{
		"Position": models.Many2OneField{String: "Fiscal Position", RelationModel: pool.AccountFiscalPosition(), JSON: "position_id" /*['account.fiscal.position']*/, Required: true, OnDelete: models.Cascade},
		"TaxSrc":   models.Many2OneField{String: "Tax on Product", RelationModel: pool.AccountTax(), JSON: "tax_src_id" /*['account.tax']*/, Required: true},
		"TaxDest":  models.Many2OneField{String: "Tax to Apply", RelationModel: pool.AccountTax(), JSON: "tax_dest_id" /*['account.tax']*/},
	})
	pool.AccountFiscalPositionTax().AddSQLConstraint( /* [('tax_src_dest_uniq'  ] */ )
	pool.AccountFiscalPositionTax().AddSQLConstraint( /* ['unique (position_id tax_src_id tax_dest_id)'  ] */ )
	pool.AccountFiscalPositionTax().AddSQLConstraint( /* ['A tax fiscal position could be defined only once time on same taxes.') ] */ )

	pool.AccountFiscalPositionAccount().DeclareModel()
	pool.AccountFiscalPositionAccount().AddFields(map[string]models.FieldDefinition{
		"Position":    models.Many2OneField{String: "Fiscal Position", RelationModel: pool.AccountFiscalPosition(), JSON: "position_id" /*['account.fiscal.position']*/, Required: true, OnDelete: models.Cascade},
		"AccountSrc":  models.Many2OneField{String: "Account on Product", RelationModel: pool.AccountAccount(), JSON: "account_src_id" /*['account.account']*/ /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/, Required: true},
		"AccountDest": models.Many2OneField{String: "Account to Use Instead", RelationModel: pool.AccountAccount(), JSON: "account_dest_id" /*['account.account']*/ /*, Filter: [('deprecated'*/ /*[ ' ']*/ /*[ False)]]*/, Required: true},
	})
	pool.AccountFiscalPositionAccount().AddSQLConstraint( /* [('account_src_dest_uniq'  ] */ )
	pool.AccountFiscalPositionAccount().AddSQLConstraint( /* ['unique (position_id account_src_id account_dest_id)'  ] */ )
	pool.AccountFiscalPositionAccount().AddSQLConstraint( /* ['An account fiscal position could be defined only once time on same accounts.') ] */ )

	pool.Partner().DeclareModel()
	pool.Partner().Methods().CreditDebitGet().DeclareMethod(
		`CreditDebitGet`,
		func(rs pool.PartnerSet) {
			//@api.multi
			/*def _credit_debit_get(self):
			  tables, where_clause, where_params = self.env['account.move.line']._query_get()
			  where_params = [tuple(self.ids)] + where_params
			  if where_clause:
			      where_clause = 'AND ' + where_clause
			  self._cr.execute("""SELECT account_move_line.partner_id, act.type, SUM(account_move_line.amount_residual)
			                FROM account_move_line
			                LEFT JOIN account_account a ON (account_move_line.account_id=a.id)
			                LEFT JOIN account_account_type act ON (a.user_type_id=act.id)
			                WHERE act.type IN ('receivable','payable')
			                AND account_move_line.partner_id IN %s
			                AND account_move_line.reconciled IS FALSE
			                """ + where_clause + """
			                GROUP BY account_move_line.partner_id, act.type
			                """, where_params)
			  for pid, type, val in self._cr.fetchall():
			      partner = self.browse(pid)
			      if type == 'receivable':
			          partner.credit = val
			      elif type == 'payable':
			          partner.debit = -val

			*/
		})
	pool.Partner().Methods().AssetDifferenceSearch().DeclareMethod(
		`AssetDifferenceSearch`,
		func(rs pool.PartnerSet, args struct {
			AccountType interface{}
			Operator    interface{}
			Operand     interface{}
		}) {
			//@api.multi
			/*def _asset_difference_search(self, account_type, operator, operand):
			  if operator not in ('<', '=', '>', '>=', '<='):
			      return []
			  if type(operand) not in (float, int):
			      return []
			  sign = 1
			  if account_type == 'payable':
			      sign = -1
			  res = self._cr.execute('''
			      SELECT partner.id
			      FROM res_partner partner
			      LEFT JOIN account_move_line aml ON aml.partner_id = partner.id
			      RIGHT JOIN account_account acc ON aml.account_id = acc.id
			      WHERE acc.internal_type = %s
			        AND NOT acc.deprecated
			      GROUP BY partner.id
			      HAVING %s * COALESCE(SUM(aml.amount_residual), 0) ''' + operator + ''' %s''', (account_type, sign, operand))
			  res = self._cr.fetchall()
			  if not res:
			      return [('id', '=', '0')]
			  return [('id', 'in', map(itemgetter(0), res))]

			*/
		})
	pool.Partner().Methods().CreditSearch().DeclareMethod(
		`CreditSearch`,
		func(rs pool.PartnerSet, args struct {
			Operator interface{}
			Operand  interface{}
		}) {
			//@api.model
			/*def _credit_search(self, operator, operand):
			  return self._asset_difference_search('receivable', operator, operand)

			*/
		})
	pool.Partner().Methods().DebitSearch().DeclareMethod(
		`DebitSearch`,
		func(rs pool.PartnerSet, args struct {
			Operator interface{}
			Operand  interface{}
		}) {
			//@api.model
			/*def _debit_search(self, operator, operand):
			  return self._asset_difference_search('payable', operator, operand)

			*/
		})
	pool.Partner().Methods().InvoiceTotal().DeclareMethod(
		`InvoiceTotal`,
		func(rs pool.PartnerSet) {
			//@api.multi
			/*def _invoice_total(self):
			  account_invoice_report = self.env['account.invoice.report']
			  if not self.ids:
			      self.total_invoiced = 0.0
			      return True

			  user_currency_id = self.env.user.company_id.currency_id.id
			  all_partners_and_children = {}
			  all_partner_ids = []
			  for partner in self:
			      # price_total is in the company currency
			      all_partners_and_children[partner] = self.search([('id', 'child_of', partner.id)]).ids
			      all_partner_ids += all_partners_and_children[partner]

			  # searching account.invoice.report via the orm is comparatively expensive
			  # (generates queries "id in []" forcing to build the full table).
			  # In simple cases where all invoices are in the same currency than the user's company
			  # access directly these elements

			  # generate where clause to include multicompany rules
			  where_query = account_invoice_report._where_calc([
			      ('partner_id', 'in', all_partner_ids), ('state', 'not in', ['draft', 'cancel']), ('company_id', '=', self.env.user.company_id.id),
			      ('type', 'in', ('out_invoice', 'out_refund'))
			  ])
			  account_invoice_report._apply_ir_rules(where_query, 'read')
			  from_clause, where_clause, where_clause_params = where_query.get_sql()

			  # price_total is in the company currency
			  query = """
			            SELECT SUM(price_total) as total, partner_id
			              FROM account_invoice_report account_invoice_report
			             WHERE %s
			             GROUP BY partner_id
			          """ % where_clause
			  self.env.cr.execute(query, where_clause_params)
			  price_totals = self.env.cr.dictfetchall()
			  for partner, child_ids in all_partners_and_children.items():
			      partner.total_invoiced = sum(price['total'] for price in price_totals if price['partner_id'] in child_ids)

			*/
		})
	pool.Partner().Methods().JournalItemCount().DeclareMethod(
		`JournalItemCount`,
		func(rs pool.PartnerSet) {
			//@api.multi
			/*def _journal_item_count(self):
			  for partner in self:
			      partner.journal_item_count = self.env['account.move.line'].search_count([('partner_id', '=', partner.id)])
			      partner.contracts_count = self.env['account.analytic.account'].search_count([('partner_id', '=', partner.id)])

			*/
		})
	pool.Partner().Methods().GetFollowupLinesDomain().DeclareMethod(
		`GetFollowupLinesDomain`,
		func(rs pool.PartnerSet, args struct {
			Date          interface{}
			OverdueOnly   interface{}
			OnlyUnblocked interface{}
		}) {
			/*def get_followup_lines_domain(self, date, overdue_only=False, only_unblocked=False):
			  domain = [('reconciled', '=', False), ('account_id.deprecated', '=', False), ('account_id.internal_type', '=', 'receivable'), '|', ('debit', '!=', 0), ('credit', '!=', 0), ('company_id', '=', self.env.user.company_id.id)]
			  if only_unblocked:
			      domain += [('blocked', '=', False)]
			  if self.ids:
			      if 'exclude_given_ids' in self._context:
			          domain += [('partner_id', 'not in', self.ids)]
			      else:
			          domain += [('partner_id', 'in', self.ids)]
			  #adding the overdue lines
			  overdue_domain = ['|', '&', ('date_maturity', '!=', False), ('date_maturity', '<', date), '&', ('date_maturity', '=', False), ('date', '<', date)]
			  if overdue_only:
			      domain += overdue_domain
			  return domain

			*/
		})
	pool.Partner().Methods().ComputeIssuedTotal().DeclareMethod(
		`ComputeIssuedTotal`,
		func(rs pool.PartnerSet) {
			//@api.multi
			/*def _compute_issued_total(self):
			  """ Returns the issued total as will be displayed on partner view """
			  today = */
		})
	pool.Partner().AddFields(map[string]models.FieldDefinition{})
	pool.Partner().Methods().ComputeHasUnreconciledEntries().DeclareMethod(
		`ComputeHasUnreconciledEntries`,
		func(rs pool.PartnerSet) {
			//@api.one
			/*def _compute_has_unreconciled_entries(self):
			  # Avoid useless work if has_unreconciled_entries is not relevant for this partner
			  if not self.active or not self.is_company and self.parent_id:
			      return
			  self.env.cr.execute(
			      """ SELECT 1 FROM(
			              SELECT
			                  p.last_time_entries_checked AS last_time_entries_checked,
			                  MAX(l.write_date) AS max_date
			              FROM
			                  account_move_line l
			                  RIGHT JOIN account_account a ON (a.id = l.account_id)
			                  RIGHT JOIN res_partner p ON (l.partner_id = p.id)
			              WHERE
			                  p.id = %s
			                  AND EXISTS (
			                      SELECT 1
			                      FROM account_move_line l
			                      WHERE l.account_id = a.id
			                      AND l.partner_id = p.id
			                      AND l.amount_residual > 0
			                  )
			                  AND EXISTS (
			                      SELECT 1
			                      FROM account_move_line l
			                      WHERE l.account_id = a.id
			                      AND l.partner_id = p.id
			                      AND l.amount_residual < 0
			                  )
			              GROUP BY p.last_time_entries_checked
			          ) as s
			          WHERE (last_time_entries_checked IS NULL OR max_date > last_time_entries_checked)
			      """, (self.id,))
			  self.has_unreconciled_entries = self.env.cr.rowcount == 1

			*/
		})
	pool.Partner().Methods().MarkAsReconciled().DeclareMethod(
		`MarkAsReconciled`,
		func(rs pool.PartnerSet) {
			//@api.multi
			/*def mark_as_reconciled(self):
			  self.env['account.partial.reconcile'].check_access_rights('write')
			  return self.sudo().write({'last_time_entries_checked': time.strftime(DEFAULT_SERVER_DATETIME_FORMAT)})

			*/
		})
	pool.Partner().Methods().GetCompanyCurrency().DeclareMethod(
		`GetCompanyCurrency`,
		func(rs pool.PartnerSet) {
			//@api.one
			/*def _get_company_currency(self):
			    if self.company_id:
			        self.currency_id = self.sudo().company_id.currency_id
			    else:
			        self.currency_id = self.env.user.company_id.currency_id

			credit = */
		})
	pool.Partner().AddFields(map[string]models.FieldDefinition{
		"Credit":                      models.FloatField{String: "Credit", Compute: pool.Partner().Methods().CreditDebitGet() /*, Search: "_credit_search"*/ /*[ string 'Total Receivable']*/, Help: "Total amount this customer owes you."},
		"Debit":                       models.FloatField{String: "Debit", Compute: pool.Partner().Methods().CreditDebitGet() /*, Search: "_debit_search"*/ /*[ string 'Total Payable']*/, Help: "Total amount you have to pay to this vendor."},
		"DebitLimit":                  models.FloatField{String: "Payable Limit')" /*['Payable Limit']*/},
		"TotalInvoiced":               models.FloatField{String: "TotalInvoiced", Compute: pool.Partner().Methods().InvoiceTotal() /*[ string "Total Invoiced"]*/ /*[ groups 'account.group_account_invoice']*/},
		"Currency":                    models.Many2OneField{String: "Currency", RelationModel: pool.Currency(), JSON: "currency_id" /*['res.currency']*/, Compute: pool.Partner().Methods().GetCompanyCurrency() /* readonly=true */, Help: "Utility field to express amount currency"},
		"ContractsCount":              models.IntegerField{String: "ContractsCount", Compute: pool.Partner().Methods().JournalItemCount() /*[ string "Contracts"]*/ /*[ type 'integer']*/},
		"JournalItemCount":            models.IntegerField{String: "JournalItemCount", Compute: pool.Partner().Methods().JournalItemCount() /*[ string "Journal Items"]*/ /*[ type "integer"]*/},
		"IssuedTotal":                 models.FloatField{String: "IssuedTotal", Compute: pool.Partner().Methods().ComputeIssuedTotal() /*[ string "Journal Items"]*/},
		"PropertyAccountPayable":      models.Many2OneField{String: "Account Payable", RelationModel: pool.AccountAccount(), JSON: "property_account_payable_id" /*['account.account']*/ /*, CompanyDependent : true*/ /*[ oldname "property_account_payable"]*/ /*, Filter: "[('internal_type'*/ /*[ ' ']*/ /*[ 'payable']*/ /*[ ('deprecated']*/ /*[ ' ']*/ /*[ False)]"]*/, Help: "This account will be used instead of the default one as the payable account for the current partner", Required: true},
		"PropertyAccountReceivable":   models.Many2OneField{String: "Account Receivable", RelationModel: pool.AccountAccount(), JSON: "property_account_receivable_id" /*['account.account']*/ /*, CompanyDependent : true*/ /*[ oldname "property_account_receivable"]*/ /*, Filter: "[('internal_type'*/ /*[ ' ']*/ /*[ 'receivable']*/ /*[ ('deprecated']*/ /*[ ' ']*/ /*[ False)]"]*/, Help: "This account will be used instead of the default one as the receivable account for the current partner", Required: true},
		"PropertyAccountPosition":     models.Many2OneField{String: "Fiscal Position", RelationModel: pool.AccountFiscalPosition(), JSON: "property_account_position_id" /*['account.fiscal.position']*/ /*, CompanyDependent : true*/, Help: "The fiscal position will determine taxes and accounts used for the partner." /*[ oldname "property_account_position"]*/},
		"PropertyPaymentTerm":         models.Many2OneField{String: "Customer Payment Terms", RelationModel: pool.AccountPaymentTerm(), JSON: "property_payment_term_id" /*['account.payment.term']*/ /*, CompanyDependent : true*/, Help: "This payment term will be used instead of the default one for sale orders and customer invoices" /*[ oldname "property_payment_term"]*/},
		"PropertySupplierPaymentTerm": models.Many2OneField{String: "Vendor Payment Terms", RelationModel: pool.AccountPaymentTerm(), JSON: "property_supplier_payment_term_id" /*['account.payment.term']*/ /*, CompanyDependent : true*/, Help: "This payment term will be used instead of the default one for purchase orders and vendor bills" /*[ oldname "property_supplier_payment_term"]*/},
		"RefCompanys":                 models.One2ManyField{String: "RefCompanyIds", RelationModel: pool.Company(), ReverseFK: "Partner", JSON: "ref_company_ids" /*['res.company']*/ /*[ 'partner_id']*/ /*[string 'Companies that refers to partner']*/ /*[ oldname "ref_companies"]*/},
		"HasUnreconciledEntries":      models.BooleanField{String: "HasUnreconciledEntries" /*[compute '_compute_has_unreconciled_entries']*/, Help: "The partner has at least one unreconciled debit and credit since last time the invoices & payments matching was performed."},
		"LastTimeEntriesChecked":      models.DateTimeField{String: "LastTimeEntriesChecked" /*[oldname 'last_reconciliation_date']*/ /*[ string 'Latest Invoices & Payments Matching Date']*/ /*[ readonly True]*/ /*[ copy False]*/, Help: `Last time the invoices & payments matching was performed for this partner. '#~#~# 'It is set either if there\'s not at least an unreconciled debit and an unreconciled credit '#~#~# 'or if you click the "Done" button.`},
		"Invoices":                    models.One2ManyField{String: "InvoiceIds", RelationModel: pool.AccountInvoice(), ReverseFK: "Partner", JSON: "invoice_ids" /*['account.invoice']*/ /*[ 'partner_id']*/ /*[string 'Invoices']*/ /* readonly */, NoCopy: true},
		"Contracts":                   models.One2ManyField{String: "ContractIds", RelationModel: pool.AccountAnalyticAccount(), ReverseFK: "Partner", JSON: "contract_ids" /*['account.analytic.account']*/ /*[ 'partner_id']*/ /*[string 'Contracts']*/ /* readonly */},
		"BankAccountCount":            models.IntegerField{String: "BankAccountCount", Compute: pool.Partner().Methods().ComputeBankCount() /*[ string "Bank"]*/},
		"Trust": models.SelectionField{String: "Degree of trust you have in this debtor", Selection: types.Selection{
			"good":   "Good Debtor",
			"normal": "Normal Debtor",
			"bad":    "Bad Debtor",
		}, /*[]*/ Default: models.DefaultValue("normal") /*[ company_dependent True]*/},
		"InvoiceWarn":    models.SelectionField{ /*invoice_warn = fields.Selection(WARNING_MESSAGE, 'Invoice', help=WARNING_HELP, required=True, default="no-message")*/ },
		"InvoiceWarnMsg": models.TextField{String: "Message for Invoice')" /*['Message for Invoice']*/},
	})
	pool.Partner().Methods().ComputeBankCount().DeclareMethod(
		`ComputeBankCount`,
		func(rs pool.PartnerSet) {
			//@api.multi
			/*def _compute_bank_count(self):
			  bank_data = self.env['res.partner.bank'].read_group([('partner_id', 'in', self.ids)], ['partner_id'], ['partner_id'])
			  mapped_data = dict([(bank['partner_id'][0], bank['partner_id_count']) for bank in bank_data])
			  for partner in self:
			      partner.bank_account_count = mapped_data.get(partner.id, 0)

			*/
		})
	pool.Partner().Methods().FindAccountingPartner().DeclareMethod(
		`FindAccountingPartner`,
		func(rs pool.PartnerSet, args struct {
			Partner interface{}
		}) {
			/*def _find_accounting_partner(self, partner):
			  ''' Find the partner for which the accounting entries will be created '''
			  return partner.commercial_partner_id

			*/
		})
	pool.Partner().Methods().CommercialFields().DeclareMethod(
		`CommercialFields`,
		func(rs pool.PartnerSet) {
			//@api.model
			/*def _commercial_fields(self):
			  return super(ResPartner, self)._commercial_fields() + \
			      ['debit_limit', 'property_account_payable_id', 'property_account_receivable_id', 'property_account_position_id',
			       'property_payment_term_id', 'property_supplier_payment_term_id', 'last_time_entries_checked']

			*/
		})
	pool.Partner().Methods().OpenPartnerHistory().DeclareMethod(
		`OpenPartnerHistory`,
		func(rs pool.PartnerSet) {
			/*def open_partner_history(self):
			  '''
			  This function returns an action that display invoices/refunds made for the given partners.
			  '''
			  action = self.env.ref('account.action_invoice_refund_out_tree').read()[0]
			  action['domain'] = literal_eval(action['domain'])
			  action['domain'].append(('partner_id', 'child_of', self.ids))
			  return action
			*/
		})

}
