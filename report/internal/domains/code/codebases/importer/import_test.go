package importer

import (
	"opg-reports/report/internal/domains/code/types"
	"opg-reports/report/packages/types/interfaces"

	"github.com/google/go-github/v84/github"
)

var _ interfaces.ImportGetterF[*github.Repository, *github.Repository, *github.Client] = Get
var _ interfaces.ImportFilterF[*github.Repository] = Filter
var _ interfaces.ImportTransformF[*types.Codebase, *github.Repository] = Transform
