package differ

type DiffedAzureApiDefinition struct {
	Resources map[string]DiffedAzureApiResource
}

type DiffedAzureApiResource struct {
	Models map[string]DiffedModelDetails
}

type DiffedModelDetails struct {
	RemovedFields []string
}
