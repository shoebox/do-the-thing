package xcode

type Config struct {
	Scheme         string
	Configuration  string
	Path           string
	CodeSign       bool
	CodeSignOption struct {
		Path string
	}
}
