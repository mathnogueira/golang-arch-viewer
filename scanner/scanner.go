package scanner

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strconv"
	"strings"

	"github.com/mathnogueira/go-arch/model"
)

func ScanDirectory(dirPath string) ([]model.Module, error) {
	fset := token.NewFileSet()
	dirPath = unquoteModulePath(dirPath)
	pkgs, err := parser.ParseDir(fset, dirPath, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ImportsOnly)
	if err != nil {
		return []model.Module{}, fmt.Errorf("coult not parse directory: %w", err)
	}

	modules := make([]model.Module, 0)

	for name, pkg := range pkgs {
		module, err := getModule(name, pkg)
		if err != nil {
			return []model.Module{}, err
		}
		modules = append(modules, module)
	}

	return modules, nil
}

func getModule(name string, pkg *ast.Package) (model.Module, error) {
	symbols := make([]model.Symbol, 0)
	importedPkgs := make([]model.Import, 0)

	for _, file := range pkg.Files {
		for _, object := range file.Scope.Objects {
			symbols = append(symbols, model.Symbol{Name: object.Name, Kind: object.Kind.String()})
		}

		for _, importedPkg := range file.Imports {
			importedPkgs = append(importedPkgs, model.Import{
				Path: unquoteModulePath(importedPkg.Path.Value),
			})
		}
	}

	return model.Module{
		Name:    name,
		Symbols: symbols,
		Imports: importedPkgs,
		UsedBy:  model.NewDependencies(),
	}, nil
}

func unquoteModulePath(in string) string {
	newDirPath, err := strconv.Unquote(in)
	if err == nil {
		in = newDirPath
	}

	return in
}
