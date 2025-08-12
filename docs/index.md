# Power BI Provider
The Power BI provider can be used to configure Power BI resources using the [Power BI REST API](https://docs.microsoft.com/en-us/rest/api/power-bi/)

-> See the [authentication guide](guides/authentication.md) for details on how to generate authentication details

## Resources
The Power BI provider supports the following resources:

### Workspace Management
- [powerbi_workspace](resources/workspace.md) - Manage Power BI workspaces
- [powerbi_workspace_access](resources/workspace_access.md) - Manage workspace user access

### Dashboard Management
- [powerbi_dashboard](resources/dashboard.md) - Manage Power BI dashboards
- [powerbi_dashboard_tile](resources/dashboard_tile.md) - Clone tiles between dashboards

### Data Management
- [powerbi_dataset](resources/dataset.md) - Manage push datasets
- [powerbi_pbix](resources/pbix.md) - Deploy PBIX files to workspaces
- [powerbi_refresh_schedule](resources/refresh_schedule.md) - Configure dataset refresh schedules

### Gateway Management
- [powerbi_gateway_datasource](resources/gateway_datasource.md) - Manage gateway data sources
- [powerbi_gateway_datasource_user](resources/gateway_datasource_user.md) - Manage datasource user access

## Data Sources
The Power BI provider supports the following data sources:

### Workspace & Content Discovery
- [powerbi_workspace](data-sources/workspace.md) - Retrieve workspace information
- [powerbi_dashboard](data-sources/dashboard.md) - Retrieve dashboard information
- [powerbi_dashboard_tiles](data-sources/dashboard_tiles.md) - List dashboard tiles

### Gateway Discovery
- [powerbi_gateway](data-sources/gateway.md) - Retrieve gateway information

## Example Usage

### Basic Configuration
```hcl
provider "powerbi" {
  tenant_id     = "1c4cc30c-271e-47f2-891e-fef13f035bc7"
  client_id     = "f9ad3042-a969-4a31-826e-856d238df3b1"
  client_secret = "u94lE93qfJSJRTEGs@Pgs]]RZzM]V?bE"
  username      = "powerbiapp@mycompany.com"
  password      = "pass@word1!"
}
```

### Complete Workspace Setup with Dashboard and Gateway
```hcl
# Create workspace
resource "powerbi_workspace" "analytics" {
  name = "Analytics Workspace"
}

# Create dashboard
resource "powerbi_dashboard" "executive" {
  name         = "Executive Dashboard"
  workspace_id = powerbi_workspace.analytics.id
}

# Find existing gateway
data "powerbi_gateway" "enterprise" {
  name = "Enterprise Gateway"
}

# Create gateway datasource
resource "powerbi_gateway_datasource" "warehouse" {
  gateway_id      = data.powerbi_gateway.enterprise.id
  datasource_name = "Data Warehouse"
  datasource_type = "Sql"

  connection_details {
    server   = "warehouse.corp.com"
    database = "EDW"
  }

  credential_type = "Windows"
  
  credential_details {
    privacy_level           = "Organizational"
    use_caller_aad_identity = true
  }
}

# Grant access to analysts
resource "powerbi_gateway_datasource_user" "analysts" {
  gateway_id              = data.powerbi_gateway.enterprise.id
  datasource_id          = powerbi_gateway_datasource.warehouse.id
  email_address          = "analysts@company.com"
  datasource_access_right = "Read"
  principal_type         = "Group"
}
```

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `client_id` - (Required) Also called Application ID. The Client ID for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_ID` Environment Variable.
* `client_secret` - (Required) Also called Application Secret. The Client Secret for the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_CLIENT_SECRET` Environment Variable.
* `tenant_id` - (Required) The Tenant ID for the tenant which contains the Azure Active Directory App Registration to use for performing Power BI REST API operations. This can also be sourced from the `POWERBI_TENANT_ID` Environment Variable.
* `password` - (Optional) The password for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_PASSWORD` Environment Variable.
* `username` - (Optional) The username for the a Power BI user to use for performing Power BI REST API operations. If provided will use resource owner password credentials flow with delegate permissions. This can also be sourced from the `POWERBI_USERNAME` Environment Variable.
<!-- /docgen -->