# pi-route

A small Go CLI that turns a pi, into a portable travel router.

It is designed for this topology:

```text
Internet uplink:
  phone USB tether, Ethernet, USB Ethernet, hotel/cafe network

Pi:
  WiFi AP for clients
  optional eth0 client LAN when eth0 is not uplink
  Tailscale node
  nftables firewall

Clients:
  connect by WiFi or Ethernet
  no Tailscale client required
```

## Security model

- No public inbound ports required.
- SSH should be key-only.
- SSH is allowed through Tailscale only by default firewall rules.
- Tunnel mode is fail-closed: clients route through `tailscale0` only.
- WAN mode can be enabled explicitly for normal internet routing.

## Supported Linux targets

- Raspberry Pi OS Lite
- Debian
- Ubuntu Server
- Arch Linux ARM

- Fedora
- openSUSE
- Alpine

The CLI detects package managers, but the actual networking stack still requires host services:

- hostapd
- dnsmasq
- nftables
- tailscale
- dhcpcd
- iw

## Build from source

```bash
git clone https://github.com/viktorybloom/pi-router.git
cd pi-router
go build -o pi-router ./cmd/pi-router
sudo install -m 755 pi-router /usr/local/bin/pi-router
```

## Config

```bash
sudo mkdir -p /usr/local/etc
sudo cp pi-router.env.example /usr/local/etc/pi-router.env
sudo nano /usr/local/etc/pi-router.env
```

Set a long random WiFi password.

## Install on a fresh Pi

```bash
sudo pi-router doctor
sudo pi-router bootstrap
sudo pi-router install
sudo reboot
```

After reboot:

```bash
sudo pi-router tailscale-up
sudo pi-router up
```

## Normal WAN routing

```bash
sudo pi-router up
```

This starts:

- WiFi AP
- Ethernet client LAN if eth0 is free
- normal WAN routing through detected uplink

## Tailscale tunnel routing

On your home node first:

```bash
sudo tailscale up \
  --advertise-exit-node \
  --advertise-routes=192.168.1.0/24 \
  --ssh=false
```

Approve the exit node and subnet route in the Tailscale admin console.

On the travel Pi, set in `/usr/local/etc/pi-router.env`:

```ini
HOME_EXIT_NODE=your-home-node-name
```

Then run:

```bash
sudo pi-router tailscale-exit
sudo pi-router route-tailscale
```

Or:

```bash
sudo pi-router tunnel
```

## Commands

```text
pi-router doctor            check distro, package manager, and required tools
pi-router bootstrap         install OS dependencies using detected package manager
pi-router install           write hostapd/dnsmasq/sysctl configs
pi-router up                start AP, optional eth LAN, and WAN routing
pi-router wifi-ap           start WiFi AP only
pi-router eth-lan           make eth0 a client LAN if not uplink
pi-router route-wan         allow clients through detected WAN/uplink
pi-router route-tailscale   fail-closed client routing via tailscale0 only
pi-router tailscale-up      tailscale up --ssh=false
pi-router tailscale-exit    use HOME_EXIT_NODE
pi-router tunnel            tailscale-exit + route-tailscale
pi-router status            show interfaces, routes, firewall, tailscale
```

## Notes

This is a host-level networking tool. Do not run the router core in Docker. It needs to control interfaces, services, routing, and firewall state.

