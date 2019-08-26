---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint"
sidebar_current: "docs-azurerm-resource-private_endpoint"
description: |-
  Manages an Azure Private Endpoint

---

# azurerm_private_endpoint

Manages an Azure Private Endpoint

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_private_endpoint" "test" {
  name                = "testprivateendpoint"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  subnet_id           = "${azurerm_subnet.test.id}"

  private_link_service_connections {
    name                    = "plsConnection"
    private_link_service_id = "${azurerm_private_link_service.test.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Private endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Private endpoint. Changing this forces a new resource to be created.

~> **NOTE:** To work around [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/5574) this property is currently treated as case-insensitive. A future version of Terraform will require that the casing is correct.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet from which the private IP will be allocated.

* `private_link_service_connections` - (Optional) A grouping of information about the connection to the remote resource.

* `manual_private_link_service_connections` - (Optional) A grouping of information about the connection to the remote resource. Used when the network admin does not have access to approve connections to the remote resource.

~> **NOTE:** Either `private_link_service_connections` or `manual_private_link_service_connections` should be set.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

Both `private_link_service_connections` block and `manual_private_link_service_connections` block support the following:

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `private_link_service_id` - (Required) The resource id of private link service.

* `group_ids` - (Optional) The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - (Optional) A message passed to the owner of the remote resource with this connection request. Restricted to 140 chars.

---

## Attributes Reference

The following attributes are exported:

* `id` - The Private Endpoint ID.
