package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func codeForExpandAndFlattenFunctions(input models.ResourceInput) (*string, error) {
	schemaModels := input.SchemaModels
	sdkModels := input.Models
	mappings := input.Details.Mappings.Fields
	resourceName := input.SdkResourceName
	methods := make([]string, 0)
	output := ""

	for modelName, schemaModel := range schemaModels {
		// code for expand from Schema (config) to SDK (payload)
		//empty := ""
		codeForExpand, err := codeForExpandSchemaToSdkType(resourceName, modelName, schemaModel, sdkModels, mappings)
		if err != nil {
			return nil, err
		}
		codeForFlatten, err := codeForFlattenSdkTypeToSchema(resourceName, modelName, schemaModel, sdkModels, mappings)
		if err != nil {
			return nil, err
		}
		//if codeForFlatten == nil {
		//	codeForFlatten = &empty
		//}
		methods = append(methods, strings.Join([]string{*codeForExpand, *codeForFlatten}, "\n\n"))
	}

	if len(methods) > 0 {
		output = strings.Join(methods, "\n")
	}

	return &output, nil
}

func codeForExpandSchemaToSdkType(resourceName string, modelName string, model resourcemanager.TerraformSchemaModelDefinition, sdkModels map[string]resourcemanager.ModelDetails, mappings []resourcemanager.FieldMappingDefinition) (*string, error) {
	// outputs func (r SomeResource) expand{input.ResourceName}ResourceSchemaTo{SdkModelName}(input {input.ResourceName}ResourceSchema) *{SdkModelForCreate} {}
	lines := make([]string, 0)
	propertyLines := make([]string, 0)
	sdkModelName := strings.ToLower(resourceName)
	modelMappingsRaw, err := findMappingsForSchemaModel(modelName, mappings)
	if err != nil {
		return nil, err
	}
	if modelMappingsRaw == nil || len(*modelMappingsRaw) == 0 {
		// nothing here... Error or silently continue?
		// Mappings that need renaming will error out here, e.g. Sku's etc.
		empty := ""
		return &empty, nil
		// return nil, fmt.Errorf("no mappings found for model %q (resource %q)", modelName, resourceName)
	}
	modelMappings := *modelMappingsRaw
	orderedFieldNames := make([]string, 0)
	for fieldName := range model.Fields {
		orderedFieldNames = append(orderedFieldNames, fieldName)
	}
	sort.Strings(orderedFieldNames)
	for _, fieldName := range orderedFieldNames {
		mapping := modelMappings[fieldName]
		switch mapping.Type {
		case resourcemanager.DirectAssignmentMappingDefinitionType:
			if field, ok := model.Fields[fieldName]; ok {
				if strings.HasSuffix(mapping.DirectAssignment.SdkModelName, "Properties") {
					line, err := expandAssignmentCodeForFieldObjectDefinition(field, mapping)
					if err != nil {
						return nil, err
					}
					propertyLines = append(propertyLines, *line)
				} else {
					line, err := expandAssignmentCodeForFieldObjectDefinition(field, mapping)
					if err != nil {
						return nil, err
					}
					lines = append(lines, *line)
				}
			}
		}
	}
	output := ""
	if len(propertyLines) > 0 {
		output = fmt.Sprintf(`
func (r %[1]s) expand%[2]sTo%[3]s(input %[2]s) *%[4]s.%[3]s {
	result := &%[4]s.%[3]s{}
%[5]s
		result.Properties = 
%[6]s
	return result
}
`, strings.TrimSuffix(modelName, "Schema"), modelName, sdkModelName, strings.ToLower(resourceName), strings.Join(lines, "\n"), strings.Join(propertyLines, "\n"))
	} else {
		output = fmt.Sprintf(`
func (r %[1]s) expand%[2]sTo%[3]s(input %[2]s) *%[4]s.%[3]s {
	result := &%[4]s.%[3]s{}
%[5]s
	
	return result
}
`, strings.TrimSuffix(modelName, "Schema"), modelName, sdkModelName, strings.ToLower(resourceName), strings.Join(lines, "\n"))
	}

	return &output, nil
}

func codeForFlattenSdkTypeToSchema(resourceName string, modelName string, model resourcemanager.TerraformSchemaModelDefinition, sdkModels map[string]resourcemanager.ModelDetails, mappings []resourcemanager.FieldMappingDefinition) (*string, error) {
	// outputs func (r SomeResource) flatten{SdkModelForRead}To{input.ResourceName}{SchemaModelName}Schema(input {input.ResourceName}ResourceSchema) *{SdkModelForCreate} {}
	lines := make([]string, 0)
	sdkModelName := "BLAHBLAHBLAH"
	modelMappingsRaw, err := findMappingsForSchemaModel(modelName, mappings)
	if err != nil {
		return nil, err
	}
	if modelMappingsRaw == nil || len(*modelMappingsRaw) == 0 {
		// nothing here... Error or continue?
		empty := ""
		return &empty, nil
		//return nil, fmt.Errorf("no mappings found for model %q (resource %q)", modelName, resourceName)
	}
	modelMappings := *modelMappingsRaw
	orderedFieldNames := make([]string, 0)
	for fieldName := range model.Fields {
		orderedFieldNames = append(orderedFieldNames, fieldName)
	}
	sort.Strings(orderedFieldNames)
	for _, fieldName := range orderedFieldNames {
		mapping := modelMappings[fieldName]
		switch mapping.Type {
		case resourcemanager.DirectAssignmentMappingDefinitionType:
			if field, ok := model.Fields[fieldName]; ok {
				line, err := flattenAssignmentCodeForFieldObjectDefinition(field, mapping, strings.ToLower(resourceName))
				if err != nil {
					return nil, err
				}
				lines = append(lines, *line)
			}
		}
	}

	output := fmt.Sprintf(`
func (r %[1]s) flatten%[3]sTo%[2]s(input *%[4]s.%[3]s, schema *%[2]s) error {
%[5]s
	
	return nil
}
`, strings.TrimSuffix(modelName, "Schema"), modelName, sdkModelName, strings.ToLower(resourceName), strings.Join(lines, "\n"))
	return &output, nil
}

func findMappingsForSchemaModel(input string, mappings []resourcemanager.FieldMappingDefinition) (*map[string]resourcemanager.FieldMappingDefinition, error) {
	output := make(map[string]resourcemanager.FieldMappingDefinition, 0)
	for _, mapping := range mappings {
		switch mapping.Type {
		case resourcemanager.DirectAssignmentMappingDefinitionType:
			if mapping.DirectAssignment == nil {
				return nil, fmt.Errorf("mapping had type `DirectAssignment`, but that block was nil")
			}
			if mapping.DirectAssignment.SchemaModelName == input {
				output[mapping.DirectAssignment.SchemaFieldPath] = mapping
			}

		default:
			return nil, fmt.Errorf("currently unsupported mapping type %q", mapping.Type)

		}
	}
	return &output, nil
}

func sdkModelFromSchemaModelName(input string, resourcePrefix string, sdkModels map[string]resourcemanager.ModelDetails) (string, resourcemanager.ModelDetails) {
	if _, ok := sdkModels[input]; ok {
		return input, sdkModels[input]
	}
	modelName := strings.TrimSuffix(input, "Schema")
	modelName = strings.TrimPrefix(modelName, resourcePrefix)
	modelName = strings.TrimPrefix(modelName, "Resource")
	return modelName, sdkModels[modelName]
}
