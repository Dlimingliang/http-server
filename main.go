package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/http"
	"os"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
	OK 			  = "200"
)

func main() {
	glog.V(0).Info("Starting http server...")
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/healthz", healthyCheckHandler)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		fmt.Println("listen error", err.Error())
		return
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	ip := remoteIp(r)
	glog.V(0).Info("client ip:", ip)
	for name, headers := range r.Header {
		for _, h := range headers {
			w.Header().Set(name, h)
		}
	}
	w.Header().Set("VERSION", os.Getenv("VERSION"))
	resp, _ := json.Marshal(map[string]string{
		"ip": ip, "statusCode": OK,
	})
	w.Write(resp)
}

func healthyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func remoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
