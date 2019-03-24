package cms_test

import (
	"testing"
	"time"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/automock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"context"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	assetstoreMock "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared/automock"
	"github.com/kyma-project/kyma/components/assetstore-controller-manager/pkg/apis/assetstore/v1alpha2"
)

func TestClusterDocsTopicResolver_ClusterDocsTopicsQuery(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		resource :=
			&v1alpha1.ClusterDocsTopic{
				ObjectMeta: v1.ObjectMeta{
					Name: "test",
				},
			}
		resources := []*v1alpha1.ClusterDocsTopic{
			resource, resource,
		}
		expected := []gqlschema.ClusterDocsTopic{
			{
				Name: "Test",
			}, {
				Name: "Test",
			},
		}

		svc := automock.NewClusterDocsTopicService()
		svc.On("List", (*string)(nil), (*string)(nil)).Return(resources, nil).Once()
		defer svc.AssertExpectations(t)

		converter := automock.NewGQLClusterDocsTopicConverter()
		converter.On("ToGQLs", resources).Return(expected, nil)
		defer converter.AssertExpectations(t)

		resolver := cms.NewClusterDocsTopicResolver(svc, nil)
		resolver.SetDocsTopicConverter(converter)

		result, err := resolver.ClusterDocsTopicsQuery(nil, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		var resources []*v1alpha1.ClusterDocsTopic

		svc := automock.NewClusterDocsTopicService()
		svc.On("List", (*string)(nil), (*string)(nil)).Return(resources, nil).Once()
		defer svc.AssertExpectations(t)
		resolver := cms.NewClusterDocsTopicResolver(svc, nil)
		var expected []gqlschema.ClusterDocsTopic

		result, err := resolver.ClusterDocsTopicsQuery(nil, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		expected := errors.New("Test")

		var resources []*v1alpha1.ClusterDocsTopic

		svc := automock.NewClusterDocsTopicService()
		svc.On("List", (*string)(nil), (*string)(nil)).Return(resources, expected).Once()
		defer svc.AssertExpectations(t)
		resolver := cms.NewClusterDocsTopicResolver(svc, nil)

		_, err := resolver.ClusterDocsTopicsQuery(nil, nil, nil)

		require.Error(t, err)
		assert.True(t, gqlerror.IsInternal(err))
	})
}

func TestClusterDocsTopicResolver_ClusterDocsTopicAssetsField(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		name := "name"
		resources := []*v1alpha2.ClusterAsset{
			{

			},
		}
		expected := new(gqlschema.JSON)
		err := expected.UnmarshalGQL(resource.Raw)
		require.NoError(t, err)

		resourceGetter := new(assetstoreMock.ClusterAssetGetter)
		resourceGetter.On("ListForDocsTopicByType", "docs-topic", []string{}).Return(resources, nil).Once()
		defer resourceGetter.AssertExpectations(t)

		resourceConverter := new(assetstoreMock.GqlClusterAssetConverter)
		resourceConverter.On("ToGQLs", "docs-topic", []string{}).Return(resource, nil).Once()
		defer resourceGetter.AssertExpectations(t)

		retriever := new(assetstoreMock.AssetStoreRetriever)
		retriever.On("ClusterAsset").Return(resourceGetter)
		retriever.On("ClusterAssetConverter").Return(resourceConverter)

		parentObj := gqlschema.ClusterDocsTopic{
			Name: name,
		}

		resolver := cms.NewClusterDocsTopicResolver(nil, retriever)

		result, err := resolver.ClusterDocsTopicAssetsField(nil, &parentObj, []string{})

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestClusterDocsTopicResolver_ClusterDocsTopicEventSubscription(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), (-24 * time.Hour))
		cancel()

		svc := automock.NewClusterDocsTopicService()
		svc.On("Subscribe", mock.Anything).Once()
		svc.On("Unsubscribe", mock.Anything).Once()
		resolver := cms.NewClusterDocsTopicResolver(svc, nil)

		_, err := resolver.ClusterDocsTopicEventSubscription(ctx)

		require.NoError(t, err)
		svc.AssertCalled(t, "Subscribe", mock.Anything)
	})

	t.Run("Unsubscribe after connection close", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), (-24 * time.Hour))
		cancel()

		svc := automock.NewClusterDocsTopicService()
		svc.On("Subscribe", mock.Anything).Once()
		svc.On("Unsubscribe", mock.Anything).Once()
		resolver := cms.NewClusterDocsTopicResolver(svc, nil)

		channel, err := resolver.ClusterDocsTopicEventSubscription(ctx)
		<-channel

		require.NoError(t, err)
		svc.AssertCalled(t, "Unsubscribe", mock.Anything)
	})
}
