package powerbiapi

import (
	"fmt"
	"net/url"
	"time"
)

// Pipeline represents a Power BI deployment pipeline
type Pipeline struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"displayName"`
	Description string          `json:"description,omitempty"`
	Stages      []PipelineStage `json:"stages,omitempty"`
	Users       []PipelineUser  `json:"users,omitempty"`
}

// PipelineStage represents a stage in a deployment pipeline
type PipelineStage struct {
	Order               int                  `json:"order"`
	StageName           string               `json:"stageName"`
	IsPublic            bool                 `json:"isPublic"`
	WorkspaceID         string               `json:"workspaceId,omitempty"`
	WorkspaceName       string               `json:"workspaceName,omitempty"`
	ArtifactsCount      int                  `json:"artifactsCount,omitempty"`
}

// PipelineUser represents a user with access to a pipeline
type PipelineUser struct {
	Identifier     string `json:"identifier"`
	AccessRight    string `json:"accessRight"`
	PrincipalType  string `json:"principalType"`
	DisplayName    string `json:"displayName,omitempty"`
	EmailAddress   string `json:"emailAddress,omitempty"`
	GraphID        string `json:"graphId,omitempty"`
	UserType       string `json:"userType,omitempty"`
}

// GetPipelinesResponse represents the response from GetPipelines API
type GetPipelinesResponse struct {
	Value []Pipeline `json:"value"`
}

// CreatePipelineRequest represents the request for creating a pipeline
type CreatePipelineRequest struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`
}

// UpdatePipelineRequest represents the request for updating a pipeline
type UpdatePipelineRequest struct {
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// AssignWorkspaceRequest represents the request for assigning a workspace to a pipeline stage
type AssignWorkspaceRequest struct {
	WorkspaceID string `json:"workspaceId"`
}

// UnassignWorkspaceRequest represents the request for unassigning a workspace from a pipeline stage
type UnassignWorkspaceRequest struct {
	WorkspaceID string `json:"workspaceId"`
}

// DeployRequest represents the request for deploying content between stages
type DeployRequest struct {
	SourceStageOrder   int                   `json:"sourceStageOrder"`
	ArtifactsToDeploy  []DeployArtifact      `json:"artifactsToDeploy,omitempty"`
	Options            *DeployOptions        `json:"options,omitempty"`
	Note               string                `json:"note,omitempty"`
}

// DeployArtifact represents an artifact to deploy
type DeployArtifact struct {
	ArtifactID   string `json:"artifactId"`
	ArtifactType string `json:"artifactType"`
}

// DeployOptions represents deployment options
type DeployOptions struct {
	AllowCreateArtifact         bool                       `json:"allowCreateArtifact,omitempty"`
	AllowOverwriteArtifact      bool                       `json:"allowOverwriteArtifact,omitempty"`
	AllowOverwriteTargetSchema  bool                       `json:"allowOverwriteTargetSchema,omitempty"`
	AllowPurgeData              bool                       `json:"allowPurgeData,omitempty"`
	AllowSkipTilesWithMissingPrerequisites bool             `json:"allowSkipTilesWithMissingPrerequisites,omitempty"`
	AllowTakeOver               bool                       `json:"allowTakeOver,omitempty"`
}

// DeployResponse represents the response from a deployment
type DeployResponse struct {
	ID string `json:"id"`
}

// PipelineOperation represents a pipeline operation
type PipelineOperation struct {
	ID                  string              `json:"id"`
	Type                string              `json:"type"`
	Status              string              `json:"status"`
	LastUpdatedTime     time.Time           `json:"lastUpdatedTime"`
	ExecutionStartTime  time.Time           `json:"executionStartTime,omitempty"`
	ExecutionEndTime    time.Time           `json:"executionEndTime,omitempty"`
	SourceStageOrder    int                 `json:"sourceStageOrder,omitempty"`
	TargetStageOrder    int                 `json:"targetStageOrder,omitempty"`
	PreDeploymentDetails interface{}        `json:"preDeploymentDetails,omitempty"`
	Note                string              `json:"note,omitempty"`
	Error               *PipelineOperationError `json:"error,omitempty"`
}

// PipelineOperationError represents an error in a pipeline operation
type PipelineOperationError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorDetails string `json:"errorDetails"`
}

// GetPipelineOperationsResponse represents the response from GetPipelineOperations API
type GetPipelineOperationsResponse struct {
	Value []PipelineOperation `json:"value"`
}

// PipelineStageArtifact represents an artifact in a pipeline stage
type PipelineStageArtifact struct {
	ArtifactID          string    `json:"artifactId"`
	ArtifactType        string    `json:"artifactType"`
	ArtifactDisplayName string    `json:"artifactDisplayName"`
	SourceArtifactID    string    `json:"sourceArtifactId,omitempty"`
	TargetArtifactID    string    `json:"targetArtifactId,omitempty"`
	LastDeploymentTime  time.Time `json:"lastDeploymentTime,omitempty"`
}

// GetPipelineStageArtifactsResponse represents the response from GetPipelineStageArtifacts API
type GetPipelineStageArtifactsResponse struct {
	Value []PipelineStageArtifact `json:"value"`
}

// AddPipelineUserRequest represents the request for adding a user to a pipeline
type AddPipelineUserRequest struct {
	Identifier    string `json:"identifier"`
	AccessRight   string `json:"accessRight"`
	PrincipalType string `json:"principalType,omitempty"`
}

// UpdatePipelineUserRequest represents the request for updating a user's pipeline access
type UpdatePipelineUserRequest struct {
	AccessRight string `json:"accessRight"`
}

// CreatePipeline creates a new deployment pipeline
func (client *Client) CreatePipeline(request CreatePipelineRequest) (*Pipeline, error) {
	var respObj Pipeline
	url := "https://api.powerbi.com/v1.0/myorg/pipelines"
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// GetPipelines returns a list of deployment pipelines
func (client *Client) GetPipelines() (*GetPipelinesResponse, error) {
	var respObj GetPipelinesResponse
	url := "https://api.powerbi.com/v1.0/myorg/pipelines"
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetPipeline returns a specific deployment pipeline
func (client *Client) GetPipeline(pipelineID string) (*Pipeline, error) {
	var respObj Pipeline
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s", url.PathEscape(pipelineID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// UpdatePipeline updates a deployment pipeline
func (client *Client) UpdatePipeline(pipelineID string, request UpdatePipelineRequest) (*Pipeline, error) {
	var respObj Pipeline
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s", url.PathEscape(pipelineID))
	err := client.doJSON("PATCH", url, request, &respObj)
	return &respObj, err
}

// DeletePipeline deletes a deployment pipeline
func (client *Client) DeletePipeline(pipelineID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s", url.PathEscape(pipelineID))
	return client.doJSON("DELETE", url, nil, nil)
}

// GetPipelineStages returns the stages of a deployment pipeline
func (client *Client) GetPipelineStages(pipelineID string) ([]PipelineStage, error) {
	pipeline, err := client.GetPipeline(pipelineID)
	if err != nil {
		return nil, err
	}
	return pipeline.Stages, nil
}

// AssignWorkspace assigns a workspace to a pipeline stage
func (client *Client) AssignWorkspace(pipelineID string, stageOrder int, request AssignWorkspaceRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/stages/%d/assignWorkspace",
		url.PathEscape(pipelineID), stageOrder)
	return client.doJSON("POST", url, request, nil)
}

// UnassignWorkspace unassigns a workspace from a pipeline stage
func (client *Client) UnassignWorkspace(pipelineID string, stageOrder int, request UnassignWorkspaceRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/stages/%d/unassignWorkspace",
		url.PathEscape(pipelineID), stageOrder)
	return client.doJSON("POST", url, request, nil)
}

// DeployAll deploys all content from source stage to target stage
func (client *Client) DeployAll(pipelineID string, request DeployRequest) (*DeployResponse, error) {
	var respObj DeployResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/deployAll",
		url.PathEscape(pipelineID))
	err := client.doJSON("POST", url, request, &respObj)
	return &respObj, err
}

// GetPipelineOperations returns operations for a deployment pipeline
func (client *Client) GetPipelineOperations(pipelineID string) (*GetPipelineOperationsResponse, error) {
	var respObj GetPipelineOperationsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/operations",
		url.PathEscape(pipelineID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetPipelineOperation returns a specific pipeline operation
func (client *Client) GetPipelineOperation(pipelineID, operationID string) (*PipelineOperation, error) {
	var respObj PipelineOperation
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/operations/%s",
		url.PathEscape(pipelineID), url.PathEscape(operationID))
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetPipelineStageArtifacts returns artifacts in a pipeline stage
func (client *Client) GetPipelineStageArtifacts(pipelineID string, stageOrder int) (*GetPipelineStageArtifactsResponse, error) {
	var respObj GetPipelineStageArtifactsResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/stages/%d/artifacts",
		url.PathEscape(pipelineID), stageOrder)
	err := client.doJSON("GET", url, nil, &respObj)
	return &respObj, err
}

// GetPipelineUsers returns users with access to a pipeline
func (client *Client) GetPipelineUsers(pipelineID string) ([]PipelineUser, error) {
	pipeline, err := client.GetPipeline(pipelineID)
	if err != nil {
		return nil, err
	}
	return pipeline.Users, nil
}

// AddPipelineUser adds a user to a pipeline
func (client *Client) AddPipelineUser(pipelineID string, request AddPipelineUserRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/users",
		url.PathEscape(pipelineID))
	return client.doJSON("POST", url, request, nil)
}

// UpdatePipelineUser updates a user's pipeline access
func (client *Client) UpdatePipelineUser(pipelineID, userID string, request UpdatePipelineUserRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/users/%s",
		url.PathEscape(pipelineID), url.PathEscape(userID))
	return client.doJSON("PATCH", url, request, nil)
}

// DeletePipelineUser removes a user from a pipeline
func (client *Client) DeletePipelineUser(pipelineID, userID string) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/pipelines/%s/users/%s",
		url.PathEscape(pipelineID), url.PathEscape(userID))
	return client.doJSON("DELETE", url, nil, nil)
}