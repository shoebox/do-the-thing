package path

import (
	"dothething/internal/api"
	"fmt"
	"path/filepath"
)

const BuildFolder = "../Build/"

type pathService struct {
	api.API
}

func NewPathService(p api.API) api.PathService {
	return pathService{API: p}
}

func (p pathService) buildFolder() string {
	return filepath.Clean(filepath.Join(p.API.Config().Path, BuildFolder))
}

func (p pathService) Archive() string {
	return filepath.Clean(filepath.Join(
		p.buildFolder(),
		fmt.Sprintf("%v-%v.xcarchive", p.API.Config().Scheme, p.API.Config().Configuration),
	))
}

func (p pathService) KeyChain() string {
	return filepath.Join(p.buildFolder(), "do-the-thing.keychain")
}

func (p pathService) ObjRoot() string {
	return fmt.Sprintf("OBJROOT=%v", filepath.Join(p.buildFolder(), "obj/"))
}

func (p pathService) SymRoot() string {
	return fmt.Sprintf("SYMROOT=%v", filepath.Join(p.buildFolder(), "sym/"))
}

func (p pathService) XCResult() string {
	return filepath.Join(
		p.buildFolder(),
		fmt.Sprintf("%v-%v.xcresult", p.API.Config().Scheme, p.API.Config().Configuration),
	)
}
