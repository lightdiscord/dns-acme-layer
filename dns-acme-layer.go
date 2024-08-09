package main

import (
	"context"
	"github.com/caddyserver/certmagic"
	"github.com/miekg/dns"
	"log"
)

func main() {
	provider := Provider{}

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = "infra-test@example.com"
	certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	certmagic.DefaultACME.DNS01Solver = &certmagic.DNS01Solver{
		DNSManager: certmagic.DNSManager{
			DNSProvider: &provider,
		},
	}

	magic := certmagic.NewDefault()

	err := magic.ManageAsync(context.Background(), []string{"wow.example.com", "*.wow.example.com"})
	if err != nil {
		return
	}

	dns.HandleFunc(".", func(writer dns.ResponseWriter, msg *dns.Msg) {
		domain := msg.Question[0].Name

		m := new(dns.Msg)
		m.SetReply(msg)
		m.Authoritative = true

		for _, entry := range provider.entries {
			if entry.record.Type != "TXT" {
				log.Fatal("expected TXT record")
			}

			if entry.record.Name+"."+entry.zone == domain {
				rr := new(dns.TXT)
				rr.Hdr = dns.RR_Header{
					Name:   domain,
					Rrtype: dns.TypeTXT,
					Class:  dns.ClassINET,
					Ttl:    uint32(entry.record.TTL),
				}
				rr.Txt = []string{entry.record.Value}

				m.Answer = append(m.Answer, rr)
			}
		}

		err := writer.WriteMsg(m)
		if err != nil {
			return
		}
	})

	srv := &dns.Server{Addr: ":7955", Net: "udp"}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
