package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMP2sVpnGateway_basic(t *testing.T) {
	dataSourceName := "data.azurerm_p2s_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceP2sVpnGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_hub_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "p2svpn_server_configuration_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.0", "101.3.0.0/16"),
				),
			},
		},
	})
}
func TestAccDataSourceAzureRMP2sVpnGateway_complete(t *testing.T) {
	dataSourceName := "data.azurerm_p2s_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceP2sVpnGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_hub_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "p2svpn_server_configuration_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.0", "101.3.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "custom_route_address_prefixes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "custom_route_address_prefixes.0", "101.168.0.6/32"),
					resource.TestCheckResourceAttr(dataSourceName, "scale_unit", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.env", "test"),
				),
			},
		},
	})
}
func TestAccDataSourceAzureRMP2sVpnGateway_update(t *testing.T) {
	dataSourceName := "data.azurerm_p2s_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceP2sVpnGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_hub_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "p2svpn_server_configuration_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.0", "101.3.0.0/16"),
				),
			},
			{
				Config: testAccDataSourceP2sVpnGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_hub_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "p2svpn_server_configuration_id"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_client_address_pool_prefixes.0", "101.3.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "custom_route_address_prefixes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "custom_route_address_prefixes.0", "101.168.0.6/32"),
					resource.TestCheckResourceAttr(dataSourceName, "scale_unit", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceP2sVpnGateway_basic(rInt int, location string) string {
	config := testAccAzureRMP2sVpnGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_p2s_vpn_gateway" "test" {
  resource_group_name = "${azurerm_p2s_vpn_gateway.test.resource_group_name}"
  name                = "${azurerm_p2s_vpn_gateway.test.name}"
}
`, config)
}

func testAccDataSourceP2sVpnGateway_complete(rInt int, location string) string {
	config := testAccAzureRMP2sVpnGateway_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_p2s_vpn_gateway" "test" {
  resource_group_name = "${azurerm_p2s_vpn_gateway.test.resource_group_name}"
  name                = "${azurerm_p2s_vpn_gateway.test.name}"
}
`, config)
}
