# pi-router

A small CLI that turns a Pi or Linux device into a portable travel router.

Designed for:

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
- WAN mode can be enabled explicitly

## Supported Linux targets

- Raspberry Pi OS Lite
- Debian
- Ubuntu Server
- Arch Linux ARM
- Fedora
- openSUSE
- Alpine

The CLI detects package managers, but the networking stack still uses native Linux services:

- `hostapd`
- `dnsmasq`
- `nftables`
- `tailscale`
- `dhcpcd`
- `iw`

---

# Build

```bash
git clone https://github.com/viktorybloom/pi-router.git
cd pi-router

go build -o pi-router ./cmd/pi-router

sudo install -m 755 pi-router /usr/local/bin/pi-router
```

---

# Configure

```bash
sudo mkdir -p /usr/local/etc

sudo cp pi-router.env.example \
  /usr/local/etc/pi-router.env

sudo nano /usr/local/etc/pi-router.env
```

Set a long random WiFi password.

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

# Usage

## Normal WAN routing

```bash
sudo pi-router up
```

Starts:

- WiFi AP
- Ethernet client LAN if `eth0` is free
- normal WAN routing through detected uplink

---

## Tailscale tunnel routing

On your home node:

```bash
sudo tailscale up \
  --advertise-exit-node \
  --advertise-routes=192.168.1.0/24 \
  --ssh=false
```

Approve:
- exit node
- subnet route

inside the Tailscale admin console.

On the travel Pi, set:

```ini
HOME_EXIT_NODE=your-home-node-name
```

Then:

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
pi-router doctor
pi-router bootstrap
pi-router install

pi-router up
pi-router wifi-ap
pi-router eth-lan

pi-router route-wan
pi-router route-tailscale

pi-router tailscale-up
pi-router tailscale-exit
pi-router tunnel

pi-router status
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
