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
	"github.com/thegeeklab/wp-matrix/matrix"
	plugin_base "github.com/thegeeklab/wp-plugin-go/v7/plugin"
	plugin_template "github.com/thegeeklab/wp-plugin-go/v7/template"
)

var ErrAuthSourceNotSet = errors.New("either username and password or userid and accesstoken are required")

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
	network, err := p.GetNetwork()
	if err != nil {
		return fmt.Errorf("error while getting network configuration: %w", err)
	}

	metadata, err := p.GetMetadata()
	if err != nil {
		return fmt.Errorf("error while getting metadata: %w", err)
	}

	msg, err := p.CreateMessage(network, metadata)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	client, err := matrix.NewClient(
		network.Context,
		p.Settings.Homeserver,
		p.Settings.RoomID,
		p.Settings.UserID,
		p.Settings.AccessToken,
		p.Settings.Username,
		p.Settings.Password,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}

	client.Message.Opt = matrix.MessageOptions{
		RoomID:         client.Message.Opt.RoomID,
		Message:        msg,
		TemplateUnsafe: p.Settings.TemplateUnsafe,
	}

	if err := client.Message.Send(network.Context); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Info().Msg("message sent successfully")

	return nil
}

// CreateMessage generates a message string based on the plugin's template and metadata.
func (p *Plugin) CreateMessage(network plugin_base.Network, metadata plugin_base.Metadata) (string, error) {
	return plugin_template.RenderTrim(network.Context, *network.Client, p.Settings.Template, metadata)
}
