package output

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
)

const (
	PASS       = "."
	FAIL       = "x"
	PENDING    = "P"
	COMPLETION = ">"
	MEASURE    = "@"
)

var (
	ColorSuccess    = color.New(color.FgCyan, color.Bold)
	ColorFail       = color.New(color.FgHiRed, color.Bold)
	ColorWarn       = color.New(color.FgYellow, color.Bold)
	ColorCompletion = color.New(color.FgGreen, color.Bold)
	ColorMeasure    = color.New(color.FgGreen, color.Bold)
	ColorPhase      = color.New(color.FgCyan, color.Bold)
)

type SimpleReporter struct {
}

func (s SimpleReporter) BuildAggregate(e LogEntry) {
}

func (s SimpleReporter) BuildTimeSummary(e LogEntry) {
	logMeasure("TIMING",
		color.New(color.Bold).Sprintf("%-20s", e.Name),
		fmt.Sprintf("%2s task(s) in %v", e.Count, fmt.Sprintf("%8s %v", e.Time, e.Unit)))
}

func (s SimpleReporter) ErrorCompile(e LogEntry)         { logError("Error compiling", e.FileName, e.Error) }
func (s SimpleReporter) ErrorCodeSign(e LogEntry)        { logError("Error codesign", e.Error, "") }
func (s SimpleReporter) ErrorClang(e LogEntry)           { logError("Error clang", e.Error, "") }
func (s SimpleReporter) ErrorFatal(e LogEntry)           { logError("Error fatal", e.Error, "") }
func (s SimpleReporter) ErrorSignature(e LogEntry)       { logError("Error signature", e.Error, "") }
func (s SimpleReporter) ErrorMissing(e LogEntry)         { logError("Error missing", e.Error, e.FileName) }
func (s SimpleReporter) ErrorLD(e LogEntry)              { logError("Error LD", e.Error, "") }
func (s SimpleReporter) ErrorUndefinedSymbol(e LogEntry) { logError("Error symbol", e.BuildArch, "") }
func (s SimpleReporter) FormatAnalyze(e LogEntry)        { logPhase("ANALYZE", e.FileName, "") }
func (s SimpleReporter) FormatAnalyzeTarget(e LogEntry)  { logPhase("ANALYZE", e.Target, "") }
func (s SimpleReporter) FormatAggregateTarget(e LogEntry) {
	log.Debug().
		Str("Configuration", e.Configuration).
		Str("Target", e.Target).
		Str("Project", e.Project).
		Msg("Aggregate")
}

func (s SimpleReporter) FormatBuildTarget(e LogEntry)       { logPhase("BUILD", e.Target, e.Configuration) }
func (s SimpleReporter) FormatCheckDependencies(e LogEntry) {}
func (s SimpleReporter) CleanRemove(e LogEntry)             { logPhase("CLEANING", "", "") }
func (s SimpleReporter) CleanTarget(e LogEntry)             { logPhase("CLEANING", e.Target, e.Project) }
func (s SimpleReporter) CodeSign(e LogEntry)                { logPhase("CODESIGN", e.FileName, "") }
func (s SimpleReporter) CodeSignTarget(e LogEntry) {
	logPhase("CODESIGN", e.FileName, fmt.Sprintf("%v %v (%v)", e.SigningIdentity, e.ProvisioningName, e.ProvisioningID))
}
func (s SimpleReporter) CompileClang(e LogEntry)         { logProgress("COMPILING", e.FileName, "") }
func (s SimpleReporter) CompileCommand(e LogEntry)       { logProgress("COMPILING", e.FileName, "") }
func (s SimpleReporter) CompileXIB(e LogEntry)           { logProgress("COMPILING", "xib", e.FileName) }
func (s SimpleReporter) CompileStoryboard(e LogEntry)    { logProgress("COMPILING", e.FileName, "") }
func (s SimpleReporter) CopyHeader(e LogEntry)           {}
func (s SimpleReporter) Copy(e LogEntry)                 {}
func (s SimpleReporter) Linking(e LogEntry)              {}
func (s SimpleReporter) GenerateDSym(e LogEntry)         { logPhase("GENERATE DSYM", e.FileName, "") }
func (s SimpleReporter) LibTool(e LogEntry)              { logPhase("LIBTOOL", e.FileName, "") }
func (s SimpleReporter) PhaseSucceeded(e LogEntry)       { logPhase("PHASE", e.Name, "completed") }
func (s SimpleReporter) PhaseScriptExecution(e LogEntry) { logPhase("SCRIPT", e.Name, "") }
func (s SimpleReporter) RunningShellCommand(e LogEntry)  {}
func (s SimpleReporter) TestCasePassed(e LogEntry)       { logSuccess("TEST PASSED", e.TestCase, e.Time) }
func (s SimpleReporter) TestFailing(e LogEntry)          { logError("TEST FAILED", e.FileName, "") }
func (s SimpleReporter) TestCase(e LogEntry) {
	if e.Status == "passed" {
		logSuccess("TEST PASSED", e.TestCase, e.Time)
	} else if e.Status == "failed" {
		logError("TEST FAILED", e.TestCase, "")
	}
}
func (s SimpleReporter) TestCaseMeasured(e LogEntry) {}
func (s SimpleReporter) TestSuiteStatus(e LogEntry) {
	if e.Status == "failed" {
		logError("SUITE FAILED", e.TestSuite, "")
	} else if e.Status == "started" {
		logMeasure("SUITE STARTED", e.TestSuite, "")
	}
}

func (s SimpleReporter) Touch(e LogEntry)     { logPhase("TOUCH", e.FileName, "") }
func (s SimpleReporter) Warning(e LogEntry)   { logWarning("WARNING", e.Message, e.FileName) }
func (s SimpleReporter) WriteAuxiliaryFiles() {}
func (s SimpleReporter) WriteFiles()          {}

func logMsg(msg string, msg1 string, msg2 string) {
	fmt.Printf("  %v %v %v\n", msg, msg1, msg2)
}

func logMsg2(c *color.Color, msg string, msg1 string, msg2 string) {
	fmt.Printf("  %-30s %v %v\n", c.Sprint(msg), msg1, msg2)
}

func logProgress(msg string, msg1 string, msg2 string) {
	logMsg2(ColorWarn, fmt.Sprintf("%v %v", COMPLETION, msg), msg1, msg2)
}

func logSuccess(msg string, msg1 string, msg2 string) {
	logMsg2(ColorSuccess, fmt.Sprintf("%v %v", COMPLETION, msg), msg1, msg2)
}

func logError(msg string, msg1 string, msg2 string) {
	logMsg2(ColorFail, fmt.Sprintf("%v %v", FAIL, msg), msg1, msg2)
}

func logPhase(msg string, msg1 string, msg2 string) {
	logMsg2(ColorPhase, fmt.Sprintf("%v %v", COMPLETION, msg), msg1, msg2)
}

func logMeasure(msg string, msg1 string, msg2 string) {
	logMsg2(ColorMeasure, fmt.Sprintf("%v %v", MEASURE, msg), msg1, msg2)
}

func logWarning(msg string, msg1 string, msg2 string) {
	logMsg2(ColorWarn, fmt.Sprintf("%v %v", FAIL, msg), msg1, msg2)
}
