package errors

type UnauthorizedError struct {
	*Err
}

// NewUnauthorizedError 创建没有权限异常
func NewUnauthorizedError() *UnauthorizedError {
	e := &Err{code: ErrStatusUnauthorized, message: ErrText[ErrStatusUnauthorized]}
	return &UnauthorizedError{e}
}
