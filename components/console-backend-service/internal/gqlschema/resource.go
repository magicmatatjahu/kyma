package gqlschema

type Resource struct {
	APIVersion        string   `json:"apiVersion"`
	Kind      string `json:"kind"`
	Metadata            ResourceMetadata   `json:"metadata"`
	RawContent            JSON   `json:"rawContent"`
}
