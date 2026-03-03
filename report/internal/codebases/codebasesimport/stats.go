package codebasesimport

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"opg-reports/report/package/cntxt"
	"opg-reports/report/package/dbx"
	"opg-reports/report/package/files"
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
	compliance_badge,
	trivy_usage,
	trivy_sbom_usage,
	trivy_locations
) VALUES (
	:codebase,
	:visibility,
	:compliance_level,
	:compliance_grade,
	:compliance_report_url,
	:compliance_badge,
	:trivy_usage,
	:trivy_sbom_usage,
	:trivy_locations
)
ON CONFLICT (codebase) DO UPDATE SET
	compliance_level=excluded.compliance_level,
	compliance_report_url=excluded.compliance_report_url,
	compliance_badge=excluded.compliance_badge,
	compliance_grade=excluded.compliance_grade,
	trivy_usage=excluded.trivy_usage,
	trivy_sbom_usage=excluded.trivy_sbom_usage,
	trivy_locations=excluded.trivy_locations,
	visibility=excluded.visibility
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
	Codebase   string `json:"codebase,omitempty"`    // full name of codebase
	Visibility string `json:"visibility,omityempty"` // visibility status

	ComplianceLevel     string `json:"compliance_level,omitempty"`      // compliance level (moj based)
	ComplianceReportUrl string `json:"compliance_report_url,omitempty"` // compliance report url
	ComplianceBadge     string `json:"compliance_badge,omitempty"`      // compliance badge url
	ComplianceGrade     int    `json:"compliance_grade,omitempty"`      // numeric version of compliance_level so sorting can be done on this

	TrivyUsage     int    `json:"trivy_usage"`      // boolean flag to show if the codebase is using trivy in workflows
	TrivySBOMUsage int    `json:"trivy_sbom_usage"` // boolean flag to show if trivy is being used to generate sboms
	TrivyLocations string `json:"trivy_locations"`  // tracks files where trivy has been utilised
}

// the badge layout puts the value in the title
var complianceRe = regexp.MustCompile(`(?m)<title>MOJ COMPLIANT:(.*)</title>`)

func handleCodebaseStats(ctx context.Context, client RepoClient, repositories []*github.Repository, in *Args) (err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "handleCodebaseStats")
	var data []*CodebaseStats = []*CodebaseStats{}
	log.Info("starting codebase stats import ...")
	// convert to local structs
	log.Debug("converting to codebase models ...")
	data, err = toCodebasesStats(ctx, client, repositories)
	if err != nil {
		return
	}
	// dump.Now(data)
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
func toCodebasesStats(ctx context.Context, client RepoClient, list []*github.Repository) (data []*CodebaseStats, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "toCodebasesStats")

	data = []*CodebaseStats{}
	log.Debug("starting ...")

	for _, repo := range list {
		log.Debug("fetching data for ... ", "repository", *repo.Name)
		var stats = &CodebaseStats{
			Codebase:   *repo.FullName,
			Visibility: *repo.Visibility,
		}
		// set compliance data
		err = setComplianceData(ctx, client, repo, stats)
		if err != nil {
			log.Error("error getting compliance data", "err", err.Error())
			return
		}

		// set trivy values
		err = setTrivyData(ctx, client, repo, stats)
		if err != nil {
			log.Error("error getting trivy data", "err", err.Error())
			return
		}

		data = append(data, stats)
	}
	log.Debug("complete.")
	return
}

// setComplianceData takes care of handling compliance stats about the code base.
//
// Sets default values for compliance data and then calls & processes the moj compliance badge to determine the
// level the codebases is at
func setComplianceData(ctx context.Context, client RepoClient, repo *github.Repository, stats *CodebaseStats) (err error) {
	var (
		log  *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "setComplianceData", "repo", *repo.Name)
		base string       = "https://github-community.service.justice.gov.uk/repository-standards"
		lvl  string       = "unknown"
	)
	log.Debug("starting ...")
	if *repo.Archived {
		log.Warn("repository is archived, skipping fetching compliance details.")
		return
	}
	// set default values
	stats.ComplianceLevel = lvl
	stats.ComplianceGrade = 1
	// set report & badge url
	stats.ComplianceReportUrl = fmt.Sprintf("%s/%s", base, *repo.Name)
	stats.ComplianceBadge = fmt.Sprintf("%s/api/%s/badge", base, *repo.Name)
	// parse the compliance badge output and set value
	lvl, err = complianceLevelFromBadge(ctx, stats.ComplianceBadge)
	if err != nil {
		return
	}
	stats.ComplianceLevel = lvl
	stats.ComplianceGrade = gradeMap[stats.ComplianceLevel]

	log.With("stats", stats).Debug("complete.")
	return
}

// complianceLevelFromBadge looks at the badge content (which is svg) and parses the title
// to find the compliance level
func complianceLevelFromBadge(ctx context.Context, badge string) (level string, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "complianceLevelFromBadge")
	var timeout = (2 * time.Second)
	level = "unknown"
	res, _, err := rest.GetStr(ctx, nil, &rest.Request{Host: badge, Timeout: timeout})
	// if its a timeout, then we wont throw and error and just carry on
	if err != nil && strings.Contains(err.Error(), "Client.Timeout exceeded") {
		log.Warn("timeout when fetching badge", "badge", badge)
		return level, nil
	}
	if err != nil {
		log.Error("error when fetching badge", "badge", badge, "err", err.Error())
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

// setTrivyData gets all files from `.github/` folder and looks at each for usage of trivy (via cli or action) and
// configures the flags accordingly.
//
// Starts in the `./.github/` directory path and recursively calls `GetContents` finding all files and returning only
// yaml / yml extensions.
//
// Looks for `trivy ` or the trivy action `aquasecurity/trivy-action` in each line of the file to decide it trivy is
// being used.
//
// If the action version is found, then it gets teh action definition (via chunking of lines) and looks for the `output:`
// property in the action, then looks for `.sbom.` being used in the output value - which is what triggers sbom generation
// for the aciton
//
// For cli usage it looks for `trivy sbom` in the line
func setTrivyData(ctx context.Context, client RepoClient, repo *github.Repository, stats *CodebaseStats) (err error) {
	var (
		log       *slog.Logger                = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "setTrivyData", "repo", *repo.Name)
		contents  []*github.RepositoryContent = []*github.RepositoryContent{}
		dir       string                      = "./.github/"
		locations string                      = ""
	)
	log.Debug("starting ...")
	// set the defaults
	stats.TrivyLocations = ""
	stats.TrivyUsage = 0
	stats.TrivySBOMUsage = 0
	//
	if *repo.Archived {
		log.Warn("repository is archived, skipping fetching trivy details.")
		return
	}
	// loop over all files and find trivy usage & sbom generation
	contents, err = getAllGithubFiles(ctx, client, repo, dir)
	if err != nil {
		return
	}
	// parse all `.github` files looking for trivy
	for _, file := range contents {
		var trivy = false
		var sbom = false
		// look for trivy data
		trivy, sbom, err = findTrivyInFile(ctx, client, repo, file)
		if err != nil {
			log.Error("error checking for trivy", "err", err.Error())
			return
		}
		// if either trivy or sbom are true, then track the file path
		if trivy || sbom {
			locations += fmt.Sprintf("%s,", *file.Path)
			stats.TrivyUsage = 1
		}
		if sbom {
			stats.TrivySBOMUsage = 1
		}
	}
	stats.TrivyLocations = locations

	log.With("stats", stats).Debug("complete.")
	return
}

// findTrivyInFile looks for trivy usage within file content by downloading it and scanning lines
func findTrivyInFile(ctx context.Context, client RepoClient, repo *github.Repository, file *github.RepositoryContent) (hasTrivy bool, hasSBOM bool, err error) {
	var (
		buff  io.ReadCloser
		lines []string
		log   *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "findTrivyInFile", "repo", *repo.Name, "file", *file.Path)
	)
	log.Debug("starting ...")
	// download the file content
	buff, _, err = client.DownloadContents(ctx, *repo.Owner.Login, *repo.Name, *file.Path, nil)
	if err != nil {
		log.Error("error downloading content", "err", err.Error())
		return
	}
	// split into lines to make searching easier
	lines = files.Lines(buff)
	// now look for trivy in the file lines
	hasTrivy, hasSBOM = checkFileLinesForTrivy(ctx, lines)

	log.Debug("complete.")
	return
}

// checkFileLinesForTrivy processes the workflow / action yanl files line by line, looking for trivy usage
// via the cli or an action.
//
// Will also check for sbom usage as well
func checkFileLinesForTrivy(ctx context.Context, lines []string) (trivy bool, sbom bool) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "checkFileLines")
	log.Debug("starting ...")
	trivy = false
	sbom = false

	// look over each line....
	//
	// - look for the cli usage first
	for i, line := range lines {
		// if the cli version is found
		if isTrivyCLIInLine(line) {
			trivy = true
			// look for sbom
			if strings.Contains(line, " sbom") {
				sbom = true
			}
		} else if isTrivyActionInLine(line) {
			trivy = true
			if l := actionHasSBOM(lines[i+1:]); l >= 0 {
				sbom = true
			}
		}

	}
	log.Debug("complete.")
	return
}

// actionHasSBOM looks for the output setting in the action config
//
// It does this by finding the end of the current action in the file and
// then scanning that chunk for the `output` attribute which should have
// a `${file}.sbom.${ext}` file name for sbom usage
func actionHasSBOM(remaining []string) (found int) {
	var nextPos = -1
	var chunk = remaining

	found = -1
	// find the next action and generate a chunk of content lines
	// to then look for hte output setting
	for i, line := range remaining {
		line = strings.TrimSpace(line)
		// fmt.Println(">>[" + line + "]<<")
		if len(line) > 0 && line[0] == '-' {
			nextPos = i
			break
		}
	}
	if nextPos >= 0 {
		chunk = remaining[0:nextPos]
	}
	for i, line := range chunk {
		line = strings.TrimSpace(line)
		// only look for `output` as there might be a space around the `:` (` : ` etc)
		if len(line) >= 6 && line[0:6] == "output" && strings.Contains(line, ".sbom.") {
			found = i
		}
	}

	return
}

// isTrivyActionInLine looks for the action version of trivy within a
// single line of a file
func isTrivyActionInLine(line string) (found bool) {
	var (
		trivyAction string = "aquasecurity/trivy-action"
		idx         int    = strings.Index(line, trivyAction)
	)
	found = false
	// if we found it, check the string before to make sure its not commentted out
	if idx >= 0 {
		found = !strings.Contains(line[0:idx], "#")
	}

	return
}

// isTrivyCLIInLine looks for cli version of trivy
func isTrivyCLIInLine(line string) (found bool) {
	var (
		trivyCLI string = "trivy "
		idx      int    = strings.Index(line, trivyCLI)
	)
	found = false
	// if we found it, check the string before to make sure its not commentted out
	if idx >= 0 {
		found = !strings.Contains(line[0:idx], "#")
	}
	return
}

// getAllGithubFiles recurisvely finds all files within a starting directory path - used to fetch all the `.github` sub files
// so we can then check content of those for trivy etc
func getAllGithubFiles(ctx context.Context, client RepoClient, repo *github.Repository, dir string) (contents []*github.RepositoryContent, err error) {
	var log *slog.Logger = cntxt.GetLogger(ctx).With("package", "codebasesimport", "func", "getAllGithubFiles", "repo", *repo.Name)
	var found = []*github.RepositoryContent{}

	log.Debug("getting files in path ... ")
	_, found, _, err = client.GetContents(ctx, *repo.Owner.Login, *repo.Name, dir, &github.RepositoryContentGetOptions{})
	// repo may not have a `.github` folder, retun nil, causing a skip rather than fatal end
	if err != nil && strings.Contains(err.Error(), "404 Not Found") {
		log.Warn("directory was not found in this repository", "dir", dir)
		return contents, nil
	} else if err != nil {
		log.Error("error getting contents", "err", err.Error())
		return
	}
	for _, d := range found {
		var c = []*github.RepositoryContent{}
		log.Debug("found content ... ", "type", *d.Type, "path", *d.Path)
		// if directory, recurse down the tree to find files in the path there
		// otherwise add the file directly to content set
		if *d.Type == "dir" {
			c, err = getAllGithubFiles(ctx, client, repo, fmt.Sprintf("./%s/", *d.Path))
			if err != nil {
				return
			}
			contents = append(contents, c...)
		} else if *d.Type == "file" && (strings.Contains(*d.Path, ".yml") || strings.Contains(*d.Path, ".yaml")) {
			contents = append(contents, d)
		}
	}
	log.Debug("complete.")
	return
}
