package main

import (
	"fmt"
	"os"

	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/project"
	"github.com/mathnogueira/go-arch/render"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get current working directory:", err.Error())
		os.Exit(1)
	}

	cfg, err := config.Load("./arch.yaml")
	if err != nil {
		fmt.Println("Cannot load the './arch.yaml' file:", err.Error())
		os.Exit(1)
	}

	tagEnricher, err := project.NewModuleEnricher(cfg.Modules)
	if err != nil {
		panic(err)
	}

	project, err := project.NewProject(wd)
	if err != nil {
		panic(err)
	}

	modules, err := project.Scan(tagEnricher)
	if err != nil {
		panic(err)
	}

	outputFile := "graph.png"
	if len(os.Args) > 1 {
		outputFile = os.Args[1]
	}

	err = render.GetRenderer(cfg).Render(modules, outputFile)
	if err != nil {
		panic(err)
	}
}
