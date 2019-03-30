package assetstore_test

import (
	"testing"

	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileConverter_ToGQL(t *testing.T) {
	t.Run("Success without filter", func(t *testing.T) {
		assetRef := &v1alpha2.AssetStatusRef{
			BaseUrl: "https://example.com",
			Assets: []string{
				"markdown.md",
				"apiSpec.json",
				"odata.xml",
			},
		}
		expected := []*assetstore.File{
			{
				URL:      "https://example.com/markdown.md",
				Metadata: map[string]interface{}{},
			},
			{
				URL:      "https://example.com/apiSpec.json",
				Metadata: map[string]interface{}{},
			},
			{
				URL:      "https://example.com/odata.xml",
				Metadata: map[string]interface{}{},
			},
		}

		svc := assetstore.NewFileService()

		result, err := svc.Extract(assetRef)
		require.NoError(t, err)

		assert.Equal(t, expected, result)
	})

	t.Run("Success with filter", func(t *testing.T) {
		assetRef := &v1alpha2.AssetStatusRef{
			BaseUrl: "https://example.com",
			Assets: []string{
				"markdown.md",
				"apiSpec.json",
				"odata.xml",
			},
		}
		expected := []*assetstore.File{
			{
				URL:      "https://example.com/markdown.md",
				Metadata: map[string]interface{}{},
			},
		}

		svc := assetstore.NewFileService()

		result, err := svc.FilterByExtensionsAndExtract(assetRef, []string{"md"})
		require.NoError(t, err)

		assert.Equal(t, expected, result)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := assetstore.NewFileService()

		result, err := svc.Extract(nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}
