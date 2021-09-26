package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/header", helloHttpHandler)
	http.HandleFunc("/healthz", healthyCheckHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("listen error", err.Error())
		return
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Welcome to http server")
}

func helloHttpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("receive request")
	for name, headers := range r.Header {
		fmt.Println("当前的header", name, headers)
		for _, h := range headers {
			fmt.Println("当前的headervalue", h)
			w.Header().Set(name, h)
		}
	}
	w.Header().Set("VERSION", os.Getenv("VERSION"))
}

func healthyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
