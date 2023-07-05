package model

type Module struct {
	Name      string
	Directory string
	Symbols   []Symbol
	Imports   []Import
	UsedBy    Dependencies
	Type      string
	Group     string
}

func (m Module) UniqueName() string {
	return m.Directory
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

func (d *Dependencies) Add(dep *Module) {
	if arr, found := d.dependenciesMap[dep.UniqueName()]; found {
		d.dependenciesMap[dep.UniqueName()] = append(arr, dep)
	} else {
		d.dependenciesMap[dep.UniqueName()] = []*Module{dep}
	}
}

func (d *Dependencies) List() []*Module {
	modules := make([]*Module, 0, len(d.dependenciesMap))
	for _, mods := range d.dependenciesMap {
		modules = append(modules, mods...)
	}

	return modules
}
