package jsonrpc2

// ErrorCode represent the int JSON RPC 2.0 error code.
type ErrorCode int

// Official JSON-RPC 2.0 Spec Error Codes and Messages
const (
	LowestReservedErrorCode  ErrorCode = -32768
	ParseErrorCode           ErrorCode = -32700
	InvalidRequestCode       ErrorCode = -32600
	MethodNotFoundCode       ErrorCode = -32601
	InvalidParamsCode        ErrorCode = -32602
	InternalErrorCode        ErrorCode = -32603
	HighestReservedErrorCode ErrorCode = -32000

	ParseErrorMessage     = "Parse error"
	InvalidRequestMessage = "Invalid Request"
	MethodNotFoundMessage = "Method not found"
	InvalidParamsMessage  = "Invalid params"
	InternalErrorMessage  = "Internal error"
)

// IsReserved returns true if c is within the reserved error code range:
// [LowestReservedErrorCode, HighestReservedErrorCode].
func (c ErrorCode) IsReserved() bool {
	return LowestReservedErrorCode <= c && c <= HighestReservedErrorCode
}
