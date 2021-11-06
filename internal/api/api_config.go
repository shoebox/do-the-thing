package api

type Config struct {
	Scheme         string
	Configuration  string
	Destination    Destination
	Path           string
	CodeSign       bool
	CodeSignOption SignConfig
	Target         string
	XCodeVersion   string
}

type SignConfig struct {
	Path                string
	CertificatePassword string
	XCConfig            string
}
