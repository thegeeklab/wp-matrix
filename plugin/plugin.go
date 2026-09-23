// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"fmt"
	"slices"

	plugin_base "github.com/thegeeklab/wp-plugin-go/v7/plugin"
	"github.com/urfave/cli/v3"
)

//go:generate go run ../hack/docs-gen/main.go -output=../docs/data/data.yaml

//nolint:lll
const DefaultMessageTemplate = `
Status: **{{ .Pipeline.Status }}**
Build: [{{ .Repository.Slug }}]({{ .Pipeline.URL }}){{ if .Curr.Branch }} ({{ .Curr.Branch }}){{ end }} by {{ .Curr.Author.Name }}
Message: {{ .Curr.Title }}{{ if .Curr.URL }} ([source]({{ .Curr.URL }})){{ end }}
`

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
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

func New(e plugin_base.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := plugin_base.Options{
		Name:        "wp-matrix",
		Description: "Send messages to a Matrix room",
		Flags: slices.Concat(
			plugin_base.LoggingFlags(plugin_base.FlagsPluginCategory),
			plugin_base.NetworkFlags(plugin_base.FlagsPluginCategory),
			Flags(p.Settings, plugin_base.FlagsPluginCategory),
		),
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

	p.Plugin = plugin_base.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
	return []cli.Flag{
		// Authentication username.
		//
		// If set, the `password` parameter is required as well.
		&cli.StringFlag{
			Name:        "username",
			Sources:     cli.EnvVars("PLUGIN_USERNAME", "MATRIX_USERNAME"),
			Usage:       "authentication username",
			Destination: &settings.Username,
			Category:    category,
		},
		// Authentication password.
		&cli.StringFlag{
			Name:        "password",
			Sources:     cli.EnvVars("PLUGIN_PASSWORD", "MATRIX_PASSWORD"),
			Usage:       "authentication password",
			Destination: &settings.Password,
			Category:    category,
		},
		// Authentication user ID.
		//
		// If set, the `access_token` parameter is required as well.
		&cli.StringFlag{
			Name:        "userid",
			Sources:     cli.EnvVars("PLUGIN_USER_ID", "PLUGIN_USERID", "MATRIX_USER_ID", "MATRIX_USERID"),
			Usage:       "authentication user ID",
			Destination: &settings.UserID,
			Category:    category,
		},
		// Authentication access token.
		&cli.StringFlag{
			Name:        "accesstoken",
			Sources:     cli.EnvVars("PLUGIN_ACCESS_TOKEN", "PLUGIN_ACCESSTOKEN", "MATRIX_ACCESS_TOKEN", "MATRIX_ACCESSTOKEN"),
			Usage:       "authentication access token",
			Destination: &settings.AccessToken,
			Category:    category,
		},
		// Matrix home server url.
		&cli.StringFlag{
			Name:        "homeserver",
			Sources:     cli.EnvVars("PLUGIN_HOMESERVER", "MATRIX_HOMESERVER"),
			Usage:       "matrix home server url",
			Value:       "https://matrix.org",
			Destination: &settings.Homeserver,
			Category:    category,
		},
		// Room ID to send messages to.
		&cli.StringFlag{
			Name:        "room_id",
			Sources:     cli.EnvVars("PLUGIN_ROOM_ID", "PLUGIN_ROOMID", "MATRIX_ROOMID", "MATRIX_ROOM_ID"),
			Usage:       "room id to send messages to",
			Destination: &settings.RoomID,
			Category:    category,
		},
		// Golang template for the message.
		//
		// The [Metadata struct](https://pkg.go.dev/github.com/thegeeklab/wp-plugin-go/plugin#Metadata) is exposed
		// to the template and all fields can be referenced. To extend the functionality,
		// [sprig functions](https://masterminds.github.io/sprig/) can also be used.
		&cli.StringFlag{
			Name:        "template",
			Sources:     cli.EnvVars("PLUGIN_TEMPLATE", "MATRIX_TEMPLATE"),
			Usage:       "golang template for the message",
			Value:       DefaultMessageTemplate,
			Destination: &settings.Template,
			Category:    category,
		},
		// Render raw HTML and potentially dangerous links in template.
		//
		// By default, raw HTML and potentially dangerous links in the template are not rendered. If inline HTML is used,
		// it may be necessary to activate this option. In such cases, please ensure that the CI configuration files
		// in the git repository are protected against malicious changes.
		&cli.BoolFlag{
			Name:        "template-unsafe",
			Sources:     cli.EnvVars("PLUGIN_TEMPLATE_UNSAFE", "MATRIX_TEMPLATE_UNSAFE"),
			Usage:       "render raw HTML and potentially dangerous links in template",
			Destination: &settings.TemplateUnsafe,
			Category:    category,
		},
	}
}
