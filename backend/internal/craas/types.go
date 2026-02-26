package craas

// DeletedImage represents an image that was successfully deleted.
type DeletedImage struct {
	Digest string   `json:"digest"`
	Tags   []string `json:"tags,omitempty"`
}

// FailedImage represents an image that failed to be deleted.
type FailedImage struct {
	Digest string   `json:"digest"`
	Tags   []string `json:"tags,omitempty"`
	Error  string   `json:"error"`
}

// CleanupResult represents the result of a cleanup operation.
type CleanupResult struct {
	Deleted []DeletedImage `json:"deleted"`
	Failed  []FailedImage  `json:"failed"`
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
