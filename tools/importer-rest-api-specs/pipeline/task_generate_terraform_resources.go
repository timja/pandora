package pipeline

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/discovery"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/schema"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/transformer"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

func (pipelineTask) generateTerraformDetails(input discovery.ServiceInput, data *models.AzureApiDefinition, logger hclog.Logger) (*models.AzureApiDefinition, error) {
	// TODO: iterate over each of the TF resources that we have in input.TerraformServiceDefinition
	// call the Schema package build that up and the other stuff..

	var apiResource map[string]models.AzureApiResource

	for key, resource := range data.Resources {
		// This is the data API name of the resource i.e. VirtualMachines
		r, err := transformer.ApiResourceFromModelResource(resource)
		if err != nil {
			return nil, err
		}

		b := schema.NewBuilder(resource.Constants, r.Schema.Models, r.Operations.Operations, r.Schema.ResourceIds)

		if t := resource.Terraform; t != nil {
			// We only generate the schema fields for resources that have existing TerraformDetails
			// Seem to only be the ones that already have a .hcl config
			var terraformDetails resourcemanager.TerraformDetails

			for k, v := range t.Resources {
				// This is the Terraform name of the resource i.e. virtual_machine - why does this need to be a map?
				// We need to add to this map any sub-schemas we find so their classes can also be generated
				logger.Info(fmt.Sprintf("Building Schema for %s", k))

				model, err := b.Build(v, logger)
				if err != nil {
					return nil, err
				}

				// Writing all of this info into an empty TerraformDetails struct for this particular resource
				v.SchemaModelName = k
				if v.SchemaModels == nil {
					v.SchemaModels = map[string]resourcemanager.TerraformSchemaModelDefinition{k: *model}
				} else {
					v.SchemaModels[k] = *model
				}

				if terraformDetails.Resources == nil {
					terraformDetails.Resources = map[string]resourcemanager.TerraformResourceDetails{k: v}
				} else {
					terraformDetails.Resources[k] = v
				}
			}

			// Adding the terraformDetails to the relevant resource, keeping the existing info on constants, models etc.
			// This feels unpleasantly hacky
			temp := resource
			temp.Terraform = &terraformDetails

			if apiResource == nil {
				apiResource = map[string]models.AzureApiResource{key: temp}
			} else {
				apiResource[key] = temp
			}
		}
	}

	data.Resources = apiResource

	return data, nil
}
