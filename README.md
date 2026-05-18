# pi-router

A small CLI that turns a Pi or Linux device into a portable travel router.

It is designed for this topology:

```text
Internet uplink:
  phone USB tether
  Ethernet
  USB Ethernet
  hotel/cafe network

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

- no public inbound ports required
- SSH should be key-only
- SSH allowed through Tailscale only by default firewall rules
- tunnel mode is fail-closed:
  clients route through `tailscale0` only
- WAN mode can be enabled explicitly for normal internet routing

## Supported Linux targets

- Raspberry Pi OS Lite
- Debian
- Ubuntu Server
- Arch Linux ARM
- Fedora
- openSUSE
- Alpine

The CLI detects package managers, but the networking stack still requires host services:

- `hostapd`
- `dnsmasq`
- `nftables`
- `tailscale`
- `dhcpcd`
- `iw`

---

# Build from source

```bash
git clone https://github.com/viktorybloom/pi-router.git
cd pi-router

go build -o pi-router ./cmd/pi-router

sudo install -m 755 pi-router /usr/local/bin/pi-router
```

---

# Config

```bash
sudo mkdir -p /usr/local/etc

sudo cp pi-router.env.example \
  /usr/local/etc/pi-router.env

sudo nano /usr/local/etc/pi-router.env
```

Set:
- a long random WiFi password
- optional `HOME_EXIT_NODE`

Example:

```ini
AP_SSID=pi_travel_router
AP_PASS=CHANGE_ME_LONG_RANDOM_PASSWORD

HOME_EXIT_NODE=your-home-node-name
```

---

# Install on a fresh Pi

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

---

# Normal WAN routing

```bash
sudo pi-router up
```

This starts:

- WiFi AP
- Ethernet client LAN if `eth0` is free
- normal WAN routing through detected uplink

---

# Tailscale tunnel routing

On your home node first:

```bash
sudo tailscale up \
  --advertise-exit-node \
  --advertise-routes=192.168.1.0/24 \
  --ssh=false
```

Approve the:
- exit node
- subnet route

inside the Tailscale admin console.

Then on the travel Pi:

```bash
sudo pi-router tailscale-exit
sudo pi-router route-tailscale
```

Or:

```bash
sudo pi-router tunnel
```

---

# Commands

```text
pi-router doctor            check distro, package manager, and required tools
pi-router bootstrap         install OS dependencies
pi-router install           write hostapd/dnsmasq/sysctl configs

pi-router up                start AP, optional eth LAN, and WAN routing
pi-router wifi-ap           start WiFi AP only
pi-router eth-lan           enable eth0 client LAN

pi-router route-wan         allow client internet through detected uplink
pi-router route-tailscale   fail-closed routing through tailscale0 only

pi-router tailscale-up      run tailscale up --ssh=false
pi-router tailscale-exit    connect to HOME_EXIT_NODE
pi-router tunnel            tailscale-exit + route-tailscale

pi-router status            show interfaces, routes, firewall, and tailscale state

pi-router uninstall         remove generated configs/services
pi-router purge             uninstall and remove runtime config
```

---

# Uninstall

Remove generated configs/services:

```bash
sudo pi-router uninstall
```

Remove config too:

```bash
sudo pi-router purge
```

Remove binary manually:

```bash
sudo rm -f /usr/local/bin/pi-router
```
