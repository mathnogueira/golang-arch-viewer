package render

import (
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/model"
)

type ImageStyler interface {
	StyleNode(module model.Module, node *cgraph.Node)
	StyleArrow(source, target model.Module, edge *cgraph.Edge)
}
