package proxy

import (
	"fmt"
	"net/http"
	"strings"
)

type ProxyHandler struct {
	Host string
	Port string
}

func (self *ProxyHandler) NewProxyHandler() error {
	addr := strings.Join([]string{self.Host, ":", self.Port}, "")

	fmt.Println(addr)

	if err := http.ListenAndServe(addr, &ProxyHandler{}); err != nil {
		return err
	}
}

func (*ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func copyHeader(duplicator http.Header, beDuplicator *http.Header) {

	for key, value := range duplicator {
		beDuplicator.Set(key, value[0])
	}
}
