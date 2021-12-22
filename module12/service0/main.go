package main

import (
	"context"
	"github.com/Dlimingliang/http-server/metrics"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	metrics.Register()
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", defaultHandler)
	mux.HandleFunc("/preStop", preStopHandler)
	mux.HandleFunc("/healthz", healthyCheckHandler)
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	processed := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		glog.V(0).Info("receive sgiterm, prepare shutdown")
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

func preStopHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(0).Info("receive preStop...")
	time.Sleep(15 * time.Second)
	glog.V(0).Info("preStop ending...")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	glog.V(0).Info("service0 handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(10, 2000)
	glog.V(0).Info("time-delay:", delay)
	time.Sleep(time.Millisecond * time.Duration(delay))

	req, err := http.NewRequest("GET", "http://service1/hello", nil)
	if err != nil {
		glog.Error("#{err}")
	}

	lowerCaseHeader := make(http.Header)
	for key, value := range r.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	glog.V(0).Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("HTTP get fail with error:", "error", err)
	}
	if resp != nil {
		resp.Write(w)
	}
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
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
