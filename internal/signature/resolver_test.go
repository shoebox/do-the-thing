package signature

import (
	"context"
	"crypto/x509"
	"dothething/internal/api"
	"dothething/internal/xcode/pbx"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

var (
	pp0 = api.ProvisioningProfile{
		BundleIdentifier: "com.tutu.toto",
		ExpirationDate:   time.Now().Add(19766 * time.Hour),
		Platform:         []string{"iOS"},
	}
	pp1 = api.ProvisioningProfile{
		BundleIdentifier: "*",
		ExpirationDate:   time.Now().Add(8766 * time.Hour),
		Platform:         []string{"iOS"},
	}
	pp2 = api.ProvisioningProfile{
		BundleIdentifier: "*",
		ExpirationDate:   time.Now().Add(9766 * time.Hour),
		Platform:         []string{"iOS"},
	}
	pp3 = api.ProvisioningProfile{
		BundleIdentifier: "com.tutu.toto2",
		ExpirationDate:   time.Now().Add(9766 * time.Hour),
		Platform:         []string{"tvOS"},
	}
)

type SignatureResolverSuite struct {
	suite.Suite
	subject signatureResolver
}

// suite testing entrypoint
func TestSignatureResolverSuite(t *testing.T) {
	suite.Run(t, new(SignatureResolverSuite))
}

func (s *SignatureResolverSuite) BeforeTest(suiteName, testName string) {
	s.subject = signatureResolver{}
}

func (s *SignatureResolverSuite) TestProvisioningResolutionForConfiguration() {
	// setup:
	list := []*api.ProvisioningProfile{&pp0, &pp1, &pp2, &pp3}

	cases := []struct {
		bu   string
		p    pbx.PBXProductType
		res  *api.ProvisioningProfile
		err  error
		list []*api.ProvisioningProfile
	}{
		{
			bu:   "com.tutu.toto",
			p:    pbx.DynamicLibrary,
			err:  &SignatureError{Msg: ErrorProvisioningProfileResolution},
			list: list,
		},
		{
			bu:   "com.tutu.toto",
			p:    pbx.Application,
			res:  &pp0,
			err:  nil,
			list: list,
		},
		{
			bu:   "com.toto.tutu",
			p:    pbx.Application,
			err:  &SignatureError{Msg: ErrorProvisioningProfileResolution},
			list: []*api.ProvisioningProfile{},
		},
	}

	for _, tc := range cases {
		s.Run(fmt.Sprintf("%v-%v", tc.bu, tc.p), func() {
			// when:
			res, err := s.subject.resolveProvisioningFileFor(
				context.Background(),
				tc.list,
				tc.bu,
				tc.p,
			)

			// then:
			s.Assert().EqualValues(tc.err, err)

			// and:
			s.Assert().EqualValues(tc.res, res)
		})
	}
}

func (s *SignatureResolverSuite) TestFindingMatchingCertificate() {
	cert1 := x509.Certificate{Raw: []byte("certificate1")}
	cert2 := x509.Certificate{Raw: []byte("certificate2")}
	c1 := api.P12Certificate{Certificate: &cert1}
	c2 := api.P12Certificate{Certificate: &cert2}
	l := []*api.P12Certificate{&c1, &c2}

	cases := []struct {
		raw      []byte
		expected *api.P12Certificate
		err      error
	}{
		{raw: []byte("certificate1"), expected: &c1},
		{raw: []byte("certificate2"), expected: &c2},
		{raw: []byte("test invalid"), expected: nil, err: errors.New("Could not find a matching certificate")},
	}

	for _, tc := range cases {
		// when:
		c, err := s.subject.findMatchingCert(l, tc.raw)

		// then:
		s.Assert().EqualValues(tc.expected, c)

		// and:
		s.Assert().EqualValues(tc.err, err)
	}

}

func (s *SignatureResolverSuite) TestCandidateBundleIdentifierSorting() {
	// setup:
	list := []*api.ProvisioningProfile{&pp0, &pp1, &pp2, &pp3}

	// when:
	sortBundleIdentifiers(list)

	// then:
	s.Assert().EqualValues([]*api.ProvisioningProfile{&pp3, &pp0, &pp1, &pp2}, list)
}

func (s *SignatureResolverSuite) TestFindFor() {
	pp0 := api.ProvisioningProfile{BundleIdentifier: "com.toto.*", Platform: []string{"iOS", "tvOS"}}
	pp1 := api.ProvisioningProfile{BundleIdentifier: "com.tutu.*", Platform: []string{"iOS"}}
	pp2 := api.ProvisioningProfile{BundleIdentifier: "com.tata.tutu.tvOS", Platform: []string{"tvOS"}}
	pp3 := api.ProvisioningProfile{BundleIdentifier: "com.tata.tutu.iOS", Platform: []string{"iOS"}}
	pp4 := api.ProvisioningProfile{BundleIdentifier: "com.tete.tutu.iOS", Platform: []string{"iOS"}}
	pp5 := api.ProvisioningProfile{BundleIdentifier: "com.txtx.tutu.iOS", Platform: []string{"invalid"}}
	pp6 := api.ProvisioningProfile{BundleIdentifier: "*", Platform: []string{"iOS"}}

	list := []*api.ProvisioningProfile{&pp0, &pp1, &pp2, &pp3, &pp4, &pp5, &pp6}

	cases := []struct {
		pps      []*api.ProvisioningProfile
		bu       string
		platform pbx.PBXProductType
		expected bool
		res      *api.ProvisioningProfile
	}{
		{pps: list, bu: "com.toto.test", platform: pbx.Application, expected: true, res: &pp0},
		{pps: list, bu: "com.toto.test.with.long.bundle.identifier", platform: pbx.Application, expected: true, res: &pp0},
		{pps: list, bu: "com.tutu.test2", platform: pbx.Application, expected: true, res: &pp1},
		{pps: list, bu: "com.tutu.test2", platform: pbx.TvExtension, expected: false},
		{pps: list, bu: "com.tata.ios", platform: pbx.Application, expected: true, res: &pp6},
		{pps: list, bu: "com.tata.tvos", platform: pbx.TvExtension, expected: false},
		{pps: list, bu: "com.tata.tutu.tvOS", platform: pbx.Application, expected: true, res: &pp2},
		{pps: list, bu: "com.tata.tutu.iOS", platform: pbx.Application, expected: true, res: &pp3},
		{pps: list, bu: "com.tata.tutu.tvOS", platform: pbx.TvExtension, expected: true, res: &pp2},
		{pps: list, bu: "com.tata.tutu.iOS", platform: pbx.TvExtension, expected: false},
		{pps: list, bu: "com.tete.tutu.iOS", platform: pbx.TvExtension, expected: false},
		{pps: list, bu: "com.tete.tutu.iOS", platform: pbx.Application, expected: true, res: &pp4},
		{pps: list, bu: "com.txtx.tutu.iOS", platform: pbx.TvExtension, expected: false},
		{pps: list, bu: "com.txtx.iOS", platform: pbx.TvExtension, expected: false},
	}

	for _, tt := range cases {
		s.Run(tt.bu, func() {
			// when:
			found, res := s.subject.findFor(tt.pps, tt.bu, tt.platform)

			// then:
			s.Assert().Equal(tt.expected, found)
			s.Assert().Equal(tt.res, res)
		})
	}
}

func (s *SignatureResolverSuite) TestContains() {
	cases := []struct {
		a  []string
		pt pbx.PBXProductType
		c  bool
	}{
		{a: []string{"tvOS", "iOS"}, pt: pbx.AppExtension, c: false},
		{a: []string{"tvOS", "iOS"}, pt: pbx.Application, c: true},
		{a: []string{"iOS"}, pt: pbx.TvExtension, c: false},
		{a: []string{"tvOS"}, pt: pbx.TvExtension, c: true},
	}

	for _, c := range cases {
		s.Run(fmt.Sprintf("%v-%v", c.a, c.pt), func() {
			// when:
			res := contains(c.a, c.pt)

			// then:
			s.Assert().EqualValues(c.c, res)
		})
	}
}
