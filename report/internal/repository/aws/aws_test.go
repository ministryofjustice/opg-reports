package aws

// check interfaces are correct
var (
	_ S3er  = &Repository{}
	_ STSer = &Repository{}
)
