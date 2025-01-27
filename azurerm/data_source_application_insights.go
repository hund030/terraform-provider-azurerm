package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmApplicationInsights() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmApplicationInsightsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"instrumentation_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"application_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmApplicationInsightsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppInsights.ComponentsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Application Insights bucket %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on Application Insights bucket %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("instrumentation_key", resp.InstrumentationKey)
	d.Set("location", resp.Location)
	d.Set("app_id", resp.AppID)
	d.Set("application_type", resp.ApplicationType)
	return tags.FlattenAndSet(d, resp.Tags)
}
