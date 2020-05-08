package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

var GroupLoadBalance *LoadBalance

func TestProxyHandler_ServeHTTP(t *testing.T) {
	GroupLoadBalance = NewLoadBalance()
	GroupLoadBalance.AddServer(NewHttpServer("http://127.0.0.1:9192", 2, 5))
	GroupLoadBalance.AddServer(NewHttpServer("http://127.0.0.1:9198", 2, 5))
	GroupLoadBalance.WatchServers()

	if err := http.ListenAndServe(":8080", &ProxyHandler{}); err != nil {
		fmt.Println(err.Error())
	}
}

type ProxyHandler struct{}

func (*ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("server err:", err.(error).Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.(error).Error()))
			return
		}
	}()

	//chrome
	if r.URL.Path == "/favicon.ico" {
		return
	}

	requestUrl, err := url.Parse(GroupLoadBalance.SelectByWeightRand().Addr)
	fmt.Println(requestUrl)

	if err != nil {
		fmt.Println(err.Error())
	}

	proxy := httputil.NewSingleHostReverseProxy(requestUrl)
	proxy.ServeHTTP(w, r)
}
