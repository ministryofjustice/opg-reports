package arguments

type constraint interface {
	*Api | *DB | *Versions | *ApiHost | *FrontHost
}

// Default returns T with populated default values that typically used
// with the type.
func Default[T constraint]() T {
	var def interface{}
	var arg T

	switch any(arg).(type) {
	case *Api:
		def = &Api{
			DB:      Default[*DB](),
			Version: Default[*Versions](),
			Info:    Default[*ApiHost](),
		}
	case *ApiHost:
		def = &ApiHost{
			host: host{
				Name:     `api`,
				Hostname: `:8081`,
			},
		}
	case *FrontHost:
		def = &FrontHost{
			host: host{
				Name:     `front`,
				Hostname: `:8080`,
			},
		}
	case *Versions:
		def = &Versions{
			Version: `0.0.1`,
			SHA:     `abcedf`,
		}
	case *DB:
		def = &DB{
			Driver:   `sqlite3`,
			Filepath: `./database/api.db`,
			Params:   `?_journal=WAL&_busy_timeout=5000&_vacuum=incremental&_synchronous=NORMAL&_cache_size=1000000000`,
		}
	}

	return def.(T)
}
