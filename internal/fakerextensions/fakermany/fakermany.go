package fakermany

import (
	"log/slog"

	"github.com/go-faker/faker/v4"
	"github.com/ministryofjustice/opg-reports/internal/dbs"
)

// Many returns multiple faked versions of T
func Fake[T dbs.Cloneable](n int) (faked []T) {
	slog.Debug("[exfaker.Many] faking many", slog.Int("n", n))

	faked = []T{}
	for i := 0; i < n; i++ {
		var item T
		var record = item.New().(T)
		if e := faker.FakeData(record); e == nil {
			faked = append(faked, record)
		} else {
			slog.Error("[exfaker.Many]", slog.String("err", e.Error()))
		}
	}
	faker.ResetUnique()
	return
}
