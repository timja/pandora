package differ

import (
	"fmt"
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
)

func (d Differ) RemovedAttributeesFromExistingAPIDefinitions(existing models.AzureApiDefinition, parsed models.AzureApiDefinition) (models.AzureApiDefinition, error) {
	var output models.AzureApiDefinition
	for resourceName, resource := range existing.Resources {
		newResource, ok := parsed.Resources[resourceName]
		if !ok {
			return output, fmt.Errorf("unable to find %q resource in new api definition", resourceName)
		}

		for oldModelName, oldModel := range resource.Models {
			diffedModels := make(map[string]DiffedModelDetails)
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
				diffedModels[oldModelName] = removedFields
			}
		}
	}

	return output, nil
}
