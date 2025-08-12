package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceDeploymentPipeline represents a Power BI deployment pipeline
func ResourceDeploymentPipeline() *schema.Resource {
	return &schema.Resource{
		Create: createDeploymentPipeline,
		Read:   readDeploymentPipeline,
		Update: updateDeploymentPipeline,
		Delete: deleteDeploymentPipeline,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
				Description:  "Display name of the deployment pipeline.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the deployment pipeline.",
			},
			// Computed fields
			"stages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of stages in the deployment pipeline.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Order of the stage in the pipeline.",
						},
						"stage_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the stage.",
						},
						"is_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the stage is public.",
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the workspace assigned to this stage.",
						},
						"workspace_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the workspace assigned to this stage.",
						},
						"artifacts_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of artifacts in this stage.",
						},
					},
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of users with access to the deployment pipeline.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifier of the user.",
						},
						"access_right": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access right of the user.",
						},
						"principal_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of principal (User, Group, or App).",
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
						"user_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of user.",
						},
					},
				},
			},
		},
	}
}

func createDeploymentPipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	request := powerbiapi.CreatePipelineRequest{
		DisplayName: d.Get("display_name").(string),
	}
	
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	
	pipeline, err := client.CreatePipeline(request)
	if err != nil {
		return fmt.Errorf("failed to create deployment pipeline: %w", err)
	}
	
	d.SetId(pipeline.ID)
	
	return readDeploymentPipeline(d, meta)
}

func readDeploymentPipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Id()
	
	pipeline, err := client.GetPipeline(pipelineID)
	if err != nil {
		// Check if pipeline was deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read deployment pipeline: %w", err)
	}
	
	d.Set("display_name", pipeline.DisplayName)
	d.Set("description", pipeline.Description)
	
	// Set stages list
	stagesList := make([]interface{}, len(pipeline.Stages))
	for i, stage := range pipeline.Stages {
		stageMap := map[string]interface{}{
			"order":           stage.Order,
			"stage_name":      stage.StageName,
			"is_public":       stage.IsPublic,
			"workspace_id":    stage.WorkspaceID,
			"workspace_name":  stage.WorkspaceName,
			"artifacts_count": stage.ArtifactsCount,
		}
		stagesList[i] = stageMap
	}
	d.Set("stages", stagesList)
	
	// Set users list
	usersList := make([]interface{}, len(pipeline.Users))
	for i, user := range pipeline.Users {
		userMap := map[string]interface{}{
			"identifier":     user.Identifier,
			"access_right":   user.AccessRight,
			"principal_type": user.PrincipalType,
			"display_name":   user.DisplayName,
			"email_address":  user.EmailAddress,
			"graph_id":       user.GraphID,
			"user_type":      user.UserType,
		}
		usersList[i] = userMap
	}
	d.Set("users", usersList)
	
	return nil
}

func updateDeploymentPipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Id()
	
	if d.HasChanges("display_name", "description") {
		request := powerbiapi.UpdatePipelineRequest{}
		
		if d.HasChange("display_name") {
			request.DisplayName = d.Get("display_name").(string)
		}
		
		if d.HasChange("description") {
			request.Description = d.Get("description").(string)
		}
		
		_, err := client.UpdatePipeline(pipelineID, request)
		if err != nil {
			return fmt.Errorf("failed to update deployment pipeline: %w", err)
		}
	}
	
	return readDeploymentPipeline(d, meta)
}

func deleteDeploymentPipeline(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Id()
	
	err := client.DeletePipeline(pipelineID)
	if err != nil {
		// Ignore 404 errors - pipeline already deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to delete deployment pipeline: %w", err)
	}
	
	return nil
}