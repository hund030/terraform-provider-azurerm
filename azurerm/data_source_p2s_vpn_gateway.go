package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmP2sVpnGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmP2sVpnGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"custom_route_address_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"p2svpn_server_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scale_unit": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),

			"virtual_hub_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpn_client_address_pool_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"vpn_client_connection_health": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocated_ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"total_egress_bytes_transferred": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_ingress_bytes_transferred": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpn_client_connections_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmP2sVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.P2sVpnGatewayClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: P2s Vpn Gateway %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading P2s Vpn Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if p2SVpnGatewayProperties := resp.P2SVpnGatewayProperties; p2SVpnGatewayProperties != nil {
		if customRoutes := p2SVpnGatewayProperties.CustomRoutes; customRoutes != nil {
			d.Set("custom_route_address_prefixes", utils.FlattenStringSlice(customRoutes.AddressPrefixes))
		}
		if p2SVpnServerConfiguration := p2SVpnGatewayProperties.P2SVpnServerConfiguration; p2SVpnServerConfiguration != nil {
			d.Set("p2svpn_server_configuration_id", p2SVpnServerConfiguration.ID)
		}
		d.Set("scale_unit", int(*p2SVpnGatewayProperties.VpnGatewayScaleUnit))
		if virtualHub := p2SVpnGatewayProperties.VirtualHub; virtualHub != nil {
			d.Set("virtual_hub_id", virtualHub.ID)
		}
		if vpnClientAddressPool := p2SVpnGatewayProperties.VpnClientAddressPool; vpnClientAddressPool != nil {
			d.Set("vpn_client_address_pool_prefixes", utils.FlattenStringSlice(vpnClientAddressPool.AddressPrefixes))
		}
		if err := d.Set("vpn_client_connection_health", flattenArmP2sVpnGatewayVpnClientConnectionHealth(p2SVpnGatewayProperties.VpnClientConnectionHealth)); err != nil {
			return fmt.Errorf("Error setting `vpn_client_connection_health`: %+v", err)
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}
