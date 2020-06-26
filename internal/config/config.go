package config

type Config struct {
	Scheme         string
	Configuration  string
	Path           string
	CodeSign       bool
	CodeSignOption SignConfig
	Target         string
}

type SignConfig struct {
	Path                string
	CertificatePassword string
	XCConfig            string
}
