// +build mock

package api

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPlistAPI struct {
	mock.Mock
}

func (m *MockPlistAPI) AddStringValue(ctx context.Context, objectId string, path string, value string) error {
	c := m.Called(ctx, objectId, path, value)
	return c.Error(0)
}

func (m *MockPlistAPI) SetStringValue(ctx context.Context, objectId string, path string, value string) error {
	c := m.Called(ctx, objectId, path, value)
	return c.Error(0)
}
