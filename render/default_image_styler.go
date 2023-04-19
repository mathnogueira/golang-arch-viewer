package render

import (
	"fmt"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/model"
)

type DefaultImageStyler struct {
	styleSpec      map[string]config.StyleSpec
	clustersConfig map[string]config.ClusterSpec
}

func NewImageStyler(config config.Config) ImageStyler {
	return &DefaultImageStyler{config.Style, config.Clusters}
}

func (s DefaultImageStyler) StyleNode(module model.Module, node *cgraph.Node) {
	node.SetStyle(cgraph.FilledNodeStyle)
	color := s.styleSpec[module.Type].Color

	node.SetFillColor(color)
	node.SetGroup(module.Group)
}

func (s DefaultImageStyler) StyleArrow(source, target model.Module, sourceCluster, targetCluster *cgraph.Graph, edge *cgraph.Edge) {
	edge.SetArrowTail(cgraph.NoneArrow)
	edge.SetArrowHead(cgraph.NormalArrow)

	if source.Group != target.Group {
		sourceClusterConfig := s.clustersConfig[source.Group]
		edge.SetColor(sourceClusterConfig.Color)
	}
}

func (s DefaultImageStyler) StyleCluster(name string, cluster *cgraph.Graph) {
	clusterConfig := s.clustersConfig[name]
	cluster.SetStyle(cgraph.FilledGraphStyle)
	cluster.SetBackgroundColor(clusterConfig.Color)
	cluster.SetLabel(fmt.Sprintf("[%s]", name))
	cluster.SetFontSize(15)
}
