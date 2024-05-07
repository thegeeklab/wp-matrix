// Copyright (c) 2020, the Drone Plugins project authors.
// Copyright (c) 2021, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	wp "github.com/thegeeklab/wp-plugin-go/v2/plugin"
)

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

func New(options wp.Options, settings *Settings) *Plugin {
	p := &Plugin{}

	options.Execute = p.run

	p.Plugin = wp.New(options)
	p.Settings = settings

	return p
}
