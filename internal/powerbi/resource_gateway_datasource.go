package powerbi

import (
	"fmt"
	"strings"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ResourceGatewayDatasource represents a Power BI gateway datasource
func ResourceGatewayDatasource() *schema.Resource {
	return &schema.Resource{
		Create: createGatewayDatasource,
		Read:   readGatewayDatasource,
		Update: updateGatewayDatasource,
		Delete: deleteGatewayDatasource,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid datasource import id format, expected 'gateway_id/datasource_id'")
				}
				d.Set("gateway_id", parts[0])
				d.SetId(parts[1])
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
			"datasource_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the datasource.",
			},
			"datasource_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of the datasource.",
				ValidateFunc: validation.StringInSlice([]string{
					"Sql", "Oracle", "OleDb", "ODBC", "SharePointList", "Web",
					"OData", "File", "Folder", "SharePointDocumentLibrary",
					"Hdfs", "AzureTable", "Exchange", "ActiveDirectory",
					"MySql", "PostgreSql", "Sybase", "DB2", "Teradata",
					"SapHana", "SapBw", "AnalysisServices", "AzureBlob",
					"AzureSql", "AzureSqlDw", "Informix", "GoogleAnalytics",
					"AmazonRedshift", "Impala", "Spark", "Smartsheet",
				}, false),
			},
			"connection_details": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Connection details for the datasource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Server name or address.",
						},
						"database": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database name.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URL for web-based datasources.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "File path for file-based datasources.",
						},
						"kind": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kind of datasource.",
						},
						"auth_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication method.",
						},
						"account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account name.",
						},
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain name.",
						},
						"email_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Email address for authentication.",
						},
						"login_server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login server.",
						},
						"class": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Class of the datasource.",
						},
					},
				},
			},
			"credential_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Basic",
				Description: "Type of credentials used for authentication.",
				ValidateFunc: validation.StringInSlice([]string{
					"Basic", "Windows", "OAuth2", "Anonymous", "Key",
				}, false),
			},
			"credential_details": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Sensitive:   true,
				Description: "Credential details for datasource authentication.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credentials": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Encrypted credentials.",
						},
						"encrypted_connection": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to use encrypted connection.",
						},
						"encryption_algorithm": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Encryption algorithm used.",
						},
						"privacy_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "None",
							Description: "Privacy level for the datasource.",
							ValidateFunc: validation.StringInSlice([]string{
								"None", "Public", "Organizational", "Private",
							}, false),
						},
						"use_caller_aad_identity": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to use caller's AAD identity.",
						},
						"use_end_user_oauth2_credentials": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to use end user OAuth2 credentials.",
						},
					},
				},
			},
			"connection_string": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Connection string for the datasource.",
			},
		},
	}
}

func createGatewayDatasource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	
	connectionDetails := buildConnectionDetails(d.Get("connection_details").([]interface{}))
	
	request := powerbiapi.CreateDatasourceRequest{
		DatasourceName:    d.Get("datasource_name").(string),
		DatasourceType:    d.Get("datasource_type").(string),
		ConnectionDetails: connectionDetails,
	}
	
	if credDetails, ok := d.GetOk("credential_details"); ok {
		request.CredentialDetails = buildCredentialDetails(credDetails.([]interface{}), d.Get("credential_type").(string))
	}
	
	datasource, err := client.CreateDatasource(gatewayID, request)
	if err != nil {
		return fmt.Errorf("failed to create gateway datasource: %w", err)
	}
	
	d.SetId(datasource.ID)
	
	return readGatewayDatasource(d, meta)
}

func readGatewayDatasource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	datasourceID := d.Id()
	
	if gatewayID == "" {
		return fmt.Errorf("gateway_id is required to read datasource")
	}
	
	datasource, err := client.GetDatasource(gatewayID, datasourceID)
	if err != nil {
		// Check if datasource was deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("failed to read gateway datasource: %w", err)
	}
	
	d.Set("gateway_id", datasource.GatewayID)
	d.Set("datasource_name", datasource.DatasourceName)
	d.Set("datasource_type", datasource.DatasourceType)
	d.Set("connection_string", datasource.ConnectionString)
	d.Set("credential_type", datasource.CredentialType)
	
	// Set connection details
	connectionDetails := []interface{}{
		map[string]interface{}{
			"server":       datasource.ConnectionDetails.Server,
			"database":     datasource.ConnectionDetails.Database,
			"url":          datasource.ConnectionDetails.URL,
			"path":         datasource.ConnectionDetails.Path,
			"kind":         datasource.ConnectionDetails.Kind,
			"auth_method":  datasource.ConnectionDetails.AuthMethod,
			"account":      datasource.ConnectionDetails.Account,
			"domain":       datasource.ConnectionDetails.Domain,
			"email_address": datasource.ConnectionDetails.EmailAddress,
			"login_server": datasource.ConnectionDetails.LoginServer,
			"class":        datasource.ConnectionDetails.Class,
		},
	}
	d.Set("connection_details", connectionDetails)
	
	return nil
}

func updateGatewayDatasource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	datasourceID := d.Id()
	
	if d.HasChange("credential_details") || d.HasChange("credential_type") {
		request := powerbiapi.UpdateDatasourceRequest{}
		
		if credDetails, ok := d.GetOk("credential_details"); ok {
			request.CredentialDetails = buildCredentialDetails(credDetails.([]interface{}), d.Get("credential_type").(string))
		}
		
		err := client.UpdateDatasource(gatewayID, datasourceID, request)
		if err != nil {
			return fmt.Errorf("failed to update gateway datasource: %w", err)
		}
	}
	
	return readGatewayDatasource(d, meta)
}

func deleteGatewayDatasource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*powerbiapi.Client)
	
	gatewayID := d.Get("gateway_id").(string)
	datasourceID := d.Id()
	
	err := client.DeleteDatasource(gatewayID, datasourceID)
	if err != nil {
		// Ignore 404 errors - datasource already deleted
		if httpErr, ok := err.(powerbiapi.HTTPUnsuccessfulError); ok {
			if httpErr.Response != nil && httpErr.Response.StatusCode == 404 {
				return nil
			}
		}
		return fmt.Errorf("failed to delete gateway datasource: %w", err)
	}
	
	return nil
}

func buildConnectionDetails(details []interface{}) powerbiapi.DatasourceConnectionDetails {
	if len(details) == 0 {
		return powerbiapi.DatasourceConnectionDetails{}
	}
	
	detail := details[0].(map[string]interface{})
	
	return powerbiapi.DatasourceConnectionDetails{
		Server:       getStringValue(detail, "server"),
		Database:     getStringValue(detail, "database"),
		URL:          getStringValue(detail, "url"),
		Path:         getStringValue(detail, "path"),
		Kind:         getStringValue(detail, "kind"),
		AuthMethod:   getStringValue(detail, "auth_method"),
		Account:      getStringValue(detail, "account"),
		Domain:       getStringValue(detail, "domain"),
		EmailAddress: getStringValue(detail, "email_address"),
		LoginServer:  getStringValue(detail, "login_server"),
		Class:        getStringValue(detail, "class"),
	}
}

func buildCredentialDetails(details []interface{}, credentialType string) *powerbiapi.DatasourceCredentialDetails {
	if len(details) == 0 {
		return nil
	}
	
	detail := details[0].(map[string]interface{})
	
	return &powerbiapi.DatasourceCredentialDetails{
		CredentialType:              credentialType,
		Credentials:                 getStringValue(detail, "credentials"),
		EncryptedConnection:         getStringValue(detail, "encrypted_connection"),
		EncryptionAlgorithm:         getStringValue(detail, "encryption_algorithm"),
		PrivacyLevel:                getStringValue(detail, "privacy_level"),
		UseCallerAADIdentity:        getBoolValue(detail, "use_caller_aad_identity"),
		UseEndUserOAuth2Credentials: getBoolValue(detail, "use_end_user_oauth2_credentials"),
	}
}

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return v.(string)
	}
	return ""
}

func getBoolValue(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok && v != nil {
		return v.(bool)
	}
	return false
}