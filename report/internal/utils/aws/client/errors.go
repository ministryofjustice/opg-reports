package client

import "errors"

var ErrLoadingConfig error = errors.New("error loading config.")
var ErrUnsupportedType error = errors.New("client type unsupported.")
