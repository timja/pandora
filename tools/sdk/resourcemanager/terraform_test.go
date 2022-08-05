package resourcemanager

import (
	"testing"
)

func TestObjectDefinitionToGolangFieldType(t *testing.T) {
	testData := []struct {
		input    TerraformSchemaFieldObjectDefinition
		expected *string
	}{
		// Simple Types
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeBoolean,
			},
			expected: stringPointer("bool"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeDateTime,
			},
			// whilst DateTime could be output as *time.Time the Go SDK outputs
			// this as a String with Get/Set methods to allow exposing this value
			// either as a raw string or by formatting the value, so we output
			// a string here rather than a *time.Time
			expected: stringPointer("string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeFloat,
			},
			expected: stringPointer("float64"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeInteger,
			},
			expected: stringPointer("int64"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeString,
			},
			expected: stringPointer("string"),
		},

		// References
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type:          TerraformSchemaFieldTypeReference,
				ReferenceName: stringPointer("SomeModel"),
			},
			expected: stringPointer("[]SomeModel"),
		},

		// Common Types
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeEdgeZone,
			},
			expected: stringPointer("string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeLocation,
			},
			expected: stringPointer("string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeIdentitySystemAssigned,
			},
			expected: stringPointer("identity.ModelSystemAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeIdentitySystemAndUserAssigned,
			},
			expected: stringPointer("identity.ModelSystemAssignedUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeIdentitySystemOrUserAssigned,
			},
			expected: stringPointer("identity.ModelSystemAssignedUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeIdentityUserAssigned,
			},
			expected: stringPointer("identity.ModelUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeResourceGroup,
			},
			expected: stringPointer("string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeTags,
			},
			expected: stringPointer("map[string]interface{}"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeZone,
			},
			expected: stringPointer("string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeZones,
			},
			expected: stringPointer("[]string"),
		},
	}
	for i, data := range testData {
		result, err := data.input.GolangFieldType()
		if err != nil {
			if data.expected == nil {
				continue
			}

			t.Fatalf("unexpected error for iteration %d: %+v", i, err)
		}
		if data.expected == nil {
			t.Fatalf("expected an error but didn't get one for iteration %d", i)
		}
		if result == nil {
			t.Fatalf("expected no error and a result but got no error and no result")
		}
		if *result != *data.expected {
			t.Fatalf("expected %q but got %q", *data.expected, *result)
		}
	}
}

func TestObjectDefinitionToGolangFieldType_Lists(t *testing.T) {
	testData := []struct {
		input    TerraformSchemaFieldObjectDefinition
		expected *string
	}{
		// Simple Types
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeBoolean,
				},
			},
			expected: stringPointer("[]bool"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeDateTime,
				},
			},
			// whilst DateTime could be output as *time.Time the Go SDK outputs
			// this as a String with Get/Set methods to allow exposing this value
			// either as a raw string or by formatting the value, so we output
			// a string here rather than a *time.Time
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeFloat,
				},
			},
			expected: stringPointer("[]float64"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeInteger,
				},
			},
			expected: stringPointer("[]int64"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeString,
				},
			},
			expected: stringPointer("[]string"),
		},

		// References
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type:          TerraformSchemaFieldTypeReference,
					ReferenceName: stringPointer("SomeModel"),
				},
			},
			expected: stringPointer("[]SomeModel"),
		},

		// Common Types
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeEdgeZone,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeLocation,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentitySystemAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelSystemAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentitySystemAndUserAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelSystemAssignedUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentitySystemOrUserAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelSystemAssignedUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentityUserAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeResourceGroup,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeTags,
				},
			},
			expected: stringPointer("[]map[string]interface{}"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeZone,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeList,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeZones,
				},
			},
			expected: stringPointer("[][]string"),
		},
	}
	for i, data := range testData {
		result, err := data.input.GolangFieldType()
		if err != nil {
			if data.expected == nil {
				continue
			}

			t.Fatalf("unexpected error for iteration %d: %+v", i, err)
		}
		if data.expected == nil {
			t.Fatalf("expected an error but didn't get one for iteration %d", i)
		}
		if result == nil {
			t.Fatalf("expected no error and a result but got no error and no result")
		}
		if *result != *data.expected {
			t.Fatalf("expected %q but got %q", *data.expected, *result)
		}
	}
}

func TestObjectDefinitionToGolangFieldType_Sets(t *testing.T) {
	testData := []struct {
		input    TerraformSchemaFieldObjectDefinition
		expected *string
	}{
		// Simple Types
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeBoolean,
				},
			},
			expected: stringPointer("[]bool"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeDateTime,
				},
			},
			// whilst DateTime could be output as *time.Time the Go SDK outputs
			// this as a String with Get/Set methods to allow exposing this value
			// either as a raw string or by formatting the value, so we output
			// a string here rather than a *time.Time
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeFloat,
				},
			},
			expected: stringPointer("[]float64"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeInteger,
				},
			},
			expected: stringPointer("[]int64"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeString,
				},
			},
			expected: stringPointer("[]string"),
		},

		// References
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type:          TerraformSchemaFieldTypeReference,
					ReferenceName: stringPointer("SomeModel"),
				},
			},
			expected: stringPointer("[]SomeModel"),
		},

		// Common Types
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeEdgeZone,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeLocation,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentitySystemAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelSystemAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentitySystemAndUserAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelSystemAssignedUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentitySystemOrUserAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelSystemAssignedUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeIdentityUserAssigned,
				},
			},
			expected: stringPointer("[]identity.ModelUserAssigned"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeResourceGroup,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeTags,
				},
			},
			expected: stringPointer("[]map[string]interface{}"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeZone,
				},
			},
			expected: stringPointer("[]string"),
		},
		{
			input: TerraformSchemaFieldObjectDefinition{
				Type: TerraformSchemaFieldTypeSet,
				NestedObject: &TerraformSchemaFieldObjectDefinition{
					Type: TerraformSchemaFieldTypeZones,
				},
			},
			expected: stringPointer("[][]string"),
		},
	}
	for i, data := range testData {
		result, err := data.input.GolangFieldType()
		if err != nil {
			if data.expected == nil {
				continue
			}

			t.Fatalf("unexpected error for iteration %d: %+v", i, err)
		}
		if data.expected == nil {
			t.Fatalf("expected an error but didn't get one for iteration %d", i)
		}
		if result == nil {
			t.Fatalf("expected no error and a result but got no error and no result")
		}
		if *result != *data.expected {
			t.Fatalf("expected %q but got %q", *data.expected, *result)
		}
	}
}

func stringPointer(input string) *string {
	return &input
}
