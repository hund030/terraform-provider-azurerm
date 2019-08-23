package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateEndpoint_basic(t *testing.T) {
	dataSourceName := "data.azurerm_private_endpoint.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()
	config := testAccDataSourceAzureRMPrivateEndpoint_basic(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("testaccendpoint%s", rs)),
					resource.TestCheckResourceAttr(dataSourceName, "location", azure.NormalizeLocation(location)),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPrivateEndpoint_basic(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d-endpoint"
  location = "%s"
}

resource "azurerm_private_endpoint" "test" {
  name                = "testaccendpoint%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  subnet_id           = "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/demo2-zhijie-westus2/providers/Microsoft.Network/virtualNetworks/zhijie-vnet/subnets/default"
  manual_private_link_service_connections {
    name = "plsConnection"
    private_link_service_id = "/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/demo2-zhijie-westus2/providers/Microsoft.Network/privateLinkServices/zhijie-pls"
  }
}

data "azurerm_private_endpoint" "test" {
	name				= "${azurerm_private_endpoint.test.name}"
	resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rString)
}
