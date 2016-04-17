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
	"github.com/npiganeau/yep-addons/base/ir"
)

func createActions() {
	ir.ActionsRegistry = ir.NewActionsCollection()

	dashboardAction := ir.BaseAction{
		ID:    "base_action_dashboard",
		Type:  ir.ACTION_ACT_WINDOW,
		Name:  "Dashboard",
		Model: "Dashboard",
	}
	ir.ActionsRegistry.AddAction(&dashboardAction)

	reportingConfigAction := ir.BaseAction{
		ID:    "base_action_reporting_config",
		Type:  ir.ACTION_ACT_WINDOW,
		Name:  "Configuration",
		Model: "DashboardSettings",
	}
	ir.ActionsRegistry.AddAction(&reportingConfigAction)

	usersAction := ir.BaseAction{
		ID:    "base_action_res_users",
		Type:  ir.ACTION_ACT_WINDOW,
		Name:  "Users",
		Model: "ResUsers",
	}
	ir.ActionsRegistry.AddAction(&usersAction)
}

func createMenus() {
	ir.TopMenu = new(ir.MenuCollection)

	// Reporting Menu
	menuReporting := ir.UiMenu{
		ID:       "base_menu_reporting",
		Name:     "Reporting",
		Sequence: 170,
	}
	ir.TopMenu.AddMenu(&menuReporting)

	menuDashboard := ir.UiMenu{
		ID:       "base_menu_reporting_dashboard",
		Name:     "Dashboards",
		Parent:   &menuReporting,
		Sequence: 0,
	}
	ir.TopMenu.AddMenu(&menuDashboard)

	menuMyDashboard := ir.UiMenu{
		ID:       "base_menu_reporting_dashboard_my",
		Name:     "My Dashboard",
		Parent:   &menuDashboard,
		Sequence: 0,
		Action:   ir.ActionsRegistry.GetActionById("base_action_dashboard"),
	}
	ir.TopMenu.AddMenu(&menuMyDashboard)

	menuConfiguration := ir.UiMenu{
		ID:       "base_menu_reporting_config",
		Name:     "Configuration",
		Parent:   &menuReporting,
		Sequence: 100,
		Action:   ir.ActionsRegistry.GetActionById("base_action_reporting_config"),
	}
	ir.TopMenu.AddMenu(&menuConfiguration)

	// Settings menu
	menuSettings := ir.UiMenu{
		ID:       "base_menu_administration",
		Name:     "Settings",
		Sequence: 255,
	}
	ir.TopMenu.AddMenu(&menuSettings)

	menuUsers := ir.UiMenu{
		ID:       "base_menu_users",
		Name:     "Users",
		Parent:   &menuSettings,
		Sequence: 4,
	}
	ir.TopMenu.AddMenu(&menuUsers)

	menuActionUsers := ir.UiMenu{
		ID:       "base_menu_action_users",
		Name:     "Users",
		Parent:   &menuUsers,
		Sequence: 1,
		Action:   ir.ActionsRegistry.GetActionById("base_action_res_users"),
	}
	ir.TopMenu.AddMenu(&menuActionUsers)
}
