// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProductPriceListWizard().DeclareTransientModel()

	pool.ProductPriceListWizard().AddFields(map[string]models.FieldDefinition{
		"PriceList": models.Many2OneField{RelationModel: pool.ProductPricelist(), Required: true},
		"Qty1":      models.IntegerField{String: "Quantity-1", Default: models.DefaultValue(1)},
		"Qty2":      models.IntegerField{String: "Quantity-2", Default: models.DefaultValue(5)},
		"Qty3":      models.IntegerField{String: "Quantity-3", Default: models.DefaultValue(10)},
		"Qty4":      models.IntegerField{String: "Quantity-4", Default: models.DefaultValue(0)},
		"Qty5":      models.IntegerField{String: "Quantity-5", Default: models.DefaultValue(0)},
	})

	pool.ProductPriceListWizard().Methods().PrintReport().DeclareMethod(
		`PrintReport returns the report action from the data in this popup (not implemented)`,
		func(rs pool.ProductPriceListWizardSet) *actions.Action {
			// TODO implement with reports
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

}
