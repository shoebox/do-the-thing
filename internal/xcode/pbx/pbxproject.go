package pbx

import (
	"errors"
)

type PBXProject struct {
	Targets []NativeTarget
}

func (p PBXProject) FindTargetByName(name string) (NativeTarget, error) {
	var err error
	var res NativeTarget
	var found bool
	for _, tgt := range p.Targets {
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
