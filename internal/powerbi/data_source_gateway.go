package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// DataSourceGateway returns gateway information
func DataSourceGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGatewayRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the gateway.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the gateway.",
				ExactlyOneOf: []string{"id", "name"},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the gateway.",
			},
			"gateway_annotation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Annotation of the gateway.",
			},
			"gateway_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the gateway.",
			},
			"gateway_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the gateway.",
			},
			"gateway_machine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Machine where the gateway is installed.",
			},
			"gateway_contact_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Contact information for the gateway.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"gateway_cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the gateway cluster.",
			},
			"gateway_cluster_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the gateway cluster.",
			},
			"public_key": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Public key information for the gateway.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exponent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exponent of the public key.",
						},
						"modulus": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modulus of the public key.",
						},
					},
				},
			},
		},
	}
}

func dataSourceGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	var gateway *powerbiapi.Gateway
	var err error
	
	if gatewayID, ok := d.GetOk("id"); ok {
		// Get gateway by ID
		gateway, err = client.GetGateway(gatewayID.(string))
		if err != nil {
			return fmt.Errorf("failed to get gateway by ID %s: %w", gatewayID, err)
		}
	} else if gatewayName, ok := d.GetOk("name"); ok {
		// Get gateway by name - need to list all gateways and find by name
		gateways, err := client.GetGateways()
		if err != nil {
			return fmt.Errorf("failed to list gateways: %w", err)
		}
		
		var foundGateway *powerbiapi.Gateway
		for _, g := range gateways.Value {
			if g.Name == gatewayName.(string) {
				foundGateway = &g
				break
			}
		}
		
		if foundGateway == nil {
			return fmt.Errorf("gateway with name '%s' not found", gatewayName)
		}
		
		gateway = foundGateway
	} else {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}
	
	d.SetId(gateway.ID)
	d.Set("id", gateway.ID)
	d.Set("name", gateway.Name)
	d.Set("type", gateway.Type)
	d.Set("gateway_annotation", gateway.GatewayAnnotation)
	d.Set("gateway_status", gateway.GatewayStatus)
	d.Set("gateway_version", gateway.GatewayVersion)
	d.Set("gateway_machine", gateway.GatewayMachine)
	d.Set("gateway_contact_info", gateway.GatewayContactInfo)
	d.Set("gateway_cluster_id", gateway.GatewayClusterId)
	d.Set("gateway_cluster_status", gateway.GatewayClusterStatus)
	
	// Set public key information
	publicKey := []interface{}{
		map[string]interface{}{
			"exponent": gateway.PublicKey.Exponent,
			"modulus":  gateway.PublicKey.Modulus,
		},
	}
	d.Set("public_key", publicKey)
	
	return nil
}