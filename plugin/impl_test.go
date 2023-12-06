package plugin

import (
	"context"
	"net/http"
	"testing"

	wp "github.com/thegeeklab/wp-plugin-go/plugin"
)

func Test_messageContent(t *testing.T) {
	//nolint:lll
	tests := []struct {
		name     string
		want     string
		unsafe   bool
		template string
		meta     wp.Metadata
	}{
		{
			name:     "render default template",
			want:     "Status: **success**\nBuild: [octocat/demo](https://ci.example.com) (main) by octobot\nMessage: feat: demo commit title ([source](https://git.example.com))",
			template: DefaultMessageTemplate,
		},
		{
			name:     "render unsafe html template",
			want:     "Status: **success**\nBuild: octocat/demo",
			unsafe:   true,
			template: "Status: **{{ .Pipeline.Status }}**<br>Build: {{ .Repository.Slug }}",
		},
		{
			name:     "render html xss template",
			want:     "Status: **success**\nBuild: octocat/demo",
			unsafe:   true,
			template: "Status: **{{ .Pipeline.Status }}**<br>Build: <a href=\"javascript:alert('XSS1')\" onmouseover=\"alert('XSS2')\">{{ .Repository.Slug }}<a>",
		},
		{
			name:     "render markdown xss template",
			want:     "Status: **success**\nBuild: octocat/demo",
			unsafe:   true,
			template: "Status: **{{ .Pipeline.Status }}**<br>Build: [{{ .Repository.Slug }}](javascript:alert(XSS1'))",
		},
	}

	options := wp.Options{
		Name:    "wp-matrix",
		Execute: func(ctx context.Context) error { return nil },
	}

	p := New(options, &Settings{})
	p.Metadata = wp.Metadata{
		Curr: wp.Commit{
			Branch: "main",
			Title:  "feat: demo commit title",
			URL:    "https://git.example.com",
			Author: wp.Author{
				Name: "octobot",
			},
		},
		Pipeline: wp.Pipeline{
			Status: "success",
			URL:    "https://ci.example.com",
		},
		Repository: wp.Repository{
			Slug: "octocat/demo",
		},
	}

	for _, tt := range tests {
		p.Settings.Template = tt.template
		p.Settings.TemplateUnsafe = tt.unsafe
		content, _ := p.messageContent(context.Background(), http.Client{})

		if content.Body != tt.want {
			t.Errorf("messageContent: %q got: %q, want: %q", tt.name, content.Body, tt.want)
		}
	}
}
