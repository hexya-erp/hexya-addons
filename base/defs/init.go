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

package defs

import (
	"github.com/npiganeau/yep/yep/ir"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/tools"
	"time"
)

func init() {
	initPartner()
	initCompany()
	initUsers()
	initFilters()
	initAttachment()
}

func PostInit() {
	env := models.NewCursorEnvironment(tools.SUPERUSER_ID)
	companyBase := ResCompany{
		Name: "Your Company",
	}
	partnerAdmin := ResPartner{
		Name:     "Administrator",
		Function: "IT Manager",
	}
	userAdmin := ResUsers{
		Name:      "Administrator",
		Active:    true,
		Company:   &companyBase,
		Login:     "admin",
		LoginDate: time.Now(),
		Password:  "admin",
		Partner:   &partnerAdmin,
		ActionId:  ir.MakeActionRef("base_action_res_users"),
	}
	env.Pool("ResPartner").Call("Create", &partnerAdmin)
	env.Pool("ResCompany").Call("Create", &companyBase)
	env.Pool("ResUsers").Call("Create", &userAdmin)
}
