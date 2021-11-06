package api

type PathService interface {
	Archive() string
	DerivedData() string
	ExportPList() string
	KeyChain() string
	ObjRoot() string
	PBXProj() string
	Package() string
	SymRoot() string
	XCResult() string
	XCodeProject() string
}
