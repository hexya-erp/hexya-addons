package account

import (
	"log"

	"github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/security"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.AccountFiscalPosition().DeclareModel()
	h.AccountFiscalPosition().SetDefaultOrder("Sequence")

	h.AccountFiscalPosition().AddFields(map[string]models.FieldDefinition{
		"Sequence": models.IntegerField{},
		"Name":     models.CharField{String: "Fiscal Position", Required: true},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "By unchecking the active field, you may hide a fiscal position without deleting it."},
		"Company": models.Many2OneField{RelationModel: h.Company()},
		"Accounts": models.One2ManyField{String: "Account Mapping",
			RelationModel: h.AccountFiscalPositionAccount(), ReverseFK: "Position", JSON: "account_ids",
			NoCopy: false},
		"Taxes": models.One2ManyField{String: "Tax Mapping", RelationModel: h.AccountFiscalPositionTax(),
			ReverseFK: "Position", JSON: "tax_ids", NoCopy: false},
		"Note": models.TextField{String: "Notes", Translate: true,
			Help: "Legal mentions that have to be printed on the invoices."},
		"AutoApply": models.BooleanField{String: "Detect Automatically",
			Help: "Apply automatically this fiscal position."},
		"VatRequired": models.BooleanField{String: "VAT required", Help: "Apply only if partner has a VAT number."},
		"Country": models.Many2OneField{String: "Country", RelationModel: h.Country(),
			OnChange: h.AccountFiscalPosition().Methods().OnchangeCountry(),
			Help:     "Apply only if delivery or invoicing country match."},
		"CountryGroup": models.Many2OneField{String: "Country Group", RelationModel: h.CountryGroup(),
			OnChange: h.AccountFiscalPosition().Methods().OnchangeCountryGroup(),
			Help:     "Apply only if delivery or invocing country match the group."},
		"States": models.Many2ManyField{String: "Federal States", RelationModel: h.CountryState(),
			JSON: "state_ids"},
		"ZipFrom": models.IntegerField{String: "Zip Range From", Default: models.DefaultValue(0),
			Constraint: h.AccountFiscalPosition().Methods().CheckZip()},
		"ZipTo": models.IntegerField{String: "Zip Range To", Default: models.DefaultValue(0),
			Constraint: h.AccountFiscalPosition().Methods().CheckZip()},
		"StatesCount": models.IntegerField{Compute: h.AccountFiscalPosition().Methods().ComputeStatesCount(),
			GoType: new(int)},
	})

	h.AccountFiscalPosition().Methods().ComputeStatesCount().DeclareMethod(
		`ComputeStatesCount returns the number of states of the partner's country'`,
		func(rs h.AccountFiscalPositionSet) *h.AccountFiscalPositionData {
			return &h.AccountFiscalPositionData{
				StatesCount: rs.Country().States().Len(),
			}
		})

	h.AccountFiscalPosition().Methods().CheckZip().DeclareMethod(
		`CheckZip fails if the zip range is the wrong way round`,
		func(rs h.AccountFiscalPositionSet) {
			if rs.ZipFrom() > rs.ZipTo() {
				log.Panic("Invalid 'Zip Range', please configure it properly.")
			}
		})

	h.AccountFiscalPosition().Methods().MapTax().DeclareMethod(
		`MapTax`,
		func(rs h.AccountFiscalPositionSet, taxes h.AccountTaxSet, product h.ProductProductSet,
			partner h.PartnerSet) h.AccountTaxSet {
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
			return h.AccountTax().NewSet(rs.Env())
		})

	h.AccountFiscalPosition().Methods().MapAccount().DeclareMethod(
		`MapAccount`,
		func(rs h.AccountFiscalPositionSet, account h.AccountAccountSet) h.AccountAccountSet {
			//@api.model
			/*def map_account(self, account):
			  for pos in self.account_ids:
			      if pos.account_src_id == account:
			          return pos.account_dest_id
			  return account

			*/
			return h.AccountAccount().NewSet(rs.Env())
		})

	h.AccountFiscalPosition().Methods().MapAccounts().DeclareMethod(
		`MapAccounts`,
		func(rs h.AccountFiscalPositionSet, accounts h.AccountAccountSet) h.AccountAccountSet {
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
			return h.AccountAccount().NewSet(rs.Env())
		})

	h.AccountFiscalPosition().Methods().OnchangeCountry().DeclareMethod(
		`OnchangeCountryId`,
		func(rs h.AccountFiscalPositionSet) (*h.AccountFiscalPositionData, []models.FieldNamer) {
			//@api.onchange('country_id')
			/*def _onchange_country_id(self):
			  if self.country_id:
			      self.zip_from = self.zip_to = self.country_group_id = False
			      self.state_ids = [(5,)]
			      self.states_count = len(self.country_id.state_ids)

			*/
			return &h.AccountFiscalPositionData{}, []models.FieldNamer{}
		})

	h.AccountFiscalPosition().Methods().OnchangeCountryGroup().DeclareMethod(
		`OnchangeCountryGroupId`,
		func(rs h.AccountFiscalPositionSet) (*h.AccountFiscalPositionData, []models.FieldNamer) {
			//@api.onchange('country_group_id')
			/*def _onchange_country_group_id(self):
			  if self.country_group_id:
			      self.zip_from = self.zip_to = self.country_id = False
			      self.state_ids = [(5,)]

			*/
			return &h.AccountFiscalPositionData{}, []models.FieldNamer{}
		})

	h.AccountFiscalPosition().Methods().GetFposByRegion().DeclareMethod(
		`GetFposByRegion`,
		func(rs h.AccountFiscalPositionSet, country h.CountrySet, state h.CountryStateSet, zipCode int64,
			vatRequired bool) h.AccountFiscalPositionSet {
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
			return h.AccountFiscalPosition().NewSet(rs.Env())
		})

	h.AccountFiscalPosition().Methods().GetFiscalPosition().DeclareMethod(
		`GetFiscalPosition`,
		func(rs h.AccountFiscalPositionSet, partner, delivery h.PartnerSet) h.AccountFiscalPositionSet {
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
			return h.AccountFiscalPosition().NewSet(rs.Env())
		})

	h.AccountFiscalPositionTax().DeclareModel()
	h.AccountFiscalPositionTax().AddFields(map[string]models.FieldDefinition{
		"Position": models.Many2OneField{String: "Fiscal Position", RelationModel: h.AccountFiscalPosition(),
			Required: true, OnDelete: models.Cascade},
		"TaxSrc": models.Many2OneField{String: "Tax on Product", RelationModel: h.AccountTax(),
			Required: true},
		"TaxDest": models.Many2OneField{String: "Tax to Apply", RelationModel: h.AccountTax()},
	})

	h.AccountFiscalPositionTax().AddSQLConstraint("tax_src_dest_uniq",
		"unique (position_id, tax_src_id, tax_dest_id)",
		"A tax fiscal position could be defined only once time on same taxes.")

	h.AccountFiscalPositionTax().Methods().NameGet().Extend("",
		func(rs h.AccountFiscalPositionTaxSet) string {
			return rs.Position().DisplayName()
		})

	h.AccountFiscalPositionAccount().DeclareModel()
	h.AccountFiscalPositionAccount().AddFields(map[string]models.FieldDefinition{
		"Position": models.Many2OneField{String: "Fiscal Position", RelationModel: h.AccountFiscalPosition(),
			Required: true, OnDelete: models.Cascade},
		"AccountSrc": models.Many2OneField{String: "Account on Product", RelationModel: h.AccountAccount(),
			Filter: q.AccountAccount().Deprecated().Equals(false), Required: true},
		"AccountDest": models.Many2OneField{String: "Account to Use Instead", RelationModel: h.AccountAccount(),
			Filter: q.AccountAccount().Deprecated().Equals(false), Required: true},
	})

	h.AccountFiscalPositionAccount().AddSQLConstraint("account_src_dest_uniq",
		"unique (position_id, account_src_id, account_dest_id)",
		"An account fiscal position could be defined only once time on same accounts.")

	h.AccountFiscalPositionAccount().Methods().NameGet().Extend("",
		func(rs h.AccountFiscalPositionAccountSet) string {
			return rs.Position().DisplayName()
		})

	h.Partner().AddFields(map[string]models.FieldDefinition{
		"Credit": models.FloatField{String: "Total Receivable",
			Compute: h.Partner().Methods().ComputeCreditDebit(), /*Search: "_credit_search"*/
			Help:    "Total amount this customer owes you."},
		"Debit": models.FloatField{String: "Total Payable",
			Compute: h.Partner().Methods().ComputeCreditDebit(), /* Search: "_debit_search"*/
			Help:    "Total amount you have to pay to this vendor."},
		"DebitLimit":    models.FloatField{String: "Payable Limit"},
		"TotalInvoiced": models.FloatField{Compute: h.Partner().Methods().ComputeTotalInvoiced()},
		"Currency": models.Many2OneField{String: "Currency", RelationModel: h.Currency(),
			Compute: h.Partner().Methods().ComputeCurrency(),
			Help:    "Utility field to express amount currency"},
		"ContractsCount": models.IntegerField{String: "Contracts",
			Compute: h.Partner().Methods().ComputeJournalItemCount(), GoType: new(int)},
		"JournalItemCount": models.IntegerField{String: "Journal Items",
			Compute: h.Partner().Methods().ComputeJournalItemCount(), GoType: new(int)},
		"IssuedTotal": models.FloatField{String: "Journal Items",
			Compute: h.Partner().Methods().ComputeIssuedTotal()},
		"PropertyAccountPayable": models.Many2OneField{String: "Account Payable",
			RelationModel: h.AccountAccount(), /*, CompanyDependent : true*/
			Filter:        q.AccountAccount().InternalType().Equals("payable").And().Deprecated().Equals(false),
			Help:          "This account will be used instead of the default one as the payable account for the current partner",
			//Required:      true,
		},
		"PropertyAccountReceivable": models.Many2OneField{String: "Account Receivable",
			RelationModel: h.AccountAccount(), /*, CompanyDependent : true*/
			Filter:        q.AccountAccount().InternalType().Equals("receivable").And().Deprecated().Equals(false),
			Help:          "This account will be used instead of the default one as the receivable account for the current partner",
			//Required:      true,
		},
		"PropertyAccountPosition": models.Many2OneField{String: "Fiscal Position",
			RelationModel: h.AccountFiscalPosition(), /*, CompanyDependent : true*/
			Help:          "The fiscal position will determine taxes and accounts used for the partner."},
		"PropertyPaymentTerm": models.Many2OneField{String: "Customer Payment Terms",
			RelationModel: h.AccountPaymentTerm(), /*, CompanyDependent : true*/
			Help:          "This payment term will be used instead of the default one for sale orders and customer invoices"},
		"PropertySupplierPaymentTerm": models.Many2OneField{String: "Vendor Payment Terms",
			RelationModel: h.AccountPaymentTerm(), /*, CompanyDependent : true*/
			Help:          "This payment term will be used instead of the default one for purchase orders and vendor bills"},
		"RefCompanies": models.One2ManyField{String: "Companies that refers to partner",
			RelationModel: h.Company(), ReverseFK: "Partner", JSON: "ref_company_ids"},
		"HasUnreconciledEntries": models.BooleanField{Compute: h.Partner().Methods().ComputeHasUnreconciledEntries(),
			Help: `The partner has at least one unreconciled debit and credit
since last time the invoices & payments matching was performed.`},
		"LastTimeEntriesChecked": models.DateTimeField{String: "Latest Invoices & Payments Matching Date",
			ReadOnly: true, NoCopy: true,
			Help: `Last time the invoices & payments matching was performed for this partner.
It is set either if there\'s not at least an unreconciled debit and an unreconciled
credit or if you click the "Done" button.`},
		"Invoices": models.One2ManyField{RelationModel: h.AccountInvoice(), ReverseFK: "Partner",
			JSON: "invoice_ids", ReadOnly: true, NoCopy: true},
		"Contracts": models.One2ManyField{RelationModel: h.AccountAnalyticAccount(), ReverseFK: "Partner",
			JSON: "contract_ids", ReadOnly: true},
		"BankAccountCount": models.IntegerField{String: "Bank",
			Compute: h.Partner().Methods().ComputeBankCount()},
		"Trust": models.SelectionField{String: "Degree of trust you have in this debtor", Selection: types.Selection{
			"good":   "Good Debtor",
			"normal": "Normal Debtor",
			"bad":    "Bad Debtor",
		}, Default: models.DefaultValue("normal") /*[ company_dependent True]*/},
		// TODO update hexya generate to master the Selection case below
		"InvoiceWarn": models.SelectionField{Selection: base.WarningMessage, String: "Invoice",
			/*Help: base.WarningHelp*/ Required: true, Default: models.DefaultValue("no-message")},
		"InvoiceWarnMsg": models.TextField{String: "Message for Invoice"},
	})

	h.Partner().Fields().TotalInvoiced().
		RevokeAccess(security.GroupEveryone, security.All).
		GrantAccess(GroupAccountInvoice, security.All)

	h.Partner().Methods().ComputeCreditDebit().DeclareMethod(
		`CreditDebitGet`,
		func(rs h.PartnerSet) *h.PartnerData {
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
			return &h.PartnerData{}
		})

	h.Partner().Methods().AssetDifferenceSearch().DeclareMethod(
		`AssetDifferenceSearch`,
		func(rs h.PartnerSet, accountType string, op operator.Operator, operand float64) q.PartnerCondition {
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
			return q.PartnerCondition{}
		})

	h.Partner().Methods().CreditSearch().DeclareMethod(
		`CreditSearch returns the condition to search on partners credits.`,
		func(rs h.PartnerSet, op operator.Operator, operand interface{}) q.PartnerCondition {
			return rs.AssetDifferenceSearch("receivable", op, operand.(float64))
		})

	h.Partner().Methods().DebitSearch().DeclareMethod(
		`DebitSearch returns the condition to search on partners debits.`,
		func(rs h.PartnerSet, op operator.Operator, operand interface{}) q.PartnerCondition {
			return rs.AssetDifferenceSearch("payable", op, operand.(float64))
		})

	h.Partner().Methods().ComputeTotalInvoiced().DeclareMethod(
		`InvoiceTotal`,
		func(rs h.PartnerSet) *h.PartnerData {
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
			return &h.PartnerData{}
		})

	h.Partner().Methods().ComputeJournalItemCount().DeclareMethod(
		`ComputeJournalItemCount`,
		func(rs h.PartnerSet) *h.PartnerData {
			//@api.multi
			/*def _journal_item_count(self):
			  for partner in self:
			      partner.journal_item_count = self.env['account.move.line'].search_count([('partner_id', '=', partner.id)])
			      partner.contracts_count = self.env['account.analytic.account'].search_count([('partner_id', '=', partner.id)])

			*/
			return &h.PartnerData{}
		})

	h.Partner().Methods().GetFollowupLinesDomain().DeclareMethod(
		`GetFollowupLinesDomain`,
		func(rs h.PartnerSet, date dates.Date, overdueOnly, onlyUnblocked bool) q.PartnerCondition {
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
			return q.PartnerCondition{}
		})

	h.Partner().Methods().ComputeIssuedTotal().DeclareMethod(
		`ComputeIssuedTotal`,
		func(rs h.PartnerSet) *h.PartnerData {
			//@api.multi
			/*def _compute_issued_total(self):
						  """ Returns the issued total as will be displayed on partner view """
						  today = fields.Date.context_today(self)
			      		  domain = self.get_followup_lines_domain(today, overdue_only=True)
			       		  for aml in self.env['account.move.line'].search(domain):
			        		    aml.partner_id.issued_total += aml.amount_residual
			*/
			return &h.PartnerData{}
		})

	h.Partner().Methods().ComputeHasUnreconciledEntries().DeclareMethod(
		`ComputeHasUnreconciledEntries`,
		func(rs h.PartnerSet) *h.PartnerData {
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
			return &h.PartnerData{}
		})

	h.Partner().Methods().MarkAsReconciled().DeclareMethod(
		`MarkAsReconciled`,
		func(rs h.PartnerSet) bool {
			//@api.multi
			/*def mark_as_reconciled(self):
			  self.env['account.partial.reconcile'].check_access_rights('write')
			  return self.sudo().write({'last_time_entries_checked': time.strftime(DEFAULT_SERVER_DATETIME_FORMAT)})

			*/
			return true
		})

	h.Partner().Methods().ComputeCurrency().DeclareMethod(
		`GetCompanyCurrency`,
		func(rs h.PartnerSet) *h.PartnerData {
			//@api.one
			/*def _get_company_currency(self):
			  if self.company_id:
			      self.currency_id = self.sudo().company_id.currency_id
			  else:
			      self.currency_id = self.env.user.company_id.currency_id
			*/
			return &h.PartnerData{}
		})

	h.Partner().Methods().ComputeBankCount().DeclareMethod(
		`ComputeBankCount`,
		func(rs h.PartnerSet) *h.PartnerData {
			//@api.multi
			/*def _compute_bank_count(self):
			  bank_data = self.env['res.partner.bank'].read_group([('partner_id', 'in', self.ids)], ['partner_id'], ['partner_id'])
			  mapped_data = dict([(bank['partner_id'][0], bank['partner_id_count']) for bank in bank_data])
			  for partner in self:
			      partner.bank_account_count = mapped_data.get(partner.id, 0)

			*/
			return &h.PartnerData{}
		})

	h.Partner().Methods().FindAccountingPartner().DeclareMethod(
		`FindAccountingPartner finds the partner for which the accounting entries will be created`,
		func(rs h.PartnerSet, partner h.PartnerSet) h.PartnerSet {
			return rs.CommercialPartner()
		})

	h.Partner().Methods().CommercialFields().Extend("",
		func(rs h.PartnerSet) []models.FieldNamer {
			//@api.model
			/*def _commercial_fields(self):
			  return super(ResPartner, self)._commercial_fields() + \
			      ['debit_limit', 'property_account_payable_id', 'property_account_receivable_id', 'property_account_position_id',
			       'property_payment_term_id', 'property_supplier_payment_term_id', 'last_time_entries_checked']

			*/
			return rs.Super().CommercialFields()
		})

	h.Partner().Methods().OpenPartnerHistory().DeclareMethod(
		`OpenPartnerHistory returns an action that display invoices/refunds made for the given partners.`,
		func(rs h.PartnerSet) *actions.Action {
			/*def open_partner_history(self):
			  '''
			  This function returns an action that display invoices/refunds made for the given partners.
			  '''
			  action = self.env.ref('account.action_invoice_refund_out_tree').read()[0]
			  action['domain'] = literal_eval(action['domain'])
			  action['domain'].append(('partner_id', 'child_of', self.ids))
			  return action
			*/
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

}
