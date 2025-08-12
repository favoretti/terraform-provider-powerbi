package powerbi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourcePipelineStage represents a Power BI deployment pipeline stage
func ResourcePipelineStage() *schema.Resource {
	return &schema.Resource{
		Create: createPipelineStage,
		Read:   readPipelineStage,
		Delete: deletePipelineStage,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid pipeline stage import id format, expected 'pipeline_id/stage_order'")
				}
				d.Set("pipeline_id", parts[0])
				stageOrder, err := strconv.Atoi(parts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid stage order: %w", err)
				}
				d.Set("stage_order", stageOrder)
				d.SetId(fmt.Sprintf("%s/%d", parts[0], stageOrder))
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the deployment pipeline.",
			},
			"stage_order": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Order of the stage in the pipeline (0-based).",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the workspace to assign to this stage.",
			},
			// Computed fields
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
	}
}

func createPipelineStage(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Get("pipeline_id").(string)
	stageOrder := d.Get("stage_order").(int)
	workspaceID := d.Get("workspace_id").(string)
	
	request := powerbiapi.AssignWorkspaceRequest{
		WorkspaceID: workspaceID,
	}
	
	err := client.AssignWorkspace(pipelineID, stageOrder, request)
	if err != nil {
		return fmt.Errorf("failed to assign workspace to pipeline stage: %w", err)
	}
	
	d.SetId(fmt.Sprintf("%s/%d", pipelineID, stageOrder))
	
	return readPipelineStage(d, meta)
}

func readPipelineStage(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Get("pipeline_id").(string)
	stageOrder := d.Get("stage_order").(int)
	
	if pipelineID == "" {
		return fmt.Errorf("pipeline_id is required to read pipeline stage")
	}
	
	stages, err := client.GetPipelineStages(pipelineID)
	if err != nil {
		return fmt.Errorf("failed to get pipeline stages: %w", err)
	}
	
	// Find the stage with the matching order
	var foundStage *powerbiapi.PipelineStage
	for _, stage := range stages {
		if stage.Order == stageOrder {
			foundStage = &stage
			break
		}
	}
	
	if foundStage == nil {
		// Stage not found or workspace unassigned
		d.SetId("")
		return nil
	}
	
	// Check if workspace is assigned
	if foundStage.WorkspaceID == "" {
		d.SetId("")
		return nil
	}
	
	d.Set("stage_name", foundStage.StageName)
	d.Set("is_public", foundStage.IsPublic)
	d.Set("workspace_id", foundStage.WorkspaceID)
	d.Set("workspace_name", foundStage.WorkspaceName)
	d.Set("artifacts_count", foundStage.ArtifactsCount)
	
	return nil
}

func deletePipelineStage(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Get("pipeline_id").(string)
	stageOrder := d.Get("stage_order").(int)
	workspaceID := d.Get("workspace_id").(string)
	
	request := powerbiapi.UnassignWorkspaceRequest{
		WorkspaceID: workspaceID,
	}
	
	err := client.UnassignWorkspace(pipelineID, stageOrder, request)
	if err != nil {
		// Ignore 404 errors - pipeline or stage already deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to unassign workspace from pipeline stage: %w", err)
	}
	
	return nil
}