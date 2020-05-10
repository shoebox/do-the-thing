package output

type LogEntry struct {
	Entry             matcherEntry
	Arg               string
	BuildArch         string
	BuildVariant      string
	Command           string
	Count             string
	Configuration     string
	FileName          string
	FilePath          string
	Name              string
	Path              string
	Project           string
	SourceFile        string
	Target            string
	TargetFile        string
	TestCase          string
	TestFailureReason string
	TestSuite         string
	Time              string
}
