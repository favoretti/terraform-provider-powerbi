package powerbi

import (
	"fmt"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceDashboardTile represents a Power BI dashboard tile
func ResourceDashboardTile() *schema.Resource {
	return &schema.Resource{
		Create: createDashboardTile,
		Read:   readDashboardTile,
		Delete: deleteDashboardTile,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid tile import id format, expected 'workspace_id/dashboard_id/tile_id'")
				}
				d.Set("workspace_id", parts[0])
				d.Set("dashboard_id", parts[1])
				d.SetId(parts[2])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the workspace containing the dashboard.",
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the dashboard to add the tile to.",
			},
			"source_dashboard_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the source dashboard to clone tile from.",
			},
			"source_tile_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the source tile to clone.",
			},
			"target_workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the target workspace (if different from source).",
			},
			"target_report_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the target report (if rebinding to different report).",
			},
			"target_model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the target model (if rebinding to different model).",
			},
			"position_conflict_action": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "Tail",
				Description: "Action to take if tile position conflicts. Options: 'Tail' or 'Abort'.",
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
		},
	}
}

func createDashboardTile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	sourceDashboardID := d.Get("source_dashboard_id").(string)
	sourceTileID := d.Get("source_tile_id").(string)
	
	request := powerbiapi.CloneTileRequest{
		TargetDashboardID:      d.Get("dashboard_id").(string),
		PositionConflictAction: d.Get("position_conflict_action").(string),
	}
	
	if v, ok := d.GetOk("target_workspace_id"); ok {
		request.TargetWorkspaceID = v.(string)
	}
	
	if v, ok := d.GetOk("target_report_id"); ok {
		request.TargetReportID = v.(string)
	}
	
	if v, ok := d.GetOk("target_model_id"); ok {
		request.TargetModelID = v.(string)
	}
	
	tile, err := client.CloneTile(workspaceID, sourceDashboardID, sourceTileID, request)
	if err != nil {
		return fmt.Errorf("failed to clone tile: %w", err)
	}
	
	d.SetId(tile.ID)
	
	return readDashboardTile(d, meta)
}

func readDashboardTile(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dashboardID := d.Get("dashboard_id").(string)
	tileID := d.Id()
	
	if workspaceID == "" || dashboardID == "" {
		return fmt.Errorf("workspace_id and dashboard_id are required to read tile")
	}
	
	tile, err := client.GetTile(workspaceID, dashboardID, tileID)
	if err != nil {
		// Check if tile was deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read tile: %w", err)
	}
	
	d.Set("title", tile.Title)
	d.Set("subtitle", tile.SubTitle)
	d.Set("embed_url", tile.EmbedURL)
	d.Set("embed_data", tile.EmbedData)
	d.Set("report_id", tile.ReportID)
	d.Set("dataset_id", tile.DatasetID)
	d.Set("row_span", tile.RowSpan)
	d.Set("col_span", tile.ColSpan)
	
	return nil
}

func deleteDashboardTile(d *schema.ResourceData, meta interface{}) error {
	// Note: Power BI API doesn't provide a direct way to delete tiles
	// Tiles are typically removed by deleting the dashboard or removing the source report
	// We'll just remove from state
	d.SetId("")
	return nil
}