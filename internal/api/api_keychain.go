package api

import "context"

type KeyChain interface {
	Create(ctx context.Context, password string) error
	Delete(ctx context.Context) error
	ImportCertificate(ctx context.Context, filePathh string, password string, identity string) error
	GetPath() string
}
