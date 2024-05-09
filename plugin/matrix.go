package plugin

import (
	"context"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

//nolint:lll
type MautrixClient interface {
	SendMessageEvent(ctx context.Context, roomID id.RoomID, eventType event.Type, contentJSON interface{}, extra ...mautrix.ReqSendEvent) (resp *mautrix.RespSendEvent, err error)
}

type MatrixClient struct {
	client  MautrixClient
	Message *MatrixMessage
}

type MatrixMessage struct {
	client MautrixClient
	Opt    MatrixMessageOpt
}

type MatrixMessageOpt struct {
	RoomID         id.RoomID
	Message        string
	TemplateUnsafe bool
}

// NewMatrixClient creates a new MatrixClient instance with the provided mautrix.Client.
func NewMatrixClient(client *mautrix.Client) *MatrixClient {
	return &MatrixClient{
		client: client,
		Message: &MatrixMessage{
			client: client,
			Opt:    MatrixMessageOpt{},
		},
	}
}

// Send sends a message to the specified room. It sanitizes the message content
// to remove potentially unsafe HTML.
func (m *MatrixMessage) Send(ctx context.Context) error {
	content := format.RenderMarkdown(m.Opt.Message, true, m.Opt.TemplateUnsafe)

	if content.FormattedBody != "" {
		content.Body = format.HTMLToMarkdown(bluemonday.UGCPolicy().Sanitize(content.FormattedBody))
		content.FormattedBody = bluemonday.UGCPolicy().Sanitize(content.FormattedBody)
	}

	_, err := m.client.SendMessageEvent(ctx, m.Opt.RoomID, event.EventMessage, content)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
