package gqlschema

type Asset struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Files     []File `json:"files"`
}