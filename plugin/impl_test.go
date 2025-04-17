package plugin

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	plugin_base "github.com/thegeeklab/wp-plugin-go/v6/plugin"
)

func Test_messageContent(t *testing.T) {
	//nolint:lll
	tests := []struct {
		name     string
		want     string
		template string
	}{
		{
			name:     "render default template",
			want:     "Status: **success**\nBuild: [octocat/demo](https://ci.example.com) (main) by octobot\nMessage: feat: demo commit title ([source](https://git.example.com))",
			template: DefaultMessageTemplate,
		},
		{
			name:     "render unsafe html template",
			want:     "Status: **success**\nBuild: octocat/demo",
			template: "Status: **{{ .Pipeline.Status }}**\nBuild: {{ .Repository.Slug }}",
		},
	}

	p := New(func(_ context.Context) error { return nil })
	p.Network = plugin_base.Network{
		Context: t.Context(),
		Client:  &http.Client{},
	}
	p.Metadata = plugin_base.Metadata{
		Curr: plugin_base.Commit{
			Branch: "main",
			Title:  "feat: demo commit title",
			URL:    "https://git.example.com",
			Author: plugin_base.Author{
				Name: "octobot",
			},
		},
		Pipeline: plugin_base.Pipeline{
			Status: "success",
			URL:    "https://ci.example.com",
		},
		Repository: plugin_base.Repository{
			Slug: "octocat/demo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.Settings.Template = tt.template

			content, err := p.CreateMessage()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, content)
		})
	}
}
