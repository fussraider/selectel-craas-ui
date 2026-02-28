package craas

import (
	"log/slog"
	"regexp"

	"github.com/generic/selectel-craas-web/internal/config"
)

// digestRegex matches standard SHA256 digests.
var digestRegex = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)

type Service struct {
	endpoint string
	logger   *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		endpoint: cfg.SelectelCraasURL,
		logger:   logger.With("service", "craas"),
	}
}
