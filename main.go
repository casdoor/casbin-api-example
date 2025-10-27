package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// CasdoorConfig holds the configuration for Casdoor API
type CasdoorConfig struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	Organization string
	Application  string
}

// CasbinRequest represents a request to Casbin API
type CasbinRequest struct {
	PermissionID string   `json:"permissionId"`
	Enforcer     string   `json:"enforcer,omitempty"`
	Params       []string `json:"params,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	Object       string   `json:"object,omitempty"`
	Action       string   `json:"action,omitempty"`
	V0           string   `json:"v0,omitempty"`
	V1           string   `json:"v1,omitempty"`
	V2           string   `json:"v2,omitempty"`
	V3           string   `json:"v3,omitempty"`
	V4           string   `json:"v4,omitempty"`
	V5           string   `json:"v5,omitempty"`
}

// CasdoorClient handles API calls to Casdoor
type CasdoorClient struct {
	config CasdoorConfig
	token  string
}

// NewCasdoorClient creates a new Casdoor client
func NewCasdoorClient(config CasdoorConfig) *CasdoorClient {
	return &CasdoorClient{
		config: config,
	}
}

// makeRequest makes an HTTP request to Casdoor API
func (c *CasdoorClient) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.config.Endpoint + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Enforce checks if a request is allowed by the policy
func (c *CasdoorClient) Enforce(permissionID string, params []string) (bool, error) {
	req := CasbinRequest{
		PermissionID: permissionID,
		Params:       params,
	}

	respBody, err := c.makeRequest("POST", "/api/enforce", req)
	if err != nil {
		return false, err
	}

	var result struct {
		Status string `json:"status"`
		Data   bool   `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return false, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Data, nil
}

// BatchEnforce checks multiple requests at once
func (c *CasdoorClient) BatchEnforce(permissionID string, paramsArray [][]string) ([]bool, error) {
	type BatchRequest struct {
		PermissionID string     `json:"permissionId"`
		ParamsArray  [][]string `json:"paramsArray"`
	}

	req := BatchRequest{
		PermissionID: permissionID,
		ParamsArray:  paramsArray,
	}

	respBody, err := c.makeRequest("POST", "/api/batch-enforce", req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string `json:"status"`
		Data   []bool `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Data, nil
}

// GetAllObjects returns all objects for a subject
func (c *CasdoorClient) GetAllObjects(permissionID, subject string) ([]string, error) {
	req := CasbinRequest{
		PermissionID: permissionID,
		Subject:      subject,
	}

	respBody, err := c.makeRequest("POST", "/api/get-all-objects", req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Data, nil
}

// GetAllActions returns all actions for a subject and object
func (c *CasdoorClient) GetAllActions(permissionID, subject, object string) ([]string, error) {
	req := CasbinRequest{
		PermissionID: permissionID,
		Subject:      subject,
		Object:       object,
	}

	respBody, err := c.makeRequest("POST", "/api/get-all-actions", req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Data, nil
}

// GetAllRoles returns all roles for a subject
func (c *CasdoorClient) GetAllRoles(permissionID, subject string) ([]string, error) {
	req := CasbinRequest{
		PermissionID: permissionID,
		Subject:      subject,
	}

	respBody, err := c.makeRequest("POST", "/api/get-all-roles", req)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Data, nil
}

func main() {
	// Load configuration from environment variables
	config := CasdoorConfig{
		Endpoint:     getEnv("CASDOOR_ENDPOINT", "http://localhost:8000"),
		ClientID:     getEnv("CASDOOR_CLIENT_ID", ""),
		ClientSecret: getEnv("CASDOOR_CLIENT_SECRET", ""),
		Organization: getEnv("CASDOOR_ORGANIZATION", "built-in"),
		Application:  getEnv("CASDOOR_APPLICATION", "app-built-in"),
	}

	client := NewCasdoorClient(config)

	fmt.Println("=== Casdoor Casbin API Example ===")
	fmt.Println()

	// Example 1: Enforce - Check if alice can read data1
	fmt.Println("Example 1: Enforce - Check if alice can read data1")
	permissionID := "built-in/permission-built-in"
	allowed, err := client.Enforce(permissionID, []string{"alice", "data1", "read"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: alice can read data1 = %v\n", allowed)
	}
	fmt.Println()

	// Example 2: BatchEnforce - Check multiple permissions at once
	fmt.Println("Example 2: BatchEnforce - Check multiple permissions")
	requests := [][]string{
		{"alice", "data1", "read"},
		{"alice", "data1", "write"},
		{"bob", "data2", "read"},
	}
	results, err := client.BatchEnforce(permissionID, requests)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		for i, result := range results {
			fmt.Printf("Request %v: %v\n", requests[i], result)
		}
	}
	fmt.Println()

	// Example 3: GetAllObjects - Get all objects alice can access
	fmt.Println("Example 3: GetAllObjects - Get all objects alice can access")
	objects, err := client.GetAllObjects(permissionID, "alice")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Objects alice can access: %v\n", objects)
	}
	fmt.Println()

	// Example 4: GetAllActions - Get all actions alice can perform on data1
	fmt.Println("Example 4: GetAllActions - Get all actions alice can perform on data1")
	actions, err := client.GetAllActions(permissionID, "alice", "data1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Actions alice can perform on data1: %v\n", actions)
	}
	fmt.Println()

	// Example 5: GetAllRoles - Get all roles for alice
	fmt.Println("Example 5: GetAllRoles - Get all roles for alice")
	roles, err := client.GetAllRoles(permissionID, "alice")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Roles for alice: %v\n", roles)
	}
	fmt.Println()

	fmt.Println("=== Demo completed ===")
	fmt.Println()
	fmt.Println("Note: This demo requires a running Casdoor instance.")
	fmt.Println("Please configure environment variables or update the config in the code.")
	fmt.Println("For more information, visit: https://casdoor.org/docs/permission/exposed-casbin-apis/")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
