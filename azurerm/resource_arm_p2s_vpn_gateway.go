package azurerm

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func resourceArmP2SVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateEndpointCreateUpdate,
		Read:   resourceArmPrivateEndpointRead,
		Update: resourceArmPrivateEndpointCreateUpdate,
		Delete: resourceArmPrivateEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),
		},
	}
}
