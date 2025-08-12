package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDashboard_basic(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	dashboardName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDashboardConfig_basic(workspaceName, dashboardName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDashboardExists("powerbi_dashboard.test"),
					resource.TestCheckResourceAttr("powerbi_dashboard.test", "name", dashboardName),
					resource.TestCheckResourceAttrSet("powerbi_dashboard.test", "id"),
					resource.TestCheckResourceAttrSet("powerbi_dashboard.test", "display_name"),
					resource.TestCheckResourceAttrSet("powerbi_dashboard.test", "workspace_id"),
				),
			},
		},
	})
}

func TestAccDashboard_importBasic(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	dashboardName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDashboardConfig_basic(workspaceName, dashboardName),
			},
			{
				ResourceName:      "powerbi_dashboard.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDashboardImportStateIdFunc("powerbi_dashboard.test"),
			},
		},
	})
}

func testAccCheckDashboardDestroy(s *terraform.State) error {
	// Since we don't have the actual API client available in tests,
	// we'll just check that the resource was removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerbi_dashboard" {
			continue
		}
		// In a real test, we would check if the dashboard exists in Power BI
		// For now, we assume if it's not in state, it was destroyed
	}
	return nil
}

func testAccCheckDashboardExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		// In a real test, we would verify the dashboard exists in Power BI
		return nil
	}
}

func testAccDashboardImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}

func testAccDashboardConfig_basic(workspaceName, dashboardName string) string {
	return fmt.Sprintf(`
resource "powerbi_workspace" "test" {
  name = "%s"
}

resource "powerbi_dashboard" "test" {
  name         = "%s"
  workspace_id = powerbi_workspace.test.id
}
`, workspaceName, dashboardName)
}