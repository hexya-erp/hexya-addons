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
	"github.com/npiganeau/yep/yep/ir"
	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools"
)

func WebClient(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Set("uid", int64(1))
	sess.Set("ID", 123)
	sess.Set("login", "admin")
	ctx, _ := json.Marshal(&tools.Context{"tz": "Europe/Paris", "lang": "en_US"})
	sess.Set("user_context", ctx)
	sess.Save()
	data := gin.H{
		"Menu": ir.MenusRegistry,
	}
	c.HTML(http.StatusOK, "web.webclient_bootstrap", data)
}

func init() {
	server := server.GetServer()
	server.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/web")
	})
	web := server.Group("/web")
	{
		web.Static("/static", "yep/server/static/web/")

		web.GET("/", WebClient)
		binary := web.Group("/binary")
		{
			binary.GET("/company_logo", CompanyLogo)
			binary.GET("/image", Image)
		}

		sess := web.Group("/session")
		{
			sess.POST("/get_session_info", GetSessionInfo)
			sess.POST("/modules", Modules)
		}

		proxy := web.Group("/proxy")
		{
			proxy.POST("/load", Load)
		}

		webClient := web.Group("/webclient")
		{
			webClient.GET("/qweb", QWeb)
			webClient.POST("/bootstrap_translations", BootstrapTranslations)
			webClient.POST("/translations", BootstrapTranslations)
			webClient.POST("/csslist", CSSList)
			webClient.POST("/jslist", JSList)
			webClient.POST("/version_info", VersionInfo)
		}
		dataset := web.Group("/dataset")
		{
			dataset.POST("/call_kw/*path", CallKW)
			dataset.POST("/search_read", SearchRead)
		}
		action := web.Group("/action")
		{
			action.POST("/load", ActionLoad)
		}
	}
}
