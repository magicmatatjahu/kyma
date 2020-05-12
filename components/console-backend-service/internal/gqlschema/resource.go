package gqlschema

type Resource struct {
	APIVersion        string   `json:"apiVersion"`
	Kind      string `json:"kind"`
	Metadata            ResourceMetadata   `json:"metadata"`
	Raw            JSON   `json:"raw"`
	Parent *Resource `json:"parent"`
}
