package pipeline

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/discovery"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/schema"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/components/transformer"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
	"log"
)

func (pipelineTask) generateTerraformDetails(input discovery.ServiceInput, data *models.AzureApiDefinition, logger hclog.Logger) (*models.AzureApiDefinition, error) {
	// TODO: iterate over each of the TF resources that we have in input.TerraformServiceDefinition
	// call the Schema package build that up and the other stuff..

	parsedResources := data.Resources

	for key, resource := range data.Resources {
		// This is the data API name of the resource i.e. VirtualMachines
		log.Printf("[DEBUG] Resource 1 %s", key)
		r, err := transformer.ApiResourceFromModelResource(resource)
		if err != nil {
			return nil, err
		}

		b := schema.NewBuilder(resource.Constants, r.Schema.Models, r.Operations.Operations, r.Schema.ResourceIds)

		if t := resource.Terraform; t != nil {
			terraformDetails := t
			//log.Printf("[DEBUG] Terraform details for %s: %+v", key, t)
			schemaModels := make(map[string]resourcemanager.TerraformSchemaModelDefinition, 0)
			//resourceDetails := make(map[string]resourcemanager.TerraformResourceDetails, 0)

			for k, v := range t.Resources {
				// This is the Terraform name of the resource i.e. virtual_machine - why does this need to be a map?
				log.Printf("[DEBUG] Resource 2 %s", k)
				logger.Info(fmt.Sprintf("Building Schema for %s", k))
				model, err := b.Build(v, logger)
				if err != nil {
					return nil, err
				}
				parsedResources[key].Terraform.Resources[k].SchemaModels = map[string]resourcemanager.TerraformSchemaModelDefinition{k: *model}
				schemaModels[k] = *model
				//resourceDetails[k].SchemaModels = schemaModels
				v.SchemaModels = map[string]resourcemanager.TerraformSchemaModelDefinition{k: *model}

				//log.Printf("[DEBUG] Resource details for %s: %+v", k, v)
			}
			//t.Resources["how tf do I get the key?"]
		}
	}

	log.Printf("[DEBUG] model: %+v", data.Resources["VirtualMachineScaleSets"].Terraform)
	return data, nil
}
