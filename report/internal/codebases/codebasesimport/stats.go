package codebasesimport

import (
	"context"
	"fmt"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/rest"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v81/github"
)

// Raw stats entry
const InsertStatsStatement string = `
INSERT INTO codebase_stats (
	codebase,
	visibility,
	compliance_level,
	compliance_grade,
	compliance_report_url,
	compliance_badge
) VALUES (
	:codebase,
	:visibility,
	:compliance_level,
	:compliance_grade,
	:compliance_report_url,
	:compliance_badge
)
ON CONFLICT (codebase) DO UPDATE SET
	visibility=excluded.visibility,
	compliance_level=excluded.compliance_level,
	compliance_report_url=excluded.compliance_report_url,
	compliance_badge=excluded.compliance_badge,
	compliance_grade=excluded.compliance_grade
RETURNING id
;
`

// leave some space incase of new grades
var gradeMap = map[string]int{
	"unknown":   1,
	"not_found": 10,
	"baseline":  20,
	"standard":  30,
	"exemplar":  40,
}

// Codebase represents a simple, joinless, db row in the cost table; used by imports and seeding commands
type CodebaseStats struct {
	Codebase            string `json:"codebase,omitempty"`              // full name of codebase
	Visibility          string `json:"visibility,omityempty"`           // visibility status
	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url
	ComplianceGrade     int    `json:"compliance_grade,omitempty"`      // numeric version of compliance_level so sorting can be done on this

}

// the badge layout puts the value in the title
var complianceRe = regexp.MustCompile(`(?m)<title>MOJ COMPLIANT:(.*)</title>`)

func handleCodebaseStats(ctx context.Context, repositories []*github.Repository, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "handleCodebaseStats")
	var data []*CodebaseStats = []*CodebaseStats{}
	log.Info("starting codebase stats import ...")
	// convert to local structs
	log.Debug("converting to codebase models ...")
	data, err = toCodebasesStats(ctx, repositories)
	if err != nil {
		return
	}
	// now write to db
	err = dbx.Insert(ctx, InsertStatsStatement, data, &dbx.InsertArgs{
		DB:     in.DB,
		Driver: in.Driver,
		Params: in.Params,
	})
	if err != nil {
		log.Error("error write data during import", "err", err.Error())
		return
	}
	log.With("count", len(data)).Debug("complete.")
	return
}

// toCodebasesStats converts the api results into local structs
func toCodebasesStats(ctx context.Context, list []*github.Repository) (data []*CodebaseStats, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "toModels")
	var base string = "https://github-community.service.justice.gov.uk/repository-standards"

	data = []*CodebaseStats{}
	log.Debug("starting ...")

	for _, item := range list {
		var repo = &CodebaseStats{
			Codebase:            *item.FullName,
			Visibility:          *item.Visibility,
			ComplianceLevel:     "unknown",
			ComplianceGrade:     1,
			ComplianceReportUrl: fmt.Sprintf("%s/%s", base, *item.Name),
			ComplianceBadge:     fmt.Sprintf("%s/api/%s/badge", base, *item.Name),
		}
		// set the compliance level
		if lvl, e := complianceLevelFromBadge(ctx, repo.ComplianceBadge); e == nil {
			repo.ComplianceLevel = lvl
		}
		// set the grade
		repo.ComplianceGrade = gradeMap[repo.ComplianceLevel]
		data = append(data, repo)
	}
	log.Debug("complete.")
	return
}

// complianceLevelFromBadge looks at the badge content (which is svg) and parses the title
// to find the compliance level
func complianceLevelFromBadge(ctx context.Context, badge string) (level string, err error) {
	var timeout = (2 * time.Second)
	level = "unknown"
	res, _, err := rest.GetStr(ctx, nil, &rest.Request{Host: badge, Timeout: timeout})
	if err != nil {
		return
	}
	// find a match
	for _, match := range complianceRe.FindAllString(res, 1) {
		level = match
	}
	// trim the extras
	level = strings.ReplaceAll(level, "<title>MOJ COMPLIANT:", "")
	level = strings.ReplaceAll(level, "</title>", "")
	// swap out not foudn for not_found to make parsing easier
	level = strings.ReplaceAll(level, "NOT FOUND", "not_found")
	level = strings.Trim(level, " ")
	// split on space
	levels := strings.Split(level, " ")
	// use. the last part only
	level = levels[len(levels)-1]
	level = strings.ToLower(level)
	return
}
