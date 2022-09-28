package resource

import (
	"fmt"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

//// TODO: unit tests for this
//

func expandAssignmentCodeForFieldObjectDefinition(left string, right string, fieldDefinition resourcemanager.TerraformSchemaFieldDefinition) (*string, error) {
	directAssignments := map[resourcemanager.TerraformSchemaFieldType]struct{}{
		resourcemanager.TerraformSchemaFieldTypeBoolean:  {},
		resourcemanager.TerraformSchemaFieldTypeDateTime: {}, // TODO: confirm
		resourcemanager.TerraformSchemaFieldTypeInteger:  {},
		resourcemanager.TerraformSchemaFieldTypeFloat:    {},
		resourcemanager.TerraformSchemaFieldTypeString:   {},
		// We're not dealing with these yet :see_no_evil:
		resourcemanager.TerraformSchemaFieldTypeList:                          {},
		resourcemanager.TerraformSchemaFieldTypeReference:                     {},
		resourcemanager.TerraformSchemaFieldTypeIdentitySystemAssigned:        {},
		resourcemanager.TerraformSchemaFieldTypeIdentitySystemAndUserAssigned: {},
		resourcemanager.TerraformSchemaFieldTypeIdentitySystemOrUserAssigned:  {},
		resourcemanager.TerraformSchemaFieldTypeIdentityUserAssigned:          {},
	}
	if _, ok := directAssignments[fieldDefinition.ObjectDefinition.Type]; ok {
		// TODO: if the field is optional, conditionally output this as a pointer
		output := fmt.Sprintf("\t\t%s = %s", left, right)
		if fieldDefinition.Optional {
			output = fmt.Sprintf("\t\t%s = utils.ToPtr(%s)", left, right)
		}
		return &output, nil
	}

	switch fieldDefinition.ObjectDefinition.Type {
	case resourcemanager.TerraformSchemaFieldTypeLocation:
		{
			output := fmt.Sprintf("\t\t%s = location.Normalize(%s)", left, right)
			if fieldDefinition.Optional {
				output = fmt.Sprintf("\t\t%s = location.NormalizeNilable(%s)", left, right)
			}
			return &output, nil
		}
	//case resourcemanager.TerraformSchemaFieldTypeIdentitySystemAssigned:
	//	{
	// 		NOTE: we need the actual model itself to determine what the type we're mapping too is.
	// 		TODO: should this be a part of the mapping for simplicities sake?
	// 		specifically here we need this for when the identity models can differ
	//	output := fmt.Sprintf("TODOormalize(%s)", mapping)
	//	return &output, nil
	//}
	case resourcemanager.TerraformSchemaFieldTypeTags:
		{
			output := fmt.Sprintf("\t\t%s = tags.Expand(%s)", left, right)
			return &output, nil
		}
	}
	return nil, fmt.Errorf("internal-error: unimplemented field type %q for expand mapping %q", string(fieldDefinition.ObjectDefinition.Type), left)
}
