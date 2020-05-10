package output

func NewMatcher(reporter reporter) []matcherEntry {
	return []matcherEntry{
		// restartingTestsMatcher   = createEntry(reporter.Copy, `^Restarting after unexpected exit or crash in.+$`)
		// writeFileMatcher         = createEntry(reporter, `^write-file\s(.*)`)
		//writeAuxiliaryFiles      = createEntry(reporter, `^Write auxiliary files`)
		createMatcherEntry(reporter.FormatAnalyze, `^Analyze(?:Shallow)?\s(?P<FilePath>.*\/(?P<FileName>.*\.(?:m|mm|cc|cpp|c|cxx)))\s*`),
		createMatcherEntry(reporter.CleanRemove, `^Clean.Remove`),
		createMatcherEntry(reporter.CleanTarget, `(?i)^=== Clean Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.CleanTarget, `^\s*Executed`),
		createMatcherEntry(reporter.CodeSign, `^CodeSign\s(?P<FilePath>(?:\\ |[^ ])*)$`),
		createMatcherEntry(reporter.CodeSign, `^CodeSign\s(?P<FilePath>(?:\\ |[^ ])*.framework)\/Versions`),
		createMatcherEntry(reporter.CompileCommand, `(?i)^Compile[\w]+\s.+?\s(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift)))\s.*`),
		createMatcherEntry(reporter.CompileCommand, `(?i)^=== Build Aggregate Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.CompileCommand, `^Compile[\w]+\s.+?\s(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift)))\s.*`),
		createMatcherEntry(reporter.CompileCommand, `^\s*(?P<Command>.*clang\s.*\s\-c\s(?P<FilePath>.*\.(?:m|mm|c|cc|cpp|cxx))\s.*\.o)$`),
		createMatcherEntry(reporter.CompileStoryboard, `^CompileStoryboard\s(?P<FilePath>.*\/(?P<FileName>[^\/].*\.storyboard))`),
		createMatcherEntry(reporter.CompileXIB, `^CompileXIB\s(?P<FilePath>.*\/(?P<FileName>.*\.xib))`),
		createMatcherEntry(reporter.Copy, `^CopyPlistFile\s(?P<SourceFile>.*\.plist)\s(?P<TargetFile>.*\.plist)`),
		createMatcherEntry(reporter.Copy, `^CopyStringsFile.*\/(?P<FileName>.*.strings)`),
		createMatcherEntry(reporter.Copy, `^CpResource\s(?P<FilePath>.*)\s\/`),
		createMatcherEntry(reporter.CopyHeader, `(?i)^CpHeader\s(?P<SourceFile>.*\.h)\s(?P<TargetFile>.*\.h)`),
		createMatcherEntry(reporter.FormatAnalyze, `^Analyze(:Shallow)?\s(?P<RelativePath>.*\/(?P<FileName>.*\.(?:m|mm|cc|cpp|c|cxx)))\s`),
		createMatcherEntry(reporter.FormatAnalyzeTarget, `^=== Analyze Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.FormatBuildTarget, `(?i)^=== Build Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.FormatCheckDependencies, `^Check dependencies`),
		createMatcherEntry(reporter.GenerateDSym, `(?i)^GenerateDsymfile \/.*\/(?P<FileName>.*\.dSym)`),
		createMatcherEntry(reporter.LibTool, `^Libtool.*\/(?P<FileName>.*\.a)`),
		createMatcherEntry(reporter.Linking, `^Ld \/?.*\/(?P<Target>.*?) (?P<BuildVariant>.*) (?P<Arch>.*)$`),
		createMatcherEntry(reporter.PhaseSucceeded, `^\*\*\s(?P<Name>.*)\sSUCCEEDED\s\*\*`),
		createMatcherEntry(reporter.PhaseScriptExecution, `^PhaseScriptExecution\s(?P<Name>(?:\\.|[^ ])+)\s(?P<FilePath>(?:\\.|[^ ])+)`),
		createMatcherEntry(reporter.RunningShellCommand, `^\s{4}(?P<Command>cd|setenv|(?:[\w\/:\\\s\-.]+?\/)?[\w\-]+)\s(?P<Arg>.*)$`),
		createMatcherEntry(reporter.TestFailing, `^\s*(?P<File>.+:\d+):\serror:\s[\+\-]\[(?P<TestSuite>.*)\s(?P<TestCase>.*)\]\s:(?:\s'.*'\s\[Failed\],)?\s(?P<Reason>.*)`),
		createMatcherEntry(reporter.TestCasePassed, `^Test Case\s'-\[(?P<TestSuite>.*)\s(?P<TestCase>.*)\]'\spassed\s\((?P<Time>\d*\.\d{3})\sseconds\)`),
		createMatcherEntry(reporter.TestCaseStarted, `^Test Case '-\[(?P<TestSuite>.*) (?P<TestCase>.*)\]' started.$`),
		createMatcherEntry(reporter.TestCasePending, `(?i)^Test Case\s'-\[(.*)\s(.*)PENDING\]'\spassed`),
		createMatcherEntry(reporter.TestCaseMeasured, `^[^:]*:[^:]*:\sTest Case\s'-\[(?P<TestSuiite>.*)\s(?P<TestCase>.*)\]'\smeasured\s\[Time,\sseconds\]\saverage:\s(?P<Time>\d*\.\d{3}),`),
		createMatcherEntry(reporter.TestSuiteComplete, `^\s*Test Suite '(?:.*\/)?(?P<Name>.*[ox]ctest.*)' (?P<Status>finished|passed|failed) at (?P<TimeStamp>.*)`),
		createMatcherEntry(reporter.TillUtif, `^TiffUtil\s(?P<FileName>.*)`),
		createMatcherEntry(reporter.Touch, `^Touch\s(?P<FilePath>.*\/(?P<FileName>.+))`),

		createMatcherEntry(reporter.BuildTimeSummary, `^\s*(?P<Name>.*)\s(?:\((?P<Count>[0-9]+) task?s)\)\s\|\s(?P<Time>\d*\.\d{3})\ssecond?s`),
	}
}
