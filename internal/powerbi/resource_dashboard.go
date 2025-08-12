package powerbi

import (
	"fmt"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceDashboard represents a Power BI dashboard
func ResourceDashboard() *schema.Resource {
	return &schema.Resource{
		Create: createDashboard,
		Read:   readDashboard,
		Delete: deleteDashboard,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid dashboard import id format, expected 'workspace_id/dashboard_id'")
				}
				d.Set("workspace_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the workspace where the dashboard will be created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "Name of the dashboard.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the dashboard.",
			},
			"is_read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the dashboard is read-only.",
			},
			"web_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Web URL of the dashboard.",
			},
			"embed_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Embed URL of the dashboard.",
			},
		},
	}
}

func createDashboard(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	name := d.Get("name").(string)
	
	dashboard, err := client.CreateDashboard(workspaceID, powerbiapi.CreateDashboardRequest{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to create dashboard: %w", err)
	}
	
	d.SetId(dashboard.ID)
	
	return readDashboard(d, meta)
}

func readDashboard(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dashboardID := d.Id()
	
	if workspaceID == "" {
		// If workspace_id is not set (e.g., during import), try to find it
		// This would require listing all workspaces and searching for the dashboard
		// For now, we'll return an error
		return fmt.Errorf("workspace_id is required to read dashboard")
	}
	
	dashboard, err := client.GetDashboard(workspaceID, dashboardID)
	if err != nil {
		// Check if dashboard was deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read dashboard: %w", err)
	}
	
	d.Set("name", dashboard.DisplayName)
	d.Set("display_name", dashboard.DisplayName)
	d.Set("is_read_only", dashboard.IsReadOnly)
	d.Set("web_url", dashboard.WebURL)
	d.Set("embed_url", dashboard.EmbedURL)
	
	return nil
}

func deleteDashboard(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dashboardID := d.Id()
	
	err := client.DeleteDashboard(workspaceID, dashboardID)
	if err != nil {
		// Ignore 404 errors - dashboard already deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}
	
	return nil
}