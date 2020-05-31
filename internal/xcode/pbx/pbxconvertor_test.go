package pbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var subject pbxConvertor

func TestVariableResolution(t *testing.T) {
	m := map[string]string{
		"BUNDLE_ID_BASE_NESTED":   "$(BUNDLE_ID_SUFFIX)",
		"BUNDLE_ID_BASE":          "toto",
		"BUNDLE_ID_APPNAME":       ".appname",
		"BUNDLE_ID_SIGNATURE":     ".sign",
		"BUNDLE_ID_SUFFIX":        ".suffix",
		"BUNDLE_ID_SUFFIX_NESTED": "$(BUNDLE_ID_SUFFIX)",

		"PRODUCT_BUNDLE_IDENTIFIER":  "$(BUNDLE_ID_BASE)$(BUNDLE_ID_APPNAME)$(BUNDLE_ID_SIGNATURE)$(BUNDLE_ID_SUFFIX_NESTED)",
		"PRODUCT_BUNDLE_IDENTIFIER2": "$(BUNDLE_ID_BASE:lower)$(BUNDLE_ID_APPNAME)$(BUNDLE_ID_SIGNATURE)$(BUNDLE_ID_SUFFIX_NESTED)",
	}

	// when:
	subject.Replace("PRODUCT_BUNDLE_IDENTIFIER", m)

	// then:
	assert.EqualValues(t, m["PRODUCT_BUNDLE_IDENTIFIER"], "toto.appname.sign.suffix")
}

func TestEntryToXCBuildConfguration(t *testing.T) {
	// setup:
	e := Entry{
		BaseConfigurationReference: "ref",
		BuildSettings: map[string]interface{}{
			"key1": 100,
			"key2": "value2",
		},
		Name: "entryName",
	}

	// when:
	res := subject.ToXCBuildConfiguration(e)

	// then:
	assert.EqualValues(t, XCBuildConfiguration{
		BuildSettings: map[string]string{
			"key1": "100",
			"key2": "value2",
		},
		BaseConfigurationReference: "ref",
		Name:                       "entryName",
	}, res)
}

func TestToNativeTarget(t *testing.T) {
	// setup:
	subject = pbxConvertor{
		PBXProjRaw{
			Objects: map[string]Entry{"value1": Entry{Name: "phaseName"}},
		}}

	e := Entry{
		BuildConfigurationList: Ref("toto"),
		BuildPhases:            []Ref{Ref("value1")},
		Name:                   "name",
		ProductName:            "productName",
		ProductInstallPath:     "productInstallPath",
		ProductType:            Framework,
	}

	// when:
	r := subject.ToNativeTarget(e)

	// then:
	assert.EqualValues(t,
		NativeTarget{
			BuildConfigurationList: XCConfigurationList{
				Reference:                   "",
				BuildConfiguration:          []XCBuildConfiguration{},
				DefaultConfigurationVisible: 0,
				DefaultConfigurationName:    "",
			},
			BuildPhases:        []PBXBuildPhase{{Name: "phaseName"}},
			Name:               "name",
			ProductName:        "productName",
			ProductInstallPath: "productInstallPath",
			ProductType:        Framework,
		}, r)
}
