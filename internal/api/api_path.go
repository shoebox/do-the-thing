package api

type PathService interface {
	Archive() string
	KeyChain() string
	ObjRoot() string
	SymRoot() string
	XCResult() string
}
