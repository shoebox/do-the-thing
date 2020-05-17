package output

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockListener struct {
	mock.Mock
}

func (m *MockListener) BuildAggregate(e LogEntry)          { m.Called(e) }
func (m *MockListener) BuildTimeSummary(e LogEntry)        { m.Called(e) }
func (m *MockListener) CleanRemove(e LogEntry)             { m.Called(e) }
func (m *MockListener) CleanTarget(e LogEntry)             { m.Called(e) }
func (m *MockListener) CodeSign(e LogEntry)                { m.Called(e) }
func (m *MockListener) CompileClang(e LogEntry)            { m.Called(e) }
func (m *MockListener) CompileCommand(e LogEntry)          { m.Called(e) }
func (m *MockListener) CompileStoryboard(e LogEntry)       { m.Called(e) }
func (m *MockListener) CompileXIB(e LogEntry)              { m.Called(e) }
func (m *MockListener) Copy(e LogEntry)                    { m.Called(e) }
func (m *MockListener) CopyHeader(e LogEntry)              { m.Called(e) }
func (m *MockListener) ErrorClang(e LogEntry)              { m.Called(e) }
func (m *MockListener) ErrorCodeSign(e LogEntry)           { m.Called(e) }
func (m *MockListener) ErrorCompile(e LogEntry)            { m.Called(e) }
func (m *MockListener) ErrorFatal(e LogEntry)              { m.Called(e) }
func (m *MockListener) ErrorLD(e LogEntry)                 { m.Called(e) }
func (m *MockListener) ErrorMissing(e LogEntry)            { m.Called(e) }
func (m *MockListener) ErrorSignature(e LogEntry)          { m.Called(e) }
func (m *MockListener) ErrorUndefinedSymbol(e LogEntry)    { m.Called(e) }
func (m *MockListener) FormatAggregateTarget(e LogEntry)   { m.Called(e) }
func (m *MockListener) FormatAnalyze(e LogEntry)           { m.Called(e) }
func (m *MockListener) FormatAnalyzeTarget(e LogEntry)     { m.Called(e) }
func (m *MockListener) FormatBuildTarget(e LogEntry)       { m.Called(e) }
func (m *MockListener) FormatCheckDependencies(e LogEntry) { m.Called(e) }
func (m *MockListener) GenerateDSym(e LogEntry)            { m.Called(e) }
func (m *MockListener) LibTool(e LogEntry)                 { m.Called(e) }
func (m *MockListener) Linking(e LogEntry)                 { m.Called(e) }
func (m *MockListener) PhaseScriptExecution(e LogEntry)    { m.Called(e) }
func (m *MockListener) PhaseSucceeded(e LogEntry)          { m.Called(e) }
func (m *MockListener) RunningShellCommand(e LogEntry)     { m.Called(e) }
func (m *MockListener) TestCase(e LogEntry)                { m.Called(e) }
func (m *MockListener) TestCaseMeasured(e LogEntry)        { m.Called(e) }
func (m *MockListener) TestCasePassed(e LogEntry)          { m.Called(e) }
func (m *MockListener) TestCasePending(e LogEntry)         { m.Called(e) }
func (m *MockListener) TestCaseStarted(e LogEntry)         { m.Called(e) }
func (m *MockListener) TestFailing(e LogEntry)             { m.Called(e) }
func (m *MockListener) TestSuiteComplete(e LogEntry)       { m.Called(e) }
func (m *MockListener) TestSuiteFailed(e LogEntry)         { m.Called(e) }
func (m *MockListener) TestSuiteStatus(e LogEntry)         { m.Called(e) }
func (m *MockListener) TillUtif(e LogEntry)                { m.Called(e) }
func (m *MockListener) Touch(e LogEntry)                   { m.Called(e) }
func (m *MockListener) Warning(e LogEntry)                 { m.Called(e) }

func TestParse(t *testing.T) {
	//
	cases := []struct {
		d string
		m string
		e LogEntry
	}{
		{
			m: "BuildTimeSummary",
			d: "CompileSwiftSources (1 task) | 13.774 seconds",
			e: LogEntry{
				Name:  "CompileSwiftSources",
				Count: "1",
				Time:  "13.774",
				Unit:  "seconds",
			},
		},
		{
			m: "BuildTimeSummary",
			d: "CompileSwiftSources (12 tasks) | 13.774 seconds",
			e: LogEntry{
				Name:  "CompileSwiftSources",
				Count: "12",
				Time:  "13.774",
				Unit:  "seconds",
			},
		},
		{
			m: "CleanTarget",
			d: "=== CLEAN TARGET ReactiveCocoa OF PROJECT ReactiveCocoa WITH CONFIGURATION Debug ===",
			e: LogEntry{
				Target:        "ReactiveCocoa",
				Project:       "ReactiveCocoa",
				Configuration: "Debug",
			},
		},
		{
			m: "CodeSign",
			d: "CodeSign build/Release/CocoaChip.app",
			e: LogEntry{
				FilePath: "build/Release/CocoaChip.app",
			},
		},
		{
			m: "CodeSign",
			d: "CodeSign build/Release/CocoaChipCore.framework/Versions/A",
			e: LogEntry{
				FilePath: "build/Release/CocoaChipCore.framework/Versions/A",
			},
		},
		{
			m: "BuildAggregate",
			d: "=== BUILD AGGREGATE TARGET Be Aggro OF PROJECT AggregateExample WITH CONFIGURATION Debug ===",
			e: LogEntry{
				Target:        "Be Aggro",
				Project:       "AggregateExample",
				Configuration: "Debug",
			},
		},
		{
			m: "CompileCommand",
			d: "CompileSwift normal x86_64 /Users/tot/Desktop/tmp/Swiftstraints/Swiftstraints/DimensionAnchor.swift (in target 'Swiftstraints iOS' from project 'Swiftstraints')",
			e: LogEntry{
				FilePath: "/Users/tot/Desktop/tmp/Swiftstraints/Swiftstraints/DimensionAnchor.swift",
				FileName: "DimensionAnchor.swift",
			},
		},
		{
			m: "CompileClang",
			d: "clang VeryLongCommandcall -c /Users/musalj/code/OSS/ObjectiveSugar/Classes/NSNumber+ObjectiveSugar.m -o /path/to//NSNumber+ObjectiveSugar.o",
			e: LogEntry{
				Command:  "clang VeryLongCommandcall -c /Users/musalj/code/OSS/ObjectiveSugar/Classes/NSNumber+ObjectiveSugar.m -o /path/to//NSNumber+ObjectiveSugar.o",
				FilePath: "/Users/musalj/code/OSS/ObjectiveSugar/Classes/NSNumber+ObjectiveSugar.m",
			},
		},
		{
			m: "CompileStoryboard",
			d: "CompileStoryboard sample/Main.storyboard",
			e: LogEntry{
				FilePath: "sample/Main.storyboard",
				FileName: "Main.storyboard",
			},
		},
		{
			m: "CompileXIB",
			d: "CompileXIB CocoaChip/en.lproj/MainMenu.xib",
			e: LogEntry{
				FilePath: "CocoaChip/en.lproj/MainMenu.xib",
				FileName: "MainMenu.xib",
			},
		},
		{
			m: "Copy",
			d: "CopyStringsFile /Users/musalj/Library/Developer/Xcode/DerivedData/ObjectiveSugar-ayzdhqmmwtqgysdpznmovjlupqjy/Build/Products/Debug-iphonesimulator/ObjectiveSugar.app/en.lproj/InfoPlist.strings ObjectiveSugar/en.lproj/InfoPlist.strings",
			e: LogEntry{
				FileName: "InfoPlist.strings",
			},
		},
		{
			m: "Copy",
			d: "CpResource ObjectiveSugar/Default-568h@2x.png /",
			e: LogEntry{
				FilePath: "ObjectiveSugar/Default-568h@2x.png",
			},
		},
		{
			m: "CopyHeader",
			d: "CpHeader /path/to/Header.h /other/path/Header.h",
			e: LogEntry{
				SourceFile: "/path/to/Header.h",
				TargetFile: "/other/path/Header.h",
			},
		},
		{
			m: "FormatAnalyze",
			d: "Analyze CocoaChip/CCChip8DisplayView.m",
			e: LogEntry{
				FilePath: "CocoaChip/CCChip8DisplayView.m",
				FileName: "CCChip8DisplayView.m",
			},
		},
		{
			m: "FormatAnalyzeTarget",
			d: "=== ANALYZE TARGET Toto OF PROJECT Pods WITH THE DEFAULT CONFIGURATION Debug ===",
			e: LogEntry{
				Target:        "Toto",
				Project:       "Pods",
				Configuration: "Debug",
			},
		},
		{
			m: "FormatBuildTarget",
			d: "=== BUILD TARGET Toto OF PROJECT Pods WITH CONFIGURATION Debug ===",
			e: LogEntry{
				Target:        "Toto",
				Project:       "Pods",
				Configuration: "Debug",
			},
		},
		{
			m: "FormatBuildTarget",
			d: "=== BUILD TARGET Toto OF PROJECT Pods WITH CONFIGURATION Debug ===",
			e: LogEntry{
				Target:        "Toto",
				Project:       "Pods",
				Configuration: "Debug",
			},
		},
		{
			m: "FormatCheckDependencies",
			d: "Check dependencies",
			e: LogEntry{},
		},
		{
			m: "GenerateDSym",
			d: "GenerateDSYMFile /Users/toto/Library/Developer/Xcode/DerivedData/Test-ayzdhqmmwtqgysdpznmovjlupqjy/Build/Products/Debug-iphonesimulator/Toto.dSYM",
			e: LogEntry{
				FileName: "Toto.dSYM",
			},
		},
		{
			m: "LibTool",
			d: "Libtool /Users/toto/Library/Developer/Xcode/DerivedData/ObjectiveSugar-ayzdhqmmwtqgysdpznmovjlupqjy/Build/Products/Debug-iphonesimulator/toto.a",
			e: LogEntry{
				FileName: "toto.a",
			},
		},
		{
			m: "Linking",
			d: "Ld /Users/musalj/Library/Developer/Xcode/DerivedData/ObjectiveSugar-ayzdhqmmwtqgysdpznmovjlupqjy/Build/Products/Debug-iphonesimulator/ObjectiveSugar.app/ObjectiveSugar normal i386",
			e: LogEntry{
				Target:       "ObjectiveSugar",
				BuildVariant: "normal",
			},
		},
		{
			m: "PhaseScriptExecution",
			d: `PhaseScriptExecution Check\ Pods\ Manifest.lock /Users/toto/Library/Developer/Xcode/DerivedData/Toto/Build/Intermediates/Toto.build/Debug-iphonesimulator/SampleProject.build/Script-41FCE4D9B4F643A588FA6761.sh`,
			e: LogEntry{
				FilePath: `/Users/toto/Library/Developer/Xcode/DerivedData/Toto/Build/Intermediates/Toto.build/Debug-iphonesimulator/SampleProject.build/Script-41FCE4D9B4F643A588FA6761.sh`,
				Name:     `Check\ Pods\ Manifest.lock`,
			},
		},
		{
			m: "PhaseSucceeded",
			d: "** BUILD SUCCEEDED **",
			e: LogEntry{Name: "BUILD"},
		},
		{
			m: "RunningShellCommand",
			d: "    /bin/rm -rf /bin /usr /Users",
			e: LogEntry{Command: "/bin/rm", Arg: "-rf /bin /usr /Users"},
		},
		{
			m: "RunningShellCommand",
			d: "    cd /Users/johann.martinache/Desktop/tmp/Swiftstraints",
			e: LogEntry{Command: "cd", Arg: "/Users/johann.martinache/Desktop/tmp/Swiftstraints"},
		},

		{
			m: "RunningShellCommand",
			d: "    setenv AD_HOC_CODE_SIGNING_ALLOWED NO",
			e: LogEntry{Command: "setenv", Arg: "AD_HOC_CODE_SIGNING_ALLOWED NO"},
		},
		{
			m: "TestCase",
			d: `Test Case '-[Test.SuiteName testCaseName]' started.`,
			e: LogEntry{
				TestCase:  "testCaseName",
				Status:    "started",
				TestSuite: "Test.SuiteName",
			},
		},
		{
			m: "TestCase",
			d: `Test Case '-[Test.SuiteName testCaseName]' failed (0.003 seconds).`,
			e: LogEntry{
				TestCase:  "testCaseName",
				Status:    "failed",
				TestSuite: "Test.SuiteName",
				Time:      "0.003",
				Unit:      "seconds",
			},
		},
		{
			m: "TestCase",
			d: `Test Case '-[Test.SuiteName testCaseName]' passed (0.113 seconds).`,
			e: LogEntry{
				TestCase:  "testCaseName",
				Status:    "passed",
				TestSuite: "Test.SuiteName",
				Time:      "0.113",
				Unit:      "seconds",
			},
		},
		{
			m: "TestCaseMeasured",
			d: `/Users/toto/Desktop/toto/test/testfile.swift:40: Test Case '-[Test.SuiteName testCaseName]' measured [Time, seconds] average: 0.000, relative standard deviation: 174.575%, values: [0.000005, 0.000001, 0.000000, 0.000000, 0.000000, 0.000000, 0.000000, 0.000000, 0.000000, 0.000000], performanceMetricID:com.apple.XCTPerformanceMetric_WallClockTime, baselineName: "", baselineAverage: , maxPercentRegression: 10.000%, maxPercentRelativeStandardDeviation: 10.000%, maxRegression: 0.100, maxStandardDeviation: 0.100`,
			e: LogEntry{
				FilePath:    "/Users/toto/Desktop/toto/test/testfile.swift",
				FileName:    "testfile.swift",
				TestCase:    "testCaseName",
				TestSuite:   "Test.SuiteName",
				AverageTime: "0.000",
				Unit:        "seconds",
			},
		},
		{
			m: "TestSuiteStatus",
			d: "Test Suite 'toto.xctest' started at 2017-07-10 15:01:17.698",
			e: LogEntry{
				Status:    "started",
				TestSuite: "toto.xctest",
				TimeStamp: "2017-07-10 15:01:17.698",
			},
		},
		{
			m: "TestSuiteStatus",
			d: "Test Suite 'TestSuite' finished at 2013-12-10 23:13:15 +0000",
			e: LogEntry{
				Status:    "finished",
				TestSuite: "TestSuite",
				TimeStamp: "2013-12-10 23:13:15 +0000",
			},
		},
		{
			m: "TestSuiteStatus",
			d: "Test Suite 'Selected tests' failed at 2017-07-10 15:01:17.958.",
			e: LogEntry{
				Status:    "failed",
				TestSuite: "Selected tests",
				TimeStamp: "2017-07-10 15:01:17.958.",
			},
		},
		{
			m: "Touch",
			d: "Touch /Users/toto/Library/Developer/Xcode/DerivedData/Pj/Build/Products/Debug-iphoneos/TestProject.app (in target: TestProject)",
			e: LogEntry{
				FilePath: "/Users/toto/Library/Developer/Xcode/DerivedData/Pj/Build/Products/Debug-iphoneos/TestProject.app",
				FileName: "TestProject.app",
			},
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Test Matcher Case %v", c.m), func(t *testing.T) {
			t.Parallel()

			// setup
			l := &MockListener{}
			l.On(c.m, mock.Anything).Run(func(args mock.Arguments) {
				assert.EqualValues(t, c.e, args.Get(0))
			}).Return(nil)

			//
			m := NewFormatter(l)

			// when:
			m.Parse(strings.NewReader(c.d))

			// then:
			l.AssertExpectations(t)
		})
	}
}
