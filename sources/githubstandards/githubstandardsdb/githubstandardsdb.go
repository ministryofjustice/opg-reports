package githubstandardsdb

import "github.com/ministryofjustice/opg-reports/pkg/datastore"

const CreateStandardsTable datastore.CreateStatement = `
CREATE TABLE standards (
    id INTEGER PRIMARY KEY,
    ts TEXT NOT NULL,
    compliant_baseline INTEGER NOT NULL DEFAULT 0,
    compliant_extended INTEGER NOT NULL DEFAULT 0,
    count_of_clones INTEGER NOT NULL DEFAULT 0,
    count_of_forks INTEGER NOT NULL DEFAULT 0,
    count_of_pull_requests INTEGER NOT NULL DEFAULT 0,
    count_of_web_hooks INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    default_branch TEXT NOT NULL,
    full_name TEXT NOT NULL,
    has_code_of_conduct                 INTEGER NOT NULL DEFAULT 0,
    has_codeowner_approval_required     INTEGER NOT NULL DEFAULT 0,
    has_contributing_guide              INTEGER NOT NULL DEFAULT 0,
    has_default_branch_of_main          INTEGER NOT NULL DEFAULT 0,
    has_default_branch_protection       INTEGER NOT NULL DEFAULT 0,
    has_delete_branch_on_merge          INTEGER NOT NULL DEFAULT 0,
    has_description                     INTEGER NOT NULL DEFAULT 0,
    has_discussions                     INTEGER NOT NULL DEFAULT 0,
    has_downloads                       INTEGER NOT NULL DEFAULT 0,
    has_issues                          INTEGER NOT NULL DEFAULT 0,
    has_license                         INTEGER NOT NULL DEFAULT 0,
    has_pages                           INTEGER NOT NULL DEFAULT 0,
    has_pull_request_approval_required  INTEGER NOT NULL DEFAULT 0,
    has_readme                          INTEGER NOT NULL DEFAULT 0,
    has_rules_enforced_for_admins       INTEGER NOT NULL DEFAULT 0,
    has_vulnerability_alerts            INTEGER NOT NULL DEFAULT 0,
    has_wiki                            INTEGER NOT NULL DEFAULT 0,
    is_archived INTEGER NOT NULL DEFAULT 0,
    is_private  INTEGER NOT NULL DEFAULT 0,
    license TEXT NOT NULL DEFAULT '',
    last_commit_date TEXT NOT NULL,
    name TEXT NOT NULL,
    owner TEXT NOT NULL,
    teams TEXT NOT NULL
) STRICT
;`

const CreateStandardsIndexIsArchived datastore.CreateStatement = `CREATE INDEX stnd_archived_idx ON standards(is_archived, name, created_at);`
const CreateStandardsIndexIsArchivedTeams datastore.CreateStatement = `CREATE INDEX stnd_archived_teams_idx ON standards(is_archived, teams, name, created_at);`
const CreateStandardsIndexTeams datastore.CreateStatement = `CREATE INDEX stnd_teams_idx ON standards(teams, name, created_at);`
const CreateStandardsIndexBaseline datastore.CreateStatement = `CREATE INDEX stnd_baseline_idx ON standards(compliant_baseline)`
const CreateStandardsIndexExtended datastore.CreateStatement = `CREATE INDEX stnd_extended_idx ON standards(compliant_extended)`

const InsertStandard datastore.InsertStatement = `
INSERT INTO standards(
    ts,
    default_branch,
    full_name,
    name,
    owner,
    license,
    last_commit_date,
    created_at,
    count_of_clones,
    count_of_forks,
    count_of_pull_requests,
    count_of_web_hooks,
    has_code_of_conduct,
    has_codeowner_approval_required,
    has_contributing_guide,
    has_default_branch_of_main,
    has_default_branch_protection,
    has_delete_branch_on_merge,
    has_description,
    has_discussions,
    has_downloads,
    has_issues,
    has_license,
    has_pages,
    has_pull_request_approval_required,
    has_readme,
    has_rules_enforced_for_admins,
    has_vulnerability_alerts,
    has_wiki,
    is_archived,
    is_private,
    compliant_baseline,
    compliant_extended,
    teams
) VALUES (
	:ts,
	:default_branch,
	:full_name,
	:name,
	:owner,
	:license,
	:last_commit_date,
	:created_at,
	:count_of_clones,
	:count_of_forks,
	:count_of_pull_requests,
	:count_of_web_hooks,
	:has_code_of_conduct,
	:has_codeowner_approval_required,
	:has_contributing_guide,
	:has_default_branch_of_main,
	:has_default_branch_protection,
	:has_delete_branch_on_merge,
	:has_description,
	:has_discussions,
	:has_downloads,
	:has_issues,
	:has_license,
	:has_pages,
	:has_pull_request_approval_required,
	:has_readme,
	:has_rules_enforced_for_admins,
	:has_vulnerability_alerts,
	:has_wiki,
	:is_archived,
	:is_private,
	:compliant_baseline,
	:compliant_extended,
	:teams
) RETURNING id;`

// RowCount returns the total number of records within the database
const RowCount datastore.SelectStatement = `
SELECT
	count(*) as row_count
FROM standards
LIMIT 1
;`

// ArchivedCount counts number of repos that are marked as archived
const ArchivedCount datastore.SelectStatement = `
SELECT
	count(*) as row_count
FROM standards
WHERE
	is_archived=1
LIMIT 1
;`

// CompliantBaselineCount counts number of records with baseline compliance
const CompliantBaselineCount datastore.SelectStatement = `
SELECT
	count(*) as row_count
FROM standards
WHERE
	compliant_baseline=1
LIMIT 1
;`

// CompliantBaselineCount counts number of records with extended compliance
const CompliantExtendedCount datastore.SelectStatement = `
SELECT
	count(*) as row_count
FROM standards
WHERE
	compliant_extended=1
LIMIT 1
;`

// All returns all items from the db
const All datastore.SelectStatement = `
SELECT
	*
FROM standards
ORDER BY name, created_at ASC
;`

const FilterByIsArchived datastore.NamedSelectStatement = `
SELECT
	*
FROM standards
WHERE
	is_archived = :archived
ORDER BY name, created_at ASC
;`

const FilterByTeam datastore.NamedSelectStatement = `
SELECT
	*
FROM standards
WHERE
	teams LIKE :team_string
ORDER BY name, created_at ASC
;`

const FilterByIsArchivedAndTeam datastore.NamedSelectStatement = `
SELECT
	*
FROM standards
WHERE
	is_archived = :archived
	teams LIKE :team_string
ORDER BY name, created_at ASC
;`

const Age datastore.SelectStatement = `
SELECT
	MIN(ts) as age
FROM standards
LIMIT 1
;`
