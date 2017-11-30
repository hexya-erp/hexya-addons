package procurement

import (
	"github.com/hexya-erp/hexya/hexya/actions"
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/pool"
)

func init() {

	pool.ProcurementOrderComputeAll().DeclareTransientModel()

	pool.ProcurementOrderComputeAll().Methods().ProcureCalculationAll().DeclareMethod(
		`ProcureCalculationAll`,
		func(rs pool.ProcurementOrderComputeAllSet) {
			models.ExecuteInNewEnvironment(rs.Env().Uid(), func(env models.Environment) {
				// TODO Avoid to run the scheduler multiple times in the same time
				companies := pool.User().NewSet(env).CurrentUser().Companies()
				for _, company := range companies.Records() {
					pool.ProcurementOrder().NewSet(env).RunScheduler(true, company)
				}
			})
		})

	pool.ProcurementOrderComputeAll().Methods().ProcureCalculation().DeclareMethod(
		`ProcureCalculation`,
		func(rs pool.ProcurementOrderComputeAllSet) *actions.Action {
			go rs.ProcureCalculationAll()
			return &actions.Action{
				Type: actions.ActionCloseWindow,
			}
		})

}
