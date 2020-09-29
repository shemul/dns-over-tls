package main

import (
	"github.com/miekg/dns"
	"github.com/shemul/dns-over-tls/config"
	"github.com/shemul/dns-over-tls/handler"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

var (
	App = cli.NewApp()
)

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

	dns.HandleFunc(".", handler.DNSHandler(conf))

	App.Name = "DNS-over-TLS Proxy over Cloudflare upstream"
	App.UsageText = "ex: go run main.go udp"
	App.Commands = []cli.Command{
		{
			Name:  "udp",
			Usage: "run the UDP/53 server",
			Action: func(c *cli.Context) {
				udp := dns.Server{Addr: conf.UPDPort, Net: "udp"}
				StartServer(&udp)
			},
		},
		{
			Name:  "tcp",
			Usage: "run the TCP/53 server",
			Action: func(c *cli.Context) {
				udp := dns.Server{Addr: conf.UPDPort, Net: "tcp"}
				StartServer(&udp)
			},
		},
	}
	err := App.Run(os.Args)
	if err != nil {
		panic(err.Error())
	}
}

func StartServer(s *dns.Server) {
	log.Printf("DNS server is running on port %v/%v", s.Net, s.Addr)
	n := ""
	if s.Net == "tcp" {
		n = "+tcp"
	}
	log.Printf("try in cli : dig +short %v google.com @localhost", n)

	err := s.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
