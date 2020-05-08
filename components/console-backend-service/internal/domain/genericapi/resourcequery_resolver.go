package genericapi

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
)

type resourceQueryResolver struct {
	services resourcesServices
	converter *resourceConverter
}

func newResourceQueryResolver(services resourcesServices, converter *resourceConverter) *resourceQueryResolver {
	return &resourceQueryResolver{
		services: services,
		converter: converter,
	}
}

func (r *resourceQueryResolver) Get(ctx context.Context, schema gqlschema.SchemaResourceInput, name string, namespace *string) (*gqlschema.Resource, error) {
	service := r.services.retrieveService(schema)
	if service == nil {
		return nil, nil
	}

	key := name
	if namespace != nil {
		key = fmt.Sprintf("%s/%s", *namespace, name)
	}

	item, exists, err := service.Informer.GetIndexer().GetByKey(key)
	if err != nil || !exists {
		return nil, err
	}

	return r.converter.ToGQL(item)
}

func (r *resourceQueryResolver) List(ctx context.Context, schema gqlschema.SchemaResourceInput, namespace *string) (gqlschema.ResourceListOutput, error) {
	service := r.services.retrieveService(schema)
	if service == nil {
		return gqlschema.ResourceListOutput{}, nil
	}

	var items []interface{}
	var err error
	if namespace != nil {
		items, err = service.Informer.GetIndexer().ByIndex("namespace", *namespace)
	} else {
		items = service.Informer.GetStore().List()
	}

	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	return r.converter.ToGQLs(items)
}

func (r *resourceQueryResolver) ResourceFields(ctx context.Context, obj *gqlschema.Resource, fields []gqlschema.ResourceFieldInput) (gqlschema.JSON, error) {
	gqlJson := gqlschema.JSON{}

	for _, field := range fields {
		pathField := strings.Split(field.Path, ".")

		val, found, err := unstructured.NestedFieldCopy(obj.RawContent, pathField...)
		if err != nil {
			return nil, err
		}
		if !found {
			continue
		}

		gqlJson[field.Key] = val
	}

	return gqlJson, nil
}

func (r *resourceQueryResolver) ResourceSubResources(ctx context.Context, obj *gqlschema.Resource, resources []gqlschema.SubResourceInput) ([]gqlschema.SubResourceOutput, error) {
	return nil, nil
}
