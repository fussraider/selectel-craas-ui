package craas

import (
	"log/slog"
	"regexp"
)

const Endpoint = "https://cr.selcloud.ru/api/v1"

// digestRegex matches standard SHA256 digests.
var digestRegex = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)

type Service struct {
	endpoint string
	logger   *slog.Logger
}

func New(logger *slog.Logger) *Service {
	return &Service{
		endpoint: Endpoint,
		logger:   logger.With("service", "craas"),
	}
}
