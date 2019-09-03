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
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"private_link_service_connections": {
				Type:          schema.TypeList,
				Optional:      true,
				MinItems:      1,
				ConflictsWith: []string{"manual_private_link_service_connections"},
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
								Type: schema.TypeString,
							},
						},

						"request_message": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Please approve my connection.",
							ValidateFunc: validate.NoEmptyStrings,
						},

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"manual_private_link_service_connections": {
				Type:          schema.TypeList,
				Optional:      true,
				MinItems:      1,
				ConflictsWith: []string{"private_link_service_connections"},
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

						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network_interfaces_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
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

	prop, err := expandPrivateEndpointProperties(d)
	if err != nil {
		return err
	}

	param := network.PrivateEndpoint{
		Name:                      &name,
		Location:                  &location,
		PrivateEndpointProperties: prop,
		Tags:                      expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, param)
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

	if err := d.Set("subnet_id", prop.Subnet.ID); err != nil {
		return fmt.Errorf("Error setting `subnet_id`: %+v", err)
	}

	if err := d.Set("manual_private_link_service_connections", flattenPrivateLinkServiceConnection(prop.ManualPrivateLinkServiceConnections)); err != nil {
		return fmt.Errorf("Error setting `manual_private_link_service_connections`: %+v", err)
	}

	if err := d.Set("private_link_service_connections", flattenPrivateLinkServiceConnection(prop.PrivateLinkServiceConnections)); err != nil {
		return fmt.Errorf("Error setting `private_link_service_connections`: %+v", err)
	}

	if err := d.Set("network_interfaces_id", flattenNetworkInterfacesID(prop.NetworkInterfaces)); err != nil {
		return fmt.Errorf("Error setting `network_interfaces_id`: %+v", err)
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
	subnetId := d.Get("subnet_id").(string)
	prop := &network.PrivateEndpointProperties{
		Subnet: &network.Subnet{
			ID: &subnetId,
		},
	}

	plsc := d.Get("private_link_service_connections").([]interface{})
	prop.PrivateLinkServiceConnections = expandPrivateLinkServiceConnection(plsc)
	mplsc := d.Get("manual_private_link_service_connections").([]interface{})
	prop.ManualPrivateLinkServiceConnections = expandPrivateLinkServiceConnection(mplsc)

	return prop, nil
}

func expandPrivateLinkServiceConnection(inputs []interface{}) *[]network.PrivateLinkServiceConnection {
	results := make([]network.PrivateLinkServiceConnection, 0, len(inputs))

	for _, item := range inputs {
		config := item.(map[string]interface{})
		name := config["name"].(string)
		privateLinkServiceID := config["private_link_service_id"].(string)
		requestMsg := config["request_message"].(string)
		gIDs := config["group_ids"].([]interface{})
		groupIDs := make([]string, 0, len(gIDs))
		for _, id := range gIDs {
			groupIDs = append(groupIDs, id.(string))
		}

		result := network.PrivateLinkServiceConnection{
			Name: &name,
			PrivateLinkServiceConnectionProperties: &network.PrivateLinkServiceConnectionProperties{
				PrivateLinkServiceID: &privateLinkServiceID,
				GroupIds:             &groupIDs,
				RequestMessage:       &requestMsg,
			},
		}

		results = append(results, result)
	}
	return &results
}

func flattenPrivateLinkServiceConnection(inputs *[]network.PrivateLinkServiceConnection) []interface{} {
	results := make([]interface{}, 0, len(*inputs))

	if inputs == nil {
		return results
	}

	for _, item := range *inputs {
		result := make(map[string]interface{})

		prop := item.PrivateLinkServiceConnectionProperties

		result["name"] = item.Name
		result["private_link_service_id"] = *prop.PrivateLinkServiceID
		result["group_ids"] = utils.FlattenStringSlice(prop.GroupIds)
		result["request_message"] = *prop.RequestMessage
		result["status"] = *prop.PrivateLinkServiceConnectionState.Status

		results = append(results, result)
	}
	return results
}

func flattenNetworkInterfacesID(inputs *[]network.Interface) []interface{} {
	results := make([]interface{}, 0, len(*inputs))

	if inputs == nil {
		return results
	}

	for _, item := range *inputs {
		results = append(results, item.ID)
	}

	return results
}
