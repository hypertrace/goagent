package grpc

import (
	"google.golang.org/grpc/codes"
)

func StatusText(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 402:
		return "Payment Required"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 406:
		return "Not Acceptable"
	case 407:
		return "Proxy Authentication Required"
	case 408:
		return "Request Timeout"
	case 409:
		return "Conflict"
	case 410:
		return "Gone"
	case 411:
		return "Length Required"
	case 412:
		return "Precondition Failed"
	case 413:
		return "Request Entity Too Large"
	case 414:
		return "Request URI Too Long"
	case 415:
		return "Unsupported Media Type"
	case 416:
		return "Requested Range Not Satisfiable"
	case 417:
		return "Expectation Failed"
	case 418:
		return "I'm a teapot"
	case 421:
		return "Misdirected Request"
	case 422:
		return "Unprocessable Entity"
	case 423:
		return "Locked"
	case 424:
		return "Failed Dependency"
	case 425:
		return "Too Early"
	case 426:
		return "Upgrade Required"
	case 428:
		return "Precondition Required"
	case 429:
		return "Too Many Requests"
	case 431:
		return "Request Header Fields Too Large"
	case 451:
		return "Unavailable For Legal Reasons"
	default:
		return "Request Error"
	}
}

// StatusCode does a best effort mapping from HTTP Request Status code to GRPC Code.
func StatusCode(code int) codes.Code {
	switch code {
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 407:
		// "Proxy Authentication Required"
		return codes.Unauthenticated
	case 408:
		// Request Timeout
		return codes.DeadlineExceeded
	case 412:
		// "Precondition Failed"
		return codes.FailedPrecondition
	case 413, // "Request Entity Too Large"
		414, // "Request URI Too Long"
		429, // "Too Many Requests"
		431: // "Request Header Fields Too Large"
		return codes.ResourceExhausted
	default:
		return codes.Unknown
	}
}
