package lib

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs"
	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
)

func Test_processStandards(t *testing.T) {
	var (
		adaptor    dbs.Adaptor
		err        error
		ctx               = context.Background()
		dir        string = "./" // t.TempDir()
		dbFile     string = filepath.Join(dir, "test.db")
		sourceFile string = "../../convertor/converted/github_standards.json" // filepath.Join(dir, "standards.json")
	)

	adaptor, err = adaptors.NewSqlite(dbFile, false)
	err = processStandards(ctx, adaptor, sourceFile)
	defer adaptor.DB().Close()

	fmt.Println(err)
}
