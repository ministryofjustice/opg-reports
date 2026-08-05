package ghfilters

import (
	"context"
	"log/slog"
	"opg-reports/app/internal/logx"
	"strings"

	"github.com/google/go-github/v87/github"
)

// --- WORKFLOW LEVEL FILTERS

// FilterWorkfowRunByPartialName looks for workflows whose name attribute contains the local name value.
//
// Used to filter all workflow runs to just path to live.
type FilterWorkfowRunByPartialName struct {
	Name string // Name is partial string match we want to match in the workflow runs name
}

// Filter checks that the workflow name contains the name value we have configured.
//
// Sets both to lowercase.
func (self *FilterWorkfowRunByPartialName) Filter(ctx context.Context, result *github.WorkflowRun) (include bool) {
	var (
		name   string       = strings.ToLower(*result.Name)
		target string       = strings.ToLower(self.Name)
		lg     *slog.Logger = logx.Default().With(
			"name", name,
			"targetName", target,
		)
	)
	include = strings.Contains(name, target)

	lg.Debug("include?", "include", include)

	return
}
