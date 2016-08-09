/*   Copyright (C) 2008-2016 by Nicolas Piganeau
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

import (
	"fmt"
	"time"

	"github.com/npiganeau/yep/yep/ir"
	"github.com/npiganeau/yep/yep/models"
)

type ResUsers struct {
	ID          int64
	LoginDate   time.Time   `yep:"type(datetime)"`
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
