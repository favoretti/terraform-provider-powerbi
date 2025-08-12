package powerbiapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	ODataContext string          `json:"@odata.context,omitempty"`
	ODataCount   int             `json:"@odata.count,omitempty"`
	ODataNextLink string         `json:"@odata.nextLink,omitempty"`
	Value        json.RawMessage `json:"value"`
}

// PaginationOptions represents options for paginated requests
type PaginationOptions struct {
	Top     int    // Number of items to return per page
	Skip    int    // Number of items to skip
	Filter  string // OData filter expression
	OrderBy string // OData orderby expression
	Select  string // OData select expression
	Expand  string // OData expand expression
}

// BuildPaginationQuery builds query parameters for pagination
func BuildPaginationQuery(options *PaginationOptions) string {
	if options == nil {
		return ""
	}

	params := url.Values{}
	
	if options.Top > 0 {
		params.Add("$top", strconv.Itoa(options.Top))
	}
	
	if options.Skip > 0 {
		params.Add("$skip", strconv.Itoa(options.Skip))
	}
	
	if options.Filter != "" {
		params.Add("$filter", options.Filter)
	}
	
	if options.OrderBy != "" {
		params.Add("$orderby", options.OrderBy)
	}
	
	if options.Select != "" {
		params.Add("$select", options.Select)
	}
	
	if options.Expand != "" {
		params.Add("$expand", options.Expand)
	}
	
	return params.Encode()
}

// GetAllPages retrieves all pages of a paginated response
func (client *Client) GetAllPages(initialURL string, result interface{}) error {
	allItems := make([]json.RawMessage, 0)
	nextURL := initialURL
	
	for nextURL != "" {
		var paginatedResp PaginatedResponse
		err := client.doJSON("GET", nextURL, nil, &paginatedResp)
		if err != nil {
			return err
		}
		
		// Parse the value array
		var items []json.RawMessage
		if err := json.Unmarshal(paginatedResp.Value, &items); err != nil {
			return fmt.Errorf("failed to parse paginated response: %w", err)
		}
		
		allItems = append(allItems, items...)
		
		// Check for next page
		if paginatedResp.ODataNextLink != "" {
			nextURL = paginatedResp.ODataNextLink
		} else {
			nextURL = ""
		}
	}
	
	// Marshal all items back into the result
	allItemsJSON, err := json.Marshal(map[string]interface{}{
		"value": allItems,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal paginated results: %w", err)
	}
	
	return json.Unmarshal(allItemsJSON, result)
}

// Enhanced list operations with pagination support

// GetGroupsWithPagination returns groups with pagination support
func (client *Client) GetGroupsWithPagination(options *PaginationOptions) (*GetGroupsResponse, error) {
	baseURL := "https://api.powerbi.com/v1.0/myorg/groups"
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetGroupsResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}

// GetDatasetsInGroupWithPagination returns datasets with pagination support
func (client *Client) GetDatasetsInGroupWithPagination(groupID string, options *PaginationOptions) (*GetDatasetsInGroupResponse, error) {
	baseURL := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets", url.PathEscape(groupID))
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetDatasetsInGroupResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}

// GetReportsInGroupWithPagination returns reports with pagination support
func (client *Client) GetReportsInGroupWithPagination(groupID string, options *PaginationOptions) (*GetReportsInGroupResponse, error) {
	baseURL := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/reports", url.PathEscape(groupID))
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetReportsInGroupResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}

// GetDashboardsWithPagination returns dashboards with pagination support
func (client *Client) GetDashboardsWithPagination(groupID string, options *PaginationOptions) (*GetDashboardsResponse, error) {
	baseURL := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dashboards", url.PathEscape(groupID))
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetDashboardsResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}

// GetDataflowsWithPagination returns dataflows with pagination support
func (client *Client) GetDataflowsWithPagination(groupID string, options *PaginationOptions) (*GetDataflowsResponse, error) {
	baseURL := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/dataflows", url.PathEscape(groupID))
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetDataflowsResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}

// GetPipelinesWithPagination returns pipelines with pagination support
func (client *Client) GetPipelinesWithPagination(options *PaginationOptions) (*GetPipelinesResponse, error) {
	baseURL := "https://api.powerbi.com/v1.0/myorg/pipelines"
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetPipelinesResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}

// GetGatewaysWithPagination returns gateways with pagination support
func (client *Client) GetGatewaysWithPagination(options *PaginationOptions) (*GetGatewaysResponse, error) {
	baseURL := "https://api.powerbi.com/v1.0/myorg/gateways"
	query := BuildPaginationQuery(options)
	
	url := baseURL
	if query != "" {
		url = fmt.Sprintf("%s?%s", baseURL, query)
	}
	
	var respObj GetGatewaysResponse
	err := client.GetAllPages(url, &respObj)
	return &respObj, err
}