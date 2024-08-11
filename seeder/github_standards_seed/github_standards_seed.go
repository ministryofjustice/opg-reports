package github_standards_seed

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/ministryofjustice/opg-reports/datastore/github_standards/ghs"
	"github.com/ministryofjustice/opg-reports/shared/fake"
)

func NewDb(ctx context.Context, dbPath string, schemaPath string) *sql.DB {
	// delete the db
	os.Remove(dbPath)
	os.WriteFile(dbPath, []byte(""), os.ModePerm)

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		slog.Error("error opening db", slog.String("err", err.Error()))
		return nil
	}
	schema, _ := os.ReadFile(schemaPath)
	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		slog.Error("error creating schema", slog.String("err", err.Error()), slog.String("schemaPath", schemaPath))
		return nil
	}
	return db

}

func Seed(ctx context.Context, db *sql.DB, counter int) (q *ghs.Queries) {

	owner := fake.String(12)
	db.Ping()
	q = ghs.New(db)

	for x := 0; x < counter; x++ {
		g := ghs.Fake()
		g.Owner = owner
		g.FullName = fmt.Sprintf("%s/%s", owner, g.Name)
		_, err := q.Insert(ctx, ghs.InsertParams{
			Uuid:                           g.Uuid,
			Ts:                             g.Ts,
			DefaultBranch:                  g.DefaultBranch,
			Owner:                          g.Owner,
			Name:                           g.Name,
			FullName:                       g.FullName,
			License:                        g.License,
			LastCommitDate:                 g.LastCommitDate,
			CreatedAt:                      g.CreatedAt,
			CountOfClones:                  g.CountOfClones,
			CountOfForks:                   g.CountOfForks,
			CountOfPullRequests:            g.CountOfPullRequests,
			CountOfWebHooks:                g.CountOfWebHooks,
			HasCodeOfConduct:               g.HasCodeOfConduct,
			HasCodeownerApprovalRequired:   g.HasCodeownerApprovalRequired,
			HasContributingGuide:           g.HasContributingGuide,
			HasDefaultBranchOfMain:         g.HasDefaultBranchOfMain,
			HasDefaultBranchProtection:     g.HasDefaultBranchProtection,
			HasDeleteBranchOnMerge:         g.HasDeleteBranchOnMerge,
			HasDescription:                 g.HasDescription,
			HasDiscussions:                 g.HasDiscussions,
			HasDownloads:                   g.HasDownloads,
			HasIssues:                      g.HasIssues,
			HasLicense:                     g.HasLicense,
			HasPages:                       g.HasPages,
			HasPullRequestApprovalRequired: g.HasPullRequestApprovalRequired,
			HasReadme:                      g.HasReadme,
			HasRulesEnforcedForAdmins:      g.HasRulesEnforcedForAdmins,
			HasVulnerabilityAlerts:         g.HasVulnerabilityAlerts,
			HasWiki:                        g.HasWiki,
			IsArchived:                     g.IsArchived,
			IsPrivate:                      g.IsPrivate,
			Teams:                          g.Teams,
		})
		if err != nil {
			slog.Error("error creating entry", slog.String("err", err.Error()))
		} else {
			slog.Debug("created entry", slog.Int("x", x), slog.Int("counter", counter))
		}
	}
	return

}
