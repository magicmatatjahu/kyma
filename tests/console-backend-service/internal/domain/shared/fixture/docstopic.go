package fixture

import "github.com/kyma-project/kyma/tests/console-backend-service/internal/domain/shared"

const (
	DocsTopicGroupName   = "ExampleGroupName"
	DocsTopicDisplayName = "Docs Topic Sample"
	DocsTopicDescription = "Docs Topic Description"
)

func DocsTopic(namespace, name string) shared.DocsTopic {
	return shared.DocsTopic{
		Name:        name,
		Namespace:   namespace,
		GroupName:   DocsTopicGroupName,
		DisplayName: DocsTopicDisplayName,
		Description: DocsTopicDescription,
	}
}
