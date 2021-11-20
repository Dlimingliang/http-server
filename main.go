package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"net"
	"net/http"
	"os"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIp       = "X-Real-IP"
	OkStr         = "200"
	Ok            = 200
)

func main() {
	glog.V(0).Info("Starting http server...")
	glog.Error("i am error message")
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/healthz", healthyCheckHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		glog.Error("listen error", err.Error())
		return
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			w.Header().Set(name, h)
		}
	}
	glog.V(0).Info("env-version:", os.Getenv("VERSION"))
	w.Header().Set("VERSION", os.Getenv("VERSION"))
	ip := remoteIp(r)
	glog.V(0).Info("client ip:", ip)
	resp, _ := json.Marshal(map[string]string{
		"ip": ip, "statusCode": OkStr,
	})
	w.Write(resp)
}

func healthyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(Ok)
}

func remoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIp); ip != "" {
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
