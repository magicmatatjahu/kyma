package spec_test

import (
	"encoding/json"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/rafter/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAsyncAPISpec_Decode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// given
		asyncAPIBytesData, err := fixAsyncAPIBytesData()
		require.NoError(t, err)
		assert.NotNil(t, asyncAPIBytesData)

		// expected
		asyncAPIData := fixAsyncAPIData()

		// when
		asyncAPI := spec.AsyncAPISpec{}
		err = asyncAPI.Decode(asyncAPIBytesData)

		// then
		require.NoError(t, err)
		assert.Equal(t, asyncAPIData.Raw, asyncAPI.Raw)
		assert.Equal(t, asyncAPIData.Data, asyncAPI.Data)
	})
}

func fixAsyncAPIData() spec.AsyncAPISpec {
	return spec.AsyncAPISpec{
		Raw: map[string]interface{}{
			"asyncapi": "2.0.0",
			"info": map[string]interface{}{
				"title": "Not example",
				"version": "1.0.0",
			},
			"channels": map[string]interface{}{
				"streetlights": map[string]interface{}{
					"publish": map[string]interface{}{
						"summary": "Inform about environmental lighting conditions of a particular streetlight.",
						"operationId": "receiveLightMeasurement",
					},
				},
			},
		},
		Data: spec.AsyncAPISpecData{
			AsyncAPI: "2.0.0",
			Channels: map[string]interface{}{
				"streetlights": map[string]interface{}{
					"publish": map[string]interface{}{
						"summary": "Inform about environmental lighting conditions of a particular streetlight.",
						"operationId": "receiveLightMeasurement",
					},
				},
			},
		},
	}
}

func fixAsyncAPIBytesData() ([]byte, error) {
	bytes, err := json.Marshal(fixAsyncAPIData().Raw)

	if err != nil {
		return nil, err
	}
	return bytes, nil
}
