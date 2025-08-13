package powerbi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceApp_byID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppConfig_byID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "description"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "published_by"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "last_update"),
				),
			},
		},
	})
}

func TestAccDataSourceApp_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppConfig_byName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "description"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "published_by"),
					resource.TestCheckResourceAttrSet("data.powerbi_app.test", "last_update"),
				),
			},
		},
	})
}

func testAccDataSourceAppConfig_byID() string {
	return `
data "powerbi_app" "test" {
  id = "f089354e-8366-4e18-aea3-4cb4a3a50b48"
}
`
}

func testAccDataSourceAppConfig_byName() string {
	return `
data "powerbi_app" "test" {
  name = "Sales Marketing"
}
`
}