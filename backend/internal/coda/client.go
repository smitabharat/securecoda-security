package coda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://coda.io/apis/v1"

type Client struct {
	token string
	http  *http.Client
}

func NewClient(token string, logger interface{}) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) doRequest(method, url string, body interface{}) ([]byte, error) {
	var buf io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, baseURL+url, buf)
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("coda api error %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// Example: List docs
func (c *Client) ListDocs() ([]map[string]interface{}, error) {
	data, err := c.doRequest("GET", "/docs", nil)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return parsed.Items, nil
}

// Get doc metadata
func (c *Client) GetDoc(docID string) (map[string]interface{}, error) {
	data, err := c.doRequest("GET", "/docs/"+docID, nil)
	if err != nil {
		return nil, err
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// List tables for a doc
func (c *Client) ListTables(docID string) ([]map[string]interface{}, error) {
	data, err := c.doRequest("GET", "/docs/"+docID+"/tables", nil)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return parsed.Items, nil
}

// Get rows for a table
func (c *Client) GetRows(docID, tableID string) ([]map[string]interface{}, error) {
	data, err := c.doRequest("GET", "/docs/"+docID+"/tables/"+tableID+"/rows", nil)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return parsed.Items, nil
}

// Export page as HTML
func (c *Client) ExportPageHTML(docID, pageID string) (string, error) {
	data, err := c.doRequest("GET", "/docs/"+docID+"/pages/"+pageID+"?export=html", nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Remediation: delete row
func (c *Client) DeleteRow(docID, tableID, rowID string) error {
	_, err := c.doRequest("DELETE", "/docs/"+docID+"/tables/"+tableID+"/rows/"+rowID, nil)
	return err
}

// Remediation: unpublish (simulate via doc update -- Coda API may vary)
func (c *Client) UnpublishDoc(docID string) error {
	// This is a placeholder: actual Coda API may require PUT/POST with proper payload
	payload := map[string]interface{}{"isPublished": false}
	_, err := c.doRequest("PUT", "/docs/"+docID, payload)
	return err
}
