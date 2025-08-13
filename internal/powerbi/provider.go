package powerbi

import (
	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider represents the powerbi terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_TENANT_ID", ""),
				Description: "The Tenant ID for the tenant which contains the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_TENANT_ID` Environment Variable",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CLIENT_ID", ""),
				Description: "Also called Application ID. The Client ID for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_ID` Environment Variable",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_CLIENT_SECRET", ""),
				Description: "Also called Application Secret. The Client Secret for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_SECRET` Environment Variable",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_USERNAME", ""),
				Description: "The username for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_USERNAME` Environment Variable",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POWERBI_PASSWORD", ""),
				Description: "The password for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_PASSWORD` Environment Variable",
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

	username, usernameOk := d.GetOk("username")
	password, passwordOk := d.GetOk("password")

	if usernameOk && passwordOk {
		return powerbiapi.NewClientWithPasswordAuth(
			d.Get("tenant_id").(string),
			d.Get("client_id").(string),
			d.Get("client_secret").(string),
			username.(string),
			password.(string),
		)
	}
	return powerbiapi.NewClientWithClientCredentialAuth(
		d.Get("tenant_id").(string),
		d.Get("client_id").(string),
		d.Get("client_secret").(string),
	)

}
