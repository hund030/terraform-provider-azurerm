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

func TestAccAzureRMP2sVpnGateway_basic(t *testing.T) {
	resourceName := "azurerm_p2s_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	config := testAccAzureRMP2sVpnGateway_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMP2sVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMP2sVpnGatewayExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMP2sVpnGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		p2sVpnGateway := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).network.P2sVpnGatewayClient

		resp, err := conn.Get(ctx, resourceGroup, p2sVpnGateway)
		if err != nil {
			return fmt.Errorf("Bad: Get on P2sVpnGatewayClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: P2s Vpn Gateway %q (resource group: %q) does not exist", p2sVpnGateway, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMP2sVpnGatewayDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).network.P2sVpnGatewayClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_p2s_vpn_gateway" {
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

func testAccAzureRMP2sVpnGateway_basic(rInt int, p2sVpnGatewaySuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-gateway"
  location = "%s"
}

resource "azurerm_p2s_vpn_gateway" "test" {
  name                = "testaccgateway%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  virtual_hub_id           = "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/demo1-zhijie-westus2/providers/Microsoft.Network/virtualHubs/zhijie-vh-westus2"
  p2s_vpn_server_configuration_id = "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/demo1-zhijie-westus2/providers/Microsoft.Network/virtualWans/zhijie-vw-westus2/p2sVpnServerConfigurations/zhijie-p2scfg-westus2"
  vpn_client_address_pool_prefixes = ["101.3.0.0/16"]
}
`, rInt, location, p2sVpnGatewaySuffix)
}
