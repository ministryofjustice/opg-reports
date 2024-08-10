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

    count_of_clones         INT(5) NOT NULL,
    count_of_forks          INT(5) NOT NULL,
    count_of_pull_requests  INT(5) NOT NULL,
    count_of_web_hooks      INT(5) NOT NULL,

    has_code_of_conduct                 INT(2) NOT NULL,
    has_codeowner_approval_required     INT(2) NOT NULL,
    has_contributing_guide              INT(2) NOT NULL,
    has_default_branch_of_main          INT(2) NOT NULL,
    has_default_branch_protection       INT(2) NOT NULL,
    has_delete_branch_on_merge          INT(2) NOT NULL,
    has_description                     INT(2) NOT NULL,
    has_discussions                     INT(2) NOT NULL,
    has_downloads                       INT(2) NOT NULL,
    has_issues                          INT(2) NOT NULL,
    has_license                         INT(2) NOT NULL,
    has_pages                           INT(2) NOT NULL,
    has_pull_request_approval_required  INT(2) NOT NULL,
    has_readme                          INT(2) NOT NULL,
    has_rules_enforced_for_admins       INT(2) NOT NULL,
    has_vulnerability_alerts            INT(2) NOT NULL,
    has_wiki                            INT(2) NOT NULL,

    is_archived INT(2) NOT NULL,
    is_private  INT(2) NOT NULL


);
