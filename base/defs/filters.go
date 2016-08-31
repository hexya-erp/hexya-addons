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
	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/ir"
	"github.com/npiganeau/yep/yep/models"
)

type IrFilters struct {
	ID        int64
	Model     string
	Domain    string
	Context   string
	Name      string
	IsDefault bool
	User      *pool.ResUsers
	ActionId  ir.ActionRef `yep:"type(char)"`
}

func initFilters() {
	models.CreateModel("IrFilters")
	models.ExtendModel("IrFilters", new(IrFilters))
}
