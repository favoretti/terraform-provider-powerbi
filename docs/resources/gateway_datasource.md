# Gateway Datasource Resource
`powerbi_gateway_datasource` represents a datasource within a Power BI gateway.

## Example Usage

### SQL Server datasource
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

  credential_type = "Windows"
  
  credential_details {
    privacy_level           = "Organizational"
    use_caller_aad_identity = true
  }
}
```

### Web API datasource
```hcl
resource "powerbi_gateway_datasource" "web_api" {
  gateway_id      = data.powerbi_gateway.enterprise.id
  datasource_name = "External API"
  datasource_type = "Web"

  connection_details {
    url = "https://api.example.com/v1"
  }

  credential_type = "OAuth2"
  
  credential_details {
    privacy_level = "Public"
  }
}
```

### Oracle datasource with basic authentication
```hcl
resource "powerbi_gateway_datasource" "oracle" {
  gateway_id      = data.powerbi_gateway.enterprise.id
  datasource_name = "Oracle DW"
  datasource_type = "Oracle"

  connection_details {
    server   = "oracle.company.com:1521"
    database = "DWPROD"
  }

  credential_type = "Basic"
  
  credential_details {
    credentials              = "encrypted_credentials_string"
    encryption_algorithm     = "RSA-OAEP"
    privacy_level           = "Organizational"
  }
}
```

~> **Security Note:** Credential details are sensitive and will be encrypted by Power BI. The `credentials` field should contain properly encrypted credentials according to the gateway's public key.

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `gateway_id` - (Required, Forces new resource) ID of the gateway.
* `datasource_name` - (Required, Forces new resource) Name of the datasource.
* `datasource_type` - (Required, Forces new resource) Type of the datasource. Supported types: `Sql`, `Oracle`, `OleDb`, `ODBC`, `SharePointList`, `Web`, `OData`, `File`, `Folder`, `SharePointDocumentLibrary`, `Hdfs`, `AzureTable`, `Exchange`, `ActiveDirectory`, `MySql`, `PostgreSql`, `Sybase`, `DB2`, `Teradata`, `SapHana`, `SapBw`, `AnalysisServices`, `AzureBlob`, `AzureSql`, `AzureSqlDw`, `Informix`, `GoogleAnalytics`, `AmazonRedshift`, `Impala`, `Spark`, `Smartsheet`.
* `connection_details` - (Required) Connection details for the datasource. A [`connection_details`](#a-connection_details-block-supports-the-following) block is defined below.
* `credential_type` - (Optional, Default: `Basic`) Type of credentials used for authentication. Options: `Basic`, `Windows`, `OAuth2`, `Anonymous`, `Key`.
* `credential_details` - (Optional) Credential details for datasource authentication. A [`credential_details`](#a-credential_details-block-supports-the-following) block is defined below.

---

#### A `connection_details` block supports the following:
* `server` - (Optional) Server name or address.
* `database` - (Optional) Database name.
* `url` - (Optional) URL for web-based datasources.
* `path` - (Optional) File path for file-based datasources.
* `kind` - (Optional) Kind of datasource.
* `auth_method` - (Optional) Authentication method.
* `account` - (Optional) Account name.
* `domain` - (Optional) Domain name.
* `email_address` - (Optional) Email address for authentication.
* `login_server` - (Optional) Login server.
* `class` - (Optional) Class of the datasource.

---

#### A `credential_details` block supports the following:
* `credentials` - (Optional) Encrypted credentials.
* `encrypted_connection` - (Optional) Whether to use encrypted connection.
* `encryption_algorithm` - (Optional) Encryption algorithm used.
* `privacy_level` - (Optional, Default: `None`) Privacy level for the datasource. Options: `None`, `Public`, `Organizational`, `Private`.
* `use_caller_aad_identity` - (Optional, Default: `false`) Whether to use caller's AAD identity.
* `use_end_user_oauth2_credentials` - (Optional, Default: `false`) Whether to use end user OAuth2 credentials.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the datasource.
<!-- docgen:ComputedParameters -->
* `connection_string` - Connection string for the datasource.
<!-- /docgen -->

## Import
Gateway datasources can be imported using the gateway ID and datasource ID separated by a forward slash:

```shell
terraform import powerbi_gateway_datasource.example gateway_id/datasource_id
```