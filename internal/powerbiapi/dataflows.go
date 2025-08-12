package powerbiapi

import (
	"fmt"
	"net/url"
	"time"
)

// Dataflow represents a Power BI dataflow
type Dataflow struct {
	ObjectID              string    `json:"objectId"`
	Name                  string    `json:"name"`
	Description           string    `json:"description,omitempty"`
	ModelURL              string    `json:"modelUrl,omitempty"`
	ConfiguredBy          string    `json:"configuredBy,omitempty"`
	ModifiedBy            string    `json:"modifiedBy,omitempty"`
	ModifiedDateTime      time.Time `json:"modifiedDateTime,omitempty"`
	Users                 []DataflowUser `json:"users,omitempty"`
	RefreshSchedule       *DataflowRefreshSchedule `json:"refreshSchedule,omitempty"`
}

// DataflowUser represents a user with access to a dataflow
type DataflowUser struct {
	Identifier    string `json:"identifier"`
	DisplayName   string `json:"displayName"`
	EmailAddress  string `json:"emailAddress"`
	GraphID       string `json:"graphId"`
	PrincipalType string `json:"principalType"`
	UserType      string `json:"userType"`
}

// GetDataflowsResponse represents the response from GetDataflows API
type GetDataflowsResponse struct {
	Value []Dataflow `json:"value"`
}

// DataflowDatasource represents a datasource in a dataflow
type DataflowDatasource struct {
	DatasourceID      string `json:"datasourceId"`
	DatasourceType    string `json:"datasourceType"`
	GatewayID         string `json:"gatewayId,omitempty"`
	ConnectionDetails interface{} `json:"connectionDetails"`
}

// GetDataflowDatasourcesResponse represents the response from GetDataflowDatasources API
type GetDataflowDatasourcesResponse struct {
	Value []DataflowDatasource `json:"value"`
}

// DataflowRefreshSchedule represents a refresh schedule for a dataflow
type DataflowRefreshSchedule struct {
	Days                        []string    `json:"days,omitempty"`
	Times                       []string    `json:"times,omitempty"`
	Enabled                     bool        `json:"enabled"`
	LocalTimeZoneID             string      `json:"localTimeZoneId,omitempty"`
	NotifyOption                string      `json:"notifyOption,omitempty"`
}

// UpdateDataflowRefreshScheduleRequest represents the request for updating a dataflow refresh schedule
type UpdateDataflowRefreshScheduleRequest struct {
	Value DataflowRefreshSchedule `json:"value"`
}

// DataflowTransaction represents a dataflow refresh transaction
type DataflowTransaction struct {
	ID               string    `json:"id"`
	RefreshType      string    `json:"refreshType"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime,omitempty"`
	Status           string    `json:"status"`
	ErrorCode        string    `json:"errorCode,omitempty"`
}

// GetDataflowTransactionsResponse represents the response from GetDataflowTransactions API
type GetDataflowTransactionsResponse struct {
	Value []DataflowTransaction `json:"value"`
}

// DataflowUpstreamDataflow represents an upstream dataflow dependency
type DataflowUpstreamDataflow struct {
	TargetDataflowID    string `json:"targetDataflowId"`
	GroupID             string `json:"groupId"`
}

// GetUpstreamDataflowsResponse represents the response from GetUpstreamDataflows API
type GetUpstreamDataflowsResponse struct {
	Value []DataflowUpstreamDataflow `json:"value"`
}

// RefreshDataflowRequest represents the request for refreshing a dataflow
type RefreshDataflowRequest struct {
	NotifyOption string `json:"notifyOption,omitempty"`
}

// UpdateDataflowRequest represents the request for updating a dataflow
type UpdateDataflowRequest struct {
	Name              string    `json:"name,omitempty"`
	Description       string    `json:"description,omitempty"`
	AllowNativeQueries bool     `json:"allowNativeQueries,omitempty"`
	ComputeEngineSettings interface{} `json:"computeEngineSettings,omitempty"`
}

// CreateDataflowRequest represents the request for creating a dataflow
type CreateDataflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Definition  string `json:"definition"`
}

// CreateDataflow creates a new dataflow in a workspace
func (client *Client) CreateDataflow(groupID string, request CreateDataflowRequest) (*Dataflow, error) {
	var respObj Dataflow
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows", url.PathEscape(groupID))
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// GetDataflows returns a list of dataflows in a workspace
func (client *Client) GetDataflows(groupID string) (*GetDataflowsResponse, error) {
	var respObj GetDataflowsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows", url.PathEscape(groupID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetDataflow returns a specific dataflow
func (client *Client) GetDataflow(groupID, dataflowID string) (*Dataflow, error) {
	var respObj Dataflow
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// UpdateDataflow updates a dataflow
func (client *Client) UpdateDataflow(groupID, dataflowID string, request UpdateDataflowRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	return client.doJSON("PATCH", url, request, nil)
}

// DeleteDataflow deletes a dataflow
func (client *Client) DeleteDataflow(groupID, dataflowID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	return client.doJSON("DELETE", url, nil, nil)
}

// GetDataflowDatasources returns datasources for a dataflow
func (client *Client) GetDataflowDatasources(groupID, dataflowID string) (*GetDataflowDatasourcesResponse, error) {
	var respObj GetDataflowDatasourcesResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/datasources",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// RefreshDataflow triggers a refresh for a dataflow
func (client *Client) RefreshDataflow(groupID, dataflowID string, request RefreshDataflowRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/refreshes",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	return client.doJSON("POST", url, request, nil)
}

// GetDataflowRefreshSchedule returns the refresh schedule for a dataflow
func (client *Client) GetDataflowRefreshSchedule(groupID, dataflowID string) (*DataflowRefreshSchedule, error) {
	var respObj DataflowRefreshSchedule
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/refreshSchedule",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// UpdateDataflowRefreshSchedule updates the refresh schedule for a dataflow
func (client *Client) UpdateDataflowRefreshSchedule(groupID, dataflowID string, request UpdateDataflowRefreshScheduleRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/refreshSchedule",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	return client.doJSON("PATCH", url, request, nil)
}

// GetDataflowTransactions returns transactions for a dataflow
func (client *Client) GetDataflowTransactions(groupID, dataflowID string) (*GetDataflowTransactionsResponse, error) {
	var respObj GetDataflowTransactionsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/transactions",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// CancelDataflowTransaction cancels a dataflow transaction
func (client *Client) CancelDataflowTransaction(groupID, dataflowID, transactionID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/transactions/%s/cancel",
		url.PathEscape(groupID), url.PathEscape(dataflowID), url.PathEscape(transactionID))
	return client.doJSON("POST", url, nil, nil)
}

// GetUpstreamDataflows returns upstream dataflows for a dataflow
func (client *Client) GetUpstreamDataflows(groupID, dataflowID string) (*GetUpstreamDataflowsResponse, error) {
	var respObj GetUpstreamDataflowsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows/%s/upstreamDataflows",
		url.PathEscape(groupID), url.PathEscape(dataflowID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}