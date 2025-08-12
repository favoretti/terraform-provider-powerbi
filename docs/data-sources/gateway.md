# Gateway Data Source
`powerbi_gateway` represents a Power BI gateway used for connecting to on-premises data sources.

## Example Usage

### Find gateway by name
```hcl
data "powerbi_gateway" "enterprise" {
  name = "Enterprise Gateway"
}

output "gateway_status" {
  value = data.powerbi_gateway.enterprise.gateway_status
}

output "gateway_version" {
  value = data.powerbi_gateway.enterprise.gateway_version
}
```

### Find gateway by ID
```hcl
data "powerbi_gateway" "existing" {
  id = "gateway-12345-abcde-67890-fghij"
}

output "gateway_machine" {
  value = data.powerbi_gateway.existing.gateway_machine
}
```

### Use gateway data for datasource creation
```hcl
data "powerbi_gateway" "corp_gateway" {
  name = "Corporate Gateway"
}

resource "powerbi_gateway_datasource" "sql_server" {
  gateway_id      = data.powerbi_gateway.corp_gateway.id
  datasource_name = "Production SQL Server"
  datasource_type = "Sql"
  
  connection_details {
    server   = "sql.company.com"
    database = "ProductionDB"
  }
}
```

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `id` - (Optional, Conflicts with: `name`) ID of the gateway.
* `name` - (Optional, Conflicts with: `id`) Name of the gateway.
<!-- /docgen -->

~> **Note:** You must specify exactly one of `id` or `name` to identify the gateway.

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
<!-- docgen:ComputedParameters -->
* `type` - Type of the gateway.
* `gateway_annotation` - Annotation of the gateway.
* `gateway_status` - Status of the gateway.
* `gateway_version` - Version of the gateway.
* `gateway_machine` - Machine where the gateway is installed.
* `gateway_contact_info` - Contact information for the gateway.
* `gateway_cluster_id` - ID of the gateway cluster.
* `gateway_cluster_status` - Status of the gateway cluster.
* `public_key` - Public key information for the gateway:
  * `exponent` - Exponent of the public key.
  * `modulus` - Modulus of the public key.
<!-- /docgen -->