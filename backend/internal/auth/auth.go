package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/generic/selectel-craas-web/internal/config"
)

type Authenticator interface {
	GetAccountToken() (string, error)
	InvalidateAccountToken()
	ListProjects(accountToken string) ([]Project, error)
	GetProjectToken(projectID string) (string, error)
	InvalidateProjectToken(projectID string)
}

type Client struct {
	cfg           *config.Config
	client        *http.Client
	AuthURL       string
	ProjURL       string
	mu            sync.Mutex
	accountToken  string
	projectTokens map[string]string
	logger        *slog.Logger
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func New(cfg *config.Config, logger *slog.Logger) *Client {
	return &Client{
		cfg:           cfg,
		client:        &http.Client{Timeout: 60 * time.Second},
		AuthURL:       cfg.SelectelAuthURL,
		ProjURL:       cfg.SelectelProjURL,
		projectTokens: make(map[string]string),
		logger:        logger.With("service", "auth"),
	}
}

func (c *Client) getAuthPayload(projectID string) map[string]interface{} {
	userAuth := map[string]interface{}{
		"methods": []string{"password"},
		"password": map[string]interface{}{
			"user": map[string]interface{}{
				"name": c.cfg.SelectelUsername,
				"domain": map[string]interface{}{
					"name": c.cfg.SelectelAccountID,
				},
				"password": c.cfg.SelectelPassword,
			},
		},
	}

	var scope map[string]interface{}
	if projectID != "" {
		scope = map[string]interface{}{
			"project": map[string]interface{}{
				"id": projectID,
			},
		}
	} else {
		scope = map[string]interface{}{
			"project": map[string]interface{}{
				"name": c.cfg.SelectelProjectName,
				"domain": map[string]interface{}{
					"name": c.cfg.SelectelAccountID,
				},
			},
		}
	}

	return map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": userAuth,
			"scope":    scope,
		},
	}
}

func (c *Client) GetAccountToken() (string, error) {
	c.mu.Lock()
	if c.accountToken != "" {
		token := c.accountToken
		c.mu.Unlock()
		c.logger.Debug("cache hit for account token")
		return token, nil
	}
	c.mu.Unlock()

	c.logger.Debug("requesting new account token")

	token, err := c.requestToken(c.getAuthPayload(""))
	if err != nil {
		c.logger.Error("failed to get account token", "error", err)
		return "", err
	}

	c.mu.Lock()
	c.accountToken = token
	c.mu.Unlock()

	c.logger.Debug("successfully acquired account token")
	return token, nil
}

func (c *Client) InvalidateAccountToken() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accountToken = ""
	c.logger.Debug("invalidated account token")
}

func (c *Client) GetProjectToken(projectID string) (string, error) {
	c.mu.Lock()
	if token, ok := c.projectTokens[projectID]; ok {
		c.mu.Unlock()
		c.logger.Debug("cache hit for project token", "project_id", projectID)
		return token, nil
	}
	c.mu.Unlock()

	c.logger.Debug("requesting new project token", "project_id", projectID)

	token, err := c.requestToken(c.getAuthPayload(projectID))
	if err != nil {
		c.logger.Error("failed to get project token", "project_id", projectID, "error", err)
		return "", err
	}

	c.mu.Lock()
	c.projectTokens[projectID] = token
	c.mu.Unlock()

	c.logger.Debug("successfully acquired project token", "project_id", projectID)
	return token, nil
}

func (c *Client) InvalidateProjectToken(projectID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.projectTokens, projectID)
	c.logger.Debug("invalidated project token", "project_id", projectID)
}

func (c *Client) requestToken(payload map[string]interface{}) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.AuthURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	c.logger.Debug("token request completed", "status", resp.StatusCode, "duration", duration)

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	token := resp.Header.Get("X-Subject-Token")
	if token == "" {
		return "", fmt.Errorf("response did not contain X-Subject-Token header")
	}

	return token, nil
}

// ListProjects lists projects accessible by the account token.
func (c *Client) ListProjects(accountToken string) ([]Project, error) {
	c.logger.Debug("listing projects")
	req, err := http.NewRequest("GET", c.ProjURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", accountToken)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.client.Do(req)
	duration := time.Since(start)
	if err != nil {
		c.logger.Error("failed to list projects request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	c.logger.Debug("list projects request completed", "status", resp.StatusCode, "duration", duration)

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list projects: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Projects []Project `json:"projects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	c.logger.Info("listed projects", "count", len(result.Projects))
	return result.Projects, nil
}
