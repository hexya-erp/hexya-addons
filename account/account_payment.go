// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountPaymentMethod().DeclareModel()
	pool.AccountPaymentMethod().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{Required: true, Translate: true},
		"Code": models.CharField{Required: true},
		"PaymentType": models.SelectionField{Selection: types.Selection{
			"inbound":  "Inbound",
			"outbound": "Outbound",
		}, Required: true},
	})

	pool.AccountAbstractPayment().DeclareMixinModel()
	pool.AccountAbstractPayment().AddFields(map[string]models.FieldDefinition{
		"PaymentType": models.SelectionField{Selection: types.Selection{
			"outbound": "Send Money",
			"inbound":  "Receive Money",
		}, Required: true},
		"PaymentMethod": models.Many2OneField{String: "Payment Method Type",
			RelationModel: pool.AccountPaymentMethod(), Required: true},
		"PaymentMethodCode": models.CharField{
			Help: "Technical field used to adapt the interface to the payment type selected." /*[ readonly True]*/},
		"PartnerType": models.SelectionField{Selection: types.Selection{
			"customer": "Customer",
			"supplier": "Vendor",
		}},
		"Partner": models.Many2OneField{RelationModel: pool.Partner()},
		"Amount": models.FloatField{String: "Payment Amount", Required: true,
			Constraint: pool.AccountAbstractPayment().Methods().CheckAmount()},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(), Required: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser().Company().Currency()
			}},
		"PaymentDate": models.DateField{Default: func(env models.Environment, vals models.FieldMap) interface{} {
			return dates.Today()
		}, Required: true, NoCopy: true},
		"Communication": models.CharField{String: "Memo"},
		"Journal": models.Many2OneField{String: "Payment Journal", RelationModel: pool.AccountJournal(),
			Required: true, Filter: pool.AccountJournal().Type().In([]string{"bank", "cash"}),
			OnChange: pool.AccountAbstractPayment().Methods().OnchangeJournal()},
		"Company": models.Many2OneField{RelationModel: pool.Company(), Related: "Journal.Company" /* readonly=true */},
		"HidePaymentMethod": models.BooleanField{
			Compute: pool.AccountAbstractPayment().Methods().ComputeHidePaymentMethod(),
			Help: `Technical field used to hide the payment method if the selected journal
has only one available which is 'manual'`},
	})

	pool.AccountAbstractPayment().Methods().CheckAmount().DeclareMethod(
		`CheckAmount`,
		func(rs pool.AccountAbstractPaymentSet) {
			//@api.constrains('amount')
			/*def _check_amount(self):
			  if not self.amount > 0.0:
			      raise ValidationError(_('The payment amount must be strictly positive.'))

			*/
		})

	pool.AccountAbstractPayment().Methods().ComputeHidePaymentMethod().DeclareMethod(
		`ComputeHidePaymentMethod`,
		func(rs pool.AccountAbstractPaymentSet) (*pool.AccountAbstractPaymentData, []models.FieldNamer) {
			//@api.depends('payment_type','journal_id')
			/*def _compute_hide_payment_method(self):
			  if not self.journal_id:
			      self.hide_payment_method = True
			      return
			  journal_payment_methods = self.payment_type == 'inbound' and self.journal_id.inbound_payment_method_ids or self.journal_id.outbound_payment_method_ids
			  self.hide_payment_method = len(journal_payment_methods) == 1 and journal_payment_methods[0].code == 'manual'

			*/
			return new(pool.AccountAbstractPaymentData), []models.FieldNamer{}
		})

	pool.AccountAbstractPayment().Methods().OnchangeJournal().DeclareMethod(
		`OnchangeJournal`,
		func(rs pool.AccountAbstractPaymentSet) (*pool.AccountAbstractPaymentData, []models.FieldNamer) {
			//@api.onchange('journal_id')
			/*def _onchange_journal(self):
			  if self.journal_id:
			      self.currency_id = self.journal_id.currency_id or self.company_id.currency_id
			      # Set default payment method (we consider the first to be the default one)
			      payment_methods = self.payment_type == 'inbound' and self.journal_id.inbound_payment_method_ids or self.journal_id.outbound_payment_method_ids
			      self.payment_method_id = payment_methods and payment_methods[0] or False
			      # Set payment method domain (restrict to methods enabled for the journal and to selected payment type)
			      payment_type = self.payment_type in ('outbound', 'transfer') and 'outbound' or 'inbound'
			      return {'domain': {'payment_method_id': [('payment_type', '=', payment_type), ('id', 'in', payment_methods.ids)]}}
			  return {}

			*/
			return new(pool.AccountAbstractPaymentData), []models.FieldNamer{}
		})

	pool.AccountAbstractPayment().Methods().GetInvoices().DeclareMethod(
		`GetInvoices`,
		func(rs pool.AccountAbstractPaymentSet) pool.AccountInvoiceSet {
			/*def _get_invoices(self):
			  """ Return the invoices of the payment. Must be overridden """
			  raise NotImplementedError

			*/
			panic("Not implemented")
		})

	pool.AccountAbstractPayment().Methods().ComputeTotalInvoicesAmount().DeclareMethod(
		`ComputeTotalInvoicesAmount`,
		func(rs pool.AccountAbstractPaymentSet) float64 {
			/*def _compute_total_invoices_amount(self):
			  """ Compute the sum of the residual of invoices, expressed in the payment currency """
			  payment_currency = self.currency_id or self.journal_id.currency_id or self.journal_id.company_id.currency_id or self.env.user.company_id.currency_id
			  invoices = self._get_invoices()

			  if all(inv.currency_id == payment_currency for inv in invoices):
			      total = sum(invoices.mapped('residual_signed'))
			  else:
			      total = 0
			      for inv in invoices:
			          if inv.company_currency_id != payment_currency:
			              total += inv.company_currency_id.with_context(date=self.payment_date).compute(inv.residual_company_signed, payment_currency)
			          else:
			              total += inv.residual_company_signed
			  return abs(total)


			*/
			return 0
		})

	pool.AccountRegisterPayments().DeclareTransientModel()
	pool.AccountRegisterPayments().InheritModel(pool.AccountAbstractPayment())

	pool.AccountRegisterPayments().Fields().PaymentType().SetOnchange(pool.AccountRegisterPayments().Methods().OnchangePaymentType())

	pool.AccountRegisterPayments().Methods().OnchangePaymentType().DeclareMethod(
		`OnchangePaymentType`,
		func(rs pool.AccountRegisterPaymentsSet) (*pool.AccountRegisterPaymentsData, []models.FieldNamer) {
			//@api.onchange('payment_type')
			/*def _onchange_payment_type(self):
			  if self.payment_type:
			      return {'domain': {'payment_method_id': [('payment_type', '=', self.payment_type)]}}

			*/
			return new(pool.AccountRegisterPaymentsData), []models.FieldNamer{}
		})

	pool.AccountRegisterPayments().Methods().GetInvoices().Extend("",
		func(rs pool.AccountRegisterPaymentsSet) pool.AccountInvoiceSet {
			/*def _get_invoices(self):
			  """ Return the invoices of the payment. Must be overridden """
				return self.env['account.invoice'].browse(self._context.get('active_ids'))
			*/
			if rs.Env().Context().HasKey("active_ids") {
				return pool.AccountInvoice().Browse(rs.Env(), rs.Env().Context().GetIntegerSlice("active_ids"))
			}
			return pool.AccountInvoice().NewSet(rs.Env())
		})

	pool.AccountRegisterPayments().Methods().DefaultGet().Extend("",
		func(rs pool.AccountRegisterPaymentsSet) models.FieldMap {
			//@api.model
			/*def default_get(self, fields):
			  rec = super(account_register_payments, self).default_get(fields)
			  context = dict(self._context or {})
			  active_model = context.get('active_model')
			  active_ids = context.get('active_ids')

			  # Checks on context parameters
			  if not active_model or not active_ids:
			      raise UserError(_("Programmation error: wizard action executed without active_model or active_ids in context."))
			  if active_model != 'account.invoice':
			      raise UserError(_("Programmation error: the expected model for this action is 'account.invoice'. The provided one is '%d'.") % active_model)

			  # Checks on received invoice records
			  invoices = self.env[active_model].browse(active_ids)
			  if any(invoice.state != 'open' for invoice in invoices):
			      raise UserError(_("You can only register payments for open invoices"))
			  if any(inv.commercial_partner_id != invoices[0].commercial_partner_id for inv in invoices):
			      raise UserError(_("In order to pay multiple invoices at once, they must belong to the same commercial partner."))
			  if any(MAP_INVOICE_TYPE_PARTNER_TYPE[inv.type] != MAP_INVOICE_TYPE_PARTNER_TYPE[invoices[0].type] for inv in invoices):
			      raise UserError(_("You cannot mix customer invoices and vendor bills in a single payment."))
			  if any(inv.currency_id != invoices[0].currency_id for inv in invoices):
			      raise UserError(_("In order to pay multiple invoices at once, they must use the same currency."))

			  total_amount = sum(inv.residual * MAP_INVOICE_TYPE_PAYMENT_SIGN[inv.type] for inv in invoices)
			  communication = ' '.join([ref for ref in invoices.mapped('reference') if ref])

			  rec.update({
			      'amount': abs(total_amount),
			      'currency_id': invoices[0].currency_id.id,
			      'payment_type': total_amount > 0 and 'inbound' or 'outbound',
			      'partner_id': invoices[0].commercial_partner_id.id,
			      'partner_type': MAP_INVOICE_TYPE_PARTNER_TYPE[invoices[0].type],
			      'communication': communication,
			  })
			  return rec

			*/
			return rs.Super().DefaultGet()
		})

	pool.AccountRegisterPayments().Methods().GetPaymentVals().DeclareMethod(
		`GetPaymentVals`,
		func(rs pool.AccountRegisterPaymentsSet) *pool.AccountPaymentData {
			/*def get_payment_vals(self):
			  """ Hook for extension """
			  return {
			      'journal_id': self.journal_id.id,
			      'payment_method_id': self.payment_method_id.id,
			      'payment_date': self.payment_date,
			      'communication': self.communication,
			      'invoice_ids': [(4, inv.id, None) for inv in self._get_invoices()],
			      'payment_type': self.payment_type,
			      'amount': self.amount,
			      'currency_id': self.currency_id.id,
			      'partner_id': self.partner_id.id,
			      'partner_type': self.partner_type,
			  }

			*/
			return new(pool.AccountPaymentData)
		})

	pool.AccountRegisterPayments().Methods().CreatePayment().DeclareMethod(
		`CreatePayment`,
		func(rs pool.AccountRegisterPaymentsSet) *actions.Action {
			payment := pool.AccountPayment().Create(rs.Env(), rs.GetPaymentVals())
			payment.Post()
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

	pool.AccountPayment().DeclareModel()
	pool.AccountPayment().InheritModel(pool.AccountAbstractPayment())
	pool.AccountPayment().SetDefaultOrder("PaymentDate DESC", "Name DESC")

	pool.AccountPayment().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Name" /*[readonly True]*/, NoCopy: true,
			Default: models.DefaultValue("Draft Payment") /* The name is attributed upon post(]*/},
		"State": models.SelectionField{String: "Status", Selection: types.Selection{
			"draft":      "Draft",
			"posted":     "Posted",
			"sent":       "Sent",
			"reconciled": "Reconciled",
		}, /*[readonly True]*/ Default: models.DefaultValue("draft"), NoCopy: true},
		"PaymentReference": models.CharField{String: "PaymentReference", NoCopy: true, /*[ readonly True]*/
			Help: "Reference of the document used to issue this payment. Eg. check number, file name, etc."},
		"MoveName": models.CharField{String: "Journal Entry Name", /*[ readonly True]*/
			Default: models.DefaultValue(false), NoCopy: true,
			Help: `Technical field holding the number given to the journal entry, automatically set when the statement
line is reconciled then stored to set the same number again if the line is cancelled,
set to draft and re-processed again." `},
		"DestinationAccount": models.Many2OneField{RelationModel: pool.AccountAccount(),
			Compute: pool.AccountPayment().Methods().ComputeDestinationAccount() /* readonly=true */},
		"DestinationJournal": models.Many2OneField{String: "Transfer To", RelationModel: pool.AccountJournal(),
			Filter: pool.AccountJournal().Type().In([]string{"bank", "cash"})},
		"Invoices": models.Many2ManyField{RelationModel: pool.AccountInvoice(), JSON: "invoice_ids",
			NoCopy: true /*[ readonly True]*/},
		"HasInvoices": models.BooleanField{Compute: pool.AccountPayment().Methods().ComputeHasInvoices(),
			Help: "Technical field used for usability purposes"},
		"PaymentDifference": models.FloatField{Compute: pool.AccountPayment().Methods().ComputePaymentDifference()},
		"PaymentDifferenceHandling": models.SelectionField{String: "Payment Difference", Selection: types.Selection{
			"open":      "Keep open",
			"reconcile": "Mark invoice as fully paid",
		}, Default: models.DefaultValue("open"), NoCopy: true},
		"WriteoffAccount": models.Many2OneField{String: "Difference Account", RelationModel: pool.AccountAccount(),
			Filter: pool.AccountAccount().Deprecated().Equals(false)},
		"MoveLines": models.One2ManyField{RelationModel: pool.AccountMoveLine(), ReverseFK: "Payment",
			JSON: "move_line_ids" /* readonly */, NoCopy: true},
	})

	pool.AccountPayment().Fields().PaymentType().
		UpdateSelection(types.Selection{"transfer": "Internal Transfer"}).
		SetOnchange(pool.AccountPayment().Methods().OnchangePaymentType())

	pool.AccountPayment().Fields().PartnerType().SetOnchange(pool.AccountPayment().Methods().OnchangePartnerType())

	pool.AccountPayment().Methods().ComputeHasInvoices().DeclareMethod(
		`ComputeHasInvoices`,
		func(rs pool.AccountPaymentSet) (*pool.AccountPaymentData, []models.FieldNamer) {
			//@api.depends('invoice_ids')
			/*def _get_has_invoices(self):
			  self.has_invoices = bool(self.invoice_ids)

			*/
			return &pool.AccountPaymentData{}, []models.FieldNamer{}
		})

	pool.AccountPayment().Methods().ComputePaymentDifference().DeclareMethod(
		`ComputePaymentDifference`,
		func(rs pool.AccountPaymentSet) (*pool.AccountPaymentData, []models.FieldNamer) {
			//@api.depends('invoice_ids','amount','payment_date','currency_id')
			/*def _compute_payment_difference(self):
			  if len(self.invoice_ids) == 0:
			      return
			  if self.invoice_ids[0].type in ['in_invoice', 'out_refund']:
			      self.payment_difference = self.amount - self._compute_total_invoices_amount()
			  else:
			      self.payment_difference = self._compute_total_invoices_amount() - self.amount
			*/
			return &pool.AccountPaymentData{}, []models.FieldNamer{}

		})

	pool.AccountPayment().Methods().ComputeDestinationAccount().DeclareMethod(
		`ComputeDestinationAccountId`,
		func(rs pool.AccountPaymentSet) (*pool.AccountPaymentData, []models.FieldNamer) {
			//@api.depends('invoice_ids','payment_type','partner_type','partner_id')
			/*def _compute_destination_account_id(self):
			  if self.invoice_ids:
			      self.destination_account_id = self.invoice_ids[0].account_id.id
			  elif self.payment_type == 'transfer':
			      if not self.company_id.transfer_account_id.id:
			          raise UserError(_('Transfer account not defined on the company.'))
			      self.destination_account_id = self.company_id.transfer_account_id.id
			  elif self.partner_id:
			      if self.partner_type == 'customer':
			          self.destination_account_id = self.partner_id.property_account_receivable_id.id
			      else:
			          self.destination_account_id = self.partner_id.property_account_payable_id.id

			*/
			return &pool.AccountPaymentData{}, []models.FieldNamer{}
		})

	pool.AccountPayment().Methods().OnchangePartnerType().DeclareMethod(
		`OnchangePartnerType`,
		func(rs pool.AccountPaymentSet) (*pool.AccountPaymentData, []models.FieldNamer) {
			//@api.onchange('partner_type')
			/*def _onchange_partner_type(self):
			  # Set partner_id domain
			  if self.partner_type:
			      return {'domain': {'partner_id': [(self.partner_type, '=', True)]}}

			*/
			return &pool.AccountPaymentData{}, []models.FieldNamer{}
		})

	pool.AccountPayment().Methods().OnchangePaymentType().DeclareMethod(
		`OnchangePaymentType`,
		func(rs pool.AccountPaymentSet) (*pool.AccountPaymentData, []models.FieldNamer) {
			//@api.onchange('payment_type')
			/*def _onchange_payment_type(self):
			    # Set partner_id domain
				if self.partner_type:
					return {'domain': {'partner_id': [(self.partner_type, '=', True)]}}

			*/
			return &pool.AccountPaymentData{}, []models.FieldNamer{}
		})

	pool.AccountPayment().Methods().DefaultGet().Extend("",
		func(rs pool.AccountPaymentSet) models.FieldMap {
			//@api.model
			/*def default_get(self, fields):
			rec = super(account_payment, self).default_get(fields)
			invoice_defaults = self.resolve_2many_commands('invoice_ids', rec.get('invoice_ids'))
			if invoice_defaults and len(invoice_defaults) == 1:
				invoice = invoice_defaults[0]
				rec['communication'] = invoice['reference'] or invoice['name'] or invoice['number']
				rec['currency_id'] = invoice['currency_id'][0]
				rec['payment_type'] = invoice['type'] in ('out_invoice', 'in_refund') and 'inbound' or 'outbound'
				rec['partner_type'] = MAP_INVOICE_TYPE_PARTNER_TYPE[invoice['type']]
				rec['partner_id'] = invoice['partner_id'][0]
				rec['amount'] = invoice['residual']
			return rec


			*/
			return rs.Super().DefaultGet()
		})

	pool.AccountPayment().Methods().GetInvoices().Extend("",
		func(rs pool.AccountPaymentSet) pool.AccountInvoiceSet {
			return rs.Invoices()
		})

	pool.AccountPayment().Methods().ButtonJournalEntries().DeclareMethod(
		`ButtonJournalEntries`,
		func(rs pool.AccountPaymentSet) *actions.Action {
			//@api.multi
			/*def button_journal_entries(self):
			  return {
			      'name': _('Journal Items'),
			      'view_type': 'form',
			      'view_mode': 'tree,form',
			      'res_model': 'account.move.line',
			      'view_id': False,
			      'type': 'ir.actions.act_window',
			      'domain': [('payment_id', 'in', self.ids)],
			  }

			*/
			return new(actions.Action)
		})

	pool.AccountPayment().Methods().ButtonInvoices().DeclareMethod(
		`ButtonInvoices`,
		func(rs pool.AccountPaymentSet) *actions.Action {
			//@api.multi
			/*def button_invoices(self):
			  return {
			      'name': _('Paid Invoices'),
			      'view_type': 'form',
			      'view_mode': 'tree,form',
			      'res_model': 'account.invoice',
			      'view_id': False,
			      'type': 'ir.actions.act_window',
			      'domain': [('id', 'in', [x.id for x in self.invoice_ids])],
			  }

			*/
			return new(actions.Action)
		})

	pool.AccountPayment().Methods().ButtonDummy().DeclareMethod(
		`ButtonDummy`,
		func(rs pool.AccountPaymentSet) bool {
			return true
		})

	pool.AccountPayment().Methods().Unreconcile().DeclareMethod(
		`Unreconcile`,
		func(rs pool.AccountPaymentSet) {
			//@api.multi
			/*def unreconcile(self):
			  """ Set back the payments in 'posted' or 'sent' state, without deleting the journal entries.
			      Called when cancelling a bank statement line linked to a pre-registered payment.
			  """
			  for payment in self:
			      if payment.payment_reference:
			          payment.write({'state': 'sent'})
			      else:
			          payment.write({'state': 'posted'})

			*/
		})

	pool.AccountPayment().Methods().Cancel().DeclareMethod(
		`Cancel`,
		func(rs pool.AccountPaymentSet) {
			//@api.multi
			/*def cancel(self):
			  for rec in self:
			      for move in rec.move_line_ids.mapped('move_id'):
			          if rec.invoice_ids:
			              move.line_ids.remove_move_reconcile()
			          move.button_cancel()
			          move.unlink()
			      rec.state = 'draft'

			*/
		})

	pool.AccountPayment().Methods().Unlink().Extend("",
		func(rs pool.AccountPaymentSet) int64 {
			//@api.multi
			/*def unlink(self):
			  if any(bool(rec.move_line_ids) for rec in self):
			      raise UserError(_("You can not delete a payment that is already posted"))
			  if any(rec.move_name for rec in self):
			      raise UserError(_('It is not allowed to delete a payment that already created a journal entry since it would create a gap in the numbering. You should create the journal entry again and cancel it thanks to a regular revert.'))
			  return super(account_payment, self).unlink()

			*/
			return rs.Super().Unlink()
		})

	pool.AccountPayment().Methods().Post().DeclareMethod(
		`Post`,
		func(rs pool.AccountPaymentSet) {
			//@api.multi
			/*def post(self):
			  """ Create the journal items for the payment and update the payment's state to 'posted'.
			      A journal entry is created containing an item in the source liquidity account (selected journal's default_debit or default_credit)
			      and another in the destination reconciliable account (see _compute_destination_account_id).
			      If invoice_ids is not empty, there will be one reconciliable move line per invoice to reconcile with.
			      If the payment is a transfer, a second journal entry is created in the destination journal to receive money from the transfer account.
			  """
			  for rec in self:

			      if rec.state != 'draft':
			          raise UserError(_("Only a draft payment can be posted. Trying to post a payment in state %s.") % rec.state)

			      if any(inv.state != 'open' for inv in rec.invoice_ids):
			          raise ValidationError(_("The payment cannot be processed because the invoice is not open!"))

			      # Use the right sequence to set the name
			      if rec.payment_type == 'transfer':
			          sequence_code = 'account.payment.transfer'
			      else:
			          if rec.partner_type == 'customer':
			              if rec.payment_type == 'inbound':
			                  sequence_code = 'account.payment.customer.invoice'
			              if rec.payment_type == 'outbound':
			                  sequence_code = 'account.payment.customer.refund'
			          if rec.partner_type == 'supplier':
			              if rec.payment_type == 'inbound':
			                  sequence_code = 'account.payment.supplier.refund'
			              if rec.payment_type == 'outbound':
			                  sequence_code = 'account.payment.supplier.invoice'
			      rec.name = self.env['ir.sequence'].with_context(ir_sequence_date=rec.payment_date).next_by_code(sequence_code)

			      # Create the journal entry
			      amount = rec.amount * (rec.payment_type in ('outbound', 'transfer') and 1 or -1)
			      move = rec._create_payment_entry(amount)

			      # In case of a transfer, the first journal entry created debited the source liquidity account and credited
			      # the transfer account. Now we debit the transfer account and credit the destination liquidity account.
			      if rec.payment_type == 'transfer':
			          transfer_credit_aml = move.line_ids.filtered(lambda r: r.account_id == rec.company_id.transfer_account_id)
			          transfer_debit_aml = rec._create_transfer_entry(amount)
			          (transfer_credit_aml + transfer_debit_aml).reconcile()

			      rec.write({'state': 'posted', 'move_name': move.name})

			*/
		})

	pool.AccountPayment().Methods().CreatePaymentEntry().DeclareMethod(
		`CreatePaymentEntry`,
		func(rs pool.AccountPaymentSet, amount float64) pool.AccountMoveSet {
			/*def _create_payment_entry(self, amount):
			  """ Create a journal entry corresponding to a payment, if the payment references invoice(s) they are reconciled.
			      Return the journal entry.
			  """
			  aml_obj = self.env['account.move.line'].with_context(check_move_validity=False)
			  invoice_currency = False
			  if self.invoice_ids and all([x.currency_id == self.invoice_ids[0].currency_id for x in self.invoice_ids]):
			      #if all the invoices selected share the same currency, record the paiement in that currency too
			      invoice_currency = self.invoice_ids[0].currency_id
			  debit, credit, amount_currency, currency_id = aml_obj.with_context(date=self.payment_date).compute_amount_fields(amount, self.currency_id, self.company_id.currency_id, invoice_currency)

			  move = self.env['account.move'].create(self._get_move_vals())

			  #Write line corresponding to invoice payment
			  counterpart_aml_dict = self._get_shared_move_line_vals(debit, credit, amount_currency, move.id, False)
			  counterpart_aml_dict.update(self._get_counterpart_move_line_vals(self.invoice_ids))
			  counterpart_aml_dict.update({'currency_id': currency_id})
			  counterpart_aml = aml_obj.create(counterpart_aml_dict)

			  #Reconcile with the invoices
			  if self.payment_difference_handling == 'reconcile' and self.payment_difference:
			      writeoff_line = self._get_shared_move_line_vals(0, 0, 0, move.id, False)
			      amount_currency_wo, currency_id = aml_obj.with_context(date=self.payment_date).compute_amount_fields(self.payment_difference, self.currency_id, self.company_id.currency_id, invoice_currency)[2:]
			      # the writeoff debit and credit must be computed from the invoice residual in company currency
			      # minus the payment amount in company currency, and not from the payment difference in the payment currency
			      # to avoid loss of precision during the currency rate computations. See revision 20935462a0cabeb45480ce70114ff2f4e91eaf79 for a detailed example.
			      total_residual_company_signed = sum(invoice.residual_company_signed for invoice in self.invoice_ids)
			      total_payment_company_signed = self.currency_id.with_context(date=self.payment_date).compute(self.amount, self.company_id.currency_id)
			      if self.invoice_ids[0].type in ['in_invoice', 'out_refund']:
			          amount_wo = total_payment_company_signed - total_residual_company_signed
			      else:
			          amount_wo = total_residual_company_signed - total_payment_company_signed
			      # Align the sign of the secondary currency writeoff amount with the sign of the writeoff
			      # amount in the company currency
			      if amount_wo > 0:
			          debit_wo = amount_wo
			          credit_wo = 0.0
			          amount_currency_wo = abs(amount_currency_wo)
			      else:
			          debit_wo = 0.0
			          credit_wo = -amount_wo
			          amount_currency_wo = -abs(amount_currency_wo)
			      writeoff_line['name'] = _('Counterpart')
			      writeoff_line['account_id'] = self.writeoff_account_id.id
			      writeoff_line['debit'] = debit_wo
			      writeoff_line['credit'] = credit_wo
			      writeoff_line['amount_currency'] = amount_currency_wo
			      writeoff_line['currency_id'] = currency_id
			      writeoff_line = aml_obj.create(writeoff_line)
			      if counterpart_aml['debit']:
			          counterpart_aml['debit'] += credit_wo - debit_wo
			      if counterpart_aml['credit']:
			          counterpart_aml['credit'] += debit_wo - credit_wo
			      counterpart_aml['amount_currency'] -= amount_currency_wo
			  self.invoice_ids.register_payment(counterpart_aml)

			  #Write counterpart lines
			  if not self.currency_id != self.company_id.currency_id:
			      amount_currency = 0
			  liquidity_aml_dict = self._get_shared_move_line_vals(credit, debit, -amount_currency, move.id, False)
			  liquidity_aml_dict.update(self._get_liquidity_move_line_vals(-amount))
			  aml_obj.create(liquidity_aml_dict)

			  move.post()
			  return move

			*/
			return pool.AccountMove().NewSet(rs.Env())
		})

	pool.AccountPayment().Methods().CreateTransferEntry().DeclareMethod(
		`CreateTransferEntry`,
		func(rs pool.AccountPaymentSet, amount float64) pool.AccountMoveLineSet {
			/*def _create_transfer_entry(self, amount):
			  """ Create the journal entry corresponding to the 'incoming money' part of an internal transfer, return the reconciliable move line
			  """
			  aml_obj = self.env['account.move.line'].with_context(check_move_validity=False)
			  debit, credit, amount_currency, dummy = aml_obj.with_context(date=self.payment_date).compute_amount_fields(amount, self.currency_id, self.company_id.currency_id)
			  amount_currency = self.destination_journal_id.currency_id and self.currency_id.with_context(date=self.payment_date).compute(amount, self.destination_journal_id.currency_id) or 0

			  dst_move = self.env['account.move'].create(self._get_move_vals(self.destination_journal_id))

			  dst_liquidity_aml_dict = self._get_shared_move_line_vals(debit, credit, amount_currency, dst_move.id)
			  dst_liquidity_aml_dict.update({
			      'name': _('Transfer from %s') % self.journal_id.name,
			      'account_id': self.destination_journal_id.default_credit_account_id.id,
			      'currency_id': self.destination_journal_id.currency_id.id,
			      'payment_id': self.id,
			      'journal_id': self.destination_journal_id.id})
			  aml_obj.create(dst_liquidity_aml_dict)

			  transfer_debit_aml_dict = self._get_shared_move_line_vals(credit, debit, 0, dst_move.id)
			  transfer_debit_aml_dict.update({
			      'name': self.name,
			      'payment_id': self.id,
			      'account_id': self.company_id.transfer_account_id.id,
			      'journal_id': self.destination_journal_id.id})
			  if self.currency_id != self.company_id.currency_id:
			      transfer_debit_aml_dict.update({
			          'currency_id': self.currency_id.id,
			          'amount_currency': -self.amount,
			      })
			  transfer_debit_aml = aml_obj.create(transfer_debit_aml_dict)
			  dst_move.post()
			  return transfer_debit_aml

			*/
			return pool.AccountMoveLine().NewSet(rs.Env())
		})

	pool.AccountPayment().Methods().GetMoveVals().DeclareMethod(
		`GetMoveVals`,
		func(rs pool.AccountPaymentSet, journal pool.AccountJournalSet) *pool.AccountMoveData {
			/*def _get_move_vals(self, journal=None):
			  """ Return dict to create the payment move
			  """
			  journal = journal or self.journal_id
			  if not journal.sequence_id:
			      raise UserError(_('Configuration Error !'), _('The journal %s does not have a sequence, please specify one.') % journal.name)
			  if not journal.sequence_id.active:
			      raise UserError(_('Configuration Error !'), _('The sequence of journal %s is deactivated.') % journal.name)
			  name = self.move_name or journal.with_context(ir_sequence_date=self.payment_date).sequence_id.next_by_id()
			  return {
			      'name': name,
			      'date': self.payment_date,
			      'ref': self.communication or '',
			      'company_id': self.company_id.id,
			      'journal_id': journal.id,
			  }

			*/
			return &pool.AccountMoveData{}
		})

	pool.AccountPayment().Methods().GetSharedMoveLineVals().DeclareMethod(
		`GetSharedMoveLineVals`,
		func(rs pool.AccountPaymentSet, debit, credit, amountCurrency float64, move pool.AccountMoveSet,
			invoice pool.AccountInvoiceSet) *pool.AccountMoveLineData {
			/*def _get_shared_move_line_vals(self, debit, credit, amount_currency, move_id, invoice_id=False):
			  """ Returns values common to both move lines (except for debit, credit and amount_currency which are reversed)
			  """
			  return {
			      'partner_id': self.payment_type in ('inbound', 'outbound') and self.env['res.partner']._find_accounting_partner(self.partner_id).id or False,
			      'invoice_id': invoice_id and invoice_id.id or False,
			      'move_id': move_id,
			      'debit': debit,
			      'credit': credit,
			      'amount_currency': amount_currency or False,
			  }

			*/
			return &pool.AccountMoveLineData{}
		})

	pool.AccountPayment().Methods().GetCounterpartMoveLineVals().DeclareMethod(
		`GetCounterpartMoveLineVals`,
		func(rs pool.AccountPaymentSet, invoice pool.AccountInvoiceSet) *pool.AccountMoveLineData {
			/*def _get_counterpart_move_line_vals(self, invoice=False):
			  if self.payment_type == 'transfer':
			      name = self.name
			  else:
			      name = ''
			      if self.partner_type == 'customer':
			          if self.payment_type == 'inbound':
			              name += _("Customer Payment")
			          elif self.payment_type == 'outbound':
			              name += _("Customer Refund")
			      elif self.partner_type == 'supplier':
			          if self.payment_type == 'inbound':
			              name += _("Vendor Refund")
			          elif self.payment_type == 'outbound':
			              name += _("Vendor Payment")
			      if invoice:
			          name += ': '
			          for inv in invoice:
			              if inv.move_id:
			                  name += inv.number + ', '
			          name = name[:len(name)-2]
			  return {
			      'name': name,
			      'account_id': self.destination_account_id.id,
			      'journal_id': self.journal_id.id,
			      'currency_id': self.currency_id != self.company_id.currency_id and self.currency_id.id or False,
			      'payment_id': self.id,
			  }

			*/
			return &pool.AccountMoveLineData{}
		})

	pool.AccountPayment().Methods().GetLiquidityMoveLineVals().DeclareMethod(
		`GetLiquidityMoveLineVals`,
		func(rs pool.AccountPaymentSet, amount float64) *pool.AccountMoveLineData {
			/*def _get_liquidity_move_line_vals(self, amount):
			  name = self.name
			  if self.payment_type == 'transfer':
			      name = _('Transfer to %s') % self.destination_journal_id.name
			  vals = {
			      'name': name,
			      'account_id': self.payment_type in ('outbound','transfer') and self.journal_id.default_debit_account_id.id or self.journal_id.default_credit_account_id.id,
			      'payment_id': self.id,
			      'journal_id': self.journal_id.id,
			      'currency_id': self.currency_id != self.company_id.currency_id and self.currency_id.id or False,
			  }

			  # If the journal has a currency specified, the journal item need to be expressed in this currency
			  if self.journal_id.currency_id and self.currency_id != self.journal_id.currency_id:
			      amount = self.currency_id.with_context(date=self.payment_date).compute(amount, self.journal_id.currency_id)
			      debit, credit, amount_currency, dummy = self.env['account.move.line'].with_context(date=self.payment_date).compute_amount_fields(amount, self.journal_id.currency_id, self.company_id.currency_id)
			      vals.update({
			          'amount_currency': amount_currency,
			          'currency_id': self.journal_id.currency_id.id,
			      })

			  return vals
			*/
			return &pool.AccountMoveLineData{}
		})

}
