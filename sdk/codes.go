package sdk

type Code uint32

const (
	// Unset is the default status code.
	Unset Code = 0
	// Error indicates the operation contains an error.
	Error Code = 1
	// Ok indicates operation has been validated by an Application developers
	// or Operator to have completed successfully, or contain no error.
	Ok Code = 2
)
