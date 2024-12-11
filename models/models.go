package models

// Full returns the known models that require a database table to be created
func Full() []interface{} {

	return []interface{}{
		&Unit{},                       // Unit is the base grouping model
		&AwsAccount{},                 // AwsAccount details attached to other aws models
		&AwsCost{},                    // AwsCosts model
		&AwsUptime{},                  // AwsUptime tracking
		&Dataset{},                    // Single record table to say if data is real or not
		&GitHubRepositoryGitHubTeam{}, // Many to many join between repo and teams
		&GitHubTeamUnit{},             // Many to many between guthub team and base units
		&GitHubTeam{},                 // GitHub team models used on other github models
		&GitHubRepository{},           // Github repo model
		&GitHubRelease{},              // Release model
		&GitHubRepositoryStandard{},   // Standards model
	}

}
