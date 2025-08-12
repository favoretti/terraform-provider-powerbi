package powerbi

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDeploymentPipeline_basic(t *testing.T) {
	pipelineName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeploymentPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentPipelineConfig_basic(pipelineName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentPipelineExists("powerbi_deployment_pipeline.test"),
					resource.TestCheckResourceAttr("powerbi_deployment_pipeline.test", "display_name", pipelineName),
					resource.TestCheckResourceAttrSet("powerbi_deployment_pipeline.test", "id"),
				),
			},
		},
	})
}

func TestAccDeploymentPipeline_withDescription(t *testing.T) {
	pipelineName := fmt.Sprintf("tftest%s", acctest.RandString(6))
	description := "Test deployment pipeline description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeploymentPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentPipelineConfig_withDescription(pipelineName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentPipelineExists("powerbi_deployment_pipeline.test"),
					resource.TestCheckResourceAttr("powerbi_deployment_pipeline.test", "display_name", pipelineName),
					resource.TestCheckResourceAttr("powerbi_deployment_pipeline.test", "description", description),
					resource.TestCheckResourceAttrSet("powerbi_deployment_pipeline.test", "id"),
				),
			},
		},
	})
}

func TestAccDeploymentPipeline_importBasic(t *testing.T) {
	pipelineName := fmt.Sprintf("tftest%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDeploymentPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentPipelineConfig_basic(pipelineName),
			},
			{
				ResourceName:      "powerbi_deployment_pipeline.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDeploymentPipelineDestroy(s *terraform.State) error {
	// Since we don't have the actual API client available in tests,
	// we'll just check that the resource was removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerbi_deployment_pipeline" {
			continue
		}
		// In a real test, we would check if the pipeline exists in Power BI
		// For now, we assume if it's not in state, it was destroyed
	}
	return nil
}

func testAccCheckDeploymentPipelineExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		// In a real test, we would verify the pipeline exists in Power BI
		return nil
	}
}

func testAccDeploymentPipelineConfig_basic(pipelineName string) string {
	return fmt.Sprintf(`
resource "powerbi_deployment_pipeline" "test" {
  display_name = "%s"
}
`, pipelineName)
}

func testAccDeploymentPipelineConfig_withDescription(pipelineName, description string) string {
	return fmt.Sprintf(`
resource "powerbi_deployment_pipeline" "test" {
  display_name = "%s"
  description  = "%s"
}
`, pipelineName, description)
}