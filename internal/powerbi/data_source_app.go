package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceApp returns a specific Power BI app
func DataSourceApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the app.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the app.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the app.",
			},
			"published_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User who published the app.",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the app was last updated.",
			},
		},
	}
}

func dataSourceAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	var app *powerbiapi.App
	var err error
	
	if appID, ok := d.GetOk("id"); ok {
		// Get app by ID
		app, err = client.GetApp(appID.(string))
		if err != nil {
			return fmt.Errorf("failed to get app by ID %s: %w", appID, err)
		}
	} else if appName, ok := d.GetOk("name"); ok {
		// Get app by name - need to list all apps and find by name
		apps, err := client.GetApps()
		if err != nil {
			return fmt.Errorf("failed to list apps: %w", err)
		}
		
		var foundApp *powerbiapi.App
		for _, a := range apps.Value {
			if a.Name == appName.(string) {
				foundApp = &a
				break
			}
		}
		
		if foundApp == nil {
			return fmt.Errorf("app with name '%s' not found", appName)
		}
		
		app = foundApp
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}
	
	d.SetId(app.ID)
	d.Set("id", app.ID)
	d.Set("name", app.Name)
	d.Set("description", app.Description)
	d.Set("published_by", app.PublishedBy)
	d.Set("last_update", app.LastUpdate)
	
	return nil
}