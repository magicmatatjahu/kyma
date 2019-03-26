package status

import (
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
)

type DocsTopicExtractor struct{}

func (e *DocsTopicExtractor) Status(status v1alpha1.CommonDocsTopicStatus) gqlschema.DocsTopicStatus {
	return gqlschema.DocsTopicStatus{
		Phase:   e.Phase(status.Phase),
		Reason:  status.Reason,
		Message: status.Message,
	}
}

func (e *DocsTopicExtractor) Phase(phase v1alpha1.DocsTopicPhase) gqlschema.DocsTopicPhaseType {
	switch phase {
	case v1alpha1.DocsTopicReady:
		return gqlschema.DocsTopicPhaseTypeReady
	case v1alpha1.DocsTopicPending:
		return gqlschema.DocsTopicPhaseTypePending
	default:
		return gqlschema.DocsTopicPhaseTypeFailed
	}
}
