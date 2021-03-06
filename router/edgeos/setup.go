package edgeos

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nextdns/nextdns/config"
	"github.com/nextdns/nextdns/router/internal"
)

type Router struct {
	DNSMasqPath     string
	ListenPort      string
	ClientReporting bool
}

func New() (*Router, bool) {
	if st, err := os.Stat("/config/scripts/post-config.d"); err != nil || !st.IsDir() {
		return nil, false
	}
	return &Router{
		DNSMasqPath: "/etc/dnsmasq.d/nextdns.conf",
		ListenPort:  "5342",
	}, true
}

func (r *Router) Configure(c *config.Config) {
	c.Listen = "127.0.0.1:" + r.ListenPort
	r.ClientReporting = c.ReportClientInfo
}

func (r *Router) Setup() error {
	if err := internal.WriteTemplate(r.DNSMasqPath, tmpl, r, 0644); err != nil {
		return err
	}

	return restartDNSMasq()
}

func (r *Router) Restore() error {
	if err := os.Remove(r.DNSMasqPath); err != nil {
		return err
	}
	return restartDNSMasq()
}

func restartDNSMasq() error {
	if err := exec.Command("sudo", "/etc/init.d/dnsmasq", "restart").Run(); err != nil {
		return fmt.Errorf("dnsmasq restart: %v", err)
	}
	return nil
}

var tmpl = `# Configuration generated by NextDNS
no-resolv
server=127.0.0.1#{{.ListenPort}}
{{- if .ClientReporting}}
add-mac
add-subnet=32,128
{{- end}}
`
