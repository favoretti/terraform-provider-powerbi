package powerbiapi

import (
	"fmt"
	"net/url"
)

// Gateway represents a Power BI gateway
type Gateway struct {
	ID                       string           `json:"id"`
	Name                     string           `json:"name"`
	Type                     string           `json:"type"`
	PublicKey                GatewayPublicKey `json:"publicKey"`
	GatewayAnnotation        string           `json:"gatewayAnnotation,omitempty"`
	GatewayStatus            string           `json:"gatewayStatus,omitempty"`
	GatewayVersion           string           `json:"gatewayVersion,omitempty"`
	GatewayMachine           string           `json:"gatewayMachine,omitempty"`
	GatewayContactInfo       []string         `json:"gatewayContactInfo,omitempty"`
	GatewayClusterId         string           `json:"gatewayClusterId,omitempty"`
	GatewayClusterStatus     string           `json:"gatewayClusterStatus,omitempty"`
}

// GatewayPublicKey represents the public key of a gateway
type GatewayPublicKey struct {
	Exponent string `json:"exponent"`
	Modulus  string `json:"modulus"`
}

// GetGatewaysResponse represents the response from GetGateways API
type GetGatewaysResponse struct {
	Value []Gateway `json:"value"`
}

// GatewayDatasource represents a data source in a gateway
type GatewayDatasource struct {
	ID                        string                     `json:"id"`
	GatewayID                 string                     `json:"gatewayId"`
	DatasourceName            string                     `json:"datasourceName"`
	DatasourceType            string                     `json:"datasourceType"`
	ConnectionString          string                     `json:"connectionString,omitempty"`
	CredentialType            string                     `json:"credentialType"`
	CredentialDetails         *DatasourceCredentialDetails `json:"credentialDetails,omitempty"`
	ConnectionDetails         DatasourceConnectionDetails `json:"connectionDetails"`
}

// DatasourceConnectionDetails represents connection details for a datasource
type DatasourceConnectionDetails struct {
	Server          string `json:"server,omitempty"`
	Database        string `json:"database,omitempty"`
	URL             string `json:"url,omitempty"`
	Path            string `json:"path,omitempty"`
	Kind            string `json:"kind,omitempty"`
	AuthMethod      string `json:"authMethod,omitempty"`
	Account         string `json:"account,omitempty"`
	Domain          string `json:"domain,omitempty"`
	EmailAddress    string `json:"emailAddress,omitempty"`
	LoginServer     string `json:"loginServer,omitempty"`
	Class           string `json:"class,omitempty"`
	AdvancedSettings interface{} `json:"advancedSettings,omitempty"`
}

// DatasourceCredentialDetails represents credential details for a datasource
type DatasourceCredentialDetails struct {
	CredentialType              string                      `json:"credentialType"`
	Credentials                 string                      `json:"credentials,omitempty"`
	EncryptedConnection         string                      `json:"encryptedConnection,omitempty"`
	EncryptionAlgorithm         string                      `json:"encryptionAlgorithm,omitempty"`
	PrivacyLevel                string                      `json:"privacyLevel,omitempty"`
	UseCallerAADIdentity        bool                        `json:"useCallerAADIdentity,omitempty"`
	UseEndUserOAuth2Credentials bool                        `json:"useEndUserOAuth2Credentials,omitempty"`
}

// CreateDatasourceRequest represents the request for creating a datasource
type CreateDatasourceRequest struct {
	DatasourceName    string                       `json:"datasourceName"`
	DatasourceType    string                       `json:"datasourceType"`
	ConnectionDetails DatasourceConnectionDetails  `json:"connectionDetails"`
	CredentialDetails *DatasourceCredentialDetails `json:"credentialDetails,omitempty"`
}

// UpdateDatasourceRequest represents the request for updating a datasource
type UpdateDatasourceRequest struct {
	CredentialDetails *DatasourceCredentialDetails `json:"credentialDetails"`
}

// GetDatasourcesResponse represents the response from GetDatasources API
type GetDatasourcesResponse struct {
	Value []GatewayDatasource `json:"value"`
}

// DatasourceUser represents a user with access to a datasource
type DatasourceUser struct {
	DatasourceAccessRight string `json:"datasourceAccessRight"`
	DisplayName           string `json:"displayName"`
	EmailAddress          string `json:"emailAddress"`
	GraphID               string `json:"graphId"`
	Identifier            string `json:"identifier"`
	PrincipalType         string `json:"principalType"`
	Profile               *UserProfile `json:"profile,omitempty"`
}

// UserProfile represents a user profile
type UserProfile struct {
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// GetDatasourceUsersResponse represents the response from GetDatasourceUsers API
type GetDatasourceUsersResponse struct {
	Value []DatasourceUser `json:"value"`
}

// AddDatasourceUserRequest represents the request for adding a datasource user
type AddDatasourceUserRequest struct {
	DatasourceAccessRight string `json:"datasourceAccessRight"`
	EmailAddress          string `json:"emailAddress,omitempty"`
	DisplayName           string `json:"displayName,omitempty"`
	Identifier            string `json:"identifier,omitempty"`
	GraphID               string `json:"graphId,omitempty"`
	PrincipalType         string `json:"principalType,omitempty"`
}

// DatasourceStatus represents the status of a datasource
type DatasourceStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// GetGateways returns a list of gateways
func (client *Client) GetGateways() (*GetGatewaysResponse, error) {
	var respObj GetGatewaysResponse
	url := "https://api.powerbi.com/v1.0/myorg/gateways"
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetGateway returns a specific gateway
func (client *Client) GetGateway(gatewayID string) (*Gateway, error) {
	var respObj Gateway
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s", url.PathEscape(gatewayID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// CreateDatasource creates a new datasource in a gateway
func (client *Client) CreateDatasource(gatewayID string, request CreateDatasourceRequest) (*GatewayDatasource, error) {
	var respObj GatewayDatasource
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources", url.PathEscape(gatewayID))
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// GetDatasources returns a list of datasources in a gateway
func (client *Client) GetDatasources(gatewayID string) (*GetDatasourcesResponse, error) {
	var respObj GetDatasourcesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources", url.PathEscape(gatewayID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetDatasource returns a specific datasource
func (client *Client) GetDatasource(gatewayID, datasourceID string) (*GatewayDatasource, error) {
	var respObj GatewayDatasource
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// UpdateDatasource updates a datasource
func (client *Client) UpdateDatasource(gatewayID, datasourceID string, request UpdateDatasourceRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	return client.doJSON("PATCH", url, request, nil)
}

// DeleteDatasource deletes a datasource
func (client *Client) DeleteDatasource(gatewayID, datasourceID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	return client.doJSON("DELETE", url, nil, nil)
}

// GetDatasourceStatus returns the status of a datasource
func (client *Client) GetDatasourceStatus(gatewayID, datasourceID string) (*DatasourceStatus, error) {
	var respObj DatasourceStatus
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/status",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetDatasourceUsers returns a list of users with access to a datasource
func (client *Client) GetDatasourceUsers(gatewayID, datasourceID string) (*GetDatasourceUsersResponse, error) {
	var respObj GetDatasourceUsersResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/users",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// AddDatasourceUser adds a user to a datasource
func (client *Client) AddDatasourceUser(gatewayID, datasourceID string, request AddDatasourceUserRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/users",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	return client.doJSON("POST", url, request, nil)
}

// DeleteDatasourceUser removes a user from a datasource
func (client *Client) DeleteDatasourceUser(gatewayID, datasourceID, userID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s/users/%s",
		url.PathEscape(gatewayID), url.PathEscape(datasourceID), url.PathEscape(userID))
	return client.doJSON("DELETE", url, nil, nil)
}