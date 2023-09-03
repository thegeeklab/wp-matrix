// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/thegeeklab/wp-matrix/plugin"
	"github.com/urfave/cli/v2"
)

//nolint:lll
const defaultTemplate = `
Status: **{{ .Pipeline.Status }}**<br/>
Build: [{{ .Repository.Slug }}]({{ .Pipeline.URL }}){{ if .Curr.Branch }} ({{ .Curr.Branch }}){{ end }} by {{ .Curr.Author.Name }}<br/>
Message: {{ .Curr.Title }}{{ if .Curr.URL }} ([source]({{ .Curr.URL }})){{ end }}
`

// settingsFlags has the cli.Flags for the plugin.Settings.
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
			EnvVars:     []string{"PLUGIN_USERID,PLUGIN_USER_ID", "MATRIX_USERID", "MATRIX_USER_ID"},
			Usage:       "authentication user id",
			Destination: &settings.UserID,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "accesstoken",
			EnvVars:     []string{"PLUGIN_ACCESSTOKEN,PLUGIN_ACCESS_TOKEN", "MATRIX_ACCESSTOKEN", "MATRIX_ACCESS_TOKEN"},
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
			Usage:       "message template",
			Value:       defaultTemplate,
			Destination: &settings.Template,
			Category:    category,
		},
	}
}
