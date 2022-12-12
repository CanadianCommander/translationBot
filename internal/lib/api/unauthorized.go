package api

// NewUnauthorizedError
// Build an unauthorized api
func NewUnauthorizedError() ApiResponse {
	return ApiResponse{
		Message: "You must authorize your request",
	}
}
