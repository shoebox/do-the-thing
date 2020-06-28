package signature

import (
	"context"
	"dothething/internal/config"
	"dothething/internal/xcode/pbx"
	"dothething/internal/xcode/project"
	"fmt"
)

type Service interface {
	Run(ctx context.Context, config config.Config) error
}

type service struct {
	project.Project
	project.ProjectService
	Resolver
}

func New(ps project.ProjectService, sr Resolver) Service {
	return service{ProjectService: ps, Resolver: sr}
}

func (s service) Run(ctx context.Context, config config.Config) error {
	var err error

	// Parsing project first
	s.Project, err = s.ProjectService.Parse(ctx)
	if err != nil {
		return err
	}

	// Validate the configuration against the project
	if err := s.Project.ValidateConfiguration(config); err != nil {
		return err
	}

	// Resolve build configuration to target
	bc, err := s.ResolveFor(config.Target, config.Configuration)
	if err != nil {
		return err
	}

	// Resolve configuration for bundle identifier
	sc, err := s.Resolver.Resolve(ctx, config.CodeSignOption.Path, bc.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"])
	if err != nil {
		return err
	}
	fmt.Println(sc)

	return nil
}

func (s service) ResolveFor(t string, config string) (pbx.XCBuildConfiguration, error) {
	var res pbx.XCBuildConfiguration
	// Resolving t
	nativeTarget, err := s.Project.Pbx.FindTargetByName(t)
	if err != nil {
		return res, err
	}

	// Resolving build configuration
	return nativeTarget.BuildConfigurationList.FindConfiguration(config)
}
