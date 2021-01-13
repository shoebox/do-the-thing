package action

import (
	"context"
	"dothething/internal/api"
	"dothething/internal/xcode"
	"fmt"

	"github.com/fatih/color"
)

//  # if ProvisionedDevices: !nil & "get-task-allow": true -> development
//  # if ProvisionedDevices: !nil & "get-task-allow": false -> ad-hoc
//  # if ProvisionedDevices: nil & "ProvisionsAllDevices": "true" -> enterprise
//  # if ProvisionedDevices: nil & ProvisionsAllDevices: nil -> app-store
type actionPackage struct {
	*api.API
}

func NewActionPackage(api *api.API) api.Action {
	return actionPackage{API: api}
}

func (a actionPackage) Run(ctx context.Context) error {
	xce := xcode.ParseXCodeBuildError(a.pack(ctx))
	if xce != nil {
		color.New(color.FgHiRed, color.Bold).Println(xce.Error())
	}

	return xce
}

func (a actionPackage) pack(ctx context.Context) error {
	// defer deletion of the keychain
	defer a.API.KeyChain.Delete(ctx)

	// Resolving signature
	if err := a.API.SignatureService.Run(ctx); err != nil {
		fmt.Println("err", err)
		return err
	}

	m, err := a.API.ExportOptionService.Compute()
	fmt.Println("method :::", m, err)
	/*
	 */

	return nil
}
