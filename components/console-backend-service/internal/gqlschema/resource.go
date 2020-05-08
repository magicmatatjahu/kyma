package gqlschema

type Resource struct {
	ApiVersion        string   `json:"apiVersion"`
	Kind      string `json:"kind"`
	Metadata            ResourceMetadata   `json:"metadata"`
	RawContent            JSON   `json:"rawContent"`
}
