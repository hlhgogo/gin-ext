package errors

type Err struct {
	code    int    // error code
	message string // error message
}

// Code 获取响应的业务错误码
func (err *Err) Code() int {
	return err.code
}

// Error
func (err *Err) Error() string {
	return err.Message()
}

// Message 获取响应的msg
func (err *Err) Message() string {
	var msg string
	if err.message != "" {
		msg = err.message
	}

	return msg
}

// NewErr 创建error对象
func NewErr(message string) *Err {
	code := ErrInternalServerError
	return &Err{
		code:    code,
		message: message,
	}
}
