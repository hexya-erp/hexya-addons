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
	"encoding/base64"
	"fmt"
	"github.com/inconshreveable/log15"
	_ "github.com/npiganeau/yep-addons/base/controllers"
	_ "github.com/npiganeau/yep-addons/base/defs"
	_ "github.com/npiganeau/yep-addons/base/methods"
	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/ir"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools"
	"io/ioutil"
)

const (
	MODULE_NAME string = "base"
	SEQUENCE    uint8  = 100
	NAME        string = "Base"
	VERSION     string = "0.1"
	CATEGORY    string = "Hidden"
	DESCRIPTION string = `
The kernel of YEP, needed for all installation
==============================================
	`
	AUTHOR     string = "NDP Systèmes"
	MAINTAINER string = "NDP Systèmes"
	WEBSITE    string = "http://www.ndp-systemes.fr"
)

var log log15.Logger

func init() {
	log = tools.GetLogger("base")
	server.RegisterModule(&server.Module{Name: MODULE_NAME, PostInit: PostInit})
}

func PostInit() {
	env := models.NewEnvironment(tools.SUPERUSER_ID)
	defer func() {
		if r := recover(); r != nil {
			env.Cr().Rollback()
			tools.LogAndPanic(log, fmt.Sprintf("%v", r))
		}
		env.Cr().Commit()
	}()
	companyBase := pool.ResCompany{
		ID:   1,
		Name: "Your Company",
	}
	partnerAdmin := pool.ResPartner{
		ID:       1,
		Name:     "Administrator",
		Function: "IT Manager",
	}
	avatarImg, _ := ioutil.ReadFile("yep/server/static/base/src/img/avatar.png")
	userAdmin := pool.ResUsers{
		ID:         1,
		Name:       "Administrator",
		Active:     true,
		Company:    pool.NewResCompanySet(env),
		Login:      "admin",
		LoginDate:  models.DateTime{},
		Password:   "admin",
		Partner:    pool.NewResPartnerSet(env),
		ActionId:   ir.MakeActionRef("base_action_res_users"),
		ImageSmall: base64.StdEncoding.EncodeToString(avatarImg),
	}
	if env.Pool("ResPartner").Filter("ID", "=", 1).SearchCount() == 0 {
		env.Pool("ResPartner").Create(&partnerAdmin)
	}
	if env.Pool("ResCompany").Filter("ID", "=", 1).SearchCount() == 0 {
		env.Pool("ResCompany").Call("Create", &companyBase)
	}
	if env.Pool("ResUsers").Filter("ID", "=", 1).SearchCount() == 0 {
		env.Pool("ResUsers").Call("Create", &userAdmin)
	}
}
