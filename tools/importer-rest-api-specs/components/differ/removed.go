package differ

import (
	"fmt"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
)

func (d Differ) RemovedAttributesFromExistingAPIDefinitions(existing models.AzureApiDefinition, parsed models.AzureApiDefinition) (DiffedAzureApiDefinition, error) {
	output := DiffedAzureApiDefinition{
		Resources: make(map[string]DiffedAzureApiResource, 0),
	}
	for resourceName, resource := range existing.Resources {
		newResource, ok := parsed.Resources[resourceName]
		if !ok {
			return output, fmt.Errorf("unable to find %q resource in new api definition", resourceName)
		}

		diffedModels := DiffedAzureApiResource{
			Models: make(map[string][]string, 0),
		}
		for oldModelName, oldModel := range resource.Models {
			newModel, ok := newResource.Models[oldModelName]
			if !ok {
				return output, fmt.Errorf("unable to find %q in %q new model definition", resourceName, resourceName)
			}

			removedFields := make([]string, 0)
			for fieldName := range oldModel.Fields {

				found := false

				for newFieldName := range newModel.Fields {
					if fieldName == newFieldName {
						found = true
						break
					}
				}

				if !found {
					removedFields = append(removedFields, fieldName)
				}
			}
			if len(removedFields) > 0 {
				diffedModels.Models[oldModelName] = removedFields
			}
		}
		if len(diffedModels.Models) > 0 {
			output.Resources[resourceName] = diffedModels
		}
	}

	return output, nil
}
