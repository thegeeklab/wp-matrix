// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/thegeeklab/wp-matrix/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
//
//go:generate go run docs.go flags.go
func settingsFlags(settings *plugin.Settings, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "username",
			EnvVars:     []string{"PLUGIN_USERNAME", "MATRIX_USERNAME"},
			Usage:       "authentication username",
			Destination: &settings.Username,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "password",
			EnvVars:     []string{"PLUGIN_PASSWORD", "MATRIX_PASSWORD"},
			Usage:       "authentication password",
			Destination: &settings.Password,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "userid",
			EnvVars:     []string{"PLUGIN_USER_ID", "PLUGIN_USERID", "MATRIX_USER_ID", "MATRIX_USERID"},
			Usage:       "authentication user ID",
			Destination: &settings.UserID,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "accesstoken",
			EnvVars:     []string{"PLUGIN_ACCESS_TOKEN", "PLUGIN_ACCESSTOKEN", "MATRIX_ACCESS_TOKEN", "MATRIX_ACCESSTOKEN"},
			Usage:       "authentication access token",
			Destination: &settings.AccessToken,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "homeserver",
			EnvVars:     []string{"PLUGIN_HOMESERVER", "MATRIX_HOMESERVER"},
			Usage:       "matrix home server url",
			Value:       "https://matrix.org",
			Destination: &settings.Homeserver,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "roomid",
			EnvVars:     []string{"PLUGIN_ROOMID", "MATRIX_ROOMID"},
			Usage:       "roomid to send messages to",
			Destination: &settings.RoomID,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "template",
			EnvVars:     []string{"PLUGIN_TEMPLATE", "MATRIX_TEMPLATE"},
			Usage:       "golang template for the message",
			Value:       plugin.DefaultMessageTemplate,
			Destination: &settings.Template,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "template-unsafe",
			EnvVars:     []string{"PLUGIN_TEMPLATE_UNSAFE", "MATRIX_TEMPLATE_UNSAFE"},
			Usage:       "render raw HTML and potentially dangerous links in template",
			Destination: &settings.TemplateUnsafe,
			Category:    category,
		},
	}
}
