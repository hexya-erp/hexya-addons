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

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/yep/server"
)

func MenuLoadNeedaction(c *gin.Context) {
	type lnaParams struct {
		MenuIds []string `json:"menu_ids"`
	}
	var params lnaParams
	server.BindRPCParams(c, &params)

	// TODO: update with real needaction support
	type lnaResponse struct {
		NeedactionEnabled bool `json:"needaction_enabled"`
		NeedactionCounter int  `json:"needaction_counter"`
	}
	res := make(map[string]lnaResponse)
	for _, menu := range params.MenuIds {
		res[menu] = lnaResponse{}
	}
	server.RPC(c, http.StatusOK, res)
}
