package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceAppReport returns reports from a Power BI app
func DataSourceAppReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppReportRead,

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
				Description: "ID of the report.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the report.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"web_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Web URL of the report.",
			},
			"embed_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Embed URL of the report.",
			},
			"dataset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dataset ID associated with the report.",
			},
		},
	}
}

func dataSourceAppReportRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	appID := d.Get("app_id").(string)
	
	var report *powerbiapi.AppReport
	var err error
	
	if reportID, ok := d.GetOk("id"); ok {
		// Get report by ID
		report, err = client.GetAppReport(appID, reportID.(string))
		if err != nil {
			return fmt.Errorf("failed to get app report by ID %s: %w", reportID, err)
		}
	} else if reportName, ok := d.GetOk("name"); ok {
		// Get report by name - need to list all reports and find by name
		reports, err := client.GetAppReports(appID)
		if err != nil {
			return fmt.Errorf("failed to list app reports: %w", err)
		}
		
		var foundReport *powerbiapi.AppReport
		for _, r := range reports.Value {
			if r.Name == reportName.(string) {
				foundReport = &r
				break
			}
		}
		
		if foundReport == nil {
			return fmt.Errorf("report with name '%s' not found in app %s", reportName, appID)
		}
		
		report = foundReport
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}
	
	d.SetId(report.ID)
	d.Set("id", report.ID)
	d.Set("name", report.Name)
	d.Set("web_url", report.WebURL)
	d.Set("embed_url", report.EmbedURL)
	d.Set("dataset_id", report.DatasetID)
	
	return nil
}