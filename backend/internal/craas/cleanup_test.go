package craas

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCleanupRepository_Security(t *testing.T) {
	// Logger for service
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Received EscapedPath: %s", r.URL.EscapedPath())

		// We use a registryID with a slash to test encoding
		// registryID = "reg/dangerous"

		if strings.Contains(r.URL.EscapedPath(), "reg/dangerous") {
			// This means the slash was NOT encoded -> Vulnerable
			t.Log("Detected unencoded slash in registry ID - Vulnerable behavior")
			w.WriteHeader(http.StatusOK)
			return
		}

		if strings.Contains(r.URL.EscapedPath(), "reg%2Fdangerous") {
			// This means the slash WAS encoded -> Secure behavior
			t.Log("Detected encoded slash in registry ID - Secure behavior")
			w.WriteHeader(http.StatusTeapot) // Use Teapot to distinguish secure path
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	svc := &Service{endpoint: ts.URL, logger: logger}

	// Malicious registry ID containing a slash
	registryID := "reg/dangerous"
	repoName := "myrepo"
	digests := []string{"sha256:123"}

	_, err := svc.CleanupRepository(context.Background(), "token", registryID, repoName, digests, false)

	if err != nil {
		if strings.Contains(err.Error(), "418") {
			// Success - Secure path taken
			return
		}
		t.Fatalf("Unexpected error: %v", err)
	}

	// If no error, it means we hit the 200 OK path => Vulnerable
	t.Fatalf("Vulnerability detected: Request path contained unencoded slash!")
}
