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

import (
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
	User      *ResUsers
	ActionId  ir.ActionRef `yep:"type(char)"`
}

func GetFilters(rs models.RecordSet, modelName, actionID string) []*IrFilters {
	var res []*IrFilters
	actRef := ir.MakeActionRef(actionID)
	rs.Filter("Model", "=", modelName).Filter("ActionId", "=", actRef.String()).Filter("User.ID", "=", rs.Env().Uid()).ReadAll(&res)
	return res
}

func initFilters() {
	models.CreateModel("IrFilters")
	models.ExtendModel("IrFilters", new(IrFilters))
	models.DeclareMethod("IrFilters", "GetFilters", GetFilters)
}
