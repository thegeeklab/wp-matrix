//go:build generate
// +build generate

package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/thegeeklab/wp-matrix/plugin"
	"github.com/thegeeklab/wp-plugin-go/docs"
	wp "github.com/thegeeklab/wp-plugin-go/plugin"
	wp_template "github.com/thegeeklab/wp-plugin-go/template"
	"github.com/urfave/cli/v2"
)

//go:embed templates/docs-data.yaml.tmpl
var yamlTemplate embed.FS

func main() {
	settings := &plugin.Settings{}
	app := &cli.App{
		Flags: settingsFlags(settings, wp.FlagsPluginCategory),
	}

	out, err := toYAML(app)
	if err != nil {
		panic(err)
	}

	fi, err := os.Create("../../docs/data/data-raw.yaml")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	if _, err := fi.WriteString(out); err != nil {
		panic(err)
	}
}

func toYAML(app *cli.App) (string, error) {
	var w bytes.Buffer

	yamlTmpl, err := template.New("docs").Funcs(wp_template.LoadFuncMap()).ParseFS(yamlTemplate, "templates/docs-data.yaml.tmpl")
	if err != nil {
		fmt.Println(yamlTmpl)
		return "", err
	}

	if err := yamlTmpl.ExecuteTemplate(&w, "docs-data.yaml.tmpl", docs.GetTemplateData(app)); err != nil {
		return "", err
	}

	return w.String(), nil
}
