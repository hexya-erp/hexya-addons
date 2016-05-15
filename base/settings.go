// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"github.com/npiganeau/yep/yep/ir"
)

/*
createSettings creates menus, actions and views of the main 'Settings' menu.
*/
func createSettings() {
	// Main menus
	menuSettings := ir.UiMenu{
		ID:       "base_menu_administration",
		Name:     "Settings",
		Sequence: 255,
	}
	ir.MenusRegistry.AddMenu(&menuSettings)

	menuUsers := ir.UiMenu{
		ID:       "base_menu_users",
		Name:     "Users",
		Parent:   &menuSettings,
		Sequence: 4,
	}
	ir.MenusRegistry.AddMenu(&menuUsers)

	// Users
	usersViewTree := ir.View{
		ID:    "base_view_users_tree",
		Name:  "res.users.tree",
		Model: "ResUsers",
		Arch: `
			<tree string="Users">
				<field name="Name"/>
				<field name="Login"/>
			<!--    <field name="Lang"/> -->
			<!--	<field name="LoginDate"/> -->
				<field name="Active"/>
			</tree>`,
	}
	ir.ViewsRegistry.AddView(&usersViewTree)

	usersViewForm := ir.View{
		ID:    "base_view_users_form",
		Name:  "res.users.form",
		Model: "ResUsers",
		Arch: `
			<form string="Users">
			<header></header>
			<sheet>
				<group>
					<field name="Name"/>
					<field name="Login"/>
				<!--    <field name="Lang"/> -->
				<!--	<field name="LoginDate"/> -->
					<field name="Active"/>
				</group>
			</sheet>
			</form>`,
	}
	ir.ViewsRegistry.AddView(&usersViewForm)

	usersViewSearch := ir.View{
		ID:    "base_view_users_search",
		Name:  "res.users.search",
		Model: "ResUsers",
		Arch: `
			<search string="Users">
            	<field name="Name" filter_domain="['|', '|', ('Name','ilike',self), ('Login','ilike',self), ('Email','ilike',self)]" string="User"/>
                <field name="CompanyIds" string="Company"/><!-- groups="base_group_multi_company"/>-->
            </search>`,
	}
	ir.ViewsRegistry.AddView(&usersViewSearch)

	usersAction := ir.BaseAction{
		ID:         "base_action_res_users",
		Type:       ir.ACTION_ACT_WINDOW,
		Name:       "Users",
		Model:      "ResUsers",
		View:       ir.MakeViewRef("base_view_users_tree"),
		SearchView: ir.MakeViewRef("base_view_users_search"),
		ViewMode:   "tree,form",
	}
	ir.ActionsRegistry.AddAction(&usersAction)

	menuActionUsers := ir.UiMenu{
		ID:       "base_menu_action_users",
		Name:     "Users",
		Parent:   &menuUsers,
		Sequence: 1,
		Action:   &usersAction,
	}
	ir.MenusRegistry.AddMenu(&menuActionUsers)

	// Partners
	partnersViewTree := ir.View{
		ID:    "base_view_partner_tree",
		Name:  "res.partner.tree",
		Model: "ResPartner",
		Arch: `
			<tree string="Users">
				<field name="Name"/>
			    <field name="Function"/>
			    <field name="Lang"/>
			    <field name="Ref"/>
			</tree>`,
	}
	ir.ViewsRegistry.AddView(&partnersViewTree)

	partnersViewSearch := ir.View{
		ID:    "base_view_partner_search",
		Name:  "res.partner.search",
		Model: "ResPartner",
		Arch: `
			<search string="Partners">
            	<field name="Name" filter_domain="['|', '|', ('Name','ilike',self), ('Email','ilike',self)]" string="Partner"/>
            </search>`,
	}
	ir.ViewsRegistry.AddView(&partnersViewSearch)

	partnersAction := ir.BaseAction{
		ID:         "base_action_res_partner",
		Type:       ir.ACTION_ACT_WINDOW,
		Name:       "Partners",
		Model:      "ResPartner",
		View:       ir.MakeViewRef("base_view_partner_tree"),
		SearchView: ir.MakeViewRef("base_view_partner_search"),
		ViewMode:   "tree",
	}
	ir.ActionsRegistry.AddAction(&partnersAction)

	menuActionPartners := ir.UiMenu{
		ID:       "base_menu_action_partner",
		Name:     "Partners",
		Parent:   &menuUsers,
		Sequence: 1,
		Action:   &partnersAction,
	}
	ir.MenusRegistry.AddMenu(&menuActionPartners)

}
