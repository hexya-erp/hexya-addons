// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package analytic

import (
	"fmt"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/operator"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.AccountAnalyticTag().DeclareModel()

	pool.AccountAnalyticTag().AddFields(map[string]models.FieldDefinition{
		"Name":  models.CharField{String: "Analytic Tag", Index: true, Required: true},
		"Color": models.IntegerField{String: "Color Index"},
	})

	pool.AccountAnalyticAccount().DeclareModel()
	pool.AccountAnalyticAccount().SetDefaultOrder("Code", "Name ASC")

	pool.AccountAnalyticAccount().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Analytic Account", Index: true, Required: true /* track_visibility 'onchange' */},
		"Code": models.CharField{String: "Reference", Index: true /*track_visibility 'onchange'*/},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If the active field is set to False, it will allow you to hide the account without removing it."},
		"Tags": models.Many2ManyField{RelationModel: pool.AccountAnalyticTag(), JSON: "tag_ids"},
		"Lines": models.One2ManyField{String: "Analytic Lines", RelationModel: pool.AccountAnalyticLine(),
			ReverseFK: "Account", JSON: "line_ids"},
		"Company": models.Many2OneField{RelationModel: pool.Company(), Required: true,
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.User().NewSet(env).CurrentUser().Company()
			}},
		"Partner": models.Many2OneField{String: "Customer", RelationModel: pool.Partner() /* track_visibility 'onchange' */},
		"Balance": models.FloatField{Compute: pool.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance()},
		"Debit":   models.FloatField{Compute: pool.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance()},
		"Credit":  models.FloatField{Compute: pool.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance()},
		"Currency": models.Many2OneField{String: "Currency", RelationModel: pool.Currency(),
			Related: "Company.Currency" /* readonly=true */},
	})

	pool.AccountAnalyticAccount().Methods().ComputeDebitCreditBalance().DeclareMethod(
		`ComputeDebitCreditBalance`,
		func(rs pool.AccountAnalyticAccountSet) (*pool.AccountAnalyticAccountData, []models.FieldNamer) {
			cond := pool.AccountAnalyticLine().Account().Equals(rs)
			if rs.Env().Context().HasKey("from_date") {
				cond = cond.And().Date().GreaterOrEqual(rs.Env().Context().GetDate("from_date"))
			}
			if rs.Env().Context().HasKey("to_date") {
				cond = cond.And().Date().LowerOrEqual(rs.Env().Context().GetDate("to_date"))
			}
			accountDebit := pool.AccountAnalyticLine().Search(rs.Env(), cond.And().Amount().Lower(0)).
				GroupBy(pool.AccountAnalyticLine().Account()).Aggregates(pool.AccountAnalyticLine().Amount())
			debitVal, _ := accountDebit[0].Values.Get("Amount", pool.AccountAnalyticLine().Underlying())
			debit := debitVal.(float64)
			accountCredit := pool.AccountAnalyticLine().Search(rs.Env(), cond.And().Amount().GreaterOrEqual(0)).
				GroupBy(pool.AccountAnalyticLine().Account()).Aggregates(pool.AccountAnalyticLine().Amount())
			creditVal, _ := accountCredit[0].Values.Get("Amount", pool.AccountAnalyticLine().Underlying())
			credit := creditVal.(float64)
			return &pool.AccountAnalyticAccountData{
					Debit:   debit,
					Credit:  credit,
					Balance: credit - debit,
				}, []models.FieldNamer{
					pool.AccountAnalyticAccount().Debit(),
					pool.AccountAnalyticAccount().Credit(),
					pool.AccountAnalyticAccount().Balance(),
				}
		})

	pool.AccountAnalyticAccount().Methods().NameGet().Extend("",
		func(rs pool.AccountAnalyticAccountSet) string {
			name := rs.Name()
			if rs.Code() != "" {
				name = fmt.Sprintf("[%s] %s", rs.Code(), rs.Name())
			}
			if !rs.Partner().IsEmpty() {
				name = fmt.Sprintf("%s - %s", rs.Name(), rs.Partner().CommercialPartner().Name())
			}
			return name
		})

	pool.AccountAnalyticAccount().Methods().SearchByName().Extend("",
		func(rs pool.AccountAnalyticAccountSet, name string, op operator.Operator, additionalCond pool.AccountAnalyticAccountCondition, limit int) pool.AccountAnalyticAccountSet {
			if !op.IsPositive() {
				return rs.Super().SearchByName(name, op, additionalCond, limit)
			}
			partners := pool.Partner().Search(rs.Env(), pool.Partner().Name().AddOperator(op, name)).Limit(limit)
			cond := pool.AccountAnalyticAccount().Code().AddOperator(op, name).Or().Name().AddOperator(op, name)
			if !partners.IsEmpty() {
				cond = cond.Or().Partner().In(partners)
			}
			return pool.AccountAnalyticAccount().Search(rs.Env(), cond.AndCond(additionalCond)).Limit(limit)
		})

	pool.AccountAnalyticLine().DeclareModel()
	pool.AccountAnalyticLine().SetDefaultOrder("Date DESC", "ID DESC")

	pool.AccountAnalyticLine().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Description", Required: true},
		"Date": models.DateField{String: "Date", Required: true, Index: true,
			Default: func(models.Environment, models.FieldMap) interface{} {
				return dates.Today()
			}},
		"Amount":     models.FloatField{Required: true, Default: models.DefaultValue(0)},
		"UnitAmount": models.FloatField{String: "Quantity", Default: models.DefaultValue(0.0)},
		"Account": models.Many2OneField{String: "Analytic Account", RelationModel: pool.AccountAnalyticAccount(),
			Required: true, OnDelete: models.Restrict, Index: true},
		"Partner": models.Many2OneField{RelationModel: pool.Partner()},
		"User": models.Many2OneField{String: "User", RelationModel: pool.User(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				user := pool.User().NewSet(env).CurrentUser()
				if env.Context().HasKey("user_id") {
					user = pool.User().Browse(env, []int64{env.Context().GetInteger("user_id")})
				}
				return user
			}},
		"Tags":     models.Many2ManyField{RelationModel: pool.AccountAnalyticTag(), JSON: "tag_ids"},
		"Company":  models.Many2OneField{RelationModel: pool.Company(), Related: "Account.Company" /* readonly=true */},
		"Currency": models.Many2OneField{RelationModel: pool.Currency(), Related: "Company.Currency" /* readonly=true */},
	})

}
