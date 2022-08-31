package pipeline

import (
	"fmt"
	"log"
	"strings"

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
			var terraformResourceDetails resourcemanager.TerraformResourceDetails

			nestedModels := make([]string, 0)
			referencedEnums := make([]string, 0)

			for k, v := range t.Resources {
				// This is the Terraform name of the resource i.e. virtual_machine - why does this need to be a map?
				// We need to add to this map any sub-schemas we find so their classes can also be generated
				logger.Info(fmt.Sprintf("Building Schema for %s", k))
				terraformResourceDetails = v
				model, err := b.Build(v, logger)
				if err != nil {
					return nil, err
				}

				nestedModels, referencedEnums = findAllNestedModelsForResource(model, resource)

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

				// Just copy Data Sources for now
				terraformDetails.DataSources = t.DataSources
			}

			terraformDetails.Schemas = map[string]resourcemanager.TerraformResourceDetails{}
			for _, v := range nestedModels {
				nestedModel, err := b.BuildNestedModelDefinition(v, terraformResourceDetails, logger)
				if err != nil || nestedModel == nil {
					continue
				}
				terraformDetails.Schemas[v] = resourcemanager.TerraformResourceDetails{
					SchemaModelName: v,
					SchemaModels:    map[string]resourcemanager.TerraformSchemaModelDefinition{v: *nestedModel},
				}
			}

			terraformDetails.Constants = map[string]resourcemanager.ConstantDetails{}
			for _, v := range referencedEnums {
				if c, ok := resource.Constants[v]; ok {
					terraformDetails.Constants[v] = c
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

	// merge the processed data back in
	for k, v := range apiResource {
		data.Resources[k] = v
	}

	return data, nil
}

func findAllNestedModelsForResource(model *resourcemanager.TerraformSchemaModelDefinition, resource models.AzureApiResource) ([]string, []string) {
	foundModels := make([]string, 0)
	foundEnums := make([]string, 0)

	for _, m := range model.Fields {
		if ref := m.ObjectDefinition.ReferenceName; ref != nil {
			if strings.EqualFold(*ref, "subresource") {
				// subresources are remote IDs, so should be references to the IDs not models here
				continue
			}
			if _, ok := resource.Models[*ref]; ok {
				// Check it's actually a model, not a reference to an enum
				foundModels = append(foundModels, *ref)
				continue
			}
			if _, ok := resource.Constants[*ref]; ok {
				foundEnums = append(foundEnums, *ref)
				continue
			}
		}
		if m.ObjectDefinition.NestedObject != nil {
			if ref := m.ObjectDefinition.NestedObject.ReferenceName; ref != nil {
				if strings.EqualFold(*ref, "subresource") {
					// subresources are remote IDs, so should be references to the IDs not models here
					continue
				}
				if _, ok := resource.Models[*ref]; ok {
					// Check it's actually a model, not a reference to an enum
					foundModels = append(foundModels, *ref)
					continue
				}
				if _, ok := resource.Constants[*ref]; ok {
					foundEnums = append(foundEnums, *ref)
					continue
				}
			}
		}
	}

	// make the list unique
	foundModels = dedupeList(foundModels)
	foundEnums = dedupeList(foundEnums)

	for _, v := range foundModels {
		if m, ok := resource.Models[v]; ok {
			for n, f := range m.Fields {
				log.Printf("processing field %q for model %q", n, v)
				if f.ObjectDefinition != nil && f.ObjectDefinition.ReferenceName != nil {
					ref := *f.ObjectDefinition.ReferenceName
					if _, constant := resource.Constants[ref]; constant {
						if _, isConstant := resource.Constants[ref]; isConstant {
							foundEnums = append(foundEnums, ref)
						}
						continue
					}
					foundModels, foundEnums = findNestedModelsForModel(ref, resource, foundModels, foundEnums)
				} else if f.ObjectDefinition != nil && f.ObjectDefinition.NestedItem != nil && f.ObjectDefinition.NestedItem.ReferenceName != nil {
					ref := *f.ObjectDefinition.NestedItem.ReferenceName
					if _, constant := resource.Constants[ref]; constant {
						if _, isConstant := resource.Constants[ref]; isConstant {
							foundEnums = append(foundEnums, ref)
						}
						continue
					}
					foundModels, foundEnums = findNestedModelsForModel(ref, resource, foundModels, foundEnums)
				}
			}
		}
	}

	return foundModels, foundEnums
}

func findNestedModelsForModel(model string, resource models.AzureApiResource, foundModels []string, foundEnums []string) ([]string, []string) {
	newRef := ""
	if m, ok := resource.Models[model]; ok {
		foundModels = append(foundModels, model)
		for _, v := range m.Fields {
			if v.ObjectDefinition != nil && v.ObjectDefinition.ReferenceName != nil {
				newRef = *v.ObjectDefinition.ReferenceName
				if _, ok := resource.Models[newRef]; !ok {
					if _, isConstant := resource.Constants[newRef]; isConstant {
						foundEnums = append(foundEnums, newRef)
					}
					continue
				}
				foundModels = append(foundModels, newRef)
			} else if v.ObjectDefinition.NestedItem != nil && v.ObjectDefinition.NestedItem.ReferenceName != nil {
				newRef = *v.ObjectDefinition.NestedItem.ReferenceName
				if _, ok := resource.Models[newRef]; !ok {
					if _, isConstant := resource.Constants[newRef]; isConstant {
						foundEnums = append(foundEnums, newRef)
					}
					continue
				}
				foundModels = append(foundModels, newRef)
			}
			if newRef != "" {
				foundModels, foundEnums = findNestedModelsForModel(newRef, resource, foundModels, foundEnums)
			}
		}
	}

	return dedupeList(foundModels), dedupeList(foundEnums)
}

func dedupeList(input []string) []string {
	uniqMap := make(map[string]struct{})
	for _, v := range input {
		uniqMap[v] = struct{}{}
	}

	uniqList := make([]string, 0, len(uniqMap))
	for v := range uniqMap {
		uniqList = append(uniqList, v)
	}

	return uniqList
}
