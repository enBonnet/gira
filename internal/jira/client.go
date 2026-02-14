package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL  string
	email    string
	apiToken string
	client   *http.Client
}

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary  string    `json:"summary"`
	Status   Status    `json:"status"`
	Priority Priority  `json:"priority"`
	Assignee *Assignee `json:"assignee"`
}

type Status struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Assignee struct {
	DisplayName string `json:"displayName"`
}

type User struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Active       bool   `json:"active"`
	TimeZone     string `json:"timeZone"`
	Locale       string `json:"locale"`
}

type SearchResponse struct {
	Issues     []Issue `json:"issues"`
	Total      int     `json:"total"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
}

func NewClient(baseURL, email, apiToken string) *Client {
	return &Client{
		baseURL:  baseURL,
		email:    email,
		apiToken: apiToken,
		client:   &http.Client{},
	}
}

func (c *Client) doRequest(method, endpoint string, body []byte) ([]byte, error) {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("jira api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) GetMyTasks(project string) ([]Issue, error) {
	jql := `assignee = currentUser() AND status in ("To Do", "BLOCKED", "In Progress")`
	if project != "" {
		jql = fmt.Sprintf(`%s AND project = "%s"`, jql, project)
	}

	searchReq := map[string]interface{}{
		"jql":    jql,
		"fields": []string{"summary", "status", "priority", "assignee"},
	}

	body, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest("POST", "/rest/api/3/search/jql", body)
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Issues, nil
}

func (c *Client) GetIssueURL(key string) string {
	return strings.TrimSuffix(c.baseURL, "/") + "/browse/" + key
}

func (c *Client) GetIssue(key string) (*Issue, error) {
	data, err := c.doRequest("GET", fmt.Sprintf("/rest/api/3/issue/%s?fields=summary,status,priority,assignee", key), nil)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) GetMyself() (*User, error) {
	data, err := c.doRequest("GET", "/rest/api/3/myself", nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
