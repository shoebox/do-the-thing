// +build mock

package api

import (
	"github.com/stretchr/testify/mock"
)

type APIMock struct {
	mock.Mock
}

func (m *APIMock) ActionArchive() Action {
	a := m.Called()
	return a.Get(0).(Action)
}

func (m *APIMock) ActionBuild() Action {
	a := m.Called()
	return a.Get(0).(Action)
}

func (m *APIMock) ActionRunTest() Action {
	a := m.Called()
	return a.Get(0).(Action)
}

func (m *APIMock) CertificateService() CertificateService {
	a := m.Called()
	return a.Get(0).(CertificateService)
}

func (m *APIMock) DestinationService() DestinationService {
	a := m.Called()
	return a.Get(0).(DestinationService)
}

func (m *APIMock) Exec() Executor {
	a := m.Called()
	return a.Get(0).(Executor)
}

func (m *APIMock) FileService() FileService {
	a := m.Called()
	return a.Get(0).(FileService)
}

func (m *APIMock) KeyChainService() KeyChain {
	a := m.Called()
	return a.Get(0).(KeyChain)
}

func (m *APIMock) PListBuddyService() PListBuddyService {
	a := m.Called()
	return a.Get(0).(PListBuddyService)
}

func (m *APIMock) ProvisioningService() ProvisioningService {
	a := m.Called()
	return a.Get(0).(ProvisioningService)
}

func (m *APIMock) SignatureResolver() SignatureResolver {
	a := m.Called()
	return a.Get(0).(SignatureResolver)
}

func (m *APIMock) SignatureService() SignatureService {
	a := m.Called()
	return a.Get(0).(SignatureService)
}

func (m *APIMock) XCodeBuildService() BuildService {
	a := m.Called()
	return a.Get(0).(BuildService)
}

func (m *APIMock) XCodeListService() ListService {
	a := m.Called()
	return a.Get(0).(ListService)
}

func (m *APIMock) XCodeProjectService() ProjectService {
	a := m.Called()
	return a.Get(0).(ProjectService)
}

func (m *APIMock) XCodeSelectService() SelectService {
	a := m.Called()
	return a.Get(0).(SelectService)
}
