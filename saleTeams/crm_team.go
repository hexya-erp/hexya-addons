// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package saleTeams

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.CRMTeam().DeclareModel()

	h.CRMTeam().Methods().GetDefaultTeam().DeclareMethod(
		`GetDefaultTeam returns the default sales team`,
		func(rs h.CRMTeamSet, user h.UserSet) h.CRMTeamSet {
			if user.IsEmpty() {
				user = h.User().NewSet(rs.Env()).CurrentUser()
			}
			var team h.CRMTeamSet
			if rs.Env().Context().HasKey("default_team_id") {
				team = h.CRMTeam().Browse(rs.Env(), []int64{rs.Env().Context().GetInteger("default_team_id")})
				if !team.IsEmpty() {
					return team
				}
			}

			company := h.Company().NewSet(rs.Env()).CompanyDefaultGet()
			cond := q.CRMTeam().User().Equals(user).Or().Members().Equals(user).AndCond(
				q.CRMTeam().Company().ChildOf(company).Or().Company().IsNull())
			team = h.CRMTeam().Search(rs.Env(), cond).Limit(1)
			if !team.IsEmpty() {
				return team
			}

			return h.CRMTeam().Search(rs.Env(),
				q.CRMTeam().HexyaExternalID().Equals("sale_teamss_team_sales_department"))
		})

	h.CRMTeam().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{String: "Sales Team", Required: true, Translate: true},
		"Active": models.BooleanField{Default: models.DefaultValue(true),
			Help: "If the active field is set to false, it will allow you to hide the sales team without removing it."},
		"Company": models.Many2OneField{RelationModel: h.Company(),
			Default: func(env models.Environment) interface{} {
				return h.Company().NewSet(env).CompanyDefaultGet()
			}},
		"User": models.Many2OneField{String: "Team Leader", RelationModel: h.User()},
		"Members": models.One2ManyField{String: "Team Members", RelationModel: h.User(), ReverseFK: "SaleTeam",
			JSON: "member_ids"},
		"ReplyTo": models.CharField{String: "Reply-To",
			Help: "The email address put in the 'Reply-To' of all emails sent by Hexya about cases in this sales team"},
		"Color": models.IntegerField{String: "Color Index", Help: "The color of the team"},
	})

	h.CRMTeam().Methods().Create().Extend("",
		func(rs h.CRMTeamSet, data *h.CRMTeamData) h.CRMTeamSet {
			return rs.WithContext("mail_create_nosubscribe", true).Super().Create(data)
		})

}
