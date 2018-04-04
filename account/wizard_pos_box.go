// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.CashBox().DeclareTransientModel()
	h.CashBox().AddFields(map[string]models.FieldDefinition{
		"Name":   models.CharField{String: "Name" /*[string 'Reason']*/, Required: true},
		"Amount": models.FloatField{String: "Amount" /*[string 'Amount']*/, Digits: nbutils.Digits{0, 0}, Required: true},
	})
	h.CashBox().Methods().Run().DeclareMethod(
		`Run`,
		func(rs h.CashBoxSet) {
			//@api.multi
			/*def run(self):
			  context = dict(self._context or {})
			  active_model = context.get('active_model', False)
			  active_ids = context.get('active_ids', [])

			  records = self.env[active_model].browse(active_ids)

			  return self._run(records)

			*/
		})
	h.CashBox().Methods().RunPrivate().DeclareMethod(
		`RunPrivate`,
		func(rs h.CashBoxSet, args struct {
			Records interface{}
		}) {
			//@api.multi
			/*def _run(self, records):
			  for box in self:
			      for record in records:
			          if not record.journal_id:
			              raise UserError(_("Please check that the field 'Journal' is set on the Bank Statement"))
			          if not record.journal_id.company_id.transfer_account_id:
			              raise UserError(_("Please check that the field 'Transfer Account' is set on the company."))
			          box._create_bank_statement_line(record)
			  return {}

			*/
		})
	h.CashBox().Methods().CreateBankStatementLine().DeclareMethod(
		`CreateBankStatementLine`,
		func(rs h.CashBoxSet, args struct {
			Record interface{}
		}) {
			//@api.one
			/*def _create_bank_statement_line(self, record):
			  if record.state == 'confirm':
			      raise UserError(_("You cannot put/take money in/out for a bank statement which is closed."))
			  values = self._calculate_values_for_statement_line(record)
			  return record.write({'line_ids': [(0, False, values)]})


			*/
		})

	h.CashBoxIn().DeclareTransientModel()
	h.CashBoxIn().AddFields(map[string]models.FieldDefinition{
		"Ref": models.CharField{String: "Reference')" /*['Reference']*/},
	})
	h.CashBoxIn().InheritModel(h.CashBox())
	h.CashBoxIn().Methods().CalculateValuesForStatementLine().DeclareMethod(
		`CalculateValuesForStatementLine`,
		func(rs h.CashBoxInSet, args struct {
			Record interface{}
		}) {
			//@api.multi
			/*def _calculate_values_for_statement_line(self, record):
			  if not record.journal_id.company_id.transfer_account_id:
			      raise UserError(_("You should have defined an 'Internal Transfer Account' in your cash register's journal!"))
			  return {
			      'date': record.date,
			      'statement_id': record.id,
			      'journal_id': record.journal_id.id,
			      'amount': self.amount or 0.0,
			      'account_id': record.journal_id.company_id.transfer_account_id.id,
			      'ref': '%s' % (self.ref or ''),
			      'name': self.name,
			  }


			*/
		})

	h.CashBoxOut().DeclareTransientModel()
	h.CashBoxOut().InheritModel(h.CashBox())
	h.CashBoxOut().Methods().CalculateValuesForStatementLine().DeclareMethod(
		`CalculateValuesForStatementLine`,
		func(rs h.CashBoxOutSet, args struct {
			Record interface{}
		}) {
			//@api.multi
			/*def _calculate_values_for_statement_line(self, record):
			  if not record.journal_id.company_id.transfer_account_id:
			      raise UserError(_("You should have defined an 'Internal Transfer Account' in your cash register's journal!"))
			  return {
			      'date': record.date,
			      'statement_id': record.id,
			      'journal_id': record.journal_id.id,
			      'amount': self.amount or 0.0,
			      'account_id': record.journal_id.company_id.transfer_account_id.id,
			      'ref': '%s' % (self.ref or ''),
			      'name': self.name,
			  }


			*/
		})

}
