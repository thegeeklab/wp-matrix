// Command docs-gen generates docs/data/data.yaml from the plugin's flag
// definitions and the long descriptions attached as doc comments above each
// flag literal.
//
// Invoke via:  go generate ./plugin/...
package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/thegeeklab/wp-matrix/plugin"
	plugin_docs "github.com/thegeeklab/wp-plugin-go/v7/docs"
	plugin_template "github.com/thegeeklab/wp-plugin-go/v7/template"
)

const yamlDescriptionIndent = "      "

func main() {
	outputFile := flag.String("output", "", "Output file path")
	sourceFile := flag.String("source", "plugin.go", "Plugin source file to scan for long descriptions")

	flag.Parse()

	if *outputFile == "" {
		log.Fatal("no output file specified")
	}

	p := plugin.New(nil)
	templateData := plugin_docs.GetTemplateDataWithSource(p.App, *sourceFile)
	descriptions := plugin_docs.LongDescriptionsFor(*sourceFile, defaultMatchers()...)

	funcs := plugin_template.LoadFuncMap()
	funcs["longDesc"] = plugin_docs.LongDescriptionFunc(descriptions, plugin_docs.ShortDescriptionFallback)
	funcs["yamlLiteral"] = yamlLiteral

	docTemplate, err := template.New("docs").Funcs(funcs).Parse(docsTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := docTemplate.Execute(&buf, templateData); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(*outputFile, buf.Bytes(), 0o600); err != nil {
		log.Fatal(err)
	}
}

// defaultMatchers returns the flag-type matchers used to extract long
// descriptions. The wp-plugin-go custom map flag types are included in
// addition to the urfave core flag types.
func defaultMatchers() []plugin_docs.FlagTypeMatcher {
	return []plugin_docs.FlagTypeMatcher{
		plugin_docs.DefaultFlagTypeMatcher,
	}
}

// yamlLiteral renders a LongDescription as the body of a YAML literal block
// scalar at the indent depth that matches the "description: |" line in
// docsTemplate.
func yamlLiteral(d *plugin_docs.LongDescription) string {
	return plugin_docs.LongDescriptionYAMLBlock(d, yamlDescriptionIndent)
}

const docsTemplate = `---
{{- if .GlobalArgs }}
properties:
{{- range $v := .GlobalArgs }}
  - name: {{ $v.Name }}
    {{- with longDesc $v }}
    description: |
{{ yamlLiteral . }}
    {{- end }}
    {{- with $v.Type }}
    type: {{ . }}
    {{- end }}
    {{- with $v.Default }}
    defaultValue: {{ . }}
    {{- end }}
    required: {{ default false $v.Required }}
{{ end -}}
{{ end -}}
`
