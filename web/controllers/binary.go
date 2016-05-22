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
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/npiganeau/yep/yep/server"
	"net/http"
	"strconv"
)

func CompanyLogo(c *gin.Context) {
	c.File("config/img/logo.png")
}

func Image(c *gin.Context) {
	model := c.Query("model")
	field := c.Query("field")
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	sess := sessions.Default(c)
	uid := sess.Get("uid").(int64)
	img, gErr := server.GetFieldValue(uid, id, model, field)
	res, err := base64.StdEncoding.DecodeString(img.(string))
	if err != nil || gErr != nil {
		c.Error(fmt.Errorf("Unable to fetch image"))
		return
	}
	c.Data(http.StatusOK, "image/png", res)
}
