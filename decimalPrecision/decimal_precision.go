package decimalPrecision

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {
	pool.DecimalPrecision().DeclareModel()
	pool.DecimalPrecision().AddFields(map[string]models.FieldDefinition{
		"Name":   models.CharField{String: "Usage", Index: true, Required: true, Unique: true},
		"Digits": models.IntegerField{String: "Digits", Required: true, Default: models.DefaultValue(2), GoType: new(int8)},
	})

	pool.DecimalPrecision().Methods().PrecisionGet().DeclareMethod(
		`PrecisionGet returns the number of digits for the given application.`,
		func(rs pool.DecimalPrecisionSet, application string) int8 {
			dp := pool.DecimalPrecision().Search(rs.Env(), pool.DecimalPrecision().Name().Equals(application))
			if dp.IsEmpty() {
				return 2
			}
			return dp.Digits()
		})
}
