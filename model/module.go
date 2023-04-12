package model

type Module struct {
	Name      string
	Directory string
	Symbols   []Symbol
	Imports   []Import
	UsedBy    Dependencies
	Tags      []string
}

func (m Module) HasTag(tag string) bool {
	for _, item := range m.Tags {
		if item == tag {
			return true
		}
	}

	return false
}

type Import struct {
	Path   string
	Module Module
}

type Symbol struct {
	Kind string
	Name string
}

type Dependencies struct {
	dependenciesMap map[string][]*Module
}

func NewDependencies() Dependencies {
	return Dependencies{dependenciesMap: make(map[string][]*Module)}
}

func (d *Dependencies) Add(moduleName string, dep *Module) {
	if arr, found := d.dependenciesMap[moduleName]; found {
		d.dependenciesMap[moduleName] = append(arr, dep)
	} else {
		d.dependenciesMap[moduleName] = []*Module{dep}
	}
}

func (d *Dependencies) List() []*Module {
	modules := make([]*Module, 0, len(d.dependenciesMap))
	for _, mods := range d.dependenciesMap {
		modules = append(modules, mods...)
	}

	return modules
}
