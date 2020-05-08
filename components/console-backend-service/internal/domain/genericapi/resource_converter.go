package genericapi

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type resourceConverter struct {}

func newResourceConverter() *resourceConverter {
	return &resourceConverter{}
}

func (c *resourceConverter) ToGQL(item interface{}) (*gqlschema.Resource, error) {
	unstructuredResource, err := c.toUnstructured(item)
	if unstructuredResource == nil || err != nil {
		return nil, err
	}

	return &gqlschema.Resource{
		APIVersion: unstructuredResource.GetAPIVersion(),
		Kind: unstructuredResource.GetKind(),
		Metadata: c.convertMetadata(unstructuredResource),
		RawContent: unstructuredResource.UnstructuredContent(),
	}, nil
}

func (c *resourceConverter) ToGQLs(items []interface{}) (gqlschema.ResourceListOutput, error) {
	output := gqlschema.ResourceListOutput{}
	resources := make([]gqlschema.Resource, 0)
	for _, item := range items {
		converted, err := c.ToGQL(item)
		if err != nil {
			return gqlschema.ResourceListOutput{}, err
		}

		if converted != nil {
			resources = append(resources, *converted)
		}
	}

	output.Items = resources
	output.ItemsCount = len(resources)

	return output, nil
}

func (c *resourceConverter) convertMetadata(unstructuredResource *unstructured.Unstructured) gqlschema.ResourceMetadata {
	ownerReferences := c.convertOwnerReferences(unstructuredResource)

	return gqlschema.ResourceMetadata{
		Name: unstructuredResource.GetName(),
		Namespace: unstructuredResource.GetNamespace(),
		GenerateName: unstructuredResource.GetGenerateName(),
		Labels: unstructuredResource.GetLabels(),
		Annotations: unstructuredResource.GetAnnotations(),
		OwnerReferences: ownerReferences,
	}
}

func (c *resourceConverter) convertOwnerReferences(unstructuredResource *unstructured.Unstructured) []gqlschema.OwnerReferenceType {
	gqlOwnerReferences := make([]gqlschema.OwnerReferenceType, 0)
	ownerReferences := unstructuredResource.GetOwnerReferences()

	for _, ownerRef := range ownerReferences {
		gqlOwnerRef := gqlschema.OwnerReferenceType{
			APIVersion: ownerRef.APIVersion,
			Kind: ownerRef.Kind,
			Name: ownerRef.Name,
			UID: string(ownerRef.UID),
			Controller: ownerRef.Controller,
			BlockOwnerDeletion: ownerRef.BlockOwnerDeletion,
		}
		gqlOwnerReferences = append(gqlOwnerReferences, gqlOwnerRef)
	}

	return gqlOwnerReferences
}

func (c *resourceConverter) toUnstructured(item interface{}) (*unstructured.Unstructured, error) {
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
