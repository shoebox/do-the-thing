package pbx

type NativeTarget struct {
	BuildConfigurationList XCConfigurationList
	BuildPhases            []PBXBuildPhase
	Dependencies           []NativeTarget
	Name                   string
	ProductInstallPath     string
	ProductName            string
	// productReference       PBXFileReference
	ProductType PBXProductType
}
