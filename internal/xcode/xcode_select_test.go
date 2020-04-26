package xcode

/*
type SelectServiceTest struct {
	suite.Suite
	mockExec *utiltest.MockExec
	subject  selectService
}

func (s *SelectServiceTest) SetupTest() {
	utiltest.SetupMockExec()
	s.mockExec = utiltest.Exec
	s.mockFileService = new(utiltest.MockFileService)
	s.service = listService{exec: s.mockExec, file: s.mockFileService}
}

func TestSelectServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SelectServiceTest))
}

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

var subject SelectService
var exec *utiltest.MockExecutor2
var mockListService *MockListService

func setupSelectTest() {
	exec = new(utiltest.MockExecutor2)
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
*/

/*
func TestVersionSelection(t *testing.T) {
	t.Run("Should handle invalid requirement", func(t *testing.T) {
		install, err := subject.Find("Invalid")
		assert.Nil(t, install)
		assert.EqualError(t, err, ErrInvalidVersion.Error())
	})

	t.Run("If a error happens during the resolution, it should be reported", func(t *testing.T) {
		setupSelectTest()
		mockListService.On("List").
			Return(nil, fmt.Errorf("Listing error"))

		install, err := subject.Find("1.2.3")
		assert.Nil(t, install)
		assert.Error(t, err, "Listing error")
	})

	t.Run("Should properly select the version if available", func(t *testing.T) {
		setupSelectTest()
		mockResults()
		installl, err := subject.Find("10.3.1")
		assert.NoError(t, err)
		assert.NotNil(t, installl)
	})

	t.Run("When could not find an instance should throw an error", func(t *testing.T) {
		setupSelectTest()
		mockResults()
		installl, err := subject.Find("10.3.3")
		assert.EqualValues(t, ErrMatchNotFound, err)
		assert.Nil(t, installl)
	})

	t.Run("Target resolution", func(t *testing.T) {
		setupSelectTest()
		mockResults()
		installl, err := subject.Find("10.3.1")
		fmt.Println("install", installl)
		assert.Nil(t, err)
		assert.NotNil(t, installl)
		assert.EqualValues(t, installl.Version, "10.3.1")
	})
}

func TestIsEqualMatch(t *testing.T) {
	install := Install{Version: "10.2.3"}

	t.Run("Matching should work as expected", func(t *testing.T) {
		v1, _ := semver.Parse("10.2.3")
		v2, _ := semver.Parse("10.2.0")

		r, err := subject.isEqualMatch(&install, v1)
		assert.True(t, r)
		assert.NoError(t, err)

		r, err = subject.isEqualMatch(&install, v2)
		assert.False(t, r)
		assert.NoError(t, err)
	})

	t.Run("Error should be handled", func(t *testing.T) {
		v3, _ := semver.Parse("")
		r, err := subject.isEqualMatch(&install, v3)
		assert.NoError(t, err)
		assert.False(t, r)
	})
}

func TestFindMatch(t *testing.T) {

	required, err := semver.Parse("10.3.1")
	assert.Nil(t, err)

	t.Run("Valid match", func(t *testing.T) {
		setupSelectTest()
		mockResults()

		install, err := subject.findMatch(required, subject.isEqualMatch)
		assert.NotNil(t, install)
		assert.NoError(t, err, "")
	})

	t.Run("should handle error", func(t *testing.T) {
		setupSelectTest()
		mockListService.On("List").
			Return(nil, fmt.Errorf("Listing error"))

		install, err := subject.findMatch(required, subject.isEqualMatch)
		assert.Nil(t, install)
		assert.EqualError(t, err, "Listing error")
	})

	t.Run("Should return ErrMatchNotFound in case of zero matchs", func(t *testing.T) {
		setupSelectTest()

		inst := &Install{}
		mockListService.On("List").Return([]*Install{inst}, nil)

		required2, err := semver.Parse("10.1.0")
		assert.Nil(t, err)

		install, err := subject.findMatch(required2, subject.isEqualMatch)
		fmt.Println(install, err)
		assert.Nil(t, install)
		assert.EqualError(t, err, ErrMatchNotFound.Error())
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
*/
