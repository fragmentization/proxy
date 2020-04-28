package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NuoMinMin/proxy"
)

func TestProxyHandler_ServeHTTP(t *testing.T) {
	if err := http.ListenAndServe(":8080", &proxy.ProxyHandler{}); err != nil {
		fmt.Println(err.Error())
	}
}
