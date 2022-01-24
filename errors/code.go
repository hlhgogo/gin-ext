package errors

const (
	Success = 10200

	ErrBadRequest          = 10400
	ErrStatusUnauthorized  = 10401
	ErrNotFound            = 10404
	ErrInternalServerError = 10500
)

var ErrText = map[int]string{
	Success:                "ok",
	ErrBadRequest:          "bad request",
	ErrStatusUnauthorized:  "status unauthorized",
	ErrNotFound:            "not found",
	ErrInternalServerError: "internal server error",
}
