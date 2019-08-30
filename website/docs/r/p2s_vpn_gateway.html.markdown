---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_p2s_vpn_gateway"
sidebar_current: "docs-azurerm-resource-p2s-vpn-gateway"
description: |-
  Manages an Azure P2s vpn gateway

---

# azurerm_p2s_vpn_gateway

Manages an Azure Batch account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_p2s_vpn_gateway" "p2sgw" {
  name                  = "testp2svpngateway"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"

  virtual_hub_id           = "${azurerm_virtual_hub.test.id}"
  p2s_vpn_server_configuration_id = "${azurerm_p2s_vpn_server_configuration.test.id}"
  vpn_client_address_pool_prefixes = ["101.3.0.0/16"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the P2s vpn gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the P2s vpn gateway. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A map of tags assigned to the P2s vpn gateway.

* `p2s_vpn_server_configuration_id` - (Required) The ID of the P2SVpnServerConfiguration to which the p2sVpnGateway is attached to.

* `virtual_hub_id` - (Required) The ID of the VirtualHub to which the gateway belongs.

* `vpn_client_address_pool_prefixes` - (Required) The reference of the address space resource which represents Address space for P2S VpnClient.

* `custom_route_address_prefixes` - (Optional) The reference of the address space resource which represents the custom routes specified by the customer for P2SVpnGateway and P2S VpnClient.

* `scale_unit` - (Optional) The scale unit for this p2s vpn gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The P2s vpn gateway ID.
