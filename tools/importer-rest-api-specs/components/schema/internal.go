package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

type operationPayloads struct {
	createModelName           string
	createPayload             resourcemanager.ModelDetails
	createPropertiesModelName string
	createPropertiesPayload   resourcemanager.ModelDetails

	readModelName           string
	readPayload             resourcemanager.ModelDetails
	readPropertiesModelName string
	readPropertiesPayload   resourcemanager.ModelDetails

	updateModelName           *string
	updatePayload             *resourcemanager.ModelDetails
	updatePropertiesModelName *string
	updatePropertiesPayload   *resourcemanager.ModelDetails
}

func (p operationPayloads) createReadUpdatePayloads() []resourcemanager.ModelDetails {
	out := []resourcemanager.ModelDetails{
		p.createPayload,
		p.readPayload,
	}
	if p.updatePayload != nil {
		out = append(out, *p.updatePayload)
	}
	return out
}

func (p operationPayloads) createReadUpdatePayloadsProperties() []resourcemanager.ModelDetails {
	out := []resourcemanager.ModelDetails{
		p.createPropertiesPayload,
		p.readPropertiesPayload,
	}
	if p.updatePropertiesPayload != nil {
		out = append(out, *p.updatePropertiesPayload)
	}

	return out
}

func (p operationPayloads) getPropertiesModelWithinModel(input resourcemanager.ModelDetails, models map[string]resourcemanager.ModelDetails) (*string, *resourcemanager.ModelDetails) {
	if props, ok := getField(input, "Properties"); ok {
		if props.ObjectDefinition.Type != resourcemanager.ReferenceApiObjectDefinitionType {
			// Chaos 2023-04-15-preview Targets resource return an empty properties object which given what the resource does is acceptable.
			if props.ObjectDefinition.Type == resourcemanager.RawObjectApiObjectDefinitionType {
				return pointer.To("properties"), pointer.To(resourcemanager.ModelDetails{})
			}

			// TODO if this returns nil it might be good to surface an error with instructions e.g. you may want raise swagger bug or consider hand writing the create?
			return nil, nil
		}

		modelName := *props.ObjectDefinition.ReferenceName
		model, ok := models[modelName]
		if !ok {
			return nil, nil
		}

		return &modelName, &model
	}

	return nil, nil
}
