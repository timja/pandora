package resource

import (
	"fmt"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func flattenAssignmentCodeForFieldObjectDefinition(left string, right string, fieldDefinition resourcemanager.TerraformSchemaFieldDefinition) (*string, error) {
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
			switch fieldDefinition.ObjectDefinition.Type {
			// TODO - Need generics in the provider for these
			case resourcemanager.TerraformSchemaFieldTypeString:
				output = fmt.Sprintf("\t\t%s = utils.NormalizeNilableString(%s)", left, right)
			case resourcemanager.TerraformSchemaFieldTypeFloat:
				output = fmt.Sprintf("\t\t%s = utils.NormalizeNilableFloat64(%s)", left, right)
			case resourcemanager.TerraformSchemaFieldTypeBoolean:
				output = fmt.Sprintf("\t\t%s = utils.NormaliseNilableBool(%s) ", left, right)
			case resourcemanager.TerraformSchemaFieldTypeInteger:
				output = fmt.Sprintf("\t\t%s = utils.NormaliseNilableInt64(%s) ", left, right)
			default:
				output = fmt.Sprintf("\t\t%s = %s", left, right) // TODO - shouldn't get here when TerraformSchemaFieldTypeReference is supported?
			}
		}
		return &output, nil
	}

	switch fieldDefinition.ObjectDefinition.Type {
	case resourcemanager.TerraformSchemaFieldTypeLocation:
		{
			output := fmt.Sprintf("%s = location.Normalize(%s)", left, right)
			if fieldDefinition.Optional {
				output = fmt.Sprintf("%s = location.NormalizeNilable(%s)", left, right)
			}

			return &output, nil
		}
	case resourcemanager.TerraformSchemaFieldTypeTags:
		{
			output := fmt.Sprintf("%s = tags.Flatten(%s)", left, right)
			return &output, nil
		}
	}

	return nil, fmt.Errorf("internal-error: unimplemented field type %q for flatten mapping %q", string(fieldDefinition.ObjectDefinition.Type), left)
}
