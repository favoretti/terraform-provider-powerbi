package powerbiapi

import (
	"fmt"
	"net/url"
)

// Dashboard represents a Power BI dashboard
type Dashboard struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	IsReadOnly  bool   `json:"isReadOnly"`
	EmbedURL    string `json:"embedUrl"`
	WebURL      string `json:"webUrl"`
}

// CreateDashboardRequest represents the request for creating a dashboard
type CreateDashboardRequest struct {
	Name string `json:"name"`
}

// GetDashboardsResponse represents the response from GetDashboards API
type GetDashboardsResponse struct {
	Value []Dashboard `json:"value"`
}

// Tile represents a dashboard tile
type Tile struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	SubTitle      string  `json:"subTitle"`
	EmbedURL      string  `json:"embedUrl"`
	EmbedData     string  `json:"embedData"`
	ReportID      string  `json:"reportId"`
	DatasetID     string  `json:"datasetId"`
	RowSpan       int     `json:"rowSpan"`
	ColSpan       int     `json:"colSpan"`
	Configuration *string `json:"configuration,omitempty"`
}

// GetTilesResponse represents the response from GetTiles API
type GetTilesResponse struct {
	Value []Tile `json:"value"`
}

// CloneTileRequest represents the request for cloning a tile
type CloneTileRequest struct {
	TargetDashboardID     string `json:"targetDashboardId"`
	TargetWorkspaceID     string `json:"targetWorkspaceId,omitempty"`
	TargetReportID        string `json:"targetReportId,omitempty"`
	TargetModelID         string `json:"targetModelId,omitempty"`
	PositionConflictAction string `json:"positionConflictAction,omitempty"`
}

// CreateDashboard creates a new dashboard in a workspace
func (client *Client) CreateDashboard(groupID string, request CreateDashboardRequest) (*Dashboard, error) {
	var respObj Dashboard
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards", url.PathEscape(groupID))
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// CreateDashboardInMyWorkspace creates a new dashboard in My Workspace
func (client *Client) CreateDashboardInMyWorkspace(request CreateDashboardRequest) (*Dashboard, error) {
	var respObj Dashboard
	url := "https://api.powerbi.com/v1.0/myorg/dashboards"
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// GetDashboards returns a list of dashboards in a workspace
func (client *Client) GetDashboards(groupID string) (*GetDashboardsResponse, error) {
	var respObj GetDashboardsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetDashboardsInMyWorkspace returns a list of dashboards in My Workspace
func (client *Client) GetDashboardsInMyWorkspace() (*GetDashboardsResponse, error) {
	var respObj GetDashboardsResponse
	url := "https://api.powerbi.com/v1.0/myorg/dashboards"
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetDashboard returns a specific dashboard
func (client *Client) GetDashboard(groupID, dashboardID string) (*Dashboard, error) {
	var respObj Dashboard
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s", 
		url.PathEscape(groupID), url.PathEscape(dashboardID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetDashboardInMyWorkspace returns a specific dashboard from My Workspace
func (client *Client) GetDashboardInMyWorkspace(dashboardID string) (*Dashboard, error) {
	var respObj Dashboard
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/dashboards/%s", url.PathEscape(dashboardID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// DeleteDashboard deletes a dashboard from a workspace
func (client *Client) DeleteDashboard(groupID, dashboardID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s",
		url.PathEscape(groupID), url.PathEscape(dashboardID))
	return client.doJSON("DELETE", url, nil, nil)
}

// DeleteDashboardInMyWorkspace deletes a dashboard from My Workspace
func (client *Client) DeleteDashboardInMyWorkspace(dashboardID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/dashboards/%s", url.PathEscape(dashboardID))
	return client.doJSON("DELETE", url, nil, nil)
}

// GetTiles returns a list of tiles in a dashboard
func (client *Client) GetTiles(groupID, dashboardID string) (*GetTilesResponse, error) {
	var respObj GetTilesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s/tiles",
		url.PathEscape(groupID), url.PathEscape(dashboardID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetTilesInMyWorkspace returns a list of tiles in a dashboard from My Workspace
func (client *Client) GetTilesInMyWorkspace(dashboardID string) (*GetTilesResponse, error) {
	var respObj GetTilesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/dashboards/%s/tiles", url.PathEscape(dashboardID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetTile returns a specific tile from a dashboard
func (client *Client) GetTile(groupID, dashboardID, tileID string) (*Tile, error) {
	var respObj Tile
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s/tiles/%s",
		url.PathEscape(groupID), url.PathEscape(dashboardID), url.PathEscape(tileID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetTileInMyWorkspace returns a specific tile from a dashboard in My Workspace
func (client *Client) GetTileInMyWorkspace(dashboardID, tileID string) (*Tile, error) {
	var respObj Tile
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/dashboards/%s/tiles/%s",
		url.PathEscape(dashboardID), url.PathEscape(tileID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// CloneTile clones a tile to another dashboard
func (client *Client) CloneTile(groupID, dashboardID, tileID string, request CloneTileRequest) (*Tile, error) {
	var respObj Tile
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards/%s/tiles/%s/Clone",
		url.PathEscape(groupID), url.PathEscape(dashboardID), url.PathEscape(tileID))
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// CloneTileInMyWorkspace clones a tile to another dashboard in My Workspace
func (client *Client) CloneTileInMyWorkspace(dashboardID, tileID string, request CloneTileRequest) (*Tile, error) {
	var respObj Tile
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/dashboards/%s/tiles/%s/Clone",
		url.PathEscape(dashboardID), url.PathEscape(tileID))
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}