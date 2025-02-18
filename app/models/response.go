package models

// WSResponse is the standardized response format.
// - Meta : *Pre-formatted response header returning data.
// - Data : *Data or list of data returned.
type WSResponse struct {
	Meta MetaResponse `json:"meta"`
	Data interface{}  `json:"data"`
}

// MetaResponse is a valid response header
// - ObjectName : *Information returned to the front end to let it know what format it is receiving.*
// - TotalCount : *Total number of records the request can return.
// - Offset : *Starting position of the list of records returned to the Front.
// - Count : *Number of records returned to the Front.
type MetaResponse struct {
	ObjectName string `json:"object_name"`
	TotalCount int    `json:"total_count"`
	Offset     int    `json:"offSet"`
	Count      int    `json:"count"`
}

// MessageTypes is an array of message types returned to the Front
// *	+ OK                  : *200*
// *	+ Created             : *201*
//   - NotModified          : *304*
//
// !	+ BadRequest          : *400*
// !	+ Unauthorized        : *401*
// !	+ PaymentRequired     : *402*
// !	- Forbidden           : *403*
// !	* NotFound            : *404*
// !	+ MethodNotAllowed    : *405*
// !  + Conflict	 	 	      : *409*
// !	+ InternalServerError : *500*
type MessageTypes struct {
	OK                  string
	Created             string
	NotModified         string
	BadRequest          string
	Unauthorized        string
	PaymentRequired     string
	Forbidden           string
	NotFound            string
	Conflict            string
	MethodNotAllowed    string
	InternalServerError string
}

// BasicResponse is a basic response.
// - Status : *http status*
// - MessageType : *typed message for Front in I18N format*.
// - Message : *response message* // - Message
type BasicResponse struct {
	Status      int    `json:"status"`
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
}

// Success is a basic response of 2xx.
func Success(status int, messageType string, msg string) *BasicResponse {
	return &BasicResponse{Status: status, MessageType: messageType, Message: msg}
}

// Redirection is a basic response of 3xx.
func Redirection(status int, messageType string, msg string) *BasicResponse {
	return &BasicResponse{Status: status, MessageType: messageType, Message: msg}
}

// KnownError is a basic response of know errors.
func KnownError(status int, messageType string, err error) *BasicResponse {
	return &BasicResponse{Status: status, MessageType: messageType, Message: err.Error()}
}

// UnknownError is a basic response of unknown errors.
func UnknownError(status int, err error) *BasicResponse {
	return KnownError(status, "error.unknown", err)
}
