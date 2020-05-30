package pbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var subject pbxConvertor

func TestVariableResolution(t *testing.T) {
	// setup:
	m := map[string]string{
		"key1":  "value1",
		"key2":  "VALUE2",
		"key3":  "$(key2)",
		"entry": "$(key1:upper):$(key2:lower):$(key2:lower)",
	}

	// when:
	res := subject.Replace(m["entry"], "entry", m)

	// then:
	assert.EqualValues(t, "VALUE1:value2:value2", res)
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
