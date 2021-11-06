package path

import (
	"dothething/internal/api"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	BuildFolder      = "../Build/"
	projectFileExt   = ".xcodeproj"
	workspaceFileExt = ".xcworkspace"
)

type pathService struct {
	*api.API
}

func NewPathService(p *api.API) api.PathService {
	return pathService{API: p}
}

func (p pathService) buildFolder() string {
	return filepath.Clean(filepath.Join(p.API.Config.Path, BuildFolder))
}

func (p pathService) Archive() string {
	return filepath.Clean(filepath.Join(
		p.buildFolder(),
		fmt.Sprintf("%v-%v-%v.xcarchive",
			p.API.Config.Target,
			p.API.Config.Scheme,
			p.API.Config.Configuration),
	))
}

func (p pathService) ExportPList() string {
	return filepath.Clean(filepath.Join(
		p.buildFolder(),
		fmt.Sprintf("%v-%v-%v-export.plist",
			p.API.Config.Target,
			p.API.Config.Scheme,
			p.API.Config.Configuration),
	))
}

func (p pathService) KeyChain() string {
	res, err := filepath.Abs(filepath.Join(p.buildFolder(), "do-the-thing.keychain"))
	if err != nil {
		fmt.Printf("Error %v", err)
	}

	return res
}

func (p pathService) Package() string {
	return filepath.Clean(filepath.Join(
		p.buildFolder(),
		fmt.Sprintf("%v-%v-%v.ipa",
			p.API.Config.Target,
			p.API.Config.Scheme,
			p.API.Config.Configuration),
	))
}

func (p pathService) ObjRoot() string {
	return fmt.Sprintf("OBJROOT=%v", filepath.Join(p.buildFolder(), "obj/"))
}

func (p pathService) SymRoot() string {
	return fmt.Sprintf("SYMROOT=%v", filepath.Join(p.buildFolder(), "sym/"))
}

func (p pathService) DerivedData() string {
	return filepath.Join(p.buildFolder(), "derived-data/")
}

func (p pathService) XCResult() string {
	return filepath.Join(
		p.buildFolder(),
		fmt.Sprintf("%v-%v-%v.xcresult",
			p.API.Config.Target,
			p.API.Config.Scheme,
			p.API.Config.Configuration,
		),
	)
}

func (p pathService) XCodeProject() string {
	// TODO: Should we use the xcworkspacedata to resolve the project path?
	if filepath.Ext(p.Config.Path) == workspaceFileExt {
		return strings.TrimSuffix(p.Config.Path, workspaceFileExt) + projectFileExt
	}

	return p.Config.Path
}

func (p pathService) PBXProj() string {
	return fmt.Sprintf("%v/project.pbxproj", p.XCodeProject())
}
