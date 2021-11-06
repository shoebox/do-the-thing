// +build mock

package api

import "github.com/stretchr/testify/mock"

type PathMock struct {
	mock.Mock
}

func (p *PathMock) Archive() string {
	c := p.Called()
	return c.String(0)
}

func (p *PathMock) ExportPList() string {
	c := p.Called()
	return c.String(0)
}

func (p *PathMock) KeyChain() string {
	c := p.Called()
	return c.String(0)
}

func (p *PathMock) ObjRoot() string {
	c := p.Called()
	return c.String(0)
}

func (p *PathMock) SymRoot() string {
	c := p.Called()
	return c.String(0)
}

func (p *PathMock) XCResult() string {
	c := p.Called()
	return c.String(0)
}
