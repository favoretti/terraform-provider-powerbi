package powerbi

import (
	"fmt"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourcePipelineOperation represents a Power BI deployment pipeline operation
func ResourcePipelineOperation() *schema.Resource {
	return &schema.Resource{
		Create: createPipelineOperation,
		Read:   readPipelineOperation,
		Delete: deletePipelineOperation,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid pipeline operation import id format, expected 'pipeline_id/operation_id'")
				}
				d.Set("pipeline_id", parts[0])
				d.SetId(parts[1])
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
			"source_stage_order": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Order of the source stage to deploy from.",
			},
			"artifacts_to_deploy": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "List of artifacts to deploy. If not specified, all artifacts will be deployed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"artifact_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the artifact to deploy.",
						},
						"artifact_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the artifact to deploy.",
							ValidateFunc: validation.StringInSlice([]string{
								"Report", "Dashboard", "Dataset", "Dataflow",
							}, false),
						},
					},
				},
			},
			"options": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Deployment options.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_create_artifact": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow creating new artifacts during deployment.",
						},
						"allow_overwrite_artifact": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow overwriting existing artifacts during deployment.",
						},
						"allow_overwrite_target_schema": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow overwriting target schema during deployment.",
						},
						"allow_purge_data": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow purging data during deployment.",
						},
						"allow_skip_tiles_with_missing_prerequisites": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow skipping tiles with missing prerequisites.",
						},
						"allow_take_over": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Allow taking over artifacts during deployment.",
						},
					},
				},
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Note for the deployment operation.",
			},
			// Computed fields
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the operation.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the operation.",
			},
			"last_updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated time of the operation.",
			},
			"execution_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Execution start time of the operation.",
			},
			"execution_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Execution end time of the operation.",
			},
			"target_stage_order": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Order of the target stage.",
			},
			"error": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Error information if the operation failed.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error code.",
						},
						"error_details": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error details.",
						},
					},
				},
			},
		},
	}
}

func createPipelineOperation(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Get("pipeline_id").(string)
	
	request := powerbiapi.DeployRequest{
		SourceStageOrder: d.Get("source_stage_order").(int),
	}
	
	if v, ok := d.GetOk("artifacts_to_deploy"); ok {
		artifacts := make([]powerbiapi.DeployArtifact, 0)
		for _, artifact := range v.([]interface{}) {
			artifactMap := artifact.(map[string]interface{})
			artifacts = append(artifacts, powerbiapi.DeployArtifact{
				ArtifactID:   artifactMap["artifact_id"].(string),
				ArtifactType: artifactMap["artifact_type"].(string),
			})
		}
		request.ArtifactsToDeploy = artifacts
	}
	
	if v, ok := d.GetOk("options"); ok && len(v.([]interface{})) > 0 {
		optionsMap := v.([]interface{})[0].(map[string]interface{})
		request.Options = &powerbiapi.DeployOptions{
			AllowCreateArtifact:                      optionsMap["allow_create_artifact"].(bool),
			AllowOverwriteArtifact:                   optionsMap["allow_overwrite_artifact"].(bool),
			AllowOverwriteTargetSchema:               optionsMap["allow_overwrite_target_schema"].(bool),
			AllowPurgeData:                          optionsMap["allow_purge_data"].(bool),
			AllowSkipTilesWithMissingPrerequisites:  optionsMap["allow_skip_tiles_with_missing_prerequisites"].(bool),
			AllowTakeOver:                           optionsMap["allow_take_over"].(bool),
		}
	}
	
	if v, ok := d.GetOk("note"); ok {
		request.Note = v.(string)
	}
	
	response, err := client.DeployAll(pipelineID, request)
	if err != nil {
		return fmt.Errorf("failed to create pipeline operation: %w", err)
	}
	
	d.SetId(response.ID)
	
	return readPipelineOperation(d, meta)
}

func readPipelineOperation(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	pipelineID := d.Get("pipeline_id").(string)
	operationID := d.Id()
	
	if pipelineID == "" {
		return fmt.Errorf("pipeline_id is required to read pipeline operation")
	}
	
	operation, err := client.GetPipelineOperation(pipelineID, operationID)
	if err != nil {
		// Check if operation was deleted or doesn't exist
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read pipeline operation: %w", err)
	}
	
	d.Set("type", operation.Type)
	d.Set("status", operation.Status)
	d.Set("target_stage_order", operation.TargetStageOrder)
	d.Set("note", operation.Note)
	
	if !operation.LastUpdatedTime.IsZero() {
		d.Set("last_updated_time", operation.LastUpdatedTime.Format("2006-01-02T15:04:05Z"))
	}
	
	if !operation.ExecutionStartTime.IsZero() {
		d.Set("execution_start_time", operation.ExecutionStartTime.Format("2006-01-02T15:04:05Z"))
	}
	
	if !operation.ExecutionEndTime.IsZero() {
		d.Set("execution_end_time", operation.ExecutionEndTime.Format("2006-01-02T15:04:05Z"))
	}
	
	// Set error information if available
	if operation.Error != nil {
		errorInfo := []interface{}{
			map[string]interface{}{
				"error_code":    operation.Error.ErrorCode,
				"error_details": operation.Error.ErrorDetails,
			},
		}
		d.Set("error", errorInfo)
	}
	
	return nil
}

func deletePipelineOperation(d *schema.ResourceData, meta interface{}) error {
	// Pipeline operations cannot be deleted, they are historical records
	// Just remove from state
	d.SetId("")
	return nil
}