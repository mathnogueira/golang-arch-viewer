package render

import (
	"github.com/mathnogueira/go-arch/model"
)

type Renderer interface {
	Render(modules []model.Module) error
}

func GetRenderer() Renderer {
	return &ImageRenderer{styler: DefaultImageStyler{}}
}
