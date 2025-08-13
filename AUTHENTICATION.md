# Power BI Terraform Provider Authentication Guide

The Power BI Terraform provider supports multiple authentication methods to provide flexibility for different deployment scenarios and security requirements.

## Authentication Methods

### 1. Service Principal with Client Secret (Recommended)

The most common and recommended method for production environments.

**Requirements:**
- Azure AD App Registration
- Client Secret
- Power BI Service Admin permissions

**Configuration:**

```hcl
provider "powerbi" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id" 
  client_secret = "your-client-secret"
}
```

**Environment Variables:**
```bash
export POWERBI_TENANT_ID="your-tenant-id"
export POWERBI_CLIENT_ID="your-client-id"
export POWERBI_CLIENT_SECRET="your-client-secret"
```

### 2. Service Principal with Certificate (Most Secure)

Certificate-based authentication provides enhanced security compared to client secrets.

**Requirements:**
- Azure AD App Registration with certificate configured
- PEM format certificate with private key
- Power BI Service Admin permissions

**Configuration:**

```hcl
provider "powerbi" {
  tenant_id         = "your-tenant-id"
  client_id         = "your-client-id"
  certificate_path  = "/path/to/certificate.pem"
  certificate_password = "cert-password" # if certificate is encrypted
}
```

**Alternatively, use base64 encoded certificate data:**

```hcl
provider "powerbi" {
  tenant_id         = "your-tenant-id"
  client_id         = "your-client-id"
  certificate_data  = "LS0tLS1CRUdJTi..." # base64 encoded certificate
}
```

**Environment Variables:**
```bash
export POWERBI_TENANT_ID="your-tenant-id"
export POWERBI_CLIENT_ID="your-client-id"
export POWERBI_CERTIFICATE_PATH="/path/to/certificate.pem"
export POWERBI_CERTIFICATE_PASSWORD="cert-password"
# OR
export POWERBI_CERTIFICATE_DATA="LS0tLS1CRUdJTi..."
```

### 3. Managed Identity (Azure Environments)

Ideal for resources running in Azure (VMs, App Service, Functions, etc.).

**Requirements:**
- Running in Azure environment
- Managed Identity enabled
- Power BI Service Admin permissions assigned to the managed identity

**Configuration:**

```hcl
provider "powerbi" {
  use_managed_identity = true
}
```

**For User-Assigned Managed Identity:**

```hcl
provider "powerbi" {
  use_managed_identity  = true
  managed_identity_id   = "your-managed-identity-client-id"
}
```

**Environment Variables:**
```bash
export POWERBI_USE_MANAGED_IDENTITY=true
export POWERBI_MANAGED_IDENTITY_ID="your-managed-identity-client-id" # optional
```

### 4. Azure CLI Authentication

Perfect for local development and testing.

**Requirements:**
- Azure CLI installed and authenticated (`az login`)
- Current user has Power BI permissions

**Configuration:**

```hcl
provider "powerbi" {
  use_azure_cli = true
}
```

**Environment Variables:**
```bash
export POWERBI_USE_AZURE_CLI=true
```

### 5. Direct Access Token

Use a pre-obtained access token directly.

**Requirements:**
- Valid Power BI access token
- Token must have appropriate scopes

**Configuration:**

```hcl
provider "powerbi" {
  access_token = "your-access-token"
}
```

**Environment Variables:**
```bash
export POWERBI_ACCESS_TOKEN="your-access-token"
```

### 6. Username/Password (Deprecated)

⚠️ **Deprecated:** This method is included for backward compatibility but is not recommended for production use.

**Configuration:**

```hcl
provider "powerbi" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  username      = "user@domain.com"
  password      = "user-password"
}
```

## Authentication Priority

When multiple authentication methods are configured, the provider uses the following priority order:

1. **Direct Access Token** (`access_token`)
2. **Managed Identity** (`use_managed_identity`)
3. **Azure CLI** (`use_azure_cli`)
4. **Certificate Authentication** (`certificate_path` or `certificate_data`)
5. **Client Secret** (`client_secret`)
6. **Username/Password** (deprecated)

## Setting up Service Principal Authentication

### Step 1: Create Azure AD App Registration

1. Navigate to [Azure Portal](https://portal.azure.com)
2. Go to **Azure Active Directory** > **App registrations**
3. Click **New registration**
4. Enter a name for your application
5. Select appropriate account types
6. Click **Register**

### Step 2: Configure Authentication Method

#### For Client Secret:
1. Go to **Certificates & secrets** in your app registration
2. Click **New client secret**
3. Enter description and set expiration
4. Copy the secret value (you won't be able to see it again)

#### For Certificate:
1. Generate a certificate:
   ```bash
   # Generate private key
   openssl genrsa -out powerbi.key 2048
   
   # Generate certificate signing request
   openssl req -new -key powerbi.key -out powerbi.csr
   
   # Generate self-signed certificate (valid for 1 year)
   openssl x509 -req -days 365 -in powerbi.csr -signkey powerbi.key -out powerbi.crt
   
   # Combine certificate and private key into PEM file
   cat powerbi.crt powerbi.key > powerbi.pem
   ```

2. Upload certificate to Azure AD:
   - Go to **Certificates & secrets** in your app registration
   - Click **Upload certificate**
   - Upload the `.crt` file

### Step 3: Configure Power BI Permissions

1. Go to [Power BI Admin Portal](https://app.powerbi.com/admin-portal)
2. Navigate to **Tenant settings**
3. Find **Developer settings**
4. Enable **Allow service principals to use Power BI APIs**
5. Add your service principal to the allowed security groups

### Step 4: Grant Workspace Access

Service principals need explicit access to workspaces:

1. Go to the Power BI workspace
2. Click **Access** 
3. Add your service principal as **Admin** or **Member**

## Environment Variables Reference

| Variable | Description | Required |
|----------|-------------|----------|
| `POWERBI_TENANT_ID` | Azure AD Tenant ID | For most auth methods |
| `POWERBI_CLIENT_ID` | Application (Client) ID | For most auth methods |
| `POWERBI_CLIENT_SECRET` | Client Secret | For client secret auth |
| `POWERBI_CERTIFICATE_PATH` | Path to certificate file | For certificate auth |
| `POWERBI_CERTIFICATE_DATA` | Base64 encoded certificate | For certificate auth |
| `POWERBI_CERTIFICATE_PASSWORD` | Certificate password | If cert is encrypted |
| `POWERBI_USE_MANAGED_IDENTITY` | Enable managed identity | For managed identity |
| `POWERBI_MANAGED_IDENTITY_ID` | User assigned MI client ID | For user-assigned MI |
| `POWERBI_USE_AZURE_CLI` | Enable Azure CLI auth | For CLI auth |
| `POWERBI_ACCESS_TOKEN` | Pre-obtained access token | For direct token auth |
| `POWERBI_USERNAME` | Username (deprecated) | For password auth |
| `POWERBI_PASSWORD` | Password (deprecated) | For password auth |

## Troubleshooting

### Common Issues

#### "Service principal not enabled"
- Enable service principals in Power BI Admin Portal
- Add service principal to allowed security groups

#### "Insufficient privileges"
- Grant service principal workspace access
- Ensure service principal has required Power BI permissions

#### "Certificate validation failed"
- Verify certificate format (PEM required)
- Check certificate expiration
- Ensure private key is included

#### "Managed identity not found"
- Verify managed identity is enabled in Azure resource
- Check managed identity has Power BI permissions

#### "Azure CLI not authenticated"
- Run `az login` to authenticate
- Verify current account has Power BI access

### Debug Tips

1. **Enable verbose logging:**
   ```bash
   export TF_LOG=DEBUG
   ```

2. **Test authentication separately:**
   ```bash
   # Test Azure CLI
   az account get-access-token --resource https://analysis.windows.net/powerbi/api
   
   # Test managed identity (from Azure resource)
   curl -H "Metadata:true" "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://analysis.windows.net/powerbi/api"
   ```

3. **Verify permissions:**
   - Use Power BI REST API directly to test permissions
   - Check workspace membership
   - Verify tenant settings

## Security Best Practices

1. **Use Certificate Authentication** over client secrets when possible
2. **Store secrets securely** (Azure Key Vault, HashiCorp Vault, etc.)
3. **Rotate credentials regularly**
4. **Use least privilege principle** - grant minimum required permissions
5. **Monitor access logs** in Azure AD and Power BI
6. **Use Managed Identity** in Azure environments
7. **Avoid hardcoding credentials** in Terraform configurations

## Example Configurations

### Production Environment (Certificate-based)

```hcl
provider "powerbi" {
  tenant_id            = var.tenant_id
  client_id            = var.client_id
  certificate_path     = var.certificate_path
  certificate_password = var.certificate_password
}
```

### Development Environment (Azure CLI)

```hcl
provider "powerbi" {
  use_azure_cli = true
}
```

### CI/CD Pipeline (Client Secret)

```hcl
provider "powerbi" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
}
```

### Azure Environment (Managed Identity)

```hcl
provider "powerbi" {
  use_managed_identity = true
}
```

For more information, see the [Microsoft Power BI REST API documentation](https://docs.microsoft.com/en-us/rest/api/power-bi/).