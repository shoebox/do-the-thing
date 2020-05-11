package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

type simplereporter struct {
}

func (s simplereporter) BuildTimeSummary(e LogEntry) {
	log.Debug().
		Str("Name", e.Name).
		Str("Task(s) count", e.Count).
		Msg("BuildTimeSummary")
}

func (s simplereporter) ErrorCompile(e LogEntry) {
	fmt.Printf("%v - Compiling %v %v\n", color.RedString("✗ ERROR"), e.FileName, e.Error)
}

func (s simplereporter) ErrorCodeSign(e LogEntry) {
	fmt.Printf("%v - Signing %v\n", color.RedString("✗ ERROR"), e.Error)
}

func (s simplereporter) ErrorClang(e LogEntry) {
}

func (s simplereporter) ErrorFatal(e LogEntry) {
	fmt.Printf("%v - %v", color.RedString("ERROR"), e.Error)
}

func (s simplereporter) ErrorSignature(e LogEntry) {
}

func (s simplereporter) ErrorMissing(e LogEntry) {
}

func (s simplereporter) ErrorLD(e LogEntry) {
}

func (s simplereporter) ErrorUndefinedSymbol(e LogEntry) {
}

func (s simplereporter) FormatAnalyze(e LogEntry) {
	log.Debug().
		Str("Path", e.Path).
		Str("FileName", e.FileName).
		Msg("Analyze")
}

func (s simplereporter) FormatAnalyzeTarget(e LogEntry) {
	log.Debug().
		Str("Configuration", e.Configuration).
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Analyze")
}

func (s simplereporter) FormatAggregateTarget(e LogEntry) {
	log.Debug().
		Str("Configuration", e.Configuration).
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Aggregate")
}

func (s simplereporter) FormatBuildTarget(e LogEntry) {
	log.Debug().
		Str("Configuration", e.Configuration).
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Building")
}

func (s simplereporter) FormatCheckDependencies(e LogEntry) {
	log.Debug().Msg("Checking dependencies")
}

func (s simplereporter) CleanRemove(e LogEntry) {
	log.Debug().Msg("Cleaning")
}

func (s simplereporter) CleanTarget(e LogEntry) {
	fmt.Printf("♻️ %v - %v\n", color.YellowString("Cleaning target"), e.Target)
	log.Debug().
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Cleaning target")
}

func (s simplereporter) CodeSign(e LogEntry) {
	log.Info().
		Str("File path", e.FilePath).
		Msg("Signin")
}

func (s simplereporter) CompileCommand(e LogEntry) {
	fmt.Printf("%v %v %v\n", color.YellowString("▸ COMPILING"), "", e.FileName)
}

func (s simplereporter) CompileStoryboard(e LogEntry) {
	fmt.Printf("%v %v %v\n", color.YellowString("▸ COMPILING"), "storyboard", e.FileName)
	log.Debug().
		Str("File name", e.FileName).
		Str("File path", e.FilePath).
		Msg("Compiling storyboard")
}

func (s simplereporter) CompileXIB(e LogEntry) {
	fmt.Printf("%v %v %v\n", color.YellowString("▸ COMPILING"), "xib", e.FileName)
}

func (s simplereporter) CopyHeader(e LogEntry) {
	log.Debug().
		Str("Source", e.SourceFile).
		Str("Target", e.TargetFile).
		Msg("Copying")
}

func (s simplereporter) Copy(e LogEntry) {
	log.Debug().
		Str("Res", e.Arg).
		Msg("Copying")
}

func (s simplereporter) GenerateDSym(e LogEntry) {
	log.Debug().
		Str("File", e.FileName).
		Msg("Generating DSYM")
}

func (s simplereporter) LibTool(e LogEntry) {
	log.Debug().
		Str("File", e.FileName).
		Msg("Building library")
}

func (s simplereporter) Linking(e LogEntry) {
	log.Debug().
		Str("Architecture", e.BuildArch).
		Str("Build variant", e.BuildVariant).
		Str("Target", e.Target).
		Msg("Linking")
}

func (s simplereporter) PhaseSucceeded(e LogEntry) {
	color.New(color.FgHiBlue).Add(color.Bold).Printf("PHASE '%v' COMPLETED\n", e.Name)
}

func (s simplereporter) PhaseScriptExecution(e LogEntry) {
	fmt.Println("PhaseScriptExecution", e.Name)
}

func (s simplereporter) RunningShellCommand(e LogEntry) {
	log.Debug().
		Str("Command", e.Command).
		Str("Arg", e.Arg).
		Msg("Running shell command")
}

func (s simplereporter) TestCasePassed(e LogEntry) {
	fmt.Printf("%v Test case - %v (%vs)\n",
		color.GreenString("✔️ PASSED"),
		e.TestCase,
		e.Time)
}

func (s simplereporter) TestCasePending(e LogEntry) {
	fmt.Printf("⏳ Test case %v - %v\n", color.GreenString("PENDING"), e.TestCase)
}

func (s simplereporter) TestCaseStarted(e LogEntry) {
}

func (s simplereporter) TestCaseMeasured(e LogEntry) {
}

func (s simplereporter) TestFailing(e LogEntry) {
	fmt.Printf("✖️ %v - %v\n", color.RedString("Test failed"), e.FileName)
}

func (s simplereporter) TestSuiteStarted(e LogEntry) {
	fmt.Printf("✖️ %v - %v\n", color.RedString("Test suite started"), e.TestSuite)
}

func (s simplereporter) TestSuiteComplete(e LogEntry) {
	log.Debug().
		Str("Name", e.TestSuite).
		Msg("Test suite started")
}

func (s simplereporter) TillUtif(e LogEntry) {
	log.Debug().
		Str("fileName", e.FileName).
		Msg("Validating")
}

func (s simplereporter) Touch(e LogEntry) {
	// panic("implement me")
}

func (s simplereporter) WriteAuxiliaryFiles() {
	panic("implement me")
}

func (s simplereporter) WriteFiles() {
	panic("implement me")
}
