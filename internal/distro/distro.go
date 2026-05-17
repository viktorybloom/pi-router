package distro

import (
	"bufio"
	"os"
	"strings"

	"github.com/viktor/pi-router/internal/system"
)

type Info struct {
	ID         string
	Name       string
	Like       string
	PkgMgr     string
	ServiceMgr string
}

func Detect() Info {
	info := Info{PkgMgr: detectPkgMgr(), ServiceMgr: detectServiceMgr()}
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return info
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := strings.Trim(parts[1], `"`)
		switch key {
		case "ID":
			info.ID = val
		case "NAME":
			info.Name = val
		case "ID_LIKE":
			info.Like = val
		}
	}
	return info
}

func detectPkgMgr() string {
	for _, pm := range []string{"apt", "pacman", "dnf", "apk", "zypper", "xbps-install"} {
		if system.Exists(pm) {
			return pm
		}
	}
	return ""
}

func detectServiceMgr() string {
	if system.Exists("systemctl") {
		return "systemd"
	}
	if system.Exists("rc-service") {
		return "openrc"
	}
	if system.Exists("service") {
		return "service"
	}
	return ""
}
