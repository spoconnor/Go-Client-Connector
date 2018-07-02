package contracts

// Object represents generic message parameters.
// In real-world application it is better to avoid such types for better
// performance.
type RpcParams map[string]interface{}

type RpcRequest struct {
	ID     int       `json:"id"`
	Method string    `json:"method"`
	Params RpcParams `json:"params"`
}

type RpcResponse struct {
	ID     int       `json:"id"`
	Result RpcParams `json:"result"`
}

type RpcError struct {
	ID    int       `json:"id"`
	Error RpcParams `json:"error"`
}
