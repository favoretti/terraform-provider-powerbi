package powerbi

import (
	"fmt"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider represents the powerbi terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_TENANT_ID", ""),
				Description: "The Tenant ID for the tenant which contains the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_TENANT_ID` Environment Variable. Required unless using Managed Identity or Azure CLI authentication.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CLIENT_ID", ""),
				Description: "Also called Application ID. The Client ID for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_ID` Environment Variable. Required unless using Managed Identity or Azure CLI authentication.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CLIENT_SECRET", ""),
				Description: "Also called Application Secret. The Client Secret for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_SECRET` Environment Variable. Cannot be used with certificate_path or certificate_data.",
				ConflictsWith: []string{"certificate_path", "certificate_data"},
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_USERNAME", ""),
				Description: "The username for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_USERNAME` Environment Variable. Deprecated: Use Service Principal authentication instead.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_PASSWORD", ""),
				Description: "The password for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_PASSWORD` Environment Variable. Deprecated: Use Service Principal authentication instead.",
			},
			"certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CERTIFICATE_PATH", ""),
				Description: "The path to a PEM or PKCS#12 certificate file to use for Service Principal authentication. This can also be sourced from the `POWERBI_CERTIFICATE_PATH` Environment Variable. Cannot be used with client_secret.",
				ConflictsWith: []string{"certificate_data", "client_secret"},
			},
			"certificate_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CERTIFICATE_DATA", ""),
				Description: "Base64 encoded PEM certificate data to use for Service Principal authentication. This can also be sourced from the `POWERBI_CERTIFICATE_DATA` Environment Variable. Cannot be used with client_secret.",
				ConflictsWith: []string{"certificate_path", "client_secret"},
			},
			"certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CERTIFICATE_PASSWORD", ""),
				Description: "The password for the certificate file (if required). This can also be sourced from the `POWERBI_CERTIFICATE_PASSWORD` Environment Variable.",
			},
			"use_managed_identity": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_USE_MANAGED_IDENTITY", false),
				Description: "Use Managed Identity for authentication. This will automatically detect if running in Azure (App Service, Function, VM, etc.) and use the appropriate managed identity endpoint. This can also be sourced from the `POWERBI_USE_MANAGED_IDENTITY` Environment Variable.",
				ConflictsWith: []string{"use_azure_cli", "access_token"},
			},
			"managed_identity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_MANAGED_IDENTITY_ID", ""),
				Description: "The User Assigned Managed Identity ID to use for authentication. Leave empty to use System Assigned Managed Identity. This can also be sourced from the `POWERBI_MANAGED_IDENTITY_ID` Environment Variable.",
			},
			"use_azure_cli": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_USE_AZURE_CLI", false),
				Description: "Use Azure CLI for authentication. The Azure CLI must be installed and logged in (`az login`). This can also be sourced from the `POWERBI_USE_AZURE_CLI` Environment Variable.",
				ConflictsWith: []string{"use_managed_identity", "access_token"},
			},
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_ACCESS_TOKEN", ""),
				Description: "A pre-obtained access token to use for authentication. This can also be sourced from the `POWERBI_ACCESS_TOKEN` Environment Variable. Note: The token must have the appropriate Power BI scopes.",
				ConflictsWith: []string{"use_managed_identity", "use_azure_cli"},
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"powerbi_workspace":                ResourceWorkspace(),
			"powerbi_pbix":                     ResourcePBIX(),
			"powerbi_refresh_schedule":         ResourceRefreshSchedule(),
			"powerbi_workspace_access":         ResourceGroupUsers(),
			"powerbi_dataset":                  ResourceDataset(),
			"powerbi_dashboard":                ResourceDashboard(),
			"powerbi_dashboard_tile":           ResourceDashboardTile(),
			"powerbi_gateway_datasource":       ResourceGatewayDatasource(),
			"powerbi_gateway_datasource_user":  ResourceGatewayDatasourceUser(),
			"powerbi_dataflow":                 ResourceDataflow(),
			"powerbi_dataflow_refresh_schedule": ResourceDataflowRefreshSchedule(),
			"powerbi_deployment_pipeline":      ResourceDeploymentPipeline(),
			"powerbi_pipeline_stage":           ResourcePipelineStage(),
			"powerbi_pipeline_operation":       ResourcePipelineOperation(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"powerbi_workspace":       DataSourceWorkspace(),
			"powerbi_dashboard":       DataSourceDashboard(),
			"powerbi_dashboard_tiles": DataSourceDashboardTiles(),
			"powerbi_gateway":         DataSourceGateway(),
			"powerbi_dataflow":        DataSourceDataflow(),
			"powerbi_app":             DataSourceApp(),
			"powerbi_app_dashboard":   DataSourceAppDashboard(),
			"powerbi_app_report":      DataSourceAppReport(),
			"powerbi_embed_token":     DataSourceEmbedToken(),
			"powerbi_template_app":    DataSourceTemplateApp(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &powerbiapi.AuthConfig{
		TenantID:            d.Get("tenant_id").(string),
		ClientID:            d.Get("client_id").(string),
		ClientSecret:        d.Get("client_secret").(string),
		Username:            d.Get("username").(string),
		Password:            d.Get("password").(string),
		CertificatePath:     d.Get("certificate_path").(string),
		CertificateData:     d.Get("certificate_data").(string),
		CertificatePassword: d.Get("certificate_password").(string),
		UseManagedIdentity:  d.Get("use_managed_identity").(bool),
		ManagedIdentityID:   d.Get("managed_identity_id").(string),
		UseAzureCLI:         d.Get("use_azure_cli").(bool),
		AccessToken:         d.Get("access_token").(string),
	}

	// Validate authentication configuration
	if err := validateAuthenticationConfig(config); err != nil {
		return nil, err
	}

	return powerbiapi.NewClientWithAuthConfig(config)
}

func validateAuthenticationConfig(config *powerbiapi.AuthConfig) error {
	// Count the number of authentication methods configured
	// Check for high-priority authentication methods first
	authMethods := 0
	var activeMethod string

	// Priority 1: Direct access token
	if config.AccessToken != "" {
		authMethods++
		activeMethod = "access_token"
	}

	// Priority 2: Managed Identity
	if config.UseManagedIdentity {
		authMethods++
		activeMethod = "managed_identity"
	}

	// Priority 3: Azure CLI
	if config.UseAzureCLI {
		authMethods++
		activeMethod = "azure_cli"
	}

	// Priority 4: Username/password (includes client_secret)
	if config.Username != "" && config.Password != "" {
		authMethods++
		activeMethod = "username_password"
	} else {
		// Priority 5: Certificate authentication
		if config.CertificatePath != "" || config.CertificateData != "" {
			authMethods++
			activeMethod = "certificate"
		}
		
		// Priority 6: Client secret (only if not username/password)
		if config.ClientSecret != "" {
			authMethods++
			activeMethod = "client_secret"
		}
	}

	// Check for multiple authentication methods
	if authMethods > 1 {
		return fmt.Errorf("multiple authentication methods configured. Please use only one authentication method")
	}

	// Check for no authentication method
	if authMethods == 0 {
		return fmt.Errorf("no authentication method configured. Please configure one of: access_token, managed_identity, azure_cli, certificate, client_secret, or username/password")
	}

	// Validate specific authentication method requirements
	switch activeMethod {
	case "certificate":
		if config.TenantID == "" {
			return fmt.Errorf("tenant_id is required when using certificate authentication")
		}
		if config.ClientID == "" {
			return fmt.Errorf("client_id is required when using certificate authentication")
		}
		if config.CertificatePath != "" && config.CertificateData != "" {
			return fmt.Errorf("certificate_path and certificate_data cannot be used together")
		}

	case "client_secret":
		if config.TenantID == "" {
			return fmt.Errorf("tenant_id is required when using client_secret authentication")
		}
		if config.ClientID == "" {
			return fmt.Errorf("client_id is required when using client_secret authentication")
		}

	case "username_password":
		if config.TenantID == "" {
			return fmt.Errorf("tenant_id is required when using username/password authentication")
		}
		if config.ClientID == "" {
			return fmt.Errorf("client_id is required when using username/password authentication")
		}
		if config.ClientSecret == "" {
			return fmt.Errorf("client_secret is required when using username/password authentication")
		}
		if config.Username == "" {
			return fmt.Errorf("username is required when using username/password authentication")
		}
		if config.Password == "" {
			return fmt.Errorf("password is required when using username/password authentication")
		}

	case "managed_identity":
		// managed_identity_id is optional
		
	case "azure_cli":
		// No additional validation required
		
	case "access_token":
		// No additional validation required
	}

	return nil
}
