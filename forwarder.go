package httptools

import (
	"net/http"
	"net/http/httputil"
)

func SSLForwarderHandleFunc(host string) func(w http.ResponseWriter, r *http.Request) {
	revproxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme="https"
			r.URL.Host=host
			r.Host=host
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		revproxy.ServeHTTP(w, r)
	}
}

func HttpForwarderHandleFunc(host string) func(w http.ResponseWriter, r *http.Request) {
	revproxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme="http"
			r.URL.Host=host
			r.Host=host
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		revproxy.ServeHTTP(w, r)
	}
}
