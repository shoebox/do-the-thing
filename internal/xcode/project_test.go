package xcode

import (
	"context"
	"dothething/internal/utiltest"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

const validResponse string = `{
				"project" : {
					"configurations" : [
						"Production Debug",
						"Production Release",
						"Staging Debug",
						"Staging Release",
						"Test Debug",
						"Test Release"
					],
					"name" : "Sample",
					"schemes" : [
						"ActionPopoverButton",
						"Sample",
						"SampleTests Prod",
						"SampleTests Stag",
						"SampleTests Test",
						"TopShelfExtension"
					],
					"targets" : [
						"ActionPopoverButton",
						"Sample",
						"SampleTests",
						"TopShelfExtension"
					]
				}
			}`

const fakePath = "/path/to/project.xcodeproj"

type projectTestSuite struct {
	suite.Suite
	ctx     context.Context
	exec    *utiltest.MockExecutor
	subject projectService
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(projectTestSuite))
}

func (s *projectTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.exec = new(utiltest.MockExecutor)
	s.subject = projectService{xcodeService: NewService(s.exec, fakePath)}
}

func (s *projectTestSuite) TestCases() {
	params := []struct {
		execErr       error
		expectedError error
		json          string
		path          string
		expectedValue *Project
	}{
		{
			execErr:       nil,
			json:          "invalid json",
			expectedError: ErrInvalidConfig,
		},
		{
			execErr:       nil,
			json:          validResponse,
			expectedError: nil,
			expectedValue: &Project{
				Configurations: []string{
					"Production Debug",
					"Production Release",
					"Staging Debug",
					"Staging Release",
					"Test Debug",
					"Test Release",
				},
				Name: "sample",
				Targets: []string{
					"ActionPopoverButton",
					"Sample",
					"SampleTests",
					"TopShelfExtension",
				},
				Schemes: []string{
					"ActionPopoverButton",
					"Sample",
					"SampleTests Prod",
					"SampleTests Stag",
					"SampleTests Test",
					"TopShelfExtension",
				},
			},
		},
		{
			execErr:       errors.New("Error calling xcode"),
			json:          "invalid json",
			expectedError: XCodebuildError{-1, -1},
		},
	}

	for index, tc := range params {
		s.Run(fmt.Sprintf("Test case %v", index), func() {
			s.SetupTest()
			s.exec.MockCommandContext(
				XCodeBuild,
				[]string{flagList, flagJSON, FlagProject, fakePath},
				tc.json,
				tc.execErr)

			p, err := s.subject.Parse(s.ctx)

			s.Assert().EqualValues(tc.expectedError, err)
			if tc.expectedValue != nil {
				s.Assert().EqualValues(tc.expectedValue.Configurations, p.Configurations)
				s.Assert().EqualValues(tc.expectedValue.Schemes, p.Schemes)
				s.Assert().EqualValues(tc.expectedValue.Targets, p.Targets)
			} else {
				s.Assert().Nil(p)
			}
		})
	}
}
