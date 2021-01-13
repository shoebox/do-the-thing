package signature

import (
	"dothething/internal/api"
	"errors"
	"fmt"
)

type exportOptionsService struct {
	API *api.API
	cfg *[]api.TargetSignatureConfig
	tgt api.TargetSignatureConfig
}

func NewExportOptionsService(api *api.API) exportOptionsService {
	return exportOptionsService{API: api}
}

func (s exportOptionsService) Compute() (string, error) {
	var res api.ExportOptions
	var err error

	// resoling the target configuration
	s.cfg = s.API.SignatureService.GetConfiguration()
	if err := s.resolveTarget(); err != nil {
		return "", err
	}

	// resolving the method for the target
	if res.Method, err = s.resolveMethod(); err != nil {
		return "", err
	}

	// resovling provisionings
	res.ProvisioningProfile = s.resolveProvisionings()

	fmt.Println("res :::", res)

	return "", nil
}

func (s exportOptionsService) resolveTarget() error {
	for _, e := range *s.cfg {
		if e.TargetName == s.API.Config.Target {
			s.tgt = e
			return nil
		}
	}

	return errors.New("not found")
}

func (s exportOptionsService) resolveMethod() (string, error) {
	for _, e := range *s.cfg {
		if e.TargetName == s.API.Config.Target {
			return s.resolveMethodForProvisioning(e.Config.ProvisioningProfile), nil
		}
	}

	return "", nil
}

func (s exportOptionsService) resolveProvisionings() map[string]string {
	res := map[string]string{}
	for _, e := range *s.cfg {
		res[e.Config.ProvisioningProfile.BundleIdentifier] = e.Config.ProvisioningProfile.UUID
	}
	return res
}

//  # if ProvisionedDevices: !nil & "get-task-allow": true -> development
//  # if ProvisionedDevices: !nil & "get-task-allow": false -> ad-hoc
//  # if ProvisionedDevices: nil & "ProvisionsAllDevices": "true" -> enterprise
//  # if ProvisionedDevices: nil & ProvisionsAllDevices: nil -> app-store
func (s exportOptionsService) resolveMethodForProvisioning(p *api.ProvisioningProfile) string {
	if p.ProvisionedDevices != nil {
		if p.Entitlements.GetAskAllow {
			return "development"
		} else {
			return "ad-hoc"
		}
	} else {
		if p.ProvisionsAllDevices != nil && *p.ProvisionsAllDevices {
			return "enterprise"
		} else {
			return "app-store"
		}
	}
}

func (s exportOptionsService) resolveProvisioningUUID() {
}
