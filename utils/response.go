package utils

// ApiResponse is used for successful API responses
type ApiResponse struct {
    Status bool        `json:"status"`
    Data   interface{} `json:"data"`
}

// ApiErrorResponse is used for error responses
type ApiErrorResponse struct {
    Status  bool   `json:"status"`
    Message string `json:"message"`
}
