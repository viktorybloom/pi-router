package main

import (
	"fmt"
	"os"

	"github.com/viktorybloom/pi-router/internal/commands"
	"github.com/viktorybloom/pi-router/internal/config"
	"github.com/viktorybloom/pi-router/internal/system"
)

func usage() {
	fmt.Println(`pi-router - portable Linux/Tailscale travel router

Usage:
  pi-router [--config path] <command>

Commands:
  doctor            check distro, package manager, and required tools
  bootstrap         install OS dependencies using detected package manager
  install           write hostapd/dnsmasq/sysctl configs
  up                start WiFi AP, optional eth LAN, and normal WAN routing
  wifi-ap           start WiFi AP only
  eth-lan           make eth0 a client LAN if it is not uplink
  route-wan         allow clients through detected WAN/uplink
  route-tailscale   fail-closed: allow clients through tailscale0 only
  tailscale-up      tailscale up --ssh=false
  tailscale-exit    use HOME_EXIT_NODE from config
  tunnel            tailscale-exit + route-tailscale
  status            show interfaces, routes, firewall, tailscale

Config default: /usr/local/etc/pi-router.env
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cfgPath := config.DefaultConfigPath
	args := os.Args[1:]
	if len(args) >= 2 && args[0] == "--config" {
		cfgPath = args[1]
		args = args[2:]
	}
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	cmd := args[0]
	if cmd != "doctor" && cmd != "status" {
		if err := system.MustRoot(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	app := commands.App{ConfigPath: cfgPath}
	var err error
	switch cmd {
	case "doctor":
		err = commands.Doctor()
	case "bootstrap":
		err = app.Bootstrap()
	case "install":
		err = app.Install()
	case "up":
		err = app.Up()
	case "wifi-ap":
		err = app.WifiAP()
	case "eth-lan":
		err = app.EthLAN()
	case "route-wan":
		err = app.RouteWAN()
	case "route-tailscale":
		err = app.RouteTunnel()
	case "tailscale-up":
		err = app.TailscaleUp()
	case "tailscale-exit":
		err = app.TailscaleExit()
	case "tunnel":
		err = app.Tunnel()
	case "status":
		err = commands.Status()
	default:
		usage()
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
