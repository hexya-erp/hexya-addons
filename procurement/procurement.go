// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package procurement

import (
	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool"
)

var (
	Actions    = types.Selection{}
	Priorities = types.Selection{
		"0": "Not Urgent",
		"1": "Normal",
		"2": "Urgent",
		"3": "Very Urgent",
	}
)

func init() {

	pool.ProcurementGroup().DeclareModel()
	pool.ProcurementGroup().SetDefaultOrder("ID DESC")

	pool.ProcurementGroup().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Reference",
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.Sequence().NewSet(env).NextByCode("procurement.group")
			}, Required: true},
		"MoveType": models.SelectionField{String: "Delivery Type", Selection: types.Selection{
			"direct": "Partial",
			"one":    "All at once",
		}, Default: models.DefaultValue("direct"), Required: true},
		"Procurements": models.One2ManyField{String: "Procurements", RelationModel: pool.ProcurementOrder(),
			ReverseFK: "Group", JSON: "procurement_ids"},
	})

	pool.ProcurementRule().DeclareModel()
	pool.ProcurementRule().SetDefaultOrder("Name")

	pool.ProcurementRule().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Name", Required: true, Translate: true,
			Help: "This field will fill the packing origin and the name of its moves"},
		"Active": models.BooleanField{String: "Active", Default: models.DefaultValue(true),
			Help: "If unchecked, it will allow you to hide the rule without removing it."},
		"GroupPropagationOption": models.SelectionField{String: "Propagation of Procurement Group",
			Selection: types.Selection{
				"none":      "Leave Empty",
				"propagate": "Propagate",
				"fixed":     "Fixed",
			}, Default: models.DefaultValue("propagate")},
		"Group":    models.Many2OneField{String: "Fixed Procurement Group", RelationModel: pool.ProcurementGroup()},
		"Action":   models.SelectionField{Selection: Actions, Required: true},
		"Sequence": models.IntegerField{String: "Sequence", Default: models.DefaultValue(20)},
		"Company":  models.Many2OneField{String: "Company", RelationModel: pool.Company()},
	})

	pool.ProcurementOrder().DeclareModel()
	pool.ProcurementOrder().SetDefaultOrder("Priority DESC", "DatePlanned", "ID ASC")
	//pool.ProcurementOrder().InheritModel(pool.MailThread())
	//pool.ProcurementOrder().InheritModel(pool.NeedActionMixin())

	pool.ProcurementOrder().AddFields(map[string]models.FieldDefinition{
		"Name": models.TextField{String: "Description", Required: true},
		"Origin": models.CharField{String: "Source Document",
			Help: "Reference of the document that created this Procurement. This is automatically completed by Hexya."},
		"Company": models.Many2OneField{RelationModel: pool.Company(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.Company().NewSet(env).CompanyDefaultGet()
			}, Required: true},
		"Priority": models.SelectionField{Selection: Priorities, Default: models.DefaultValue("1"),
			Required: true, Index: true /*track_visibility='onchange')*/},
		"DatePlanned": models.DateTimeField{String: "Scheduled Date",
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return dates.Now()
			}, Required: true, Index: true /*[ track_visibility 'onchange']*/},
		"Group": models.Many2OneField{String: "Procurement Group", RelationModel: pool.ProcurementGroup()},
		"Rule": models.Many2OneField{RelationModel: pool.ProcurementRule(), /*[ track_visibility 'onchange']*/
			Help: `Chosen rule for the procurement resolution. Usually chosen by the system but can be manually set by
the procurement manager to force an unusual behavior.`},
		"Product": models.Many2OneField{String: "Product", RelationModel: pool.ProductProduct(),
			Required: true, OnChange: pool.ProcurementOrder().Methods().OnchangeProduct(),
			/* readonly=true */ /*[ states {'confirmed': [('readonly']*/ /*[ False)]}]*/},
		"ProductQty": models.FloatField{String: "Quantity",
			Digits:   decimalPrecision.GetPrecision("Product Unit of Measure"), /*[ readonly True]*/
			Required: true /*[ states {'confirmed': [('readonly']*/ /*[ False)]}]*/},
		"ProductUom": models.Many2OneField{String: "Product Unit of Measure", RelationModel: pool.ProductUom(), /* readonly=true */
			Required: true /*[ states {'confirmed': [('readonly']*/ /*[ False)]}]*/},
		"State": models.SelectionField{String: "Status", Selection: types.Selection{
			"cancel":    "Cancelled",
			"confirmed": "Confirmed",
			"exception": "Exception",
			"running":   "Running",
			"done":      "Done",
		}, Default: models.DefaultValue("confirmed"), NoCopy: true, Required: true,
		/*[ track_visibility 'onchange']*/},
	})

	//pool.ProcurementOrder().Methods().NeedactionDomainGet().DeclareMethod(
	//	`NeedactionDomainGet`,
	//	func(rs pool.ProcurementOrderSet) {
	//		//@api.model
	//		/*def _needaction_domain_get(self):
	//			  return [('state', '=', 'exception')]
	//
	//		  */})

	pool.ProcurementOrder().Methods().Create().Extend("",
		func(rs pool.ProcurementOrderSet, data *pool.ProcurementOrderData) pool.ProcurementOrderSet {
			procurement := rs.Super().Create(data)
			if !rs.Env().Context().HasKey("procurement_autorun_defer") {
				procurement.Run(false)
			}
			return procurement
		})

	pool.ProcurementOrder().Methods().Unlink().Extend("",
		func(rs pool.ProcurementOrderSet) int64 {
			for _, proc := range rs.Records() {
				if proc.State() == "cancel" {
					panic(rs.T("You cannot delete procurements that are in cancel state."))
				}
			}
			return rs.Super().Unlink()
		})

	pool.ProcurementOrder().Methods().DoViewProcurements().DeclareMethod(
		`DoViewProcurements returns an action that display existing procurement orders
			  of same procurement group of given ids`,
		func(rs pool.ProcurementOrderSet) *actions.Action {
			action := actions.Registry.GetById("procurement_do_view_procurements")
			action.Domain = "[('group_id', 'in', self.mapped('group_id').ids)]"
			return action
		})

	pool.ProcurementOrder().Methods().OnchangeProduct().DeclareMethod(
		`OnchangeProduct updates the UI when the user changes product`,
		func(rs pool.ProcurementOrderSet) (*pool.ProcurementOrderData, []models.FieldNamer) {
			if !rs.Product().IsEmpty() {
				return &pool.ProcurementOrderData{
					ProductUom: rs.Product().Uom(),
				}, []models.FieldNamer{pool.ProcurementOrder().ProductUom()}
			}
			return &pool.ProcurementOrderData{}, []models.FieldNamer{}
		})

	pool.ProcurementOrder().Methods().Cancel().DeclareMethod(
		`Cancel these procurements`,
		func(rs pool.ProcurementOrderSet) bool {
			toCancel := rs.Search(pool.ProcurementOrder().State().NotEquals("done"))
			if !toCancel.IsEmpty() {
				return toCancel.Write(&pool.ProcurementOrderData{State: "cancel"})
			}
			return false
		})

	pool.ProcurementOrder().Methods().ResetToConfirmed().DeclareMethod(
		`ResetToConfirmed sets this procurement back to the confirmed state.`,
		func(rs pool.ProcurementOrderSet) bool {
			return rs.Write(&pool.ProcurementOrderData{State: "confirmed"})
		})

	pool.ProcurementOrder().Methods().Run().DeclareMethod(
		`Run resolves these procurements`,
		func(rs pool.ProcurementOrderSet, autocommit bool) bool {
			runProcurement := func(proc pool.ProcurementOrderSet) bool {
				if !proc.Assign() {
					//proc.MessagePost(rs.T("No rule matching this procurement"))
					proc.Write(&pool.ProcurementOrderData{State: "exception"})
					return false
				}
				res := proc.RunPrivate()
				if !res {
					proc.Write(&pool.ProcurementOrderData{State: "exception"})
					return false
				}
				proc.Write(&pool.ProcurementOrderData{State: "running"})
				return true
			}

			rs.Load("State")
			for _, procurement := range rs.Records() {
				if procurement.State() == "running" || procurement.State() == "done" {
					continue
				}
				if autocommit {
					models.ExecuteInNewEnvironment(rs.Env().Uid(), func(env models.Environment) {
						runProcurement(procurement.WithEnv(env))
					})
				}
				runProcurement(procurement)
			}
			return true
		})

	pool.ProcurementOrder().Methods().Check().DeclareMethod(
		`Check updates the state of fulfilled procurements to done.`,
		func(rs pool.ProcurementOrderSet, autocommit bool) bool {
			checkProcurement := func(proc pool.ProcurementOrderSet) bool {
				if !proc.CheckPrivate() {
					return false
				}
				proc.Write(&pool.ProcurementOrderData{State: "done"})
				return true
			}
			rs.Load("State")
			for _, procurement := range rs.Records() {
				if procurement.State() == "cancel" || procurement.State() == "done" {
					continue
				}
				if autocommit {
					models.ExecuteInNewEnvironment(rs.Env().Uid(), func(env models.Environment) {
						checkProcurement(procurement.WithEnv(env))
					})
				}
				checkProcurement(procurement)
			}
			return true
		})

	pool.ProcurementOrder().Methods().FindSuitableRule().DeclareMethod(
		`FindSuitableRuler eturns a procurement.rule that depicts what to do with the given procurement
			  in order to complete its needs.`,
		func(rs pool.ProcurementOrderSet) pool.ProcurementRuleSet {
			return pool.ProcurementRule().NewSet(rs.Env())
		})

	pool.ProcurementOrder().Methods().Assign().DeclareMethod(
		`Assign check what to do with the given procurement in order to complete its needs.
			  It returns False if no solution is found, otherwise it stores the matching rule (if any) and
			  returns True.`,
		func(rs pool.ProcurementOrderSet) bool {
			// if the procurement already has a rule assigned, we keep it (it has a higher priority as it may have
			// been chosen manually)
			if !rs.Rule().IsEmpty() {
				return true
			}
			if rs.Product().Type() == "service" || rs.Product().Type() == "digital" {
				return false
			}
			rule := rs.FindSuitableRule()
			if !rule.IsEmpty() {
				rs.Write(&pool.ProcurementOrderData{Rule: rule})
				return true
			}
			return false
		})

	pool.ProcurementOrder().Methods().RunPrivate().DeclareMethod(
		`RunPrivate implements the resolution of the given procurement.
		It returns true if the resolution of the procurement was a success, false otherwise to set it in exception`,
		func(rs pool.ProcurementOrderSet) bool {
			return true
		})

	pool.ProcurementOrder().Methods().CheckPrivate().DeclareMethod(
		`CheckPrivate returns True if the given procurement is fulfilled, False otherwise`,
		func(rs pool.ProcurementOrderSet) bool {
			return false
		})

	pool.ProcurementOrder().Methods().RunScheduler().DeclareMethod(
		`RunScheduler calls the scheduler to check the procurement order.
              This is intented to be done for all existing companies at the
              same time, so we're running all the methods as SUPERUSER to
              avoid intercompany and access rights issues.

			  If useNewCursor is set, each procurement run is done in a new
              environment with a new cursor. This is appropriate for batch jobs only.`,
		func(rs pool.ProcurementOrderSet, useNewCursor bool, company pool.CompanySet) {
			procurementSudo := pool.ProcurementOrder().NewSet(rs.Env()).Sudo()
			// Run confirmed procurements
			cond := pool.ProcurementOrder().State().Equals("confirmed")
			if !company.IsEmpty() {
				cond = cond.And().Company().Equals(company)
			}
			for procurements := procurementSudo.Search(cond); !procurements.IsEmpty(); {
				procurements.Run(useNewCursor)
				newCond := cond.And().ID().NotIn(procurements.Ids())
				procurements = procurementSudo.Search(newCond)
			}

			// Check done procurements
			cond = pool.ProcurementOrder().State().Equals("running")
			if !company.IsEmpty() {
				cond = cond.And().Company().Equals(company)
			}
			for procurements := procurementSudo.Search(cond); !procurements.IsEmpty(); {
				procurements.Check(useNewCursor)
				newCond := cond.And().ID().NotIn(procurements.Ids())
				procurements = procurementSudo.Search(newCond)
			}
		})

}
