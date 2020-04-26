package output

import (
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
	log.Info().
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
	log.Info().
		Str("File", e.FileName).
		Msg("Compiling")
}

func (s simplereporter) CompileStoryboard(e LogEntry) {
	log.Info().
		Str("File name", e.FileName).
		Str("File path", e.FilePath).
		Msg("Compiling storyboard")
}

func (s simplereporter) CompileXIB(e LogEntry) {
	log.Info().
		Str("File name", e.FileName).
		Str("File path", e.FilePath).
		Msg("Compiling XIB")
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

func (s simplereporter) RunningShellCommand(e LogEntry) {
	log.Debug().
		Str("Command", e.Command).
		Str("Arg", e.Arg).
		Msg("Running shell command")
}

func (s simplereporter) TestPassed(e LogEntry) {
	log.Info().Str("Name", e.TestCase).Msg("✅ Test Pass")
}

func (s simplereporter) TestFailing(e LogEntry) {
	log.Warn().Str("Name", e.TestCase).Msg("❌ Test failed")
	log.Debug().
		Str("FilePath", e.FilePath).
		Str("Test suite", e.TestSuite).
		Str("Test case", e.TestCase).
		Str("Reason", e.TestFailureReason).
		Msg("Failing test")
}

func (s simplereporter) TestSuiteStart(e LogEntry) {
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
