package dataapigenerator

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/parser/cleanup"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func codeForTerraformSchemaDefinition(terraformNamespace string, details resourcemanager.TerraformResourceDetails) string {
	// TODO: output the Field Object Definition(details.SchemaModels[""].Fields[""].ObjectDefinition
	// using the FieldObjectDefinitionType in a method basically duplicating dotNetTypeNameForCustomType
	// @mbfrahry ^^^

	// TODO: schema models are available in details.SchemaModels
	// TODO: the main schema name is available in details.SchemaModelName

	// TODO: output the real schema

	schema := ""
	className := ""
	if details.DisplayName != "" {
		className = fmt.Sprintf("%sResourceSchema", details.ResourceName)
	} else {
		className = fmt.Sprintf("%sSchema", details.SchemaModelName)
	}

	for _, v := range details.SchemaModels {
		// Make the schema ordered
		fieldList := make([]string, 0, len(v.Fields))
		for f := range v.Fields {
			fieldList = append(fieldList, f)
		}
		sort.Strings(fieldList)
		for _, field := range fieldList {
			def := v.Fields[field]
			//}
			//for field, def := range v.Fields {
			fieldIsId := false
			fieldIsList := false
			fieldIsEnum := def.ObjectDefinition.Type == "Enum"
			if def.ObjectDefinition.NestedObject != nil {
				fieldIsEnum = def.ObjectDefinition.NestedObject.Type == "Enum"
			}
			// TODO - refactor all this into some form of "cleanup" function(s)
			operations := "{ get; set; }" // hardcoded for now - work this out later when we have enough data
			fieldType := strings.ToLower(string(def.ObjectDefinition.Type))
			fieldIsList = strings.EqualFold(fieldType, "list")
			ref := ""
			if def.ObjectDefinition.ReferenceName != nil && !fieldIsList {
				ref = *def.ObjectDefinition.ReferenceName
			} else if def.ObjectDefinition.NestedObject != nil && def.ObjectDefinition.NestedObject.ReferenceName != nil {
				ref = *def.ObjectDefinition.NestedObject.ReferenceName
			} else if def.ObjectDefinition.NestedObject != nil && def.ObjectDefinition.NestedObject.Type != "" {
				fieldType = string(def.ObjectDefinition.NestedObject.Type)
			}

			// subresources are references to other IDs
			// TODO - Should be references to the remote ID, or just strings here?
			if strings.EqualFold(ref, "subresource") {
				log.Printf("%+v", ref)
				fieldIsId = true
				ref = cleanup.NormalizeName(field)
				fieldType = "string"
			}
			// fixup type names
			if strings.EqualFold(fieldType, "integer") {
				fieldType = "int"
			}

			if strings.EqualFold(fieldType, "boolean") {
				fieldType = "bool"
			}

			if strings.EqualFold(fieldType, "string") {
				fieldType = "string"
			}

			if fieldIsList {
				if ref != "" {
					if !fieldIsId && !fieldIsEnum {
						fieldType = fmt.Sprintf("List<%sSchema>?", ref)
					} else if fieldIsEnum {
						fieldType = fmt.Sprintf("List<%sConstant>?", ref)
					} else {
						fieldType = fmt.Sprintf("List<%s>?", fieldType)
					}
				} else {
					fieldType = fmt.Sprintf("List<%s>?", fieldType) // TODO deal with primitives?
				}
			}

			// We have custom type for tags
			if strings.EqualFold(fieldType, "tags") {
				fieldType = "CustomTypes.Tags"
			}

			// We have custom type for location
			if strings.EqualFold(fieldType, "location") {
				fieldType = "CustomTypes.Location"
			}

			if strings.EqualFold(field, "identity") {
				switch fieldType {
				// TODO - the other identity types...
				case "identitysystemanduserassigned":
					fieldType = "CustomTypes.LegacySystemAndUserAssignedIdentityMap"
				}
			}

			if strings.EqualFold(fieldType, "datetime") {
				fieldType = "string"
			}

			if fieldType == "reference" && def.ObjectDefinition.ReferenceName != nil {
				fieldType = *def.ObjectDefinition.ReferenceName
			}

			schema += fmt.Sprintf("\t[HclName(\"%s\")]\n", field)

			if def.ForceNew {
				schema += fmt.Sprintf("\t[ForceNew]\n")
			}
			if def.Required {
				schema += fmt.Sprintf("\t[Required]\n")
			}
			if def.Computed {
				schema += fmt.Sprintf("\t[Computed]\n")
			}
			if def.Optional {
				schema += fmt.Sprintf("\t[Optional]\n")
			}

			if def.ObjectDefinition.ReferenceName != nil && !fieldIsId {
				if fieldIsEnum {
					schema += fmt.Sprintf("\tpublic %s %s %s\n", fmt.Sprintf("%sConstant?", ref), ref, operations)
				} else {
					schema += fmt.Sprintf("\tpublic %s %s %s\n", fmt.Sprintf("%sSchema", fieldType), ref, operations)
				}
			} else {
				schema += fmt.Sprintf("\tpublic %s %s %s\n", fieldType, cleanup.NormalizeName(field), operations)
			}

			schema += fmt.Sprintf("\n")
		}
	}
	log.Print(schema)

	return fmt.Sprintf(`using System.Collections.Generic;
using Pandora.Definitions.Attributes;

namespace %[1]s;

public class %[2]s
{
%s
}
`, terraformNamespace, className, schema)
}

//[HclName("location")]
//[ForceNew]
//[Required]
//public CustomTypes.Location Location { get; set; }
//
//[HclName("name")]
//[ForceNew]
//[Required]
//public string Name { get; set; }
//
//[HclName("tags")]
//[Optional]
//public CustomTypes.Tags Tags { get; set; }
//
//[HclName("host_id")]
//[Optional]
//[ForceNew]
//public string HostId { get; set; }
