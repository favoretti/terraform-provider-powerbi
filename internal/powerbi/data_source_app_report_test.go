package powerbi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAppReport_byID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppReportConfig_byID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "web_url"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "embed_url"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "dataset_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAppReport_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppReportConfig_byName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "web_url"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "embed_url"),
					resource.TestCheckResourceAttrSet("data.powerbi_app_report.test", "dataset_id"),
				),
			},
		},
	})
}

func testAccDataSourceAppReportConfig_byID() string {
	return `
data "powerbi_app_report" "test" {
  app_id = "f089354e-8366-4e18-aea3-4cb4a3a50b48"
  id     = "5b218778-e7a5-4d73-8187-f10824047715"
}
`
}

func testAccDataSourceAppReportConfig_byName() string {
	return `
data "powerbi_app_report" "test" {
  app_id = "f089354e-8366-4e18-aea3-4cb4a3a50b48"
  name   = "Customer Profitability Sample"
}
`
}