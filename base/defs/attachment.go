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
	"github.com/npiganeau/yep/yep/models"
)

type IrAttachment struct {
	ID          int64
	Name        string `yep:"string(Attachment Name)"`
	DatasFname  string `yep:"string(File Name)"`
	Description string
	//ResName     string      `yep:"string(Resource Name);compute(NameGetResName);store(true)"`
	ResModel string             `yep:"string(Resource Model);help(The database object this attachment will be attached to)"`
	ResId    int64              `yep:"string(Resource ID);help(The record id this is attached to)"`
	Company  pool.ResCompanySet `yep:"type(many2one)"`
	Type     string             `yep:"help(Binary File or URL)"`
	Url      string
	//Datas       string      `yep:"compute(DataGet);string(File Content)"`
	StoreFname string `yep:"string(Stored Filename)"`
	DbDatas    string `yep:"string(Database Data)"`
	FileSize   int    `yep:"string(File Size)"`
}

func initAttachment() {
	models.CreateModel("IrAttachment")
	models.ExtendModel("IrAttachment", new(IrAttachment))
}
