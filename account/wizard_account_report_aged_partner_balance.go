// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package account

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.AccountAgedTrialBalance().DeclareTransientModel()
	h.AccountAgedTrialBalance().AddFields(map[string]models.FieldDefinition{
		"PeriodLength": models.IntegerField{String: "PeriodLength" /*[string 'Period Length (days)']*/, Required: true, Default: models.DefaultValue(30)},
		"Journals":     models.Many2ManyField{String: "Journals", RelationModel: h.AccountJournal(), JSON: "journal_ids" /*['account.journal']*/ /*[ required True]*/},
		"DateFrom":     models.DateField{String: "DateFrom" /*[default lambda *a: time.strftime('%Y-%m-%d']*/},
	})
	h.AccountAgedTrialBalance().Methods().PrintReport().DeclareMethod(
		`PrintReport`,
		func(rs h.AccountAgedTrialBalanceSet, args struct {
			Data interface{}
		}) {
			/*def _print_report(self, data):
			  res = {}
			  data = self.pre_print_report(data)
			  data['form'].update(self.read(['period_length'])[0])
			  period_length = data['form']['period_length']
			  if period_length<=0:
			      raise UserError(_('You must set a period length greater than 0.'))
			  if not data['form']['date_from']:
			      raise UserError(_('You must set a start date.'))

			  start = datetime.strptime(data['form']['date_from'], "%Y-%m-%d")

			  for i in range(5)[::-1]:
			      stop = start - relativedelta(days=period_length - 1)
			      res[str(i)] = {
			          'name': (i!=0 and (str((5-(i+1)) * period_length) + '-' + str((5-i) * period_length)) or ('+'+str(4 * period_length))),
			          'stop': start.strftime('%Y-%m-%d'),
			          'start': (i!=0 and stop.strftime('%Y-%m-%d') or False),
			      }
			      start = stop - relativedelta(days=1)
			  data['form'].update(res)
			  return self.env['report'].with_context(landscape=True).get_action(self, 'account.report_agedpartnerbalance', data=data)
			*/
		})

}
