package core

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os/exec"
)

func BackupVolume(ctx context.Context, volName string, s3storage *S3Storage) error {
	// 内存管道逻辑
	pr, pw := io.Pipe()
	go func() {
		gw := gzip.NewWriter(pw)

		cmd := exec.CommandContext(ctx, "podman", "volume", "export", volName)
		cmd.Stdout = gw

		err := cmd.Run()
		gw.Close()

		if err != nil {
			pw.CloseWithError(err)
		} else {
			pw.Close()
		}
	}()

	err := s3storage.UploadStream(ctx, "container-volume", volName+".tar.gz", pr)
	if err != nil {
		return fmt.Errorf("传输失败: %w", err)
	}

	fmt.Println("卷备份成功")
	return nil
}
