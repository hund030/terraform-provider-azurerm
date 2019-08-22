package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmPrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateEndpointRead,

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

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_link_service_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_link_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"request_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"manual_private_link_service_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_link_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"request_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Network Interface %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	prop := resp.PrivateEndpointProperties
	if prop == nil {
		return fmt.Errorf("Error reading PrivateEndpointProperties")
	}

	if prop.Subnet != nil {
		if err := d.Set("subnet_id", prop.Subnet.ID); err != nil {
			return fmt.Errorf("Error setting `subnet_id`: %+v", err)
		}
	}

	plsc := prop.ManualPrivateLinkServiceConnections
	if prop.ManualPrivateLinkServiceConnections == nil {
		plsc = prop.PrivateLinkServiceConnections
	}

	flat := make([]interface{}, 0, len(*plsc))
	for _, c := range *plsc {
		v := make(map[string]interface{})

		prop := c.PrivateLinkServiceConnectionProperties

		v["name"] = c.Name
		v["private_link_service_id"] = *prop.PrivateLinkServiceID
		v["group_ids"] = utils.FlattenStringSlice(prop.GroupIds)
		v["request_message"] = *prop.RequestMessage

		flat = append(flat, v)
	}

	if prop.ManualPrivateLinkServiceConnections != nil {
		if err := d.Set("manual_private_link_service_connections", flat); err != nil {
			return fmt.Errorf("Error setting `manual_private_link_service_connections`: %+v", err)
		}
	} else {
		if err := d.Set("private_link_service_connections", flat); err != nil {
			return fmt.Errorf("Error setting `private_link_service_connections`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return nil
}
