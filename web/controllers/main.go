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
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep-addons/base"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func WebClient(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Set("uid", int64(1))
	sess.Set("ID", 123)
	sess.Set("login", "admin")
	ctx, _ := json.Marshal(&models.Context{"tz": "Europe/Paris", "lang": "en_US"})
	sess.Set("user_context", ctx)
	sess.Save()
	data := gin.H{
		"Menu": base.TopMenu,
	}
	c.HTML(http.StatusOK, "web.webclient_bootstrap", data)
}

func CompanyLogo(c *gin.Context) {
	c.File("config/img/logo.png")
}

func SessionInfo(sess sessions.Session) gin.H {
	var userContext models.Context
	if sess.Get("uid") != nil && sess.Get("user_context") != nil {
		if json.Unmarshal(sess.Get("user_context").([]byte), &userContext) != nil {
			userContext = models.Context{}
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
	server.RPC(c, http.StatusOK, server.Modules)
}

func Load(c *gin.Context) {
	qwebParams := struct {
		Path string `json:"path"`
	}{}
	server.BindRPCParams(c, &qwebParams)
	path, _ := url.ParseRequestURI(qwebParams.Path)
	targetURL := tools.AbsolutizeURL(c.Request, path.RequestURI())
	resp, err := http.Get(targetURL)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	server.RPC(c, http.StatusOK, string(body))
}

func QWeb(c *gin.Context) {
	mods := strings.Split(c.Query("mods"), ",")
	fileNames := tools.ListStaticFiles("src/xml", mods)
	res, _ := tools.ConcatXML(fileNames)
	c.String(http.StatusOK, string(res))
}

func BootstrapTranslations(c *gin.Context) {
	res := gin.H{
		"lang_parameters": tools.LangParameters{
			DateFormat:   "%m/%d/%Y",
			Direction:    tools.LANG_DIRECTION_LTR,
			ThousandsSep: ",",
			TimeFormat:   "%H:%M:%S",
			DecimalPoint: ".",
			ID:           1,
			Grouping:     "[]",
		},
		"modules": gin.H{},
	}
	server.RPC(c, http.StatusOK, res)
}

func CSSList(c *gin.Context) {
	Params := struct {
		Mods string `json:"mods"`
	}{}
	server.BindRPCParams(c, &Params)
	mods := strings.Split(Params.Mods, ",")
	fileNames := tools.ListStaticFiles("src/css", mods)
	server.RPC(c, http.StatusOK, fileNames)
}

func JSList(c *gin.Context) {
	Params := struct {
		Mods string `json:"mods"`
	}{}
	server.BindRPCParams(c, &Params)
	mods := strings.Split(Params.Mods, ",")
	fileNames := tools.ListStaticFiles("src/js", mods)
	server.RPC(c, http.StatusOK, fileNames)
}

func VersionInfo(c *gin.Context) {
	data := gin.H{
		"server_serie":        "8.0",
		"server_version_info": []int8{8, 0, 0, 0, 0},
		"server_version":      "8.0",
		"protocol":            1,
	}
	server.RPC(c, http.StatusOK, data)
}

func CallKW(c *gin.Context) {
	sess := sessions.Default(c)
	uid := sess.Get("uid").(int64)
	var params server.CallParams
	server.BindRPCParams(c, &params)
	res := server.Execute(uid, params)
	server.RPC(c, http.StatusOK, res)
}

func ActionLoad(c *gin.Context) {
	params := struct {
		ActionID          string `json:"action_id"`
		AdditionalContext string `json:"additional_context"`
	}{}
	server.BindRPCParams(c, &params)
	action := base.ActionsRegistry.GetActionById(params.ActionID)
	server.RPC(c, http.StatusOK, action)
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
		web.GET("/binary/company_logo", CompanyLogo)

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
		}
		action := web.Group("/action")
		{
			action.POST("/load", ActionLoad)
		}
	}
}
