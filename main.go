package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/mathnogueira/go-arch/model"
	"github.com/mathnogueira/go-arch/project"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	project, err := project.NewProject(wd)
	if err != nil {
		panic(err)
	}

	modules, err := project.Scan()
	if err != nil {
		panic(err)
	}

	for _, module := range modules {
		fmt.Printf("Module %s\n", module.Name)
		for _, usedBy := range module.UsedBy.List() {
			fmt.Printf("\tUsed by %s\n", usedBy.Name)
		}
	}

	err = renderGraph(modules)
	if err != nil {
		panic(err)
	}
}

func renderGraph(modules []model.Module) error {
	g := graphviz.New()
	graph, err := g.Graph(graphviz.Directed, graphviz.StrictDirected)
	if err != nil {
		return err
	}

	defer func() {
		graph.Close()
		g.Close()
	}()

	graphNodes := make(map[string]*cgraph.Node, len(modules))

	for _, module := range modules {
		graphNodes[module.Name], _ = graph.CreateNode(module.Name)
	}

	for _, module := range modules {
		for _, usedModule := range module.UsedBy.List() {
			source := graphNodes[usedModule.Name]
			target := graphNodes[module.Name]
			graph.CreateEdge("", source, target)
		}
	}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "graph.png"); err != nil {
		log.Fatal(err)
	}

	return nil
}
