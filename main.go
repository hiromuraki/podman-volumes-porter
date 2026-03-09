package main

import (
	"context"
	"fmt"
	"podman-volumes-porter/internal/core"
	"time"
)

func main() {
	s3storage := core.S3Storage{
		Endpoint:  "http://localhost:8333",
		AccessKey: "MySeaweedAccessKey",
		SecretKey: "MySeaweedSecretKey123",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	err := core.RestoreVolume(ctx, "seaweed-config-2", "seaweed-config.tar.gz", s3storage)
	// err := core.BackupVolume(ctx, "seaweed-config", &s3storage)
	if err != nil {
		fmt.Print(err.Error())
	}
}
