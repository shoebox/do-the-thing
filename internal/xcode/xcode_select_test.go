package xcode

import (
	"context"
	"dothething/internal/utiltest"
	"errors"

	"testing"

	"github.com/blang/semver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
)

type selectSuite struct {
	suite.Suite
	listContext     context.Context
	mockExec        *utiltest.MockExecutor2
	mockFileService *utiltest.MockFileService
	subject         selectService
	listService     *mockListService
}

type mockListService struct {
	mock.Mock
}

func (m *mockListService) List(ctx context.Context) ([]*Install, error) {
	r := m.Called(ctx)
	return r.Get(0).([]*Install), r.Error(1)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestXCodeSelectSuite(t *testing.T) {
	suite.Run(t, new(selectSuite))
}

func (s *selectSuite) SetupTest() {
	s.mockExec = new(utiltest.MockExecutor2)
	s.listContext = context.Background()
	s.mockFileService = new(utiltest.MockFileService)
	s.listService = new(mockListService)
	s.subject = selectService{Executor: s.mockExec, list: s.listService}
}

var install1 = Install{
	Path:          "",
	BundleVersion: "",
	Version:       "11.3.1",
}

func (s *selectSuite) TestErroHandlingWhileListing() {
	s.listService.On("List", mock.Anything).
		Return([]*Install{}, errors.New("error text"))

	i, err := s.subject.Find(s.listContext, "10.0.0")
	s.Assert().Nil(i)
	s.Assert().EqualError(err, ErrMatchNotFound.Error())
}

func (s *selectSuite) TestFind() {
	s.listService.On("List", mock.Anything).
		Return([]*Install{
			&install1,
		}, nil)

	s.Run("Should handle invalid requirement", func() {
		install, err := s.subject.Find(s.listContext, "Invalid")
		s.Assert().Nil(install)
		s.Assert().EqualError(err, ErrInvalidVersion.Error())
	})

	s.Run("Should handle invalid requirement", func() {
		install, err := s.subject.Find(s.listContext, "12.3.1")
		s.Assert().Nil(install)
		s.Assert().EqualError(err, ErrMatchNotFound.Error())
	})

	s.Run("Should handle invalid requirement", func() {
		install, err := s.subject.Find(s.listContext, "11.3.1")
		s.Assert().NoError(err)
		s.Assert().EqualValues(&install1, install)
	})
}

func (s *selectSuite) TestIsEqualMatch() {
	install := Install{Version: "10.2.3"}

	s.Run("Matching should work as expected", func() {
		v1, _ := semver.Parse("10.2.3")
		v2, _ := semver.Parse("10.2.0")

		r, err := s.subject.isEqualMatch(&install, v1)
		s.Assert().True(r)
		s.Assert().NoError(err)

		r, err = s.subject.isEqualMatch(&install, v2)
		s.Assert().False(r)
		s.Assert().NoError(err)
	})

	s.Run("Error should be handled", func() {
		v3, _ := semver.Parse("")
		r, err := s.subject.isEqualMatch(&install, v3)
		s.Assert().NoError(err)
		s.Assert().False(r)
	})
}

func TestSortInstalls(t *testing.T) {

	t.Run("Simple sorting", func(t *testing.T) {
		install1 := Install{Version: "10.1.2"}
		install2 := Install{Version: "10.0.1"}
		install3 := Install{Version: "10.1.1"}

		c := append([]*Install{}, &install1, &install2, &install3)

		sortInstalls(c)

		assert.Equal(t, c[0], &install1)
		assert.Equal(t, c[1], &install3)
		assert.Equal(t, c[2], &install2)
	})

	t.Run("Error parsing should be handle and result sorted down", func(t *testing.T) {
		install1 := Install{Version: "10.1.2"}
		install2 := Install{Version: "10.0.1"}
		install3 := Install{Version: "10.1.1"}
		empty := Install{}

		c := append([]*Install{}, &install1, &install2, &install3, &empty, &empty)
		sortInstalls(c)

		assert.Equal(t, c[0], &install1)
		assert.Equal(t, c[1], &install3)
		assert.Equal(t, c[2], &install2)
		assert.Equal(t, c[3], &empty)
	})

}

func TestCompareInstall(t *testing.T) {
	empty := Install{}
	install := Install{Version: "10.1.2"}

	t.Run("Should handle invalid install1", func(t *testing.T) {
		res := compareInstall(&empty, &install)
		assert.False(t, res)
	})

	t.Run("Should handle invalid install2", func(t *testing.T) {
		res := compareInstall(&install, &empty)
		assert.False(t, res)
	})
}
