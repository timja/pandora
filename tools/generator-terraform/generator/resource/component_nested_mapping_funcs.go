package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/pandora/tools/generator-terraform/generator/resource/mappings"

	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
)

func codeForNestedMappingFunctions(input models.ResourceInput) (*string, error) {
	schemaModel, ok := input.SchemaModels[input.SchemaModelName]
	if !ok {
		return nil, fmt.Errorf("the top-level schema model %q was not found", input.SchemaModelName)
	}

	createOperation, ok := input.Operations[input.Details.CreateMethod.MethodName]
	if !ok {
		return nil, fmt.Errorf("the model for Create %q was not found", input.Details.CreateMethod.MethodName)
	}
	createModelName := *createOperation.RequestObject.ReferenceName
	createModel, ok := input.Models[createModelName]
	if !ok {
		return nil, fmt.Errorf("the model %q used for the Create Operation %q was not found", createModelName, input.Details.CreateMethod.MethodName)
	}

	readOperation, ok := input.Operations[input.Details.ReadMethod.MethodName]
	if !ok {
		return nil, fmt.Errorf("the model for Read %q was not found", input.Details.ReadMethod.MethodName)
	}
	readModelName := *readOperation.ResponseObject.ReferenceName
	readModel, ok := input.Models[readModelName]
	if !ok {
		return nil, fmt.Errorf("the model %q used for the Create Operation %q was not found", readModelName, input.Details.ReadMethod.MethodName)
	}

	functions := make([]string, 0)
	for schemaFieldName, schemaField := range schemaModel.Fields {
		if mapping := schemaField.Mappings.SdkPathForCreate; mapping != nil && strings.Contains(*mapping, ".") {
			// TODO: Update shouldn't duplicate the Expand functions, keep track of what's generated
			helper := mappings.NestedMappingFunctionHelper{
				SdkServiceName:  input.SdkServiceName,
				SdkResourceName: input.SdkResourceName,
				SdkApiVersion:   input.SdkApiVersion,

				Models:                input.Models,
				TerraformSchemaModels: input.SchemaModels,
				ResourceName:          input.ResourceTypeName,

				Mapping:                  *mapping,
				FunctionName:             schemaFieldName,
				ModelName:                createModelName,
				Model:                    createModel,
				TerraformSchemaName:      input.SchemaModelName,
				TerraformSchemaModel:     schemaModel,
				TerraformSchemaFieldName: schemaFieldName,
				TerraformSchemaField:     schemaField,
			}
			function, err := helper.CodeForNestedExpandMappingFunction()
			if err != nil {
				return nil, fmt.Errorf("generating expand nested mapping functions for schema field %q with mapping %q: %+v", schemaFieldName, *mapping, err)
			}
			functions = append(functions, *function)
		}

		if mapping := schemaField.Mappings.SdkPathForRead; mapping != nil && strings.Contains(*mapping, ".") {
			helper := mappings.NestedMappingFunctionHelper{
				SdkServiceName:  input.SdkServiceName,
				SdkResourceName: input.SdkResourceName,
				SdkApiVersion:   input.SdkApiVersion,

				Models:                input.Models,
				TerraformSchemaModels: input.SchemaModels,
				ResourceName:          input.ResourceTypeName,

				Mapping:                  *mapping,
				FunctionName:             schemaFieldName,
				ModelName:                readModelName,
				Model:                    readModel,
				TerraformSchemaName:      input.SchemaModelName,
				TerraformSchemaModel:     schemaModel,
				TerraformSchemaFieldName: schemaFieldName,
				TerraformSchemaField:     schemaField,
			}
			function, err := helper.CodeForNestedFlattenMappingFunction()
			if err != nil {
				return nil, fmt.Errorf("generating flatten nested mapping functions for schema field %q with mapping %q: %+v", schemaFieldName, *mapping, err)
			}
			functions = append(functions, *function)
		}
	}
	sort.Strings(functions)

	output := strings.Join(functions, "\n")
	return &output, nil
}
