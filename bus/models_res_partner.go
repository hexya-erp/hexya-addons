//filetovalid
package bus

import (
	"fmt"

	"OdooToHexya/output2/typedef"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

func init() {

	h.Partner().AddFields(map[string]models.FieldDefinition{
		"ImStatus": models.CharField{
			String:  "IM Status",
			Compute: h.Partner().Methods().ComputeImStatus(),
		},
	})

	/*VALID*/
	h.Partner().Methods().ComputeImStatus().DeclareMethod(
		`ComputeImStatus`,
		func(rs h.PartnerSet) *h.PartnerData {
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
			return &h.PartnerData{ImStatus: res}
		})

	/*TOVALID*/
	h.Partner().Methods().ImSearch().DeclareMethod(
		`	Search partner with a name and return its id, name and im_status.
				Note : the user must be logged
            	:param name : the partner name to search
            	:param limit : the limit of result to return`,
		func(rs h.PartnerSet, name string, limit int) []h.PartnerData {
			// This method is supposed to be used only in the context of channel creation or
			// extension via an invite. As both of these actions require the 'create' access
			// right, we check this specific ACL.

			// h.MailChannel should not exist -> condition removed
			name = "%" + name + "%"
			excluded := rs.User().Partner().Ids() // not sure
			var res []models.FieldMap
			rs.Env().Cr().Select(&res, `
                SELECT
                    U.id as user_id,
                    P.id as id,
                    P.name as name,
                    CASE WHEN B.last_poll IS NULL THEN 'offline'
                         WHEN age(now() AT TIME ZONE 'UTC', B.last_poll) > interval ? THEN 'offline'
                         WHEN age(now() AT TIME ZONE 'UTC', B.last_presence) > interval ? THEN 'away'
                         ELSE 'online'
                    END as im_status
                FROM user U
                    JOIN partner P ON P.id = U.partner_id
                    LEFT JOIN bus_presence B ON B.user_id = U.id
                WHERE P.name ILIKE ?
                    AND P.id NOT IN ?
                    AND U.active = 't'
                LIMIT ?
            `,
				fmt.Sprintf("%s seconds", busTypes.DisconnectionTimer),
				fmt.Sprintf("%s seconds", busTypes.AwayTimer),
				name, excluded, limit)
			var out []h.PartnerData
			for _, r := range res {
				out = append(out, h.PartnerData{
					ID:       r["id"].(int64),
					User:     h.User().Browse(rs.Env(), r["user_id"].([]int64)),
					Name:     r["name"].(string),
					ImStatus: r["im_status"].(string),
				})
			}
			return out
		})

}
