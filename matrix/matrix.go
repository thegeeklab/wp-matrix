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

// NewClient creates a new Matrix client with the given parameters and joins the specified room.
// It authenticates the user if the userID and token are not provided, and returns a Client struct
// that can be used to send messages to the room.
func NewClient(ctx context.Context, url, roomID, userID, token, username, password string) (*Client, error) {
	muid := id.NewUserID(EnsurePrefix("@", userID), url)

	c, err := mautrix.NewClient(url, muid, token)
	if err != nil {
		return nil, err
	}

	if userID == "" || token == "" {
		_, err := c.Login(
			ctx,
			&mautrix.ReqLogin{
				Type:                     "m.login.password",
				Identifier:               mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: username},
				Password:                 password,
				InitialDeviceDisplayName: "Woodpecker CI",
				StoreCredentials:         true,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate user: %w", err)
		}
	}

	joinResp, err := c.JoinRoom(ctx, EnsurePrefix("!", roomID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to join room: %w", err)
	}

	return &Client{
		client: c,
		Message: &Message{
			client: c,
			Opt: MessageOptions{
				RoomID: joinResp.RoomID,
			},
		},
	}, nil
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
		return err
	}

	return nil
}
