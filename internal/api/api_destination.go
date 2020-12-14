package api

import "context"

// Destination available destination for the scheme
type Destination struct {
	Name     string
	Platform string
	ID       string
	OS       string
}

// DestinationService destination service definition
type DestinationService interface {
	Boot(ctx context.Context, d Destination) error
	List(ctx context.Context, scheme string) ([]Destination, error)
	ShutDown(ctx context.Context, d Destination) error
}
