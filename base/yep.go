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

package base

import (
	_ "github.com/npiganeau/yep-addons/base/controllers"
	"github.com/npiganeau/yep-addons/base/defs"
	"github.com/npiganeau/yep/yep/server"
)

const (
	MODULE_NAME string = "base"
	SEQUENCE    uint8  = 100
	NAME        string = "Base"
	VERSION     string = "0.1"
	CATEGORY    string = "Hidden"
	DESCRIPTION string = `
The kernel of YEP, needed for all installation
==============================================
	`
	AUTHOR     string = "NDP Systèmes"
	MAINTAINER string = "NDP Systèmes"
	WEBSITE    string = "http://www.ndp-systemes.fr"
)

func init() {
	server.RegisterModule(&server.Module{Name: MODULE_NAME, PostInit: PostInit})
}

func PostInit() {
	defs.PostInit()
}
