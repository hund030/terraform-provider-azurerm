package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func resourceArmP2sVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmP2sVpnGatewayCreateUpdate,
		Read:   resourceArmP2sVpnGatewayRead,
		Update: resourceArmP2sVpnGatewayCreateUpdate,
		Delete: resourceArmP2sVpnGatewayDelete,
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

			"p2s_vpn_server_configuration_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"custom_route_address_prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"vpn_client_address_pool_prefixes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"scale_unit": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmP2sVpnGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.P2sVpnGatewayClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM p2s vpn gateway creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing P2s Vpn Gateway %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_p2s_vpn_gateway", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	prop, err := expandP2sVpnGatewayProperties(d)
	if err != nil {
		return err
	}

	param := network.P2SVpnGateway{
		Name:                    &name,
		Location:                &location,
		P2SVpnGatewayProperties: prop,
		Tags:                    expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, param)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating P2s Vpn Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of P2s Vpn Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read P2s Vpn Gateway %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmP2sVpnGatewayRead(d, meta)
}

func resourceArmP2sVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceArmP2sVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.P2sVpnGatewayClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["p2sVpnGateways"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting P2s Vpn Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of P2s Vpn Gateways %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandP2sVpnGatewayProperties(d *schema.ResourceData) (*network.P2SVpnGatewayProperties, error) {
	virtualHubId := d.Get("virtual_hub_id").(string)
	p2sVpnServerCfg := d.Get("p2s_vpn_server_configuration_id").(string)

	prop := &network.P2SVpnGatewayProperties{
		VirtualHub: &network.SubResource{
			ID: &virtualHubId,
		},
		P2SVpnServerConfiguration: &network.SubResource{
			ID: &p2sVpnServerCfg,
		},
	}

	if v, ok := d.GetOk("scale_unit"); ok {
		scaleUnit := int32(v.(int))
		prop.VpnGatewayScaleUnit = &scaleUnit
	}

	if r, ok := d.GetOk("custom_route_address_prefixes"); ok {
		var customRoutes []string
		r := r.(*schema.Set).List()
		for _, v := range r {
			s := v.(string)
			customRoutes = append(customRoutes, s)
		}
		prop.CustomRoutes = &network.AddressSpace{
			AddressPrefixes: &customRoutes,
		}
	}

	if r, ok := d.GetOk("vpn_client_address_pool_prefixes"); ok {
		var vpnClient []string
		r := r.(*schema.Set).List()
		for _, v := range r {
			s := v.(string)
			vpnClient = append(vpnClient, s)
		}
		prop.VpnClientAddressPool = &network.AddressSpace{
			AddressPrefixes: &vpnClient,
		}
	}

	return prop, nil
}
