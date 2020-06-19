package pbx

import "errors"

type XCConfigurationList struct {
	Reference                   string
	BuildConfiguration          []XCBuildConfiguration
	DefaultConfigurationVisible int
	DefaultConfigurationName    string
}

type XCBuildConfiguration struct {
	Name                       string
	BuildSettings              map[string]string
	BaseConfigurationReference string
}

func (xc XCConfigurationList) FindConfiguration(name string) (XCBuildConfiguration, error) {
	var res XCBuildConfiguration
	var err error
	var found bool
	for _, b := range xc.BuildConfiguration {
		if b.Name == name {
			res = b
			found = true
			break
		}
	}

	if !found {
		err = errors.New("Missing configuration")
	}

	return res, err
}
