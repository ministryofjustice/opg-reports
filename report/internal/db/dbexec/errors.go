package dbexec

import "errors"

var ErrTransactionBeginFailed = errors.New("transaction begin failed with error.")
var ErrTransactionExecFailed = errors.New("transaction exec failed with error.")
var ErrTransactionCommitFailed = errors.New("transaction commit failed with error.")
