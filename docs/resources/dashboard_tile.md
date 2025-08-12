# Dashboard Tile Resource
`powerbi_dashboard_tile` represents a tile cloned from one dashboard to another within Power BI workspaces.

## Example Usage

### Basic tile cloning
```hcl
resource "powerbi_workspace" "source" {
  name = "Source Workspace"
}

resource "powerbi_workspace" "target" {
  name = "Target Workspace"
}

resource "powerbi_dashboard" "source_dashboard" {
  name         = "Source Dashboard"
  workspace_id = powerbi_workspace.source.id
}

resource "powerbi_dashboard" "target_dashboard" {
  name         = "Target Dashboard"
  workspace_id = powerbi_workspace.target.id
}

resource "powerbi_dashboard_tile" "cloned_tile" {
  workspace_id            = powerbi_workspace.source.id
  dashboard_id           = powerbi_dashboard.target_dashboard.id
  source_dashboard_id    = powerbi_dashboard.source_dashboard.id
  source_tile_id         = "tile-12345"
  target_workspace_id    = powerbi_workspace.target.id
  position_conflict_action = "Tail"
}
```

### Cross-workspace tile cloning with rebinding
```hcl
resource "powerbi_dashboard_tile" "rebounded_tile" {
  workspace_id            = powerbi_workspace.source.id
  dashboard_id           = powerbi_dashboard.target_dashboard.id
  source_dashboard_id    = powerbi_dashboard.source_dashboard.id
  source_tile_id         = "tile-12345"
  target_workspace_id    = powerbi_workspace.target.id
  target_report_id       = "new-report-id"
  target_model_id        = "new-model-id"
  position_conflict_action = "Abort"
}
```

~> **Note:** Dashboard tiles cannot be updated after creation. Any changes will force a new resource to be created.

~> **Power BI API Limitation:** The Power BI API does not provide a direct way to delete individual tiles. When this resource is destroyed, it will only be removed from Terraform state.

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `workspace_id` - (Required, Forces new resource) ID of the workspace containing the dashboard.
* `dashboard_id` - (Required, Forces new resource) ID of the dashboard to add the tile to.
* `source_dashboard_id` - (Required, Forces new resource) ID of the source dashboard to clone tile from.
* `source_tile_id` - (Required, Forces new resource) ID of the source tile to clone.
* `target_workspace_id` - (Optional, Forces new resource) ID of the target workspace (if different from source).
* `target_report_id` - (Optional, Forces new resource) ID of the target report (if rebinding to different report).
* `target_model_id` - (Optional, Forces new resource) ID of the target model (if rebinding to different model).
* `position_conflict_action` - (Optional, Default: `Tail`, Forces new resource) Action to take if tile position conflicts. Options: `Tail` or `Abort`.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported in addition to the arguments listed above:
* `id` - The ID of the tile.
<!-- docgen:ComputedParameters -->
* `title` - Title of the tile.
* `subtitle` - Subtitle of the tile.
* `embed_url` - Embed URL of the tile.
* `embed_data` - Embed data of the tile.
* `report_id` - Report ID associated with the tile.
* `dataset_id` - Dataset ID associated with the tile.
* `row_span` - Number of rows the tile spans.
* `col_span` - Number of columns the tile spans.
<!-- /docgen -->

## Import
Dashboard tiles can be imported using the workspace ID, dashboard ID, and tile ID separated by forward slashes:

```shell
terraform import powerbi_dashboard_tile.example workspace_id/dashboard_id/tile_id
```