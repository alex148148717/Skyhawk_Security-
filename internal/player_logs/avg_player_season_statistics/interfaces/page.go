package interfaces

import (
	_ "embed"
	"html/template"
)

//go:embed template.html
var pageTemplateByte string

type PageTemplate struct {
	pageTpl *template.Template
}

func NewPageTemplate() (*PageTemplate, error) {

	var err error
	pageTpl, err := template.New("playerPage").Parse(pageTemplateByte)
	if err != nil {
		return nil, err
	}
	return &PageTemplate{pageTpl: pageTpl}, nil
}
