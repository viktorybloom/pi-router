package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/viktorybloom/pi-router/internal/config"
	"github.com/viktorybloom/pi-router/internal/distro"
	"github.com/viktorybloom/pi-router/internal/firewall"
	"github.com/viktorybloom/pi-router/internal/network"
	"github.com/viktorybloom/pi-router/internal/packages"
	"github.com/viktorybloom/pi-router/internal/services"
	"github.com/viktorybloom/pi-router/internal/system"
	"github.com/viktorybloom/pi-router/internal/tailscale"
)

type App struct{ ConfigPath string }

func (a App) cfg() (config.Config, error) { return config.Load(a.ConfigPath) }

func Doctor() error {
	info := distro.Detect()
	fmt.Printf("OS: %s (%s)\n", info.Name, info.ID)
	fmt.Printf("Package manager: %s\n", info.PkgMgr)
	fmt.Printf("Service manager: %s\n", info.ServiceMgr)
	for _, bin := range []string{"hostapd", "dnsmasq", "nft", "tailscale", "dhcpcd", "iw", "ip"} {
		fmt.Printf("%-10s %v\n", bin, system.Exists(bin))
	}
	return nil
}

func (a App) Bootstrap() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	info := distro.Detect()
	return packages.Bootstrap(info.PkgMgr, cfg.InstallTailscale)
}

func (a App) Install() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(config.DefaultConfigPath), 0755); err != nil {
		return err
	}
	if err := config.Save(config.DefaultConfigPath, cfg); err != nil {
		return err
	}

	_ = system.Run("systemctl", "stop", "hostapd")
	_ = system.Run("systemctl", "stop", "dnsmasq")

	if err := os.WriteFile("/etc/sysctl.d/90-pi-router.conf", []byte("net.ipv4.ip_forward=1\nnet.ipv4.conf.all.rp_filter=1\nnet.ipv4.conf.default.rp_filter=1\nnet.ipv4.icmp_echo_ignore_broadcasts=1\nnet.ipv4.conf.all.accept_redirects=0\nnet.ipv4.conf.all.send_redirects=0\n"), 0644); err != nil {
		return err
	}
	_ = system.Run("sysctl", "--system")

	hostapd := fmt.Sprintf(`interface=wlan0
driver=nl80211
ssid=%s
hw_mode=g
channel=%s
wmm_enabled=1
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=%s
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP
country_code=%s
`, cfg.APSSID, cfg.APChannel, cfg.APPass, cfg.WifiCountry)
	if err := os.MkdirAll("/etc/hostapd", 0755); err != nil {
		return err
	}
	if err := os.WriteFile("/etc/hostapd/hostapd.conf", []byte(hostapd), 0600); err != nil {
		return err
	}
	_ = os.WriteFile("/etc/default/hostapd", []byte("DAEMON_CONF=\"/etc/hostapd/hostapd.conf\"\n"), 0644)

	_ = os.WriteFile("/etc/dnsmasq.conf", []byte("conf-dir=/etc/dnsmasq.d,.bak\n"), 0644)
	_ = os.MkdirAll("/etc/dnsmasq.d", 0755)
	dns := fmt.Sprintf(`interface=wlan0
dhcp-range=%s,%s,255.255.255.0,24h

interface=eth0
dhcp-range=%s,%s,255.255.255.0,24h

domain-needed
bogus-priv
`, cfg.APDHCPStart, cfg.APDHCPEnd, cfg.EthClientDHCPStart, cfg.EthClientDHCPEnd)
	if err := os.WriteFile("/etc/dnsmasq.d/pi-router.conf", []byte(dns), 0644); err != nil {
		return err
	}

	info := distro.Detect()
	sm := services.Manager{Kind: info.ServiceMgr}
	_ = sm.Enable("nftables")
	_ = sm.Enable("fail2ban")
	fmt.Println("Installed configs. Reboot recommended.")
	return nil
}

func (a App) WifiAP() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	info := distro.Detect()
	sm := services.Manager{Kind: info.ServiceMgr}
	_ = sm.Stop("hostapd")
	_ = sm.Stop("dnsmasq")
	if err := network.AppendManagedBlock("/etc/dhcpcd.conf", "pi-router wlan0", fmt.Sprintf("interface wlan0\nstatic ip_address=%s/24\nnohook wpa_supplicant", cfg.APIP)); err != nil {
		return err
	}
	network.ResetIface("wlan0")
	_ = sm.Restart("dhcpcd")
	_ = sm.Start("dnsmasq")
	_ = sm.Start("hostapd")
	fmt.Println("WiFi AP active:", cfg.APSSID)
	return nil
}

func (a App) EthLAN() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	if network.Uplink(cfg) == "eth0" {
		return fmt.Errorf("eth0 is uplink; cannot also be LAN")
	}
	info := distro.Detect()
	sm := services.Manager{Kind: info.ServiceMgr}
	if err := network.AppendManagedBlock("/etc/dhcpcd.conf", "pi-router eth0", fmt.Sprintf("interface eth0\nstatic ip_address=%s/24", cfg.EthClientIP)); err != nil {
		return err
	}
	network.ResetIface("eth0")
	_ = sm.Restart("dhcpcd")
	_ = sm.Restart("dnsmasq")
	fmt.Println("Ethernet LAN active on eth0")
	return nil
}

func (a App) RouteWAN() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	cfg.FailClosed = false
	return firewall.Apply(cfg, firewall.ModeWAN)
}

func (a App) RouteTunnel() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	cfg.FailClosed = true
	return firewall.Apply(cfg, firewall.ModeTunnel)
}

func (a App) TailscaleUp() error { return tailscale.Up() }
func (a App) TailscaleExit() error {
	cfg, err := a.cfg()
	if err != nil {
		return err
	}
	return tailscale.UseExit(cfg)
}

func (a App) Up() error {
	if err := a.WifiAP(); err != nil {
		return err
	}
	cfg, _ := a.cfg()
	if network.Uplink(cfg) != "eth0" {
		_ = a.EthLAN()
	}
	return a.RouteWAN()
}

func (a App) Tunnel() error {
	if err := a.TailscaleExit(); err != nil {
		return err
	}
	return a.RouteTunnel()
}

func Status() error {
	fmt.Println("===== interfaces =====")
	_ = system.Run("ip", "-br", "addr")
	fmt.Println("\n===== routes =====")
	_ = system.Run("ip", "route")
	fmt.Println("\n===== nftables =====")
	_ = system.Run("nft", "list", "ruleset")
	fmt.Println("\n===== tailscale =====")
	_ = system.Run("tailscale", "status")
	return nil
}
