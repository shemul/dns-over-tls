package config

import "time"

type Config struct {
	UpStreamResolverIp   string
	UpStreamResolverPort string
	TCPPort              string
	UPDPort              string
	UpstreamTimeout      time.Duration
}
