package files

import "errors"

var ErrFileExists = errors.New("file exists, cannot copy to this lcoation,")
