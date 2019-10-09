package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceArmAutomationVariableBool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAutomationVariableBoolRead,

		Schema: datasourceAutomationVariableCommonSchema(schema.TypeBool),
	}
}

func dataSourceArmAutomationVariableBoolRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceAutomationVariableRead(d, meta, "Bool")
}
