package nifi

import (
	"text/template"
	"bytes"

	"github.com/ottom8/hadoop-ottom8r/logger"
)

const tmplPostPGTemplateBody = `
{
  "name": "{{.Name}}",
  "description": "{{.Description}}",
  "snippetId": "{{.SnippetId}}"
}
`

type postPGTemplateBody struct {
	Name        string
	Description string
	SnippetId   float64
}

func doProcessGroupTemplate(postBody *postPGTemplateBody) string {
	var out bytes.Buffer

	tmpl, _ := template.New("Process Group Template").Parse(tmplPostPGTemplateBody)
	if err := tmpl.Execute(&out, postBody); err  != nil {
		logger.Fatal(err.Error())
	}
	return out.String()
}

