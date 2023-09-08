service "ChaosStudio" {
  terraform_package = "chaosstudio"

  api "2023-04-15-preview" {
    package "Targets" {
      definition "chaos_studio_target" {
        id = "/{scope}/providers/Microsoft.Chaos/targets/{targetName}"
        display_name = "Chaos Studio Target"
        website_subcategory = "Chaos Studio"
        description = "Manages a Chaos Studio Target"
        generate_create = false
        generate_update = false
        test_data {
          basic_variables {
            strings = {
              "name" = "Microsoft-StorageAccount"
              "scope_id" = "storage_account_id"
            }
          }
        }
        # TODO schema overrides
        # todo in the importer when we write the terraformdetails
#        schema_overrides = {
#          "name" = "target_type"
#          "scope" = "target_resource_id"
#        }
        # TODO documentation overrides
#        documentation_overrides = {
#          "target_type" = "The name of the Chaos Studio Target. This has the format of [publisher]-[targetType] e.g. `Microsoft-StorageAccount`. Add Az CLI command here. Force new bit appended here automatically"
#        }
      }
    }
#    package "Capabilities" {
#      definition "chaos_studio_capability" {
#        id = "/{scope}/providers/Microsoft.Chaos/targets/{targetName}/capabilities/{capabilityName}"
#        display_name = "Chaos Studio Capability"
#        website_subcategory = "Chaos Studio"
#        description = "Manages a Chaos Studio Capability"
#        generate_create = false
#        generate_update = false
#        test_data {
#          basic_variables {
#            strings = {
#              "name" = "TODO capability name"
#              "chaos_target_id" = "chaos_target_id"
#            }
#          }
#        }
#      }
#    }
    ## discussion with tom: probably needs to be handwritten
    # Users would only be able to enable one capability on a target per resource
    # by that I mean if they wanted to enable several capabilities on a target they'd need to define
    # multiple of these resources using the same chaos_target_id
    # I actually think a better use experience would be to provide a list of capabilities you'd want to have enabled on a single target
    # this means this would need to be completely handwritten
  }
  # allow generation of scoped resources - DONE
  # overwrite method with generate_create = false - DONE
  # scope property name? target_resource_id
  # create method needs to be manually written, what is the order of events - write create method manually first or after?
  # handwrite a chaos_studio_target_resource_gen.go containing bare minimum e.g. resource struct, model struct - then generate
}
# flat additional properties in the swagger parser as a dictionary map[string]interface{}

# allow empty properties payload toggle for resource
#Importer for Service "ChaosStudio".Importer for API Version "2023-04-15-preview":      💥 Error: generating Data API Definitions for Service "ChaosStudio" / API Version "2023-04-15-preview": generating API Versions: generating Resource "Targets" (Namespace "Pandora.Definitions.ResourceManager.ChaosStudio.v2023_04_15_preview.Targets"): generating code for model "Target" in "Pandora.Definitions.ResourceManager.ChaosStudio.v2023_04_15_preview.Targets": generating code for field "Properties": a dictionary must have a reference or a nested item but got neither
#Error: parsing data for Service "ChaosStudio": generating API Definitions for Service "ChaosStudio" / Version "2023-04-15-preview": generating Data API Definitions for Service "ChaosStudio" / API Version "2023-04-15-preview": generating API Versions: generating Resource "Targets" (Namespace "Pandora.Definitions.ResourceManager.ChaosStudio.v2023_04_15_preview.Targets"): generating code for model "Target" in "Pandora.Definitions.ResourceManager.ChaosStudio.v2023_04_15_preview.Targets": generating code for field "Properties": a dictionary must have a reference or a nested item but got neither
#fix remaining errors where when writing into the data API format dictionary string <-> object