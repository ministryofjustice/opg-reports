package global

var teams []string = []string{
	"team-a",
	"team-b",
	"team-c",
	"team-d",
	"team-e",
	"team-f",
}

// environments are shared account environment names to use in seeded data
var environments []string = []string{
	"development",
	"pre-production",
	"integrations",
	"production",
}

// fix regions to used for seeds
var regions []string = []string{
	"eu-west-1",
	"eu-west-2",
	"us-east-1",
	"NoRegion",
}

// aws service that we use for seeding
var services []string = []string{
	"Amazon Relational Database Service",
	"Amazon Simple Storage Service",
	"AmazonCloudWatch",
	"Amazon Elastic Load Balancing",
	"AWS Shield",
	"AWS Config",
	"AWS CloudTrail",
	"AWS Key Management Service",
	"Amazon Virtual Private Cloud",
	"Amazon Elastic Container Service",
	"EC2 - Other",
}
