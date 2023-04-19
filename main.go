package main

import (
	"os"

	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/project"
	"github.com/mathnogueira/go-arch/render"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	wd = "/home/matheus/kubeshop/tracetest/server"

	cfg, err := config.Load("./arch.yaml")
	if err != nil {
		panic(err)
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

	err = render.GetRenderer(cfg).Render(modules)
	if err != nil {
		panic(err)
	}
}
