package pbx

func (prj PBXProjRaw) Parse() PBXProject {
	c := NewConvertor(prj)

	var tgs []NativeTarget

	for _, tgt := range prj.GetRoot().Targets.GetList(prj) {
		tgs = append(tgs, c.ToNativeTarget(tgt))
	}

	return PBXProject{Targets: tgs}
}
