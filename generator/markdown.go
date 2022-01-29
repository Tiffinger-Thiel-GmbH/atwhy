package generator

import (
	"github.com/Tiffinger-Thiel-GmbH/atwhy/template"
	"io"
)

type Markdown struct{}

func (m Markdown) Generate(template template.Markdown, writer io.Writer) error {
	err := template.Execute(writer)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}

func (m Markdown) Ext() string {
	return ".md"
}
