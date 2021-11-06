package util

import (
	"context"
	"dothething/internal/api"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAPI struct {
	mock.Mock
	api.API
	XcodeSelectService MockSelectService
}

type MockSelectService struct {
	mock.Mock
}

func (m MockSelectService) Find(ctx context.Context) (*api.Install, error) {
	args := m.Called(ctx)
	return args.Get(0).(*api.Install), args.Error(1)
}

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup:
	install := &api.Install{Path: "/path/to/xcode"}
	service := new(MockSelectService)
	service.On("Find", mock.Anything).Return(install, nil)

	// when:
	e := NewExecutor(&api.API{XcodeSelectService: service})
	cmd, err := e.XCodeCommandContext(ctx, "hello", "world")

	// then:
	assert.NoError(t, err)
	w := (*cmd).(*cmdWrapper)
	assert.Equal(t, "DEVELOPER_DIR", w.Env[0])
}
