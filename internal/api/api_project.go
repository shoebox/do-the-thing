package api

import (
	"context"
	"dothething/internal/xcode/pbx"
)

type ProjectService interface {
	Parse(ctx context.Context) (Project, error)
}

// Project datas
type Project struct {
	Configurations []string `json:"configurations"`
	Name           string   `json:"name"`
	Pbx            pbx.PBXProject
	Schemes        []string `json:"schemes"`
	Targets        []string `json:"targets"`
}
