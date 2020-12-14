package destination

import "fmt"

func NewActionError(action string, destination string) error {
	return fmt.Errorf("Failed to %s destination: %v", action, destination)
}

func NewShutDownError(destination string) error {
	return NewActionError("Shutdown", destination)
}

func NewBootError(destination string) error {
	return NewActionError("Boot", destination)
}

func NewListingError() error {
	return fmt.Errorf("Failed to list destinations")
}
