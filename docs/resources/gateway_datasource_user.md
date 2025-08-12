# Gateway Datasource User Resource
`powerbi_gateway_datasource_user` represents user access permissions for a gateway datasource within Power BI.

## Example Usage

### Grant user access to datasource
```hcl
data "powerbi_gateway" "enterprise" {
  name = "Enterprise Gateway"
}

resource "powerbi_gateway_datasource" "sql_server" {
  gateway_id      = data.powerbi_gateway.enterprise.id
  datasource_name = "Production SQL Server"
  datasource_type = "Sql"
  
  connection_details {
    server   = "sql.company.com"
    database = "ProductionDB"
  }
}

resource "powerbi_gateway_datasource_user" "analyst" {
  gateway_id              = data.powerbi_gateway.enterprise.id
  datasource_id          = powerbi_gateway_datasource.sql_server.id
  email_address          = "data.analyst@company.com"
  datasource_access_right = "Read"
  principal_type         = "User"
}
```

### Grant group access to datasource
```hcl
resource "powerbi_gateway_datasource_user" "analysts_group" {
  gateway_id              = data.powerbi_gateway.enterprise.id
  datasource_id          = powerbi_gateway_datasource.sql_server.id
  identifier             = "analysts@company.com"
  datasource_access_right = "ReadOverrideEffectiveIdentity"
  principal_type         = "Group"
  display_name           = "Data Analysts Group"
}
```

### Grant service principal access
```hcl
resource "powerbi_gateway_datasource_user" "service_app" {
  gateway_id              = data.powerbi_gateway.enterprise.id
  datasource_id          = powerbi_gateway_datasource.sql_server.id
  graph_id               = "12345678-1234-1234-1234-123456789abc"
  datasource_access_right = "Read"
  principal_type         = "App"
}
```

~> **Note:** Gateway datasource users cannot be updated after creation. Any changes will force a new resource to be created.

~> **Important:** You must specify exactly one of `email_address`, `identifier`, or `graph_id` to identify the user or principal.

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `gateway_id` - (Required, Forces new resource) ID of the gateway.
* `datasource_id` - (Required, Forces new resource) ID of the datasource.
* `email_address` - (Optional, Forces new resource, Conflicts with: `identifier`, `graph_id`) Email address of the user.
* `identifier` - (Optional, Forces new resource, Conflicts with: `email_address`, `graph_id`) Identifier of the user.
* `graph_id` - (Optional, Forces new resource, Conflicts with: `email_address`, `identifier`) Graph ID of the user.
* `display_name` - (Optional, Forces new resource) Display name of the user.
* `principal_type` - (Optional, Default: `User`, Forces new resource) Type of principal. Options: `User`, `Group`, `App`.
* `datasource_access_right` - (Required, Forces new resource) Access right for the datasource. Options: `Read`, `ReadOverrideEffectiveIdentity`.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the user assignment.
<!-- docgen:ComputedParameters -->
* `computed_display_name` - Computed display name from the API response.
* `computed_email_address` - Computed email address from the API response.
* `computed_identifier` - Computed identifier from the API response.
* `computed_graph_id` - Computed graph ID from the API response.
<!-- /docgen -->

## Import
Gateway datasource users can be imported using the gateway ID, datasource ID, and user ID separated by forward slashes:

```shell
terraform import powerbi_gateway_datasource_user.example gateway_id/datasource_id/user_id
```

Note: The user ID for import purposes is typically the email address, identifier, or graph ID used to create the assignment.