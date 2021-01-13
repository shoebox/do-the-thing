package util

import (
	"context"
	"dothething/internal/api"
	"fmt"
)

const buddy = "/usr/libexec/PlistBuddy"

type plistBuddy struct {
	*api.API
}

func NewPListBuddy(api *api.API) plistBuddy {
	return plistBuddy{api}
}

func (p plistBuddy) computePath(action, objectId, path, value string) string {
	return fmt.Sprintf("%s :objects:%v:%v",
		action,
		objectId,
		path)
}

func (p plistBuddy) AddStringValue(ctx context.Context, objectId string, path string, value string) error {
	return p.API.Exec.CommandContext(ctx,
		buddy,
		"-c",
		fmt.Sprintf("%v string %v", p.computePath("Add", objectId, path, value), value),
		fmt.Sprintf("%v/project.pbxproj", p.Config.Path)).Run()
}

func (p plistBuddy) SetStringValue(ctx context.Context, objectId string, path string, value string) error {
	action := fmt.Sprintf("Set :objects:%v:%v %v",
		objectId,
		path,
		value)

	return p.API.Exec.CommandContext(ctx,
		buddy,
		"-c",
		action,
		fmt.Sprintf("%v/project.pbxproj", p.Config.Path)).Run()
}
