package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
	"github.com/golang/glog"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
	"github.com/pkg/errors"
)

type assetResolver struct {
	assetSvc assetSvc
	assetConverter gqlAssetConverter
	fileSvc *fileService
	fileConverter *fileConverter
}

func newAssetResolver(assetService *assetService) *assetResolver {
	return &assetResolver{
		assetSvc: assetService,
		assetConverter: &assetConverter{},
		fileSvc: &fileService{},
		fileConverter: &fileConverter{},
	}
}

func (r *assetResolver) AssetFilesField(ctx context.Context, obj *gqlschema.Asset, filterExtensions []string) ([]gqlschema.File, error) {
	if obj == nil {
		glog.Error(errors.New("%s cannot be empty in order to resolve `files` field"), pretty.Asset)
		return nil, gqlerror.NewInternal()
	}

	asset, err := r.assetSvc.Find(obj.Namespace, obj.Name)
	if err != nil {
		if module.IsDisabledModuleError(err) {
			return nil, err
		}
		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", pretty.Asset, pretty.Asset, obj.Name))
		return nil, gqlerror.New(err, pretty.Asset)
	}

	items, err := r.fileSvc.FilterByExtensions(&asset.Status.AssetRef, filterExtensions)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", pretty.Files, pretty.Asset, obj.Name))
		return nil, gqlerror.New(err, pretty.Asset)
	}

	files, err := r.fileConverter.ToGQLs(items)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while converting %s", pretty.Files))
		return nil, gqlerror.New(err, pretty.Asset)
	}

	return files, nil
}
