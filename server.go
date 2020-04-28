package proxy

import (
	"net/http"

	"fmt"
	"net/url"

	"net/http/httputil"
)

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

	requestUrl, err := url.Parse(LB.SelectByWeightRand().Addr)
	fmt.Println(requestUrl)

	if err != nil {
		fmt.Println(err.Error())
	}

	proxy := httputil.NewSingleHostReverseProxy(requestUrl)
	proxy.ServeHTTP(w, r)
}
