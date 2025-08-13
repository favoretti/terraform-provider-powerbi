package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceEmbedToken_report(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEmbedTokenConfig_report(workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "token"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "token_id"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "expiration"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "expires_on"),
					resource.TestCheckResourceAttr("data.powerbi_embed_token.test", "type", "report"),
					resource.TestCheckResourceAttr("data.powerbi_embed_token.test", "access_level", "View"),
				),
			},
		},
	})
}

func TestAccDataSourceEmbedToken_dataset(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEmbedTokenConfig_dataset(workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "token"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "token_id"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "expiration"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "expires_on"),
					resource.TestCheckResourceAttr("data.powerbi_embed_token.test", "type", "dataset"),
					resource.TestCheckResourceAttr("data.powerbi_embed_token.test", "access_level", "Edit"),
				),
			},
		},
	})
}

func TestAccDataSourceEmbedToken_dashboard(t *testing.T) {
	workspaceName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEmbedTokenConfig_dashboard(workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "token"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "token_id"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "expiration"),
					resource.TestCheckResourceAttrSet("data.powerbi_embed_token.test", "expires_on"),
					resource.TestCheckResourceAttr("data.powerbi_embed_token.test", "type", "dashboard"),
				),
			},
		},
	})
}

func testAccDataSourceEmbedTokenConfig_report(workspaceName string) string {
	return fmt.Sprintf(`
resource "powerbi_workspace" "test" {
  name = "%s"
}

data "powerbi_embed_token" "test" {
  workspace_id = powerbi_workspace.test.id
  type         = "report"
  resource_id  = "5b218778-e7a5-4d73-8187-f10824047715"
  access_level = "View"
}
`, workspaceName)
}

func testAccDataSourceEmbedTokenConfig_dataset(workspaceName string) string {
	return fmt.Sprintf(`
resource "powerbi_workspace" "test" {
  name = "%s"
}

data "powerbi_embed_token" "test" {
  workspace_id = powerbi_workspace.test.id
  type         = "dataset"
  resource_id  = "cfafbeb1-8037-4d0c-896e-a46fb27ff229"
  access_level = "Edit"
}
`, workspaceName)
}

func testAccDataSourceEmbedTokenConfig_dashboard(workspaceName string) string {
	return fmt.Sprintf(`
resource "powerbi_workspace" "test" {
  name = "%s"
}

data "powerbi_embed_token" "test" {
  workspace_id = powerbi_workspace.test.id
  type         = "dashboard"
  resource_id  = "69ffaa6c-b36d-4d01-96f5-1ed67c64d4af"
  access_level = "View"
}
`, workspaceName)
}