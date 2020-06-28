package signature

import (
	"context"
	"dothething/internal/config"
	"dothething/internal/xcode/project"
	"fmt"
)

type Service interface {
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

	//
	s.ResolveFor(config.Target, config.Configuration)

	return nil
}

func (s service) ResolveFor(t string, config string) ([]string, error) {
	var res []string
	// Resolving t
	nativeTarget, err := s.Project.Pbx.FindTargetByName(t)
	if err != nil {
		return res, err
	}

	// Resolving build configuration
	bc, err := nativeTarget.BuildConfigurationList.FindConfiguration(config)
	if err != nil {
		return res, err
	}

	fmt.Println(bc, err)

	return res, nil
}
