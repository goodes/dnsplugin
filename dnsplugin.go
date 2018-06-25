// Package dnsplugin implements a plugin that returns details about the resolving
// querying it.
package dnsplugin

import (
	"context"
	"github.com/coredns/coredns/request"

	"fmt"
	"github.com/miekg/dns"
	"net"
	"strconv"
)

// PhDNS is a struct
type PhDNS struct{}

// Name is a method
func (phd PhDNS) Name() string { return "dnsplugin" }

// recordAAAA is a method
func recordAAAA(ip string, state request.Request) dns.RR {
	var rr dns.RR
	rr = new(dns.AAAA)
	rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
	rr.(*dns.AAAA).AAAA = net.ParseIP(ip)
	return rr
}

// recordA is a method
func recordA(ip string, state request.Request) dns.RR {
	var rr dns.RR
	rr = new(dns.A)
	rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
	rr.(*dns.A).A = net.ParseIP(ip).To4()
	return rr
}

// ServeDNS is a method
func (phd PhDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	ip := state.IP()
	var rr dns.RR

	switch state.Family() {
	case 1:
		//switch state.QType() {
		//	case dns.TypeAAAA:
		//		rr = recordAAAA(ip, state)
		//	default:
		rr = recordA(ip, state)
		//}

	case 2:
		rr = recordA(ip, state)
		//rr = recordAAAA(ip, state)
	}

	srv := new(dns.SRV)
	srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
	if state.QName() == "." {
		srv.Hdr.Name = "_" + state.Proto() + state.QName()
	}
	port, _ := strconv.Atoi(state.Port())
	srv.Port = uint16(port)
	srv.Target = "."

	cname := new(dns.CNAME)
	cname.Hdr = dns.RR_Header{Name: "test." + state.QName(), Rrtype: dns.TypeCNAME, Class: state.QClass()}
	cname.Target = "a.b.c.d"
	a.Extra = []dns.RR{rr, srv, cname}
	for _, entry := range a.Extra {
		fmt.Println(entry)
	}
	fmt.Println("-------------------")

	state.SizeAndDo(a)
	w.WriteMsg(a)

	return 0, nil
}
