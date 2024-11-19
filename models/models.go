package models

// All returns a list of all models
func All() []interface{} {

	return []interface{}{
		&Unit{},             // Unit is the base grouping model
		&AwsAccount{},       // AwsAccount details attached to other aws models
		&AwsCost{},          // AwsCosts model
		&AwsUptime{},        // AwsUptime tracking
		&GitHubTeam{},       // GitHub team models used on other github models
		&GitHubRepository{}, // Github repo model
		&GitHubRelease{},
		&GitHubRepositoryStandard{},
	}

}
