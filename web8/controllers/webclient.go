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
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/yep/server"
	"github.com/npiganeau/yep/yep/tools"
	"net/http"
	"strings"
)

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
			Direction:    tools.LangDirectionLTR,
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
