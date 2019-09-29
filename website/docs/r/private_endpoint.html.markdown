---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint"
sidebar_current: "docs-azurerm-resource-private-endpoint"
description: |-
  Manage Azure PrivateEndpoint instance.
---

# azurerm_private_endpoint

Manage Azure PrivateEndpoint instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acctestRG"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name     = "acctestvnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                 = "testSubnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.1.0/24"
  private_link_service_network_policies = "Disabled"
  private_endpoint_network_policies = "Disabled"
}

resource "azurerm_public_ip" "example" {
  name                = "testPip"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "testLb"
  sku                 = "Standard"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  frontend_ip_configuration {
    name                 = "${azurerm_public_ip.example.name}"
    public_ip_address_id = "${azurerm_public_ip.example.id}"
  }
}

resource "azurerm_private_link_service" "example" {
  name = "testpls"
  location = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  nat_ip_configuration {
    name = "${azurerm_public_ip.example.name}"
    subnet_id = "${azurerm_subnet.example.id}"
  }
  load_balancer_frontend_ip_configuration_ids = ["${azurerm_lb.test.frontend_ip_configuration.0.id}"]
}

resource "azurerm_private_endpoint" "example" {
  name                = "testpe"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  subnet_id           = "${azurerm_subnet.example.id}"
  tags = {
    env = "test"
  }

  private_link_service_connections {
    name = "testplsconnection"
    private_link_service_id = "${azurerm_private_link_service.example.id}"
    request_message         = "Please approve my connection"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet from which the private IP will be allocated.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `manual_private_link_service_connections` - (Optional) A grouping of information about the connection to the remote resource. Used when the network admin does not have access to approve connections to the remote resource. One or more `manual_private_link_service_connection` block defined below.

* `private_link_service_connections` - (Optional) A grouping of information about the connection to the remote resource. One or more `private_link_service_connection` block defined below.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `manual_private_link_service_connection` block supports the following:

* `private_link_service_id` - (Required) The resource id of private link service.

* `group_ids` - (Optional) The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - (Optional) A message passed to the owner of the remote resource with this connection request. Restricted to 140 chars.

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

---

The `private_link_service_connection` block supports the following:

* `private_link_service_id` - (Required) The resource id of private link service.

* `group_ids` - (Optional) The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - (Optional) A message passed to the owner of the remote resource with this connection request. Restricted to 140 chars.

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.

## Attributes Reference

The following attributes are exported:

* `network_interfaces` - Gets an array of references to the network interfaces created for this private endpoint. One or more `network_interface` block defined below.

---

The `network_interface` block contains the following:

* `id` - Resource ID.

## Import

Private Endpoint can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/privateEndpoints/example-private-endpoint
```
