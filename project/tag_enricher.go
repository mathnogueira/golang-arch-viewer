package project

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mathnogueira/go-arch/config"
	"github.com/mathnogueira/go-arch/model"
)

type ModuleSpecEnricher struct {
	spec map[*regexp.Regexp]config.ModuleSpec
}

func NewModuleEnricher(modules map[string]config.ModuleSpec) (ModuleEnricher, error) {
	processedModules := make(map[*regexp.Regexp]config.ModuleSpec, len(modules))
	for query, module := range modules {
		if strings.HasSuffix(query, "/*") {
			query = strings.ReplaceAll(query, "/*", "/?(.*)")
		}

		regexQuery := strings.ReplaceAll(query, "/", "\\/")
		regexQuery = fmt.Sprintf("^%s$", regexQuery)
		regex, err := regexp.Compile(regexQuery)
		if err != nil {
			return nil, fmt.Errorf("could not convert tag query to regex: (%s: %s): %w", query, regexQuery, err)
		}

		processedModules[regex] = module
	}

	return &ModuleSpecEnricher{processedModules}, nil
}

func (t *ModuleSpecEnricher) Enrich(project *Project, module *model.Module) {
	relativePath := strings.TrimPrefix(module.Directory, project.RootDir)
	relativePath = strings.TrimPrefix(relativePath, "/")

	for regex, spec := range t.spec {
		if regex.Match([]byte(relativePath)) {
			module.Type = spec.Type
			module.Group = spec.Group
		}
	}
}
