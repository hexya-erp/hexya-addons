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

var TopMenu *ir.MenuCollection
var ActionsRegistry *ir.ActionsCollection

func createActions() {
	ActionsRegistry = ir.NewActionsCollection()

	dashboardAction := ir.BaseAction{
		ID:    "base.action_dashboard",
		Type:  ir.ACTION_ACT_WINDOW,
		Name:  "Dashboard",
		Model: "Dashboard",
	}
	ActionsRegistry.AddAction(&dashboardAction)

	reportingConfigAction := ir.BaseAction{
		ID:    "base.action_reporting_config",
		Type:  ir.ACTION_ACT_WINDOW,
		Name:  "Configuration",
		Model: "DashboardSettings",
	}
	ActionsRegistry.AddAction(&reportingConfigAction)

	usersAction := ir.BaseAction{
		ID:    "base.action_res_users",
		Type:  ir.ACTION_ACT_WINDOW,
		Name:  "Users",
		Model: "ResUsers",
	}
	ActionsRegistry.AddAction(&usersAction)
}

func createMenus() {
	TopMenu = new(ir.MenuCollection)

	// Reporting Menu
	menuReporting := ir.UiMenu{
		ID:       "base.menu_reporting",
		Name:     "Reporting",
		Sequence: 170,
	}
	TopMenu.AddMenu(&menuReporting)

	menuDashboard := ir.UiMenu{
		ID:       "base.menu_reporting_dashboard",
		Name:     "Dashboards",
		Parent:   &menuReporting,
		Sequence: 0,
	}
	TopMenu.AddMenu(&menuDashboard)

	menuMyDashboard := ir.UiMenu{
		ID:       "base.menu_reporting_dashboard_my",
		Name:     "My Dashboard",
		Parent:   &menuDashboard,
		Sequence: 0,
		Action:   ActionsRegistry.GetActionById("base.action_dashboard"),
	}
	TopMenu.AddMenu(&menuMyDashboard)

	menuConfiguration := ir.UiMenu{
		ID:       "base.menu_reporting_config",
		Name:     "Configuration",
		Parent:   &menuReporting,
		Sequence: 100,
		Action:   ActionsRegistry.GetActionById("base.action_reporting_config"),
	}
	TopMenu.AddMenu(&menuConfiguration)

	// Settings menu
	menuSettings := ir.UiMenu{
		ID:       "base.menu_administration",
		Name:     "Settings",
		Sequence: 255,
	}
	TopMenu.AddMenu(&menuSettings)

	menuUsers := ir.UiMenu{
		ID:       "base.menu_users",
		Name:     "Users",
		Parent:   &menuSettings,
		Sequence: 4,
	}
	TopMenu.AddMenu(&menuUsers)

	menuActionUsers := ir.UiMenu{
		ID:       "base.menu_action_users",
		Name:     "Users",
		Parent:   &menuUsers,
		Sequence: 1,
		Action:   ActionsRegistry.GetActionById("base.action_res_users"),
	}
	TopMenu.AddMenu(&menuActionUsers)
}
