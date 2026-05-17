package firewall

import (
	"fmt"
	"os"

	"github.com/viktorybloom/pi-router/internal/config"
	"github.com/viktorybloom/pi-router/internal/network"
	"github.com/viktorybloom/pi-router/internal/system"
)

type Mode string

const (
	ModeWAN    Mode = "wan"
	ModeTunnel Mode = "tunnel"
)

func Apply(cfg config.Config, mode Mode) error {
	wan := network.Uplink(cfg)
	failClosed := cfg.FailClosed || mode == ModeTunnel
	natOut := "tailscale0"
	forwardExtra := ""
	natExtra := ""
	if !failClosed && wan != "" {
		forwardExtra = fmt.Sprintf(`    iif "wlan0" oif "%s" accept
    iif "eth0" oif "%s" accept`, wan, wan)
		natExtra = fmt.Sprintf(`    oif "%s" masquerade`, wan)
		natOut = wan
	}
	_ = natOut
	sshClient := ""
	if cfg.AllowClientSSH {
		sshClient = `    iif "wlan0" tcp dport 22 accept
    iif "eth0" tcp dport 22 accept`
	}
	content := fmt.Sprintf(`flush ruleset

table inet filter {
  chain input {
    type filter hook input priority 0;
    policy drop;

    iif "lo" accept
    ct state established,related accept
    iif "tailscale0" accept
    ip protocol icmp accept
    ip6 nexthdr ipv6-icmp accept

    iif "wlan0" udp dport {53,67,68} accept
    iif "wlan0" tcp dport 53 accept
    iif "eth0" udp dport {53,67,68} accept
    iif "eth0" tcp dport 53 accept
    iif "tailscale0" tcp dport 22 accept
%s
  }

  chain forward {
    type filter hook forward priority 0;
    policy drop;

    ct state established,related accept
    iif "wlan0" oif "tailscale0" accept
    iif "eth0" oif "tailscale0" accept
%s
  }

  chain output {
    type filter hook output priority 0;
    policy accept;
  }
}

table ip nat {
  chain postrouting {
    type nat hook postrouting priority 100;
    oif "tailscale0" masquerade
%s
  }
}
`, sshClient, forwardExtra, natExtra)
	if err := os.WriteFile("/etc/nftables.conf", []byte(content), 0644); err != nil {
		return err
	}
	return system.Run("systemctl", "restart", "nftables")
}
