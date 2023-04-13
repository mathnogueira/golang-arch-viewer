package render

import (
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/model"
)

type DefaultImageStyler struct{}

func (s DefaultImageStyler) StyleNode(module model.Module, node *cgraph.Node) {
	node.SetStyle(cgraph.FilledNodeStyle)

	if module.Type == "core" {
		node.SetFillColor("#c9e6be")
	} else if module.Type == "infrastructure" {
		node.SetFillColor("#a1bbe3")
	} else if module.Type == "mixed" {
		node.SetFillColor("#de957f")
	}

	node.SetGroup(module.Group)
}

func (s DefaultImageStyler) StyleArrow(source, target model.Module, edge *cgraph.Edge) {
	edge.SetArrowTail(cgraph.NoneArrow)
	edge.SetArrowHead(cgraph.NormalArrow)
}
