package dataapigenerator

import (
	"fmt"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/parser/cleanup"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
	"log"
	"strings"
)

func codeForTerraformSchemaDefinition(terraformNamespace string, details resourcemanager.TerraformResourceDetails) string {
	// TODO: output the Field Object Definition(details.SchemaModels[""].Fields[""].ObjectDefinition
	// using the FieldObjectDefinitionType in a method basically duplicating dotNetTypeNameForCustomType
	// @mbfrahry ^^^

	// TODO: schema models are available in details.SchemaModels
	// TODO: the main schema name is available in details.SchemaModelName

	// TODO: output the real schema

	schema := ""
	for _, v := range details.SchemaModels {
		for field, def := range v.Fields {
			// TODO - refactor all this into some form of "cleanup" function(s)
			operations := "{ get; set; }" // hardcoded for now - work this out later when we have enough data
			fieldType := strings.ToLower(string(def.ObjectDefinition.Type))
			ref := ""
			if def.ObjectDefinition.ReferenceName != nil {
				ref = *def.ObjectDefinition.ReferenceName
			}

			// subresources are references to other IDs
			// TODO - Should be references to the remote ID, or just strings here?
			if strings.EqualFold(ref, "subresource") {
				log.Printf("%+v", ref)
				ref = cleanup.NormalizeName(field)
				fieldType = ref
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

			// fixup type names
			if strings.EqualFold(fieldType, "integer") {
				fieldType = "int"
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
			if def.ObjectDefinition.ReferenceName != nil {
				schema += fmt.Sprintf("\tpublic %s %s %s\n", fieldType, ref, operations)
			} else {
				schema += fmt.Sprintf("\tpublic %s %s %s\n", fieldType, cleanup.NormalizeName(field), operations)
			}

			schema += fmt.Sprintf("\n")
		}
	}
	log.Print(schema)

	return fmt.Sprintf(`using Pandora.Definitions.Attributes;

namespace %[1]s;

public class %[2]sResourceSchema
{

%s

}
`, terraformNamespace, details.ResourceName, schema)
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
