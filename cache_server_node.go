package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func boom() {

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Host)
		// create the reverse proxy
		url, err := url.Parse("http://" + req.Host)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(url)

		// Update the headers to allow for SSL redirection
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(res, req)

	})

	http.ListenAndServe(":8080", nil)
}
