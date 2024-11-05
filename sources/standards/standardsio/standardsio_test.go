package standardsio_test

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/sources/standards/standardsio"
)

func TestStandardsIOResolver(t *testing.T) {
	// just test as got a custom resolver
	var _ huma.Resolver = (*standardsio.StandardsInput)(nil)
}
