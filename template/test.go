package template

import "text/template"

func TestTemplate(fields Markdown, templateStr string) Markdown {
	fields.template = template.Must(template.New("test").Parse(templateStr))
	return fields
}
