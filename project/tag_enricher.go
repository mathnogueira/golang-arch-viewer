package project

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mathnogueira/go-arch/model"
)

type TagEnricher struct {
	tags map[*regexp.Regexp][]string
}

func NewTagEnricher(tags map[string][]string) (ModuleEnricher, error) {
	processedTags := make(map[*regexp.Regexp][]string, len(tags))
	for query, tags := range tags {
		if strings.HasSuffix(query, "/*") {
			query = strings.ReplaceAll(query, "/*", "/?(.*)")
		}

		regexQuery := strings.ReplaceAll(query, "/", "\\/")
		regexQuery = fmt.Sprintf("^%s$", regexQuery)
		regex, err := regexp.Compile(regexQuery)
		if err != nil {
			return nil, fmt.Errorf("could not convert tag query to regex: (%s: %s): %w", query, regexQuery, err)
		}

		processedTags[regex] = tags
	}

	return &TagEnricher{processedTags}, nil
}

func (t *TagEnricher) Enrich(project *Project, module *model.Module) {
	relativePath := strings.TrimPrefix(module.Directory, project.RootDir)
	relativePath = strings.TrimPrefix(relativePath, "/")

	for regex, tags := range t.tags {
		if regex.Match([]byte(relativePath)) {
			module.Tags = tags
		}
	}
}
