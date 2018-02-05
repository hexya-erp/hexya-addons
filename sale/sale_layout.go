// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool/h"
)

func init() {

	h.SaleLayoutCategory().DeclareModel()
	h.SaleLayoutCategory().SetDefaultOrder("Sequence", "ID")

	h.SaleLayoutCategory().AddFields(map[string]models.FieldDefinition{
		"Name":      models.CharField{String: "Name", Required: true, Translate: true},
		"Sequence":  models.IntegerField{String: "Sequence", Required: true, Default: models.DefaultValue(10)},
		"Subtotal":  models.BooleanField{String: "Add subtotal", Default: models.DefaultValue(true)},
		"Pagebreak": models.BooleanField{String: "Add pagebreak"},
	})

}
