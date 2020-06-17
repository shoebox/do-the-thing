package pbx

import "errors"

type NativeTarget struct {
	BuildConfigurationList XCConfigurationList
	BuildPhases            []PBXBuildPhase
	// Dependencies           []PBXTargetDependency
	Name               string
	ProductInstallPath string
	ProductName        string
	// productReference       PBXFileReference
	ProductType PBXProductType
}

func FindTargetByName(list []NativeTarget, name string) (NativeTarget, error) {
	var err error
	var res NativeTarget
	var found bool
	for _, tgt := range list {
		if tgt.Name == name {
			found = true
			res = tgt
			break
		}
	}

	if !found {
		err = errors.New("Missing target")
	}

	return res, err
}
