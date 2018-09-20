package contracts

// Object represents generic message parameters.
// In real-world application it is better to avoid such types for better
// performance.
type RpcParams map[string]interface{}

type RpcRequest struct {
	ID     int       `json:"Id"`
	Method string    `json:"Method"`
	Params RpcParams `json:"Params"`
}

type RpcResponse struct {
	ID     int         `json:"Id"`
	Result interface{} `json:"Result"`
	Error  RpcError    `json:"Error"`
}

type RpcError struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

const (
	/// <summary>Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.</summary>
	ParserError = -32700
	/// <summary>The JSON sent is not a valid Request object.</summary>
	InvalidRequest = -32600
	/// <summary>The method does not exist / is not available.</summary>
	MethodNotFound = -32601
	/// <summary>Invalid method parameter(s).</summary>
	InvalidParams = -32602
	/// <summary>Internal JSON-RPC error.</summary>
	InternalError = -32603
	/// <summary>Reserved for implementation-defined server-errors.</summary>
	ServerError = -32000
)
