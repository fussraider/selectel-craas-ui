package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func New(cfg *config.Config) *Client {
	return &Client{
		cfg:           cfg,
		client:        &http.Client{Timeout: 10 * time.Second},
		AuthURL:       "https://cloud.api.selcloud.ru/identity/v3/auth/tokens",
		ProjURL:       "https://cloud.api.selcloud.ru/identity/v3/auth/projects",
		projectTokens: make(map[string]string),
	}
}

func (c *Client) GetAccountToken() (string, error) {
	c.mu.Lock()
	if c.accountToken != "" {
		token := c.accountToken
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	payload := map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
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
			},
			"scope": map[string]interface{}{
				"domain": map[string]interface{}{
					"name": c.cfg.SelectelAccountID,
				},
			},
		},
	}

	token, err := c.requestToken(payload)
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.accountToken = token
	c.mu.Unlock()
	return token, nil
}

func (c *Client) InvalidateAccountToken() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accountToken = ""
}

func (c *Client) GetProjectToken(projectID string) (string, error) {
	c.mu.Lock()
	if token, ok := c.projectTokens[projectID]; ok {
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	payload := map[string]interface{}{
		"auth": map[string]interface{}{
			"identity": map[string]interface{}{
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
			},
			"scope": map[string]interface{}{
				"project": map[string]interface{}{
					"id": projectID,
				},
			},
		},
	}

	token, err := c.requestToken(payload)
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.projectTokens[projectID] = token
	c.mu.Unlock()
	return token, nil
}

func (c *Client) InvalidateProjectToken(projectID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.projectTokens, projectID)
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

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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
	req, err := http.NewRequest("GET", c.ProjURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", accountToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

	return result.Projects, nil
}
