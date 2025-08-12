package powerbi

import (
	"fmt"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceDataflow represents a Power BI dataflow
func ResourceDataflow() *schema.Resource {
	return &schema.Resource{
		Create: createDataflow,
		Read:   readDataflow,
		Update: updateDataflow,
		Delete: deleteDataflow,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid dataflow import id format, expected 'workspace_id/dataflow_id'")
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
				Description: "ID of the workspace where the dataflow will be created.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "Name of the dataflow.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the dataflow.",
			},
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "JSON definition of the dataflow schema.",
			},
			"allow_native_queries": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to allow native queries in the dataflow.",
			},
			// Computed fields
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
		},
	}
}

func createDataflow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	
	request := powerbiapi.CreateDataflowRequest{
		Name: d.Get("name").(string),
	}
	
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	
	if v, ok := d.GetOk("definition"); ok {
		request.Definition = v.(string)
	}
	
	dataflow, err := client.CreateDataflow(workspaceID, request)
	if err != nil {
		return fmt.Errorf("failed to create dataflow: %w", err)
	}
	
	d.SetId(dataflow.ObjectID)
	
	return readDataflow(d, meta)
}

func readDataflow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Id()
	
	if workspaceID == "" {
		return fmt.Errorf("workspace_id is required to read dataflow")
	}
	
	dataflow, err := client.GetDataflow(workspaceID, dataflowID)
	if err != nil {
		// Check if dataflow was deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read dataflow: %w", err)
	}
	
	d.Set("name", dataflow.Name)
	d.Set("description", dataflow.Description)
	d.Set("model_url", dataflow.ModelURL)
	d.Set("configured_by", dataflow.ConfiguredBy)
	d.Set("modified_by", dataflow.ModifiedBy)
	
	if !dataflow.ModifiedDateTime.IsZero() {
		d.Set("modified_date_time", dataflow.ModifiedDateTime.Format("2006-01-02T15:04:05Z"))
	}
	
	return nil
}

func updateDataflow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Id()
	
	if d.HasChanges("name", "description", "allow_native_queries") {
		request := powerbiapi.UpdateDataflowRequest{}
		
		if d.HasChange("name") {
			request.Name = d.Get("name").(string)
		}
		
		if d.HasChange("description") {
			request.Description = d.Get("description").(string)
		}
		
		if d.HasChange("allow_native_queries") {
			request.AllowNativeQueries = d.Get("allow_native_queries").(bool)
		}
		
		err := client.UpdateDataflow(workspaceID, dataflowID, request)
		if err != nil {
			return fmt.Errorf("failed to update dataflow: %w", err)
		}
	}
	
	return readDataflow(d, meta)
}

func deleteDataflow(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Id()
	
	err := client.DeleteDataflow(workspaceID, dataflowID)
	if err != nil {
		// Ignore 404 errors - dataflow already deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to delete dataflow: %w", err)
	}
	
	return nil
}