package dbselects

import "errors"

var (
	ErrTransactionBeginFailed = errors.New("transaction begin failed with error.")
	ErrPreparedStmtFailed     = errors.New("prepared stmt failed with error.")
	ErrMissingResults         = errors.New("error with returned results.")
	ErrFailedTx               = errors.New("error comitting txn.")
	ErrMissingTable           = errors.New("missing table.")
)
