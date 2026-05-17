package services

import "github.com/viktorybloom/pi-router/internal/system"

type Manager struct{ Kind string }

func (m Manager) Start(name string) error   { return m.exec("start", name) }
func (m Manager) Stop(name string) error    { return m.exec("stop", name) }
func (m Manager) Restart(name string) error { return m.exec("restart", name) }
func (m Manager) Enable(name string) error  { return m.exec("enable", name) }

func (m Manager) exec(action, name string) error {
	switch m.Kind {
	case "systemd":
		return system.Run("systemctl", action, name)
	case "openrc":
		if action == "enable" {
			return system.Run("rc-update", "add", name, "default")
		}
		return system.Run("rc-service", name, action)
	default:
		if action == "enable" {
			return nil
		}
		return system.Run("service", name, action)
	}
}
