package output

type LogEntry struct {
	Arg               string
	AverageTime       string
	BuildArch         string
	BuildVariant      string
	Command           string
	Configuration     string
	Count             string
	Entry             matcherEntry
	Error             string
	FileName          string
	FilePath          string
	Message           string
	Name              string
	Path              string
	Project           string
	SigningIdentity   string
	ProvisioningName  string
	ProvisioningID    string
	SourceFile        string
	Status            string
	Target            string
	TargetFile        string
	TestCase          string
	TestFailureReason string
	TestSuite         string
	TimeStamp         string
	Time              string
	Unit              string
}
