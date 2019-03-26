package cms

import (
	"testing"
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/graphql"
	"github.com/stretchr/testify/require"
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/client"
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/resource"
	"fmt"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/domain/shared/wait"
	"github.com/stretchr/testify/assert"
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/domain/shared"
	"github.com/kyma-project/kyma/tests/console-backend-service"
)

const (
	clusterDocsTopicName1 = "ExampleClusterDocsTopic1"
	clusterDocsTopicName2 = "ExampleClusterDocsTopic2"
	clusterDocsTopicName3 = "ExampleClusterDocsTopic3"
)

type ClusterDocsTopicEvent struct {
	Type             string
	ClusterDocsTopic shared.ClusterDocsTopic
}

func TestClusterDocsTopicsQueries(t *testing.T) {
	c, err := graphql.New()
	require.NoError(t, err)

	cmsCli, _, err := client.NewCmsClientWithConfig()
	require.NoError(t, err)

	subscription := subscribeClusterDocsTopic(c, clusterDocsTopicEventDetailsFields())
	defer subscription.Close()

	clusterDocsTopicClient := resource.NewClusterDocsTopic(cmsCli, t.Logf)

	createClusterDocsTopic(t, clusterDocsTopicClient, clusterDocsTopicName1, "1")

	t.Log(fmt.Sprintf("Check subscription event of clusterDocsTopic %s created", clusterDocsTopicName1))
	expectedEvent := clusterDocsTopicEvent("ADD", fixClusterDocsTopic(clusterDocsTopicName1))
	event, err := readClusterDocsTopicEvent(subscription)
	assert.NoError(t, err)
	checkClusterDocsTopicEvent(t, expectedEvent, event)

	createClusterDocsTopic(t, clusterDocsTopicClient, clusterDocsTopicName2, "2")
	createClusterDocsTopic(t, clusterDocsTopicClient, clusterDocsTopicName3, "3")

	waitForClusterDocsTopic(t, clusterDocsTopicClient, clusterDocsTopicName1)

	t.Log(fmt.Sprintf("Check subscription event of clusterDocsTopic %s updated", clusterDocsTopicName1))
	expectedEvent = clusterDocsTopicEvent("UPDATE", fixClusterDocsTopic(clusterDocsTopicName1))
	event, err = readClusterDocsTopicEvent(subscription)
	assert.NoError(t, err)
	checkClusterDocsTopicEvent(t, expectedEvent, event)

	waitForClusterDocsTopic(t, clusterDocsTopicClient, clusterDocsTopicName2)
	waitForClusterDocsTopic(t, clusterDocsTopicClient, clusterDocsTopicName3)
}

func createClusterDocsTopic(t *testing.T, client *resource.ClusterDocsTopic, name, order string) {
	t.Log(fmt.Sprintf("Create clusterDocsTopic %s", name))
	err := client.Create(fixClusterDocsTopicMeta(name, order), fixCommonClusterDocsTopicSpec())
	require.NoError(t, err)
}

func waitForClusterDocsTopic(t *testing.T, client *resource.ClusterDocsTopic, name string) {
	t.Log(fmt.Sprintf("Wait for clusterDocsTopic Ready %s", name))
	err := wait.ForClusterDocsTopicReady(name, client.Get)
	require.NoError(t, err)
}

func subscribeClusterDocsTopic(c *graphql.Client, resourceDetailsQuery string) *graphql.Subscription {
	query := fmt.Sprintf(`
		subscription {
			clusterDocsTopicEvent {
				%s
			}
		}
	`, resourceDetailsQuery)
	req := graphql.NewRequest(query)

	return c.Subscribe(req)
}

func clusterDocsTopicDetailsFields() string {
	return `
		name
    	groupName
    	assets {
			name
			type
			files {
				url
				metadata
			}
		}
    	displayName
    	description
	`
}

func clusterDocsTopicEventDetailsFields() string {
	return fmt.Sprintf(`
        type
        clusterDocsTopic {
			%s
        }
    `, clusterDocsTopicDetailsFields())
}

func clusterDocsTopicEvent(eventType string, clusterDocsTopic shared.ClusterDocsTopic) ClusterDocsTopicEvent {
	return ClusterDocsTopicEvent{
		Type:             eventType,
		ClusterDocsTopic: clusterDocsTopic,
	}
}

func readClusterDocsTopicEvent(sub *graphql.Subscription) (ClusterDocsTopicEvent, error) {
	type Response struct {
		ClusterDocsTopicEvent ClusterDocsTopicEvent
	}

	var clusterDocsTopicEvent Response
	err := sub.Next(&clusterDocsTopicEvent, tester.DefaultSubscriptionTimeout)

	return clusterDocsTopicEvent.ClusterDocsTopicEvent, err
}

func checkClusterDocsTopicEvent(t *testing.T, expected, actual ClusterDocsTopicEvent) {
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.ClusterDocsTopic.Name, actual.ClusterDocsTopic.Name)
}

func fixClusterDocsTopic(name string) shared.ClusterDocsTopic {
	return shared.ClusterDocsTopic{
		Name: name,
	}
}

func fixClusterDocsTopicMeta(name, order string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: name,
		Labels: map[string]string{
			OrderLabel: order,
		},
	}
}

func fixCommonClusterDocsTopicSpec() v1alpha1.CommonDocsTopicSpec {
	return v1alpha1.CommonDocsTopicSpec{
		DisplayName: "Docs Topic Sample",
		Description: "Docs Topic Description",
		Sources: map[string]v1alpha1.Source{
			"openapi": {
				Mode: v1alpha1.DocsTopicSingle,
				URL:  "https://petstore.swagger.io/v2/swagger.json",
			},
		},
	}
}
