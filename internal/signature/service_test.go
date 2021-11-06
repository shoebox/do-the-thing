package signature

/*
type signatureServiceSuite struct {
	suite.Suite
	API     *api.API
	plist   *api.MockPlistAPI
	subject *signatureService
}

// suite testing entrypoint
func TestSignatureServiceSuite(t *testing.T) {
	suite.Run(t, new(signatureServiceSuite))
}

func (s *signatureServiceSuite) BeforeTest(suiteName, testName string) {
	s.plist = new(api.MockPlistAPI)
	s.API = &api.API{
		PlistBuddyService: s.plist,
	}

	s.subject = new(signatureService)
	s.subject.API = s.API
}

func (s *signatureServiceSuite) TestConfigureBuildSettings() {
	// setup:
	key1 := "key1"
	key2 := "key2"
	m := map[string]string{
		key1: "value1",
		key2: "value2",
	}
	ctx := context.Background()
	xcb := pbx.XCBuildConfiguration{
		Reference: "RefID",
		BuildSettings: map[string]string{
			key2: s.subject.buildSettingsPath("key2-value"),
		},
	}

	s.plist.On(
		"AddStringValue",
		ctx,
		xcb.Reference,
		s.subject.buildSettingsPath(key1),
		"value1",
	).Return(nil)

	s.plist.On(
		"SetStringValue",
		ctx,
		xcb.Reference,
		s.subject.buildSettingsPath(key2),
		"value2",
	).Return(nil)

	// when:
	err := s.subject.configureBuildSetting(ctx, xcb, m)

	// then:
	s.plist.AssertExpectations(s.T())

	// and:
	s.Assert().NoError(err)
}

func (s *signatureServiceSuite) TestShouldConfigureDependencies() {
	for _, pt := range pbx.ProductTypes {
		called := false

		s.subject.configureDependencies(pbx.NativeTarget{
			Dependencies: []pbx.NativeTarget{{ProductType: pt}},
		}, func(s string) error {
			called = true
			return nil
		})

		if pt == pbx.Application || pt == pbx.TvExtension {
			s.Assert().True(called)
		} else {
			s.Assert().False(called)
		}
	}
}
*/
