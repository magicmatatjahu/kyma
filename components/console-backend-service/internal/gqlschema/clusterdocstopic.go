package gqlschema

type ClusterDocsTopic struct {
	Name        string         `json:"name"`
	GroupName   string         `json:"groupName"`
	Assets      []ClusterAsset `json:"assets"`
	DisplayName string         `json:"displayName"`
	Description string         `json:"description"`
}