// Package ghfilters contains some standard filters used when querying the github api
// to handle elements that the api endpoints / sdk does not provide directly
package ghfilters

import (
	"context"
	"fmt"
	"opg-reports/app/internal/logx"
	"strings"

	"github.com/google/go-github/v87/github"
)

// ExcludeArchivedRepository is a very simple filter to remove archieved repositories from the
// original data set
type ExcludeArchivedRepository struct{}

// Filter checks the archive status of the repository and returns the inverse value.
//
// When `Archived == true`, `include = false`
func (self *ExcludeArchivedRepository) Filter(ctx context.Context, result *github.Repository) (include bool) {
	var _, lg = logx.Get(ctx)

	include = !*result.Archived
	lg.Debug(fmt.Sprintf("[%s] archived: [%v], include = [%v] ", *result.FullName, *result.Archived, include))

	return
}

// FilterByRepositoryName will only return a repository whose short name exactly matches
// the name property of this filter - allowing to find a specific repo out of a larger
// set
type FilterByRepositoryName struct {
	Name string // the is the short name of the repository we're looking for
}

// Filter checks that the repository name exactly matches the set value and only returns
// true for those than do.
//
// Sets both to lowercase.
func (self *FilterByRepositoryName) Filter(ctx context.Context, result *github.Repository) (include bool) {
	var (
		_, lg  = logx.Get(ctx)
		name   = strings.ToLower(*result.Name)
		target = strings.ToLower(self.Name)
	)

	include = (name == target)
	lg.Debug(fmt.Sprintf("[%s] repo name match: ([%s] == [%s]), include = [%v] ", *result.FullName, target, name, include))
	return
}

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
