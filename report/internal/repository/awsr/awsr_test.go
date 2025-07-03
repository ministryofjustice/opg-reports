package awsr

// check interfaces are correct
var (
	_ S3Repository  = &Repository{}
	_ STSRepository = &Repository{}
)
