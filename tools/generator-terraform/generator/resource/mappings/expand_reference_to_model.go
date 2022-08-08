package mappings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

var _ expandDefinition = expandReferenceToModel{}

type expandReferenceToModel struct {
}

func (expandReferenceToModel) isApplicable(input resourcemanager.TerraformSchemaFieldDefinition, output resourcemanager.FieldDetails) bool {
	if input.ObjectDefinition.Type != resourcemanager.TerraformSchemaFieldTypeReference {
		return false
	}
	if output.ObjectDefinition.Type != resourcemanager.ReferenceApiObjectDefinitionType {
		return false
	}

	return true
}

func (expandReferenceToModel) mappingCode(input resourcemanager.TerraformSchemaFieldDefinition, output resourcemanager.FieldDetails, metadata NestedMappingFunctionHelper) (*string, error) {
	sdkResourceName := strings.ToLower(metadata.SdkResourceName)
	inputType, err := input.ObjectDefinition.GolangFieldType()
	if err != nil {
		return nil, fmt.Errorf("determining Go Field Type for input ObjectDefinition: %+v", err)
	}

	outputType, err := output.ObjectDefinition.GolangTypeName(&sdkResourceName)
	if err != nil {
		return nil, fmt.Errorf("determining Go Type Name for output ObjectDefinition: %+v", err)
	}

	if input.Optional {
		out := fmt.Sprintf(`
func (r %[1]sResource) expand%[2]s(input *%[3]s) %[4]s {
	output := %[4]s{}
	if input == nil {
		return output
	}

	// map the fields within 'input'

	return output
}
`, metadata.ResourceName, metadata.FunctionName, *inputType, *outputType)

		if output.Optional {
			out = fmt.Sprintf(`
func (r %[1]sResource) expand%[2]s(input *%[3]s) *[]%[4]s {
	output := make([]%[4]s, 0)
	if input == nil {
		return nil
	}

	for _, v := range *input {
		item := %[4]s{}
		// map the fields within 'v'
		output = append(output, item)
	}

	return &output
}
`, metadata.ResourceName, metadata.FunctionName, *inputType, *outputType)
		}
		return &out, nil
	}

	out := fmt.Sprintf(`
func (r %[1]sResource) expand%[2]s(input %[3]s) %[4]s {
	output := %[4]s{}

	// map the fields within 'input'

	return output
}
`, metadata.ResourceName, metadata.FunctionName, *inputType, *outputType)

	if output.Optional {
		out = fmt.Sprintf(`
func (r %[1]sResource) expand%[2]s(input %[3]s) *%[4]s {
	output := %[4]s{}

	// map the fields within 'input'

	return &output
}
`, metadata.ResourceName, metadata.FunctionName, *inputType, *outputType)
	}

	return &out, nil
}
