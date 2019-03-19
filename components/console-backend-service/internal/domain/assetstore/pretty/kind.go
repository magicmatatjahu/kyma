package pretty

type Kind int

const (
	Asset Kind = iota
	ClusterAsset
)

func (k Kind) String() string {
	switch k {
	case Asset:
		return "Asset"
	case ClusterAsset:
		return "Cluster Asset"
	default:
		return ""
	}
}
