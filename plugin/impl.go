// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/thegeeklab/wp-plugin-go/v2/template"
	"maunium.net/go/mautrix"
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
	muid := id.NewUserID(EnsurePrefix("@", p.Settings.UserID), p.Settings.Homeserver)

	matrix, err := mautrix.NewClient(p.Settings.Homeserver, muid, p.Settings.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	if p.Settings.UserID == "" || p.Settings.AccessToken == "" {
		_, err := matrix.Login(
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

	joinResp, err := matrix.JoinRoom(p.Network.Context, EnsurePrefix("!", p.Settings.RoomID), "", nil)
	if err != nil {
		return fmt.Errorf("failed to join room: %w", err)
	}

	msg, err := p.CreateMessage()
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	client := NewMatrixClient(matrix)
	client.Message.Opt = MatrixMessageOpt{
		RoomID:         joinResp.RoomID,
		Message:        msg,
		TemplateUnsafe: p.Settings.TemplateUnsafe,
	}

	if err := client.Message.Send(p.Network.Context); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Info().Msg("message sent successfully")

	return nil
}

// CreateMessage generates a message string based on the plugin's template and metadata.
func (p *Plugin) CreateMessage() (string, error) {
	return template.RenderTrim(p.Network.Context, *p.Network.Client, p.Settings.Template, p.Metadata)
}
