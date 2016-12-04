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

package methods

import (
	"fmt"
	"github.com/npiganeau/yep/pool"
	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/models/types"
)

// UsersNameGet is the NameGet implementation for users
func UsersNameGet(rs pool.ResUsersSet) string {
	res := rs.Super()
	return fmt.Sprintf("%s (%s)", res, rs.Login())
}

// UsersContextGet returns a context with the user's lang, tz and uid
// This method must be called on a singleton.
func UsersContextGet(rs pool.ResUsersSet) *types.Context {
	rs.EnsureOne()
	res := types.NewContext()
	res = res.WithKey("lang", rs.Lang())
	res = res.WithKey("tz", rs.TZ())
	res = res.WithKey("uid", rs.ID())
	return res
}

func initUsers() {
	models.ExtendMethod("ResUsers", "NameGet", UsersNameGet)
	models.CreateMethod("ResUsers", "ContextGet", UsersContextGet)
}
