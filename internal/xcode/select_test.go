package xcode

import (
	"context"
	"dothething/internal/api"
	"errors"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockListService struct{ mock.Mock }

func (l mockListService) List(ctx context.Context) ([]*api.Install, error) {
	a := l.Called(ctx)
	return a.Get(0).([]*api.Install), a.Error(1)
}

type selectSuite struct {
	suite.Suite
	subject selectService
	ls      *mockListService
	API     *api.API
}

func (s *selectSuite) SetupTest() {
	s.ls = new(mockListService)
	s.API = &api.API{
		XcodeListService: s.ls,
	}

	//
	s.subject = selectService{API: s.API}
	s.subject.API = s.API
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestXCodeSelectSuite(t *testing.T) {
	suite.Run(t, new(selectSuite))
}

var i1 = api.Install{Version: "13.1.2"}
var i2 = api.Install{Version: "12.0.9"}
var i3 = api.Install{Version: "16.0.9"}

func (s *selectSuite) TestErroHandlingWhileListing() {

	// setup:

	cases := []struct {
		req      string
		installs []*api.Install
		res      *api.Install
		err      error
		lerr     error
	}{
		{installs: []*api.Install{&i1, &i2, &i3}, res: &i3, err: nil, req: ">=16.0.0"},
		{installs: []*api.Install{&i1, &i2, &i3}, res: &i3, err: nil, req: ">12.0.0"},
		{installs: []*api.Install{&i1, &i2, &i3}, res: &i1, err: nil, req: ">12.0.0 <14.0.0"},
		{installs: []*api.Install{&i1, &i2, &i3}, res: nil, err: ErrMatchNotFound, req: "<4.0.0"},
		{installs: []*api.Install{&i1, &i2, &i3}, res: nil, err: ErrParsing, req: "wrong"},
	}

	for _, tc := range cases {
		s.Run("-", func() {
			// setup:
			s.ls.On("List", mock.Anything).Return(tc.installs, tc.lerr)

			// when:
			i, err := s.subject.Find(context.Background(), tc.req)

			// then:
			s.Assert().EqualValues(tc.err, err)

			// and:
			s.Assert().EqualValues(tc.res, i)
		})
	}
}

func (s *selectSuite) TestIsEqualMatch() {
	cases := []struct {
		av   string // available version
		req  string // required version
		res  bool   // exepecteds result
		err  error  // error or not ?
		name string
	}{
		{name: "range - matching", req: ">10.2.0 <10.2.4", av: "10.2.3", res: true},
		{name: "range - not matching", req: ">10.2.0 <10.2.4", av: "10.3.3", res: false},
		{name: "absolute match - not matching", req: "=10.2.3", av: "10.3.3", res: false},
		{name: "absolute match - matching", req: "=10.2.3", av: "10.2.3", res: true},
		{name: "absolute match - matching", req: "=10.2.3", av: "toto", res: false, err: errors.New("No Major.Minor.Patch elements found")},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			// setup: parsing the range
			rg, err := semver.ParseRange(c.req)
			s.Assert().NoError(err)

			// creating the install mock
			i := api.Install{Version: c.av}

			// when:
			r, err := s.subject.isMatchingRequirement(&i, rg)

			// then:
			s.Assert().Equal(c.err, err)
			s.Assert().Equal(c.res, r)
		})
	}
}

func TestCompareInstall(t *testing.T) {
	cases := []struct {
		Name string
		V1   api.Install
		V2   api.Install
		res  bool
	}{
		{Name: "Invalid Install", V1: api.Install{}, V2: api.Install{Version: "10.1.3"}, res: false},
		{Name: "Invalid Install", V1: api.Install{Version: "10.1.1"}, V2: api.Install{}, res: false},
		{Name: "Superior", V1: api.Install{Version: "10.3.1"}, V2: api.Install{Version: "10.1.3"}, res: true},
		{Name: "Inferior", V1: api.Install{Version: "10.1.3"}, V2: api.Install{Version: "10.1.4"}, res: false},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			res := compareInstall(&c.V1, &c.V2)
			assert.EqualValues(t, c.res, res)
		})
	}
}

func TestSortInstalls(t *testing.T) {
	// setup:
	cases := []*api.Install{}
	c := []string{"10.1.2", "10.0.1", "10.1.1"}
	for _, v := range c {
		i := api.Install{Version: v}
		cases = append(cases, &i)
	}

	// when:
	sortInstalls(cases)

	// then:
	assert.Equal(t, cases[0].Version, "10.1.2")
	assert.Equal(t, cases[1].Version, "10.1.1")
	assert.Equal(t, cases[2].Version, "10.0.1")
}
