//filetovalid
package bus

import (
	"OdooToHexya/output2/typedef"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.User().AddFields(map[string]models.FieldDefinition{
		"ImStatus": models.CharField{
			String:  "IM Status",
			Compute: h.User().Methods().ComputeImStatus(),
		},
	})

	/*TOVALID*/
	h.User().Methods().ComputeImStatus().DeclareMethod(
		`ComputeImStatus`,
		func(rs h.PartnerSet) *h.UserData {
			var res string
			pres := h.BusPresence().Search(rs.Env(), q.BusPresence().UserFilteredOn(q.User().Partner().Equals(rs)))
			switch {
			case dates.Now().Sub(pres.LastPoll()) > busTypes.DisconnectionTimer:
				res = "offline"
			case dates.Now().Sub(pres.LastPresence()) > busTypes.AwayTimer:
				res = "away"
			default:
				res = "online"
			}
			return &h.UserData{ImStatus: res}
		})

}
