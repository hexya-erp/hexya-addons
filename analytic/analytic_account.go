// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package analytic

import (
	"fmt"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.AccountAnalyticTag().DeclareModel()

	h.AccountAnalyticTag().AddFields(map[string]models.FieldDefinition{
		"Name":  models.CharField{String: "Analytic Tag", Index: true, Required: true},
		"Color": models.IntegerField{String: "Color Index"},
	})

	h.AccountAnalyticAccount().DeclareModel()
	h.AccountAnalyticAccount().SetDefaultOrder("Code", "Name ASC")

	h.AccountAnalyticAccount().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Analytic Account", Index: true, Required: true /* track_visibility 'onchange' */},
		"Code": models.CharField{String: "Reference", Index: true /*track_visibility 'onchange'*/},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If the active field is set to False, it will allow you to hide the account without removing it."},
		"Tags": models.Many2ManyField{RelationModel: h.AccountAnalyticTag(), JSON: "tag_ids"},
		"Lines": models.One2ManyField{String: "Analytic Lines", RelationModel: h.AccountAnalyticLine(),
			ReverseFK: "Account", JSON: "line_ids"},
		"Company": models.Many2OneField{RelationModel: h.Company(), Required: true,
			Default: func(env models.Environment) interface{} {
				return h.User().NewSet(env).CurrentUser().Company()
			}},
		"Partner": models.Many2OneField{String: "Customer", RelationModel: h.Partner() /* track_visibility 'onchange' */},
		"Balance": models.FloatField{Compute: h.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance()},
		"Debit":   models.FloatField{Compute: h.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance()},
		"Credit":  models.FloatField{Compute: h.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance()},
		"Currency": models.Many2OneField{String: "Currency", RelationModel: h.Currency(),
			Related: "Company.Currency", ReadOnly: true},
	})

	h.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance().DeclareMethod(
		`ComputeDebitCreditBalance`,
		func(rs h.AccountAnalyticAccountSet) *h.AccountAnalyticAccountData {
			cond := q.AccountAnalyticLine().Account().Equals(rs)
			if rs.Env().Context().HasKey("from_date") {
				cond = cond.And().Date().GreaterOrEqual(rs.Env().Context().GetDate("from_date"))
			}
			if rs.Env().Context().HasKey("to_date") {
				cond = cond.And().Date().LowerOrEqual(rs.Env().Context().GetDate("to_date"))
			}
			accountDebit := h.AccountAnalyticLine().Search(rs.Env(), cond.And().Amount().Lower(0)).
				GroupBy(h.AccountAnalyticLine().Account()).Aggregates(h.AccountAnalyticLine().Amount())
			debitVal, _ := accountDebit[0].Values.Get("Amount", h.AccountAnalyticLine().Underlying())
			debit := debitVal.(float64)
			accountCredit := h.AccountAnalyticLine().Search(rs.Env(), cond.And().Amount().GreaterOrEqual(0)).
				GroupBy(h.AccountAnalyticLine().Account()).Aggregates(h.AccountAnalyticLine().Amount())
			creditVal, _ := accountCredit[0].Values.Get("Amount", h.AccountAnalyticLine().Underlying())
			credit := creditVal.(float64)
			return &h.AccountAnalyticAccountData{
				Debit:   debit,
				Credit:  credit,
				Balance: credit - debit,
			}
		})

	h.AccountAnalyticAccount().Methods().NameGet().Extend("",
		func(rs h.AccountAnalyticAccountSet) string {
			name := rs.Name()
			if rs.Code() != "" {
				name = fmt.Sprintf("[%s] %s", rs.Code(), rs.Name())
			}
			if !rs.Partner().IsEmpty() {
				name = fmt.Sprintf("%s - %s", rs.Name(), rs.Partner().CommercialPartner().Name())
			}
			return name
		})

	h.AccountAnalyticAccount().Methods().SearchByName().Extend("",
		func(rs h.AccountAnalyticAccountSet, name string, op operator.Operator, additionalCond q.AccountAnalyticAccountCondition, limit int) h.AccountAnalyticAccountSet {
			if !op.IsPositive() {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			partners := h.Partner().Search(rs.Env(), q.Partner().Name().AddOperator(op, name)).Limit(limit)
			cond := q.AccountAnalyticAccount().Code().AddOperator(op, name).Or().Name().AddOperator(op, name)
			if !partners.IsEmpty() {
				cond = cond.Or().Partner().In(partners)
			}
			return h.AccountAnalyticAccount().Search(rs.Env(), cond.AndCond(additionalCond)).Limit(limit)
		})

	h.AccountAnalyticLine().DeclareModel()
	h.AccountAnalyticLine().SetDefaultOrder("Date DESC", "ID DESC")

	h.AccountAnalyticLine().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Description", Required: true},
		"Date": models.DateField{String: "Date", Required: true, Index: true,
			Default: func(models.Environment) interface{} {
				return dates.Today()
			}},
		"Amount":     models.FloatField{Required: true, Default: models.DefaultValue(0)},
		"UnitAmount": models.FloatField{String: "Quantity", Default: models.DefaultValue(0.0)},
		"Account": models.Many2OneField{String: "Analytic Account", RelationModel: h.AccountAnalyticAccount(),
			Required: true, OnDelete: models.Restrict, Index: true},
		"Partner": models.Many2OneField{RelationModel: h.Partner()},
		"User": models.Many2OneField{String: "User", RelationModel: h.User(),
			Default: func(env models.Environment) interface{} {
				user := h.User().NewSet(env).CurrentUser()
				if env.Context().HasKey("user_id") {
					user = h.User().Browse(env, []int64{env.Context().GetInteger("user_id")})
				}
				return user
			}},
		"Tags":     models.Many2ManyField{RelationModel: h.AccountAnalyticTag(), JSON: "tag_ids"},
		"Company":  models.Many2OneField{RelationModel: h.Company(), Related: "Account.Company", ReadOnly: true},
		"Currency": models.Many2OneField{RelationModel: h.Currency(), Related: "Company.Currency", ReadOnly: true},
	})

}
