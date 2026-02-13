package types

// Common types

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Registry struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	Status    string `json:"status"`
}

type Repository struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	UpdatedAt string `json:"updatedAt"`
}

type Image struct {
	Digest    string   `json:"digest"`
	Tags      []string `json:"tags"`
	Size      int64    `json:"size"`
	CreatedAt string   `json:"createdAt"`
}
