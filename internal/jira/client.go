package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Client struct {
	url      string
	email    string
	apiToken string
	client   *http.Client
}

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		Status      struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

type IssueResponse struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		Status      struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

func NewClient(url, email, apiToken string) *Client {
	return &Client{
		url:      url,
		email:    email,
		apiToken: apiToken,
		client:   &http.Client{},
	}
}

func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	url := fmt.Sprintf("%s/rest/api/3/issue/%s", c.url, issueKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.email, c.apiToken)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jira API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var issueResp IssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	issue := &Issue{
		Key: issueResp.Key,
	}
	issue.Fields.Summary = issueResp.Fields.Summary
	issue.Fields.Description = issueResp.Fields.Description
	issue.Fields.Status.Name = issueResp.Fields.Status.Name

	return issue, nil
}

func ExtractJiraKey(text string) string {
	// Match Jira issue keys like PROJ-123, ABC-456, etc.
	re := regexp.MustCompile(`([A-Z]+-\d+)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}
