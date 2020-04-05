package xcode

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockXCodeBuildService struct {
	mock.Mock
}

func (m *MockXCodeBuildService) List(ctx context.Context) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockXCodeBuildService) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	args := m.Called(scheme)
	return args.String(0), args.Error(1)
}

func (m *MockXCodeBuildService) Run(ctx context.Context, arg ...string) (string, error) {
	args := m.Called(arg)
	return args.String(0), args.Error(1)
}
