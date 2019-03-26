package shared

type ClusterAsset struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	files  []File      `json:"files"`
	Status AssetStatus `json:"status"`
}