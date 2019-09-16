---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_p2s_vpn_gateway"
sidebar_current: "docs-azurerm-datasource-p2s-vpn-gateway"
description: |-
  Gets information about an existing P2s Vpn Gateway
---

# Data Source: azurerm_p2s_vpn_gateway

Use this data source to access information about an existing P2s Vpn Gateway.

## Example Usage

```hcl
data "azurerm_p2s_vpn_gateway" "example" {
  resource_group_name = "example-rg"
  name                = "acctestgateway"
}

output "virtual_hub_id" {
  value = "${data.azurerm_p2s_vpn_gateway.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the gateway.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - Resource ID.

* `location` - Resource location.

* `custom_route_address_prefixes` - A list of address blocks reserved for this virtual network in CIDR notation.

* `p2svpn_server_configuration_id` - The P2SVpnServerConfiguration to which the p2sVpnGateway is attached to.

* `scale_unit` - The scale unit for this p2s vpn gateway.

* `virtual_hub_id` - The VirtualHub to which the gateway belongs.

* `vpn_client_address_pool_prefixes` - A list of address blocks reserved for this virtual network in CIDR notation.

* `vpn_client_connection_health` - One `vpn_client_connection_health` block defined below.

* `tags` - Resource tags.

---

The `vpn_client_connection_health` block contains the following:

* `total_ingress_bytes_transferred` - Total of the Ingress Bytes Transferred in this P2S Vpn connection.

* `total_egress_bytes_transferred` - Total of the Egress Bytes Transferred in this connection.

* `vpn_client_connections_count` - The total of p2s vpn clients connected at this time to this P2SVpnGateway.

* `allocated_ip_addresses` - List of allocated ip addresses to the connected p2s vpn clients.
