package resource

import (
	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func codeForExpandAndFlattenFunctions(input models.ResourceInput) (*string, error) {
	mappings := input.Details.Mappings

	// code for expand from Schema (config) to SDK (payload)
	codeForExpand, err := expandSchemaToSdkType(mappings)
	if err != nil {
		return nil, err
	}
	// code for flatten from SDK (resp) to Schema (model)
	codeForFlatten, err := flattenSdkTypeToSchema(mappings)
	return nil, nil
}

func expandSchemaToSdkType(mappings resourcemanager.MappingDefinition) (*string, error) {

	// outputs func (r SomeResource) expand{input.ResourceName}ResourceSchemaTo{SdkModelName}(input {input.ResourceName}ResourceSchema) *{SdkModelForCreate} {}

	return nil, nil
}

func flattenSdkTypeToSchema(mappings resourcemanager.MappingDefinition) (*string, error) {

	// outputs func (r SomeResource) flatten{SdkModelForRead}To{input.ResourceName}{SchemaModelName}Schema(input {input.ResourceName}ResourceSchema) *{SdkModelForCreate} {}

	return nil, nil
}
