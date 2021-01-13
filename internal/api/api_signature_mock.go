// +build mock

package api

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type SignatureServiceMock struct {
	mock.Mock
}

func (m *SignatureServiceMock) Run(ctx context.Context) error {
	c := m.Called()
	return c.Error(0)
}

func (m *SignatureServiceMock) GetConfiguration() *[]TargetSignatureConfig {
	c := m.Called()
	return c.Get(0).(*[]TargetSignatureConfig)
}
