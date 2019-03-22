package cms_test

import (
	"time"
	"testing"
	"k8s.io/client-go/rest"
	"github.com/stretchr/testify/require"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms"
	"github.com/stretchr/testify/assert"
	"context"
)

const testTimes = 3
const informerResyncPeriod = 10 * time.Second

func TestPluggableContainer(t *testing.T) {
	pluggable, err := cms.New(&rest.Config{}, informerResyncPeriod, nil)
	require.NoError(t, err)

	pluggable.SetFakeClient()

	for i := 0; i < testTimes; i++ {
		require.NotPanics(t, func() {
			err := pluggable.Enable()
			require.NoError(t, err)
			<-pluggable.Pluggable.SyncCh

			checkExportedFields(t, pluggable, true)
		})
		require.NotPanics(t, func() {
			err := pluggable.Disable()
			require.NoError(t, err)

			checkExportedFields(t, pluggable, false)
		})
	}
}

func checkExportedFields(t *testing.T, resolver *cms.PluggableContainer, enabled bool) {
	assert.NotNil(t, resolver.Resolver)
	require.NotNil(t, resolver.CmsRetriever)
	assert.NotNil(t, resolver.CmsRetriever.ClusterDocsTopicGetter)
	assert.NotNil(t, resolver.CmsRetriever.DocsTopicGetter)
	assert.NotNil(t, resolver.CmsRetriever.GqlClusterDocsTopicConverter)
	assert.NotNil(t, resolver.CmsRetriever.GqlDocsTopicConverter)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	val, err := resolver.Resolver.ClusterDocsTopicsQuery(ctx, nil, nil)
	if enabled {
		require.NoError(t, err)
	} else {
		require.Error(t, err)
		require.Nil(t, val)
	}
}
