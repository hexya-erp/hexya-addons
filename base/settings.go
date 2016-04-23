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
				<field name="Lang"/>
				<field name="LoginDate"/>
			</tree>`,
	}
	ir.ViewsRegistry.AddView(&usersViewTree)

	usersViewSearch := ir.View{
		ID:    "base_view_users_search",
		Name:  "res.users.search",
		Model: "ResUsers",
		Arch: `
			<search string="Users">
            	<field name="name" filter_domain="['|', '|', ('name','ilike',self), ('login','ilike',self), ('email','ilike',self)]" string="User"/>
                <field name="company_ids" string="Company" groups="base_group_multi_company"/>
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
		ViewMode:   "list",
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
}
