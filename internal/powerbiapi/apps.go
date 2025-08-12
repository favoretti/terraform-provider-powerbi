package powerbiapi

import (
	"fmt"
	"net/url"
)

// App represents a Power BI app
type App struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	PublishedBy       string `json:"publishedBy,omitempty"`
	LastUpdate        string `json:"lastUpdate,omitempty"`
}

// GetAppsResponse represents the response from GetApps API
type GetAppsResponse struct {
	Value []App `json:"value"`
}

// AppDashboard represents a dashboard within an app
type AppDashboard struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	EmbedURL    string `json:"embedUrl"`
	IsReadOnly  bool   `json:"isReadOnly"`
	WebURL      string `json:"webUrl"`
}

// GetAppDashboardsResponse represents the response from GetAppDashboards API
type GetAppDashboardsResponse struct {
	Value []AppDashboard `json:"value"`
}

// AppReport represents a report within an app
type AppReport struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	WebURL    string `json:"webUrl"`
	EmbedURL  string `json:"embedUrl"`
	DatasetID string `json:"datasetId"`
}

// GetAppReportsResponse represents the response from GetAppReports API
type GetAppReportsResponse struct {
	Value []AppReport `json:"value"`
}

// AppTile represents a tile within an app dashboard
type AppTile struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	SubTitle      string `json:"subTitle"`
	EmbedURL      string `json:"embedUrl"`
	EmbedData     string `json:"embedData"`
	RowSpan       int    `json:"rowSpan"`
	ColSpan       int    `json:"colSpan"`
	ReportID      string `json:"reportId"`
	DatasetID     string `json:"datasetId"`
}

// GetAppTilesResponse represents the response from GetAppTiles API
type GetAppTilesResponse struct {
	Value []AppTile `json:"value"`
}

// GetApps returns a list of installed apps
func (client *Client) GetApps() (*GetAppsResponse, error) {
	var respObj GetAppsResponse
	url := "https://api.powerbi.com/v1.0/myorg/apps"
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetApp returns a specific installed app
func (client *Client) GetApp(appID string) (*App, error) {
	var respObj App
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s", url.PathEscape(appID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetAppDashboards returns a list of dashboards from an app
func (client *Client) GetAppDashboards(appID string) (*GetAppDashboardsResponse, error) {
	var respObj GetAppDashboardsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s/dashboards", url.PathEscape(appID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetAppDashboard returns a specific dashboard from an app
func (client *Client) GetAppDashboard(appID, dashboardID string) (*AppDashboard, error) {
	var respObj AppDashboard
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s/dashboards/%s", 
		url.PathEscape(appID), url.PathEscape(dashboardID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetAppReports returns a list of reports from an app
func (client *Client) GetAppReports(appID string) (*GetAppReportsResponse, error) {
	var respObj GetAppReportsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s/reports", url.PathEscape(appID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetAppReport returns a specific report from an app
func (client *Client) GetAppReport(appID, reportID string) (*AppReport, error) {
	var respObj AppReport
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s/reports/%s", 
		url.PathEscape(appID), url.PathEscape(reportID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetAppTiles returns a list of tiles from an app dashboard
func (client *Client) GetAppTiles(appID, dashboardID string) (*GetAppTilesResponse, error) {
	var respObj GetAppTilesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s/dashboards/%s/tiles", 
		url.PathEscape(appID), url.PathEscape(dashboardID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetAppTile returns a specific tile from an app dashboard
func (client *Client) GetAppTile(appID, dashboardID, tileID string) (*AppTile, error) {
	var respObj AppTile
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/apps/%s/dashboards/%s/tiles/%s", 
		url.PathEscape(appID), url.PathEscape(dashboardID), url.PathEscape(tileID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}