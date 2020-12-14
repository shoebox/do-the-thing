package api

import (
	"context"

	"github.com/stretchr/testify/mock"
)

//
type ListServiceMock struct {
	mock.Mock
}

func (l *ListServiceMock) List(ctx context.Context) ([]*Install, error) {
	c := l.Called(ctx)
	return c.Get(0).([]*Install), c.Error(1)
}

//
type SelectServiceMock struct {
	mock.Mock
}

func (m *SelectServiceMock) Find(ctx context.Context, version string) (*Install, error) {
	c := m.Called(ctx, version)
	return c.Get(0).(*Install), c.Error(1)
}

//
type BuildServiceMock struct {
	mock.Mock
}

func (m *BuildServiceMock) List(ctx context.Context) (string, error) {
	c := m.Called(ctx)
	return c.String(0), c.Error(1)
}

func (m *BuildServiceMock) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	c := m.Called(ctx, scheme)
	return c.String(0), c.Error(1)
}

func (m *BuildServiceMock) GetArg() string {
	c := m.Called()
	return c.String(0)
}

func (m *BuildServiceMock) GetProjectPath() string {
	c := m.Called()
	return c.String(0)
}
