package powerbiapi

import (
	"fmt"
	"net/url"
	"time"
)

// EmbedToken represents a Power BI embed token
type EmbedToken struct {
	Token      string    `json:"token"`
	TokenID    string    `json:"tokenId"`
	Expiration time.Time `json:"expiration"`
}

// GenerateTokenRequest represents the request to generate an embed token
type GenerateTokenRequest struct {
	AccessLevel string   `json:"accessLevel,omitempty"`
	DatasetIds  []string `json:"datasetIds,omitempty"`
	ReportIds   []string `json:"reportIds,omitempty"`
	TargetWorkspaces []TargetWorkspace `json:"targetWorkspaces,omitempty"`
}

// TargetWorkspace represents a target workspace for embed token
type TargetWorkspace struct {
	ID string `json:"id"`
}

// GenerateTokenResponse represents the response from generate token API
type GenerateTokenResponse struct {
	Token      string `json:"token"`
	TokenID    string `json:"tokenId"`
	Expiration string `json:"expiration"`
}

// GenerateEmbedToken generates an embed token for reports
func (client *Client) GenerateEmbedToken(workspaceID string, request GenerateTokenRequest) (*GenerateTokenResponse, error) {
	var respObj GenerateTokenResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports/GenerateToken", url.PathEscape(workspaceID))
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}

// GenerateEmbedTokenForReport generates an embed token for a specific report
func (client *Client) GenerateEmbedTokenForReport(workspaceID, reportID string, request GenerateTokenRequest) (*GenerateTokenResponse, error) {
	var respObj GenerateTokenResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports/%s/GenerateToken", 
		url.PathEscape(workspaceID), url.PathEscape(reportID))
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}

// GenerateEmbedTokenForDataset generates an embed token for a dataset
func (client *Client) GenerateEmbedTokenForDataset(workspaceID, datasetID string, request GenerateTokenRequest) (*GenerateTokenResponse, error) {
	var respObj GenerateTokenResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/GenerateToken", 
		url.PathEscape(workspaceID), url.PathEscape(datasetID))
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}

// GenerateEmbedTokenForDashboard generates an embed token for a dashboard
func (client *Client) GenerateEmbedTokenForDashboard(workspaceID, dashboardID string, request GenerateTokenRequest) (*GenerateTokenResponse, error) {
	var respObj GenerateTokenResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s/GenerateToken", 
		url.PathEscape(workspaceID), url.PathEscape(dashboardID))
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}

// GenerateEmbedTokenForTile generates an embed token for a dashboard tile
func (client *Client) GenerateEmbedTokenForTile(workspaceID, dashboardID, tileID string, request GenerateTokenRequest) (*GenerateTokenResponse, error) {
	var respObj GenerateTokenResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s/tiles/%s/GenerateToken", 
		url.PathEscape(workspaceID), url.PathEscape(dashboardID), url.PathEscape(tileID))
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}