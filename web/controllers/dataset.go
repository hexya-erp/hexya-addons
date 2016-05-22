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

package controllers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/yep/server"
)

func CallKW(c *gin.Context) {
	sess := sessions.Default(c)
	uid := sess.Get("uid").(int64)
	var params server.CallParams
	server.BindRPCParams(c, &params)
	res, err := server.Execute(uid, params)
	server.RPC(c, http.StatusOK, res, err)
}

func SearchRead(c *gin.Context) {
	sess := sessions.Default(c)
	uid := sess.Get("uid").(int64)
	var params server.SearchReadParams
	server.BindRPCParams(c, &params)
	res, err := server.SearchRead(uid, params)
	server.RPC(c, http.StatusOK, res, err)
}
