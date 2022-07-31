package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"httpserver/metrics"
	"io"
	"k8s.io/apimachinery/pkg/util/rand"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const PNameVersion = "VERSION"

func main() {
	flag.Set("v", "4")
	//metric register prometheus
	metrics.Register()

	mux := http.NewServeMux()

	//注册系统SIGTERM信号HOOK
	handleSigterm()
	//注册路由处理器
	log.Printf("main---------start to handler http request")
	mux.HandleFunc("/", reqHandler)
	log.Printf("111")
	mux.HandleFunc("/healthz", health)
	log.Printf("222")

	mux.Handle("/metrics", promhttp.Handler())
	//log.Printf("222-metrics")

	//监听80端口,使用GO内置的HTTP Server
	svr := http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	svr.ListenAndServe()
	//go func() {
	//	if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: $s \n", err)
	//	}
	//}()
	//err := http.ListenAndServe(":80", nil)
	log.Printf("333, %v", svr)
	//if err != nil {
	//	log.Printf("startup http server failed, %v", err)
	//}
	log.Printf("end main ---------start to handler http request")

}

/**
1.接收客户端 request，并将 request 中带的 header 写入 response header
2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4.当访问 localhost/healthz 时，应返回200
*/
func reqHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("start to handler http request")

	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(10, 200)
	time.Sleep(time.Millisecond * time.Duration(delay))
	writeBackReqHeader(w, r)
	writeVersionHeader(w, r)
	logRemoteIPAndStatus(w, r)
	io.WriteString(w, "simple http server for fun!\n")
	io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	log.Printf("handle request done")
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

//1.接收客户端 request，并将 request 中带的 header 写入 response header
func writeBackReqHeader(w http.ResponseWriter, r *http.Request) {
	log.Printf("fun1: write back request header to response")
	for k, v := range r.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}
}

//2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
func writeVersionHeader(w http.ResponseWriter, r *http.Request) {
	versionValue := os.Getenv(PNameVersion)
	log.Printf("env param: version = %s ", versionValue)
	w.Header().Set(PNameVersion, versionValue)
}

//3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
func logRemoteIPAndStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Printf("remote addr:%s, resp status code %d", r.RemoteAddr, http.StatusOK)
}

//4.当访问 localhost/healthz 时，应返回200
func health(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle request: /healthz, and response 200_OK")
	reqHandler(w, r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Custom-Header", "Awesome")
	io.WriteString(w, "200")
}

//5.处理SIGTERM信号
func handleSigterm() {
	log.Printf("register sigterm signal handler")
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			//handle SIGINT
		case syscall.SIGTERM:
			//handle SIGTERM
			time.Sleep(2 * time.Second)
			log.Printf("handle system sigterm signal for 2 seconds")
			//system exit normally
			os.Exit(0)

		}
	}()
}

//6. handler metrics
func metricHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handle request: /metrics, and response 200_OK")
	reqHandler(w, r)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Custom-Header", "Awesome")
	io.WriteString(w, "200")
}
