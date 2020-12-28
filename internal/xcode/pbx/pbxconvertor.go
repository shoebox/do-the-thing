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
		Dependencies:           c.ToDependencies(e.Dependencies),
		Ref:                    e.Ref,
		Name:                   e.Name,
		ProductName:            e.ProductName,
		ProductInstallPath:     e.ProductInstallPath,
		ProductType:            PBXProductType(e.ProductType),
	}
}

func (c pbxConvertor) ToDependencies(a ArrayRef) []NativeTarget {
	var res []NativeTarget
	for _, e := range a.GetList(c.p) {
		res = append(res, c.ToNativeTarget(e.Target.Get(c.p)))
	}
	return res
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
		Reference:                  e.Ref,
		BuildSettings:              c.ToStringMap(e.BuildSettings),
		BaseConfigurationReference: e.BaseConfigurationReference,
		Name:                       e.Name,
	}
}

func (c pbxConvertor) ToStringMap(m map[string]interface{}) map[string]string {
	res := map[string]string{}
	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	for k, _ := range m {
		c.Replace(k, res)
	}

	return res
}
func (c pbxConvertor) Replace(key string, m map[string]string) {
	r := regexp.MustCompile(`\$\(([^\)\:]+):?(lower|upper)?\)`)

	for _, e := range r.FindAllStringSubmatch(m[key], -1) {
		k := e[1]
		mod := e[2]
		c.Replace(k, m)

		val := m[key]
		if mod == "lower" {
			val = strings.ToLower(val)
		} else if mod == "upper" {
			val = strings.ToUpper(val)
		}

		m[key] = strings.ReplaceAll(m[key], e[0], m[k])
	}
}
