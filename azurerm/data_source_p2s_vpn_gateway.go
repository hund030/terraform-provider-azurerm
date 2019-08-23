package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmP2sVpnGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmP2sVpnGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"p2s_vpn_server_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_hub_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"custom_route_address_prefixes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"vpn_client_address_pool_prefixes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"scale_unit": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmP2sVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.P2sVpnGatewayClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["p2sVpnGateways"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on P2s Vpn Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	prop := resp.P2SVpnGatewayProperties
	if prop == nil {
		return fmt.Errorf("Error reading P2SVpnGatewayProperties")
	}

	if err := d.Set("p2s_vpn_server_configuration_id", prop.P2SVpnServerConfiguration.ID); err != nil {
		return fmt.Errorf("Error setting `p2s_vpn_server_configuration_id`: %+v", err)
	}

	if err := d.Set("virtual_hub_id", prop.VirtualHub.ID); err != nil {
		return fmt.Errorf("Error setting `virtual_hub_id`: %+v", err)
	}

	if addr := prop.CustomRoutes.AddressPrefixes; addr != nil {
		if err := d.Set("custom_route_address_prefixes", addr); err != nil {
			return fmt.Errorf("Error setting `custom_route_address_prefixes`: %+v", err)
		}
	}

	if addr := prop.VpnClientAddressPool.AddressPrefixes; addr != nil {
		if err := d.Set("vpn_client_address_pool_prefixes", addr); err != nil {
			return fmt.Errorf("Error setting `vpn_client_address_pool_prefixes`: %+v", err)
		}
	}

	if prop.VpnGatewayScaleUnit != nil {
		if err := d.Set("scale_unit", prop.VpnGatewayScaleUnit); err != nil {
			return fmt.Errorf("Error setting `scale_unit`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
