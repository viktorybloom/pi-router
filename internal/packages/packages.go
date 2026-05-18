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

		if err := system.Run(
			"apt",
			append([]string{"install", "-y"}, filtered...)...,
		); err != nil {
			return err
		}

	case "pacman":
		if err := system.Run(
			"pacman",
			append([]string{"-Syu", "--needed", "--noconfirm"}, filtered...)...,
		); err != nil {
			return err
		}

	case "dnf":
		if err := system.Run(
			"dnf",
			append([]string{"install", "-y"}, filtered...)...,
		); err != nil {
			return err
		}

	case "zypper":
		if err := system.Run(
			"zypper",
			append([]string{"--non-interactive", "install"}, filtered...)...,
		); err != nil {
			return err
		}

	case "apk":
		if err := system.Run("apk", "update"); err != nil {
			return err
		}

		if err := system.Run(
			"apk",
			append([]string{"add"}, filtered...)...,
		); err != nil {
			return err
		}

	default:
		return fmt.Errorf("no installer implementation for %s", pm)
	}

	if installTailscale {
		return InstallTailscale()
	}

	return nil
}

func InstallTailscale() error {
	if system.CommandExists("tailscale") {
		return nil
	}

	return system.Run(
		"sh",
		"-c",
		"curl -fsSL https://tailscale.com/install.sh | sh",
	)
}
