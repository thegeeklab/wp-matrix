// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"github.com/thegeeklab/wp-plugin-go/v2/template"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

var ErrAuthSourceNotSet = errors.New("either username and password or userid and accesstoken are required")

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if (p.Settings.Username == "" || p.Settings.Password == "") &&
		(p.Settings.UserID == "" || p.Settings.AccessToken == "") {
		return ErrAuthSourceNotSet
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	muid := id.NewUserID(prepend("@", p.Settings.UserID), p.Settings.Homeserver)

	client, err := mautrix.NewClient(p.Settings.Homeserver, muid, p.Settings.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	if p.Settings.UserID == "" || p.Settings.AccessToken == "" {
		_, err := client.Login(
			p.Network.Context,
			&mautrix.ReqLogin{
				Type:                     "m.login.password",
				Identifier:               mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: p.Settings.Username},
				Password:                 p.Settings.Password,
				InitialDeviceDisplayName: "Woodpecker CI",
				StoreCredentials:         true,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to authenticate user: %w", err)
		}
	}

	log.Info().Msg("logged in successfully")

	joined, err := client.JoinRoom(p.Network.Context, prepend("!", p.Settings.RoomID), "", nil)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	content, err := p.messageContent(p.Network.Context, *p.Network.Client)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	if _, err := client.SendMessageEvent(p.Network.Context, joined.RoomID, event.EventMessage, content); err != nil {
		return fmt.Errorf("failed to submit message: %w", err)
	}

	log.Info().Msg("message sent successfully")

	return nil
}

func (p *Plugin) messageContent(ctx context.Context, client http.Client) (event.MessageEventContent, error) {
	message, err := template.RenderTrim(ctx, client, p.Settings.Template, p.Metadata)
	if err != nil {
		return event.MessageEventContent{}, err
	}

	content := format.RenderMarkdown(message, true, p.Settings.TemplateUnsafe)

	safeBody := format.HTMLToMarkdown(bluemonday.UGCPolicy().Sanitize(content.FormattedBody))
	if content.Body != safeBody {
		content.Body = safeBody
	}

	if content.FormattedBody != "" {
		content.FormattedBody = bluemonday.UGCPolicy().Sanitize(content.FormattedBody)
	}

	return content, nil
}

func prepend(prefix, input string) string {
	if strings.TrimSpace(input) == "" {
		return input
	}

	if strings.HasPrefix(input, prefix) {
		return input
	}

	return prefix + input
}
