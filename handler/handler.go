package handler

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/shemul/dns-over-tls/config"
	"log"
	"net"
)

// This handler will exchange the dns queries to cloudfare DoT for answering the dns queries
func DNSHandler(conf config.Config) func(writer dns.ResponseWriter, msg *dns.Msg) {
	return func(writer dns.ResponseWriter, msg *dns.Msg) {
		qString := ""
		for _, v := range msg.Question {
			qString += v.String()
		}
		c := new(dns.Client)
		c.Net = "tcp-tls"
		c.Dialer = &net.Dialer{
			Timeout: conf.UpstreamTimeout,
		}

		ans, rtt, err := c.Exchange(msg, fmt.Sprintf("%v:%v", conf.UpStreamResolverIp, conf.UpStreamResolverPort))
		if err != nil {
			log.Printf("failed to communicate with upstream: %s", err)
			return
		}
		log.Printf("Answer for '%s' received in %s", qString, rtt.String())
		writer.WriteMsg(ans)
	}
}
