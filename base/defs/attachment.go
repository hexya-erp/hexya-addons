/*   Copyright (C) 2008-2016 by Nicolas Piganeau and the TS2 team
 *   (See AUTHORS file)
 *
 *   This program is free software; you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation; either version 2 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program; if not, write to the
 *   Free Software Foundation, Inc.,
 *   59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.
 */

package defs

import "github.com/npiganeau/yep/yep/models"

type IrAttachment struct {
	ID          int64
	Name        string `yep:"string(Attachment Name)"`
	DatasFname  string `yep:"string(File Name)"`
	Description string
	//ResName     string      `yep:"string(Resource Name);compute(NameGetResName);store(true)"`
	ResModel string      `yep:"string(Resource Model);help(The database object this attachment will be attached to)"`
	ResId    int64       `yep:"string(Resource ID);help(The record id this is attached to)"`
	Company  *ResCompany `orm:"rel(fk)"`
	Type     string      `yep:"help(Binary File or URL)"`
	Url      string      `orm:"size(1024)"`
	//Datas       string      `yep:"compute(DataGet);string(File Content)"`
	StoreFname string `yep:"string(Stored Filename)"`
	DbDatas    string `yep:"string(Database Data)"`
	FileSize   int    `yep:"string(File Size)"`
}

func initAttachment() {
	models.CreateModel("IrAttachment")
	models.ExtendModel("IrAttachment", new(IrAttachment))
}
