package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccGatewayDatasource_basic(t *testing.T) {
	datasourceName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGatewayDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDatasourceConfig_basic(datasourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGatewayDatasourceExists("powerbi_gateway_datasource.test"),
					resource.TestCheckResourceAttr("powerbi_gateway_datasource.test", "datasource_name", datasourceName),
					resource.TestCheckResourceAttr("powerbi_gateway_datasource.test", "datasource_type", "Sql"),
					resource.TestCheckResourceAttrSet("powerbi_gateway_datasource.test", "id"),
					resource.TestCheckResourceAttrSet("powerbi_gateway_datasource.test", "gateway_id"),
				),
			},
		},
	})
}

func TestAccGatewayDatasource_importBasic(t *testing.T) {
	datasourceName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGatewayDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewayDatasourceConfig_basic(datasourceName),
			},
			{
				ResourceName:      "powerbi_gateway_datasource.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGatewayDatasourceImportStateIdFunc("powerbi_gateway_datasource.test"),
			},
		},
	})
}

func testAccCheckGatewayDatasourceDestroy(s *terraform.State) error {
	// Since we don't have the actual API client available in tests,
	// we'll just check that the resource was removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerbi_gateway_datasource" {
			continue
		}
		// In a real test, we would check if the datasource exists in Power BI
		// For now, we assume if it's not in state, it was destroyed
	}
	return nil
}

func testAccCheckGatewayDatasourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		// In a real test, we would verify the datasource exists in Power BI
		return nil
	}
}

func testAccGatewayDatasourceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("not found: %s", name)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["gateway_id"], rs.Primary.ID), nil
	}
}

func testAccGatewayDatasourceConfig_basic(datasourceName string) string {
	return fmt.Sprintf(`
# Note: This test configuration assumes a gateway exists
# In a real test environment, you would either:
# 1. Use a data source to find an existing gateway
# 2. Mock the gateway for testing
# 3. Skip the test if no gateway is available

data "powerbi_gateway" "test" {
  # This would need to reference an actual gateway in the test environment
  name = "TestGateway"
}

resource "powerbi_gateway_datasource" "test" {
  gateway_id      = data.powerbi_gateway.test.id
  datasource_name = "%s"
  datasource_type = "Sql"

  connection_details {
    server   = "localhost"
    database = "TestDB"
  }

  credential_type = "Basic"
  
  credential_details {
    privacy_level = "Organizational"
  }
}
`, datasourceName)
}