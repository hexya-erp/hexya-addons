//filevalid
package bus

import (
	"time"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.BusPresence().DeclareModel()

	h.BusPresence().AddFields(map[string]models.FieldDefinition{
		"User": models.Many2OneField{
			RelationModel: h.User(),
			String:        "Users",
			Required:      true,
			Index:         true,
			OnDelete:      models.Cascade,
		},
		"LastPoll": models.DateTimeField{
			String:  "Last Poll",
			Default: models.DefaultValue(dates.Now()),
		},
		"LastPresence": models.DateTimeField{
			String:  "Last Presence",
			Default: models.DefaultValue(dates.Now()),
		},
		"Status": models.SelectionField{
			Selection: types.Selection{
				"online":  "Online",
				"away":    "Away",
				"offline": "Offline",
			},
			String:  "IM Status",
			Default: models.DefaultValue("offline"),
		},
	})

	h.BusPresence().AddSQLConstraint("BusUserPresenceUnique", "unique(User)", "A user can only have one IM status.")

	h.BusPresence().Methods().Update().DeclareMethod(
		`	User Presence
				Its status is 'online', 'away' or 'offline'. This model should be a one2one, but is not
				attached to res_users to avoid database concurrence errors. Since the 'update' method is executed
				at each poll, if the user have multiple opened tabs, concurrence errors can happend, but are 'muted-logged'.
			 `,
		func(rs h.BusPresenceSet, inactivity_period time.Duration) {
			presence := h.BusPresence().NewSet(rs.Env()).Search(q.BusPresence().ID().Equals(rs.Env().Uid())).Limit(1)
			// compute last_presence timestamp
			values := h.BusPresenceData{}
			lastPresence := dates.Now().Add(-inactivity_period)
			values.LastPoll = dates.Now()
			// update the presence or a create a new one
			if presence.IsEmpty() { // create a new presence for the user
				values.User = h.User().Browse(rs.Env(), []int64{rs.Env().Uid()})
				values.LastPresence = lastPresence
				rs.Create(&values)
			} else { // update the last_presence if necessary, and write values
				if presence.LastPresence().Lower(lastPresence) {
					values.LastPresence = lastPresence
				}
				presence.Write(&values)
			}
		})

}
