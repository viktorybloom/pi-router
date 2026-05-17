package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const DefaultConfigPath = "/usr/local/etc/pi-router.env"

type Config struct {
	APSSID             string
	APPass             string
	APChannel          string
	APIP               string
	APDHCPStart        string
	APDHCPEnd          string
	EthClientIP        string
	EthClientDHCPStart string
	EthClientDHCPEnd   string
	WifiCountry        string
	InstallTailscale   bool
	UplinkIF           string
	HomeExitNode       string
	AllowClientSSH     bool
	FailClosed         bool
}

func Default() Config {
	return Config{
		APSSID:             "pi_travel_router",
		APPass:             "CHANGE_ME_LONG_RANDOM_PASSWORD",
		APChannel:          "6",
		APIP:               "192.168.50.1",
		APDHCPStart:        "192.168.50.10",
		APDHCPEnd:          "192.168.50.100",
		EthClientIP:        "192.168.60.1",
		EthClientDHCPStart: "192.168.60.10",
		EthClientDHCPEnd:   "192.168.60.100",
		WifiCountry:        "AU",
		InstallTailscale:   true,
		AllowClientSSH:     false,
		FailClosed:         true,
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		path = DefaultConfigPath
	}
	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	values := map[string]string{}
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		values[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	if err := s.Err(); err != nil {
		return cfg, err
	}

	get := func(k, fallback string) string {
		if v, ok := values[k]; ok {
			return v
		}
		return fallback
	}
	boolv := func(k string, fallback bool) bool {
		v, ok := values[k]
		if !ok {
			return fallback
		}
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "1", "yes", "y", "on":
			return true
		default:
			return false
		}
	}

	cfg.APSSID = get("AP_SSID", cfg.APSSID)
	cfg.APPass = get("AP_PASS", cfg.APPass)
	cfg.APChannel = get("AP_CHANNEL", cfg.APChannel)
	cfg.APIP = get("AP_IP", cfg.APIP)
	cfg.APDHCPStart = get("AP_DHCP_START", cfg.APDHCPStart)
	cfg.APDHCPEnd = get("AP_DHCP_END", cfg.APDHCPEnd)
	cfg.EthClientIP = get("ETH_CLIENT_IP", cfg.EthClientIP)
	cfg.EthClientDHCPStart = get("ETH_CLIENT_DHCP_START", cfg.EthClientDHCPStart)
	cfg.EthClientDHCPEnd = get("ETH_CLIENT_DHCP_END", cfg.EthClientDHCPEnd)
	cfg.WifiCountry = get("WIFI_COUNTRY", cfg.WifiCountry)
	cfg.InstallTailscale = boolv("INSTALL_TAILSCALE", cfg.InstallTailscale)
	cfg.UplinkIF = get("UPLINK_IF", cfg.UplinkIF)
	cfg.HomeExitNode = get("HOME_EXIT_NODE", cfg.HomeExitNode)
	cfg.AllowClientSSH = boolv("ALLOW_CLIENT_SSH", cfg.AllowClientSSH)
	cfg.FailClosed = boolv("FAIL_CLOSED", cfg.FailClosed)
	return cfg, nil
}

func Save(path string, cfg Config) error {
	if path == "" {
		path = DefaultConfigPath
	}
	content := fmt.Sprintf(`AP_SSID=%s
AP_PASS=%s
AP_CHANNEL=%s
AP_IP=%s
AP_DHCP_START=%s
AP_DHCP_END=%s
ETH_CLIENT_IP=%s
ETH_CLIENT_DHCP_START=%s
ETH_CLIENT_DHCP_END=%s
WIFI_COUNTRY=%s
INSTALL_TAILSCALE=%t
UPLINK_IF=%s
HOME_EXIT_NODE=%s
ALLOW_CLIENT_SSH=%t
FAIL_CLOSED=%t
`, cfg.APSSID, cfg.APPass, cfg.APChannel, cfg.APIP, cfg.APDHCPStart, cfg.APDHCPEnd, cfg.EthClientIP, cfg.EthClientDHCPStart, cfg.EthClientDHCPEnd, cfg.WifiCountry, cfg.InstallTailscale, cfg.UplinkIF, cfg.HomeExitNode, cfg.AllowClientSSH, cfg.FailClosed)
	return os.WriteFile(path, []byte(content), 0600)
}
