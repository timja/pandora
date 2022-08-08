package mappings

import "github.com/hashicorp/pandora/tools/sdk/resourcemanager"

type expandDefinition interface {
	isApplicable(input resourcemanager.TerraformSchemaFieldDefinition, output resourcemanager.FieldDetails) bool

	mappingCode(input resourcemanager.TerraformSchemaFieldDefinition, output resourcemanager.FieldDetails, mapping NestedMappingFunctionHelper) (*string, error)
}

var expanders = []expandDefinition{
	expandReferenceToModel{},
	expandReferenceToListOfModel{},
}
