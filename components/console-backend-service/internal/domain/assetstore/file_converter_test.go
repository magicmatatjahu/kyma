package assetstore

import (
	"testing"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileConverter_ToGQL(t *testing.T) {
	t.Run("All properties are given", func(t *testing.T) {
		converter := fileConverter{}

		item := fixFile()
		expected := gqlschema.File{
			URL:      "ExampleUrl",
			Metadata: map[string]interface{}{},
		}

		result, err := converter.ToGQL(item)
		require.NoError(t, err)
		assert.Equal(t, &expected, result)
	})

	t.Run("Empty", func(t *testing.T) {
		converter := &fileConverter{}
		_, err := converter.ToGQL(&File{})
		require.NoError(t, err)
	})

	t.Run("Nil", func(t *testing.T) {
		converter := &fileConverter{}
		item, err := converter.ToGQL(nil)

		require.NoError(t, err)
		assert.Nil(t, item)
	})
}

func TestFileConverter_ToGQLs(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		files := []*File{
			fixFile(),
			fixFile(),
		}

		converter := fileConverter{}
		result, err := converter.ToGQLs(files)

		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "ExampleUrl", result[0].URL)
	})

	t.Run("Empty", func(t *testing.T) {
		var files []*File

		converter := fileConverter{}
		result, err := converter.ToGQLs(files)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("With nil", func(t *testing.T) {
		files := []*File{
			nil,
			fixFile(),
			nil,
		}

		converter := fileConverter{}
		result, err := converter.ToGQLs(files)

		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "ExampleUrl", result[0].URL)
	})
}

func fixFile() *File {
	return &File{
		URL:      "ExampleUrl",
		Metadata: map[string]interface{}{},
	}
}
