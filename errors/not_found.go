package errors

type ErrNotFoundError struct {
	*Err
}

// NewErrNotFoundError 创建页面没有找到异常
func NewErrNotFoundError() *ErrNotFoundError {
	e := &Err{code: ErrNotFound, message: ErrText[ErrNotFound]}
	return &ErrNotFoundError{e}
}
