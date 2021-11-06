package cmd

import (
	"context"
	"dothething/internal/api"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type CLI interface {
	Run() error
}

type menu struct {
	*api.API
}

func New(a *api.API) CLI {
	return menu{API: a}
}

func (m menu) Run() error {
	app := cli.App{
		Name:    "do-the-thing",
		Version: "0.1.0",
	}

	app.Commands = []*cli.Command{
		{Name: "build", Action: m.buildCommand},
		{Name: "package", Action: m.packageCommand},
		{Name: "archive", Action: m.archiveCommand},
		{Name: "test", Action: m.testCommand},
	}

	app.Flags = []cli.Flag{
		&cli.PathFlag{Name: "project", Destination: &m.API.Config.Path},
		&cli.StringFlag{Name: "xcodeVersion", Destination: &m.API.Config.XCodeVersion, EnvVars: []string{"XCODE_VERSION"}},
		&cli.StringFlag{Name: "buildScheme", Destination: &m.API.Config.Scheme},
		&cli.StringFlag{Name: "buildConfiguration", Destination: &m.API.Config.Configuration},
		&cli.StringFlag{Name: "target", Destination: &m.API.Config.Target},
		&cli.StringFlag{Name: "signatureFilesPath", Destination: &m.API.Config.CodeSignOption.Path},
		&cli.StringFlag{Name: "certificatePassword", Destination: &m.API.Config.CodeSignOption.CertificatePassword},
	}

	err := app.Run(os.Args)
	if err != nil {
		return fmt.Errorf("run error  %v", err)
	}

	return nil
}

func (m menu) archiveCommand(c *cli.Context) error {
	return m.runAction(m.API.ActionArchive)
}

func (m menu) buildCommand(c *cli.Context) error {
	return m.runAction(m.API.ActionBuild)
}

func (m menu) packageCommand(c *cli.Context) error {
	return m.runAction(m.API.ActionPack)
}

func (m menu) testCommand(c *cli.Context) error {
	return m.runAction(m.API.ActionRunTest)
}

func (m menu) runAction(action api.Action) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	return action.Run(ctx)
}
