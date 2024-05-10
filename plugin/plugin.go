// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"fmt"

	wp "github.com/thegeeklab/wp-plugin-go/v2/plugin"
	"github.com/urfave/cli/v2"
)

//go:generate mockery
//go:generate go run ../internal/doc/main.go -output=../docs/data/data-raw.yaml

//nolint:lll
const DefaultMessageTemplate = `
Status: **{{ .Pipeline.Status }}**
Build: [{{ .Repository.Slug }}]({{ .Pipeline.URL }}){{ if .Curr.Branch }} ({{ .Curr.Branch }}){{ end }} by {{ .Curr.Author.Name }}
Message: {{ .Curr.Title }}{{ if .Curr.URL }} ([source]({{ .Curr.URL }})){{ end }}
`

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
	Settings *Settings
}

// Settings for the plugin.
type Settings struct {
	Username       string
	Password       string
	UserID         string
	AccessToken    string
	Homeserver     string
	RoomID         string
	Template       string
	TemplateUnsafe bool
}

func New(e wp.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := wp.Options{
		Name:                "wp-matrix",
		Description:         "Send messages to a Matrix room",
		Flags:               Flags(p.Settings, wp.FlagsPluginCategory),
		Execute:             p.run,
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e != nil {
		options.Execute = e
	}

	p.Plugin = wp.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
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
			Value:       DefaultMessageTemplate,
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
