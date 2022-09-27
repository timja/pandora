package resource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

//// TODO: unit tests for this
//

const (
	topLevelFromPath           = "payload."
	topLevelPropertiesFromPath = "payload.Properties."
	topLevelToPath             = "model."
)

func expandAssignmentCodeForCreateField(fieldMapping resourcemanager.FieldMappingDefinition, field resourcemanager.TerraformSchemaFieldDefinition, modelName string) (*string, error) {
	if strings.Contains(fieldMapping.From.SchemaFieldPath, ".") {
		// TODO - Pure guesswork right now - revisit when nested mappings are a thing
		toPath := fieldMapping.To.SdkModelName
		if fieldMapping.To.SdkFieldPath != "" {
			toPath = strings.Join([]string{fieldMapping.To.SdkFieldPath, fieldMapping.To.SdkModelName}, ".")
		}
		assignmentCode := fmt.Sprintf("r.expand%[1]s(%s%[2]s)", fieldMapping.To.SdkModelName, topLevelToPath, toPath)
		fromPath := fieldMapping.To.SdkModelName
		if fieldMapping.From.SchemaFieldPath != "" {
			fromPath = strings.Join([]string{fieldMapping.From.SchemaFieldPath, fieldMapping.From.SchemaModelName}, ".")
		}
		output := fmt.Sprintf("%s.%s = %s", topLevelToPath, fromPath, assignmentCode)
		return &output, nil
	}

	assignmentCode, err := expandAssignmentCodeForFieldObjectDefinition(fmt.Sprintf("model.%[1]s", fieldMapping.To.SdkFieldPath), field)
	if err != nil {
		return nil, fmt.Errorf("building expand assignment code for top level field %q: %+v", fieldMapping.To.SdkFieldPath, err)
	}
	output := ""
	switch field.ObjectDefinition.Type {
	case resourcemanager.TerraformSchemaFieldTypeLocation, resourcemanager.TerraformSchemaFieldTypeZones, resourcemanager.TerraformSchemaFieldTypeZone, resourcemanager.TerraformSchemaFieldTypeTags:
		output = fmt.Sprintf("%s%s = %s", topLevelFromPath, fieldMapping.From.SchemaFieldPath, *assignmentCode)
	default:
		output = fmt.Sprintf("%s%s = %s", topLevelPropertiesFromPath, fieldMapping.From.SchemaFieldPath, *assignmentCode)

	}
	return &output, nil
}

//func expandAssignmentCodeForCreateField(assignmentVariable string, schemaFieldName string, field resourcemanager.TerraformSchemaFieldDefinition, currentModel resourcemanager.ModelDetails, models map[string]resourcemanager.ModelDetails) (*string, error) {
//	// if it's a nested mapping (e.g. `Properties.Foo`) we need to pass `Properties` to
//	// the expand function, which in turn needs to check if `Foo` is nil (and return
//	// whatever it needs too)
//	topLevelFieldMapping := *field.Mappings.SdkPathForCreate
//	if strings.Contains(topLevelFieldMapping, ".") {
//		split := strings.Split(topLevelFieldMapping, ".")
//		topLevelFieldMapping = split[0]
//
//		// TODO: generate that method which needs to split/nil-check on
//		// remainingMapping := strings.Join(split[1:], ".")
//
//		assignmentCode := fmt.Sprintf("r.expand%[1]s(config.%[2]s)", schemaFieldName, topLevelFieldMapping)
//		output := fmt.Sprintf("// TODO: - %s = %s", assignmentVariable, assignmentCode)
//		return &output, nil
//	}
//
//	assignmentCode, err := expandAssignmentCodeForFieldObjectDefinition(fmt.Sprintf("config.%[1]s", schemaFieldName), field.ObjectDefinition)
//	if err != nil {
//		return nil, fmt.Errorf("building expand assignment code for top level field %q: %+v", schemaFieldName, err)
//	}
//
//	output := fmt.Sprintf("%s = %s", assignmentVariable, *assignmentCode)
//	return &output, nil
//}

//func expandAssignmentCodeForUpdateField(assignmentVariable string, schemaFieldName string, field resourcemanager.TerraformSchemaFieldDefinition, currentModel resourcemanager.ModelDetails, models map[string]resourcemanager.ModelDetails) (*string, error) {
//	// if it's a nested mapping (e.g. `Properties.Foo`) we need to pass `Properties` to
//	// the expand function, which in turn needs to check if `Foo` is nil (and return
//	// whatever it needs too)
//	topLevelFieldMapping := *field.Mappings.SdkPathForUpdate
//	if strings.Contains(topLevelFieldMapping, ".") {
//		split := strings.Split(topLevelFieldMapping, ".")
//		topLevelFieldMapping = split[0]
//
//		// TODO: generate that method which needs to split/nil-check on
//		// remainingMapping := strings.Join(split[1:], ".")
//
//		assignmentCode := fmt.Sprintf("r.expand%[1]s(config.%[2]s)", schemaFieldName, topLevelFieldMapping)
//		output := fmt.Sprintf("// TODO: - %s = %s", assignmentVariable, assignmentCode)
//		return &output, nil
//	}
//
//	assignmentCode, err := expandAssignmentCodeForFieldObjectDefinition(fmt.Sprintf("config.%[1]s", schemaFieldName), field.ObjectDefinition)
//	if err != nil {
//		return nil, fmt.Errorf("building expand assignment code for top level field %q: %+v", schemaFieldName, err)
//	}
//
//	output := fmt.Sprintf("%s = %s", assignmentVariable, *assignmentCode)
//	return &output, nil
//}
//
func expandAssignmentCodeForFieldObjectDefinition(mapping string, fieldDefinition resourcemanager.TerraformSchemaFieldDefinition) (*string, error) {
	directAssignments := map[resourcemanager.TerraformSchemaFieldType]struct{}{
		resourcemanager.TerraformSchemaFieldTypeBoolean:  {},
		resourcemanager.TerraformSchemaFieldTypeDateTime: {}, // TODO: confirm
		resourcemanager.TerraformSchemaFieldTypeInteger:  {},
		resourcemanager.TerraformSchemaFieldTypeFloat:    {},
		resourcemanager.TerraformSchemaFieldTypeString:   {},
	}
	if _, ok := directAssignments[fieldDefinition.ObjectDefinition.Type]; ok {
		// TODO: if the field is optional, conditionally output this as a pointer
		if fieldDefinition.Optional {
			mapping = fmt.Sprintf("utils.String(%s)", mapping)
		}
		return &mapping, nil
	}

	switch fieldDefinition.ObjectDefinition.Type {
	case resourcemanager.TerraformSchemaFieldTypeLocation:
		{
			output := fmt.Sprintf("location.Normalize(%s)", mapping)
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
			output := fmt.Sprintf("tags.Expand(%s)", mapping)
			return &output, nil
		}
	}
	//return nil, fmt.Errorf("internal-error: unimplemented field type %q for expand mapping %q", string(fieldDefinition.ObjectDefinition.Type), mapping)
	// TODO - Hack for Gen until there's real mappings to work with on actual resources
	empty := ""
	return &empty, nil
}
