package api

import "context"

type PListBuddyService interface {
	AddStringValue(ctx context.Context, objectId string, path string, value string) error
	SetStringValue(ctx context.Context, objectId string, path string, value string) error
}
