package docs

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
	"github.com/hashicorp/pandora/tools/sdk/testhelpers"
)

// TODO: Lists
// TODO: above/below

func TestComponentBlocks_ModelsSingle(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SchemaModelName:  "TopLevelModelResourceSchema",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"TopLevelModelResourceSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"RequiredNestedItem": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("RequiredNestedSchema"),
						},
						Required: true,
						HclName:  "required_nested_item",
					},
				},
			},
			"RequiredNestedSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalItem": {
						HclName: "optional_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for optional_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"ComputedItem": {
						HclName: "computed_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_item.",
						},
					},
					"RequiredItem": {
						HclName: "required_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for required_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"BooleanItem": {
						HclName: "boolean_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeBoolean,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for boolean_item.",
						},
					},
				},
			},
		},
	}
	actual, err := codeForBlocksReference(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	expected := strings.ReplaceAll(`## Blocks Reference

### 'required_nested_item' Block

The 'required_nested_item' block supports the following arguments:

* 'required_item' - (Required) Description for required_item. Possible values are 'string1', 'string2' and 'string3'.

* 'boolean_item' - (Optional) Description for boolean_item. Possible values are 'true' and 'false'.

* 'optional_item' - (Optional) Description for optional_item. Possible values are 'string1', 'string2' and 'string3'.

The 'required_nested_item' block exports the following attributes:

* 'computed_item' - Description for computed_item.
`, "'", "`")

	testhelpers.AssertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentBlocks_ModelsMultiple(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SchemaModelName:  "TopLevelModelResourceSchema",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"TopLevelModelResourceSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalNestedItem": {
						Optional: true,
						HclName:  "optional_nested_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("NestedSchema"),
						},
					},
					"RequiredInteger": {
						HclName: "required_integer",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeInteger,
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for required_integer.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{1, 2, 3},
							},
						},
					},
					"OptionalInteger": {
						HclName: "optional_integer",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeInteger,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for optional_integer.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{4, 5, 6},
							},
						},
					},
					"ComputedInteger": {
						HclName: "computed_integer",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeInteger,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_integer.",
						},
					},
					"RequiredNestedItem": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("RequiredNestedSchema"),
						},
						Required: true,
						HclName:  "required_nested_item",
					},
					"RequiredString": {
						HclName: "required_string",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for required_string.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"OptionalString": {
						HclName: "optional_string",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for optional_string.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"ComputedString": {
						HclName: "computed_string",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_string.",
						},
					},
				},
			},
			"NestedSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"NestedItem": {
						HclName: "nested_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for nested_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"AnotherNestedItem": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("AnotherNestedSchema"),
						},
						Optional: true,
						HclName:  "another_nested_item",
					},
					"ZAnotherNestedItem2": {
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("ZAnotherNestedSchema2"),
						},
						Optional: true,
						HclName:  "z_another_nested_item_2",
					},
					"ComputedItem": {
						HclName: "computed_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_string.",
						},
					},
				},
			},
			"RequiredNestedSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalItem": {
						HclName: "optional_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for optional_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"ComputedItem": {
						HclName: "computed_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_item.",
						},
					},
					"RequiredItem": {
						HclName: "required_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for required_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
					"BooleanItem": {
						HclName: "boolean_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeBoolean,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for boolean_item.",
						},
					},
				},
			},
			"AnotherNestedSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalItem": {
						HclName: "optional_item",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for optional_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string1", "string2", "string3"},
							},
						},
					},
				},
			},
			"ZAnotherNestedSchema2": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OptionalItem2": {
						HclName: "optional_item_2",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for optional_item.",
						},
						Validation: &resourcemanager.TerraformSchemaValidationDefinition{
							Type: "PossibleValues",
							PossibleValues: &resourcemanager.TerraformSchemaValidationPossibleValuesDefinition{
								Values: []interface{}{"string4", "string5", "string6"},
							},
						},
					},
				},
			},
		},
	}
	actual, err := codeForBlocksReference(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	expected := strings.ReplaceAll(`## Blocks Reference

### 'another_nested_item' Block

The 'another_nested_item' block supports the following arguments:

* 'optional_item' - (Optional) Description for optional_item. Possible values are 'string1', 'string2' and 'string3'.

---

### 'optional_nested_item' Block

The 'optional_nested_item' block supports the following arguments:

* 'another_nested_item' - (Optional) An 'another_nested_item' block as defined below. 

* 'nested_item' - (Optional) Description for nested_item. Possible values are 'string1', 'string2' and 'string3'.

* 'z_another_nested_item_2' - (Optional) A 'z_another_nested_item_2' block as defined below. 

The 'optional_nested_item' block exports the following attributes:

* 'computed_item' - Description for computed_string.

---

### 'required_nested_item' Block

The 'required_nested_item' block supports the following arguments:

* 'required_item' - (Required) Description for required_item. Possible values are 'string1', 'string2' and 'string3'.

* 'boolean_item' - (Optional) Description for boolean_item. Possible values are 'true' and 'false'.

* 'optional_item' - (Optional) Description for optional_item. Possible values are 'string1', 'string2' and 'string3'.

The 'required_nested_item' block exports the following attributes:

* 'computed_item' - Description for computed_item.

---

### 'z_another_nested_item_2' Block

The 'z_another_nested_item_2' block supports the following arguments:

* 'optional_item_2' - (Optional) Description for optional_item. Possible values are 'string4', 'string5' and 'string6'.

`, "'", "`")

	testhelpers.AssertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentBlocks_ModelsMultipleNestingTheSameModelShouldBeDeduped(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SchemaModelName:  "TopLevelModelResourceSchema",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"TopLevelModelResourceSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"First": {
						HclName: "first",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("First"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for first.",
						},
					},
					"Second": {
						HclName: "second",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("Second"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for second.",
						},
					},
				},
			},
			"Example": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"SomeField": {
						HclName: "some_field",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for some_field.",
						},
					},
					"ComputedField": {
						HclName: "computed_field",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_field.",
						},
					},
				},
			},
			"First": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Example": {
						HclName: "example",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("Example"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for example.",
						},
					},
				},
			},
			"Second": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Example": {
						HclName: "example",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("Example"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for example.",
						},
					},
				},
			},
		},
	}
	actual, err := codeForBlocksReference(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	expected := strings.ReplaceAll(`## Blocks Reference

### 'example' Block

The 'example' block supports the following arguments:

* 'some_field' - (Optional) Description for some_field. 

The 'example' block exports the following attributes:

* 'computed_field' - Description for computed_field.

---

### 'first' Block

The 'first' block supports the following arguments:

* 'example' - (Required) An 'example' block as defined below. Description for example.

---

### 'second' Block

The 'second' block supports the following arguments:

* 'example' - (Required) An 'example' block as defined below. Description for example.

`, "'", "`")

	testhelpers.AssertTemplatedCodeMatches(t, expected, *actual)
}

func TestComponentBlocks_ModelsMultipleNestingDifferentModelsWithTheSameNameShouldBeOutputUniquely(t *testing.T) {
	input := models.ResourceInput{
		ResourceTypeName: "Example",
		SchemaModelName:  "TopLevelModelResourceSchema",
		SchemaModels: map[string]resourcemanager.TerraformSchemaModelDefinition{
			"TopLevelModelResourceSchema": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"First": {
						HclName: "first",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("First"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for first.",
						},
					},
					"Second": {
						HclName: "second",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("Second"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for second.",
						},
					},
				},
			},
			"Example1": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"SomeField": {
						HclName: "some_field",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for some_field.",
						},
					},
					"ComputedField": {
						HclName: "computed_field",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_field.",
						},
					},
				},
			},
			"Example2": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"OtherField": {
						HclName: "other_field",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Optional: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for other_field.",
						},
					},
					"ComputedField": {
						HclName: "computed_field",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type: resourcemanager.TerraformSchemaFieldTypeString,
						},
						Computed: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for computed_field.",
						},
					},
				},
			},
			"First": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Example": {
						HclName: "example",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("Example1"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for example.",
						},
					},
				},
			},
			"Second": {
				Fields: map[string]resourcemanager.TerraformSchemaFieldDefinition{
					"Example": {
						HclName: "example",
						ObjectDefinition: resourcemanager.TerraformSchemaFieldObjectDefinition{
							Type:          resourcemanager.TerraformSchemaFieldTypeReference,
							ReferenceName: pointer.To("Example2"),
						},
						Required: true,
						Documentation: resourcemanager.TerraformSchemaDocumentationDefinition{
							Markdown: "Description for example.",
						},
					},
				},
			},
		},
	}
	actual, err := codeForBlocksReference(input)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	expected := strings.ReplaceAll(`## Blocks Reference

### 'example' Block (within the 'first' block)

The 'example' block within the 'first' block supports the following arguments:

* 'some_field' - (Optional) Description for some_field. 

The 'example' block within the 'first' block exports the following attributes:

* 'computed_field' - Description for computed_field.

---

### 'example' Block (within the 'second' block)

The 'example' block within the 'second' block supports the following arguments:

* 'other_field' - (Optional) Description for other_field. 

The 'example' block within the 'second' block exports the following attributes:

* 'computed_field' - Description for computed_field.

---

### 'first' Block

The 'first' block supports the following arguments:

* 'example' - (Required) An 'example' block as defined below. Description for example.

---

### 'second' Block

The 'second' block supports the following arguments:

* 'example' - (Required) An 'example' block as defined below. Description for example.

`, "'", "`")

	testhelpers.AssertTemplatedCodeMatches(t, expected, *actual)
}
