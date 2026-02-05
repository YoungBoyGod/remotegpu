package errors

// CodeX 2026-02-05: Agent task queue error codes.
const (
	ErrAttemptMismatch = 30001
	ErrTaskImmutable   = 30002
	ErrLeaseExpired    = 30003
	ErrTaskNotFound    = 30004
	ErrInvalidParams   = 30005
	ErrInternal        = 30099
)

var errorMsg = map[int]string{
	ErrAttemptMismatch: "attempt mismatch",
	ErrTaskImmutable:   "task immutable",
	ErrLeaseExpired:    "lease expired",
	ErrTaskNotFound:    "task not found",
	ErrInvalidParams:   "invalid params",
	ErrInternal:        "internal error",
}

// CodeX 2026-02-05: Message returns the default error message for an agent error code.
func Message(code int) string {
	if msg, ok := errorMsg[code]; ok {
		return msg
	}
	return "unknown error"
}
