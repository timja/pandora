package dataapigenerator

import (
	"fmt"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

// testTemplateForCustomType will be an ever expanding list of common schemas that can be pulled from to build out the config for generated tests
func testTemplateForCustomType(input resourcemanager.TerraformSchemaFieldType, randomInt int, packageName string) (*string, error) {
	var nilableType = func(in string) (*string, error) {
		return &in, nil
	}

	switch input {
	case resourcemanager.TerraformSchemaFieldTypeResourceGroup:
		// todo pass in location?
		return nilableType(fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%s-%d"
  location = "westus"
}
`, packageName, randomInt))

	case resourcemanager.TerraformSchemaFieldTypeEdgeZone:
		return nilableType(fmt.Sprintf(`
data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}
`))

	case resourcemanager.TerraformSchemaFieldTypeIdentitySystemAndUserAssigned, resourcemanager.TerraformSchemaFieldTypeIdentityUserAssigned:
		return nilableType(fmt.Sprintf(`
resource "azurerm_user_assigned_identity" "test" {
	name                = "acctest-%d"
	resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
}
`, randomInt))
	}

	return nil, fmt.Errorf("unmapped Custom Type %q", string(input))
}
