package project

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mathnogueira/go-arch/model"
	"github.com/mathnogueira/go-arch/scanner"
)

type ModuleEnricher interface {
	Enrich(project *Project, module *model.Module)
}

func (p *Project) Scan(enrichers ...ModuleEnricher) ([]model.Module, error) {
	p.scannedDirectories = make(map[string]bool)
	scannedModules := make(map[string]*model.Module, 0)
	err := p.scanDirectory(p.RootDir, scannedModules)
	if err != nil {
		return []model.Module{}, err
	}

	for _, module := range scannedModules {
		for _, importedModule := range module.Imports {
			importedModulePath := strings.ReplaceAll(importedModule.Path, p.ModuleName, p.RootDir)
			if referencedModule, found := scannedModules[importedModulePath]; found {
				referencedModule.UsedBy.Add(module)
			}
		}
	}

	modules := make([]model.Module, 0)
	for _, module := range scannedModules {
		for _, enricher := range enrichers {
			enricher.Enrich(p, module)
		}

		modules = append(modules, *module)
	}

	modules = p.ensureNoNameCollision(modules)

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Directory <= modules[j].Directory
	})

	return modules, nil
}

func (p *Project) scanDirectory(path string, scannedModules map[string]*model.Module) error {
	if _, found := p.scannedDirectories[path]; found {
		return nil
	}

	modules, err := scanner.ScanDirectory(path)
	if err != nil {
		return fmt.Errorf("could not find any module in the project root: %w", err)
	}

	p.scannedDirectories[path] = true

	for _, module := range modules {
		scannedModules[path] = &module

		for _, importedModule := range module.Imports {
			if !strings.HasPrefix(importedModule.Path, p.ModuleName) {
				continue
			}

			modulePath := strings.ReplaceAll(importedModule.Path, p.ModuleName, p.RootDir)
			p.scanDirectory(modulePath, scannedModules)
		}
	}

	return nil
}

func (p *Project) ensureNoNameCollision(modules []model.Module) []model.Module {
	namedModules := make(map[string][]model.Module, 0)
	for _, module := range modules {
		namedModules[module.Name] = append(namedModules[module.Name], module)
	}

	newModules := make([]model.Module, 0, len(modules))
	for _, values := range namedModules {
		for _, module := range values {
			if len(values) > 1 {
				// conflict
				path := strings.Split(module.Directory, "/")
				module.Name = strings.Join(path[len(path)-2:], "/")
			}

			newModules = append(newModules, module)
		}
	}

	return newModules
}
