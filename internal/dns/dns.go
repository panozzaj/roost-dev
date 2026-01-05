package dns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

// Server is a simple DNS server that responds to all queries with 127.0.0.1
type Server struct {
	port   int
	tld    string
	server *dns.Server
}

// New creates a new DNS server
func New(port int, tld string) *Server {
	return &Server{
		port: port,
		tld:  tld,
	}
}

// Start starts the DNS server
func (s *Server) Start() error {
	s.server = &dns.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", s.port),
		Net:  "udp",
	}

	dns.HandleFunc(s.tld+".", s.handleQuery)
	dns.HandleFunc(".", s.handleQuery) // Handle all queries

	return s.server.ListenAndServe()
}

// Stop stops the DNS server
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Shutdown()
	}
}

// handleQuery responds to DNS queries with 127.0.0.1
func (s *Server) handleQuery(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA:
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP("127.0.0.1"),
			})
		case dns.TypeAAAA:
			m.Answer = append(m.Answer, &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				AAAA: net.ParseIP("::1"),
			})
		}
	}

	w.WriteMsg(m)
}
