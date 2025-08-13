package powerbiapi

import (
	"fmt"
	"net/url"
)

// TemplateApp represents a Power BI template app
type TemplateApp struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	PublisherName   string `json:"publisherName,omitempty"`
	PublisherEmail  string `json:"publisherEmail,omitempty"`
	SupportContact  string `json:"supportContact,omitempty"`
	Version         string `json:"version,omitempty"`
	LogoURL         string `json:"logoUrl,omitempty"`
	PackageURL      string `json:"packageUrl,omitempty"`
}

// GetTemplateAppsResponse represents the response from GetTemplateApps API
type GetTemplateAppsResponse struct {
	Value []TemplateApp `json:"value"`
}

// TemplateAppInstallation represents an installation of a template app
type TemplateAppInstallation struct {
	ID            string `json:"id"`
	TemplateAppID string `json:"templateAppId"`
	WorkspaceID   string `json:"workspaceId"`
	AppID         string `json:"appId,omitempty"`
	PackageKey    string `json:"packageKey,omitempty"`
	OwnerTenantID string `json:"ownerTenantId,omitempty"`
	Config        map[string]interface{} `json:"config,omitempty"`
}

// InstallTemplateAppRequest represents the request to install a template app
type InstallTemplateAppRequest struct {
	TemplateAppID string                 `json:"templateAppId"`
	WorkspaceID   string                 `json:"workspaceId"`
	PackageKey    string                 `json:"packageKey,omitempty"`
	Config        map[string]interface{} `json:"config,omitempty"`
}

// InstallTemplateAppResponse represents the response from installing a template app
type InstallTemplateAppResponse struct {
	ID string `json:"id"`
}

// GetTemplateApps returns a list of available template apps
func (client *Client) GetTemplateApps() (*GetTemplateAppsResponse, error) {
	var respObj GetTemplateAppsResponse
	url := "https://api.powerbi.com/v1.0/myorg/templateApps"
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetTemplateApp returns a specific template app
func (client *Client) GetTemplateApp(templateAppID string) (*TemplateApp, error) {
	var respObj TemplateApp
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/templateApps/%s", url.PathEscape(templateAppID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// InstallTemplateApp installs a template app to a workspace
func (client *Client) InstallTemplateApp(request InstallTemplateAppRequest) (*InstallTemplateAppResponse, error) {
	var respObj InstallTemplateAppResponse
	url := "https://api.powerbi.com/v1.0/myorg/templateApps/install"
	err := client.doJSON("POST", url, &request, &respObj)
	return &respObj, err
}

// GetTemplateAppInstallation returns details of a template app installation
func (client *Client) GetTemplateAppInstallation(installationID string) (*TemplateAppInstallation, error) {
	var respObj TemplateAppInstallation
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/templateApps/installations/%s", url.PathEscape(installationID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// UninstallTemplateApp uninstalls a template app
func (client *Client) UninstallTemplateApp(installationID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/templateApps/installations/%s", url.PathEscape(installationID))
	return client.doJSON("DELETE", url, nil, nil)
}