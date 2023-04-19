package render

import (
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/model"
)

type ImageRenderer struct {
	nodes        map[string]*cgraph.Node
	groups       map[string]*cgraph.Graph
	clusterEdges map[string]*cgraph.Edge
	styler       ImageStyler
}

func (ir *ImageRenderer) Render(modules []model.Module) error {
	g := graphviz.New()
	graph, err := g.Graph(graphviz.Directed, graphviz.StrictDirected)
	if err != nil {
		return err
	}

	defer func() {
		graph.Close()
		g.Close()
	}()

	graph.SetCompound(true)
	graph.SetNodeSeparator(1)
	graph.SetRankSeparator(2)
	graph.SetSplines("ortho")
	graph.SetOutputOrder(cgraph.BreadthFirst)
	graph.SetOrdering(cgraph.InOrdering)

	ir.nodes = make(map[string]*cgraph.Node, len(modules))
	ir.groups = make(map[string]*cgraph.Graph, 0)
	ir.clusterEdges = make(map[string]*cgraph.Edge, 0)

	for _, module := range modules {
		ir.renderModule(graph, module)
	}

	for _, module := range modules {
		for _, usedModule := range module.UsedBy.List() {
			ir.renderDependency(graph, module, *usedModule)
		}
	}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (ir *ImageRenderer) renderModule(graph *cgraph.Graph, module model.Module) error {
	cluster := ir.getGroupCluster(graph, module)
	node, _ := cluster.CreateNode(module.Name)
	ir.styler.StyleNode(module, node)
	ir.nodes[module.Name] = node

	return nil
}

func (ir *ImageRenderer) getGroupCluster(graph *cgraph.Graph, module model.Module) *cgraph.Graph {
	if clusterGraph, exists := ir.groups[module.Group]; exists {
		return clusterGraph
	}

	clusterGraph := graph.SubGraph(fmt.Sprintf("cluster_%s", module.Group), 1)

	ir.styler.StyleCluster(module.Group, clusterGraph)
	ir.groups[module.Group] = clusterGraph

	return clusterGraph
}

func (ir *ImageRenderer) renderDependency(graph *cgraph.Graph, target, source model.Module) error {
	targetNode := ir.nodes[target.Name]
	sourceNode := ir.nodes[source.Name]
	edgeName := fmt.Sprintf("%s -> %s", source.Group, target.Group)

	if target.Group != source.Group {
		if _, exists := ir.clusterEdges[edgeName]; exists {
			// An edge already exists and we don't need to render a new one
			return nil
		}
	}

	sourceCluster := ir.getGroupCluster(graph, source)
	targetCluster := ir.getGroupCluster(graph, target)

	if target.Group == source.Group {
		edge, err := sourceCluster.CreateEdge("", sourceNode, targetNode)
		if err != nil {
			return err
		}
		ir.styler.StyleArrow(source, target, sourceCluster, targetCluster, edge)
	} else {
		edge, err := graph.CreateEdge("", sourceNode, targetNode)
		if err != nil {
			return err
		}

		ir.styler.StyleArrow(source, target, sourceCluster, targetCluster, edge)
		edge.SetLogicalHead(fmt.Sprintf("cluster_%s", target.Group))
		edge.SetLogicalTail(fmt.Sprintf("cluster_%s", source.Group))
		edge.SetMinLen(2)

		ir.clusterEdges[edgeName] = edge
	}

	return nil
}
