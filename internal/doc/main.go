//go:build generate
// +build generate

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/thegeeklab/wp-matrix/plugin"
	"github.com/thegeeklab/wp-plugin-go/v2/docs"
	"github.com/thegeeklab/wp-plugin-go/v2/template"
)

func main() {
	tmpl := "https://raw.githubusercontent.com/thegeeklab/woodpecker-plugins/main/templates/docs-data.yaml.tmpl"
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	p := plugin.New(nil)

	out, err := template.Render(context.Background(), client, tmpl, docs.GetTemplateData(p.App))
	if err != nil {
		panic(err)
	}

	outputFile := flag.String("output", "", "Output file path")
	flag.Parse()

	if *outputFile == "" {
		panic("no output file specified")
	}

	err = os.WriteFile(*outputFile, []byte(out), 0o644)
	if err != nil {
		panic(err)
	}
}
