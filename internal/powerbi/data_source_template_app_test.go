package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTemplateApp_byID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTemplateAppConfig_byID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "description"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "publisher_name"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "version"),
				),
			},
		},
	})
}

func TestAccDataSourceTemplateApp_byName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTemplateAppConfig_byName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "description"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "publisher_name"),
					resource.TestCheckResourceAttrSet("data.powerbi_template_app.test", "version"),
				),
			},
		},
	})
}

func testAccDataSourceTemplateAppConfig_byID() string {
	return `
data "powerbi_template_app" "test" {
  id = "f089354e-8366-4e18-aea3-4cb4a3a50b48"
}
`
}

func testAccDataSourceTemplateAppConfig_byName() string {
	return `
data "powerbi_template_app" "test" {
  name = "Sales Marketing"
}
`
}