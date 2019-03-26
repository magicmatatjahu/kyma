package shared

type Asset struct {
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Type      string      `json:"type"`
	files     []File      `json:"files"`
	Status    AssetStatus `json:"status"`
}

type AssetStatus struct {
	Phase   AssetPhaseType `json:"phase"`
	Reason  string         `json:"reason"`
	Message string         `json:"message"`
}

type AssetPhaseType string

const (
	AssetPhaseTypeReady   AssetPhaseType = "READY"
	AssetPhaseTypePending AssetPhaseType = "PENDING"
	AssetPhaseTypeFailed  AssetPhaseType = "FAILED"
)
