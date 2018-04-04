// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/hexya-erp/hexya-addons/account/accounttypes"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.AccountAccountType().DeclareModel()
	h.AccountAccountType().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Account Type", Required: true, Translate: true},
		"IncludeInitialBalance": models.BooleanField{String: "Bring Accounts Balance Forward",
			Help: `Used in reports to know if we should consider journal items from the beginning of time instead of
from the fiscal year only. Account types that should be reset to zero at each new fiscal year
(like expenses, revenue..) should not have this option set.`},
		"Type": models.SelectionField{String: "Type", Selection: types.Selection{
			"other":      "Regular",
			"receivable": "Receivable",
			"payable":    "Payable",
			"liquidity":  "Liquidity",
		}, Required: true, Default: models.DefaultValue("other"),
			Help: `The 'Internal Type' is used for features available on different types of accounts:
- liquidity type is for cash or bank accounts
- payable/receivable is for vendor/customer accounts.`},
		"Note": models.TextField{String: "Description"},
	})

	h.AccountAccountTag().DeclareModel()
	h.AccountAccountTag().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Name", Required: true},
		"Applicability": models.SelectionField{String: "Applicability", Selection: types.Selection{
			"accounts": "Accounts",
			"taxes":    "Taxes",
		}, Required: true, Default: models.DefaultValue("accounts")},
		"Color": models.IntegerField{String: "Color Index"},
	})

	h.AccountAccount().DeclareModel()
	h.AccountAccount().SetDefaultOrder("Code")

	h.AccountAccount().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Required: true, Index: true},
		"Currency": models.Many2OneField{String: "Account Currency", RelationModel: h.Currency(),
			Help: "Forces all moves for this account to have this account currency."},
		"Code":       models.CharField{Size: 64, Required: true, Index: true},
		"Deprecated": models.BooleanField{Index: true, Default: models.DefaultValue(false)},
		"UserType": models.Many2OneField{String: "Type", RelationModel: h.AccountAccountType(),
			Required: true, Help: `Account Type is used for information purpose, to generate country-specific
legal reports, and set the rules to close a fiscal year and generate opening entries.`},
		"InternalType": models.SelectionField{Related: "UserType.Type", ReadOnly: true,
			Constraint: h.AccountAccount().Methods().CheckReconcile(),
			OnChange:   h.AccountAccount().Methods().OnchangeInternalType()},
		"LastTimeEntriesChecked": models.DateTimeField{String: "Latest Invoices & Payments Matching Date",
			ReadOnly: true, NoCopy: true,
			Help: `Last time the invoices & payments matching was performed on this account.
It is set either if there's not at least an unreconciled debit and an unreconciled credit
or if you click the "Done" button.`},
		"Reconcile": models.BooleanField{String: "Allow Reconciliation", Default: models.DefaultValue(false),
			Constraint: h.AccountAccount().Methods().CheckReconcile(),
			Help:       "Check this box if this account allows invoices & payments matching of journal items."},
		"Taxes": models.Many2ManyField{String: "Default Taxes", RelationModel: h.AccountTax(), JSON: "tax_ids"},
		"Note":  models.TextField{String: "Internal Notes"},
		"Company": models.Many2OneField{RelationModel: h.Company(), Required: true,
			Default: func(env models.Environment) interface{} {
				return h.Company().NewSet(env).CompanyDefaultGet()
			}},
		"Tags": models.Many2ManyField{RelationModel: h.AccountAccountTag(), JSON: "tag_ids",
			Help: "Optional tags you may want to assign for custom reporting"},
	})

	h.AccountAccount().AddSQLConstraint("code_company_uniq", "unique (code,company_id)",
		"The code of the account must be unique per company !")

	h.AccountAccount().Methods().CheckReconcile().DeclareMethod(
		`CheckReconcile`,
		func(rs h.AccountAccountSet) {
			//@api.constrains('internal_type','reconcile')
			/*def _check_reconcile(self):
			  for account in self:
			      if account.internal_type in ('receivable', 'payable') and account.reconcile == False:
			          raise ValidationError(_('You cannot have a receivable/payable account that is not reconciliable. (account code: %s)') % account.code)
			*/
		})

	h.AccountAccount().Methods().DefaultGet().Extend("",
		func(rs h.AccountAccountSet) models.FieldMap {
			//@api.model
			/*def default_get(self, default_fields):
			  """If we're creating a new account through a many2one, there are chances that we typed the account code
			  instead of its name. In that case, switch both fields values.
			  """
			  default_name = self._context.get('default_name')
			  default_code = self._context.get('default_code')
			  if default_name and not default_code:
			      try:
			          default_code = int(default_name)
			      except ValueError:
			          pass
			      if default_code:
			          default_name = False
			  contextual_self = self.with_context(default_name=default_name, default_code=default_code)
			  return super(AccountAccount, contextual_self).default_get(default_fields)

			*/
			return rs.Super().DefaultGet()
		})

	h.AccountAccount().Methods().SearchByName().Extend("",
		func(rs h.AccountAccountSet, name string, op operator.Operator, additionalCond q.AccountAccountCondition, limit int) h.AccountAccountSet {
			//@api.model
			/*def name_search(self, name, args=None, operator='ilike', limit=100):
			  args = args or []
			  domain = []
			  if name:
			      domain = ['|', ('code', '=ilike', name + '%'), ('name', operator, name)]
			      if operator in expression.NEGATIVE_TERM_OPERATORS:
			          domain = ['&', '!'] + domain[1:]
			  accounts = self.search(domain + args, limit=limit)
			  return accounts.name_get()

			*/
			return rs.Super().SearchByName(name, op, additionalCond, limit)
		})

	h.AccountAccount().Methods().OnchangeInternalType().DeclareMethod(
		`OnchangeInternalType`,
		func(rs h.AccountAccountSet) (*h.AccountAccountSet, []models.FieldNamer) {
			//@api.onchange('internal_type')
			/*def onchange_internal_type(self):
			  if self.internal_type in ('receivable', 'payable'):
			      self.reconcile = True

			*/
			return new(h.AccountAccountSet), []models.FieldNamer{}
		})

	h.AccountAccount().Methods().NameGet().Extend("",
		func(rs h.AccountAccountSet) string {
			//@api.depends('name','code')
			/*def name_get(self):
			  result = []
			  for account in self:
			      name = account.code + ' ' + account.name
			      result.append((account.id, name))
			  return result

			*/
			return rs.Super().NameGet()
		})

	h.AccountAccount().Methods().Copy().Extend("",
		func(rs h.AccountAccountSet, overrides *h.AccountAccountData, fieldsToReset ...models.FieldNamer) h.AccountAccountSet {
			//@api.returns('self',lambdavalue:value.id)
			/*def copy(self, default=None):
			  default = dict(default or {})
			  default.setdefault('code', _("%s (copy)") % (self.code or ''))
			  return super(AccountAccount, self).copy(default)

			*/
			return rs.Super().Copy(overrides, fieldsToReset...)
		})

	h.AccountAccount().Methods().Write().Extend("",
		func(rs h.AccountAccountSet, vals *h.AccountAccountData, fieldsToReset ...models.FieldNamer) bool {
			//@api.multi
			/*def write(self, vals):
			  # Dont allow changing the company_id when account_move_line already exist
			  if vals.get('company_id', False):
			      move_lines = self.env['account.move.line'].search([('account_id', 'in', self.ids)], limit=1)
			      for account in self:
			          if (account.company_id.id <> vals['company_id']) and move_lines:
			              raise UserError(_('You cannot change the owner company of an account that already contains journal items.'))
			  # If user change the reconcile flag, all aml should be recomputed for that account and this is very costly.
			  # So to prevent some bugs we add a constraint saying that you cannot change the reconcile field if there is any aml existing
			  # for that account.
			  if vals.get('reconcile'):
			      move_lines = self.env['account.move.line'].search([('account_id', 'in', self.ids)], limit=1)
			      if len(move_lines):
			          raise UserError(_('You cannot change the value of the reconciliation on this account as it already has some moves'))
			  return super(AccountAccount, self).write(vals)

			*/
			return rs.Super().Write(vals, fieldsToReset...)
		})

	h.AccountAccount().Methods().Unlink().Extend("",
		func(rs h.AccountAccountSet) int64 {
			//@api.multi
			/*def unlink(self):
			  if self.env['account.move.line'].search([('account_id', 'in', self.ids)], limit=1):
			      raise UserError(_('You cannot do that on an account that contains journal items.'))
			  #Checking whether the account is set as a property to any Partner or not
			  values = ['account.account,%s' % (account_id,) for account_id in self.ids]
			  partner_prop_acc = self.env['ir.property'].search([('value_reference', 'in', values)], limit=1)
			  if partner_prop_acc:
			      raise UserError(_('You cannot remove/deactivate an account which is set on a customer or vendor.'))
			  return super(AccountAccount, self).unlink()

			*/
			return rs.Super().Unlink()
		})

	h.AccountAccount().Methods().MarkAsReconciled().DeclareMethod(
		`MarkAsReconciled`,
		func(rs h.AccountAccountSet) bool {
			//@api.multi
			/*def mark_as_reconciled(self):
			  return self.write({'last_time_entries_checked': time.strftime(DEFAULT_SERVER_DATETIME_FORMAT)})

			*/
			return true
		})

	h.AccountAccount().Methods().ActionOpenReconcile().DeclareMethod(
		`ActionOpenReconcile`,
		func(rs h.AccountAccountSet) *actions.Action {
			//@api.multi
			/*def action_open_reconcile(self):
			  self.ensure_one()
			  # Open reconciliation view for this account
			  if self.internal_type == 'payable':
			      action_context = {'show_mode_selector': False, 'mode': 'suppliers'}
			  elif self.internal_type == 'receivable':
			      action_context = {'show_mode_selector': False, 'mode': 'customers'}
			  else:
			      action_context = {'show_mode_selector': False, 'account_ids': [self.id,]}
			  return {
			      'type': 'ir.actions.client',
			      'tag': 'manual_reconciliation_view',
			      'context': action_context,
			  }


			*/
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	h.AccountJournal().DeclareModel()
	h.AccountJournal().SetDefaultOrder("Sequence", "Type", "Code")

	h.AccountJournal().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Journal Name", Required: true},
		"Code": models.CharField{String: "Short Code", Size: 5, Required: true,
			Help: "The journal entries of this journal will be named using this prefix."},
		"Type": models.SelectionField{Selection: types.Selection{
			"sale":     "Sale",
			"purchase": "Purchase",
			"cash":     "Cash",
			"bank":     "Bank",
			"general":  "Miscellaneous",
		}, Required: true,
			Constraint: h.AccountJournal().Methods().CheckBankAccount(),
			Help: `Select 'Sale' for customer invoices journals.
Select 'Purchase' for vendor bills journals.
Select 'Cash' or 'Bank' for journals that are used in customer or vendor payments.
Select 'General' for miscellaneous operations journals.`},
		"TypeControls": models.Many2ManyField{String: "Account Types Allowed",
			RelationModel: h.AccountAccountType(), JSON: "type_control_ids"},
		"AccountControls": models.Many2ManyField{String: "Accounts Allowed", RelationModel: h.AccountAccount(),
			JSON: "account_control_ids", Filter: q.AccountAccount().Deprecated().Equals(false)},
		"DefaultCreditAccount": models.Many2OneField{RelationModel: h.AccountAccount(),
			Constraint: h.AccountJournal().Methods().CheckCurrency(),
			OnChange:   h.AccountJournal().Methods().OnchangeCreditAccountId(),
			Filter:     q.AccountAccount().Deprecated().Equals(false),
			Help:       "It acts as a default account for credit amount"},
		"DefaultDebitAccount": models.Many2OneField{RelationModel: h.AccountAccount(),
			Constraint: h.AccountJournal().Methods().CheckCurrency(),
			OnChange:   h.AccountJournal().Methods().OnchangeDebitAccountId(),
			Filter:     q.AccountAccount().Deprecated().Equals(false),
			Help:       "It acts as a default account for debit amount"},
		"UpdatePosted": models.BooleanField{String: "Allow Cancelling Entries",
			Help: `Check this box if you want to allow the cancellation the entries related to this journal or
of the invoice related to this journal`},
		"GroupInvoiceLines": models.BooleanField{
			Help: `If this box is checked, the system will try to group the accounting lines when generating
them from invoices.`},
		"EntrySequence": models.Many2OneField{RelationModel: h.Sequence(), JSON: "sequence_id",
			Help:     "This field contains the information related to the numbering of the journal entries of this journal.",
			Required: true, NoCopy: true},
		"RefundEntrySequence": models.Many2OneField{RelationModel: h.Sequence(), JSON: "refund_sequence_id",
			Help:   "This field contains the information related to the numbering of the refund entries of this journal.",
			NoCopy: true},
		"Sequence": models.IntegerField{
			Help:    "Used to order Journals in the dashboard view', default=10",
			Default: models.DefaultValue(10)},
		"Currency": models.Many2OneField{RelationModel: h.Currency(),
			Constraint: h.AccountJournal().Methods().CheckCurrency(),
			Help:       "The currency used to enter statement"},
		"Company": models.Many2OneField{RelationModel: h.Company(), Required: true, Index: true,
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company()
			}, Help: "Company related to this journal"},
		"RefundSequence": models.BooleanField{String: "Dedicated Refund Sequence",
			Help: `Check this box if you don't want to share the
same sequence for invoices and refunds made from this journal`,
			Default: models.DefaultValue(false)},
		"InboundPaymentMethods": models.Many2ManyField{String: "Debit Methods",
			RelationModel: h.AccountPaymentMethod(), JSON: "inbound_payment_method_ids",
			Filter: q.AccountPaymentMethod().PaymentType().Equals("inbound"),
			Default: func(env models.Environment) interface{} {
				return h.AccountPaymentMethod().Search(env,
					q.AccountPaymentMethod().HexyaExternalID().Equals("account_account_payment_method_manual_in"))
			},
			Help: `Means of payment for collecting money.
Hexya modules offer various payments handling facilities,
but you can always use the 'Manual' payment method in order
to manage payments outside of the software.`},
		"OutboundPaymentMethods": models.Many2ManyField{String: "Payment Methods",
			RelationModel: h.AccountPaymentMethod(), JSON: "outbound_payment_method_ids",
			Filter: q.AccountPaymentMethod().PaymentType().Equals("outbound"),
			Default: func(env models.Environment) interface{} {
				return h.AccountPaymentMethod().Search(env,
					q.AccountPaymentMethod().HexyaExternalID().Equals("account_account_payment_method_manual_out"))
			}, Help: `Means of payment for sending money.
Hexya modules offer various payments handling facilities
but you can always use the 'Manual' payment method in order
to manage payments outside of the software.`},
		"AtLeastOneInbound": models.BooleanField{Compute: h.AccountJournal().Methods().MethodsCompute(),
			Stored: true},
		"AtLeastOneOutbound": models.BooleanField{Compute: h.AccountJournal().Methods().MethodsCompute(),
			Stored: true},
		"ProfitAccount": models.Many2OneField{RelationModel: h.AccountAccount(),
			Filter: q.AccountAccount().Deprecated().Equals(false),
			Help:   "Used to register a profit when the ending balance of a cash register differs from what the system computes"},
		"LossAccount": models.Many2OneField{RelationModel: h.AccountAccount(),
			Filter: q.AccountAccount().Deprecated().Equals(false),
			Help:   "Used to register a loss when the ending balance of a cash register differs from what the system computes"},
		"BelongsToCompany": models.BooleanField{String: "Belong to the user's current company",
			Compute: h.AccountJournal().Methods().BelongToCompany() /*[ search "_search_company_journals"]*/},

		"BankAccount": models.Many2OneField{RelationModel: h.BankAccount(), OnDelete: models.Restrict,
			Constraint: h.AccountJournal().Methods().CheckBankAccount(),
			NoCopy:     true},
		"DisplayOnFooter": models.BooleanField{String: "Show in Invoices Footer",
			Help: "Display this bank account on the footer of printed documents like invoices and sales orders."},
		"BankStatementsSource": models.SelectionField{String: "Bank Feeds", Selection: types.Selection{
			"manual": "Record Manually",
		}},
		//"BankAccNumber": models.CharField{Related: "BankAccount.Name"},
		//"Bank":          models.Many2OneField{RelationModel: h.Bank(), Related: "BankAccount.Bank"},
	})

	h.AccountJournal().AddSQLConstraint("code_company_uniq", "unique (code, name, company_id)",
		"The code and name of the journal must be unique per company !'")

	h.AccountJournal().Methods().CheckCurrency().DeclareMethod(
		`CheckCurrency`,
		func(rs h.AccountJournalSet) {
			//@api.constrains('currency_id','default_credit_account_id','default_debit_account_id')
			/*def _check_currency(self):
			  if self.currency_id:
			      if self.default_credit_account_id and not self.default_credit_account_id.currency_id.id == self.currency_id.id:
			          raise ValidationError(_('Configuration error!\nThe currency of the journal should be the same than the default credit account.'))
			      if self.default_debit_account_id and not self.default_debit_account_id.currency_id.id == self.currency_id.id:
			          raise ValidationError(_('Configuration error!\nThe currency of the journal should be the same than the default debit account.'))

			*/
		})

	h.AccountJournal().Methods().CheckBankAccount().DeclareMethod(
		`CheckBankAccount`,
		func(rs h.AccountJournalSet) {
			//@api.constrains('type','bank_account_id')
			/*def _check_bank_account(self):
			  if self.type == 'bank' and self.bank_account_id:
			      if self.bank_account_id.company_id != self.company_id:
			          raise ValidationError(_('The bank account of a bank journal must belong to the same company (%s).') % self.company_id.name)
			      # A bank account can belong to a customer/supplier, in which case their partner_id is the customer/supplier.
			      # Or they are part of a bank journal and their partner_id must be the company's partner_id.
			      if self.bank_account_id.partner_id != self.company_id.partner_id:
			          raise ValidationError(_('The holder of a journal\'s bank account must be the company (%s).') % self.company_id.name)

			*/
		})

	h.AccountJournal().Methods().OnchangeDebitAccountId().DeclareMethod(
		`OnchangeDebitAccountId`,
		func(rs h.AccountJournalSet) (*h.AccountJournalData, []models.FieldNamer) {
			//@api.onchange('default_debit_account_id')
			/*def onchange_debit_account_id(self):
			  if not self.default_credit_account_id:
			      self.default_credit_account_id = self.default_debit_account_id

			*/
			return new(h.AccountJournalData), []models.FieldNamer{}
		})

	h.AccountJournal().Methods().OnchangeCreditAccountId().DeclareMethod(
		`OnchangeCreditAccountId`,
		func(rs h.AccountJournalSet) (*h.AccountJournalData, []models.FieldNamer) {
			//@api.onchange('default_credit_account_id')
			/*def onchange_credit_account_id(self):
			  if not self.default_debit_account_id:
			      self.default_debit_account_id = self.default_credit_account_id

			*/
			return new(h.AccountJournalData), []models.FieldNamer{}
		})

	h.AccountJournal().Methods().Unlink().Extend("",
		func(rs h.AccountJournalSet) int64 {
			//@api.multi
			/*def unlink(self):
			bank_accounts = self.env['res.partner.bank'].browse()
			for bank_account in self.mapped('bank_account_id'):
				accounts = self.search([('bank_account_id', '=', bank_account.id)])
				if accounts <= self:
					bank_accounts += bank_account
			ret = super(AccountJournal, self).unlink()
			bank_accounts.unlink()
			return ret
			*/
			return rs.Super().Unlink()
		})

	h.AccountJournal().Methods().Copy().Extend("",
		func(rs h.AccountJournalSet, overrides *h.AccountJournalData, fieldsToReset ...models.FieldNamer) h.AccountJournalSet {
			//@api.returns('self',lambdavalue:value.id)
			/*def copy(self, default=None):
			default = dict(default or {})
			default.update(
				code=_("%s (copy)") % (self.code or ''),
				name=_("%s (copy)") % (self.name or ''))
			return super(AccountJournal, self).copy(default)
			*/
			return rs.Super().Copy(overrides, fieldsToReset...)
		})

	h.AccountJournal().Methods().Write().Extend("",
		func(rs h.AccountJournalSet, vals *h.AccountJournalData, fieldsToReset ...models.FieldNamer) bool {
			//@api.multi
			/*def write(self, vals):
			for journal in self:
				if ('company_id' in vals and journal.company_id.id != vals['company_id']):
					if self.env['account.move'].search([('journal_id', 'in', self.ids)], limit=1):
						raise UserError(_('This journal already contains items, therefore you cannot modify its company.'))
				if ('code' in vals and journal.code != vals['code']):
					if self.env['account.move'].search([('journal_id', 'in', self.ids)], limit=1):
						raise UserError(_('This journal already contains items, therefore you cannot modify its short name.'))
					new_prefix = self._get_sequence_prefix(vals['code'], refund=False)
					journal.sequence_id.write({'prefix': new_prefix})
					if journal.refund_sequence_id:
						new_prefix = self._get_sequence_prefix(vals['code'], refund=True)
						journal.refund_sequence_id.write({'prefix': new_prefix})
				if 'currency_id' in vals:
					if not 'default_debit_account_id' in vals and self.default_debit_account_id:
						self.default_debit_account_id.currency_id = vals['currency_id']
					if not 'default_credit_account_id' in vals and self.default_credit_account_id:
						self.default_credit_account_id.currency_id = vals['currency_id']
				if 'bank_acc_number' in vals and not vals.get('bank_acc_number') and journal.bank_account_id:
					raise UserError(_('You cannot empty the account number once set.\nIf you would like to delete the account number, you can do it from the Bank Accounts list.'))
			result = super(AccountJournal, self).write(vals)

			# Create the bank_account_id if necessary
			if 'bank_acc_number' in vals:
				for journal in self.filtered(lambda r: r.type == 'bank' and not r.bank_account_id):
					journal.set_bank_account(vals.get('bank_acc_number'), vals.get('bank_id'))
			# create the relevant refund sequence
			if vals.get('refund_sequence'):
				for journal in self.filtered(lambda j: j.type in ('sale', 'purchase') and not j.refund_sequence_id):
					journal_vals = {
						'name': journal.name,
						'company_id': journal.company_id.id,
						'code': journal.code
					}
					journal.refund_sequence_id = self.sudo()._create_sequence(journal_vals, refund=True).id

			return result
			*/
			return rs.Super().Write(vals, fieldsToReset...)
		})

	h.AccountJournal().Methods().GetSequencePrefix().DeclareMethod(
		`GetSequencePrefix returns the prefix of the sequence for the given code.`,
		func(rs h.AccountJournalSet, code string, refund bool) string {
			prefix := strings.ToUpper(code)
			if refund {
				prefix = "R" + prefix
			}
			return prefix + "/%(range_year)%/"
		})

	h.AccountJournal().Methods().CreateSequence().DeclareMethod(
		`CreateSequence creates new no_gap entry sequence for every new Journal`,
		func(rs h.AccountJournalSet, vals *h.AccountJournalData, refund bool) h.SequenceSet {
			prefix := rs.GetSequencePrefix(vals.Code, refund)
			name := vals.Name
			if refund {
				name = rs.T("%s: Refund", name)
			}
			seq := h.SequenceData{
				Name:            name,
				Implementation:  "no_gap",
				Prefix:          prefix,
				Padding:         4,
				NumberIncrement: 1,
				UseDateRange:    true,
				Company:         vals.Company,
			}
			return h.Sequence().Create(rs.Env(), &seq)
		})

	h.AccountJournal().Methods().PrepareLiquidityAccount().DeclareMethod(
		`PrepareLiquidityAccount prepares the value to use for the creation of the default debit and credit accounts of a
			  liquidity journal (created through the wizard of generating COA from templates for example).`,
		func(rs h.AccountJournalSet, name string, company h.CompanySet, currency h.CurrencySet, accType string) *h.AccountAccountData {
			// Seek the next available number for the account code
			codeDigits := company.AccountsCodeDigits()
			accountCodePrefix := company.BankAccountCodePrefix()
			if accType == "cash" {
				if company.CashAccountCodePrefix() != "" {
					accountCodePrefix = company.CashAccountCodePrefix()
				}
			}
			var (
				flag    bool
				newCode string
			)
			for num := 1; num < 100; num++ {
				newCode = strings.Replace(
					fmt.Sprintf("%-[1]*[2]d%d", codeDigits-1, accountCodePrefix, num), " ", "0", -1)
				rec := h.AccountAccount().Search(rs.Env(),
					q.AccountAccount().Code().Equals(newCode).And().Company().Equals(company)).Limit(1)
				if rec.IsEmpty() {
					flag = true
					break
				}
			}
			if !flag {
				panic(rs.T("Cannot generate an unused account code."))
			}
			liquidityType := h.AccountAccountType().Search(rs.Env(),
				q.AccountAccountType().HexyaExternalID().Equals("account.data_account_type_liquidity"))

			return &h.AccountAccountData{
				Name:     name,
				Currency: currency,
				Code:     newCode,
				UserType: liquidityType,
				Company:  company,
			}
		})

	h.AccountJournal().Methods().Create().Extend("",
		func(rs h.AccountJournalSet, vals *h.AccountJournalData) h.AccountJournalSet {
			company := vals.Company
			if company.IsEmpty() {
				company = h.User().NewSet(rs.Env()).CurrentUser().Company()
			}
			if vals.Type == "bank" || vals.Type == "cash" {
				// # For convenience, the name can be inferred from account number
				// if not vals.get('name') and 'bank_acc_number' in vals:
				//    vals['name'] = vals['bank_acc_number']
				if vals.Code == "" {
					journalCodeBase := "BNK"
					if vals.Type == "cash" {
						journalCodeBase = "CSH"
					}
					journals := h.AccountJournal().Search(rs.Env(),
						q.AccountJournal().Code().Like(journalCodeBase+"%").And().Company().Equals(company))
					journalCodes := make(map[string]bool)
					for _, j := range journals.Records() {
						journalCodes[j.Code()] = true
					}
					for num := 1; num < 100; num++ {
						// journal_code has a maximal size of 5, hence we can enforce the boundary num < 100
						jCode := journalCodeBase + strconv.Itoa(num)
						if _, exists := journalCodes[jCode]; !exists {
							vals.Code = jCode
							break
						}
					}
					if vals.Code == "" {
						panic(rs.T("Cannot generate an unused journal code. Please fill the 'Shortcode' field."))
					}
				}
				// Create a default debit/credit account if not given
				defaultAccount := vals.DefaultDebitAccount
				if defaultAccount.IsEmpty() {
					defaultAccount = vals.DefaultCreditAccount
				}
				if defaultAccount.IsEmpty() {
					accountVals := rs.PrepareLiquidityAccount(vals.Name, company, vals.Currency, vals.Type)
					defaultAccount = h.AccountAccount().Create(rs.Env(), accountVals)
					vals.DefaultDebitAccount = defaultAccount
					vals.DefaultCreditAccount = defaultAccount
				}

			}
			// We just need to create the relevant sequences according to the chosen options
			if vals.EntrySequence.IsEmpty() {
				vals.EntrySequence = rs.Sudo().CreateSequence(vals, false)
			}
			if (vals.Type == "sale" || vals.Type == "purchase") && vals.RefundSequence && vals.RefundEntrySequence.IsEmpty() {
				vals.RefundEntrySequence = rs.Sudo().CreateSequence(vals, true)
			}
			journal := rs.Super().Create(vals)

			/*

			  # Create the bank_account_id if necessary
			  if journal.type == 'bank' and not journal.bank_account_id and vals.get('bank_acc_number'):
			      journal.set_bank_account(vals.get('bank_acc_number'), vals.get('bank_id'))

			  return journal

			*/
			return journal
		})

	h.AccountJournal().Methods().DefineBankAccount().DeclareMethod(
		`DefineBankAccount`,
		func(rs h.AccountJournalSet, accNumber string, bank h.BankSet) {
			/*def set_bank_account(self, acc_number, bank_id=None):
			  """ Create a res.partner.bank and set it as value of the  field bank_account_id """
			  self.ensure_one()
			  self.bank_account_id = self.env['res.partner.bank'].create({
			      'acc_number': acc_number,
			      'bank_id': bank_id,
			      'company_id': self.company_id.id,
			      'currency_id': self.currency_id.id,
			      'partner_id': self.company_id.partner_id.id,
			  }).id

			*/
		})

	h.AccountJournal().Methods().NameGet().Extend("",
		func(rs h.AccountJournalSet) string {
			//@api.depends('name','currency_id','company_id','company_id.currency_id')
			/*def name_get(self):
			res = []
			for journal in self:
				currency = journal.currency_id or journal.company_id.currency_id
				name = "%s (%s)" % (journal.name, currency.name)
				res += [(journal.id, name)]
			return res

			*/
			return rs.Super().NameGet()
		})

	h.AccountJournal().Methods().SearchByName().Extend("",
		func(rs h.AccountJournalSet, name string, op operator.Operator, additionalCond q.AccountJournalCondition, limit int) h.AccountJournalSet {
			//@api.model
			/*def name_search(self, name, args=None, operator='ilike', limit=100):
			args = args or []
			connector = '|'
			if operator in expression.NEGATIVE_TERM_OPERATORS:
				connector = '&'
			recs = self.search([connector, ('code', operator, name), ('name', operator, name)] + args, limit=limit)
			return recs.name_get()

			*/
			return rs.Super().SearchByName(name, op, additionalCond, limit)
		})

	h.AccountJournal().Methods().BelongToCompany().DeclareMethod(
		`BelongToCompany`,
		func(rs h.AccountJournalSet) *h.AccountJournalData {
			//@api.depends('company_id')
			/*def _belong_to_company(self):
			  for journal in self:
			      journal.belong_to_company = (journal.company_id.id == self.env.user.company_id.id)

			*/
			return new(h.AccountJournalData)
		})

	/*
		h.AccountJournal().Methods().SearchCompanyJournals().DeclareMethod(
			`SearchCompanyJournals`,
			func(rs h.AccountJournalSet, op operator.Operator, value string)
			}) {
				//@api.multi
	*/
	/*def _search_company_journals(self, operator, value):
	  if value:
	      recs = self.search([('company_id', operator, self.env.user.company_id.id)])
	  elif operator == '=':
	      recs = self.search([('company_id', '!=', self.env.user.company_id.id)])
	  else:
	      recs = self.search([('company_id', operator, self.env.user.company_id.id)])
	  return [('id', 'in', [x.id for x in recs])]

	*/ /*

		})
	*/

	h.AccountJournal().Methods().MethodsCompute().DeclareMethod(
		`MethodsCompute`,
		func(rs h.AccountJournalSet) *h.AccountJournalData {
			//@api.depends('inbound_payment_method_ids','outbound_payment_method_ids')
			/*def _methods_compute(self):
			  for journal in self:
			      journal.at_least_one_inbound = bool(len(journal.inbound_payment_method_ids))
			      journal.at_least_one_outbound = bool(len(journal.outbound_payment_method_ids))


			*/
			return new(h.AccountJournalData)
		})

	h.BankAccount().AddFields(map[string]models.FieldDefinition{
		"Journal": models.One2ManyField{RelationModel: h.AccountJournal(), ReverseFK: "BankAccount",
			JSON: "journal_id", Filter: q.AccountJournal().Type().Equals("bank"), ReadOnly: true,
			Help:       "The accounting journal corresponding to this bank account.",
			Constraint: h.BankAccount().Methods().CheckJournal()},
	})

	h.BankAccount().Methods().CheckJournal().DeclareMethod(
		`CheckJournal`,
		func(rs h.BankAccountSet) {
			//@api.constrains('journal_id')
			/*def _check_journal_id(self):
			  if len(self.journal_id) > 1:
			      raise ValidationError(_('A bank account can only belong to one journal.'))

			*/
		})

	h.AccountTaxGroup().DeclareModel()
	h.AccountTaxGroup().SetDefaultOrder("Sequence ASC")

	h.AccountTaxGroup().AddFields(map[string]models.FieldDefinition{
		"Name":     models.CharField{Required: true, Translate: true},
		"Sequence": models.IntegerField{Default: models.DefaultValue(10)},
	})

	h.AccountTax().DeclareModel()
	h.AccountTax().SetDefaultOrder("Sequence")

	h.AccountTax().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Tax Name", Required: true, Translate: true},
		"TypeTaxUse": models.SelectionField{String: "Tax Scope", Selection: types.Selection{
			"sale":     "Sales",
			"purchase": "Purchases",
			"none":     "None",
		}, Required: true, Default: models.DefaultValue("sale"),
			Constraint: h.AccountTax().Methods().CheckChildrenScope(),
			Help: `Determines where the tax is selectable.
Note: 'None' means a tax can't be used by itself however it can still be used in a group.`},
		"TaxAdjustment": models.BooleanField{String: "TaxAdjustment",
			Help: `Set this field to true if this tax can be used in the tax adjustment wizard,
used to manually fill some data in the tax declaration`},
		"AmountType": models.SelectionField{String: "Tax Computation", Selection: types.Selection{
			"group":    "Group of Taxes",
			"fixed":    "Fixed",
			"percent":  "Percentage of Price",
			"division": "Percentage of Price Tax Included",
		}, Required: true, Default: models.DefaultValue("percent")},
		"Active": models.BooleanField{String: "Active", Default: models.DefaultValue(true),
			Help: "Set active to false to hide the tax without removing it."},
		"Company": models.Many2OneField{RelationModel: h.Company(), Required: true,
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company()
			}},
		"ChildrenTaxes": models.Many2ManyField{RelationModel: h.AccountTax(), JSON: "children_tax_ids",
			M2MTheirField: "ChildTax", M2MOurField: "ParentTax",
			Constraint: h.AccountTax().Methods().CheckChildrenScope()},
		"Sequence": models.IntegerField{Required: true, GoType: new(int), Default: models.DefaultValue(1),
			Help: "The sequence field is used to define order in which the tax lines are applied."},
		"Amount": models.FloatField{Required: true, Digits: nbutils.Digits{Precision: 16, Scale: 4},
			OnChange: h.AccountTax().Methods().OnchangeAmount()},
		"Account": models.Many2OneField{String: "Tax Account",
			RelationModel: h.AccountAccount(), Filter: q.AccountAccount().Deprecated().Equals(false),
			OnDelete: models.Restrict,
			OnChange: h.AccountTax().Methods().OnchangeAccount(),
			Help:     "Account that will be set on invoice tax lines for invoices. Leave empty to use the expense account."},
		"RefundAccount": models.Many2OneField{String: "Tax Account on Refunds",
			RelationModel: h.AccountAccount(), Filter: q.AccountAccount().Deprecated().Equals(false),
			OnDelete: models.Restrict,
			Help:     "Account that will be set on invoice tax lines for refunds. Leave empty to use the expense account."},
		"Description": models.CharField{String: "Label on Invoices", Translate: true},
		"PriceInclude": models.BooleanField{String: "Included in Price", Default: models.DefaultValue(false),
			OnChange: h.AccountTax().Methods().OnchangePriceInclude(),
			Help:     "Check this if the price you use on the product and invoices includes this tax."},
		"IncludeBaseAmount": models.BooleanField{String: "Affect Base of Subsequent Taxes",
			Default: models.DefaultValue(false),
			Help:    "If set, taxes which are computed after this one will be computed based on the price tax included."},
		"Analytic": models.BooleanField{String: "Include in Analytic Cost",
			Help: `If set, the amount computed by this tax will be assigned
to the same analytic account as the invoice line (if any)`},
		"Tags": models.Many2ManyField{String: "Tags", RelationModel: h.AccountAccountTag(), JSON: "tag_ids",
			Help: "Optional tags you may want to assign for custom reporting"},
		"TaxGroup": models.Many2OneField{RelationModel: h.AccountTaxGroup(),
			Default: func(env models.Environment) interface{} {
				return h.AccountTaxGroup().NewSet(env).SearchAll().Limit(1)
			}, Required: true},
	})

	h.AccountTax().AddSQLConstraint("name_company_uniq", "unique(name, company_id, type_tax_use)",
		"Tax names must be unique !")

	h.AccountTax().Methods().Unlink().Extend("",
		func(rs h.AccountTaxSet) int64 {
			//@api.multi
			/*def unlink(self):
			  company_id = self.env.user.company_id.id
			  ir_values = self.env['ir.values']
			  supplier_taxes_id = set(ir_values.get_default('product.template', 'supplier_taxes_id', company_id=company_id) or [])
			  deleted_sup_tax = self.filtered(lambda tax: tax.id in supplier_taxes_id)
			  if deleted_sup_tax:
			      ir_values.sudo().set_default('product.template', "supplier_taxes_id", list(supplier_taxes_id - set(deleted_sup_tax.ids)), for_all_users=True, company_id=company_id)
			  taxes_id = set(self.env['ir.values'].get_default('product.template', 'taxes_id', company_id=company_id) or [])
			  deleted_tax = self.filtered(lambda tax: tax.id in taxes_id)
			  if deleted_tax:
			      ir_values.sudo().set_default('product.template', "taxes_id", list(taxes_id - set(deleted_tax.ids)), for_all_users=True, company_id=company_id)
			  return super(AccountTax, self).unlink()

			*/
			return rs.Super().Unlink()
		})

	h.AccountTax().Methods().CheckChildrenScope().DeclareMethod(
		`CheckChildrenScope`,
		func(rs h.AccountTaxSet) {
			//@api.constrains('children_tax_ids','type_tax_use')
			/*def _check_children_scope(self):
			  if not all(child.type_tax_use in ('none', self.type_tax_use) for child in self.children_tax_ids):
			      raise ValidationError(_('The application scope of taxes in a group must be either the same as the group or "None".'))

			*/
		})

	h.AccountTax().Methods().Copy().Extend("",
		func(rs h.AccountTaxSet, overrides *h.AccountTaxData, fieldsToReset ...models.FieldNamer) h.AccountTaxSet {
			//@api.returns('self',lambdavalue:value.id)
			/*def copy(self, default=None):
			default = dict(default or {}, name=_("%s (Copy)") % self.name)
			return super(AccountTax, self).copy(default=default)

			*/
			return rs.Super().Copy(overrides, fieldsToReset...)
		})

	h.AccountTax().Methods().SearchByName().Extend("",
		func(rs h.AccountTaxSet, name string, op operator.Operator, additionalCond q.AccountTaxCondition, limit int) h.AccountTaxSet {
			//@api.model
			/*def name_search(self, name, args=None, operator='ilike', limit=100):
			""" Returns a list of tupples containing id, name, as internally it is called {def name_get}
				result format: {[(id, name), (id, name), ...]}
			"""
			args = args or []
			if operator in expression.NEGATIVE_TERM_OPERATORS:
				domain = [('description', operator, name), ('name', operator, name)]
			else:
				domain = ['|', ('description', operator, name), ('name', operator, name)]
			taxes = self.search(expression.AND([domain, args]), limit=limit)
			return taxes.name_get()

			*/
			return rs.Super().SearchByName(name, op, additionalCond, limit)
		})

	h.AccountTax().Methods().Search().Extend("",
		func(rs h.AccountTaxSet, cond q.AccountTaxCondition) h.AccountTaxSet {
			//@api.model
			/*def search(self, args, offset=0, limit=None, order=None, count=False):
			  context = self._context or {}

			  if context.get('type'):
			      if context.get('type') in ('out_invoice', 'out_refund'):
			          args += [('type_tax_use', '=', 'sale')]
			      elif context.get('type') in ('in_invoice', 'in_refund'):
			          args += [('type_tax_use', '=', 'purchase')]

			  if context.get('journal_id'):
			      journal = self.env['account.journal'].browse(context.get('journal_id'))
			      if journal.type in ('sale', 'purchase'):
			          args += [('type_tax_use', '=', journal.type)]

			  return super(AccountTax, self).search(args, offset, limit, order, count=count)

			*/
			return rs.Super().Search(cond)
		})

	h.AccountTax().Methods().OnchangeAmount().DeclareMethod(
		`OnchangeAmount`,
		func(rs h.AccountTaxSet) (*h.AccountTaxData, []models.FieldNamer) {
			//@api.onchange('amount')
			/*def onchange_amount(self):
			  if self.amount_type in ('percent', 'division') and self.amount != 0.0 and not self.description:
			      self.description = "{0:.4g}%".format(self.amount)

			*/
			return new(h.AccountTaxData), []models.FieldNamer{}
		})

	h.AccountTax().Methods().OnchangeAccount().DeclareMethod(
		`OnchangeAccount`,
		func(rs h.AccountTaxSet) (*h.AccountTaxData, []models.FieldNamer) {
			//@api.onchange('account_id')
			/*def onchange_account_id(self):
			  self.refund_account_id = self.account_id

			*/
			return new(h.AccountTaxData), []models.FieldNamer{}
		})

	h.AccountTax().Methods().OnchangePriceInclude().DeclareMethod(
		`OnchangePriceInclude`,
		func(rs h.AccountTaxSet) (*h.AccountTaxData, []models.FieldNamer) {
			//@api.onchange('price_include')
			/*def onchange_price_include(self):
			  if self.price_include:
			      self.include_base_amount = True

			*/
			return new(h.AccountTaxData), []models.FieldNamer{}
		})

	h.AccountTax().Methods().GetGroupingKey().DeclareMethod(
		`GetGroupingKey`,
		func(rs h.AccountTaxSet, invoiceTaxVal *h.AccountInvoiceTaxData) string {
			/*def get_grouping_key(self, invoice_tax_val):
			  """ Returns a string that will be used to group account.invoice.tax sharing the same properties"""
			  self.ensure_one()
			  return str(invoice_tax_val['tax_id']) + '-' + str(invoice_tax_val['account_id']) + '-' + str(invoice_tax_val['account_analytic_id'])

			*/
			return ""
		})

	h.AccountTax().Methods().ComputeAmount().DeclareMethod(
		`ComputeAmount returns the amount of a single tax.

		baseAmount is the actual amount on which the tax is applied, which is priceUnit * quantity eventually
		affected by previous taxes (if tax is include_base_amount XOR price_include)`,
		func(rs h.AccountTaxSet, baseAmount, priceUnit, quantity float64, product h.ProductProductSet, partner h.PartnerSet) float64 {
			rs.EnsureOne()
			if rs.AmountType() == "fixed" {
				// Use Copysign to take into account the sign of the base amount which includes the sign
				// of the quantity and the sign of the priceUnit
				// Amount is the fixed price for the tax, it can be negative
				// Base amount included the sign of the quantity and the sign of the unit price and when
				// a product is returned, it can be done either by changing the sign of quantity or by changing the
				// sign of the price unit.
				// When the price unit is equal to 0, the sign of the quantity is absorbed in base_amount then
				// a "else" case is needed.
				if baseAmount != 0 {
					return math.Copysign(quantity, baseAmount) * rs.Amount()
				}
				return quantity * rs.Amount()
			}
			if (rs.AmountType() == "percent" && !rs.PriceInclude()) || (rs.AmountType() == "division" && rs.PriceInclude()) {
				return baseAmount * rs.Amount() / 100
			}
			if rs.AmountType() == "percent" && rs.PriceInclude() {
				return baseAmount - (baseAmount / (1 + rs.Amount()/100))
			}
			if rs.AmountType() == "division" && !rs.PriceInclude() {
				return baseAmount/(1-rs.Amount()/100) - baseAmount
			}
			log.Fatal("Unhandled tax type", "tax", rs.ID(), "type", rs.AmountType(), "priceInclude", rs.PriceInclude())
			panic("Unhandled tax type")
		})

	h.AccountTax().Methods().JSONFriendlyComputeAll().DeclareMethod(
		`JSONFriendlyComputeAll`,
		func(rs h.AccountTaxSet, priceUnit float64, currencyID int64, quantity float64, productID int64, partnerID int64) float64 {
			//@api.multi
			/*def json_friendly_compute_all(self, price_unit, currency_id=None, quantity=1.0, product_id=None, partner_id=None):
			  """ Just converts parameters in browse records and calls for compute_all, because js widgets can't serialize browse records """
			  if currency_id:
			      currency_id = self.env['res.currency'].browse(currency_id)
			  if product_id:
			      product_id = self.env['product.product'].browse(product_id)
			  if partner_id:
			      partner_id = self.env['res.partner'].browse(partner_id)
			  return self.compute_all(price_unit, currency=currency_id, quantity=quantity, product=product_id, partner=partner_id)

			*/
			return 0
		})

	h.AccountTax().Methods().ComputeAll().DeclareMethod(
		`ComputeAll returns all information required to apply taxes (in self + their children in case of a tax goup).
			      We consider the sequence of the parent for group of taxes.
			          Eg. considering letters as taxes and alphabetic order as sequence :
			          [G, B([A, D, F]), E, C] will be computed as [A, D, F, C, E, G]

			  RETURN:

                   0.0,                 # Base

			       0.0,                 # Total without taxes

			       0.0,                 # Total with taxes

                   []AppliedTaxData     # One struct for each tax in rs and their children
			  } `,
		func(rs h.AccountTaxSet, priceUnit float64, currency h.CurrencySet, quantity float64,
			product h.ProductProductSet, partner h.PartnerSet) (float64, float64, float64, []accounttypes.AppliedTaxData) {

			company := rs.Company()
			if rs.IsEmpty() {
				company = h.User().NewSet(rs.Env()).CurrentUser().Company()
			}
			if currency.IsEmpty() {
				currency = company.Currency()
			}
			var taxes []accounttypes.AppliedTaxData
			// By default, for each tax, tax amount will first be computed
			// and rounded at the 'Account' decimal precision for each
			// PO/SO/invoice line and then these rounded amounts will be
			// summed, leading to the total amount for that tax. But, if the
			// company has tax_calculation_rounding_method = round_globally,
			// we still follow the same method, but we use a much larger
			// precision when we round the tax amount for each line (we use
			// the 'Account' decimal precision + 5), and that way it's like
			// rounding after the sum of the tax amounts of each line

			dp := currency.DecimalPlaces()
			// In some cases, it is necessary to force/prevent the rounding of the tax and the total
			// amounts. For example, in SO/PO line, we don't want to round the price unit at the
			// precision of the currency.
			// The context key 'round' allows to force the standard behavior.
			roundTax := true
			if company.TaxCalculationRoundingMethod() == "round_globally" {
				roundTax = false
			}
			roundTotal := true
			if rs.Env().Context().HasKey("round") {
				roundTax = rs.Env().Context().GetBool("round")
				roundTotal = rs.Env().Context().GetBool("round")
			}
			if !roundTax {
				dp += 5
			}
			prec := math.Pow10(-dp)
			totalExcluded := nbutils.Round(priceUnit*quantity, prec)
			totalIncluded := nbutils.Round(priceUnit*quantity, prec)
			base := nbutils.Round(priceUnit*quantity, prec)
			baseValues := rs.Env().Context().GetFloatSlice("base_values")
			if len(baseValues) != 0 {
				totalExcluded = baseValues[0]
				totalIncluded = baseValues[1]
				base = baseValues[2]
			}
			// Sorting key is mandatory in this case. When no key is provided, sorted() will perform a
			// search. However, the search method is overridden in account.tax in order to add a domain
			// depending on the context. This domain might filter out some taxes from self, e.g. in the
			// case of group taxes.
			taxRecords := rs.Records()
			sort.Slice(taxRecords, func(i, j int) bool {
				return taxRecords[i].Sequence() < taxRecords[j].Sequence()
			})
			for _, tax := range taxRecords {
				if tax.AmountType() == "group" {
					children := tax.ChildrenTaxes().WithContext("base_values", []float64{totalExcluded, totalIncluded, base})
					retBase, retExcl, retIncl, retTaxes := children.ComputeAll(priceUnit, currency, quantity, product, partner)
					totalExcluded = retExcl
					if tax.IncludeBaseAmount() {
						base = retBase
					}
					totalIncluded = retIncl
					taxes = append(taxes, retTaxes...)
					continue
				}

				taxAmount := tax.ComputeAmount(base, priceUnit, quantity, product, partner)
				if roundTax {
					taxAmount = nbutils.Round(taxAmount, prec)
				} else {
					taxAmount = currency.Round(taxAmount)
				}

				if tax.PriceInclude() {
					totalExcluded -= taxAmount
					base -= taxAmount
				} else {
					totalIncluded += taxAmount
				}

				// Keep base amount used for the current tax
				taxBase := base

				if tax.IncludeBaseAmount() {
					base += taxAmount
				}

				taxes = append(taxes, accounttypes.AppliedTaxData{
					ID:              tax.ID(),
					Name:            tax.WithContext("lang", partner.Lang()).Name(),
					Amount:          taxAmount,
					Base:            taxBase,
					Sequence:        tax.Sequence(),
					AccountID:       tax.Account().ID(),
					RefundAccountID: tax.RefundAccount().ID(),
					Analytic:        tax.Analytic(),
				})
			}

			if roundTotal {
				totalIncluded = currency.Round(totalIncluded)
				totalExcluded = currency.Round(totalExcluded)
			}
			sort.Slice(taxes, func(i, j int) bool {
				return taxes[i].Sequence < taxes[j].Sequence
			})
			return base, totalExcluded, totalIncluded, taxes
		})

	h.AccountTax().Methods().FixTaxIncludedPrice().DeclareMethod(
		`FixTaxIncludedPrice`,
		func(rs h.AccountTaxSet, price float64, prodTaxes, lineTaxes h.AccountTaxSet) float64 {
			//@api.model
			/*def _fix_tax_included_price(self, price, prod_taxes, line_taxes):
			  """Subtract tax amount from price when corresponding "price included" taxes do not apply"""
			  # FIXME get currency in param?
			  incl_tax = prod_taxes.filtered(lambda tax: tax not in line_taxes and tax.price_include)
			  if incl_tax:
			      return incl_tax.compute_all(price)['total_excluded']
			  return price

			*/
			return 0
		})

	h.AccountReconcileModel().DeclareModel()

	h.AccountReconcileModel().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Button Label", Required: true,
			OnChange: h.AccountReconcileModel().Methods().OnchangeName()},
		"Sequence":      models.IntegerField{Required: true, Default: models.DefaultValue(10)},
		"HasSecondLine": models.BooleanField{String: "Add a second line", Default: models.DefaultValue(false)},
		"Company": models.Many2OneField{RelationModel: h.Company(), Required: true,
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company()
			}},
		"Account": models.Many2OneField{RelationModel: h.AccountAccount(),
			OnDelete: models.Cascade,
			Filter:   q.AccountAccount().Deprecated().Equals(false)},
		"Journal": models.Many2OneField{RelationModel: h.AccountJournal(),
			OnDelete: models.Cascade, Help: "This field is ignored in a bank statement reconciliation."},
		"Label": models.CharField{String: "Journal Item Label"},
		"AmountType": models.SelectionField{Selection: types.Selection{
			"fixed":      "Fixed",
			"percentage": "Percentage of balance"}, Required: true, Default: models.DefaultValue("percentage")},
		"Amount": models.FloatField{Required: true, Default: models.DefaultValue(100.0),
			Help: "Fixed amount will count as a debit if it is negative, as a credit if it is positive."},
		"Tax": models.Many2OneField{String: "Tax", RelationModel: h.AccountTax(),
			OnDelete: models.Restrict, Filter: q.AccountTax().TypeTaxUse().Equals("purchase")},
		"AnalyticAccount": models.Many2OneField{RelationModel: h.AccountAnalyticAccount(),
			OnDelete: models.SetNull},
		"SecondAccount": models.Many2OneField{RelationModel: h.AccountAccount(),
			OnDelete: models.Cascade, Filter: q.AccountAccount().Deprecated().Equals(false),
		},
		"SecondJournal": models.Many2OneField{RelationModel: h.AccountJournal(),
			OnDelete: models.Cascade, Help: "This field is ignored in a bank statement reconciliation."},
		"SecondLabel": models.CharField{String: "Second Journal Item Label"},
		"SecondAmountType": models.SelectionField{Selection: types.Selection{
			"fixed":      "Fixed",
			"percentage": "Percentage of balance"}, Required: true, Default: models.DefaultValue("percentage")},
		"SecondAmount": models.FloatField{Required: true, Default: models.DefaultValue(100.0),
			Help: "Fixed amount will count as a debit if it is negative, as a credit if it is positive."},
		"SecondTax": models.Many2OneField{RelationModel: h.AccountTax(),
			OnDelete: models.Restrict, Filter: q.AccountTax().TypeTaxUse().Equals("purchase")},
		"SecondAnalyticAccount": models.Many2OneField{RelationModel: h.AccountAnalyticAccount(),
			OnDelete: models.SetNull},
	})

	h.AccountReconcileModel().Methods().OnchangeName().DeclareMethod(
		`OnchangeName`,
		func(rs h.AccountReconcileModelSet) (*h.AccountReconcileModelData, []models.FieldNamer) {
			//@api.onchange('name')
			/*def onchange_name(self):
			  self.label = self.name
			*/
			return new(h.AccountReconcileModelData), []models.FieldNamer{}
		})

}
