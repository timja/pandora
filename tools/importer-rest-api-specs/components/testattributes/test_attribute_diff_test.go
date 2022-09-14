package testattributes

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
	"testing"
)

func TestDiff_TopLevelFields(t *testing.T) {
	testData := []struct {
		old      string
		new      resourcemanager.TerraformSchemaModelDefinition
		expected []string
	}{
		{
			old: `
	required_bool_attribute  = false
	required_float_attribute = 10.1
	optional_int_attribute = 15
	optional_string_attribute = "foo"
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"RequiredBoolAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Boolean",
						},
						HclName:  "required_bool_attribute",
						Required: true,
					},
					"RequiredFloatAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Float",
						},
						HclName:  "required_float_attribute",
						Required: true,
					},
					"OptionalIntAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Integer",
						},
						HclName:  "optional_int_attribute",
						Optional: true,
					},
					"OptionalStringAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "String",
						},
						HclName:  "optional_string_attribute",
						Optional: true,
					},
				},
			},
			expected: []string{},
		},
		{
			old: `
	required_bool_attribute  = false
	required_float_attribute = 10.1
	optional_int_attribute = 15
	optional_string_attribute = "foo"
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"RequiredFloatAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Float",
						},
						HclName:  "required_float_attribute",
						Required: true,
					},
					"OptionalIntAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Integer",
						},
						HclName:  "optional_int_attribute",
						Optional: true,
					},
					"OptionalStringAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "String",
						},
						HclName:  "optional_string_attribute",
						Optional: true,
					},
				},
			},
			expected: []string{"required_bool_attribute"},
		},
		{
			old: `
	required_bool_attribute  = false
	required_float_attribute = 10.1
	optional_int_attribute = 15
	optional_string_attribute = "foo"
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"RequiredFloatAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Float",
						},
						HclName:  "required_float_attribute",
						Required: true,
					},
					"OptionalIntAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Integer",
						},
						HclName:  "optional_int_attribute",
						Optional: true,
					},
				},
			},
			expected: []string{"optional_string_attribute", "required_bool_attribute"},
		},
		{
			old: `
	required_bool_attribute  = false
	required_float_attribute = 10.1
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"RequiredBoolAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Boolean",
						},
						HclName:  "required_bool_attribute",
						Required: true,
					},
					"RequiredFloatAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Float",
						},
						HclName:  "required_float_attribute",
						Required: true,
					},
					"OptionalIntAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "Integer",
						},
						HclName:  "optional_int_attribute",
						Optional: true,
					},
					"OptionalStringAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: "String",
						},
						HclName:  "optional_string_attribute",
						Optional: true,
					},
				},
			},
			expected: []string{},
		},
	}

	for i, testCase := range testData {
		file := hclwrite.NewEmptyFile()
		helper := TestAttributesHelpers{
			SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{},
		}
		err := helper.GetAttributesForTests(testCase.new, *file.Body(), false)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}

		diff, err := testAttributeDiff(testCase.old, *file.Body())
		if err != nil {
			t.Fatalf("unexpected error for index %d", i)
		}
		if !cmp.Equal(diff, testCase.expected, cmpopts.SortSlices(func(a, b string) bool { return a < b })) {
			t.Fatalf("diffs didn't match\nExpected: %+v\nActual: %+v", testCase.expected, diff)
		}
	}
}

func TestDiff_NestedFields(t *testing.T) {
	testData := []struct {
		old      string
		new      resourcemanager.TerraformSchemaModelDefinition
		expected []string
		helper   TestAttributesHelpers
	}{
		{
			old: `
	optional_set_attribute {
		required_nested_string = "foo"
		optional_nested_string = "foo"
	}
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalSetAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeList,
							NestedObject: &resourcemanager.TerraformSchemaFieldObjectDefinition{
								Type:          resourcemanager.TerraformSchemaFieldTypeReference,
								ReferenceName: stringPointer("NestedSchemaModel"),
							},
						},
						HclName:  "optional_set_attribute",
						Optional: true,
					},
				},
			},
			helper: TestAttributesHelpers{
				SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
					"NestedSchemaModel": {
						Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
							"RequiredNestedString": {
								ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
									Type: resourcemanager.TerraformSchemaFieldTypeString,
								},
								HclName:  "required_nested_string",
								Required: true,
							},
							"OptionalNestedString": {
								ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
									Type: resourcemanager.TerraformSchemaFieldTypeString,
								},
								HclName:  "optional_nested_string",
								Optional: true,
							},
						},
					},
				},
			},
			expected: []string{},
		},
		{
			old: `
	optional_set_attribute {
		required_nested_bool   = false
		required_nested_string = "foo"
		optional_nested_string = "foo"
	}
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{},
			},
			helper: TestAttributesHelpers{
				SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
					"NestedSchemaModel": {
						Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
							"RequiredNestedString": {
								ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
									Type: resourcemanager.TerraformSchemaFieldTypeString,
								},
								HclName:  "required_nested_string",
								Required: true,
							},
							"OptionalNestedString": {
								ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
									Type: resourcemanager.TerraformSchemaFieldTypeString,
								},
								HclName:  "optional_nested_string",
								Optional: true,
							},
						},
					},
				},
			},
			expected: []string{"optional_set_attribute"},
		},
		{
			old: `
	optional_set_attribute {
		required_nested_string = "foo"
		optional_nested_string = "foo"
	}
`,
			new: resourcemanager.TerraformSchemaModelDefinition{
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalSetAttribute": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeList,
							NestedObject: &resourcemanager.TerraformSchemaFieldObjectDefinition{
								Type:          resourcemanager.TerraformSchemaFieldTypeReference,
								ReferenceName: stringPointer("NestedSchemaModel"),
							},
						},
						HclName:  "optional_set_attribute",
						Optional: true,
					},
				},
			},
			helper: TestAttributesHelpers{
				SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
					"NestedSchemaModel": {
						Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
							"RequiredNestedString": {
								ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
									Type: resourcemanager.TerraformSchemaFieldTypeString,
								},
								HclName:  "required_nested_string",
								Required: true,
							},
						},
					},
				},
			},
			expected: []string{"optional_set_attribute.optional_nested_string"},
		},
	}

	for i, testCase := range testData {
		file := hclwrite.NewEmptyFile()
		err := testCase.helper.GetAttributesForTests(testCase.new, *file.Body(), false)
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}

		diff, err := testAttributeDiff(testCase.old, *file.Body())
		if err != nil {
			t.Fatalf("unexpected error for index %d", i)
		}
		if !cmp.Equal(diff, testCase.expected, cmpopts.SortSlices(func(a, b string) bool { return a < b })) {
			t.Fatalf("diffs didn't match\nExpected: %+v\nActual: %+v", testCase.expected, diff)
		}
	}
}
