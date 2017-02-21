package nifi

import (
	"text/template"

	"github.com/go-resty/resty"
)

const postTemplateBody = `
{
  "name": "{{.Name}}",
  "description": "{{.Description}}"
}
`

func doPostTemplate(processGroupId string) *resty.Response {
	t := template.New("Process Group Template")
	t, _ = t.Parse(postTemplateBody)
	myResp := Call(restHandler(postProcessGroupTemplate),
		Request{Id: "root", Body: })
}

