package powerbi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceDataflowRefreshSchedule represents a Power BI dataflow refresh schedule
func ResourceDataflowRefreshSchedule() *schema.Resource {
	return &schema.Resource{
		Create: createDataflowRefreshSchedule,
		Read:   readDataflowRefreshSchedule,
		Update: updateDataflowRefreshSchedule,
		Delete: deleteDataflowRefreshSchedule,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid dataflow refresh schedule import id format, expected 'workspace_id/dataflow_id'")
				}
				d.Set("workspace_id", parts[0])
				d.Set("dataflow_id", parts[1])
				d.SetId(fmt.Sprintf("%s/%s", parts[0], parts[1]))
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the workspace containing the dataflow.",
			},
			"dataflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the dataflow.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the refresh schedule is enabled.",
			},
			"days": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Days of the week when the dataflow should be refreshed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
					}, false),
				},
			},
			"times": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Times of day when the dataflow should be refreshed (in HH:MM format).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`),
						"time must be in HH:MM format (24-hour)",
					),
				},
			},
			"local_time_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time zone ID for the refresh schedule.",
			},
			"notify_option": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NoNotification",
				Description: "Notification option for refresh failures.",
				ValidateFunc: validation.StringInSlice([]string{
					"NoNotification", "MailOnFailure", "MailOnCompletion",
				}, false),
			},
		},
	}
}

func createDataflowRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Get("dataflow_id").(string)
	
	schedule := buildDataflowRefreshSchedule(d)
	
	request := powerbiapi.UpdateDataflowRefreshScheduleRequest{
		Value: schedule,
	}
	
	err := client.UpdateDataflowRefreshSchedule(workspaceID, dataflowID, request)
	if err != nil {
		return fmt.Errorf("failed to create dataflow refresh schedule: %w", err)
	}
	
	d.SetId(fmt.Sprintf("%s/%s", workspaceID, dataflowID))
	
	return readDataflowRefreshSchedule(d, meta)
}

func readDataflowRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Get("dataflow_id").(string)
	
	if workspaceID == "" || dataflowID == "" {
		return fmt.Errorf("workspace_id and dataflow_id are required to read dataflow refresh schedule")
	}
	
	schedule, err := client.GetDataflowRefreshSchedule(workspaceID, dataflowID)
	if err != nil {
		// Check if schedule was deleted (or dataflow doesn't exist)
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read dataflow refresh schedule: %w", err)
	}
	
	d.Set("enabled", schedule.Enabled)
	d.Set("days", schedule.Days)
	d.Set("times", schedule.Times)
	d.Set("local_time_zone_id", schedule.LocalTimeZoneID)
	d.Set("notify_option", schedule.NotifyOption)
	
	return nil
}

func updateDataflowRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Get("dataflow_id").(string)
	
	schedule := buildDataflowRefreshSchedule(d)
	
	request := powerbiapi.UpdateDataflowRefreshScheduleRequest{
		Value: schedule,
	}
	
	err := client.UpdateDataflowRefreshSchedule(workspaceID, dataflowID, request)
	if err != nil {
		return fmt.Errorf("failed to update dataflow refresh schedule: %w", err)
	}
	
	return readDataflowRefreshSchedule(d, meta)
}

func deleteDataflowRefreshSchedule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	workspaceID := d.Get("workspace_id").(string)
	dataflowID := d.Get("dataflow_id").(string)
	
	// Disable the refresh schedule instead of deleting it
	schedule := powerbiapi.DataflowRefreshSchedule{
		Enabled: false,
	}
	
	request := powerbiapi.UpdateDataflowRefreshScheduleRequest{
		Value: schedule,
	}
	
	err := client.UpdateDataflowRefreshSchedule(workspaceID, dataflowID, request)
	if err != nil {
		// Ignore 404 errors - dataflow or schedule already deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to disable dataflow refresh schedule: %w", err)
	}
	
	return nil
}

func buildDataflowRefreshSchedule(d *schema.ResourceData) powerbiapi.DataflowRefreshSchedule {
	schedule := powerbiapi.DataflowRefreshSchedule{
		Enabled: d.Get("enabled").(bool),
	}
	
	if v, ok := d.GetOk("days"); ok {
		days := make([]string, 0)
		for _, day := range v.([]interface{}) {
			days = append(days, day.(string))
		}
		schedule.Days = days
	}
	
	if v, ok := d.GetOk("times"); ok {
		times := make([]string, 0)
		for _, time := range v.([]interface{}) {
			times = append(times, time.(string))
		}
		schedule.Times = times
	}
	
	if v, ok := d.GetOk("local_time_zone_id"); ok {
		schedule.LocalTimeZoneID = v.(string)
	}
	
	if v, ok := d.GetOk("notify_option"); ok {
		schedule.NotifyOption = v.(string)
	}
	
	return schedule
}