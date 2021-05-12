package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8080"
	fmt.Println("WAPot v1.0\n\n[+] Listening on port", port)
	http.Handle("/", http.FileServer(http.Dir("/var/www/html")))
	if err := http.ListenAndServe(port, logRequest(http.DefaultServeMux)); err != nil {
		panic(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("%s %s %s \"%s\" \"%s\" \n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Method, req.URL.RequestURI(), req.UserAgent())
		handler.ServeHTTP(w, req)
	})
}
