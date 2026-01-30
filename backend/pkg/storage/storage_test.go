package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLocalStorage_Upload(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "storage_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建本地存储实例
	storage, err := NewLocalStorage("test-local", tmpDir)
	if err != nil {
		t.Fatalf("创建本地存储失败: %v", err)
	}

	ctx := context.Background()

	// 测试上传
	content := []byte("Hello, World!")
	reader := bytes.NewReader(content)

	err = storage.Upload(ctx, "test.txt", reader, int64(len(content)), nil)
	if err != nil {
		t.Fatalf("上传文件失败: %v", err)
	}

	// 验证文件存在
	exists, err := storage.Exists(ctx, "test.txt")
	if err != nil {
		t.Fatalf("检查文件存在失败: %v", err)
	}
	if !exists {
		t.Fatal("文件应该存在")
	}

	// 验证文件内容
	filePath := filepath.Join(tmpDir, "test.txt")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("文件内容不匹配: got %s, want %s", data, content)
	}

	t.Log("本地存储上传测试通过")
}

func TestS3Storage_Upload(t *testing.T) {
	// 从环境变量获取S3配置，如果没有配置则跳过测试
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	// 如果没有配置环境变量，使用默认的RustFS测试配置
	if endpoint == "" {
		endpoint = "http://192.168.10.210:9000"
		accessKey = "minioadmin"
		secretKey = "minioadmin"
		bucket = "test-bucket"
	}

	// 创建S3存储实例
	storage, err := NewS3Storage("test-s3", "s3", endpoint, accessKey, secretKey, bucket, "")
	if err != nil {
		t.Skipf("跳过S3测试，无法连接到S3服务: %v", err)
	}

	ctx := context.Background()

	// 生成唯一的测试文件名
	testFileName := fmt.Sprintf("test_%d.txt", time.Now().UnixNano())
	content := []byte("Hello, S3 Storage!")
	reader := bytes.NewReader(content)

	// 测试上传
	err = storage.Upload(ctx, testFileName, reader, int64(len(content)), nil)
	if err != nil {
		t.Skipf("跳过S3测试，上传失败（可能S3服务不可用）: %v", err)
	}

	// 验证文件存在
	exists, err := storage.Exists(ctx, testFileName)
	if err != nil {
		t.Fatalf("检查文件存在失败: %v", err)
	}
	if !exists {
		t.Fatal("文件应该存在")
	}

	// 验证文件内容
	downloadReader, fileInfo, err := storage.Download(ctx, testFileName)
	if err != nil {
		t.Fatalf("下载文件失败: %v", err)
	}
	defer downloadReader.Close()

	downloadedContent, err := io.ReadAll(downloadReader)
	if err != nil {
		t.Fatalf("读取下载内容失败: %v", err)
	}

	if string(downloadedContent) != string(content) {
		t.Fatalf("文件内容不匹配: got %s, want %s", downloadedContent, content)
	}

	t.Logf("文件信息: name=%s, size=%d", fileInfo.Name, fileInfo.Size)

	// 清理：删除测试文件
	err = storage.Delete(ctx, testFileName)
	if err != nil {
		t.Logf("警告：删除测试文件失败: %v", err)
	}

	t.Log("S3存储上传测试通过")
}
