package craas

// CleanupResult represents the result of a cleanup operation.
type CleanupResult struct {
	Deleted []interface{} `json:"deleted"`
	Failed  []interface{} `json:"failed"`
}

// CleanupRequest represents the request body for cleanup operation.
type CleanupRequest struct {
	Digests   []string `json:"digests"`
	DisableGC bool     `json:"disable_gc"`
	Tags      []string `json:"tags,omitempty"`
}

// GCInfo represents the garbage collection size information.
type GCInfo struct {
	SizeNonReferenced int64 `json:"sizeNonReferenced"`
	SizeSummary       int64 `json:"sizeSummary"`
	SizeUntagged      int64 `json:"sizeUntagged"`
}
