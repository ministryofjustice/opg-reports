package costsio_test

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/ministryofjustice/opg-reports/sources/costs/costsio"
)

func TestCostsIOResolver(t *testing.T) {
	// just test as got a custom resolver
	var _ huma.Resolver = (*costsio.GroupedDateRangeInput)(nil)

}
