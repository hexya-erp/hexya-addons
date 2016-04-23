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
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools"
)

func SessionInfo(sess sessions.Session) gin.H {
	var userContext tools.Context
	if sess.Get("uid") != nil && sess.Get("user_context") != nil {
		if json.Unmarshal(sess.Get("user_context").([]byte), &userContext) != nil {
			userContext = tools.Context{}
		}
	}
	return gin.H{
		"session_id":   sess.Get("ID"),
		"uid":          sess.Get("uid"),
		"user_context": userContext,
		"db":           "default",
		"username":     sess.Get("login"),
		"company_id":   1,
	}
}

func GetSessionInfo(c *gin.Context) {
	sess := sessions.Default(c)
	server.RPC(c, http.StatusOK, SessionInfo(sess))
}

func Modules(c *gin.Context) {
	mods := make([]string, len(server.Modules))
	for i, m := range server.Modules {
		mods[i] = m.Name
	}
	server.RPC(c, http.StatusOK, mods)
}
