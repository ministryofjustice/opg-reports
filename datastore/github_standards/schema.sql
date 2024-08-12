CREATE TABLE github_standards (
    -- -- generated ones
    -- auto inc primary key
    id INTEGER PRIMARY KEY,
    -- timestamp for when we generated this record
    ts TEXT NOT NULL,

    -- -- from github data
    -- compliance flags
    compliant_baseline INTEGER NOT NULL DEFAULT 0,
    compliant_extended INTEGER NOT NULL DEFAULT 0,
    -- counter of how many clones this repo has had recently
    count_of_clones INTEGER NOT NULL DEFAULT 0,
    -- count of how forks
    count_of_forks INTEGER NOT NULL DEFAULT 0,
    -- how many open pull requests
    count_of_pull_requests INTEGER NOT NULL DEFAULT 0,
    -- how many webhooks
    count_of_web_hooks INTEGER NOT NULL DEFAULT 0,
    -- the created_at property of the repo
    created_at TEXT NOT NULL,
    -- default branch name, likely main | master
    default_branch TEXT NOT NULL,
    -- <owner>/<name>
    full_name TEXT NOT NULL,
    -- boolean flags
    -- we generate the value of these based on properties of the repository
    -- see commands.github_standards module func mapFromApi
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
    -- is
    is_archived INTEGER NOT NULL DEFAULT 0,
    is_private  INTEGER NOT NULL DEFAULT 0,
    -- what license, likely MIT | GPL
    license TEXT NOT NULL DEFAULT '',
    -- last commit from the default_branch
    last_commit_date TEXT NOT NULL,
    -- repositories name slug
    name TEXT NOT NULL,
    -- owner slug
    owner TEXT NOT NULL,
    -- all teams for this repo
    teams TEXT NOT NULL
) STRICT;
-- always sorted by name, created_at ASC
CREATE INDEX ghs_sort_idx ON github_standards(name, created_at);
-- used by ArchivedFilter to find and sort archived matched rows
CREATE INDEX ghs_archived_idx ON github_standards(is_archived, name, created_at);
-- used by TeamFilter
CREATE INDEX ghs_teams_idx ON github_standards(teams, name, created_at);
-- used ArchivedTeamFilter
CREATE INDEX ghs_archived_teams_idx ON github_standards(is_archived, teams, name, created_at);
--
CREATE INDEX ghs_baseline_idx ON github_standards(compliant_baseline);
CREATE INDEX ghs_extended_idx ON github_standards(compliant_extended);
