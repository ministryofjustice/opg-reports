package ghfilters

import (
	"context"
	"fmt"
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
		_, lg  = logx.Get(ctx)
		name   = strings.ToLower(*result.Name)
		target = strings.ToLower(self.Name)
	)
	include = strings.Contains(name, target)

	lg.Debug(fmt.Sprintf("[%d] workflow name match: ([%s] == [%s]), include = [%v] ", *result.ID, target, name, include))

	return
}
