package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMP2sVpnGateway_basic(t *testing.T) {
	dataSourceName := "data.azurerm_p2s_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMP2sVpnGateway_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("testaccgateway%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "location", azure.NormalizeLocation(location)),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMP2sVpnGateway_basic(rInt int, rString string, location string) string {
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

data "azurerm_p2s_vpn_gateway" "test" {
	name				= "${azurerm_p2s_vpn_gateway.test.name}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rString)
}
