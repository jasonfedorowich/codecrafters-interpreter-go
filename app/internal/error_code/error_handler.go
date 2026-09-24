package error_code

import (
	"fmt"
	"os"
)

type ErrorHandler struct {
	ErrorCode int
}

func (eh *ErrorHandler) InError() bool {
	return eh.ErrorCode != 0
}

func (eh *ErrorHandler) Error(errorCode int, format string, a ...any) {
	eh.ErrorCode = errorCode
	fmt.Fprintf(os.Stderr, format, a...)
}
