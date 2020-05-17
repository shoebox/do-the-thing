package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

const (
	PASS       = "✓"
	FAIL       = "✗"
	PENDING    = "⧖"
	COMPLETION = "▸"
	MEASURE    = "◷"

	ASCII_PASS       = "."
	ASCII_FAIL       = "x"
	ASCII_PENDING    = "P"
	ASCII_COMPLETION = ">"
	ASCII_MEASURE    = 'T'
)

type SimpleReporter struct {
}

func (s SimpleReporter) BuildAggregate(e LogEntry) {
}

func (s SimpleReporter) BuildTimeSummary(e LogEntry) {
	logMeasure("TIMING",
		e.Name,
		fmt.Sprintf("%v task(s) in %v", e.Count, fmt.Sprintf("%v %v", e.Time, e.Unit)))
}

func (s SimpleReporter) ErrorCompile(e LogEntry) {
	logError("Compiling", e.FileName, e.Error)
}

func (s SimpleReporter) ErrorCodeSign(e LogEntry) {
	fmt.Printf("%v - Signing %v\n", color.RedString("✗ ERROR"), e.Error)
}

func (s SimpleReporter) ErrorClang(e LogEntry) {
}

func (s SimpleReporter) ErrorFatal(e LogEntry) {
	fmt.Printf("%v - %v", color.RedString("ERROR"), e.Error)
}

func (s SimpleReporter) ErrorSignature(e LogEntry) {
}

func (s SimpleReporter) ErrorMissing(e LogEntry) {
}

func (s SimpleReporter) ErrorLD(e LogEntry) {
}

func (s SimpleReporter) ErrorUndefinedSymbol(e LogEntry) {
}

func (s SimpleReporter) FormatAnalyze(e LogEntry) {
	log.Debug().
		Str("Path", e.Path).
		Str("FileName", e.FileName).
		Msg("Analyze")
}

func (s SimpleReporter) FormatAnalyzeTarget(e LogEntry) {
	log.Debug().
		Str("Configuration", e.Configuration).
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Analyze")
}

func (s SimpleReporter) FormatAggregateTarget(e LogEntry) {
	log.Debug().
		Str("Configuration", e.Configuration).
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Aggregate")
}

func (s SimpleReporter) FormatBuildTarget(e LogEntry) {
	logPhase("BUILD TARGET", e.Target, e.Configuration)
}

func (s SimpleReporter) FormatCheckDependencies(e LogEntry) {
	log.Debug().Msg("Checking dependencies")
}

func (s SimpleReporter) CleanRemove(e LogEntry) {
	log.Debug().Msg("Cleaning")
}

func (s SimpleReporter) CleanTarget(e LogEntry) {
	fmt.Printf("♻️ %v - %v\n", color.YellowString("Cleaning target"), e.Target)
	log.Debug().
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Cleaning target")
}

func (s SimpleReporter) CodeSign(e LogEntry) {
	logPhase("CODESIGN", e.FileName, "")
}

func (s SimpleReporter) CompileClang(e LogEntry) {
	logPhase("COMPILING", e.FileName, "")
}

func (s SimpleReporter) CompileCommand(e LogEntry) {
	logPhase("COMPILING", e.FileName, "")
}

func (s SimpleReporter) CompileStoryboard(e LogEntry) {
	logPhase("COMPILING", "Storyboard", e.FileName)
}

func (s SimpleReporter) CompileXIB(e LogEntry) {
	logPhase("COMPILING", "xib", e.FileName)
}

func (s SimpleReporter) CopyHeader(e LogEntry) {
}

func (s SimpleReporter) Copy(e LogEntry) {
}

func (s SimpleReporter) GenerateDSym(e LogEntry) {
	log.Debug().
		Str("File", e.FileName).
		Msg("Generating DSYM")
}

func (s SimpleReporter) LibTool(e LogEntry) {
	logPhase("LIBTOOL", e.FileName, "")
}

func (s SimpleReporter) Linking(e LogEntry) {
	log.Debug().
		Str("Architecture", e.BuildArch).
		Str("Build variant", e.BuildVariant).
		Str("Target", e.Target).
		Msg("Linking")
}

func (s SimpleReporter) PhaseSucceeded(e LogEntry) {
	color.New(color.FgHiBlue).Add(color.Bold).Printf("PHASE '%v' COMPLETED\n", e.Name)
}

func (s SimpleReporter) PhaseScriptExecution(e LogEntry) {
	logPhase("SCRIPT", e.Name, "")
}

func (s SimpleReporter) RunningShellCommand(e LogEntry) {
}

func (s SimpleReporter) TestCasePassed(e LogEntry) {
	logSuccess("TEST PASSED", e.TestCase, e.Time)
}

func (s SimpleReporter) TestCase(e LogEntry) {
	if e.Status == "passed" {
		logSuccess("TEST PASSED", e.TestCase, e.Time)
	} else if e.Status == "failed" {
		logError("TEST FAILED", e.TestCase, "")
	}
}

func (s SimpleReporter) TestCaseMeasured(e LogEntry) {
}

func (s SimpleReporter) TestFailing(e LogEntry) {
	logError("TEST FAILED", e.FileName, "")
}

func (s SimpleReporter) TestSuiteStatus(e LogEntry) {
	if e.Status == "failed" {
		logError("SUITE FAILED", e.TestSuite, "")
	} else {
		//.Printf(" Test suite %v - %v\n", e.Status, e.TestSuite)
	}
}

func (s SimpleReporter) Touch(e LogEntry) {
	logPhase("TOUCH", e.FileName, "")
}

func (s SimpleReporter) Warning(e LogEntry) {
	logWarning("WARNING", e.Message, e.FileName)
}

func (s SimpleReporter) WriteAuxiliaryFiles() {
	panic("implement me")
}

func (s SimpleReporter) WriteFiles() {
	panic("implement me")
}

func logSuccess(msg string, msg1 string, msg2 string) {
	fmt.Printf("  %v %v %v\n", color.GreenString("%v %-10s", PASS, msg), msg1, msg2)
}

func logError(msg string, msg1 string, msg2 string) {
	fmt.Printf("  %v %v %v\n", color.RedString("%v %-10s", FAIL, msg), msg1, msg2)
}

func logPhase(msg string, msg1 string, msg2 string) {
	fmt.Printf("  %v %v %v\n", color.YellowString("%v %-10s", COMPLETION, msg), msg1, msg2)
}

func logMeasure(msg string, msg1 string, msg2 string) {
	fmt.Printf("  %v %-30s %v\n", color.GreenString("%v %-10s", MEASURE, msg), msg1, msg2)
}

func logWarning(msg string, msg1 string, msg2 string) {
	fmt.Printf("  %v %-30s %v\n", color.YellowString("%v %-10s", FAIL, msg), msg1, msg2)
}
