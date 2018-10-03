package busTypes

import "time"

var Timeout = time.Second * 50
var DisconnectionTimer = Timeout + time.Second*5
var AwayTimer = time.Second * 1800 // 30 minutes

type BusTypesNotification struct {
	Id      int64
	Channel string
	Message string
}
