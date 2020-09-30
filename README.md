
## Background  
  
Our applications don't handle DNS-over-TLS by default. But there are some hacks to enable it. But for now lets assume we don't want to hack our clients first. So our task is to design and create a simple DNS to DNS-over-TLS proxy that we could use to enable our application to query a DNS-over-TLS server. So our DNS queries will be secure.  
  
![alt text](https://i.imgur.com/rm6cQwv.jpg "Without DoT")  
  
  
  
## DNS over TLS Proxy (DoT)   
So, A DNS over TLS proxy that accepts simple (conventional) DNS requests and proxy it to a DNS servers running with DNS over TLS (DoT) (eg. cloudflare). So this sidecar DNS proxy will proxify our DNS queries to a DoT server.

![alt text](https://i.imgur.com/gjoygas.jpg "Title")

## Getting started

This program is writen in Golang and depend on the [miekg/dns](https://github.com/miekg/dns) library. Miekg/dns library is used for great projects as coredns.

Also a Docker image is available on [DockerHub](https://hub.docker.com/repository/docker/shemul/dns-over-tls)

**UDP**

    docker run -it -p 53:53/UDP shemul/dns-over-tls:latest /bin/app udp
  
  To test `dig +short  google.com @localhost`

**TCP**

    docker run -it -p 53:53 shemul/dns-over-tls:latest /bin/app tcp

to test `dig +short +tcp google.com @localhost`

## Implementation:
For simplicity I assumed there will be some configerations for this app. that can be come from Env or any config source. here is my simple config for now. 

    conf := config.Config{  
       //Cloudflare's dns as resolver  
      UpStreamResolverIp:   "1.1.1.1",  
      UpStreamResolverPort: "853",  
      TCPPort:              ":53",  
      UPDPort:              ":53",  
      UpstreamTimeout:      time.Millisecond * 3000,  
    }
and 

    ans, rtt, err := c.Exchange(msg, fmt.Sprintf("%v:%v", conf.UpStreamResolverIp, conf.UpStreamResolverPort))
this `Exchange` method actually performs a synchronous query to a DoT server when the DNS client type is `tcp-tls`  to get Answer from DoT

## Security Concerns:

This proxy allow us encrypted connection to upstream DoT servers, but all the traffic until this service, including its responses to clients, still not secure. For example, if you host this service in a public address and your DNS client points to it over public internet access, you can be a victim of a [man in the middle attack](https://en.wikipedia.org/wiki/Man-in-the-middle_attack).

Another thing is when this proxy will be deployed in a DHCP server. then this proxy will act as a default DNS for that DHCP server. So all of our DNS queries will hit our proxy unencrypted and insecure. 

## Distributed environment:

Runnig this container as a daemonset (so that it runs on every node) with  `hostNetwork: true` in Kubernetes. Then every Node that uses the localhost address as its own NameServer. Kubernetes uses DNS service that resolves the cluster internal names and the Pods of this DNS service will talk with proxy daemonset directly. 