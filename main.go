package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
)

var server_udp *dns.Server
var server_tcp *dns.Server

func main() {

	udp := dns.Server{Addr: ":53", Net: "udp"}
	//tcp :=dns.Server{Addr: ":53", Net: "tcp"}

	dns.HandleFunc(".", func(writer dns.ResponseWriter, msg *dns.Msg) {
		println("got")
		rString := ""
		for _, v := range msg.Question {
			rString += v.String()
		}

		c := new(dns.Client)
		c.Net = "tcp-tls"
		c.Dialer = &net.Dialer{
			Timeout: 2000,
		}

		a, rtt, err := c.Exchange(msg, fmt.Sprintf("%s:%s", "1.1.1.1", "853"))
		if err != nil {
			log.Printf("failed to communicate with upstream: %s", err)
			return
		}
		log.Printf("Answer for '%s' received in %s", rString, rtt.String())
		writer.WriteMsg(a)
	})

	//err := tcp.ListenAndServe()
	err := udp.ListenAndServe()
	if err != nil {
		println(err.Error())
	}

	//defer tcp.Shutdown()
	defer udp.Shutdown()

}
