package signature

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
	"fmt"
)

type service struct {
	api.API
}

func New(api api.API) api.SignatureService {
	return service{api}
}

func (s service) Run(ctx context.Context, target string, bcName string, path string, project api.Project) error {
	var err error

	// Resolve build configuration to target
	bc, err := s.ResolveFor(project, target, bcName)
	if err != nil {
		return err
	}

	// Resolve configuration for bundle identifier
	sc, err := s.API.
		SignatureResolver().
		Resolve(ctx, path, bc.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"])
	if err != nil {
		return err
	}
	fmt.Println(sc)

	return nil
}

func (s service) ResolveFor(pj api.Project, t string, config string) (pbx.XCBuildConfiguration, error) {
	var res pbx.XCBuildConfiguration
	// Resolving target by name
	nativeTarget, err := pj.Pbx.FindTargetByName(t)
	if err != nil {
		return res, err
	}

	// Resolving build configuration
	return nativeTarget.BuildConfigurationList.FindConfiguration(config)
}
