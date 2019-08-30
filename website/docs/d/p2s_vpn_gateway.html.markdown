---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_p2s_vpn_gateway"
sidebar_current: "docs-azurerm-datasource-p2s_vpn_gateway"
description: |-
  Get information about an existing P2s vpn gateway

---

# Data Source: azurerm_private_endpoint

Use this data source to access information about an existing P2s vpn gateway.

## Example Usage

```hcl
data "azurerm_p2s_vpn_gateway" "test" {
    name                = "testp2svpngateway"
    resource_group_name = "test"
}

output "virtual_hub_id" {
    value = "${data.azurerm_p2s_vpn_gateway.test.virtual_hub_id}"
}
```

## Argument Reference

* `name` - (Required) The name of the P2s vpn gateway.

* `resource_group_name` - (Required) The Name of the Resource Group where this P2s vpn gateway exists.

## Attributes Reference

The following attributes are exported:

* `id` - The P2s vpn gateway ID.

* `name` - The P2s vpn gateway name.

* `location` - The Azure Region in which this P2s vpn gateway exists.

* `tags` - A map of tags assigned to the P2s vpn gateway.

* `p2s_vpn_server_configuration_id` - The ID of the P2SVpnServerConfiguration to which the p2sVpnGateway is attached to.

* `virtual_hub_id` - The ID of the VirtualHub to which the gateway belongs.

* `vpn_client_address_pool_prefixes` - The reference of the address space resource which represents Address space for P2S VpnClient.

* `custom_route_address_prefixes` - The reference of the address space resource which represents the custom routes specified by the customer for P2SVpnGateway and P2S VpnClient.

* `scale_unit` - The scale unit for this p2s vpn gateway.
