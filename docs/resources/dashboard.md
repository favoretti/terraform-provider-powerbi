# Dashboard Resource
`powerbi_dashboard` represents a dashboard within a Power BI workspace.

## Example Usage
```hcl
resource "powerbi_workspace" "example" {
  name = "Example Workspace"
}

resource "powerbi_dashboard" "sales" {
  name         = "Sales Dashboard"
  workspace_id = powerbi_workspace.example.id
}
```

~> **Note:** Dashboards cannot be updated after creation. Any changes to the name will force a new resource to be created.

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `name` - (Required, Forces new resource) Name of the dashboard.
* `workspace_id` - (Required, Forces new resource) ID of the workspace where the dashboard will be created.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the dashboard.
<!-- docgen:ComputedParameters -->
* `display_name` - Display name of the dashboard.
* `is_read_only` - Whether the dashboard is read-only.
* `web_url` - Web URL of the dashboard.
* `embed_url` - Embed URL of the dashboard.
<!-- /docgen -->

## Import
Dashboards can be imported using the workspace ID and dashboard ID separated by a forward slash:

```shell
terraform import powerbi_dashboard.example workspace_id/dashboard_id
```