package bitcast_go

import "errors"

var (
	ErrKeyIsEmpty            = errors.New("the key is empty")
	ErrIndexUpdateFailed     = errors.New("failed to updata index")
	ErrKeyNotFound           = errors.New("key not found in database")
	ErrDataFileNotFound      = errors.New("data file is not found")
	ErrDataDirectryCorrupted = errors.New("the data base directory maybe corrupted")
	ErrExceedMaxBatchNum     = errors.New("exceed the max batch num")
	ErrMergeIsProgress       = errors.New("merge is in progress,try again later")
)
