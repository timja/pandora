package mappings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

type NestedMappingFunctionHelper struct {
	SdkServiceName  string
	SdkResourceName string
	SdkApiVersion   string

	Models                map[string]resourcemanager.ModelDetails
	TerraformSchemaModels map[string]resourcemanager.TerraformSchemaModelDefinition
	ResourceName          string

	Mapping                  string
	FunctionName             string
	ModelName                string
	Model                    resourcemanager.ModelDetails
	TerraformSchemaName      string
	TerraformSchemaModel     resourcemanager.TerraformSchemaModelDefinition
	TerraformSchemaFieldName string
	TerraformSchemaField     resourcemanager.TerraformSchemaFieldDefinition
}

func (h NestedMappingFunctionHelper) CodeForNestedExpandMappingFunction() (*string, error) {
	split := strings.Split(h.Mapping, ".")
	modelFieldName := split[0]
	fieldFromModel, ok := h.Model.Fields[modelFieldName]
	if !ok {
		return nil, fmt.Errorf("couldn't find the field %q in model %q used in mapping %q", modelFieldName, h.ModelName, h.Mapping)
	}

	lines := ""
	for _, v := range expanders {
		if !v.isApplicable(h.TerraformSchemaField, fieldFromModel) {
			continue
		}

		mappingCode, err := v.mappingCode(h.TerraformSchemaField, fieldFromModel, h)
		if err != nil {
			return nil, fmt.Errorf("building mapping code for the field %q in model %q used in mapping %q", modelFieldName, h.ModelName, h.Mapping)
		}
		lines = *mappingCode
	}
	if lines == "" {
		return nil, fmt.Errorf("no expand mappings found for Field %q to SDK Model %q", h.TerraformSchemaField.ObjectDefinition.Type, fieldFromModel.ObjectDefinition.String())
	}

	return &lines, nil
}

func (h NestedMappingFunctionHelper) CodeForNestedFlattenMappingFunction() (*string, error) {
	split := strings.Split(h.Mapping, ".")

	// find the model associated with this mapping
	modelFieldName := split[0]
	fieldFromModel, ok := h.Model.Fields[modelFieldName]
	if !ok {
		return nil, fmt.Errorf("couldn't find the field %q in model %q used in mapping %q", modelFieldName, h.ModelName, h.Mapping)
	}

	sdkResourceName := strings.ToLower(h.SdkResourceName)
	returnType, err := fieldFromModel.ObjectDefinition.GolangTypeName(&sdkResourceName)
	if err != nil {
		return nil, fmt.Errorf("determining Go Type Name for ObjectDefinition for Field %q in Model %q: %+v", modelFieldName, h.ModelName, err)
	}

	inputType, err := h.TerraformSchemaField.ObjectDefinition.GolangFieldType()
	if err != nil {
		return nil, fmt.Errorf("determining Go Field Type for FieldObjectDefinition for Field %q in Model %q: %+v", modelFieldName, h.ModelName, err)
	}

	if h.TerraformSchemaField.Optional {
		v := fmt.Sprintf("*%s", *inputType)
		inputType = &v
	}
	if fieldFromModel.Optional {
		v := fmt.Sprintf("*%s", *returnType)
		returnType = &v
	}

	// TODO: the type of function is going to need to change dynamically depending on _what_ is being output
	// e.g. list/basic type/optional/required etc

	output := fmt.Sprintf(`
func (r %[1]sResource) flatten%[2]s(input %[3]s) %[4]s {
	var output %[4]s

	return output
}
`, h.ResourceName, h.FunctionName, *returnType, *inputType)
	return &output, nil
}
