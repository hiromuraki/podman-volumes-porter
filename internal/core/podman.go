package core

import (
	"context"
	"os/exec"
)

func VolumeExists(ctx context.Context, volumeName string) bool {
	cmd := exec.CommandContext(ctx, "podman", "volume", "import", volumeName, "-")
	err := cmd.Run()
	return err == nil
}
