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
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/server"
	"github.com/npiganeau/yep/tools"
	"github.com/npiganeau/yep/yep/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func WebClient(c *gin.Context) {
	c.HTML(http.StatusOK, "web.webclient_bootstrap", gin.H{})
}

func CompanyLogo(c *gin.Context) {
	c.File("config/img/logo.png")
}

func SessionInfo(sess sessions.Session) gin.H {
	var userContext models.Context
	if sess.Get("uid") != nil {
		userContext = sess.Get("user_context").(models.Context)
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

func Load(c *gin.Context) {
	var req server.RequestRPC
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Set("id", req.ID)
	qwebParams := struct {
		Path string `json:"path"`
	}{}
	if err := json.Unmarshal(req.Params, &qwebParams); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
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
			DateFormat: "%m/%d/%Y",
			Direction: tools.LANG_DIRECTION_LTR,
			ThousandsSep: ",",
			TimeFormat: "%H:%M:%S",
			DecimalPoint: ".",
			ID: 1,
			Grouping: "[]",
		},
		"modules": gin.H{},
	}
	server.RPC(c, http.StatusOK, res)
}

func init() {
	server := server.GetServer()
	server.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/web")
	})
	web := server.Group("/web")
	{
		web.Static("/static", "server/static/web/")
		web.GET("/", WebClient)
		web.GET("/binary/company_logo", CompanyLogo)

		sess := web.Group("/session")
		{
			sess.POST("/get_session_info", GetSessionInfo)
		}

		proxy := web.Group("/proxy")
		{
			proxy.POST("/load", Load)
		}

		webClient := web.Group("/webclient")
		{
			webClient.GET("/qweb", QWeb)
			webClient.POST("/bootstrap_translations", BootstrapTranslations)
		}
	}
}
