package packages

import (
	"fmt"

	"github.com/viktorybloom/pi-router/internal/system"
)

type Set struct {
	Hostapd       string
	Dnsmasq       string
	Nftables      string
	Curl          string
	DHCPCD        string
	IW            string
	WirelessTools string
	Unattended    string
	Fail2ban      string
	Tailscale     string
}

var packageMap = map[string]Set{
	"apt": {
		Hostapd:       "hostapd",
		Dnsmasq:       "dnsmasq",
		Nftables:      "nftables",
		Curl:          "curl",
		DHCPCD:        "dhcpcd5",
		IW:            "iw",
		WirelessTools: "wireless-tools",
		Unattended:    "unattended-upgrades",
		Fail2ban:      "fail2ban",
		Tailscale:     "",
	},
	"pacman": {
		Hostapd:       "hostapd",
		Dnsmasq:       "dnsmasq",
		Nftables:      "nftables",
		Curl:          "curl",
		DHCPCD:        "dhcpcd",
		IW:            "iw",
		WirelessTools: "wireless_tools",
		Unattended:    "",
		Fail2ban:      "fail2ban",
		Tailscale:     "tailscale",
	},
	"dnf": {
		Hostapd:       "hostapd",
		Dnsmasq:       "dnsmasq",
		Nftables:      "nftables",
		Curl:          "curl",
		DHCPCD:        "dhcpcd",
		IW:            "iw",
		WirelessTools: "wireless-tools",
		Unattended:    "",
		Fail2ban:      "fail2ban",
		Tailscale:     "tailscale",
	},
	"zypper": {
		Hostapd:       "hostapd",
		Dnsmasq:       "dnsmasq",
		Nftables:      "nftables",
		Curl:          "curl",
		DHCPCD:        "dhcpcd",
		IW:            "iw",
		WirelessTools: "wireless-tools",
		Unattended:    "",
		Fail2ban:      "fail2ban",
		Tailscale:     "tailscale",
	},
	"apk": {
		Hostapd:       "hostapd",
		Dnsmasq:       "dnsmasq",
		Nftables:      "nftables",
		Curl:          "curl",
		DHCPCD:        "dhcpcd",
		IW:            "iw",
		WirelessTools: "wireless-tools",
		Unattended:    "",
		Fail2ban:      "fail2ban",
		Tailscale:     "tailscale",
	},
}

func Bootstrap(pm string, installTailscale bool) error {
	pkgs, ok := packageMap[pm]
	if !ok || pm == "" {
		return fmt.Errorf("unsupported or unknown package manager: %s", pm)
	}

	list := []string{
		pkgs.Hostapd,
		pkgs.Dnsmasq,
		pkgs.Nftables,
		pkgs.Curl,
		pkgs.DHCPCD,
		pkgs.IW,
		pkgs.WirelessTools,
		pkgs.Unattended,
		pkgs.Fail2ban,
	}

	if installTailscale && pkgs.Tailscale != "" {
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
		if err := system.Run("apt", append([]string{"install", "-y"}, filtered...)...); err != nil {
			return err
		}
		if installTailscale {
			return InstallTailscaleApt()
		}
		return nil

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

func InstallTailscaleApt() error {
	if system.CommandExists("tailscale") {
		return nil
	}

	// Safer long-term TODO:
	// replace this with Tailscale's distro repo/key setup.
	// This is pragmatic and matches Tailscale's official installer path.
	return system.Run("sh", "-c", "curl -fsSL https://tailscale.com/install.sh | sh")
}
