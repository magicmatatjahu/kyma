package status

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

type AssetExtractor struct{}

func (e *AssetExtractor) Status(status v1alpha2.CommonAssetStatus) gqlschema.AssetStatus {
	 return gqlschema.AssetStatus{
	 	Phase: e.Phase(status.Phase),
		Reason: status.Reason,
		Message: status.Message,
	}
}

func (e *AssetExtractor) Phase(phase v1alpha2.AssetPhase) gqlschema.AssetPhaseType {
	switch phase {
	case v1alpha2.AssetReady:
		return gqlschema.AssetPhaseTypeReady
	case v1alpha2.AssetPending:
		return gqlschema.AssetPhaseTypePending
	default:
		return gqlschema.AssetPhaseTypeFailed
	}
}