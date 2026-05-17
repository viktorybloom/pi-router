package tailscale

import (
	"fmt"
	"github.com/viktorybloom/pi-router/internal/config"
	"github.com/viktorybloom/pi-router/internal/system"
)

func Up() error { return system.Run("tailscale", "up", "--ssh=false") }

func UseExit(cfg config.Config) error {
	if cfg.HomeExitNode == "" {
		return fmt.Errorf("HOME_EXIT_NODE is blank")
	}
	return system.Run("tailscale", "up", "--exit-node="+cfg.HomeExitNode, "--exit-node-allow-lan-access", "--ssh=false")
}
