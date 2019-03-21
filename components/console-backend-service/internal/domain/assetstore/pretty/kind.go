package pretty

type Kind int

const (
	Asset Kind = iota
	Assets
	ClusterAsset
	ClusterAssets
)

func (k Kind) String() string {
	switch k {
	case Asset:
		return "Asset"
	case Assets:
		return "Assets"
	case ClusterAsset:
		return "Cluster Asset"
	case ClusterAssets:
		return "Cluster Assets"
	default:
		return ""
	}
}
