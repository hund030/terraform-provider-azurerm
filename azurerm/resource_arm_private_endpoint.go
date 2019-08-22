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

func resourceArmPrivateEndpoint() *schema.Resource {
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

			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"private_link_service_connections": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"private_link_service_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 0,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},

						"request_message": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Please approve my connection.",
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"manual_private_link_service_connections": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"private_link_service_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},

						"request_message": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Please approve my connection.",
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPrivateEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM private endpoint creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private Endpoint %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_endpoint", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	peProperties, pePropsErr := expandPrivateEndpointProperties(d)
	if pePropsErr != nil {
		return pePropsErr
	}

	pe := network.PrivateEndpoint{
		Name:                      &name,
		Location:                  &location,
		PrivateEndpointProperties: peProperties,
		Tags:                      expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, pe)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Private Endpoint %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Private Endpoint %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Private Endpoint %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPrivateEndpointRead(d, meta)
}

func resourceArmPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["privateEndpoints"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Private Endpoint %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	prop := resp.PrivateEndpointProperties
	if prop == nil {
		return fmt.Errorf("Error reading PrivateEndpointProperties")
	}

	if prop.Subnet != nil {
		if err := d.Set("subnet_id", prop.Subnet.ID); err != nil {
			return fmt.Errorf("Error setting `subnet_id`: %+v", err)
		}
	}

	if prop.ManualPrivateLinkServiceConnections != nil {
		if err := d.Set("manual_private_link_service_connections", flattenPrivateLinkServiceConnection(prop.ManualPrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `manual_private_link_service_connections`: %+v", err)
		}
	} else {
		if err := d.Set("private_link_service_connections", flattenPrivateLinkServiceConnection(prop.PrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `private_link_service_connections`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmPrivateEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["privateEndpoints"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Private Endpoints %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Private Endpoints %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandPrivateEndpointProperties(d *schema.ResourceData) (*network.PrivateEndpointProperties, error) {
	properties := &network.PrivateEndpointProperties{}

	if v, ok := d.GetOk("subnet_id"); ok {
		subnetID := v.(string)
		subnet := network.Subnet{ID: &subnetID}
		properties.Subnet = &subnet
	}

	if v, ok := d.GetOk("private_link_service_connections"); ok {
		plsc := v.([]interface{})
		properties.PrivateLinkServiceConnections = expandPrivateLinkServiceConnection(plsc)
	}
	if v, ok := d.GetOk("manual_private_link_service_connections"); ok {
		plsc := v.([]interface{})
		properties.ManualPrivateLinkServiceConnections = expandPrivateLinkServiceConnection(plsc)
	}

	return properties, nil
}

func expandPrivateLinkServiceConnection(connection []interface{}) *[]network.PrivateLinkServiceConnection {
	plsConnections := make([]network.PrivateLinkServiceConnection, 0, len(connection))

	for _, v := range connection {
		config := v.(map[string]interface{})
		name := config["name"].(string)
		privateLinkServiceID := config["private_link_service_id"].(string)
		requestMsg := config["request_message"].(string)
		gIDs := config["group_ids"].([]interface{})
		groupIDs := make([]string, 0, len(gIDs))
		for _, id := range gIDs {
			groupIDs = append(groupIDs, id.(string))
		}

		plsConnection := network.PrivateLinkServiceConnection{
			Name: &name,
			PrivateLinkServiceConnectionProperties: &network.PrivateLinkServiceConnectionProperties{
				PrivateLinkServiceID: &privateLinkServiceID,
				GroupIds:             &groupIDs,
				RequestMessage:       &requestMsg,
			},
		}

		plsConnections = append(plsConnections, plsConnection)
	}
	return &plsConnections
}

func flattenPrivateLinkServiceConnection(plsc *[]network.PrivateLinkServiceConnection) []interface{} {
	flat := make([]interface{}, 0, len(*plsc))

	if plsc == nil {
		return flat
	}

	for _, c := range *plsc {
		v := make(map[string]interface{})

		prop := c.PrivateLinkServiceConnectionProperties

		v["name"] = *prop.Name
		v["private_link_service_id"] = *prop.PrivateLinkServiceID
		v["group_ids"] = utils.FlattenStringSlice(prop.GroupIds)
		v["request_message"] = *prop.RequestMessage

		flat = append(flat, v)
	}
	return flat
}
