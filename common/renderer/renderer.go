package renderer

import "io"

type Renderer interface {
	Render(templateName string, w io.Writer, data interface{}) (err error)
}
