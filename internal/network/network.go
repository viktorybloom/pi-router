package network

import (
	"fmt"
	"os"
	"strings"

	"github.com/viktor/pi-router/internal/config"
	"github.com/viktor/pi-router/internal/system"
)

func Uplink(cfg config.Config) string {
	if cfg.UplinkIF != "" {
		if system.Exists("ip") {
			if _, err := system.Output("ip", "link", "show", cfg.UplinkIF); err == nil {
				return cfg.UplinkIF
			}
		}
	}
	for _, iface := range []string{"usb0", "eth1"} {
		if hasIPv4(iface) {
			return iface
		}
	}
	out, _ := system.Output("sh", "-c", "ip -o link show | awk -F': ' '{print $2}' | grep '^enx' || true")
	for _, iface := range strings.Fields(out) {
		if hasIPv4(iface) {
			return iface
		}
	}
	if hasIPv4("eth0") {
		return "eth0"
	}
	out, _ = system.Output("sh", "-c", "ip route | awk '/^default/ {print $5; exit}'")
	return strings.TrimSpace(out)
}

func hasIPv4(iface string) bool {
	out, err := system.Output("ip", "-4", "addr", "show", iface)
	return err == nil && strings.Contains(out, "inet ")
}

func RemoveManagedBlock(path, name string) error {
	b, _ := os.ReadFile(path)
	start := "# >>> " + name + " >>>"
	end := "# <<< " + name + " <<<"
	lines := strings.Split(string(b), "\n")
	out := []string{}
	skip := false
	for _, line := range lines {
		if strings.TrimSpace(line) == start {
			skip = true
			continue
		}
		if strings.TrimSpace(line) == end {
			skip = false
			continue
		}
		if !skip {
			out = append(out, line)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0644)
}

func AppendManagedBlock(path, name, content string) error {
	_ = RemoveManagedBlock(path, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\n# >>> %s >>>\n%s\n# <<< %s <<<\n", name, strings.TrimSpace(content), name)
	return err
}

func ResetIface(iface string) {
	_ = system.Run("ip", "link", "set", iface, "down")
	_ = system.Run("ip", "addr", "flush", "dev", iface)
	_ = system.Run("ip", "link", "set", iface, "up")
}
