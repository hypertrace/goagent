package sdk

type Code uint32

const (
	// Unset is the default status code.
	StatusCodeUnset Code = 0
	// Error indicates the operation contains an error.
	StatusCodeError Code = 1
	// Ok indicates operation has been validated by an Application developers
	// or Operator to have completed successfully, or contain no error.
	StatusCodeOk Code = 2
)
