package conn

import (
	"fmt"
	"os"
	"path/filepath"
)

func SqlitePath(db string, params string) string {
	os.MkdirAll(filepath.Dir(db), os.ModePerm)

	if params == "" {
		params = "?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000"
	}
	return fmt.Sprintf("%s%s", db, params)
}
