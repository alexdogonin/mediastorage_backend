package router

import "net/http"

func NewRouter(service interface{}) http.Handler {
	mux := http.NewServeMux()

	return mux
}
