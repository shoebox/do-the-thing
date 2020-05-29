package pbx

import (
	"fmt"
	"regexp"
	"strings"
)

type PBXConvertor interface {
	ToNativeTarget(e Entry) NativeTarget
}

func NewConvertor(p PBXProjRaw) pbxConvertor {
	return pbxConvertor{p}
}

type pbxConvertor struct {
	p PBXProjRaw
}

func (c pbxConvertor) ToNativeTarget(e Entry) NativeTarget {
	return NativeTarget{
		BuildConfigurationList: c.ToXCConfigurationList(e.BuildConfigurationList.Get(c.p)),
		BuildPhases:            c.ToBuildPhases(e.BuildPhases),
		Name:                   e.Name,
		ProductName:            e.ProductName,
		ProductInstallPath:     e.ProductInstallPath,
		ProductType:            PBXProductType(e.ProductType),
	}
}

func (c pbxConvertor) ToBuildPhases(a ArrayRef) []PBXBuildPhase {
	res := []PBXBuildPhase{}
	for _, p := range a.GetList(c.p) {
		res = append(res, PBXBuildPhase{Name: p.Name})
	}
	return res
}

func (c pbxConvertor) ToXCConfigurationList(e Entry) XCConfigurationList {
	return XCConfigurationList{
		BuildConfiguration:       c.ToXCConfigurationArray(e.BuildConfigurations),
		Reference:                e.Ref,
		DefaultConfigurationName: e.DefaultConfigurationName,
	}
}

func (c pbxConvertor) ToXCConfigurationArray(e ArrayRef) []XCBuildConfiguration {
	res := []XCBuildConfiguration{}
	for _, cfg := range e.GetList(c.p) {
		res = append(res, c.ToXCBuildConfiguration(cfg))
	}

	return res
}

func (c pbxConvertor) ToXCBuildConfiguration(e Entry) XCBuildConfiguration {
	return XCBuildConfiguration{
		BuildSettings:              c.ToStringMap(e.BuildSettings),
		BaseConfigurationReference: e.BaseConfigurationReference,
		Name:                       e.Name,
	}
}

func (c pbxConvertor) ToStringMap(m map[string]interface{}) map[string]string {
	res := map[string]string{}
	for k, v := range m {
		res[k] = c.Replace(fmt.Sprintf("%v", v), k, res)
	}

	return res
}

func (c pbxConvertor) Replace(into string, key string, m map[string]string) string {
	r := regexp.MustCompile(`\$\(([^\)\:]+):?(lower|upper)?\)`)

	for _, e := range r.FindAllStringSubmatch(into, -1) {
		k := e[1]
		mod := e[2]

		val := c.Replace(m[k], k, m)

		if mod == "lower" {
			val = strings.ToLower(val)
		} else if mod == "upper" {
			val = strings.ToUpper(val)
		}

		into = strings.ReplaceAll(into, e[0], val)
	}
	return into
}
