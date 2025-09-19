package awsr

// check interfaces are correct
var (
	_ RepositoryS3           = &Repository{}
	_ RepositorySTS          = &Repository{}
	_ RepositoryCostExplorer = &Repository{}
	_ RespositoryCloudwatch  = &Repository{}
)
