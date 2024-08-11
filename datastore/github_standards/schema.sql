CREATE TABLE github_standards (
    -- uid
    uuid                VARCHAR(36) PRIMARY KEY,
    ts                  VARCHAR(50) NOT NULL,
    default_branch      VARCHAR(50) NOT NULL,
    full_name           VARCHAR(255) NOT NULL,
    name                VARCHAR(255) NOT NULL,
    owner               VARCHAR(255) NOT NULL,
    license             VARCHAR(50) NOT NULL,
    last_commit_date    VARCHAR(30) NOT NULL,
    created_at          VARCHAR(30) NOT NULL,

    count_of_clones         INT(5) NOT NULL DEFAUlt 0,
    count_of_forks          INT(5) NOT NULL DEFAUlt 0,
    count_of_pull_requests  INT(5) NOT NULL DEFAUlt 0,
    count_of_web_hooks      INT(5) NOT NULL DEFAUlt 0,

    -- boolean flags for property checks
    has_code_of_conduct                 INT(2) NOT NULL  DEFAUlt 0,
    has_codeowner_approval_required     INT(2) NOT NULL  DEFAUlt 0,
    has_contributing_guide              INT(2) NOT NULL  DEFAUlt 0,
    has_default_branch_of_main          INT(2) NOT NULL  DEFAUlt 0,
    has_default_branch_protection       INT(2) NOT NULL  DEFAUlt 0,
    has_delete_branch_on_merge          INT(2) NOT NULL  DEFAUlt 0,
    has_description                     INT(2) NOT NULL  DEFAUlt 0,
    has_discussions                     INT(2) NOT NULL  DEFAUlt 0,
    has_downloads                       INT(2) NOT NULL  DEFAUlt 0,
    has_issues                          INT(2) NOT NULL  DEFAUlt 0,
    has_license                         INT(2) NOT NULL  DEFAUlt 0,
    has_pages                           INT(2) NOT NULL  DEFAUlt 0,
    has_pull_request_approval_required  INT(2) NOT NULL  DEFAUlt 0,
    has_readme                          INT(2) NOT NULL  DEFAUlt 0,
    has_rules_enforced_for_admins       INT(2) NOT NULL  DEFAUlt 0,
    has_vulnerability_alerts            INT(2) NOT NULL  DEFAUlt 0,
    has_wiki                            INT(2) NOT NULL  DEFAUlt 0,

    is_archived INT(2) NOT NULL  DEFAUlt 0,
    is_private  INT(2) NOT NULL  DEFAUlt 0,
    -- all teams for this repo
    teams TEXT NOT NULL

);
-- always sorted by name, created_at ASC
CREATE INDEX ghs_sort_idx ON github_standards(name, created_at);
-- used by ArchivedFilter to find and sort archived matched rows
CREATE INDEX ghs_archived_idx ON github_standards(is_archived, name, created_at);
-- used by TeamFilter
CREATE INDEX ghs_teams_idx ON github_standards(teams, name, created_at);
--  used ArchivedTeamFilter
CREATE INDEX ghs_archived_teams_idx ON github_standards(is_archived, teams, name, created_at);
