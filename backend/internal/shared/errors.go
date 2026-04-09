package shared

import (
	"net/http"
)

type AppError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e AppError) Error() string {
	return e.Message
}

var (
	ErrNetworkNotFound = &AppError{http.StatusNotFound, "network_not_found", "network not found"}
	ErrPeerNotFound    = &AppError{http.StatusNotFound, "peer_not_found", "peer not found"}
)
