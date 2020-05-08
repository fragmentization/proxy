package proxy

import (
	"math"
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

	//向下取整
	newFailWeight := int(math.Floor(float64(server.Weight)) * (1 / server.FailFactor))
	if newFailWeight == 0 {
		newFailWeight = 1
	}

	server.FailWeight += newFailWeight

	if server.FailWeight > server.Weight {
		server.FailWeight = server.Weight
	}
}

func (this *HttpChecker) success(server *HttpServer) {
	server.FailWeight = 0
}
