package ghfilters

import (
	"context"
	"log/slog"
	"opg-reports/app/internal/logx"
	"strings"

	"github.com/google/go-github/v87/github"
)

// --- REPOSITORY LEVEL FILTERS

// ExcludeArchivedRepository is a very simple filter to remove archieved repositories from the
// original data set
type ExcludeArchivedRepository struct{}

// Filter checks the archive status of the repository and returns the inverse value.
//
// When `Archived == true`, `include = false`
func (self *ExcludeArchivedRepository) Filter(ctx context.Context, result *github.Repository) (include bool) {
	var lg *slog.Logger = logx.Default().With(
		"repository", *result.FullName,
		"archived", *result.Archived,
	)

	include = !*result.Archived

	lg.Debug("include?", "include", include)

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
		name   string       = strings.ToLower(*result.Name)
		target string       = strings.ToLower(self.Name)
		lg     *slog.Logger = logx.Default().With(
			"repository", *result.FullName,
			"name", name,
			"targetName", target,
		)
	)

	include = (name == target)
	lg.Debug("include?", "include", include)

	return
}
