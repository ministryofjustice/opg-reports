package releasesio_test

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/sources/releases/releasesio"
)

func TestReleasesIOResolver(t *testing.T) {
	// just test as got a custom resolver
	var _ huma.Resolver = (*releasesio.ReleasesInput)(nil)
}
