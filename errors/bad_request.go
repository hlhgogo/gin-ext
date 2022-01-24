package errors

type BadRequestError struct {
	*Err
}

// NewBadRequestError 创建业务异常
func NewBadRequestError(errMsg string, errCode ...int) *BadRequestError {
	code := ErrBadRequest
	if len(errCode) > 0 {
		code = errCode[0]
	}
	e := &Err{code: code, message: errMsg}
	return &BadRequestError{e}
}
