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
