package project

import (
	"context"
	"dothething/internal/utiltest"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"dothething/internal/xcode"
	"dothething/internal/xcode/pbx"

	"github.com/stretchr/testify/assert"
)

var ps projectService

func TestParsing(t *testing.T) {
	// setup:
	f, err := os.Open("testdata/project.pbxproj")
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)

	// when:
	pj, err := ps.decodeProject(b)

	// then:
	assert.EqualValues(t, pj.Targets[0],
		makeTarget(
			"Swiftstraints iOS",
			pbx.Framework,
			"Swiftstraints",
			4,
			map[string]string{
				"CLANG_ENABLE_MODULES":              "YES",
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "",
				"DEFINES_MODULE":                    "YES",
				"DYLIB_COMPATIBILITY_VERSION":       "1",
				"DYLIB_CURRENT_VERSION":             "1",
				"DYLIB_INSTALL_NAME_BASE":           "@rpath",
				"INFOPLIST_FILE":                    "Swiftstraints/Info.plist",
				"INSTALL_PATH":                      "/Frameworks",
				"IPHONEOS_DEPLOYMENT_TARGET":        "8.0",
				"LD_RUNPATH_SEARCH_PATHS":           " @executable_path/Frameworks @loader_path/Frameworks",
				"PRODUCT_BUNDLE_IDENTIFIER":         "com.skyvive.$(PRODUCT_NAME:rfc1034identifier)",
				"PRODUCT_NAME":                      "Swiftstraints",
				"SKIP_INSTALL":                      "YES",
				"SWIFT_OPTIMIZATION_LEVEL":          "-Onone",
				"SWIFT_SWIFT3_OBJC_INFERENCE":       "Default",
				"SWIFT_VERSION":                     "5.0",
			},
			map[string]string{
				"CLANG_ENABLE_MODULES":              "YES",
				"DEFINES_MODULE":                    "YES",
				"DYLIB_COMPATIBILITY_VERSION":       "1",
				"DYLIB_CURRENT_VERSION":             "1",
				"DYLIB_INSTALL_NAME_BASE":           "@rpath",
				"INFOPLIST_FILE":                    "Swiftstraints/Info.plist",
				"INSTALL_PATH":                      "/Frameworks",
				"IPHONEOS_DEPLOYMENT_TARGET":        "8.0",
				"LD_RUNPATH_SEARCH_PATHS":           " @executable_path/Frameworks @loader_path/Frameworks",
				"PRODUCT_BUNDLE_IDENTIFIER":         "com.skyvive.$(PRODUCT_NAME:rfc1034identifier)",
				"PRODUCT_NAME":                      "Swiftstraints",
				"SKIP_INSTALL":                      "YES",
				"SWIFT_SWIFT3_OBJC_INFERENCE":       "Default",
				"SWIFT_VERSION":                     "5.0",
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]": "",
			},
		),
	)

	// and:
	assert.NoError(t, err)
}

func makeTarget(name string,
	t string,
	productName string,
	phaseCount int,
	mDebug map[string]string,
	mRelease map[string]string,
) pbx.NativeTarget {

	phases := make([]pbx.PBXBuildPhase, phaseCount)

	return pbx.NativeTarget{
		BuildConfigurationList: pbx.XCConfigurationList{
			Reference: "",
			BuildConfiguration: []pbx.XCBuildConfiguration{
				{Name: "Debug", BuildSettings: mDebug},
				{Name: "Release", BuildSettings: mRelease},
			},
			DefaultConfigurationVisible: 0,
			DefaultConfigurationName:    "Release",
		},
		BuildPhases:        phases,
		Name:               name,
		ProductInstallPath: "",
		ProductName:        productName,
		ProductType:        t,
	}
}

func TestCases(t *testing.T) {
	params := []struct {
		execErr       error
		expectedError error
		path          string
		list          *list
		project       *Project
		workspace     *Project
	}{
		{
			expectedError: ErrInvalidConfig,
		},
		{
			project: &projectSample.Project,
			list:    &projectSample,
		},
		{
			workspace: &workspaceSample.Workspace,
			list:      &workspaceSample,
		},
		{
			execErr:       errors.New("Error calling xcode"),
			expectedError: xcode.NewError(-1),
		},
	}

	const fakePath = "/path/to/project.xcodeproj"

	for index, tc := range params {
		t.Run(fmt.Sprintf("Test case %v", index), func(t *testing.T) {
			// setup:
			exec := new(utiltest.MockExecutor)
			subject := projectService{xcodeService: xcode.NewService(exec, fakePath)}

			raw := "invalid json"
			if tc.list != nil {
				b, err := json.Marshal(tc.list)
				assert.NoError(t, err)
				raw = string(b)
			}

			exec.MockCommandContext(
				xcode.Build,
				[]string{xcode.FlagList, xcode.FlagJSON, xcode.FlagProject, fakePath},
				raw,
				tc.execErr)

			// when: Resolving project
			p, err := subject.resolveProject(context.Background())
			assert.EqualValues(t, tc.expectedError, err)

			// then: Should return the project defintion
			if tc.project != nil {
				assert.EqualValues(t, tc.project, &p)
			}

			// and: Should return a workspace in that case
			if tc.workspace != nil {
				assert.EqualValues(t, tc.workspace, &p)
			}
		})
	}
}

var projectSample = list{Project: project}
var workspaceSample = list{Workspace: project}

var project = Project{
	Configurations: []string{"Config1", "Config2"},
	Name:           "projectName",
	Schemes:        []string{"scheme1", "scheme2"},
	Targets:        []string{"target1", "target2"},
}
