package dataapigenerator

import (
	"fmt"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
	"log"
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
				schema += fmt.Sprintf("\tpublic %s %s\n", def.ObjectDefinition.Type, *def.ObjectDefinition.ReferenceName)
			}
			schema += fmt.Sprintf("\n")
		}
	}
	log.Print(schema)

	return fmt.Sprintf(`using Pandora.Definitions.Attributes;

namespace %[1]s;

public class %[2]sResourceSchema
{
	// TODO: populate with a real schema

    [HclName("location")]
    [ForceNew]
    [Required]
    public CustomTypes.Location Location { get; set; }

    [HclName("name")]
    [ForceNew]
    [Required]
    public string Name { get; set; } 
    
    [HclName("tags")]
    [Optional]
    public CustomTypes.Tags Tags { get; set; }

    [HclName("host_id")]
    [Optional]
	[ForceNew]
    public string HostId { get; set; }

%s

}
`, terraformNamespace, details.ResourceName, schema)
}
