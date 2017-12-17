// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.Company().AddFields(map[string]models.FieldDefinition{
		"SaleNote": models.TextField{String: "Default Terms and Conditions", Translate: true},
	})

}
