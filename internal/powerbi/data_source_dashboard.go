package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceDashboard returns a specific dashboard from a workspace
func DataSourceDashboard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDashboardRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the workspace containing the dashboard.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the dashboard.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the dashboard.",
				ExactlyOneOf: []string{"id", "name"},
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

func dataSourceDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	
	var dashboard *powerbiapi.Dashboard
	var err error
	
	if dashboardID, ok := d.GetOk("id"); ok {
		// Get dashboard by ID
		dashboard, err = client.GetDashboard(workspaceID, dashboardID.(string))
		if err != nil {
			return fmt.Errorf("failed to get dashboard by ID %s: %w", dashboardID, err)
		}
	} else if dashboardName, ok := d.GetOk("name"); ok {
		// Get dashboard by name - need to list all dashboards and find by name
		dashboards, err := client.GetDashboards(workspaceID)
		if err != nil {
			return fmt.Errorf("failed to list dashboards: %w", err)
		}
		
		var foundDashboard *powerbiapi.Dashboard
		for _, d := range dashboards.Value {
			if d.DisplayName == dashboardName.(string) {
				foundDashboard = &d
				break
			}
		}
		
		if foundDashboard == nil {
			return fmt.Errorf("dashboard with name '%s' not found in workspace %s", dashboardName, workspaceID)
		}
		
		dashboard = foundDashboard
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}
	
	d.SetId(dashboard.ID)
	d.Set("id", dashboard.ID)
	d.Set("name", dashboard.DisplayName)
	d.Set("display_name", dashboard.DisplayName)
	d.Set("is_read_only", dashboard.IsReadOnly)
	d.Set("web_url", dashboard.WebURL)
	d.Set("embed_url", dashboard.EmbedURL)
	
	return nil
}