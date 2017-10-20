// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package decimalPrecision

import (
	// decimalPrecision depends on base module
	_ "github.com/hexya-erp/hexya-base/base"
	"github.com/hexya-erp/hexya/hexya/server"
	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
)

const MODULE_NAME string = "decimalPrecision"

var Precisions = map[string]nbutils.Digits{}

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PostInit: func() {},
	})
}

// GetPrecision returns the precision for the given application
func GetPrecision(app string) nbutils.Digits {
	res, exists := Precisions[app]
	if !exists {
		return nbutils.Digits{Precision: 16, Scale: 2}
	}
	return res
}
