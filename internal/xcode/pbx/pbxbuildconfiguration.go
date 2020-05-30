package pbx

type XCConfigurationList struct {
	Reference                   string
	BuildConfiguration          []XCBuildConfiguration
	DefaultConfigurationVisible int
	DefaultConfigurationName    string
}

type XCBuildConfiguration struct {
	BuildSettings              map[string]string
	BaseConfigurationReference string
	Name                       string
}
