package output

func NewMatcher(reporter reporter) []entry {
	return []entry{
		// restartingTestsMatcher   = createEntry(reporter.Copy, `^Restarting after unexpected exit or crash in.+$`)
		// writeFileMatcher         = createEntry(reporter, `^write-file\s(.*)`)
		//writeAuxiliaryFiles      = createEntry(reporter, `^Write auxiliary files`)
		createEntry(reporter.FormatAnalyze, `^Analyze(?:Shallow)?\s(?P<FilePath>.*\/(?P<FileName>.*\.(?:m|mm|cc|cpp|c|cxx)))\s*`),
		createEntry(reporter.CleanRemove, `^Clean.Remove`),
		createEntry(reporter.CleanTarget, `(?i)^=== Clean Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith Configuration\s(?P<Configuration>.*)\s===`),
		createEntry(reporter.CleanTarget, `^\s*Executed`),
		createEntry(reporter.TestPassed, `^\s*Test Case\s'-\[(?P<TestSuite>.*)\s(?P<TestCase>.*)\]'\spassed\s\((?P<Time>\d*\.\d{3})\sseconds\)`),
		createEntry(reporter.CodeSign, `^CodeSign\s(?P<FilePath>(?:\\ |[^ ])*)$`),
		createEntry(reporter.CodeSign, `^CodeSign\s(?P<FilePath>(?:\\ |[^ ])*.framework)\/Versions`),
		createEntry(reporter.CompileCommand, `(?i)^Compile[\w]+\s.+?\s(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift)))\s.*`),
		createEntry(reporter.CompileCommand, `(?i)^=== Build Aggregate Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createEntry(reporter.CompileCommand, `^Compile[\w]+\s.+?\s(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift)))\s.*`),
		createEntry(reporter.CompileCommand, `^\s*(?P<Command>.*clang\s.*\s\-c\s(?P<FilePath>.*\.(?:m|mm|c|cc|cpp|cxx))\s.*\.o)$`),
		createEntry(reporter.CompileStoryboard, `^CompileStoryboard\s(?P<FilePath>.*\/(?P<FileName>[^\/].*\.storyboard))`),
		createEntry(reporter.CompileXIB, `^CompileXIB\s(?P<FilePath>.*\/(?P<FileName>.*\.xib))`),
		createEntry(reporter.Copy, `^CopyPlistFile\s(?P<SourceFile>.*\.plist)\s(?P<TargetFile>.*\.plist)`),
		createEntry(reporter.Copy, `^CopyStringsFile.*\/(?P<FileName>.*.strings)`),
		createEntry(reporter.Copy, `^CpResource\s(?P<FilePath>.*)\s\/`),
		createEntry(reporter.CopyHeader, `(?i)^CpHeader\s(?P<SourceFile>.*\.h)\s(?P<TargetFile>.*\.h)`),
		createEntry(reporter.FormatAnalyze, `^Analyze(:Shallow)?\s(?P<RelativePath>.*\/(?P<FileName>.*\.(?:m|mm|cc|cpp|c|cxx)))\s`),
		createEntry(reporter.FormatAnalyzeTarget, `^=== Analyze Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createEntry(reporter.FormatBuildTarget, `(?i)^=== Build Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createEntry(reporter.FormatCheckDependencies, `^Check dependencies`),
		createEntry(reporter.GenerateDSym, `(?i)^GenerateDsymfile \/.*\/(?P<FileName>.*\.dSym)`),
		createEntry(reporter.LibTool, `^Libtool.*\/(?P<FileName>.*\.a)`),
		createEntry(reporter.Linking, `^Ld \/?.*\/(?P<Target>.*?) (?P<BuildVariant>.*) (?P<Arch>.*)$`),
		createEntry(reporter.RunningShellCommand, `^\s{4}(?P<Command>cd|setenv|(?:[\w\/:\\\s\-.]+?\/)?[\w\-]+)\s(?P<Arg>.*)$`),
		createEntry(reporter.TestFailing, `^\s*(?P<File>.+:\d+):\serror:\s[\+\-]\[(?P<TestSuite>.*)\s(?P<TestCase>.*)\]\s:(?:\s'.*'\s\[Failed\],)?\s(?P<Reason>.*)`),
		createEntry(reporter.TestSuiteStart, `^\s*Test Suite '(:.*\/)?(?P<FileName>.*[ox]ctest.*)' started at(?P<TimeStamp>.*)`),
		createEntry(reporter.TestSuiteStart, `^\s*Test Suite '(?P<TestSuite>.*)' started at`),
		createEntry(reporter.TillUtif, `^TiffUtil\s(?P<FileName>.*)`),
		createEntry(reporter.Touch, `^Touch\s(?P<FilePath>.*\/(?P<FileName>.+))`),

		createEntry(reporter.BuildTimeSummary, `^\s*(?P<Name>.*)\s(?:\((?P<Count>[0-9]+) task?s)\)\s\|\s(?P<Time>\d*\.\d{3})\ssecond?s`),
	}
}
