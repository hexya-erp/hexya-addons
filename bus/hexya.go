//filevalid

package bus

import (
	"github.com/hexya-erp/hexya-base/web/controllers"
	"github.com/hexya-erp/hexya/hexya/server"
)

const MODULE_NAME string = "bus"

func init() {
	server.RegisterModule(&server.Module{
		Name:     MODULE_NAME,
		PreInit:  func() {},
		PostInit: func() {},
	})
	controllers.BackendJS = append(controllers.BackendJS,
		"/static/bus/src/js/bus.js",
	)

}
