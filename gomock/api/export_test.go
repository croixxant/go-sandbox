package api

import "net/http"

func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	server.router.ServeHTTP(w, req)
}
