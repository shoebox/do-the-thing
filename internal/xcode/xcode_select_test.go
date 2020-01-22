package xcode

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"

	"dothething/internal/utiltest"
)

type MockListService struct {
	mock.Mock
}

func (m *MockListService) List() ([]*Install, error) {
	args := m.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	res2 := args.Get(0)
	s, ok := res2.([]*Install)
	if ok {
		return s, nil
	}
	return nil, nil
}

var subject *XCodeSelectService
var exec *utiltest.MockExec
var mockListService *MockListService

func setupSelectTest() {
	exec = new(utiltest.MockExec)
	mockListService = new(MockListService)
	subject = NewSelectService(mockListService, exec)
}

func mockResults() {
	mockListService.On("List").Return([]*Install{
		&Install{Path: "/Applications/Xcode 1.2.3.app",
			BundleVersion: "1.2.3-snapshot",
			Version:       "1.2.3"},

		&Install{Path: "/Applications/Xcode 10.3.1app",
			BundleVersion: "10.3.1-snapshot",
			Version:       "10.3.1"},

		&Install{Path: "/Applications/Xcode.app",
			BundleVersion: "7.1-snapshot",
			Version:       "7.1"},
	}, nil)

}

func TestVersionSelection(t *testing.T) {
	t.Run("If a error happens during the resolution, it should be reported", func(t *testing.T) {
		setupSelectTest()
		mockListService.On("List").
			Return(nil, fmt.Errorf("Listing error"))

		install, err := subject.SelectVersion("1.2.3")
		assert.Nil(t, install)
		assert.Error(t, err, "Listing error")
	})

	t.Run("Should properly select the version if available", func(t *testing.T) {
		setupSelectTest()
		mockResults()
		installl, err := subject.SelectVersion("10.3.1")
		assert.NoError(t, err)
		assert.NotNil(t, installl)
	})

	t.Run("When could not find an instance should throw an error", func(t *testing.T) {
		setupSelectTest()
		mockResults()
		installl, err := subject.SelectVersion("10.3.3")
		assert.EqualValues(t, ErrMatchNotFound, err)
		assert.Nil(t, installl)
	})

	t.Run("Should expect an error in case of invalid version selection", func(t *testing.T) {
		setupSelectTest()
		mockResults()

		install, err := subject.SelectVersion("Toto")
		assert.EqualError(t, err, "Invalid version")
		assert.Nil(t, install)
	})
}

func TestIsEqualMatch(t *testing.T) {
	install := Install{Version: "10.2.3"}

	t.Run("Matching should work as expected", func(t *testing.T) {
		r, err := subject.isEqualMatch(&install, "10.2.3")
		assert.True(t, r)
		assert.NoError(t, err)

		r, err = subject.isEqualMatch(&install, "10.2.0")
		assert.False(t, r)
		assert.NoError(t, err)
	})

	t.Run("Error should be handled", func(t *testing.T) {
		r, err := subject.isEqualMatch(&install, "")
		assert.Error(t, err, "t")
		assert.False(t, r)
	})
}

func TestFindMatch(t *testing.T) {
	t.Run("Valid match", func(t *testing.T) {
		setupSelectTest()
		mockResults()

		install, err := subject.findMatch("10.3.1", subject.isEqualMatch)
		assert.NotNil(t, install)
		assert.NoError(t, err, "")
	})

	t.Run("should handle error", func(t *testing.T) {
		setupSelectTest()
		mockListService.On("List").
			Return(nil, fmt.Errorf("Listing error"))

		install, err := subject.findMatch("10.3.1", subject.isEqualMatch)
		assert.Nil(t, install)
		assert.EqualError(t, err, "Listing error")
	})
}
