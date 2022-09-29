package resource

import (
	"fmt"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func flattenAssignmentCodeForFieldObjectDefinition(fieldDefinition resourcemanager.TerraformSchemaFieldDefinition, mapping resourcemanager.FieldMappingDefinition, resourceName string) (*string, error) {
	left := fmt.Sprintf("schema.%s", mapping.DirectAssignment.SchemaFieldPath)
	right := fmt.Sprintf("input.%s", mapping.DirectAssignment.SdkFieldPath)
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
				output = fmt.Sprintf("\t\t%s = %s %s", left, right, "// TODO - unhandled SchemaFieldType Case") // TODO - unhandled SchemaFieldType Case
			}
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
	case resourcemanager.TerraformSchemaFieldTypeTags:
		{
			output := fmt.Sprintf("\t\t%s = tags.Flatten(%s)", left, right)
			return &output, nil
		}
	case resourcemanager.TerraformSchemaFieldTypeReference:
		right = fmt.Sprintf("flatten%[1]sTo%[2]s(%[3]s, &%[4]s)", mapping.DirectAssignment.SdkModelName, *fieldDefinition.ObjectDefinition.ReferenceName, left, right)
		{
			output := fmt.Sprintf(`		schema.%[1]s = %[2]s{}
		if err := %[3]s; err != nil {
			return nil
		}`, mapping.DirectAssignment.SdkFieldPath, resourceName, right)
			return &output, nil
		}
	case resourcemanager.TerraformSchemaFieldTypeList:
		if fieldDefinition.ObjectDefinition.NestedObject == nil {
			return nil, fmt.Errorf("generating flatten for %q (model %q), type was List but NestedObject was nil", mapping.DirectAssignment.SchemaFieldPath, mapping.DirectAssignment.SchemaModelName)
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
					return nil, fmt.Errorf("generating flatten for %q (model %q), type was List with type `Reference` but NestedObject ReferenceName was nil", mapping.DirectAssignment.SchemaFieldPath, mapping.DirectAssignment.SchemaModelName)
				}
				output := fmt.Sprintf("\t\t%[1]s = flatten%[2]sTo%[3]s(%[4]s)", left, mapping.DirectAssignment.SdkModelName, *fieldDefinition.ObjectDefinition.NestedObject.ReferenceName, right)
				return &output, nil

			}
		}
	}

	return nil, fmt.Errorf("internal-error: unimplemented field type %q for flatten mapping %q", string(fieldDefinition.ObjectDefinition.Type), left)
}
