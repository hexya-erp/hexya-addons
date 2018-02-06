// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package procurement

import (
	"github.com/hexya-erp/hexya-addons/decimalPrecision"
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
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

	h.ProcurementGroup().DeclareModel()
	h.ProcurementGroup().SetDefaultOrder("ID DESC")

	h.ProcurementGroup().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Reference",
			Default: func(env models.Environment) interface{} {
				return h.Sequence().NewSet(env).NextByCode("procurement.group")
			}, Required: true},
		"MoveType": models.SelectionField{String: "Delivery Type", Selection: types.Selection{
			"direct": "Partial",
			"one":    "All at once",
		}, Default: models.DefaultValue("direct"), Required: true},
		"Procurements": models.One2ManyField{String: "Procurements", RelationModel: h.ProcurementOrder(),
			ReverseFK: "Group", JSON: "procurement_ids"},
	})

	h.ProcurementRule().DeclareModel()
	h.ProcurementRule().SetDefaultOrder("Name")

	h.ProcurementRule().AddFields(map[string]models.FieldDefinition{
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
		"Group":    models.Many2OneField{String: "Fixed Procurement Group", RelationModel: h.ProcurementGroup()},
		"Action":   models.SelectionField{Selection: Actions, Required: true},
		"Sequence": models.IntegerField{String: "Sequence", Default: models.DefaultValue(20)},
		"Company":  models.Many2OneField{String: "Company", RelationModel: h.Company()},
	})

	h.ProcurementOrder().DeclareModel()
	h.ProcurementOrder().SetDefaultOrder("Priority DESC", "DatePlanned", "ID ASC")
	//h.ProcurementOrder().InheritModel(h.MailThread())
	//h.ProcurementOrder().InheritModel(h.NeedActionMixin())

	h.ProcurementOrder().AddFields(map[string]models.FieldDefinition{
		"Name": models.TextField{String: "Description", Required: true},
		"Origin": models.CharField{String: "Source Document",
			Help: "Reference of the document that created this Procurement. This is automatically completed by Hexya."},
		"Company": models.Many2OneField{RelationModel: h.Company(),
			Default: func(env models.Environment) interface{} {
				return h.Company().NewSet(env).CompanyDefaultGet()
			}, Required: true},
		"Priority": models.SelectionField{Selection: Priorities, Default: models.DefaultValue("1"),
			Required: true, Index: true /*track_visibility='onchange')*/},
		"DatePlanned": models.DateTimeField{String: "Scheduled Date",
			Default: func(env models.Environment) interface{} {
				return dates.Now()
			}, Required: true, Index: true /*[ track_visibility 'onchange']*/},
		"Group": models.Many2OneField{String: "Procurement Group", RelationModel: h.ProcurementGroup()},
		"Rule": models.Many2OneField{RelationModel: h.ProcurementRule(), /*[ track_visibility 'onchange']*/
			Help: `Chosen rule for the procurement resolution. Usually chosen by the system but can be manually set by
the procurement manager to force an unusual behavior.`},
		"Product": models.Many2OneField{String: "Product", RelationModel: h.ProductProduct(),
			Required: true, OnChange: h.ProcurementOrder().Methods().OnchangeProduct(),
			/* readonly=true */ /*[ states {'confirmed': [('readonly']*/ /*[ False)]}]*/},
		"ProductQty": models.FloatField{String: "Quantity",
			Digits:   decimalPrecision.GetPrecision("Product Unit of Measure"), /*[ readonly True]*/
			Required: true /*[ states {'confirmed': [('readonly']*/ /*[ False)]}]*/},
		"ProductUom": models.Many2OneField{String: "Product Unit of Measure", RelationModel: h.ProductUom(), /* readonly=true */
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

	//h.ProcurementOrder().Methods().NeedactionDomainGet().DeclareMethod(
	//	`NeedactionDomainGet`,
	//	func(rs h.ProcurementOrderSet) {
	//		//@api.model
	//		/*def _needaction_domain_get(self):
	//			  return [('state', '=', 'exception')]
	//
	//		  */})

	h.ProcurementOrder().Methods().Create().Extend("",
		func(rs h.ProcurementOrderSet, data *h.ProcurementOrderData) h.ProcurementOrderSet {
			procurement := rs.Super().Create(data)
			if !rs.Env().Context().HasKey("procurement_autorun_defer") {
				procurement.Run(false)
			}
			return procurement
		})

	h.ProcurementOrder().Methods().Unlink().Extend("",
		func(rs h.ProcurementOrderSet) int64 {
			for _, proc := range rs.Records() {
				if proc.State() == "cancel" {
					panic(rs.T("You cannot delete procurements that are in cancel state."))
				}
			}
			return rs.Super().Unlink()
		})

	h.ProcurementOrder().Methods().DoViewProcurements().DeclareMethod(
		`DoViewProcurements returns an action that display existing procurement orders
			  of same procurement group of given ids`,
		func(rs h.ProcurementOrderSet) *actions.Action {
			action := actions.Registry.GetById("procurement_do_view_procurements")
			action.Domain = "[('group_id', 'in', self.mapped('group_id').ids)]"
			return action
		})

	h.ProcurementOrder().Methods().OnchangeProduct().DeclareMethod(
		`OnchangeProduct updates the UI when the user changes product`,
		func(rs h.ProcurementOrderSet) (*h.ProcurementOrderData, []models.FieldNamer) {
			if !rs.Product().IsEmpty() {
				return &h.ProcurementOrderData{
					ProductUom: rs.Product().Uom(),
				}, []models.FieldNamer{h.ProcurementOrder().ProductUom()}
			}
			return &h.ProcurementOrderData{}, []models.FieldNamer{}
		})

	h.ProcurementOrder().Methods().Cancel().DeclareMethod(
		`Cancel these procurements`,
		func(rs h.ProcurementOrderSet) bool {
			toCancel := rs.Search(q.ProcurementOrder().State().NotEquals("done"))
			if !toCancel.IsEmpty() {
				return toCancel.Write(&h.ProcurementOrderData{State: "cancel"})
			}
			return false
		})

	h.ProcurementOrder().Methods().ResetToConfirmed().DeclareMethod(
		`ResetToConfirmed sets this procurement back to the confirmed state.`,
		func(rs h.ProcurementOrderSet) bool {
			return rs.Write(&h.ProcurementOrderData{State: "confirmed"})
		})

	h.ProcurementOrder().Methods().Run().DeclareMethod(
		`Run resolves these procurements`,
		func(rs h.ProcurementOrderSet, autocommit bool) bool {
			runProcurement := func(proc h.ProcurementOrderSet) bool {
				if !proc.Assign() {
					//proc.MessagePost(rs.T("No rule matching this procurement"))
					proc.Write(&h.ProcurementOrderData{State: "exception"})
					return false
				}
				res := proc.RunPrivate()
				if !res {
					proc.Write(&h.ProcurementOrderData{State: "exception"})
					return false
				}
				proc.Write(&h.ProcurementOrderData{State: "running"})
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

	h.ProcurementOrder().Methods().Check().DeclareMethod(
		`Check updates the state of fulfilled procurements to done.`,
		func(rs h.ProcurementOrderSet, autocommit bool) bool {
			checkProcurement := func(proc h.ProcurementOrderSet) bool {
				if !proc.CheckPrivate() {
					return false
				}
				proc.Write(&h.ProcurementOrderData{State: "done"})
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

	h.ProcurementOrder().Methods().FindSuitableRule().DeclareMethod(
		`FindSuitableRuler eturns a procurement.rule that depicts what to do with the given procurement
			  in order to complete its needs.`,
		func(rs h.ProcurementOrderSet) h.ProcurementRuleSet {
			return h.ProcurementRule().NewSet(rs.Env())
		})

	h.ProcurementOrder().Methods().Assign().DeclareMethod(
		`Assign check what to do with the given procurement in order to complete its needs.
			  It returns False if no solution is found, otherwise it stores the matching rule (if any) and
			  returns True.`,
		func(rs h.ProcurementOrderSet) bool {
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
				rs.Write(&h.ProcurementOrderData{Rule: rule})
				return true
			}
			return false
		})

	h.ProcurementOrder().Methods().RunPrivate().DeclareMethod(
		`RunPrivate implements the resolution of the given procurement.
		It returns true if the resolution of the procurement was a success, false otherwise to set it in exception`,
		func(rs h.ProcurementOrderSet) bool {
			return true
		})

	h.ProcurementOrder().Methods().CheckPrivate().DeclareMethod(
		`CheckPrivate returns True if the given procurement is fulfilled, False otherwise`,
		func(rs h.ProcurementOrderSet) bool {
			return false
		})

	h.ProcurementOrder().Methods().RunScheduler().DeclareMethod(
		`RunScheduler calls the scheduler to check the procurement order.
              This is intented to be done for all existing companies at the
              same time, so we're running all the methods as SUPERUSER to
              avoid intercompany and access rights issues.

			  If useNewCursor is set, each procurement run is done in a new
              environment with a new cursor. This is appropriate for batch jobs only.`,
		func(rs h.ProcurementOrderSet, useNewCursor bool, company h.CompanySet) {
			procurementSudo := h.ProcurementOrder().NewSet(rs.Env()).Sudo()
			// Run confirmed procurements
			cond := q.ProcurementOrder().State().Equals("confirmed")
			if !company.IsEmpty() {
				cond = cond.And().Company().Equals(company)
			}
			for procurements := procurementSudo.Search(cond); !procurements.IsEmpty(); {
				procurements.Run(useNewCursor)
				newCond := cond.And().ID().NotIn(procurements.Ids())
				procurements = procurementSudo.Search(newCond)
			}

			// Check done procurements
			cond = q.ProcurementOrder().State().Equals("running")
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
