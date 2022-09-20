package differ

type DiffedAzureApiDefinition struct {
	Resources map[string]DiffedAzureApiResource
}

// todo change []string to a struct so we can hold more info if we want to
type DiffedAzureApiResource struct {
	Models map[string][]string
}
