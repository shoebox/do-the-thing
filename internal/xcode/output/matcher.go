package output

func NewMatcher(reporter reporter) []matcherEntry {
	return []matcherEntry{
		createMatcherEntry(reporter.BuildTimeSummary, `^(?P<Name>(\w+))\s\((?P<Count>\d+)\stask(?:s)?\)\s\|\s(?P<Time>[\d.]+)\s(?P<Unit>\w+)$`),
		createMatcherEntry(reporter.CleanTarget, `(?i)^=== Clean Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.CodeSign, `^CodeSign\s(?P<FilePath>(?:\\ |[^ ])*)$`),
		createMatcherEntry(reporter.CodeSign, `^CodeSign\s(?P<FilePath>(?:\\ |[^ ])*.framework)\/Versions`),
		createMatcherEntry(reporter.BuildAggregate, `(?i)^=== Build Aggregate Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.CompileCommand, `(?i)^Compile[\w]+\s.+?\s(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift)))\s.*`),
		createMatcherEntry(reporter.CompileCommand, `^Compile[\w]+\s.+?\s(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift)))\s.*`),
		createMatcherEntry(reporter.CompileClang, `^\s*(?P<Command>.*clang\s.*\s\-c\s(?P<FilePath>.*\.(?:m|mm|c|cc|cpp|cxx))\s.*\.o)$`),
		createMatcherEntry(reporter.CompileStoryboard, `^CompileStoryboard\s(?P<FilePath>.*\/(?P<FileName>[^\/].*\.storyboard))`),
		createMatcherEntry(reporter.CompileXIB, `^CompileXIB\s(?P<FilePath>.*\/(?P<FileName>.*\.xib))`),
		createMatcherEntry(reporter.Copy, `^CopyPlistFile\s(?P<SourceFile>.*\.plist)\s(?P<TargetFile>.*\.plist)`),
		createMatcherEntry(reporter.Copy, `^CopyStringsFile.*\/(?P<FileName>.*.strings)`),
		createMatcherEntry(reporter.Copy, `^CpResource\s(?P<FilePath>.*)\s\/`),
		createMatcherEntry(reporter.CopyHeader, `(?i)^CpHeader\s(?P<SourceFile>.*\.h)\s(?P<TargetFile>.*\.h)`),
		createMatcherEntry(reporter.FormatAnalyze, `^Analyze(:Shallow)?\s(?P<FilePath>.*\/(?P<FileName>.*\.(?:m|mm|cc|cpp|c|cxx)))`),
		createMatcherEntry(reporter.FormatAnalyzeTarget, `(?i)^=== Analyze Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.FormatBuildTarget, `(?i)^=== Build Target\s(?P<Target>.*)\sOf Project\s(?P<Project>.*)\sWith.*Configuration\s(?P<Configuration>.*)\s===`),
		createMatcherEntry(reporter.FormatCheckDependencies, `^Check dependencies`),
		createMatcherEntry(reporter.GenerateDSym, `(?i)^GenerateDsymfile \/.*\/(?P<FileName>.*\.dSym)`),
		createMatcherEntry(reporter.LibTool, `^Libtool.*\/(?P<FileName>.*\.a)`),
		createMatcherEntry(reporter.Linking, `^Ld \/?.*\/(?P<Target>.*?) (?P<BuildVariant>.*) (?P<Arch>.*)$`),
		createMatcherEntry(reporter.PhaseScriptExecution, `(?i)^PhaseScriptExecution\s(?P<Name>(?:\\.|[^ ])+)\s(?P<FilePath>(?:\\.|[^ ])+)`),
		createMatcherEntry(reporter.PhaseSucceeded, `^\*\*\s(?P<Name>.*)\sSUCCEEDED\s\*\*`),
		createMatcherEntry(reporter.RunningShellCommand, `^\s{4}(?P<Command>cd|setenv|(?:[\w\/:\\\s\-.]+?\/)?[\w\-]+)\s(?P<Arg>.*)$`),
		createMatcherEntry(reporter.TestCase, `^Test Case\s'-\[(?P<TestSuite>.*)\s(?P<TestCase>.*)\]'\s(?P<Status>started|failed|passed)(?:\.|\s\(?(?P<Time>\d*\.\d{3}) (?P<Unit>\w+)\)\.)$`),
		createMatcherEntry(reporter.TestCaseMeasured, `^(?P<FilePath>(?:\\.|[^ ])+\/(?P<FileName>(?:\\.|[^ ])+\.(?:m|mm|c|cc|cpp|cxx|swift))):\w+:\sTest Case\s'-\[(?P<TestSuite>.*)\s(?P<TestCase>.*)\]'\smeasured\s(?:\[Time,\s(?P<Unit>\w+)\]\saverage:\s(?P<AverageTime>\d*\.\d{3}),)`),
		createMatcherEntry(reporter.TestSuiteStatus, `^Test Suite\s'(?P<TestSuite>\w+?.\w+)'\s(?P<Status>started|finished|failed)\sat\s(?P<TimeStamp>.*)?`),
		createMatcherEntry(reporter.Touch, `^Touch\s((?P<FilePath>.*\/(?P<FileName>.+))\s\(.+)$`),

		// Errors matchers
		// TODO: Undefined and duplicated symbols
		createMatcherEntry(reporter.ErrorClang, `^(clang: error:(?P<Error>.*))$`),
		createMatcherEntry(reporter.ErrorCodeSign, `^(Code\s?Sign error:(?P<Error>.*|Code signing is required for product type .* in SDK .*)|No profile matching .* found:.*|Provisioning profile .* doesn't .*|Swift is unavailable on .*|.?Use Legacy Swift Language Version.*)$`),
		createMatcherEntry(reporter.ErrorCompile, `^(\/.+\/(?P<FileName>.*):.*:.*):\s(?:fatal\s)?error:\s(?P<Error>.*)$`),
		createMatcherEntry(reporter.ErrorFatal, `^(?:(fatal\s)?error:(?P<Error>.*))$`),
		createMatcherEntry(reporter.ErrorLD, `^(?:ld:(?P<Error>.*))`),
		createMatcherEntry(reporter.ErrorMissing, `^<unknown>:0:\s(?:error:\s(?P<Error>.*))\s'(?P<FilePath>\/.+\/(?P<FileName>.*\..*))'$`),
		createMatcherEntry(reporter.ErrorSignature, `^(?P<Error>.*requires a provisioning profile.*|No certificate matching.*)$`),
		createMatcherEntry(reporter.ErrorSignature, `^(?P<Error>.*requires a provisioning profile.*|No certificate matching.*)$`),
		createMatcherEntry(reporter.ErrorUndefinedSymbol, `^Undefined symbols for architecture (?P<Arch>.*):$`),

		// Warning matchers
		createMatcherEntry(reporter.Warning, `^(?P<FilePath>\/.+\/(?P<FileName>.*):.*:.*):\swarning:\s(?P<Message>.*)$`),
		createMatcherEntry(reporter.Warning, `^(ld: )warning: (?P<Message>.*)`),
		createMatcherEntry(reporter.Warning, `^warning:\s(?P<Message>.*)$`),
	}
}
