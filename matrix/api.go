package matrix

import (
	"context"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

//nolint:lll
type APIClient interface {
	SendMessageEvent(ctx context.Context, roomID id.RoomID, eventType event.Type, contentJSON interface{}, extra ...mautrix.ReqSendEvent) (resp *mautrix.RespSendEvent, err error)
}
