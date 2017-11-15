// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountAccountType().DeclareModel()
	pool.AccountAccountType().AddFields(map[string]models.FieldDefinition{
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

	pool.AccountAccountTag().DeclareModel()
	pool.AccountAccountTag().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Name", Required: true},
		"Applicability": models.SelectionField{String: "Applicability", Selection: types.Selection{
			"accounts": "Accounts",
			"taxes":    "Taxes",
		}, Required: true, Default: models.DefaultValue("accounts")},
		"Color": models.IntegerField{String: "Color Index"},
	})

	pool.AccountAccount().DeclareModel()
	pool.AccountAccount().SetDefaultOrder("Code")

	pool.AccountAccount().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Required: true, Index: true},
		"Currency": models.Many2OneField{String: "Account Currency", RelationModel: pool.Currency(),
			Help: "Forces all moves for this account to have this account currency."},
		"Code":       models.CharField{Size: 64, Required: true, Index: true},
		"Deprecated": models.BooleanField{Index: true, Default: models.DefaultValue(false)},
		"UserType": models.Many2OneField{String: "Type", RelationModel: pool.AccountAccountType(),
			Required: true, Help: `Account Type is used for information purpose, to generate country-specific
legal reports, and set the rules to close a fiscal year and generate opening entries.`},
		"InternalType": models.SelectionField{Related: "userType.Type", /*readonly=True)*/
			Constraint: pool.AccountAccount().Methods().CheckReconcile(),
			OnChange:   pool.AccountAccount().Methods().OnchangeInternalType()},
		"LastTimeEntriesChecked": models.DateTimeField{String: "Latest Invoices & Payments Matching Date", /*[ readonly True]*/
			Help: `Last time the invoices & payments matching was performed on this account.
It is set either if there's not at least an unreconciled debit and an unreconciled credit
or if you click the "Done" button.`, NoCopy: true},
		"Reconcile": models.BooleanField{String: "Allow Reconciliation", Default: models.DefaultValue(false),
			Constraint: pool.AccountAccount().Methods().CheckReconcile(),
			Help:       "Check this box if this account allows invoices & payments matching of journal items."},
		"Taxes": models.Many2ManyField{String: "Default Taxes", RelationModel: pool.AccountTax(), JSON: "tax_ids"},
		"Note":  models.TextField{String: "Internal Notes"},
		"Company": models.Many2OneField{RelationModel: pool.Company(), Required: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.Company().NewSet(env).CompanyDefaultGet()
			}},
		"Tags": models.Many2ManyField{RelationModel: pool.AccountAccountTag(), JSON: "tag_ids",
			Help: "Optional tags you may want to assign for custom reporting"},
	})

	pool.AccountAccount().AddSQLConstraint("code_company_uniq", "unique (code,company_id)",
		"The code of the account must be unique per company !")

	pool.AccountAccount().Methods().CheckReconcile().DeclareMethod(
		`CheckReconcile`,
		func(rs pool.AccountAccountSet) {
			//@api.constrains('internal_type','reconcile')
			/*def _check_reconcile(self):
			  for account in self:
			      if account.internal_type in ('receivable', 'payable') and account.reconcile == False:
			          raise ValidationError(_('You cannot have a receivable/payable account that is not reconciliable. (account code: %s)') % account.code)
			*/
		})

	pool.AccountAccount().Methods().DefaultGet().Extend("",
		func(rs pool.AccountAccountSet) models.FieldMap {
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

	pool.AccountAccount().Methods().SearchByName().Extend("",
		func(rs pool.AccountAccountSet, name string, op operator.Operator, additionalCond pool.AccountAccountCondition, limit int) pool.AccountAccountSet {
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

	pool.AccountAccount().Methods().OnchangeInternalType().DeclareMethod(
		`OnchangeInternalType`,
		func(rs pool.AccountAccountSet) (*pool.AccountAccountSet, []models.FieldNamer) {
			//@api.onchange('internal_type')
			/*def onchange_internal_type(self):
			  if self.internal_type in ('receivable', 'payable'):
			      self.reconcile = True

			*/
			return new(pool.AccountAccountSet), []models.FieldNamer{}
		})

	pool.AccountAccount().Methods().NameGet().Extend("",
		func(rs pool.AccountAccountSet) string {
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

	pool.AccountAccount().Methods().Copy().Extend("",
		func(rs pool.AccountAccountSet, overrides *pool.AccountAccountData, fieldsToReset ...models.FieldNamer) pool.AccountAccountSet {
			//@api.returns('self',lambdavalue:value.id)
			/*def copy(self, default=None):
			  default = dict(default or {})
			  default.setdefault('code', _("%s (copy)") % (self.code or ''))
			  return super(AccountAccount, self).copy(default)

			*/
			return rs.Super().Copy(overrides, fieldsToReset...)
		})

	pool.AccountAccount().Methods().Write().Extend("",
		func(rs pool.AccountAccountSet, vals *pool.AccountAccountData, fieldsToReset ...models.FieldNamer) bool {
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

	pool.AccountAccount().Methods().Unlink().Extend("",
		func(rs pool.AccountAccountSet) int64 {
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

	pool.AccountAccount().Methods().MarkAsReconciled().DeclareMethod(
		`MarkAsReconciled`,
		func(rs pool.AccountAccountSet) bool {
			//@api.multi
			/*def mark_as_reconciled(self):
			  return self.write({'last_time_entries_checked': time.strftime(DEFAULT_SERVER_DATETIME_FORMAT)})

			*/
			return true
		})

	pool.AccountAccount().Methods().ActionOpenReconcile().DeclareMethod(
		`ActionOpenReconcile`,
		func(rs pool.AccountAccountSet) *actions.Action {
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

	pool.AccountJournal().DeclareModel()
	pool.AccountJournal().SetDefaultOrder("Sequence", "Type", "Code")

	pool.AccountJournal().AddFields(map[string]models.FieldDefinition{
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
			Constraint: pool.AccountJournal().Methods().CheckBankAccount(),
			Help: `Select 'Sale' for customer invoices journals.
Select 'Purchase' for vendor bills journals.
Select 'Cash' or 'Bank' for journals that are used in customer or vendor payments.
Select 'General' for miscellaneous operations journals.`},
		"TypeControls": models.Many2ManyField{String: "Account Types Allowed",
			RelationModel: pool.AccountAccountType(), JSON: "type_control_ids"},
		"AccountControls": models.Many2ManyField{String: "Accounts Allowed", RelationModel: pool.AccountAccount(),
			JSON: "account_control_ids", Filter: pool.AccountAccount().Deprecated().Equals(false)},
		"DefaultCreditAccount": models.Many2OneField{RelationModel: pool.AccountAccount(),
			Constraint: pool.AccountJournal().Methods().CheckCurrency(),
			OnChange:   pool.AccountJournal().Methods().OnchangeCreditAccountId(),
			Filter:     pool.AccountAccount().Deprecated().Equals(false),
			Help:       "It acts as a default account for credit amount"},
		"DefaultDebitAccount": models.Many2OneField{RelationModel: pool.AccountAccount(),
			Constraint: pool.AccountJournal().Methods().CheckCurrency(),
			OnChange:   pool.AccountJournal().Methods().OnchangeDebitAccountId(),
			Filter:     pool.AccountAccount().Deprecated().Equals(false),
			Help:       "It acts as a default account for debit amount"},
		"UpdatePosted": models.BooleanField{String: "Allow Cancelling Entries",
			Help: `Check this box if you want to allow the cancellation the entries related to this journal or
of the invoice related to this journal`},
		"GroupInvoiceLines": models.BooleanField{
			Help: `If this box is checked, the system will try to group the accounting lines when generating
them from invoices.`},
		"EntrySequence": models.Many2OneField{RelationModel: pool.Sequence(), JSON: "sequence_id",
			Help:     "This field contains the information related to the numbering of the journal entries of this journal.",
			Required: true, NoCopy: true},
		"RefundEntrySequence": models.Many2OneField{RelationModel: pool.Sequence(), JSON: "refund_sequence_id",
			Help:   "This field contains the information related to the numbering of the refund entries of this journal.",
			NoCopy: true},
		"Sequence": models.IntegerField{
			Help:    "Used to order Journals in the dashboard view', default=10",
			Default: models.DefaultValue(10)},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(),
			Constraint: pool.AccountJournal().Methods().CheckCurrency(),
			Help:       "The currency used to enter statement"},
		"Company": models.Many2OneField{RelationModel: pool.Company(), Required: true, Index: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser().Company()
			}, Help: "Company related to this journal"},
		"RefundSequence": models.BooleanField{String: "Dedicated Refund Sequence",
			Help: `Check this box if you don't want to share the
same sequence for invoices and refunds made from this journal`,
			Default: models.DefaultValue(false)},
		"InboundPaymentMethods": models.Many2ManyField{String: "Debit Methods",
			RelationModel: pool.AccountPaymentMethod(), JSON: "inbound_payment_method_ids",
			Filter: pool.AccountPaymentMethod().PaymentType().Equals("inbound"),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.AccountPaymentMethod().Search(env,
					pool.AccountPaymentMethod().HexyaExternalID().Equals("account_account_payment_method_manual_in"))
			},
			Help: `Means of payment for collecting money.
Hexya modules offer various payments handling facilities,
but you can always use the 'Manual' payment method in order
to manage payments outside of the software.`},
		"OutboundPaymentMethods": models.Many2ManyField{String: "Payment Methods",
			RelationModel: pool.AccountPaymentMethod(), JSON: "outbound_payment_method_ids",
			Filter: pool.AccountPaymentMethod().PaymentType().Equals("outbound"),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.AccountPaymentMethod().Search(env,
					pool.AccountPaymentMethod().HexyaExternalID().Equals("account_account_payment_method_manual_out"))
			}, Help: `Means of payment for sending money.
Hexya modules offer various payments handling facilities
but you can always use the 'Manual' payment method in order
to manage payments outside of the software.`},
		"AtLeastOneInbound": models.BooleanField{Compute: pool.AccountJournal().Methods().MethodsCompute(),
			Stored: true},
		"AtLeastOneOutbound": models.BooleanField{Compute: pool.AccountJournal().Methods().MethodsCompute(),
			Stored: true},
		"ProfitAccount": models.Many2OneField{RelationModel: pool.AccountAccount(),
			Filter: pool.AccountAccount().Deprecated().Equals(false),
			Help:   "Used to register a profit when the ending balance of a cash register differs from what the system computes"},
		"LossAccount": models.Many2OneField{RelationModel: pool.AccountAccount(),
			Filter: pool.AccountAccount().Deprecated().Equals(false),
			Help:   "Used to register a loss when the ending balance of a cash register differs from what the system computes"},
		"BelongsToCompany": models.BooleanField{String: "Belong to the user's current company",
			Compute: pool.AccountJournal().Methods().BelongToCompany() /*[ search "_search_company_journals"]*/},

		"BankAccount": models.Many2OneField{RelationModel: pool.BankAccount(), OnDelete: models.Restrict,
			Constraint: pool.AccountJournal().Methods().CheckBankAccount(),
			NoCopy:     true},
		"DisplayOnFooter": models.BooleanField{String: "Show in Invoices Footer",
			Help: "Display this bank account on the footer of printed documents like invoices and sales orders."},
		"BankStatementsSource": models.SelectionField{String: "Bank Feeds", Selection: types.Selection{
			"manual": "Record Manually",
		}},
		"BankAccNumber": models.CharField{Related: "BankAccount.AccNumber"},
		"Bank":          models.Many2OneField{RelationModel: pool.Bank(), Related: "BankAccount.Bank"},
	})

	pool.AccountJournal().AddSQLConstraint("code_company_uniq", "unique (code, name, company_id)",
		"The code and name of the journal must be unique per company !'")

	pool.AccountJournal().Methods().CheckCurrency().DeclareMethod(
		`CheckCurrency`,
		func(rs pool.AccountJournalSet) {
			//@api.constrains('currency_id','default_credit_account_id','default_debit_account_id')
			/*def _check_currency(self):
			  if self.currency_id:
			      if self.default_credit_account_id and not self.default_credit_account_id.currency_id.id == self.currency_id.id:
			          raise ValidationError(_('Configuration error!\nThe currency of the journal should be the same than the default credit account.'))
			      if self.default_debit_account_id and not self.default_debit_account_id.currency_id.id == self.currency_id.id:
			          raise ValidationError(_('Configuration error!\nThe currency of the journal should be the same than the default debit account.'))

			*/
		})

	pool.AccountJournal().Methods().CheckBankAccount().DeclareMethod(
		`CheckBankAccount`,
		func(rs pool.AccountJournalSet) {
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

	pool.AccountJournal().Methods().OnchangeDebitAccountId().DeclareMethod(
		`OnchangeDebitAccountId`,
		func(rs pool.AccountJournalSet) (*pool.AccountJournalData, []models.FieldNamer) {
			//@api.onchange('default_debit_account_id')
			/*def onchange_debit_account_id(self):
			  if not self.default_credit_account_id:
			      self.default_credit_account_id = self.default_debit_account_id

			*/
			return new(pool.AccountJournalData), []models.FieldNamer{}
		})

	pool.AccountJournal().Methods().OnchangeCreditAccountId().DeclareMethod(
		`OnchangeCreditAccountId`,
		func(rs pool.AccountJournalSet) (*pool.AccountJournalData, []models.FieldNamer) {
			//@api.onchange('default_credit_account_id')
			/*def onchange_credit_account_id(self):
			  if not self.default_debit_account_id:
			      self.default_debit_account_id = self.default_credit_account_id

			*/
			return new(pool.AccountJournalData), []models.FieldNamer{}
		})

	pool.AccountJournal().Methods().Unlink().Extend("",
		func(rs pool.AccountJournalSet) int64 {
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

	pool.AccountJournal().Methods().Copy().Extend("",
		func(rs pool.AccountJournalSet, overrides *pool.AccountJournalData, fieldsToReset ...models.FieldNamer) pool.AccountJournalSet {
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

	pool.AccountJournal().Methods().Write().Extend("",
		func(rs pool.AccountJournalSet, vals *pool.AccountJournalData, fieldsToReset ...models.FieldNamer) bool {
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

	pool.AccountJournal().Methods().GetSequencePrefix().DeclareMethod(
		`GetSequencePrefix`,
		func(rs pool.AccountJournalSet, code string, refund bool) string {
			//@api.model
			/*def _get_sequence_prefix(self, code, refund=False):
						  prefix = code.upper()
						  if refund:
						      prefix = 'R' + prefix
						  return prefix + '/%(range_year)s/'

			x			*/
			return ""
		})

	pool.AccountJournal().Methods().CreateSequence().DeclareMethod(
		`CreateSequence`,
		func(rs pool.AccountJournalSet, refund bool, vals *pool.AccountJournalData, fieldsToReset ...models.FieldNamer) {
			//@api.model
			/*def _create_sequence(self, vals, refund=False):
			  """ Create new no_gap entry sequence for every new Journal"""
			  prefix = self._get_sequence_prefix(vals['code'], refund)
			  seq = {
			      'name': refund and vals['name'] + _(': Refund') or vals['name'],
			      'implementation': 'no_gap',
			      'prefix': prefix,
			      'padding': 4,
			      'number_increment': 1,
			      'use_date_range': True,
			  }
			  if 'company_id' in vals:
			      seq['company_id'] = vals['company_id']
			  return self.env['ir.sequence'].create(seq)

			*/
		})

	pool.AccountJournal().Methods().PrepareLiquidityAccount().DeclareMethod(
		`PrepareLiquidityAccount`,
		func(rs pool.AccountJournalSet, name string, company pool.CompanySet, currency pool.CurrencySet, accType string) *pool.AccountAccountData {
			//@api.model
			/*def _prepare_liquidity_account(self, name, company, currency_id, type):
			  '''
			  This function prepares the value to use for the creation of the default debit and credit accounts of a
			  liquidity journal (created through the wizard of generating COA from templates for example).

			  :param name: name of the bank account
			  :param company: company for which the wizard is running
			  :param currency_id: ID of the currency in wich is the bank account
			  :param type: either 'cash' or 'bank'
			  :return: mapping of field names and values
			  :rtype: dict
			  '''

			  # Seek the next available number for the account code
			  code_digits = company.accounts_code_digits or 0
			  if type == 'bank':
			      account_code_prefix = company.bank_account_code_prefix or ''
			  else:
			      account_code_prefix = company.cash_account_code_prefix or company.bank_account_code_prefix or ''
			  for num in xrange(1, 100):
			      new_code = str(account_code_prefix.ljust(code_digits - 1, '0')) + str(num)
			      rec = self.env['account.account'].search([('code', '=', new_code), ('company_id', '=', company.id)], limit=1)
			      if not rec:
			          break
			  else:
			      raise UserError(_('Cannot generate an unused account code.'))

			  liquidity_type = self.env.ref('account.data_account_type_liquidity')
			  return {
			          'name': name,
			          'currency_id': currency_id or False,
			          'code': new_code,
			          'user_type_id': liquidity_type and liquidity_type.id or False,
			          'company_id': company.id,
			  }

			*/
			return &pool.AccountAccountData{}
		})

	pool.AccountJournal().Methods().Create().Extend("",
		func(rs pool.AccountJournalSet, vals *pool.AccountJournalData) pool.AccountJournalSet {
			//@api.model
			/*def create(self, vals):
			  company_id = vals.get('company_id', self.env.user.company_id.id)
			  if vals.get('type') in ('bank', 'cash'):
			      # For convenience, the name can be inferred from account number
			      if not vals.get('name') and 'bank_acc_number' in vals:
			          vals['name'] = vals['bank_acc_number']

			      # If no code provided, loop to find next available journal code
			      if not vals.get('code'):
			          journal_code_base = (vals['type'] == 'cash' and 'CSH' or 'BNK')
			          journals = self.env['account.journal'].search([('code', 'like', journal_code_base + '%'), ('company_id', '=', company_id)])
			          for num in xrange(1, 100):
			              # journal_code has a maximal size of 5, hence we can enforce the boundary num < 100
			              journal_code = journal_code_base + str(num)
			              if journal_code not in journals.mapped('code'):
			                  vals['code'] = journal_code
			                  break
			          else:
			              raise UserError(_("Cannot generate an unused journal code. Please fill the 'Shortcode' field."))

			      # Create a default debit/credit account if not given
			      default_account = vals.get('default_debit_account_id') or vals.get('default_credit_account_id')
			      if not default_account:
			          company = self.env['res.company'].browse(company_id)
			          account_vals = self._prepare_liquidity_account(vals.get('name'), company, vals.get('currency_id'), vals.get('type'))
			          default_account = self.env['account.account'].create(account_vals)
			          vals['default_debit_account_id'] = default_account.id
			          vals['default_credit_account_id'] = default_account.id

			  # We just need to create the relevant sequences according to the chosen options
			  if not vals.get('sequence_id'):
			      vals.update({'sequence_id': self.sudo()._create_sequence(vals).id})
			  if vals.get('type') in ('sale', 'purchase') and vals.get('refund_sequence') and not vals.get('refund_sequence_id'):
			      vals.update({'refund_sequence_id': self.sudo()._create_sequence(vals, refund=True).id})

			  journal = super(AccountJournal, self).create(vals)

			  # Create the bank_account_id if necessary
			  if journal.type == 'bank' and not journal.bank_account_id and vals.get('bank_acc_number'):
			      journal.set_bank_account(vals.get('bank_acc_number'), vals.get('bank_id'))

			  return journal

			*/
			return rs.Super().Create(vals)
		})

	pool.AccountJournal().Methods().DefineBankAccount().DeclareMethod(
		`DefineBankAccount`,
		func(rs pool.AccountJournalSet, accNumber string, bank pool.BankSet) {
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

	pool.AccountJournal().Methods().NameGet().Extend("",
		func(rs pool.AccountJournalSet) string {
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

	pool.AccountJournal().Methods().SearchByName().Extend("",
		func(rs pool.AccountJournalSet, name string, op operator.Operator, additionalCond pool.AccountJournalCondition, limit int) pool.AccountJournalSet {
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

	pool.AccountJournal().Methods().BelongToCompany().DeclareMethod(
		`BelongToCompany`,
		func(rs pool.AccountJournalSet) (*pool.AccountJournalData, []models.FieldNamer) {
			//@api.depends('company_id')
			/*def _belong_to_company(self):
			  for journal in self:
			      journal.belong_to_company = (journal.company_id.id == self.env.user.company_id.id)

			*/
			return new(pool.AccountJournalData), []models.FieldNamer{}
		})

	/*
		pool.AccountJournal().Methods().SearchCompanyJournals().DeclareMethod(
			`SearchCompanyJournals`,
			func(rs pool.AccountJournalSet, op operator.Operator, value string)
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

	pool.AccountJournal().Methods().MethodsCompute().DeclareMethod(
		`MethodsCompute`,
		func(rs pool.AccountJournalSet) (*pool.AccountJournalData, []models.FieldNamer) {
			//@api.depends('inbound_payment_method_ids','outbound_payment_method_ids')
			/*def _methods_compute(self):
			  for journal in self:
			      journal.at_least_one_inbound = bool(len(journal.inbound_payment_method_ids))
			      journal.at_least_one_outbound = bool(len(journal.outbound_payment_method_ids))


			*/
			return new(pool.AccountJournalData), []models.FieldNamer{}
		})

	pool.BankAccount().AddFields(map[string]models.FieldDefinition{
		"Journal": models.One2ManyField{RelationModel: pool.AccountJournal(), ReverseFK: "BankAccount",
			JSON: "journal_id", Filter: pool.AccountJournal().Type().Equals("bank"), /* readonly=True */
			Help:       "The accounting journal corresponding to this bank account.",
			Constraint: pool.BankAccount().Methods().CheckJournal()},
	})

	pool.BankAccount().Methods().CheckJournal().DeclareMethod(
		`CheckJournal`,
		func(rs pool.BankAccountSet) {
			//@api.constrains('journal_id')
			/*def _check_journal_id(self):
			  if len(self.journal_id) > 1:
			      raise ValidationError(_('A bank account can only belong to one journal.'))

			*/
		})

	pool.AccountTaxGroup().DeclareModel()
	pool.AccountTaxGroup().SetDefaultOrder("Sequence ASC")

	pool.AccountTaxGroup().AddFields(map[string]models.FieldDefinition{
		"Name":     models.CharField{Required: true, Translate: true},
		"Sequence": models.IntegerField{Default: models.DefaultValue(10)},
	})

	pool.AccountTax().DeclareModel()
	pool.AccountTax().SetDefaultOrder("Sequence")

	pool.AccountTax().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Tax Name", Required: true, Translate: true},
		"TypeTaxUse": models.SelectionField{String: "Tax Scope", Selection: types.Selection{
			"sale":     "Sales",
			"purchase": "Purchases",
			"none":     "None",
		}, Required: true, Default: models.DefaultValue("sale"),
			Constraint: pool.AccountTax().Methods().CheckChildrenScope(),
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
		"Company": models.Many2OneField{RelationModel: pool.Company(), Required: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser().Company()
			}},
		"ChildrenTaxes": models.Many2ManyField{RelationModel: pool.AccountTax(), JSON: "children_tax_ids",
			Constraint: pool.AccountTax().Methods().CheckChildrenScope()},
		"Sequence": models.IntegerField{Required: true, Default: models.DefaultValue(1),
			Help: "The sequence field is used to define order in which the tax lines are applied."},
		"Amount": models.FloatField{Required: true, Digits: nbutils.Digits{Precision: 16, Scale: 4},
			OnChange: pool.AccountTax().Methods().OnchangeAmount()},
		"Account": models.Many2OneField{String: "Tax Account",
			RelationModel: pool.AccountAccount(), Filter: pool.AccountAccount().Deprecated().Equals(false),
			OnDelete: models.Restrict,
			OnChange: pool.AccountTax().Methods().OnchangeAccount(),
			Help:     "Account that will be set on invoice tax lines for invoices. Leave empty to use the expense account."},
		"RefundAccount": models.Many2OneField{String: "Tax Account on Refunds",
			RelationModel: pool.AccountAccount(), Filter: pool.AccountAccount().Deprecated().Equals(false),
			OnDelete: models.Restrict,
			Help:     "Account that will be set on invoice tax lines for refunds. Leave empty to use the expense account."},
		"Description": models.CharField{String: "Label on Invoices", Translate: true},
		"PriceInclude": models.BooleanField{String: "Included in Price", Default: models.DefaultValue(false),
			OnChange: pool.AccountTax().Methods().OnchangePriceInclude(),
			Help:     "Check this if the price you use on the product and invoices includes this tax."},
		"IncludeBaseAmount": models.BooleanField{String: "Affect Base of Subsequent Taxes",
			Default: models.DefaultValue(false),
			Help:    "If set, taxes which are computed after this one will be computed based on the price tax included."},
		"Analytic": models.BooleanField{String: "Include in Analytic Cost",
			Help: `If set, the amount computed by this tax will be assigned
to the same analytic account as the invoice line (if any)`},
		"Tags": models.Many2ManyField{String: "Tags", RelationModel: pool.AccountAccountTag(), JSON: "tag_ids",
			Help: "Optional tags you may want to assign for custom reporting"},
		"TaxGroup": models.Many2OneField{RelationModel: pool.AccountTaxGroup(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.AccountTaxGroup().NewSet(env).SearchAll().Limit(1)
			}, Required: true},
	})

	pool.AccountTax().AddSQLConstraint("name_company_uniq", "unique(name, company_id, type_tax_use)",
		"Tax names must be unique !")

	pool.AccountTax().Methods().Unlink().Extend("",
		func(rs pool.AccountTaxSet) int64 {
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

	pool.AccountTax().Methods().CheckChildrenScope().DeclareMethod(
		`CheckChildrenScope`,
		func(rs pool.AccountTaxSet) {
			//@api.constrains('children_tax_ids','type_tax_use')
			/*def _check_children_scope(self):
			  if not all(child.type_tax_use in ('none', self.type_tax_use) for child in self.children_tax_ids):
			      raise ValidationError(_('The application scope of taxes in a group must be either the same as the group or "None".'))

			*/
		})

	pool.AccountTax().Methods().Copy().Extend("",
		func(rs pool.AccountTaxSet, overrides *pool.AccountTaxData, fieldsToReset ...models.FieldNamer) pool.AccountTaxSet {
			//@api.returns('self',lambdavalue:value.id)
			/*def copy(self, default=None):
			default = dict(default or {}, name=_("%s (Copy)") % self.name)
			return super(AccountTax, self).copy(default=default)

			*/
			return rs.Super().Copy(overrides, fieldsToReset...)
		})

	pool.AccountTax().Methods().SearchByName().Extend("",
		func(rs pool.AccountTaxSet, name string, op operator.Operator, additionalCond pool.AccountTaxCondition, limit int) pool.AccountTaxSet {
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

	pool.AccountTax().Methods().Search().Extend("",
		func(rs pool.AccountTaxSet, cond pool.AccountTaxCondition) pool.AccountTaxSet {
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

	pool.AccountTax().Methods().OnchangeAmount().DeclareMethod(
		`OnchangeAmount`,
		func(rs pool.AccountTaxSet) (*pool.AccountTaxData, []models.FieldNamer) {
			//@api.onchange('amount')
			/*def onchange_amount(self):
			  if self.amount_type in ('percent', 'division') and self.amount != 0.0 and not self.description:
			      self.description = "{0:.4g}%".format(self.amount)

			*/
			return new(pool.AccountTaxData), []models.FieldNamer{}
		})

	pool.AccountTax().Methods().OnchangeAccount().DeclareMethod(
		`OnchangeAccount`,
		func(rs pool.AccountTaxSet) (*pool.AccountTaxData, []models.FieldNamer) {
			//@api.onchange('account_id')
			/*def onchange_account_id(self):
			  self.refund_account_id = self.account_id

			*/
			return new(pool.AccountTaxData), []models.FieldNamer{}
		})

	pool.AccountTax().Methods().OnchangePriceInclude().DeclareMethod(
		`OnchangePriceInclude`,
		func(rs pool.AccountTaxSet) (*pool.AccountTaxData, []models.FieldNamer) {
			//@api.onchange('price_include')
			/*def onchange_price_include(self):
			  if self.price_include:
			      self.include_base_amount = True

			*/
			return new(pool.AccountTaxData), []models.FieldNamer{}
		})

	pool.AccountTax().Methods().GetGroupingKey().DeclareMethod(
		`GetGroupingKey`,
		func(rs pool.AccountTaxSet, invoiceTaxVal *pool.AccountInvoiceTaxData) string {
			/*def get_grouping_key(self, invoice_tax_val):
			  """ Returns a string that will be used to group account.invoice.tax sharing the same properties"""
			  self.ensure_one()
			  return str(invoice_tax_val['tax_id']) + '-' + str(invoice_tax_val['account_id']) + '-' + str(invoice_tax_val['account_analytic_id'])

			*/
			return ""
		})

	pool.AccountTax().Methods().ComputeAmount().DeclareMethod(
		`ComputeAmount`,
		func(rs pool.AccountTaxSet, baseAmount, priceUnit, quantity float64, product pool.ProductProductSet, partner pool.PartnerSet) float64 {
			/*def _compute_amount(self, base_amount, price_unit, quantity=1.0, product=None, partner=None):
			  """ Returns the amount of a single tax. base_amount is the actual amount on which the tax is applied, which is
			      price_unit * quantity eventually affected by previous taxes (if tax is include_base_amount XOR price_include)
			  """
			  self.ensure_one()
			  if self.amount_type == 'fixed':
			      # Use copysign to take into account the sign of the base amount which includes the sign
			      # of the quantity and the sign of the price_unit
			      # Amount is the fixed price for the tax, it can be negative
			      # Base amount included the sign of the quantity and the sign of the unit price and when
			      # a product is returned, it can be done either by changing the sign of quantity or by changing the
			      # sign of the price unit.
			      # When the price unit is equal to 0, the sign of the quantity is absorbed in base_amount then
			      # a "else" case is needed.
			      if base_amount:
			          return math.copysign(quantity, base_amount) * self.amount
			      else:
			          return quantity * self.amount
			  if (self.amount_type == 'percent' and not self.price_include) or (self.amount_type == 'division' and self.price_include):
			      return base_amount * self.amount / 100
			  if self.amount_type == 'percent' and self.price_include:
			      return base_amount - (base_amount / (1 + self.amount / 100))
			  if self.amount_type == 'division' and not self.price_include:
			      return base_amount / (1 - self.amount / 100) - base_amount

			*/
			return 0
		})

	pool.AccountTax().Methods().JSONFriendlyComputeAll().DeclareMethod(
		`JSONFriendlyComputeAll`,
		func(rs pool.AccountTaxSet, priceUnit float64, currencyID int64, quantity float64, productID int64, partnerID int64) float64 {
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

	pool.AccountTax().Methods().ComputeAll().DeclareMethod(
		`ComputeAll`,
		func(rs pool.AccountTaxSet, priceUnit float64, currency pool.CurrencySet, quantity float64, product pool.ProductProductSet, partner pool.PartnerSet) (float64, float64, pool.AccountTaxSet) {
			//@api.multi
			/*def compute_all(self, price_unit, currency=None, quantity=1.0, product=None, partner=None):
			  """ Returns all information required to apply taxes (in self + their children in case of a tax goup).
			      We consider the sequence of the parent for group of taxes.
			          Eg. considering letters as taxes and alphabetic order as sequence :
			          [G, B([A, D, F]), E, C] will be computed as [A, D, F, C, E, G]

			  RETURN: {
			      'total_excluded': 0.0,    # Total without taxes
			      'total_included': 0.0,    # Total with taxes
			      'taxes': [{               # One dict for each tax in self and their children
			          'id': int,
			          'name': str,
			          'amount': float,
			          'sequence': int,
			          'account_id': int,
			          'refund_account_id': int,
			          'analytic': boolean,
			      }]
			  } """
			  if len(self) == 0:
			      company_id = self.env.user.company_id
			  else:
			      company_id = self[0].company_id
			  if not currency:
			      currency = company_id.currency_id
			  taxes = []
			  # By default, for each tax, tax amount will first be computed
			  # and rounded at the 'Account' decimal precision for each
			  # PO/SO/invoice line and then these rounded amounts will be
			  # summed, leading to the total amount for that tax. But, if the
			  # company has tax_calculation_rounding_method = round_globally,
			  # we still follow the same method, but we use a much larger
			  # precision when we round the tax amount for each line (we use
			  # the 'Account' decimal precision + 5), and that way it's like
			  # rounding after the sum of the tax amounts of each line
			  prec = currency.decimal_places

			  # In some cases, it is necessary to force/prevent the rounding of the tax and the total
			  # amounts. For example, in SO/PO line, we don't want to round the price unit at the
			  # precision of the currency.
			  # The context key 'round' allows to force the standard behavior.
			  round_tax = False if company_id.tax_calculation_rounding_method == 'round_globally' else True
			  round_total = True
			  if 'round' in self.env.context:
			      round_tax = bool(self.env.context['round'])
			      round_total = bool(self.env.context['round'])

			  if not round_tax:
			      prec += 5

			  base_values = self.env.context.get('base_values')
			  if not base_values:
			      total_excluded = total_included = base = round(price_unit * quantity, prec)
			  else:
			      total_excluded, total_included, base = base_values

			  # Sorting key is mandatory in this case. When no key is provided, sorted() will perform a
			  # search. However, the search method is overridden in account.tax in order to add a domain
			  # depending on the context. This domain might filter out some taxes from self, e.g. in the
			  # case of group taxes.
			  for tax in self.sorted(key=lambda r: r.sequence):
			      if tax.amount_type == 'group':
			          children = tax.children_tax_ids.with_context(base_values=(total_excluded, total_included, base))
			          ret = children.compute_all(price_unit, currency, quantity, product, partner)
			          total_excluded = ret['total_excluded']
			          base = ret['base'] if tax.include_base_amount else base
			          total_included = ret['total_included']
			          tax_amount = total_included - total_excluded
			          taxes += ret['taxes']
			          continue

			      tax_amount = tax._compute_amount(base, price_unit, quantity, product, partner)
			      if not round_tax:
			          tax_amount = round(tax_amount, prec)
			      else:
			          tax_amount = currency.round(tax_amount)

			      if tax.price_include:
			          total_excluded -= tax_amount
			          base -= tax_amount
			      else:
			          total_included += tax_amount

			      # Keep base amount used for the current tax
			      tax_base = base

			      if tax.include_base_amount:
			          base += tax_amount

			      taxes.append({
			          'id': tax.id,
			          'name': tax.with_context(**{'lang': partner.lang} if partner else {}).name,
			          'amount': tax_amount,
			          'base': tax_base,
			          'sequence': tax.sequence,
			          'account_id': tax.account_id.id,
			          'refund_account_id': tax.refund_account_id.id,
			          'analytic': tax.analytic,
			      })

			  return {
			      'taxes': sorted(taxes, key=lambda k: k['sequence']),
			      'total_excluded': currency.round(total_excluded) if round_total else total_excluded,
			      'total_included': currency.round(total_included) if round_total else total_included,
			      'base': base,
			  }

			*/
			return 0, 0, pool.AccountTax().NewSet(rs.Env())
		})

	pool.AccountTax().Methods().FixTaxIncludedPrice().DeclareMethod(
		`FixTaxIncludedPrice`,
		func(rs pool.AccountTaxSet, price float64, prodTaxes, lineTaxes pool.AccountTaxSet) float64 {
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

	pool.AccountReconcileModel().DeclareModel()

	pool.AccountReconcileModel().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Button Label", Required: true,
			OnChange: pool.AccountReconcileModel().Methods().OnchangeName()},
		"Sequence":      models.IntegerField{Required: true, Default: models.DefaultValue(10)},
		"HasSecondLine": models.BooleanField{String: "Add a second line", Default: models.DefaultValue(false)},
		"Company": models.Many2OneField{RelationModel: pool.Company(), Required: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser().Company()
			}},
		"Account": models.Many2OneField{RelationModel: pool.AccountAccount(),
			OnDelete: models.Cascade,
			Filter:   pool.AccountAccount().Deprecated().Equals(false)},
		"Journal": models.Many2OneField{RelationModel: pool.AccountJournal(),
			OnDelete: models.Cascade, Help: "This field is ignored in a bank statement reconciliation."},
		"Label": models.CharField{String: "Journal Item Label"},
		"AmountType": models.SelectionField{Selection: types.Selection{
			"fixed":      "Fixed",
			"percentage": "Percentage of balance"}, Required: true, Default: models.DefaultValue("percentage")},
		"Amount": models.FloatField{Required: true, Default: models.DefaultValue(100.0),
			Help: "Fixed amount will count as a debit if it is negative, as a credit if it is positive."},
		"Tax": models.Many2OneField{String: "Tax", RelationModel: pool.AccountTax(),
			OnDelete: models.Restrict, Filter: pool.AccountTax().TypeTaxUse().Equals("purchase")},
		"AnalyticAccount": models.Many2OneField{RelationModel: pool.AccountAnalyticAccount(),
			OnDelete: models.SetNull},
		"SecondAccount": models.Many2OneField{RelationModel: pool.AccountAccount(),
			OnDelete: models.Cascade, Filter: pool.AccountAccount().Deprecated().Equals(false),
		},
		"SecondJournal": models.Many2OneField{RelationModel: pool.AccountJournal(),
			OnDelete: models.Cascade, Help: "This field is ignored in a bank statement reconciliation."},
		"SecondLabel": models.CharField{String: "Second Journal Item Label"},
		"SecondAmountType": models.SelectionField{Selection: types.Selection{
			"fixed":      "Fixed",
			"percentage": "Percentage of balance"}, Required: true, Default: models.DefaultValue("percentage")},
		"SecondAmount": models.FloatField{Required: true, Default: models.DefaultValue(100.0),
			Help: "Fixed amount will count as a debit if it is negative, as a credit if it is positive."},
		"SecondTax": models.Many2OneField{RelationModel: pool.AccountTax(),
			OnDelete: models.Restrict, Filter: pool.AccountTax().TypeTaxUse().Equals("purchase")},
		"SecondAnalyticAccount": models.Many2OneField{RelationModel: pool.AccountAnalyticAccount(),
			OnDelete: models.SetNull},
	})

	pool.AccountReconcileModel().Methods().OnchangeName().DeclareMethod(
		`OnchangeName`,
		func(rs pool.AccountReconcileModelSet) (*pool.AccountReconcileModelData, []models.FieldNamer) {
			//@api.onchange('name')
			/*def onchange_name(self):
			  self.label = self.name
			*/
			return new(pool.AccountReconcileModelData), []models.FieldNamer{}
		})

}
