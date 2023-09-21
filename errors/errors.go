package errors

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

// Define Errors
var (
	ErrDatabase            = status.New(106001, "database error").Err()
	ErrMarshal             = status.New(106002, "marshal failed").Err()
	ErrUnmarshal           = status.New(106003, "unmarshal failed").Err()
	ErrInternalServerError = status.New(106004, "internal server error").Err()
	ErrBadRequest          = status.New(106005, "bad request").Err()
	ErrSystemError         = status.New(106006, "system error").Err()
)

// ErrorWithMessage wraps detail error
func ErrorWithMessage(err error, message string) error {
	s, ok := status.FromError(err)
	if !ok {
		return errors.Wrap(err, message)
	}

	grpcStatus := status.New(s.Code(), s.Message()+": "+message)
	return grpcStatus.Err()
}
