package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceDashboardTiles returns a list of tiles from a dashboard
func DataSourceDashboardTiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDashboardTilesRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the workspace containing the dashboard.",
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the dashboard.",
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
						"configuration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration of the tile (if available).",
						},
					},
				},
			},
		},
	}
}

func dataSourceDashboardTilesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dashboardID := d.Get("dashboard_id").(string)
	
	tiles, err := client.GetTiles(workspaceID, dashboardID)
	if err != nil {
		return fmt.Errorf("failed to get tiles for dashboard %s: %w", dashboardID, err)
	}
	
	tilesList := make([]interface{}, len(tiles.Value))
	for i, tile := range tiles.Value {
		tileMap := map[string]interface{}{
			"id":         tile.ID,
			"title":      tile.Title,
			"subtitle":   tile.SubTitle,
			"embed_url":  tile.EmbedURL,
			"embed_data": tile.EmbedData,
			"report_id":  tile.ReportID,
			"dataset_id": tile.DatasetID,
			"row_span":   tile.RowSpan,
			"col_span":   tile.ColSpan,
		}
		
		if tile.Configuration != nil {
			tileMap["configuration"] = *tile.Configuration
		}
		
		tilesList[i] = tileMap
	}
	
	d.Set("tiles", tilesList)
	d.SetId(fmt.Sprintf("%s/%s", workspaceID, dashboardID))
	
	return nil
}