package action

import (
	"dothething/internal/api"
	"dothething/internal/xcode/output"
)

func RunCmd(cmd api.Cmd) error {
	pout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	perr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		f := output.NewFormatter(output.SimpleReporter{})
		f.Parse(pout)
		f.Parse(perr)
	}()

	//
	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
