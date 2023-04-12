package main

import (
	"fmt"

	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/project"
	"github.com/mathnogueira/go-arch/render"
)

func main() {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	wd := "/home/matheus/kubeshop/tracetest/server"

	cfg, err := config.Load("./arch.yaml")
	if err != nil {
		panic(err)
	}

	tagEnricher, err := project.NewTagEnricher(cfg.Tags)
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

	for _, module := range modules {
		fmt.Printf("Module %s\n", module.Name)
		for _, usedBy := range module.UsedBy.List() {
			fmt.Printf("\tUsed by %s\n", usedBy.Name)
		}
	}

	err = render.GetRenderer().Render(modules)
	if err != nil {
		panic(err)
	}
}
