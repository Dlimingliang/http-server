package main

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIp       = "X-Real-IP"
	OkStr         = "200"
	Ok            = 200
)

func main() {

	glog.V(0).Info("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/healthz", healthyCheckHandler)

	srv := http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	processed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); nil != err {
			glog.Error("server shutdown failed, err: %v\n", err)
		}
		glog.V(0).Info("server gracefully shutdown")
		close(processed)
	}()

	err := srv.ListenAndServe()
	if http.ErrServerClosed != err {
		glog.Error("server not gracefully shutdown, err :%v\n", err)
	}

	<-processed
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
	time.Sleep(5 * time.Second)
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
