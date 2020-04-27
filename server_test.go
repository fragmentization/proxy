package proxy

import (
	"fmt"
	"net/http"
	"testing"
)

func TestProxyHandler_ServeHTTP(t *testing.T) {
	if err := http.ListenAndServe(":8080", &ProxyHandler{}); err != nil {
		fmt.Println(err.Error())
	}
}
