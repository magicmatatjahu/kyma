package genericapi

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type ResourceConverter struct {}

func NewResourceConverter() *ResourceConverter {
	return &ResourceConverter{}
}

func (c *ResourceConverter) ToGQL(item interface{}, parent *gqlschema.Resource) (*gqlschema.Resource, error) {
	if item == nil {
		return nil, nil
	}

	unstructuredResource, err := c.toUnstructured(item)
	if unstructuredResource == nil || err != nil {
		return nil, err
	}

	return &gqlschema.Resource{
		APIVersion: unstructuredResource.GetAPIVersion(),
		Kind: unstructuredResource.GetKind(),
		Metadata: c.convertMetadata(unstructuredResource),
		Raw: unstructuredResource.UnstructuredContent(),
		Parent: parent,
	}, nil
}

func (c *ResourceConverter) ToGQLs(items []interface{}, parent *gqlschema.Resource) (gqlschema.ResourceListOutput, error) {
	output := gqlschema.ResourceListOutput{}

	resources := make([]gqlschema.Resource, 0)
	for _, item := range items {
		converted, err := c.ToGQL(item, parent)
		if err != nil {
			return output, err
		}

		if converted != nil {
			resources = append(resources, *converted)
		}
	}

	output.Items = resources
	output.TotalCount = len(resources)
	return output, nil
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

func (c *ResourceConverter) toUnstructured(item interface{}) (*unstructured.Unstructured, error) {
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
