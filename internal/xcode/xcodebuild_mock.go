// +build mock

package xcode

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type XCodeBuildMock struct {
	mock.Mock
}

func (x *XCodeBuildMock) List(ctx context.Context) (string, error) {
	c := x.Called(ctx)
	return c.Get(0).(string), c.Error(1)
}

func (x *XCodeBuildMock) ShowDestinations(ctx context.Context, scheme string) (string, error) {
	c := x.Called(ctx, scheme)
	return c.Get(0).(string), c.Error(1)
}

func (x *XCodeBuildMock) GetArg() string {
	c := x.Called()
	return c.Get(0).(string)
}

func (x *XCodeBuildMock) GetProjectPath() string {
	c := x.Called()
	return c.Get(0).(string)
}
