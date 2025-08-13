package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceTemplateApp returns information about Power BI template apps
func DataSourceTemplateApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateAppRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the template app.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the template app.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the template app.",
			},
			"publisher_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the publisher.",
			},
			"publisher_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email of the publisher.",
			},
			"support_contact": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Support contact for the template app.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the template app.",
			},
			"logo_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the template app logo.",
			},
			"package_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the template app package.",
			},
		},
	}
}

func dataSourceTemplateAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)

	var templateApp *powerbiapi.TemplateApp
	var err error

	if templateAppID, ok := d.GetOk("id"); ok {
		// Get template app by ID
		templateApp, err = client.GetTemplateApp(templateAppID.(string))
		if err != nil {
			return fmt.Errorf("failed to get template app by ID %s: %w", templateAppID, err)
		}
	} else if templateAppName, ok := d.GetOk("name"); ok {
		// Get template app by name - need to list all template apps and find by name
		templateApps, err := client.GetTemplateApps()
		if err != nil {
			return fmt.Errorf("failed to list template apps: %w", err)
		}

		var foundTemplateApp *powerbiapi.TemplateApp
		for _, ta := range templateApps.Value {
			if ta.Name == templateAppName.(string) {
				foundTemplateApp = &ta
				break
			}
		}

		if foundTemplateApp == nil {
			return fmt.Errorf("template app with name '%s' not found", templateAppName)
		}

		templateApp = foundTemplateApp
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}

	d.SetId(templateApp.ID)
	d.Set("id", templateApp.ID)
	d.Set("name", templateApp.Name)
	d.Set("description", templateApp.Description)
	d.Set("publisher_name", templateApp.PublisherName)
	d.Set("publisher_email", templateApp.PublisherEmail)
	d.Set("support_contact", templateApp.SupportContact)
	d.Set("version", templateApp.Version)
	d.Set("logo_url", templateApp.LogoURL)
	d.Set("package_url", templateApp.PackageURL)

	return nil
}