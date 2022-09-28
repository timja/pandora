package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

type readFunctionComponents struct {
	constants       map[string]resourcemanager.ConstantDetails
	idParseLine     string
	mappings        resourcemanager.MappingDefinition
	readMethod      resourcemanager.MethodDefinition
	readOperation   resourcemanager.ApiOperation
	resourceId      resourcemanager.ResourceIdDefinition
	schemaModelName string
	sdkResourceName string
	terraformModel  resourcemanager.TerraformSchemaModelDefinition
	topLevelModel   resourcemanager.ModelDetails
}

func readFunctionForResource(input models.ResourceInput) (*string, error) {
	if !input.Details.ReadMethod.Generate {
		return nil, nil
	}

	readOperation, ok := input.Operations[input.Details.ReadMethod.MethodName]
	if !ok {
		return nil, fmt.Errorf("couldn't find read operation named %q", input.Details.ReadMethod.MethodName)
	}

	idParseLine, err := input.ParseResourceIdFuncName()
	if err != nil {
		return nil, fmt.Errorf("determining Parse function name for Resource ID: %+v", err)
	}

	terraformModel, ok := input.SchemaModels[input.SchemaModelName]
	if !ok {
		return nil, fmt.Errorf("the Schema Model named %q was not found", input.SchemaModelName)
	}

	resourceId, ok := input.ResourceIds[input.Details.ResourceIdName]
	if !ok {
		return nil, fmt.Errorf("the Resource ID named %q was not found", input.Details.ResourceIdName)
	}

	// at this point we only support References being returned, so this is safe
	topLevelModel, ok := input.Models[*readOperation.ResponseObject.ReferenceName]
	if !ok {
		return nil, fmt.Errorf("the top-level model %q used in the response was not found", *readOperation.ResponseObject.ReferenceName)
	}

	helper := readFunctionComponents{
		constants:       input.Constants,
		idParseLine:     *idParseLine,
		mappings:        input.Details.Mappings,
		readMethod:      input.Details.ReadMethod,
		readOperation:   readOperation,
		resourceId:      resourceId,
		schemaModelName: input.SchemaModelName,
		sdkResourceName: input.SdkResourceName,
		terraformModel:  terraformModel,
		topLevelModel:   topLevelModel,
	}
	components := []func() (*string, error){
		helper.codeForIDParser,
		helper.codeForGet,
		helper.codeForModelAssignments,
	}
	lines := make([]string, 0)
	for i, component := range components {
		result, err := component()
		if err != nil {
			return nil, fmt.Errorf("running component %d: %+v", i, err)
		}

		lines = append(lines, *result)
	}

	output := fmt.Sprintf(`
func (r %[1]sResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: %[2]d * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.%[3]s.%[4]s
			schema := %[5]s{}

			%[6]s

			return metadata.Encode(&schema)
		},
	}
}
`, input.ResourceTypeName, input.Details.ReadMethod.TimeoutInMinutes, input.ServiceName, input.SdkResourceName, input.SchemaModelName, strings.Join(lines, "\n\n"))
	return &output, nil
}

func (c readFunctionComponents) codeForIDParser() (*string, error) {
	output := fmt.Sprintf(`
			id, err := %[1]s(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
`, c.idParseLine)
	return &output, nil
}

func (c readFunctionComponents) codeForGet() (*string, error) {
	methodArguments := argumentsForApiOperationMethod(c.readOperation, c.sdkResourceName, c.readMethod.MethodName, true)
	output := fmt.Sprintf(`
			resp, err := client.%[1]s(%[2]s)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %%s: %%+v", *id, err)
			}
`, c.readMethod.MethodName, methodArguments)
	return &output, nil
}

func (c readFunctionComponents) codeForModelAssignments() (*string, error) {
	resourceIdMappings, err := c.codeForResourceIdMappings()
	if err != nil {
		return nil, fmt.Errorf("building code for resource id mappings: %+v", err)
	}
	// first map all the Resource ID segments across
	if len(c.terraformModel.Fields) == 0 && len(*resourceIdMappings) == 0 {
		// Nothing to see here, move along
		empty := ""
		return &empty, nil
	}
	// then output the top-level mappings, which'll call into nested items as required

	topLevelMappings, err := c.codeForTopLevelMappings()
	if err != nil {
		return nil, fmt.Errorf("building code for top-level field mappings: %+v", err)
	}
	output := fmt.Sprintf(`
			if model := resp.Model; model != nil {
				%[1]s
				%[2]s
			}
`, *resourceIdMappings, *topLevelMappings)
	return &output, nil
}

func (c readFunctionComponents) codeForResourceIdMappings() (*string, error) {
	lines := make([]string, 0)

	// TODO: note that when there's a parent ID field present we'll need to call `parent.NewParentID(..).ID()`
	// to get the right URI
	for _, v := range c.resourceId.Segments {
		if v.Type == resourcemanager.StaticSegment || v.Type == resourcemanager.SubscriptionIdSegment {
			continue
		}

		for _, resourceIdMapping := range c.mappings.ResourceId {
			if resourceIdMapping.SegmentName != v.Name {
				continue
			}

			// Constants are output into the Schema as their native types (e.g. int/float/string) so we need to convert prior to assigning
			if v.ConstantReference != nil {
				constant, ok := c.constants[*v.ConstantReference]
				if !ok {
					return nil, fmt.Errorf("the constant %q referenced in Resource ID Segment %q was not found", *v.ConstantReference, v.Name)
				}
				constantGoTypeName, err := golangFieldTypeFromConstantType(constant.Type)
				if err != nil {
					return nil, fmt.Errorf("determining Golang Type name for Constant Type %q: %+v", string(constant.Type), err)
				}
				lines = append(lines, fmt.Sprintf("schema.%s = %s(id.%s)", resourceIdMapping.SchemaFieldName, *constantGoTypeName, strings.Title(resourceIdMapping.SegmentName)))
			} else {
				lines = append(lines, fmt.Sprintf("schema.%s = id.%s", resourceIdMapping.SchemaFieldName, strings.Title(resourceIdMapping.SegmentName)))
			}
			break
		}
	}

	sort.Strings(lines)

	output := strings.Join(lines, "\n")
	return &output, nil
}

func (c readFunctionComponents) codeForTopLevelMappings() (*string, error) {
	// TODO: tests for this
	mappings := make([]string, 0)
	mappingsMap := make(map[string]string, 0)
	schemaPrefix := fmt.Sprintf("schema.")
	hasProperties := false
	for _, v := range c.mappings.Read {
		fieldName := v.DirectAssignment.SdkFieldPath
		if strings.HasSuffix(v.DirectAssignment.SchemaModelName, "Properties") {
			hasProperties = true
		}
		if v.DirectAssignment.SchemaModelName != c.sdkResourceName {
			// We only care about top level Items here...
			continue
		}
		modelPrefix := fmt.Sprintf("model.%s", v.DirectAssignment.SchemaFieldPath)
		temp, err := flattenAssignmentCodeForField(v, c.terraformModel.Fields[fieldName], schemaPrefix, modelPrefix)
		if err != nil {
			return nil, err
		}
		mappingsMap[fieldName] = *temp
	}

	orderedFieldNames := make([]string, 0)
	for fieldName := range c.terraformModel.Fields {
		orderedFieldNames = append(orderedFieldNames, fieldName)
	}
	sort.Strings(orderedFieldNames)
	for _, tfFieldName := range orderedFieldNames {
		if _, ok := mappingsMap[tfFieldName]; ok {
			mappings = append(mappings, fmt.Sprintf("%s", mappingsMap[tfFieldName]))
		}
	}

	if hasProperties {
		propsMappings := make([]string, 0)
		propsMappingsMap := make(map[string]string, 0)
		propertiesCode := ""
		// We only care if there's a top level Properties model here
		for _, v := range c.mappings.Read {
			fieldName := v.DirectAssignment.SdkFieldPath
			if v.DirectAssignment.SchemaModelName == fmt.Sprintf("%sProperties", c.sdkResourceName) {
				modelPrefix := fmt.Sprintf("model.Properties.%s", v.DirectAssignment.SchemaFieldPath)
				temp, err := flattenAssignmentCodeForField(v, c.terraformModel.Fields[fieldName], schemaPrefix, modelPrefix)
				if err != nil {
					return nil, err
				}
				propsMappingsMap[fieldName] = *temp
			}
		}
		orderedPropsNames := make([]string, 0)
		for fieldName := range c.terraformModel.Fields {
			orderedPropsNames = append(orderedPropsNames, fieldName)
		}
		sort.Strings(orderedPropsNames)
		for _, tfFieldName := range orderedPropsNames {
			if _, ok := propsMappingsMap[tfFieldName]; ok {
				propsMappings = append(propsMappings, fmt.Sprintf("%s", propsMappingsMap[tfFieldName]))
			}
		}

		propertiesCode = fmt.Sprintf(`
		if props := model.Properties; props != nil {
			%[1]s
		}
`, strings.Join(propsMappings, "\n"))
		mappings = append(mappings, propertiesCode)
	}

	output := strings.Join(mappings, "\n")

	return &output, nil
}
