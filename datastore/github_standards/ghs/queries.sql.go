// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package ghs

import (
	"context"
)

const all = `-- name: All :many
SELECT uuid, ts, default_branch, full_name, name, owner, license, last_commit_date, created_at, count_of_clones, count_of_forks, count_of_pull_requests, count_of_web_hooks, has_code_of_conduct, has_codeowner_approval_required, has_contributing_guide, has_default_branch_of_main, has_default_branch_protection, has_delete_branch_on_merge, has_description, has_discussions, has_downloads, has_issues, has_license, has_pages, has_pull_request_approval_required, has_readme, has_rules_enforced_for_admins, has_vulnerability_alerts, has_wiki, is_archived, is_private FROM github_standards
ORDER BY name, created_at ASC
`

func (q *Queries) All(ctx context.Context) ([]GithubStandard, error) {
	rows, err := q.db.QueryContext(ctx, all)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GithubStandard
	for rows.Next() {
		var i GithubStandard
		if err := rows.Scan(
			&i.Uuid,
			&i.Ts,
			&i.DefaultBranch,
			&i.FullName,
			&i.Name,
			&i.Owner,
			&i.License,
			&i.LastCommitDate,
			&i.CreatedAt,
			&i.CountOfClones,
			&i.CountOfForks,
			&i.CountOfPullRequests,
			&i.CountOfWebHooks,
			&i.HasCodeOfConduct,
			&i.HasCodeownerApprovalRequired,
			&i.HasContributingGuide,
			&i.HasDefaultBranchOfMain,
			&i.HasDefaultBranchProtection,
			&i.HasDeleteBranchOnMerge,
			&i.HasDescription,
			&i.HasDiscussions,
			&i.HasDownloads,
			&i.HasIssues,
			&i.HasLicense,
			&i.HasPages,
			&i.HasPullRequestApprovalRequired,
			&i.HasReadme,
			&i.HasRulesEnforcedForAdmins,
			&i.HasVulnerabilityAlerts,
			&i.HasWiki,
			&i.IsArchived,
			&i.IsPrivate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const archived = `-- name: Archived :many
SELECT uuid, ts, default_branch, full_name, name, owner, license, last_commit_date, created_at, count_of_clones, count_of_forks, count_of_pull_requests, count_of_web_hooks, has_code_of_conduct, has_codeowner_approval_required, has_contributing_guide, has_default_branch_of_main, has_default_branch_protection, has_delete_branch_on_merge, has_description, has_discussions, has_downloads, has_issues, has_license, has_pages, has_pull_request_approval_required, has_readme, has_rules_enforced_for_admins, has_vulnerability_alerts, has_wiki, is_archived, is_private FROM github_standards
WHERE is_archived = ?
ORDER BY name, created_at ASC
`

func (q *Queries) Archived(ctx context.Context, isArchived int) ([]GithubStandard, error) {
	rows, err := q.db.QueryContext(ctx, archived, isArchived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GithubStandard
	for rows.Next() {
		var i GithubStandard
		if err := rows.Scan(
			&i.Uuid,
			&i.Ts,
			&i.DefaultBranch,
			&i.FullName,
			&i.Name,
			&i.Owner,
			&i.License,
			&i.LastCommitDate,
			&i.CreatedAt,
			&i.CountOfClones,
			&i.CountOfForks,
			&i.CountOfPullRequests,
			&i.CountOfWebHooks,
			&i.HasCodeOfConduct,
			&i.HasCodeownerApprovalRequired,
			&i.HasContributingGuide,
			&i.HasDefaultBranchOfMain,
			&i.HasDefaultBranchProtection,
			&i.HasDeleteBranchOnMerge,
			&i.HasDescription,
			&i.HasDiscussions,
			&i.HasDownloads,
			&i.HasIssues,
			&i.HasLicense,
			&i.HasPages,
			&i.HasPullRequestApprovalRequired,
			&i.HasReadme,
			&i.HasRulesEnforcedForAdmins,
			&i.HasVulnerabilityAlerts,
			&i.HasWiki,
			&i.IsArchived,
			&i.IsPrivate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
