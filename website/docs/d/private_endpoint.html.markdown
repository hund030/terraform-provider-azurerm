---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint"
sidebar_current: "docs-azurerm-datasource-private_endpoint"
description: |-
  Get information about an existing Private Endpoint

---

# Data Source: azurerm_private_endpoint

Use this data source to access information about an existing Batch Account.

## Example Usage

```hcl
data "azurerm_batch_account" "test" {
    name                = "testprivateendpoint"
    resource_group_name = "test"
}

output "" {
    value = ""
}
```

## Argument Reference

* `name` - (Required) The name of the Private Endpoint.

* `resource_group_name` - (Required) The Name of the Resource Group where this Private endpoint exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Private endpoint ID.

* `name` - The Private endpoint name.

* `location` - The Azure Region in which this Private endpoint exists.

* `subnet_id` - The ID of the subnet from which the private IP will be allocated.

* `private_link_service_connections` - A grouping of information about the connection to the remote resource.

* `manual_private_link_service_connections` - A grouping of information about the connection to the remote resource. Used when the network admin does not have access to approve connections to the remote resource.

* `tags` - A map of tags assigned to the Private endpoint.

---

Both `private_link_service_connections` block and `manual_private_link_service_connections` block support the following:

* `private_link_service_id` - The resource id of private link service.

* `group_ids` - The ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.

* `request_message` - A message passed to the owner of the remote resource with this connection request. Restricted to 140 chars.

---
