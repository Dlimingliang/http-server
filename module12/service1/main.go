package main

import (
	"context"
	"github.com/golang/glog"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	glog.V(0).Info("Starting service1 server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", defaultHandler)
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

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	glog.V(0).Info("service1 handler")
	delay := randInt(10, 2000)
	glog.V(0).Info("time-delay:", delay)

	req, err := http.NewRequest("GET", "http://service2/hello", nil)
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
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	w.Write(body)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func healthyCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
