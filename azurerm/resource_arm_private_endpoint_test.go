package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateEndpoint_basic(t *testing.T) {
	resourceName := "azurerm_private_endpoint.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	config := testAccAzureRMPrivateEndpoint_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateEndpointExists(resourceName),
					// resource.TestCheckResourceAttr(resourceName, ""),
				),
			},
		},
	})
}

func testCheckAzureRMPrivateEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		privateEndpoint := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).network.PrivateEndpointClient
		fmt.Println(privateEndpoint)

		resp, err := conn.Get(ctx, resourceGroup, privateEndpoint, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on privateEndpointClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private Endpoint %q (resource group: %q) does not exist", privateEndpoint, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateEndpointDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).batch.AccountClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_endpoint" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMPrivateEndpoint_basic(rInt int, privateEndpointSuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-privateendpoint"
  location = "%s"
}

resource "azurerm_private_endpoint" "test" {
  name                = "testaccendpoint%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  }
`, rInt, location, privateEndpointSuffix)
}
