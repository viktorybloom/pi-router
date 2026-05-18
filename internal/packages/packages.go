package packages

import (
	"fmt"

	"github.com/viktorybloom/pi-router/internal/system"
)

type Set struct{ Hostapd, Dnsmasq, Nftables, Curl, DHCPCD, IW, WirelessTools, Unattended, Fail2ban, Tailscale string }

var packageMap = map[string]Set{
	"apt":    {"hostapd", "dnsmasq", "nftables", "curl", "dhcpcd5", "iw", "wireless-tools", "unattended-upgrades", "fail2ban", "tailscale"},
	"pacman": {"hostapd", "dnsmasq", "nftables", "curl", "dhcpcd", "iw", "wireless_tools", "", "fail2ban", "tailscale"},
	"dnf":    {"hostapd", "dnsmasq", "nftables", "curl", "dhcpcd", "iw", "wireless-tools", "", "fail2ban", "tailscale"},
	"zypper": {"hostapd", "dnsmasq", "nftables", "curl", "dhcpcd", "iw", "wireless-tools", "", "fail2ban", "tailscale"},
	"apk":    {"hostapd", "dnsmasq", "nftables", "curl", "dhcpcd", "iw", "wireless-tools", "", "fail2ban", "tailscale"},
}

func Bootstrap(pm string, installTailscale bool) error {
	pkgs, ok := packageMap[pm]
	if !ok || pm == "" {
		return fmt.Errorf("unsupported or unknown package manager: %s", pm)
	}
	list := []string{pkgs.Hostapd, pkgs.Dnsmasq, pkgs.Nftables, pkgs.Curl, pkgs.DHCPCD, pkgs.IW, pkgs.WirelessTools, pkgs.Unattended, pkgs.Fail2ban}
	if installTailscale {
		list = append(list, pkgs.Tailscale)
	}
	filtered := []string{}
	for _, p := range list {
		if p != "" {
			filtered = append(filtered, p)
		}
	}

	switch pm {
	case "apt":
		if err := system.Run("apt", "update"); err != nil {
			return err
		}
		return system.Run("apt", append([]string{"install", "-y"}, filtered...)...)
	case "pacman":
		return system.Run("pacman", append([]string{"-Syu", "--needed", "--noconfirm"}, filtered...)...)
	case "dnf":
		return system.Run("dnf", append([]string{"install", "-y"}, filtered...)...)
	case "zypper":
		return system.Run("zypper", append([]string{"--non-interactive", "install"}, filtered...)...)
	case "apk":
		if err := system.Run("apk", "update"); err != nil {
			return err
		}
		return system.Run("apk", append([]string{"add"}, filtered...)...)
	default:
		return fmt.Errorf("no installer implementation for %s", pm)
	}
}
