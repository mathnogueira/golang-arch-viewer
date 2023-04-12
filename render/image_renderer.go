package render

import (
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/model"
)

type ImageRenderer struct {
	nodes  map[string]*cgraph.Node
	styler ImageStyler
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

	ir.nodes = make(map[string]*cgraph.Node, len(modules))

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
	node, _ := graph.CreateNode(module.Name)
	ir.styler.StyleNode(module, node)
	ir.nodes[module.Name] = node

	return nil
}

func (ir *ImageRenderer) renderDependency(graph *cgraph.Graph, source, target model.Module) error {
	sourceNode := ir.nodes[source.Name]
	targetNode := ir.nodes[target.Name]

	edge, err := graph.CreateEdge("", targetNode, sourceNode)
	ir.styler.StyleArrow(source, target, edge)

	return err
}
