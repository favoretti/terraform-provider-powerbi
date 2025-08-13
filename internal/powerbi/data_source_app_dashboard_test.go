package powerbi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAppDashboard_byID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppDashboardConfig_byID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "display_name"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "embed_url"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "is_read_only"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "web_url"),
				),
			},
		},
	})
}

func TestAccDataSourceAppDashboard_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppDashboardConfig_byName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "display_name"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "embed_url"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "is_read_only"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_dashboard.test", "web_url"),
				),
			},
		},
	})
}

func testAccDataSourceAppDashboardConfig_byID() string {
	return `
data "powerbi_app_dashboard" "test" {
  app_id = "f089354e-8366-4e18-aea3-4cb4a3a50b48"
  id     = "69ffaa6c-b36d-4d01-96f5-1ed67c64d4af"
}
`
}

func testAccDataSourceAppDashboardConfig_byName() string {
	return `
data "powerbi_app_dashboard" "test" {
  app_id       = "f089354e-8366-4e18-aea3-4cb4a3a50b48"
  display_name = "Team Scorecard"
}
`
}