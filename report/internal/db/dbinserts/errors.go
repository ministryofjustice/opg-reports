package dbinserts

import "errors"

var ErrTransactionBeginFailed = errors.New("transaction begin failed with error.")
var ErrPreparedInsertFailed = errors.New("prepared insert stmt failed with error.")
var ErrGetContextFailed = errors.New("stmt context failed with error.")
var ErrTransactionExecFailed = errors.New("transaction insert failed with error.")
var ErrTransactionCommitFailed = errors.New("transaction commit failed with error.")
var ErrMissingResults = errors.New("error with returned results.")
