package main

import (
	"github.com/miekg/dns"
	"github.com/shemul/dns-over-tls/config"
	"github.com/shemul/dns-over-tls/handler"
	"time"
)

var serverUdp *dns.Server

//var serverTcp *dns.Server

func main() {

	conf := config.Config{
		//Cloudflare's dns as resolver
		UpStreamResolverIp:   "1.1.1.1",
		UpStreamResolverPort: "853",
		TCPPort:              ":53",
		UPDPort:              ":53",
		UpstreamTimeout:      time.Millisecond * 3000,
	}

	udp := dns.Server{Addr: conf.UPDPort, Net: "udp"}
	tcp := dns.Server{Addr: conf.TCPPort, Net: "tcp-tls"}

	dns.HandleFunc(".", handler.DNSHandler(conf))

	StartServer(udp)
	StartServer(tcp)
}

func StartServer(s dns.Server) {

	err := s.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
