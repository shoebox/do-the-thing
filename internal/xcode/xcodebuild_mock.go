package xcode

import "github.com/stretchr/testify/mock"

type MockXCodeBuildService struct {
	mock.Mock
}

func (m *MockXCodeBuildService) List() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockXCodeBuildService) ShowDestinations(scheme string) (string, error) {
	args := m.Called(scheme)
	return args.String(0), args.Error(1)
}

func (m *MockXCodeBuildService) Run(arg ...string) (string, error) {
	args := m.Called(arg)
	return args.String(0), args.Error(1)
}
