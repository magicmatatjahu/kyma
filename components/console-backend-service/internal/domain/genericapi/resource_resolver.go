package genericapi

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
)

type ResourceResolver struct {
	services ResourcesServices
	converter *ResourceConverter
	pager *ResourcePager
	sort *ResourceSort
	filter *ResourceFilter
}

func NewResourceResolver(services ResourcesServices, converter *ResourceConverter, pager *ResourcePager, sort *ResourceSort, filter *ResourceFilter) *ResourceResolver {
	return &ResourceResolver{
		services: services,
		converter: converter,
		pager: pager,
		sort: sort,
		filter: filter,
	}
}

func (r *ResourceResolver) Get(ctx context.Context, schema string, name string, namespace *string) (*gqlschema.Resource, error) {
	return r.get(ctx, schema, name, namespace)
}

func (r *ResourceResolver) List(ctx context.Context, schema string, namespace *string, pager *gqlschema.ResourcePager, options *gqlschema.ResourceListOptions) (gqlschema.ResourceListOutput, error) {
	return r.list(ctx, schema, namespace, pager, options)
}

func (r *ResourceResolver) Watch(ctx context.Context, schema string, namespace, name *string) (<-chan gqlschema.ResourceEvent, error) {
	service := r.services.Get(schema)
	if service == nil {
		return nil, nil
	}

	channel := make(chan gqlschema.ResourceEvent, 1)
	filter := func(entity interface{}) bool {
		if entity == nil {
			return false
		}
		return true
	}

	listener := NewResourceListener(channel, filter, r.converter)
	service.Subscribe(listener)
	go func() {
		defer close(channel)
		defer service.Unsubscribe(listener)
		<-ctx.Done()
	}()

	return channel, nil
}


func (r *ResourceResolver) ResourceSpec(ctx context.Context, obj *gqlschema.Resource, fields []gqlschema.ResourceFieldInput, rootField *string) (gqlschema.JSON, error) {
	gqlJSON := gqlschema.JSON{}

	for _, field := range fields {
		pathField := strings.Split(field.Path, ".")
		if rootField != nil && *rootField != "" {
			pathField = r.prependPath(pathField, *rootField)
		}

		val, found, err := unstructured.NestedFieldCopy(obj.Raw, pathField...)
		if err != nil {
			return nil, err
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

		gqlJSON[key] = val
	}

	return gqlJSON, nil
}

func (r *ResourceResolver) ResourceSubResource(ctx context.Context, parent *gqlschema.Resource, schema string, name string, namespace *string) (*gqlschema.Resource, error) {
	parsedName := r.parseArgValueForSubResource(parent, &name)
	parsedNamespace := ""
	if namespace != nil {
		parsedNamespace = *namespace
	}
	parsedNamespace = r.parseArgValueForSubResource(parent, &parsedNamespace)

	return r.get(ctx, schema, parsedName, &parsedNamespace)
}

func (r *ResourceResolver) ResourceSubResources(ctx context.Context, parent *gqlschema.Resource, schema string, namespace *string, pager *gqlschema.ResourcePager, options *gqlschema.ResourceListOptions) (gqlschema.ResourceListOutput, error) {
	parsedNamespace := ""
	if namespace != nil {
		parsedNamespace = *namespace
	}
	parsedNamespace = r.parseArgValueForSubResource(parent, &parsedNamespace)

	return r.list(ctx, schema, &parsedNamespace, pager, options)
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
				deepParent = &(*parent)
			} else if deepParent.Parent != nil {
				deepParent = &(*deepParent.Parent)
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
	val, found, err := unstructured.NestedString(deepParent.Raw, pathField...)
	if err != nil || !found {
		return deepPath
	}
	return val
}

func (r *ResourceResolver) get(ctx context.Context, schema string, name string, namespace *string) (*gqlschema.Resource, error) {
	service := r.services.Get(schema)
	if service == nil {
		return nil, nil
	}

	item, err := service.Get(namespace, name)
	if err != nil {
		return nil, err
	}

	unstructuredItem, err := r.converter.ToUnstructured(item)
	if err != nil {
		return nil, err
	}

	gqlItem, err := r.converter.ToGQL(unstructuredItem, nil)
	if err != nil {
		return nil, err
	}

	return gqlItem, nil
}

func (r *ResourceResolver) list(ctx context.Context, schema string, namespace *string, pager *gqlschema.ResourcePager, options *gqlschema.ResourceListOptions) (gqlschema.ResourceListOutput, error) {
	service := r.services.Get(schema)
	if service == nil {
		return gqlschema.ResourceListOutput{}, nil
	}

	items, err := service.List(namespace, options)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	unstructuredItems, err := r.converter.ToUnstructuredList(items)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	filteredItems, err := r.withFiltering(unstructuredItems, options)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	sortedItems, err := r.withSorting(filteredItems, options)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	paginatedItems, err := r.withPagination(sortedItems, pager)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	gqlItems, err := r.converter.ToGQLs(paginatedItems, nil)
	if err != nil {
		return gqlschema.ResourceListOutput{}, err
	}

	return r.converter.ToListOutput(gqlItems, false), nil
}

func (r *ResourceResolver) withFiltering(items []unstructured.Unstructured, options *gqlschema.ResourceListOptions) ([]unstructured.Unstructured, error) {
	if options != nil && options.Filter != nil && len(options.Filter) > 0 {
		return r.filter.Filter(items, options.Filter)
	}
	return items, nil
}

func (r *ResourceResolver) withSorting(items []unstructured.Unstructured, options *gqlschema.ResourceListOptions) ([]unstructured.Unstructured, error) {
	if options != nil && options.Sort != nil && len(options.Sort) > 0 {
		return r.sort.Sort(items, options.Sort)
	}
	return r.sort.Sort(items, nil)
}

func (r *ResourceResolver) withPagination(items []unstructured.Unstructured, pager *gqlschema.ResourcePager) ([]unstructured.Unstructured, error) {
	if pager != nil {
		return r.pager.Page(items, *pager)
	}
	return items, nil
}
