package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
)

type fileConverter struct {}

func (c *fileConverter) ToGQL(file *File) (*gqlschema.File, error) {
	if file == nil {
		return nil, nil
	}

	metadata := c.extractMetadata(file.Metadata)

	result := gqlschema.File{
		URL: file.URL,
		Metadata: metadata,
	}
	return &result, nil
}

func (c *fileConverter) ToGQLs(files []*File) ([]gqlschema.File, error) {
	var result []gqlschema.File
	for _, u := range files {
		converted, err := c.ToGQL(u)
		if err != nil {
			return nil, err
		}

		if converted != nil {
			result = append(result, *converted)
		}
	}
	return result, nil
}

func (c *fileConverter) extractMetadata(metadata map[string]interface{}) gqlschema.JSON {
	if metadata == nil {
		return nil
	}

	result := make(gqlschema.JSON)
	result = metadata

	return result
}