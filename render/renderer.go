package render

import (
	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/model"
)

type Renderer interface {
	Render(modules []model.Module) error
}

func GetRenderer(config config.Config) Renderer {
	return &ImageRenderer{styler: NewImageStyler(config)}
}
