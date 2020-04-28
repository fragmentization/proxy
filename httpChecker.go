package proxy

import (
	"fmt"
	"net/http"
	"time"
)

type HttpChecker struct {
	Servers HttpServers
}

func NewHttpChecker(servers HttpServers) *HttpChecker {
	return &HttpChecker{
		Servers: servers,
	}
}

func (this *HttpChecker) Check(timeout time.Duration) {
	client := http.Client{}
	for _, server := range this.Servers {
		response, err := client.Head(server.Addr)
		if response != nil {
			defer response.Body.Close()
		}

		if err != nil {
			this.fail(server)
			continue
		}

		if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusInternalServerError {
			this.success(server)
		} else {
			this.fail(server)
		}
	}
}

func (this *HttpChecker) fail(server *HttpServer) {
	fmt.Println(server.Addr, "fail")
	server.FailWeight += server.Weight * (1 / SumWeight)
}

func (this *HttpChecker) success(server *HttpServer) {
	fmt.Println(server.Addr, "success")
	server.FailWeight = 0
}
