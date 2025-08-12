package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDataflow_byID(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	dataflowName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataflowConfig_byID(workspaceName, dataflowName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "id"),
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "name"),
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "workspace_id"),
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "configured_by"),
				),
			},
		},
	})
}

func TestAccDataSourceDataflow_byName(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	dataflowName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataflowConfig_byName(workspaceName, dataflowName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "id"),
					resource.TestCheckResourceAttr("data.powerbi_dataflow.test", "name", dataflowName),
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "workspace_id"),
					resource.TestCheckResourceAttrSet("data.powerbi_dataflow.test", "configured_by"),
				),
			},
		},
	})
}

func testAccDataSourceDataflowConfig_byID(workspaceName, dataflowName string) string {
	return fmt.Sprintf(`
resource "powerbi_workspace" "test" {
  name = "%s"
}

resource "powerbi_dataflow" "test_dataflow" {
  name         = "%s"
  workspace_id = powerbi_workspace.test.id
}

data "powerbi_dataflow" "test" {
  workspace_id = powerbi_workspace.test.id
  id           = powerbi_dataflow.test_dataflow.id
}
`, workspaceName, dataflowName)
}

func testAccDataSourceDataflowConfig_byName(workspaceName, dataflowName string) string {
	return fmt.Sprintf(`
resource "powerbi_workspace" "test" {
  name = "%s"
}

resource "powerbi_dataflow" "test_dataflow" {
  name         = "%s"
  workspace_id = powerbi_workspace.test.id
}

data "powerbi_dataflow" "test" {
  workspace_id = powerbi_workspace.test.id
  name         = "%s"
}
`, workspaceName, dataflowName, dataflowName)
}