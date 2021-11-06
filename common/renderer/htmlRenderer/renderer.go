package htmlrenderer

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

type renderer struct {
	template *template.Template
}

func New(path string) (r *renderer, err error) {
	t, err := template.ParseGlob(filepath.Join(path, "*.tmpl"))
	if err != nil {
		return &renderer{}, err
	}

	return &renderer{
		template: t,
	}, nil
}

func (r *renderer) Render(templateName string, w io.Writer, data interface{}) (err error) {
	t := r.template.Lookup(templateName)
	if t == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
