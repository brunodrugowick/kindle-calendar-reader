package api

import (
	"net/http"
)

type Api interface {
	HandleRequests(w http.ResponseWriter, r *http.Request)
	GetPath() string
}
