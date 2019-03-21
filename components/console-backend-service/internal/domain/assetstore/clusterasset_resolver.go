package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
	"github.com/pkg/errors"
	"github.com/golang/glog"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
)

type clusterAssetResolver struct {
	clusterAssetSvc clusterAssetSvc
	clusterAssetConverter gqlClusterAssetConverter
	fileSvc *fileService
	fileConverter *fileConverter
}

func newClusterAssetResolver(clusterAssetService *clusterAssetService) *clusterAssetResolver {
	return &clusterAssetResolver{
		clusterAssetSvc: clusterAssetService,
		clusterAssetConverter: &clusterAssetConverter{},
		fileSvc: &fileService{},
		fileConverter: &fileConverter{},
	}
}

func (r *clusterAssetResolver) ClusterAssetFilesField(ctx context.Context, obj *gqlschema.ClusterAsset, filterExtensions []string) ([]gqlschema.File, error) {
	if obj == nil {
		glog.Error(errors.New("%s cannot be empty in order to resolve `files` field"), pretty.ClusterAsset)
		return nil, gqlerror.NewInternal()
	}

	asset, err := r.clusterAssetSvc.Find(obj.Name)
	if err != nil {
		if module.IsDisabledModuleError(err) {
			return nil, err
		}
		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", pretty.ClusterAsset, pretty.ClusterAsset, obj.Name))
		return nil, gqlerror.New(err, pretty.ClusterAsset)
	}

	items, err := r.fileSvc.FilterByExtensions(&asset.Status.AssetRef, filterExtensions)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", pretty.Files, pretty.ClusterAsset, obj.Name))
		return nil, gqlerror.New(err, pretty.ClusterAsset)
	}

	files, err := r.fileConverter.ToGQLs(items)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while converting %s", pretty.Files))
		return nil, gqlerror.New(err, pretty.ClusterAsset)
	}

	return files, nil
}
