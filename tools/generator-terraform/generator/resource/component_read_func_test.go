package resource

import (
	"github.com/hashicorp/pandora/tools/generator-terraform/featureflags"
	"testing"

	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

// TODO: re-introduce Mappings for Schema <-> SDK

func TestComponentReadFunc_CommonId_Disabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SdkResourceName:  "SdkResource",
		ServiceName:      "Resources",
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         false,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Name",
						SegmentName:     "resourceGroupName",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				CommonAlias: stringPointer("Subscription"),
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	if actual != nil {
		t.Fatalf("expected `actual` to be nil but got %q", *actual)
	}
}

func TestComponentReadFunc_RegularResourceId_Disabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SdkResourceName:  "SdkResource",
		ServiceName:      "Resources",
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         false,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Name",
						SegmentName:     "resourceGroupName",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				Segments: []resourcemanager.ResourceIdSegment{},
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	if actual != nil {
		t.Fatalf("expected `actual` to be nil but got %q", *actual)
	}
}

func TestComponentReadFunc_CommonId_Enabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SdkResourceName:  "SdkResource",
		ServiceName:      "Resources",
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         true,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Name",
						SegmentName:     "resourceGroupName",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
				ResponseObject: &resourcemanager.ApiObjectDefinition{
					Type:          resourcemanager.ReferenceApiObjectDefinitionType,
					ReferenceName: stringPointer("GetModel"),
				},
			},
		},
		Models: map[string]resourcemanager.ModelDetails{
			"GetModel": {
				Fields: map[string]resourcemanager.FieldDetails{
					"Name": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "name",
					},
					"SomeSdkField": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "someSdkField",
					},
				},
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				CommonAlias: stringPointer("Subscription"),
				Id:          "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}",
				Segments: []resourcemanager.ResourceIdSegment{
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "subscriptions",
						FixedValue: stringPointer("subscriptions"),
					},
					{
						Type: resourcemanager.SubscriptionIdSegment,
						Name: "subscriptionId",
					},
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "resourceGroups",
						FixedValue: stringPointer("resourceGroups"),
					},
					{
						Type: resourcemanager.ResourceGroupSegment,
						Name: "resourceGroupName",
					},
				},
			},
		},
		SchemaModelName: "ExampleModel",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"ExampleModel": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Name": {
						HclName: "name",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
					},
					"SomeField": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
						HclName:  "some_field",
					},
				},
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
func (r ExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        Timeout: 10 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resources.SdkResource
			schema := ExampleModel{}
			id, err := commonids.ParseSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if model := resp.Model; model != nil {
				schema.Name = id.ResourceGroupName
			}
			return metadata.Encode(&schema)
        },
	}
}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_CommonId_Options_Enabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName:   "Example",
		SdkResourceName:    "SdkResource",
		ServiceName:        "Resources",
		ServicePackageName: "sdkservicepackage",
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         true,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Name",
						SegmentName:     "resourceGroupName",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
				ResponseObject: &resourcemanager.ApiObjectDefinition{
					Type:          resourcemanager.ReferenceApiObjectDefinitionType,
					ReferenceName: stringPointer("GetModel"),
				},
				Options: map[string]resourcemanager.ApiOperationOption{
					"SomeOption": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						HeaderName: stringPointer("X-Some-Option"),
						Required:   false,
					},
				},
			},
		},
		Models: map[string]resourcemanager.ModelDetails{
			"GetModel": {
				Fields: map[string]resourcemanager.FieldDetails{
					"Name": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "name",
					},
					"SomeSdkField": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "someSdkField",
					},
				},
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				CommonAlias: stringPointer("Subscription"),
				Id:          "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}",
				Segments: []resourcemanager.ResourceIdSegment{
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "subscriptions",
						FixedValue: stringPointer("subscriptions"),
					},
					{
						Type: resourcemanager.SubscriptionIdSegment,
						Name: "subscriptionId",
					},
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "resourceGroups",
						FixedValue: stringPointer("resourceGroups"),
					},
					{
						Type: resourcemanager.ResourceGroupSegment,
						Name: "resourceGroupName",
					},
				},
			},
		},
		SchemaModelName: "ExampleModel",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"ExampleModel": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Name": {
						HclName: "name",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
					},
					"SomeField": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
						HclName:  "some_field",
					},
				},
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
func (r ExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        Timeout: 10 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resources.SdkResource
			schema := ExampleModel{}
			id, err := commonids.ParseSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, sdkresource.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if model := resp.Model; model != nil {
				schema.Name = id.ResourceGroupName
			}
			return metadata.Encode(&schema)
        },
	}
}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_RegularResourceId_Enabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SdkResourceName:  "SdkResource",
		ServiceName:      "Resources",
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         true,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Name",
						SegmentName:     "resourceGroupName",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
				ResponseObject: &resourcemanager.ApiObjectDefinition{
					Type:          resourcemanager.ReferenceApiObjectDefinitionType,
					ReferenceName: stringPointer("GetModel"),
				},
			},
		},
		Models: map[string]resourcemanager.ModelDetails{
			"GetModel": {
				Fields: map[string]resourcemanager.FieldDetails{
					"Name": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "name",
					},
					"SomeSdkField": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "someSdkField",
					},
				},
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				Id: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}",
				Segments: []resourcemanager.ResourceIdSegment{
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "subscriptions",
						FixedValue: stringPointer("subscriptions"),
					},
					{
						Type: resourcemanager.SubscriptionIdSegment,
						Name: "subscriptionId",
					},
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "resourceGroups",
						FixedValue: stringPointer("resourceGroups"),
					},
					{
						Type: resourcemanager.ResourceGroupSegment,
						Name: "resourceGroupName",
					},
				},
			},
		},
		SchemaModelName: "ExampleModel",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"ExampleModel": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Name": {
						HclName: "name",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
					},
					"SomeField": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
						HclName:  "some_field",
					},
				},
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
func (r ExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        Timeout: 10 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resources.SdkResource
			schema := ExampleModel{}
			id, err := sdkresource.ParseCustomSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if model := resp.Model; model != nil {
				schema.Name = id.ResourceGroupName
			}
			return metadata.Encode(&schema)
        },
	}
}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_RegularResourceId_Constant_Enabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SdkResourceName:  "SdkResource",
		ServiceName:      "Resources",
		Constants: map[string]resourcemanager.ConstantDetails{
			"AnimalType": {
				Type: resourcemanager.StringConstant,
				Values: map[string]string{
					"Cow":   "Cow",
					"Panda": "Panda",
				},
			},
		},
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         true,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Animal",
						SegmentName:     "animalType",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
				ResponseObject: &resourcemanager.ApiObjectDefinition{
					Type:          resourcemanager.ReferenceApiObjectDefinitionType,
					ReferenceName: stringPointer("GetModel"),
				},
			},
		},
		Models: map[string]resourcemanager.ModelDetails{
			"GetModel": {
				Fields: map[string]resourcemanager.FieldDetails{
					"Name": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "name",
					},
					"SomeSdkField": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "someSdkField",
					},
				},
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				Id: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}",
				Segments: []resourcemanager.ResourceIdSegment{
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "subscriptions",
						FixedValue: stringPointer("subscriptions"),
					},
					{
						Type: resourcemanager.SubscriptionIdSegment,
						Name: "subscriptionId",
					},
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "animals",
						FixedValue: stringPointer("animals"),
					},
					{
						Type:              resourcemanager.ConstantSegment,
						Name:              "animalType",
						ConstantReference: stringPointer("AnimalType"),
					},
				},
			},
		},
		SchemaModelName: "ExampleModel",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"ExampleModel": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Animal": {
						HclName: "animal",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
					},
					"SomeField": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
						HclName:  "some_field",
					},
				},
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
func (r ExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        Timeout: 10 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resources.SdkResource
			schema := ExampleModel{}
			id, err := sdkresource.ParseCustomSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if model := resp.Model; model != nil {
				schema.Animal = string(id.AnimalType)
			}
			return metadata.Encode(&schema)
        },
	}
}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_RegularResourceId_Options_Enabled(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName:   "Example",
		SdkResourceName:    "SdkResource",
		ServiceName:        "Resources",
		ServicePackageName: "sdkservicepackage",
		Details: resourcemanager.TerraformResourceDetails{
			ReadMethod: resourcemanager.MethodDefinition{
				Generate:         true,
				MethodName:       "Get",
				TimeoutInMinutes: 10,
			},
			ResourceIdName: "CustomSubscriptionId",
			Mappings: resourcemanager.MappingDefinition{
				ResourceId: []resourcemanager.ResourceIdMappingDefinition{
					{
						SchemaFieldName: "Name",
						SegmentName:     "resourceGroupName",
					},
				},
			},
		},
		Operations: map[string]resourcemanager.ApiOperation{
			"Get": {
				LongRunning:    false,
				ResourceIdName: stringPointer("CustomSubscriptionId"),
				ResponseObject: &resourcemanager.ApiObjectDefinition{
					Type:          resourcemanager.ReferenceApiObjectDefinitionType,
					ReferenceName: stringPointer("GetModel"),
				},
				Options: map[string]resourcemanager.ApiOperationOption{
					"SomeOption": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						HeaderName: stringPointer("X-Some-Option"),
						Required:   false,
					},
				},
			},
		},
		Models: map[string]resourcemanager.ModelDetails{
			"GetModel": {
				Fields: map[string]resourcemanager.FieldDetails{
					"Name": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "name",
					},
					"SomeSdkField": {
						ObjectDefinition: resourcemanager.ApiObjectDefinition{
							Type: resourcemanager.StringApiObjectDefinitionType,
						},
						Required: true,
						JsonName: "someSdkField",
					},
				},
			},
		},
		ResourceIds: map[string]resourcemanager.ResourceIdDefinition{
			"CustomSubscriptionId": {
				Id: "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}",
				Segments: []resourcemanager.ResourceIdSegment{
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "subscriptions",
						FixedValue: stringPointer("subscriptions"),
					},
					{
						Type: resourcemanager.SubscriptionIdSegment,
						Name: "subscriptionId",
					},
					{
						Type:       resourcemanager.StaticSegment,
						Name:       "resourceGroups",
						FixedValue: stringPointer("resourceGroups"),
					},
					{
						Type: resourcemanager.ResourceGroupSegment,
						Name: "resourceGroupName",
					},
				},
			},
		},
		SchemaModelName: "ExampleModel",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"ExampleModel": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Name": {
						HclName: "name",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
					},
					"SomeField": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						ForceNew: true,
						HclName:  "some_field",
					},
				},
			},
		},
	}
	actual, err := readFunctionForResource(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
func (r ExampleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
        Timeout: 10 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resources.SdkResource
			schema := ExampleModel{}
			id, err := sdkresource.ParseCustomSubscriptionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id, sdkresource.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if model := resp.Model; model != nil {
				schema.Name = id.ResourceGroupName
			}
			return metadata.Encode(&schema)
        },
	}
}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_CodeForIDParser(t *testing.T) {
	actual, err := readFunctionComponents{
		idParseLine: "plz.Parse",
	}.codeForIDParser()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
	id, err := plz.Parse(metadata.ResourceData.Id())
	if err != nil {
		return err
	}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_CodeForGet(t *testing.T) {
	actual, err := readFunctionComponents{
		readMethod: resourcemanager.MethodDefinition{
			Generate:         true,
			MethodName:       "Get",
			TimeoutInMinutes: 5,
		},
		readOperation: resourcemanager.ApiOperation{
			ResourceIdName: stringPointer("SomeResourceId"),
		},
		sdkResourceName: "SdkResource",
	}.codeForGet()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
		resp, err := client.Get(ctx, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return metadata.MarkAsGone(*id)
			}
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_CodeForGet_Options(t *testing.T) {
	actual, err := readFunctionComponents{
		readMethod: resourcemanager.MethodDefinition{
			Generate:         true,
			MethodName:       "Get",
			TimeoutInMinutes: 5,
		},
		readOperation: resourcemanager.ApiOperation{
			ResourceIdName: stringPointer("SomeResourceId"),
			Options: map[string]resourcemanager.ApiOperationOption{
				"Example": {},
			},
		},
		sdkResourceName: "SdkResource",
	}.codeForGet()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := `
		resp, err := client.Get(ctx, *id, sdkresource.DefaultGetOperationOptions())
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return metadata.MarkAsGone(*id)
			}
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
`
	assertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentReadFunc_MappingsToSdk_NoFields(t *testing.T) {
	// TODO: remove this once the feature-flag is properly threaded through
	if !featureflags.OutputMappings {
		t.Skip("skipping since `featureflags.OutputMappings` is disabled")
	}

	actual, err := readFunctionComponents{
		terraformModel: resourcemanager.TerraformSchemaModelDefinition{
			Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{},
		},
	}.codeForModelAssignments()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	expected := ``
	assertTemplatedCodeMatches(t, expected, *actual)

}

func TestComponentReadFunc_MappingsToSdk_FromTopLevel(t *testing.T) {
	actual, err := readFunctionComponents{
		sdkResourceName: "ResourceGroup",
		terraformModel: resourcemanager.TerraformSchemaModelDefinition{
			Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
				"Description": {
					ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
						Type: resourcemanager.TerraformSchemaFieldTypeString,
					},
					Optional: true,
				},
				"AnOtherSetting": {
					ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
						Type: resourcemanager.TerraformSchemaFieldTypeString,
					},
					Optional: true,
				},
				"Location": {
					ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
						Type: resourcemanager.TerraformSchemaFieldTypeLocation,
					},
					Required: true,
				},
				"Tags": {
					ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
						Type: resourcemanager.TerraformSchemaFieldTypeTags,
					},
					Required: true,
				},
			},
		},
		mappings: resourcemanager.MappingDefinition{
			Read: []resourcemanager.FieldMappingDefinition{
				{
					Type: resourcemanager.DirectAssignmentMappingDefinitionType,
					DirectAssignment: &resourcemanager.FieldMappingDirectAssignmentDefinition{
						SchemaModelName: "ResourceGroup",
						SchemaFieldPath: "Location",
						SdkFieldPath:    "Location",
						SdkModelName:    "ResourceGroupResourceSchema",
					},
				},
				{
					Type: resourcemanager.DirectAssignmentMappingDefinitionType,
					DirectAssignment: &resourcemanager.FieldMappingDirectAssignmentDefinition{
						SchemaFieldPath: "Tags",
						SchemaModelName: "ResourceGroup",
						SdkFieldPath:    "Tags",
						SdkModelName:    "ResourceGroupResourceSchema",
					},
				},
				{
					Type: resourcemanager.DirectAssignmentMappingDefinitionType,
					DirectAssignment: &resourcemanager.FieldMappingDirectAssignmentDefinition{
						SchemaFieldPath: "Description",
						SchemaModelName: "ResourceGroupProperties",
						SdkFieldPath:    "Description",
						SdkModelName:    "ResourceGroupResourceSchema",
					},
				},
				{
					Type: resourcemanager.DirectAssignmentMappingDefinitionType,
					DirectAssignment: &resourcemanager.FieldMappingDirectAssignmentDefinition{
						SchemaFieldPath: "AnOtherSetting",
						SchemaModelName: "ResourceGroupProperties",
						SdkFieldPath:    "AnOtherSetting",
						SdkModelName:    "ResourceGroupResourceSchema",
					},
				},
			},
		},
	}.codeForModelAssignments()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	expected := `
	if model := resp.Model; model != nil {
        schema.Location = location.Normalize(model.Location)
        schema.Tags = tags.Flatten(model.Tags)
        if props := model.Properties; props != nil {
    	    schema.AnOtherSetting = model.Properties.AnOtherSetting
        	schema.Description = model.Properties.Description
        	}
        }
`
	assertTemplatedCodeMatches(t, expected, *actual)
}
