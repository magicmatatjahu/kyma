package genericapi

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type ResourceConverter struct {
	pager *ResourcePager
}

func NewResourceConverter(pager *ResourcePager) *ResourceConverter {
	return &ResourceConverter{
		pager: pager,
	}
}

func (c *ResourceConverter) ToGQL(item *unstructured.Unstructured, parent *gqlschema.Resource) (*gqlschema.Resource, error) {
	if item == nil {
		return nil, nil
	}

	return &gqlschema.Resource{
		APIVersion: item.GetAPIVersion(),
		Kind: item.GetKind(),
		Metadata: c.convertMetadata(item),
		Raw: item.UnstructuredContent(),
		Parent: parent,
	}, nil
}

func (c *ResourceConverter) ToGQLs(items []unstructured.Unstructured, parent *gqlschema.Resource) ([]gqlschema.Resource, error) {
	resources := make([]gqlschema.Resource, 0)
	for _, item := range items {
		converted, err := c.ToGQL(&item, parent)
		if err != nil {
			return []gqlschema.Resource{}, err
		}

		if converted != nil {
			resources = append(resources, *converted)
		}
	}
	return resources, nil
}

func (c *ResourceConverter) ToListOutput(gqlItems []gqlschema.Resource, withPagination bool) gqlschema.ResourceListOutput {
	return gqlschema.ResourceListOutput{
		Edges: c.toEdges(gqlItems, withPagination),
		Nodes: gqlItems,
		TotalCount: len(gqlItems),
	}
}

func (c *ResourceConverter) ToUnstructured(item interface{}) (*unstructured.Unstructured, error) {
	if item == nil {
		return nil, nil
	}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(item)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s to unstructured", item)
	}
	if len(u) == 0 {
		return nil, nil
	}

	return &unstructured.Unstructured{Object: u}, nil
}

func (c *ResourceConverter) ToUnstructuredList(items []interface{}) ([]unstructured.Unstructured, error) {
	uList := make([]unstructured.Unstructured, 0)
	for _, item := range items {
		u, err := c.ToUnstructured(item)
		if err != nil {
			return []unstructured.Unstructured{}, err
		}

		if u != nil {
			uList = append(uList, *u)
		}
	}
	return uList, nil
}

func (c *ResourceConverter) convertMetadata(unstructuredResource *unstructured.Unstructured) gqlschema.ResourceMetadata {
	return gqlschema.ResourceMetadata{
		Name: unstructuredResource.GetName(),
		Namespace: unstructuredResource.GetNamespace(),
		Labels: unstructuredResource.GetLabels(),
		Annotations: unstructuredResource.GetAnnotations(),
		GenerateName: unstructuredResource.GetGenerateName(),
		UID: string(unstructuredResource.GetUID()),
		OwnerReferences: c.convertOwnerReferences(unstructuredResource.GetOwnerReferences()),
		ResourceVersion: unstructuredResource.GetResourceVersion(),
		Generation: int(unstructuredResource.GetGeneration()),
		SelfLink: unstructuredResource.GetSelfLink(),
		Continue: unstructuredResource.GetContinue(),
		RemainingItemCount: c.convertToInt64Ptr(unstructuredResource.GetRemainingItemCount()),
		CreationTimestamp: unstructuredResource.GetCreationTimestamp().Time,
		DeletionTimestamp: c.convertToTimestampPtr(unstructuredResource.GetDeletionTimestamp()),
		DeletionGracePeriodSeconds: c.convertToInt64Ptr(unstructuredResource.GetDeletionGracePeriodSeconds()),
		Finalizers: unstructuredResource.GetFinalizers(),
		ClusterName: unstructuredResource.GetClusterName(),
		ManagedFields: c.convertManagedFieldsEntry(unstructuredResource.GetManagedFields()),
	}
}

func (c *ResourceConverter) convertOwnerReferences(owners []v1.OwnerReference) []gqlschema.OwnerReferenceType {
	gqlOwnerReferences := make([]gqlschema.OwnerReferenceType, 0)
	for _, ownerRef := range owners {
		gqlOwnerReferences = append(gqlOwnerReferences, gqlschema.OwnerReferenceType{
			APIVersion: ownerRef.APIVersion,
			Kind: ownerRef.Kind,
			Name: ownerRef.Name,
			UID: string(ownerRef.UID),
			Controller: ownerRef.Controller,
			BlockOwnerDeletion: ownerRef.BlockOwnerDeletion,
		})
	}
	return gqlOwnerReferences
}

func (c *ResourceConverter) convertManagedFieldsEntry(fields []v1.ManagedFieldsEntry) []gqlschema.ManagedField {
	gqlOwnerReferences := make([]gqlschema.ManagedField, 0)
	for _, field := range fields {
		var fieldsV1 *string
		if field.FieldsV1 != nil {
			raw := string(field.FieldsV1.Raw)
			fieldsV1 = &raw
		}

		gqlOwnerReferences = append(gqlOwnerReferences, gqlschema.ManagedField{
			Manager: field.Manager,
			Operation: gqlschema.ManagedFieldsOperationType(field.Operation),
			APIVersion: field.APIVersion,
			Time: c.convertToTimestampPtr(field.Time),
			FieldsType: field.FieldsType,
			FieldsV1: fieldsV1,
		})
	}
	return gqlOwnerReferences
}

func (c *ResourceConverter) convertToInt64Ptr(num *int64) *int {
	var ptr *int
	if num != nil {
		intNum := int(*num)
		ptr = &intNum
	}
	return ptr
}

func (c *ResourceConverter) convertToTimestampPtr(timestamp *v1.Time) *time.Time {
	var t *time.Time
	if !timestamp.IsZero() {
		t = &timestamp.Time
	}
	return t
}

func (c *ResourceConverter) toEdges(resources []gqlschema.Resource, withPagination bool) []gqlschema.ResourceListEdges {
	resourcesLength := len(resources)
	edges := make([]gqlschema.ResourceListEdges, 0)
	for idx, resource := range resources {
		var prev, next *gqlschema.Resource
		if idx != 0 {
			prev = &resources[idx - 1]
		}
		if idx < resourcesLength  {
			next = &resources[idx + 1]
		}

		var cursor *string
		if withPagination {
			cursorUID := c.pager.EncodeNextCursor(resource.Metadata.UID)
			cursor = &cursorUID
		}

		edges = append(edges, gqlschema.ResourceListEdges{
			Prev: prev,
			Node: &resource,
			Next: next,
			Cursor: cursor,
		})
	}
	return edges
}
