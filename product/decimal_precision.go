package product

import (
	"log"

	"github.com/hexya-erp/hexya/hexya/tools/nbutils"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.DecimalPrecision().Fields().Digits().SetConstraint(pool.DecimalPrecision().Methods().CheckMainCurrencyRounding())

	pool.DecimalPrecision().Methods().CheckMainCurrencyRounding().DeclareMethod(
		`CheckMainCurrencyRounding checks that the precision of the "Account" application
		is less than the rounding factor of the company's currency`,
		func(rs pool.DecimalPrecisionSet) {
			for _, dp := range rs.Records() {
				if dp.Name() != "Account" {
					continue
				}
				currentUser := pool.User().NewSet(rs.Env()).CurrentUser()
				if a, b := nbutils.Compare(currentUser.Company().Currency().Rounding(), float64(10^(-dp.Digits())), nbutils.Digits{Precision: 6}); !a && !b {
					log.Panic(rs.T("You cannot define the decimal precision of 'Account' as greater than the rounding factor of the company's main currency"))
				}
			}
		})

}
