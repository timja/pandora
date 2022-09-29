package resource

import (
	"fmt"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

//// TODO: unit tests for this
//

func expandAssignmentCodeForFieldObjectDefinition(fieldDefinition resourcemanager.TerraformSchemaFieldDefinition, mapping resourcemanager.FieldMappingDefinition) (*string, error) {
	left := fmt.Sprintf("result.%s", mapping.DirectAssignment.SdkFieldPath)
	right := fmt.Sprintf("input.%s", mapping.DirectAssignment.SchemaFieldPath)
	directAssignments := map[resourcemanager.TerraformSchemaFieldType]struct{}{
		resourcemanager.TerraformSchemaFieldTypeBoolean:  {},
		resourcemanager.TerraformSchemaFieldTypeDateTime: {}, // TODO: confirm
		resourcemanager.TerraformSchemaFieldTypeInteger:  {},
		resourcemanager.TerraformSchemaFieldTypeFloat:    {},
		resourcemanager.TerraformSchemaFieldTypeString:   {},
		// We're not dealing with these yet :see_no_evil:
		resourcemanager.TerraformSchemaFieldTypeIdentitySystemAssigned:        {},
		resourcemanager.TerraformSchemaFieldTypeIdentitySystemAndUserAssigned: {},
		resourcemanager.TerraformSchemaFieldTypeIdentitySystemOrUserAssigned:  {},
		resourcemanager.TerraformSchemaFieldTypeIdentityUserAssigned:          {},
	}
	if _, ok := directAssignments[fieldDefinition.ObjectDefinition.Type]; ok {
		// TODO: if the field is optional, conditionally output this as a pointer
		output := fmt.Sprintf("\t\t%s = %s", left, right)
		if fieldDefinition.Optional {
			output = fmt.Sprintf("\t\t%s = utils.ToPtr(%s) %s", left, right, "// TODO - unhandled SchemaFieldType Case")
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
	case resourcemanager.TerraformSchemaFieldTypeReference:
		{
			right = fmt.Sprintf("expand%[1]sTo%[2]s(%[3]s)", mapping.DirectAssignment.SchemaModelName, mapping.DirectAssignment.SdkModelName, right)
			output := fmt.Sprintf(`		if %s, err := %s; err != nil {
			return err
		}`, left, right)
			return &output, nil
		}
	case resourcemanager.TerraformSchemaFieldTypeList:
		{
			if fieldDefinition.ObjectDefinition.NestedObject == nil {
				return nil, fmt.Errorf("generating expand for %q (model %q), type was List but NestedObject was nil", mapping.DirectAssignment.SchemaFieldPath, mapping.DirectAssignment.SchemaModelName)
			}

			switch fieldDefinition.ObjectDefinition.NestedObject.Type {
			case resourcemanager.TerraformSchemaFieldTypeString, resourcemanager.TerraformSchemaFieldTypeFloat, resourcemanager.TerraformSchemaFieldTypeInteger:
				{
					output := fmt.Sprintf("\t\t%s = %s", left, right)
					return &output, nil
				}
			case resourcemanager.TerraformSchemaFieldTypeReference:
				{
					if fieldDefinition.ObjectDefinition.NestedObject.ReferenceName == nil {
						return nil, fmt.Errorf("generating expand for %q (model %q), type was List with type `Reference` but NestedObject ReferenceName was nil", mapping.DirectAssignment.SchemaFieldPath, mapping.DirectAssignment.SchemaModelName)
					}
					output := fmt.Sprintf("\t\t%[1]s = expand%[2]sTo%[3]s(%[4]s)", left, *fieldDefinition.ObjectDefinition.NestedObject.ReferenceName, mapping.DirectAssignment.SdkModelName, right)
					return &output, nil
				}

			default:
				return nil, fmt.Errorf("generating list expansion for %q (model %q), unsupported NestedObject.Type, got %q", mapping.DirectAssignment.SchemaFieldPath, mapping.DirectAssignment.SchemaModelName, fieldDefinition.ObjectDefinition.NestedObject.Type)
			}
		}

	}
	return nil, fmt.Errorf("internal-error: unimplemented field type %q for expand mapping %q", string(fieldDefinition.ObjectDefinition.Type), left)
}
