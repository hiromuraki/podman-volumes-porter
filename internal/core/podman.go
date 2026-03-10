package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Podman struct{}

var podman Podman

func (p Podman) volumeExists(ctx context.Context, volumeName string) bool {
	cmd := exec.CommandContext(ctx, "podman", "volume", "import", volumeName, "-")
	err := cmd.Run()
	return err == nil
}

func (p Podman) importVolume(ctx context.Context, volumeName string, source io.Reader) error {
	if err := exec.CommandContext(ctx, "podman", "volume", "create", volumeName).Run(); err != nil {
		return fmt.Errorf("无法创建卷 %s", volumeName)
	}

	cmd := exec.CommandContext(ctx, "podman", "volume", "import", volumeName, "-")
	cmd.Stdin = source
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p Podman) deleteVolume(ctx context.Context, volumeName string) error {
	return exec.CommandContext(ctx, "podman", "volume", "rm", "-f", volumeName).Run()
}

func (p Podman) getAllVolumeNames() []string {
	cmd := exec.Command("podman", "volume", "ls", "--format", "{{.Name}}")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")

	var volumeNames []string
	for _, line := range lines {
		if name := strings.TrimSpace(line); name != "" {
			volumeNames = append(volumeNames, name)
		}
	}
	return volumeNames
}
