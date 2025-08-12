package powerbi

import (
	"fmt"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceGatewayDatasourceUser represents a Power BI gateway datasource user
func ResourceGatewayDatasourceUser() *schema.Resource {
	return &schema.Resource{
		Create: createGatewayDatasourceUser,
		Read:   readGatewayDatasourceUser,
		Delete: deleteGatewayDatasourceUser,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid datasource user import id format, expected 'gateway_id/datasource_id/user_id'")
				}
				d.Set("gateway_id", parts[0])
				d.Set("datasource_id", parts[1])
				d.SetId(parts[2])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the gateway.",
			},
			"datasource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the datasource.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Email address of the user.",
				ExactlyOneOf: []string{"email_address", "identifier", "graph_id"},
			},
			"identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Identifier of the user.",
				ExactlyOneOf: []string{"email_address", "identifier", "graph_id"},
			},
			"graph_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Graph ID of the user.",
				ExactlyOneOf: []string{"email_address", "identifier", "graph_id"},
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Display name of the user.",
			},
			"principal_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "User",
				Description: "Type of principal (User, Group, or App).",
				ValidateFunc: validation.StringInSlice([]string{
					"User", "Group", "App",
				}, false),
			},
			"datasource_access_right": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Access right for the datasource.",
				ValidateFunc: validation.StringInSlice([]string{
					"Read", "ReadOverrideEffectiveIdentity",
				}, false),
			},
			// Computed fields
			"computed_display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed display name from the API response.",
			},
			"computed_email_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed email address from the API response.",
			},
			"computed_identifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed identifier from the API response.",
			},
			"computed_graph_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed graph ID from the API response.",
			},
		},
	}
}

func createGatewayDatasourceUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	datasourceID := d.Get("datasource_id").(string)
	
	request := powerbiapi.AddDatasourceUserRequest{
		DatasourceAccessRight: d.Get("datasource_access_right").(string),
		PrincipalType:         d.Get("principal_type").(string),
	}
	
	if v, ok := d.GetOk("email_address"); ok {
		request.EmailAddress = v.(string)
	}
	
	if v, ok := d.GetOk("identifier"); ok {
		request.Identifier = v.(string)
	}
	
	if v, ok := d.GetOk("graph_id"); ok {
		request.GraphID = v.(string)
	}
	
	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = v.(string)
	}
	
	err := client.AddDatasourceUser(gatewayID, datasourceID, request)
	if err != nil {
		return fmt.Errorf("failed to add datasource user: %w", err)
	}
	
	// Generate a unique ID for this user assignment
	userID := generateDatasourceUserID(request)
	d.SetId(userID)
	
	return readGatewayDatasourceUser(d, meta)
}

func readGatewayDatasourceUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	datasourceID := d.Get("datasource_id").(string)
	
	if gatewayID == "" || datasourceID == "" {
		return fmt.Errorf("gateway_id and datasource_id are required to read datasource user")
	}
	
	users, err := client.GetDatasourceUsers(gatewayID, datasourceID)
	if err != nil {
		return fmt.Errorf("failed to get datasource users: %w", err)
	}
	
	// Find the user based on the provided identifier
	var foundUser *powerbiapi.DatasourceUser
	for _, user := range users.Value {
		if matchesDatasourceUser(d, &user) {
			foundUser = &user
			break
		}
	}
	
	if foundUser == nil {
		// User not found, remove from state
		d.SetId("")
		return nil
	}
	
	// Set computed values
	d.Set("computed_display_name", foundUser.DisplayName)
	d.Set("computed_email_address", foundUser.EmailAddress)
	d.Set("computed_identifier", foundUser.Identifier)
	d.Set("computed_graph_id", foundUser.GraphID)
	d.Set("datasource_access_right", foundUser.DatasourceAccessRight)
	d.Set("principal_type", foundUser.PrincipalType)
	
	return nil
}

func deleteGatewayDatasourceUser(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	datasourceID := d.Get("datasource_id").(string)
	userID := d.Id()
	
	err := client.DeleteDatasourceUser(gatewayID, datasourceID, userID)
	if err != nil {
		// Ignore 404 errors - user already removed
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to delete datasource user: %w", err)
	}
	
	return nil
}

// generateDatasourceUserID creates a unique identifier for the user assignment
func generateDatasourceUserID(request powerbiapi.AddDatasourceUserRequest) string {
	if request.EmailAddress != "" {
		return fmt.Sprintf("email_%s", request.EmailAddress)
	}
	if request.Identifier != "" {
		return fmt.Sprintf("id_%s", request.Identifier)
	}
	if request.GraphID != "" {
		return fmt.Sprintf("graph_%s", request.GraphID)
	}
	if request.DisplayName != "" {
		return fmt.Sprintf("name_%s", request.DisplayName)
	}
	return "unknown"
}

// matchesDatasourceUser checks if a user from the API matches the resource definition
func matchesDatasourceUser(d *schema.ResourceData, user *powerbiapi.DatasourceUser) bool {
	if v, ok := d.GetOk("email_address"); ok {
		return user.EmailAddress == v.(string)
	}
	if v, ok := d.GetOk("identifier"); ok {
		return user.Identifier == v.(string)
	}
	if v, ok := d.GetOk("graph_id"); ok {
		return user.GraphID == v.(string)
	}
	if v, ok := d.GetOk("display_name"); ok {
		return user.DisplayName == v.(string)
	}
	return false
}