package pbx

type Entry struct {
	Ref string `plist:"reference"`

	// Common
	FileRef string `plist:"fileRef"`
	Isa     string `plist:"isa"`

	BuildConfigurationList Ref      `plist:"buildConfigurationList"`
	BuildPhases            ArrayRef `plist:"buildPhases"`
	Dependencies           ArrayRef `plist:"dependencies"`
	Name                   string   `plist:"name"`
	Path                   string   `plist:"path"`
	ProductName            string   `plist:"productName"`

	// PBXFileReference
	ExplicitFileType  string `plist:"explicitFileType"`
	LastKnownFileType string `plist:"lastKnownFileType"`
	SourceTree        string `plist:"sourceTree"`

	// PBXFrameworksBuildPhase
	BuildActionMask string `plist:"buildActionMask"`
	Files           []string `plist:"files"`

	// PBXGroup
	Children []string `plist:"children"`

	// PBXNativeTarget
	ProductInstallPath string `plist:"productInstallPath"`
	ProductReference   string `plist:"productReference"`
	ProductType        string `plist:"productType"`

	// PBXBuildFile
	Settings map[string]interface{} `plist:"settings"`

	// PBXProject
	CompatibilityVersion string   `json:"compatibilityVersion"`
	DevelopmentRegion    string   `plist:"developmentRegion"`
	MainGroup            string   `plist:"mainGroup"`
	ProductRefGroup      string   `plist:"productRefGroup"`
	ProjectDirPath       string   `plist:"projectDirPath"`
	ProjectReferences    string   `plist:"projectReferences"`
	Targets              ArrayRef `plist:"targets"`

	// PBXShellScriptBuildPhase
	InputPaths  []string `plist:"inputPaths"`
	OutputPaths []string `plist:"outputPaths"`
	ShellPath   string   `plist:"shellPath"`
	ShellScript string   `plist:"shellScript"`

	// PBXTargetDependency
	TargetProxy string `plist:"targetProxy"`
	Target      Ref    `plist:"target"`

	// XCBuildConfiguration
	BuildSettings map[string]interface{} `plist:"buildSettings"`

	// XCConfigurationList
	DefaultConfigurationName   string   `plist:"defaultConfigurationName"`
	BaseConfigurationReference string   `plist:"baseConfigurationReference"`
	BuildConfigurations        ArrayRef `plist:"buildConfigurations"`
}

type PBXProjRaw struct {
	ArchiveVersion string           `plist:"archiveVersion"`
	Objects        map[string]Entry `plist:"objects"`
	RootObject     string           `plist:"rootObject"`
}

func (p PBXProjRaw) GetRoot() Entry {
	return p.Objects[p.RootObject]
}

type Ref string

func (r Ref) Get(proj PBXProjRaw) Entry {
	return proj.Objects[string(r)]
}

type ArrayRef []Ref

func (a ArrayRef) GetList(p PBXProjRaw) []Entry {
	var res []Entry
	for _, ref := range a {
		d := ref.Get(p)
		d.Ref = string(ref)
		res = append(res, d)
	}

	return res
}
