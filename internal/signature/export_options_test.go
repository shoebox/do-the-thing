package signature

import (
	"dothething/internal/api"
	"testing"

	"github.com/stretchr/testify/suite"
)

type exportOptionsPlistSuite struct {
	suite.Suite
	API     *api.API
	subject *exportOptionsService
}

// suite testing entrypoint
func TestOptionalPlistSuite(t *testing.T) {
	suite.Run(t, new(exportOptionsPlistSuite))
}

/*
func (s *exportOptionsPlistSuite) BeforeTest(suiteName, testName string) {
	s.API = &api.API{
		PlistBuddyService: new(api.MockPlistAPI),
		SignatureService:  new(api.SignatureServiceMock),
	}

	s.subject = &exportOptionsService{API: s.API}
}
func (s *exportOptionsPlistSuite) TestResolveMethodForProvisioning() {
	t := true
	cases := []struct {
		pd  *[]string
		ta  bool
		all *bool
		res string
	}{
		{pd: &[]string{"UUID"}, ta: true, res: "development"},
		{pd: &[]string{"UUID"}, ta: false, res: "ad-hoc"},
		{all: &t, res: "enterprise"},
		{res: "app-store"},
	}

	for _, c := range cases {
		fmt.Println(c)
		s.Run(c.res, func() {
			// setup
			p := api.ProvisioningProfile{
				ProvisionedDevices:   c.pd,
				ProvisionsAllDevices: c.all,
			}
			p.Entitlements.GetAskAllow = c.ta

			// when:
			m := s.subject.resolveMethodForProvisioning(&p)

			// then:
			s.Assert().EqualValues(c.res, m)
		})
	}
}
*/
