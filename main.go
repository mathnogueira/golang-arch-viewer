package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/project"
	"github.com/mathnogueira/go-arch/render"
)

var (
	projectPath string
	outputFile  string
)

func main() {
	wd, err := getWorkDir()
	if err != nil {
		fmt.Println("Cannot get current working directory:", err.Error())
		os.Exit(1)
	}

	cfg, err := config.Load(path.Join(wd, "./arch.yaml"))
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

	err = render.GetRenderer(cfg).Render(modules, outputFile)
	if err != nil {
		panic(err)
	}
}

func getWorkDir() (string, error) {
	if projectPath == "" {
		return os.Getwd()
	}

	return projectPath, nil
}

func init() {
	flag.StringVar(&projectPath, "p", "", "path to the project")
	flag.StringVar(&outputFile, "o", "graph.png", "path to graph file")

	flag.Parse()
}
