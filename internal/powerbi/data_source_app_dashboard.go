package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceAppDashboard returns dashboards from a Power BI app
func DataSourceAppDashboard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppDashboardRead,

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the app.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the dashboard.",
				ExactlyOneOf: []string{"id", "display_name"},
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display name of the dashboard.",
				ExactlyOneOf: []string{"id", "display_name"},
			},
			"embed_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Embed URL of the dashboard.",
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
			"tiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of tiles in the dashboard.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the tile.",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Title of the tile.",
						},
						"subtitle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subtitle of the tile.",
						},
						"embed_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Embed URL of the tile.",
						},
						"embed_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Embed data of the tile.",
						},
						"row_span": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of rows the tile spans.",
						},
						"col_span": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of columns the tile spans.",
						},
						"report_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Report ID associated with the tile.",
						},
						"dataset_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dataset ID associated with the tile.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAppDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	appID := d.Get("app_id").(string)
	
	var dashboard *powerbiapi.AppDashboard
	var err error
	
	if dashboardID, ok := d.GetOk("id"); ok {
		// Get dashboard by ID
		dashboard, err = client.GetAppDashboard(appID, dashboardID.(string))
		if err != nil {
			return fmt.Errorf("failed to get app dashboard by ID %s: %w", dashboardID, err)
		}
	} else if dashboardName, ok := d.GetOk("display_name"); ok {
		// Get dashboard by name - need to list all dashboards and find by name
		dashboards, err := client.GetAppDashboards(appID)
		if err != nil {
			return fmt.Errorf("failed to list app dashboards: %w", err)
		}
		
		var foundDashboard *powerbiapi.AppDashboard
		for _, db := range dashboards.Value {
			if db.DisplayName == dashboardName.(string) {
				foundDashboard = &db
				break
			}
		}
		
		if foundDashboard == nil {
			return fmt.Errorf("dashboard with name '%s' not found in app %s", dashboardName, appID)
		}
		
		dashboard = foundDashboard
	} else {
		return fmt.Errorf("either 'id' or 'display_name' must be specified")
	}
	
	d.SetId(dashboard.ID)
	d.Set("id", dashboard.ID)
	d.Set("display_name", dashboard.DisplayName)
	d.Set("embed_url", dashboard.EmbedURL)
	d.Set("is_read_only", dashboard.IsReadOnly)
	d.Set("web_url", dashboard.WebURL)
	
	// Get tiles for this dashboard
	tiles, err := client.GetAppTiles(appID, dashboard.ID)
	if err != nil {
		return fmt.Errorf("failed to get tiles for app dashboard %s: %w", dashboard.ID, err)
	}
	
	tilesList := make([]interface{}, len(tiles.Value))
	for i, tile := range tiles.Value {
		tileMap := map[string]interface{}{
			"id":         tile.ID,
			"title":      tile.Title,
			"subtitle":   tile.SubTitle,
			"embed_url":  tile.EmbedURL,
			"embed_data": tile.EmbedData,
			"row_span":   tile.RowSpan,
			"col_span":   tile.ColSpan,
			"report_id":  tile.ReportID,
			"dataset_id": tile.DatasetID,
		}
		tilesList[i] = tileMap
	}
	d.Set("tiles", tilesList)
	
	return nil
}