package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceDataflow returns a specific dataflow from a workspace
func DataSourceDataflow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataflowRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the workspace containing the dataflow.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the dataflow.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the dataflow.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the dataflow.",
			},
			"model_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the dataflow model.",
			},
			"configured_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User who configured the dataflow.",
			},
			"modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User who last modified the dataflow.",
			},
			"modified_date_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the dataflow was last modified.",
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of users with access to the dataflow.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the user.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the user.",
						},
						"email_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address of the user.",
						},
						"graph_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Graph ID of the user.",
						},
						"principal_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of principal (User, Group, or App).",
						},
						"user_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of user.",
						},
					},
				},
			},
			"refresh_schedule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Refresh schedule configuration.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the refresh schedule is enabled.",
						},
						"days": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Days of the week when the dataflow should be refreshed.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"times": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Times of day when the dataflow should be refreshed.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"local_time_zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time zone ID for the refresh schedule.",
						},
						"notify_option": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notification option for refresh failures.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDataflowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	
	var dataflow *powerbiapi.Dataflow
	var err error
	
	if dataflowID, ok := d.GetOk("id"); ok {
		// Get dataflow by ID
		dataflow, err = client.GetDataflow(workspaceID, dataflowID.(string))
		if err != nil {
			return fmt.Errorf("failed to get dataflow by ID %s: %w", dataflowID, err)
		}
	} else if dataflowName, ok := d.GetOk("name"); ok {
		// Get dataflow by name - need to list all dataflows and find by name
		dataflows, err := client.GetDataflows(workspaceID)
		if err != nil {
			return fmt.Errorf("failed to list dataflows: %w", err)
		}
		
		var foundDataflow *powerbiapi.Dataflow
		for _, df := range dataflows.Value {
			if df.Name == dataflowName.(string) {
				foundDataflow = &df
				break
			}
		}
		
		if foundDataflow == nil {
			return fmt.Errorf("dataflow with name '%s' not found in workspace %s", dataflowName, workspaceID)
		}
		
		dataflow = foundDataflow
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}
	
	d.SetId(dataflow.ObjectID)
	d.Set("id", dataflow.ObjectID)
	d.Set("name", dataflow.Name)
	d.Set("description", dataflow.Description)
	d.Set("model_url", dataflow.ModelURL)
	d.Set("configured_by", dataflow.ConfiguredBy)
	d.Set("modified_by", dataflow.ModifiedBy)
	
	if !dataflow.ModifiedDateTime.IsZero() {
		d.Set("modified_date_time", dataflow.ModifiedDateTime.Format("2006-01-02T15:04:05Z"))
	}
	
	// Set users list
	usersList := make([]interface{}, len(dataflow.Users))
	for i, user := range dataflow.Users {
		userMap := map[string]interface{}{
			"identifier":     user.Identifier,
			"display_name":   user.DisplayName,
			"email_address":  user.EmailAddress,
			"graph_id":       user.GraphID,
			"principal_type": user.PrincipalType,
			"user_type":      user.UserType,
		}
		usersList[i] = userMap
	}
	d.Set("users", usersList)
	
	// Set refresh schedule if available
	if dataflow.RefreshSchedule != nil {
		refreshSchedule := []interface{}{
			map[string]interface{}{
				"enabled":              dataflow.RefreshSchedule.Enabled,
				"days":                 dataflow.RefreshSchedule.Days,
				"times":                dataflow.RefreshSchedule.Times,
				"local_time_zone_id":   dataflow.RefreshSchedule.LocalTimeZoneID,
				"notify_option":        dataflow.RefreshSchedule.NotifyOption,
			},
		}
		d.Set("refresh_schedule", refreshSchedule)
	}
	
	return nil
}