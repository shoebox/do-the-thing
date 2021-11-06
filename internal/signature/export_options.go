package signature

import (
	"bytes"
	"dothething/internal/api"
	"errors"
	"os"
	"path/filepath"

	"howett.net/plist"
)

type exportOptionsService struct {
	API *api.API
	cfg *[]api.TargetSignatureConfig
	tgt api.TargetSignatureConfig
}

func NewExportOptionsService(api *api.API) exportOptionsService {
	return exportOptionsService{API: api}
}

func (s exportOptionsService) createExportOptions() (*api.ExportOptions, error) {
	// resolve xcode target
	tgt, err := s.resolveTarget()
	if err != nil {
		return nil, NewSignatureError(err, ErrorExportOptions)
	}

	// resulting export options struct
	var res = api.ExportOptions{
		SigningStyle:        "manual",
		SigningCertificate:  tgt.Config.Cert.Issuer.CommonName,
		ProvisioningProfile: s.resolveProvisionings(),
		TeamID:              tgt.Config.ProvisioningProfile.Entitlements.TeamID,
	}

	return &res, err
}

func (s exportOptionsService) Compute() error {
	// resoling the target configuration
	s.cfg = s.API.SignatureService.GetConfiguration()

	// create basic unpopulated export options object
	res, err := s.createExportOptions()
	if err != nil {
		return err
	}

	// resolving the method for the target
	if res.Method, err = s.resolveMethod(); err != nil {
		return NewSignatureError(err, ErrorExportOptions)
	}

	// enabled by default for AppStore signing method
	if res.Method == "app-store" {
		res.UploadBitCode = true
		res.UploadSymbols = true
	} else {
		// TODO: configurable via CLI
	}

	// encoding the struct to plist format
	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	if err := encoder.Encode(res); err != nil {
		return NewSignatureError(err, ErrorExportOptions)
	}

	// and exporting it to the destination file
	if err := s.exportToFile(&buf); err != nil {
		return err
	}

	return nil
}

func (s exportOptionsService) exportToFile(buf *bytes.Buffer) error {
	// create
	dir := filepath.Dir(s.API.PathService.ExportPList())
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModeSticky|os.ModePerm); err != nil {
			return NewSignatureError(err, ErrorExportOptions)
		}
	}

	// creating the output file
	f, err := os.Create(s.API.PathService.ExportPList())
	if err != nil {
		return NewSignatureError(err, ErrorExportOptions)
	}
	defer f.Close()

	// Writing the buffer to the output file
	if _, err := buf.WriteTo(f); err != nil {
		return NewSignatureError(err, ErrorExportOptions)
	}

	return nil
}

func (s exportOptionsService) resolveTarget() (*api.TargetSignatureConfig, error) {
	for _, e := range *s.cfg {
		if e.TargetName == s.API.Config.Target {
			s.tgt = e
			return &e, nil
		}
	}

	return nil, errors.New("not found")
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
