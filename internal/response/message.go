package response

var DefaultMessages = map[int]string{

	// 2xx Success
	StatusOK:                  "Success",
	StatusCreated:             "Resource created successfully",
	StatusAccepted:            "Request accepted",
	StatusNoContent:           "No content",

	// 3xx Redirection
	StatusMovedPermanently:    "Resource moved permanently",
	StatusFound:               "Resource found",
	StatusNotModified:         "Resource not modified",

	// 4xx Client Errors
	StatusBadRequest:          "Bad request",
	StatusUnauthorized:       "Unauthorized",
	StatusPaymentRequired:    "Payment required",
	StatusForbidden:          "Forbidden",
	StatusNotFound:           "Resource not found",
	StatusMethodNotAllowed:   "Method not allowed",
	StatusConflict:           "Conflict occurred",
	StatusUnprocessableEntity:"Validation failed",
	StatusTooManyRequests:    "Too many requests",

	// 5xx Server Errors
	StatusInternalServerError:"Internal server error",
	StatusNotImplemented:     "Feature not implemented",
	StatusBadGateway:         "Bad gateway",
	StatusServiceUnavailable: "Service unavailable",
	StatusGatewayTimeout:     "Gateway timeout",
}