# Dashboard Tiles Data Source
`powerbi_dashboard_tiles` represents all tiles within a Power BI dashboard.

## Example Usage

### Get all tiles from a dashboard
```hcl
data "powerbi_workspace" "example" {
  name = "Example Workspace"
}

data "powerbi_dashboard" "sales" {
  workspace_id = data.powerbi_workspace.example.id
  name         = "Sales Dashboard"
}

data "powerbi_dashboard_tiles" "sales_tiles" {
  workspace_id  = data.powerbi_workspace.example.id
  dashboard_id  = data.powerbi_dashboard.sales.id
}

output "tile_count" {
  value = length(data.powerbi_dashboard_tiles.sales_tiles.tiles)
}

output "tile_titles" {
  value = [for tile in data.powerbi_dashboard_tiles.sales_tiles.tiles : tile.title]
}
```

### Use tiles data to create conditional resources
```hcl
locals {
  revenue_tiles = [
    for tile in data.powerbi_dashboard_tiles.sales_tiles.tiles : tile
    if can(regex("(?i)revenue", tile.title))
  ]
}

output "revenue_tile_ids" {
  value = [for tile in local.revenue_tiles : tile.id]
}
```

## Argument Reference
#### The following arguments are supported:
<!-- docgen:NonComputedParameters -->
* `workspace_id` - (Required) ID of the workspace containing the dashboard.
* `dashboard_id` - (Required) ID of the dashboard.
<!-- /docgen -->

## Attributes Reference
#### The following attributes are exported:
<!-- docgen:ComputedParameters -->
* `tiles` - List of tiles in the dashboard. Each tile contains the following attributes:
  * `id` - ID of the tile.
  * `title` - Title of the tile.
  * `subtitle` - Subtitle of the tile.
  * `embed_url` - Embed URL of the tile.
  * `embed_data` - Embed data of the tile.
  * `report_id` - Report ID associated with the tile.
  * `dataset_id` - Dataset ID associated with the tile.
  * `row_span` - Number of rows the tile spans.
  * `col_span` - Number of columns the tile spans.
  * `configuration` - Configuration of the tile (if available).
<!-- /docgen -->