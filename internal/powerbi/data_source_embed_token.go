package powerbi

import (
	"fmt"
	"time"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// DataSourceEmbedToken generates embed tokens for Power BI content
func DataSourceEmbedToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmbedTokenRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the workspace containing the content.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of content to generate token for.",
				ValidateFunc: validation.StringInSlice([]string{
					"report", "dataset", "dashboard", "tile",
				}, false),
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the resource (report, dataset, dashboard, or tile).",
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dashboard ID (required when type is 'tile').",
			},
			"access_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "View",
				Description: "Access level for the token.",
				ValidateFunc: validation.StringInSlice([]string{
					"View", "Edit", "Create",
				}, false),
			},
			"dataset_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of dataset IDs to include in the token.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"report_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of report IDs to include in the token.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_workspaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of target workspace IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The generated embed token.",
			},
			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The token ID.",
			},
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Token expiration date and time.",
			},
			"expires_on": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Token expiration as Unix timestamp.",
			},
		},
	}
}

func dataSourceEmbedTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	workspaceID := d.Get("workspace_id").(string)
	resourceType := d.Get("type").(string)
	resourceID := d.Get("resource_id").(string)
	accessLevel := d.Get("access_level").(string)

	// Build the request
	request := powerbiapi.GenerateTokenRequest{
		AccessLevel: accessLevel,
	}

	// Add dataset IDs if provided
	if datasetIdsRaw, ok := d.GetOk("dataset_ids"); ok {
		datasetIds := make([]string, 0)
		for _, id := range datasetIdsRaw.([]interface{}) {
			datasetIds = append(datasetIds, id.(string))
		}
		request.DatasetIds = datasetIds
	}

	// Add report IDs if provided
	if reportIdsRaw, ok := d.GetOk("report_ids"); ok {
		reportIds := make([]string, 0)
		for _, id := range reportIdsRaw.([]interface{}) {
			reportIds = append(reportIds, id.(string))
		}
		request.ReportIds = reportIds
	}

	// Add target workspaces if provided
	if targetWorkspacesRaw, ok := d.GetOk("target_workspaces"); ok {
		targetWorkspaces := make([]powerbiapi.TargetWorkspace, 0)
		for _, id := range targetWorkspacesRaw.([]interface{}) {
			targetWorkspaces = append(targetWorkspaces, powerbiapi.TargetWorkspace{
				ID: id.(string),
			})
		}
		request.TargetWorkspaces = targetWorkspaces
	}

	var response *powerbiapi.GenerateTokenResponse
	var err error

	// Generate token based on type
	switch resourceType {
	case "report":
		response, err = client.GenerateEmbedTokenForReport(workspaceID, resourceID, request)
	case "dataset":
		response, err = client.GenerateEmbedTokenForDataset(workspaceID, resourceID, request)
	case "dashboard":
		response, err = client.GenerateEmbedTokenForDashboard(workspaceID, resourceID, request)
	case "tile":
		dashboardID, ok := d.GetOk("dashboard_id")
		if !ok {
			return fmt.Errorf("dashboard_id is required when type is 'tile'")
		}
		response, err = client.GenerateEmbedTokenForTile(workspaceID, dashboardID.(string), resourceID, request)
	default:
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	if err != nil {
		return fmt.Errorf("failed to generate embed token for %s %s: %w", resourceType, resourceID, err)
	}

	// Parse expiration time
	expiration, err := time.Parse(time.RFC3339, response.Expiration)
	if err != nil {
		return fmt.Errorf("failed to parse expiration time: %w", err)
	}

	// Set the resource ID as a combination of workspace, type, and resource
	d.SetId(fmt.Sprintf("%s/%s/%s", workspaceID, resourceType, resourceID))
	d.Set("token", response.Token)
	d.Set("token_id", response.TokenID)
	d.Set("expiration", response.Expiration)
	d.Set("expires_on", expiration.Unix())

	return nil
}