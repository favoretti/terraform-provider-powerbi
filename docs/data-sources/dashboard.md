# Dashboard Data Source
`powerbi_dashboard` represents a dashboard within a Power BI workspace.

## Example Usage

### Find dashboard by name
```hcl
data "powerbi_workspace" "example" {
  name = "Example Workspace"
}

data "powerbi_dashboard" "sales" {
  workspace_id = data.powerbi_workspace.example.id
  name         = "Sales Dashboard"
}

output "dashboard_url" {
  value = data.powerbi_dashboard.sales.web_url
}
```

### Find dashboard by ID
```hcl
data "powerbi_dashboard" "existing" {
  workspace_id = "workspace-12345"
  id           = "dashboard-67890"
}
```

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `workspace_id` - (Required) ID of the workspace containing the dashboard.
* `id` - (Optional, Conflicts with: `name`) ID of the dashboard.
* `name` - (Optional, Conflicts with: `id`) Name of the dashboard.
<!-- /docgen -->

~> **Note:** You must specify exactly one of `id` or `name` to identify the dashboard.

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
<!-- docgen:ComputedParameters -->
* `display_name` - Display name of the dashboard.
* `is_read_only` - Whether the dashboard is read-only.
* `web_url` - Web URL of the dashboard.
* `embed_url` - Embed URL of the dashboard.
<!-- /docgen -->