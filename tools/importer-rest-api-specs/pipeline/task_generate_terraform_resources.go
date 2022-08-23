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

			nestedModels := make([]string, 0)

			for k, v := range t.Resources {
				// This is the Terraform name of the resource i.e. virtual_machine - why does this need to be a map?
				// We need to add to this map any sub-schemas we find so their classes can also be generated
				logger.Info(fmt.Sprintf("Building Schema for %s", k))

				model, err := b.Build(v, logger)
				if err != nil {
					return nil, err
				}

				nestedModels = findAllNestedModelsForResource(model, resource)

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

			terraformDetails.Schemas = map[string]resourcemanager.TerraformResourceDetails{}
			for _, v := range nestedModels {
				nestedModel, err := b.BuildNestedModelDefinition(v, logger)
				if err != nil {
					continue
				}
				terraformDetails.Schemas[v] = resourcemanager.TerraformResourceDetails{
					SchemaModelName: v,
					SchemaModels:    map[string]resourcemanager.TerraformSchemaModelDefinition{v: *nestedModel},
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

func findAllNestedModelsForResource(model *resourcemanager.TerraformSchemaModelDefinition, resource models.AzureApiResource) []string {
	temp := make([]string, 0)

	for _, m := range model.Fields {
		if m.ObjectDefinition.ReferenceName != nil {
			temp = append(temp, *m.ObjectDefinition.ReferenceName)
		}
	}

	// This needs to be more clever - sub-nested models have a hierarchy to walk to get to
	//for _, i := range resource.Models {
	//	for _, v := range i.Fields {
	//		if obj := v.ObjectDefinition; obj != nil && obj.ReferenceName != nil {
	//			temp = append(temp, *obj.ReferenceName)
	//		}
	//	}
	//}

	uniqMap := make(map[string]struct{})
	for _, v := range temp {
		uniqMap[v] = struct{}{}
	}

	nestedModels := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		nestedModels = append(nestedModels, v)
	}

	return nestedModels
}

func findNestedModelsForModel(input models.ModelDetails) []string {
	result := make([]string, 0)
	for _, v := range input.Fields {
		if obj := v.ObjectDefinition; obj != nil && obj.ReferenceName != nil {
			result = append(result, *obj.ReferenceName)
		}
	}
	return result
}
