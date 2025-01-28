package cstmerr

// Custom HTTP Error for REST HTTP
// that has Code and Message and Error fields
// for user friendly answer and server logs
type HTTPErr struct {
	Code    int         `json:"code"`    // http code
	Message string      `json:"message"` // message for client
	Err     error       `json:"-"`       // message for server logs
	Details interface{} `json:"details,omitempty"`
}

func (e *HTTPErr) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func New(statusCode int, message string, err error, details ...interface{}) *HTTPErr {
	httpErr := &HTTPErr{
		Code:    statusCode,
		Message: message,
		Err:     err,
	}
	if len(details) != 0 {
		httpErr.Details = details
	}

	return httpErr
}
