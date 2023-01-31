package project

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type Project struct {
	ModuleName         string
	RootDir            string
	scannedDirectories map[string]bool
}

func NewProject(rootPath string) (Project, error) {
	moduleName, err := getRootModuleName(rootPath)
	if err != nil {
		return Project{}, err
	}

	return Project{
		RootDir:            rootPath,
		ModuleName:         moduleName,
		scannedDirectories: make(map[string]bool),
	}, nil
}

func getRootModuleName(rootPath string) (string, error) {
	goModFileContent, err := ioutil.ReadFile(filepath.Join(rootPath, "go.mod"))
	if err != nil {
		return "", errors.New("go.mod not found")
	}

	re := regexp.MustCompile(`module (.*)`)
	match := re.FindStringSubmatch(string(goModFileContent))
	return match[1], nil
}
