// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.CRMTeam().DeclareModel()

	pool.CRMTeam().Methods().GetDefaultTeam().DeclareMethod(
		`GetDefaultTeam returns the default sales team`,
		func(rs pool.CRMTeamSet, user pool.UserSet) pool.CRMTeamSet {
			if user.IsEmpty() {
				user = pool.User().NewSet(rs.Env()).CurrentUser()
			}
			var team pool.CRMTeamSet
			if rs.Env().Context().HasKey("default_team_id") {
				team = pool.CRMTeam().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("default_team_id")})
				if !team.IsEmpty() {
					return team
				}
			}

			company := pool.Company().NewSet(rs.Env()).CompanyDefaultGet()
			cond := pool.CRMTeam().User().Equals(user).Or().Members().Equals(user).AndCond(
				pool.CRMTeam().Company().ChildOf(company).Or().Company().IsNull())
			team = pool.CRMTeam().Search(rs.Env(), cond).Limit(1)
			if !team.IsEmpty() {
				return team
			}

			return pool.CRMTeam().Search(rs.Env(),
				pool.CRMTeam().HexyaExternalID().Equals("sale_teamss_team_sales_department"))
		})

	pool.CRMTeam().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Sales Team", Required: true, Translate: true},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If the active field is set to false, it will allow you to hide the sales team without removing it."},
		"Company": models.Many2OneField{RelationModel: pool.Company(),
			Default: func(env models.Environment, vals models.FieldMap) interface{} {
				return pool.Company().NewSet(env).CompanyDefaultGet()
			}},
		"User": models.Many2OneField{String: "Team Leader", RelationModel: pool.User()},
		"Members": models.One2ManyField{String: "Team Members", RelationModel: pool.User(), ReverseFK: "SaleTeam",
			JSON: "member_ids"},
		"ReplyTo": models.CharField{String: "Reply-To",
			Help: "The email address put in the 'Reply-To' of all emails sent by Hexya about cases in this sales team"},
		"Color": models.IntegerField{String: "Color Index", Help: "The color of the team"},
	})

	pool.CRMTeam().Methods().Create().Extend("",
		func(rs pool.CRMTeamSet, data *pool.CRMTeamData) pool.CRMTeamSet {
			return rs.WithContext("mail_create_nosubscribe", true).Super().Create(data)
		})

}
