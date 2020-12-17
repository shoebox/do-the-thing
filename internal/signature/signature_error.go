package signature

import "fmt"

const (
	ErrorBuildConfigurationResolution  = "Failed to resolve build configuration"
	ErrorCertificateImport             = "Failed to import certificate"
	ErrorCertificateResolution         = "Failed to resolve matching certificate"
	ErrorProvisioningInstall           = "Failed to install provisioining profile"
	ErrorProvisioningProfileResolution = "Failed to resolve matching provisioning profile"
	ErrorTargetResolution              = "Failed to resolve target"
)

type SignatureError struct {
	Msg string
	error
}

func NewSignatureError(err error, msg string) *SignatureError {
	return &SignatureError{error: err, Msg: msg}
}

func (e SignatureError) Error() string {
	if e.error != nil {
		return fmt.Sprintf("%v (Error : %v)", e.Msg, e.error)
	}

	return e.Msg
}
