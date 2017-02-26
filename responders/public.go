package responders

import "net/http"

type public struct {
	dir     http.Dir
	handler http.Handler
}

// PublicResponder returns a new PublicResponder
func PublicResponder(docroot string) *public {
	dir := http.Dir(docroot)
	handler := http.FileServer(dir)

	return &public{dir, handler}
}

func (p public) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler.ServeHTTP(w, r)
}
