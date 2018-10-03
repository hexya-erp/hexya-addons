//filetovalid
package bus

import (
	"OdooToHexya/output2/typedef"
	"container/list"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types/dates"
	"github.com/hexya-erp/hexya/pool/h"
	"github.com/hexya-erp/hexya/pool/q"
)

//#----------------------------------------------------------

//# Dispatcher

//#----------------------------------------------------------
/*TOVALID*/
//dispatch = None

//if not odoo.multi_process or odoo.evented:

//    # We only use the event dispatcher in threaded and gevent mode

//    dispatch = ImDispatch()

/*TOVALID*/
//class ImDispatch(object):
/*TOVALID*/
//    def __init__(self):

//        self.channels = {}

//        self.started = False
/*TOVALID*/
//    def poll(self, dbname, channels, last, options=None, timeout=TIMEOUT):

//        if options is None:

//            options = {}

//        # Dont hang ctrl-c for a poll request, we need to bypass private

//        # attribute access because we dont know before starting the thread that

//        # it will handle a longpolling request

//        if not odoo.evented:

//            current = threading.current_thread()

//            current._Thread__daemonic = True

//            # rename the thread to avoid tests waiting for a longpolling

//            current.setName("openerp.longpolling.request.%s" % current.ident)

//        registry = odoo.registry(dbname)

//        # immediatly returns if past notifications exist

//        with registry.cursor() as cr:

//            env = api.Environment(cr, SUPERUSER_ID, {})

//            notifications = env['bus.bus'].poll(channels, last, options)

//        # immediatly returns in peek mode

//        if options.get('peek'):

//            return dict(notifications=notifications, channels=channels)

//        # or wait for future ones

//        if not notifications:

//            if not self.started:

//                # Lazy start of events listener

//                self.start()

//            event = self.Event()

//            for channel in channels:

//                self.channels.setdefault(hashable(channel), []).append(event)

//            try:

//                event.wait(timeout=timeout)

//                with registry.cursor() as cr:

//                    env = api.Environment(cr, SUPERUSER_ID, {})

//                    notifications = env['bus.bus'].poll(

//                        channels, last, options, force_status=True)

//            except Exception:

//                # timeout

//                pass

//        return notifications
/*TOVALID*/
//    def loop(self):

//        """ Dispatch postgres notifications to the relevant polling threads/greenlets """

//        _logger.info("Bus.loop listen imbus on db postgres")

//        with odoo.sql_db.db_connect('postgres').cursor() as cr:

//            conn = cr._cnx

//            cr.execute("listen imbus")

//            cr.commit()

//            while True:

//                if select.select([conn], [], [], TIMEOUT) == ([], [], []):

//                    pass

//                else:

//                    conn.poll()

//                    channels = []

//                    while conn.notifies:

//                        channels.extend(json.loads(

//                            conn.notifies.pop().payload))

//                    # dispatch to local threads/greenlets

//                    events = set()

//                    for channel in channels:

//                        events.update(self.channels.pop(hashable(channel), []))

//                    for event in events:

//                        event.set()
/*TOVALID*/
//    def run(self):

//        while True:

//            try:

//                self.loop()

//            except Exception, e:

//                _logger.exception("Bus.loop error, sleep and retry")

//                time.sleep(TIMEOUT)
/*TOVALID*/
//    def start(self):

//        if odoo.evented:

//            # gevent mode

//            import gevent

//            self.Event = gevent.event.Event

//            gevent.spawn(self.run)

//        else:

//            # threaded mode

//            self.Event = threading.Event

//            t = threading.Thread(name="%s.Bus" % __name__, target=self.run)

//            t.daemon = True

//            t.start()

//        self.started = True

//        return self

func init() {

	h.BusBus().DeclareModel()
	h.BusBus().AddFields(map[string]models.FieldDefinition{
		"Channel": models.CharField{
			String: "Channel",
		},
		"Message": models.CharField{
			String: "Message",
		},
	})

	h.BusBus().Fields().CreateDate().SetString("Create date")

	/*VALID*/
	h.BusBus().Methods().Gc().DeclareMethod(
		`Gc`,
		func(rs h.BusBusSet) int64 {
			timeoutAgo := dates.Now().Add(time.Second * (busTypes.Timeout * -2))
			query := q.BusBus().CreateDate().Lower(timeoutAgo)
			return h.BusBus().NewSet(rs.Env()).Sudo().Search(query).Unlink()
		})

	/*VALID*/
	h.BusBus().Methods().SendMany().DeclareMethod(
		`Sendmany`,
		func(rs h.BusBusSet, notifications []busTypes.BusTypesNotification) {
			channels := list.New()
			for _, notif := range notifications {
				channels.PushBack(notif.Channel)
				values := h.BusBusData{
					Channel: notif.Channel,
					Message: notif.Message,
				}
				rs.Sudo().Create(&values)
				if rand.Float32() < 0.01 {
					rs.Gc()
				}
			}
			/*   TODO NYI
			if channels.Len() > 0 {
				//	We have to wait until the notifications are commited in database.
				//	When calling `NOTIFY imbus`, some concurrent threads will be
				//	awakened and will fetch the notification in the bus table. If the
				//	transaction is not commited yet, there will be nothing to fetch,
				//	and the longpolling will return no notification.

				//  def notify():
				//      with odoo.sql_db.db_connect("postgres").cursor() as cr:
				//          cr.execute("notify imbus, %s",                       (json_dump(list(channels)),))
				//  self._cr.after("commit", notify)
			}*/
		})

	/*VALID*/
	h.BusBus().Methods().SendOne().DeclareMethod(
		`SendOne`,
		func(rs h.BusBusSet, channel, message string) {
			rs.SendMany([]busTypes.BusTypesNotification{{
				Channel: channel,
				Message: message,
			}})
		})

	/*VALID*/
	h.BusBus().Methods().Poll().DeclareMethod(
		`Poll`,
		func(rs h.BusBusSet, channels []string, last int64, options map[string]interface{}, forceStatus bool) []busTypes.BusTypesNotification {
			var domain q.BusBusCondition
			// first poll return the notification in the 'buffer'
			if last == 0 {
				timeoutAgo := dates.Now().Add(time.Second * -busTypes.Timeout)
				domain = q.BusBus().CreateDate().Greater(timeoutAgo)

			} else { // else returns the unread notifications
				domain = q.BusBus().ID().Greater(last)
			}
			domain = domain.And().Channel().In(channels)
			notifications := rs.Sudo().Search(domain)
			// list of notification to return
			var result []busTypes.BusTypesNotification
			for _, notif := range notifications.Records() {
				result = append(result, busTypes.BusTypesNotification{
					Id:      notif.ID(),
					Channel: notif.Channel(),
					Message: notif.Message(),
				})
			}
			if len(result) > 0 || forceStatus {
				partnerIds := options["bus_presence_partner_ids"].([]int64)
				if len(partnerIds) > 0 {
					partners := h.Partner().Browse(rs.Env(), partnerIds)
					for _, partner := range partners.Records() {
						marshall, err := json.Marshal(map[string]interface{}{
							"id":        partner.ID(),
							"im_status": partner.ImStatus(),
						})
						if err != nil {
							panic(err)
						}
						result = append(result, busTypes.BusTypesNotification{
							Id:      -1,
							Channel: "bus.presence",
							Message: string(marshall),
						})
					}
				}
			}
			return result
		})
}
