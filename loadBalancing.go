package proxy

import (
	"math/rand"
	"sort"
	"time"
)

type HttpServers []*HttpServer

func (p HttpServers) Len() int           { return len(p) }
func (p HttpServers) Less(i, j int) bool { return p[i].CurrentWeight > p[j].CurrentWeight }
func (p HttpServers) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type HttpServer struct {
	Addr          string
	Weight        int
	CurrentWeight int
	FailWeight    int
	FailFactor    float64
}

type LoadBalance struct {
	Servers HttpServers
}

var LB *LoadBalance
var ServerIndices []int
var SumWeight int //the total weight

func init() {
	rand.Seed(time.Now().UnixNano())
	LB = NewLoadBalance()
	LB.AddServer(NewHttpServer("http://127.0.0.1:9091", 2))
	LB.AddServer(NewHttpServer("http://127.0.0.1:9092", 1))

	for index, server := range LB.Servers {
		if server.Weight > 0 {
			for i := 0; i < server.Weight; i++ {
				ServerIndices = append(ServerIndices, index)
			}
			SumWeight += server.Weight
		}

	}

}

func NewHttpServer(addr string, weight int) *HttpServer {
	return &HttpServer{
		Addr:          addr,
		Weight:        weight,
		CurrentWeight: 0,
		FailFactor:    5,
	}
}

func NewLoadBalance() *LoadBalance {
	return &LoadBalance{
		Servers: make(HttpServers, 0),
	}
}

func (this *LoadBalance) AddServer(server *HttpServer) {
	this.Servers = append(this.Servers, server)
}

func (this *LoadBalance) getSumWeight() (sum int) {
	for _, server := range LB.Servers {
		newWeight := server.Weight - server.FailWeight
		if newWeight > 0 {
			sum = sum + newWeight
		}
	}
	return
}

//平滑加权轮询
func (this *LoadBalance) SelectByWeightRand() *HttpServer {
	for _, server := range this.Servers {
		server.CurrentWeight += (server.Weight - server.FailWeight)
	}

	sort.Sort(this.Servers)
	maxWeightServer := this.Servers[0]

	maxWeightServer.CurrentWeight -= this.getSumWeight()

	go Heartbeat(LB.Servers)

	//strCurrentWeight := ""
	//for _, server := range this.Servers {
	//	strCurrentWeight += fmt.Sprintf("%d,", server.CurrentWeight)
	//}

	return maxWeightServer
}

//heartbeat
func Heartbeat(servers HttpServers) {
	t := time.NewTicker(time.Second * 5)
	httpCheck := NewHttpChecker(servers)
	for {
		select {
		case <-t.C:
			httpCheck.Check(time.Second * 2)

		}
	}
}
