package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/mathnogueira/go-arch/model"
	"github.com/mathnogueira/go-arch/scanner"
)

func (p *Project) Scan() ([]model.Module, error) {
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
				referencedModule.UsedBy.Add(module.Name, module)
			}
		}
	}

	modules := make([]model.Module, 0)
	for _, module := range scannedModules {
		modules = append(modules, *module)
	}

	return modules, nil
}

func (p *Project) scanDirectory(path string, scannedModules map[string]*model.Module) error {
	if _, found := p.scannedDirectories[path]; found {
		fmt.Printf("\t\t%s was skipped\n", path)
		return nil
	}

	fmt.Printf("Scanning directory %s\n", path)
	modules, err := scanner.ScanDirectory(path)
	if err != nil {
		return fmt.Errorf("could not find any module in the project root: %w", err)
	}

	wd, _ := os.Getwd()

	p.scannedDirectories[path] = true

	for _, module := range modules {
		scannedModules[path] = &module

		for _, importedModule := range module.Imports {
			if !strings.HasPrefix(importedModule.Path, p.ModuleName) {
				continue
			}

			modulePath := strings.ReplaceAll(importedModule.Path, p.ModuleName, wd)
			fmt.Printf("\tImports %s\n", modulePath)
		}

		for _, importedModule := range module.Imports {
			if !strings.HasPrefix(importedModule.Path, p.ModuleName) {
				continue
			}

			modulePath := strings.ReplaceAll(importedModule.Path, p.ModuleName, wd)
			p.scanDirectory(modulePath, scannedModules)
		}
	}

	return nil
}
