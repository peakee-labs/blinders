package transport

type Response struct {
	// The status code of the response. Currently, codes 200 and 400 are supported by the transport.
	//
	// SUPPORTED CODES: [ 200: Successful, 400: Failed ]
	Code int
	// If the request is marked as failed from the target, the response body will include the 'error' field that contains the message passed from the target to the caller.
	//
	// This field contains an error message
	Message      string
	ResponseBody []byte // This field will contains response byte
}
