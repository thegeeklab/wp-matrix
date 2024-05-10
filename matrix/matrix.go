package matrix

import (
	"context"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

type Client struct {
	client  APIClient
	Message *Message
}

type Message struct {
	client APIClient
	Opt    MessageOptions
}

type MessageOptions struct {
	RoomID         id.RoomID
	Message        string
	TemplateUnsafe bool
}

// NewClient creates a new Client instance with the provided mautrix.Client.
func NewClient(client *mautrix.Client) *Client {
	return &Client{
		client: client,
		Message: &Message{
			client: client,
			Opt:    MessageOptions{},
		},
	}
}

// Send sends a message to the specified room. It sanitizes the message content
// to remove potentially unsafe HTML.
func (m *Message) Send(ctx context.Context) error {
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
