package response

const (
	// 2xx Success
	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNoContent            = 204

	// 3xx Redirection
	StatusMovedPermanently     = 301
	StatusFound                = 302
	StatusNotModified          = 304

	// 4xx Client Errors
	StatusBadRequest           = 400
	StatusUnauthorized         = 401
	StatusPaymentRequired      = 402
	StatusForbidden            = 403
	StatusNotFound             = 404
	StatusMethodNotAllowed     = 405
	StatusConflict             = 409
	StatusUnprocessableEntity  = 422
	StatusTooManyRequests      = 429

	// 5xx Server Errors
	StatusInternalServerError  = 500
	StatusNotImplemented       = 501
	StatusBadGateway           = 502
	StatusServiceUnavailable   = 503
	StatusGatewayTimeout       = 504
)