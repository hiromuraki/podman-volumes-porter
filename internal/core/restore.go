package core

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"os/exec"
)

func RestoreVolume(ctx context.Context, volName, key string, s3storage S3Storage) error {
	fmt.Printf("📥 正在从 s3://%s/%s 获取流式数据...\n", "container-volume", key)
	objReader, err := s3storage.DownloadStream(ctx, "container-volume", key)
	if err != nil {
		return err
	}

	gzReader, err := gzip.NewReader(objReader)
	if err != nil {
		return fmt.Errorf("初始化 Gzip 解压器失败: %w", err)
	}
	defer gzReader.Close()

	cmd := exec.CommandContext(ctx, "podman", "volume", "import", volName, "-")
	cmd.Stdin = gzReader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("正在拉起 Podman 并注入数据到卷 [%s]...\n", volName)

	if !VolumeExists(ctx, volName) {
		fmt.Printf("⚠️ 卷 %s 不存在，正在尝试创建...\n", volName)
		err := exec.CommandContext(ctx, "podman", "volume", "create", volName).Run()
		if err != nil {
			return fmt.Errorf("无法创建卷 %s", volName)
		}
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Podman 卷导入失败: %w", err)
	}

	fmt.Println("卷恢复成功")
	return nil
}
