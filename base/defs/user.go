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
	"fmt"

	"github.com/npiganeau/yep/yep/ir"
	"github.com/npiganeau/yep/yep/models"
)

type ResUsers struct {
	ID          int64
	LoginDate   models.DateTime
	Partner     *ResPartner `yep:"inherits"`
	Name        string
	Login       string
	Password    string
	NewPassword string
	Signature   string
	Active      bool
	ActionId    ir.ActionRef `yep:"type(char)"`
	//GroupsID *ir.Group
	Company    *ResCompany
	CompanyIds []*ResCompany `yep:"json(company_ids);type(many2many)"`
	ImageSmall string
}

func NameGet(rs models.RecordSet) string {
	res := rs.Super()
	user := struct {
		ID    int64
		Login string
	}{}
	rs.ReadOne(&user)
	return fmt.Sprintf("%s (%s)", res, user.Login)
}

func initUsers() {
	models.CreateModel("ResUsers")
	models.ExtendModel("ResUsers", new(ResUsers))
	models.DeclareMethod("ResUsers", "NameGet", NameGet)
}
