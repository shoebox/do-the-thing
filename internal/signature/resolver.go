package signature

import (
	"context"
	"dothething/internal/xcode"
	"dothething/internal/xcode/pbx"
	"dothething/internal/xcode/project"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type SignatureResolver interface {
	Resolve(ctx context.Context, c xcode.Config, p project.Project)
}

func NewSignatureResolver(p ProvisioningService) SignatureResolver {
	return signatureResolver{p}
}

type signatureResolver struct {
	p ProvisioningService
}

func isProvisioningFile(info os.FileInfo) bool {
	return info.Mode().IsRegular() &&
		!info.IsDir() &&
		strings.HasSuffix(info.Name(), ".mobileprovision")
}

func (r signatureResolver) Resolve(ctx context.Context, c xcode.Config, p project.Project) {
	// TODO: Temp test path
	//pps := r.resolveProvisioningFilesInFolder(ctx, c.CodeSignOption.Path)

	var nativeTarget pbx.NativeTarget
	for _, tgt := range p.Pbx.Targets {
		if tgt.Name == c.Target {
			nativeTarget = tgt
			break
		}
	}

	var bc pbx.XCBuildConfiguration
	for _, b := range nativeTarget.BuildConfigurationList.BuildConfiguration {
		if b.Name == c.Configuration {
			bc = b
			break
		}
	}
	team := bc.BuildSettings["DEVELOPMENT_TEAM"]
	pbi := bc.BuildSettings["PRODUCT_BUNDLE_IDENTIFIER"]
	fmt.Println(team, pbi)

	//for _, pp := range pps {
	//	fmt.Println(pp.Entitlements.AppID)
	//}
}

func (r signatureResolver) resolveProvisioningFilesInFolder(ctx context.Context, root string) []ProvisioningProfile {
	g, ctx := errgroup.WithContext(ctx)
	paths := make(chan string)
	cp := make(chan ProvisioningProfile)
	g.Go(func() error {
		defer close(paths)
		return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err == nil && isProvisioningFile(info) {
				p := path
				g.Go(func() error {
					dpp, err := r.p.Decode(ctx, p)
					if err != nil {
						return err
					}

					select {
					case cp <- dpp:
					case <-ctx.Done():
						return ctx.Err()
					}
					return nil
				})
			}
			return nil
		})
	})

	go func() {
		g.Wait()
		close(cp)
	}()

	res := []ProvisioningProfile{}
	for pp := range cp {
		res = append(res, pp)
	}

	return res
}

func parseProvisioning(r io.ByteReader) {
}
