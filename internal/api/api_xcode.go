package api

import "context"

// ListService basic interface
type ListService interface {
	List(ctx context.Context) ([]*Install, error)
}

// Install xcode installation definition
type Install struct {
	DevPath       string
	Path          string
	BundleVersion string
	Version       string
}

// SelectService The XCode version selection service interface
type SelectService interface {
	Find(ctx context.Context) (*Install, error)
}

// XCodeBuildService service definition
type BuildService interface {
	List(ctx context.Context) (string, error)
	ShowDestinations(ctx context.Context, scheme string) (string, error)
	GetArg() string
}
