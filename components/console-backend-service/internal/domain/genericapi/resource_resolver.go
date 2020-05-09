package genericapi

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
	"strings"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/resource"
	"github.com/kyma-project/kyma/components/function-controller/pkg/apis/serverless/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceResolver struct {
	services ResourcesServices
	converter *ResourceConverter
}

func NewResourceResolver(serviceFactory *resource.ServiceFactory) *ResourceResolver {
	schemas := []schema.GroupVersionResource{
		{
			Version:  "v1beta1",
			Group:    "servicecatalog.k8s.io",
			Resource: "serviceinstances",
		},
		{
			Version:  "v1beta1",
			Group:    "servicecatalog.k8s.io",
			Resource: "clusterserviceclasses",
		},
		{
			Version:  v1alpha1.GroupVersion.Version,
			Group:    v1alpha1.GroupVersion.Group,
			Resource: "functions",
		},
	}
	services := NewResourceServices(serviceFactory, schemas)
	converter := NewResourceConverter()

	return &ResourceResolver{
		services: services,
		converter: converter,
	}
}

func (r *ResourceResolver) Get(ctx context.Context, schema gqlschema.SchemaResourceInput, name string, namespace *string) (*gqlschema.Resource, error) {
	service := r.services.Get(schema)
	if service == nil {
		return nil, nil
	}

	item, err := service.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	return r.converter.ToGQL(item, nil)
}

func (r *ResourceResolver) List(ctx context.Context, schema gqlschema.SchemaResourceInput, namespace *string) (gqlschema.ResourceListOutput, error) {
	service := r.services.Get(schema)
	if service == nil {
		return gqlschema.ResourceListOutput{}, nil
	}

	items, err := service.List(namespace)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	return r.converter.ToGQLs(items, nil)
}

func (r *ResourceResolver) ResourceSpec(ctx context.Context, obj *gqlschema.Resource, fields []gqlschema.ResourceFieldInput, rootField *string) (gqlschema.ResourceSpecOutput, error) {
	gqlJson := gqlschema.ResourceSpecOutput{
		Data: map[string]interface{}{},
	}

	for _, field := range fields {
		pathField := strings.Split(field.Path, ".")
		if rootField != nil && *rootField != "" {
			pathField = r.prependPath(pathField, *rootField)
		}

		val, found, err := unstructured.NestedFieldCopy(obj.RawContent, pathField...)
		if err != nil {
			return gqlJson, err
		}
		if !found {
			continue
		}

		key := ""
		if field.Key != nil && *field.Key != "" {
			key = *field.Key
		} else {
			key = pathField[len(pathField) - 1]
		}

		gqlJson.Data[key] = val
	}

	return gqlJson, nil
}

func (r *ResourceResolver) ResourceSubResource(ctx context.Context, parent *gqlschema.Resource, schema gqlschema.SchemaResourceInput, name string, namespace *string) (*gqlschema.Resource, error) {
	service := r.services.Get(schema)
	if service == nil {
		return nil, nil
	}

	parsedName := r.parseArgValueForSubResource(parent, &name)
	parsedNamespace := ""
	if namespace != nil {
		parsedNamespace = *namespace
	}
	parsedNamespace = r.parseArgValueForSubResource(parent, &parsedNamespace)

	item, err := service.Get(&parsedNamespace, parsedName)
	if err != nil {
		return nil, err
	}

	return r.converter.ToGQL(item, nil)
}

func (r *ResourceResolver) ResourceSubResources(ctx context.Context, parent *gqlschema.Resource, schema gqlschema.SchemaResourceInput, namespace *string) (gqlschema.ResourceListOutput, error) {
	service := r.services.Get(schema)
	if service == nil {
		return gqlschema.ResourceListOutput{}, nil
	}

	parsedNamespace := ""
	if namespace != nil {
		parsedNamespace = *namespace
	}
	parsedNamespace = r.parseArgValueForSubResource(parent, &parsedNamespace)

	items, err := service.List(&parsedNamespace)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	return r.converter.ToGQLs(items, parent)
}

func (r *ResourceResolver) prependPath(paths []string, path string) []string {
	paths = append(paths, "")
	copy(paths[1:], paths)
	paths[0] = path
	return paths
}

func (r *ResourceResolver) parseArgValueForSubResource(parent *gqlschema.Resource, arg *string) string {
	if arg == nil {
		return ""
	}
	if parent == nil {
		return *arg
	}

	// for $parent
	deepPath := *arg
	var deepParent *gqlschema.Resource = nil
	for {
		if strings.HasPrefix(deepPath, "$parent.") {
			if deepParent == nil {
				deepParent = parent
			} else if deepParent.Parent != nil {
				deepParent = deepParent.Parent
			}

			deepPath = strings.TrimPrefix(deepPath, "$parent.")
			continue
		}

		break
	}

	if deepParent == nil {
		return deepPath
	}

	pathField := strings.Split(deepPath, ".")
	log.Print(pathField)
	val, found, err := unstructured.NestedString(deepParent.RawContent, pathField...)
	if err != nil || !found {
		return deepPath
	}
	return val
}
