package render

import (
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/model"
)

type DefaultImageStyler struct{}

func (s DefaultImageStyler) StyleNode(module model.Module, node *cgraph.Node) {
	fillColor := "#e0e0e0"

	if module.HasTag("core") {
		fillColor = "#c9e6be"
	} else if module.HasTag("infrastructure") {
		fillColor = "#a1bbe3"
	}

	node.SetStyle(cgraph.FilledNodeStyle)
	node.SetFillColor(fillColor)
}

func (s DefaultImageStyler) StyleArrow(source, target model.Module, edge *cgraph.Edge) {
	edge.SetArrowTail(cgraph.NoneArrow)
	edge.SetArrowHead(cgraph.NormalArrow)
}
