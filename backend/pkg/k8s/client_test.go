package k8s

import (
	"testing"
	"time"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

// TestConfigValidate 测试配置验证
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid in-cluster config",
			config: &Config{
				InCluster: true,
				Namespace: "default",
				Timeout:   30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid kubeconfig",
			config: &Config{
				InCluster:  false,
				KubeConfig: "/tmp/test-kubeconfig",
				Namespace:  "default",
				Timeout:    30 * time.Second,
			},
			wantErr: true, // 文件不存在
			errMsg:  "kubeconfig file not found",
		},
		{
			name: "missing kubeconfig path",
			config: &Config{
				InCluster: false,
				Namespace: "default",
				Timeout:   30 * time.Second,
			},
			wantErr: true,
			errMsg:  "kubeconfig path is required",
		},
		{
			name: "empty namespace defaults to default",
			config: &Config{
				InCluster: true,
				Namespace: "",
				Timeout:   30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "zero timeout defaults to 30s",
			config: &Config{
				InCluster: true,
				Namespace: "default",
				Timeout:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				// 验证默认值
				if tt.config.Namespace == "" {
					assert.Equal(t, "default", tt.config.Namespace)
				}
				if tt.config.Timeout == 0 {
					assert.Equal(t, 30*time.Second, tt.config.Timeout)
				}
			}
		})
	}
}

// TestWrapError 测试错误包装
func TestWrapError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
		want    string
	}{
		{
			name:    "wrap error",
			err:     ErrPodNotFound,
			message: "failed to get pod",
			want:    "failed to get pod: pod not found",
		},
		{
			name:    "nil error",
			err:     nil,
			message: "some message",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.message)
			if tt.err == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Contains(t, result.Error(), tt.want)
			}
		})
	}
}

// TestWrapErrorf 测试格式化错误包装
func TestWrapErrorf(t *testing.T) {
	// 测试正常情况
	err := ErrPodCreationFailed
	result := WrapErrorf(err, "failed to create pod %s in namespace %s", "test-pod", "default")

	assert.NotNil(t, result)
	assert.Contains(t, result.Error(), "failed to create pod test-pod in namespace default")
	assert.Contains(t, result.Error(), "pod creation failed")

	// 测试nil error
	result = WrapErrorf(nil, "some message %s", "arg")
	assert.Nil(t, result)
}

// TestIsPodNotFound 测试Pod不存在错误判断
func TestIsPodNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is pod not found",
			err:  ErrPodNotFound,
			want: true,
		},
		{
			name: "wrapped pod not found",
			err:  WrapError(ErrPodNotFound, "failed to get pod"),
			want: true,
		},
		{
			name: "other error",
			err:  ErrConnectionFailed,
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPodNotFound(tt.err)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestIsConnectionFailed 测试连接失败错误判断
func TestIsConnectionFailed(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is connection failed",
			err:  ErrConnectionFailed,
			want: true,
		},
		{
			name: "wrapped connection failed",
			err:  WrapError(ErrConnectionFailed, "failed to connect"),
			want: true,
		},
		{
			name: "other error",
			err:  ErrPodNotFound,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsConnectionFailed(tt.err)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestConfigString 测试配置字符串表示
func TestConfigString(t *testing.T) {
	cfg := &Config{
		Namespace: "test-ns",
		InCluster: true,
		Timeout:   30 * time.Second,
	}

	result := cfg.String()
	assert.Contains(t, result, "test-ns")
	assert.Contains(t, result, "true")
	assert.Contains(t, result, "30s")
}

// TestNewClientWithInvalidConfig 测试使用无效配置创建客户端
func TestNewClientWithInvalidConfig(t *testing.T) {
	// 测试nil配置
	client, err := NewClient(nil)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "config is nil")

	// 测试无效配置（缺少kubeconfig路径）
	cfg := &Config{
		InCluster: false,
		Namespace: "default",
		Timeout:   30 * time.Second,
	}
	client, err = NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

// TestClientMethods 测试Client的基本方法
func TestClientMethods(t *testing.T) {
	// 创建fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// 使用fake clientset创建client
	client := NewClientWithClientset(fakeClientset, "test-namespace")
	defer client.Close()

	// 测试GetClientset
	assert.NotNil(t, client.GetClientset())
	assert.Equal(t, fakeClientset, client.GetClientset())

	// 测试GetNamespace
	assert.Equal(t, "test-namespace", client.GetNamespace())

	// 测试GetContext
	assert.NotNil(t, client.GetContext())

	// 测试GetContextWithTimeout
	ctx, cancel := client.GetContextWithTimeout(5 * time.Second)
	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)
	cancel()
}

// TestClientPing 测试Ping方法
func TestClientPing(t *testing.T) {
	// 创建fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// 使用fake clientset创建client
	client := NewClientWithClientset(fakeClientset, "default")
	defer client.Close()

	// 测试Ping（fake clientset应该能正常返回）
	err := client.Ping()
	assert.NoError(t, err)
}

// TestClientClose 测试Close方法
func TestClientClose(t *testing.T) {
	// 创建fake clientset
	fakeClientset := fake.NewSimpleClientset()

	// 使用fake clientset创建client
	client := NewClientWithClientset(fakeClientset, "default")

	// 测试Close
	client.Close()

	// 验证context已被取消
	select {
	case <-client.GetContext().Done():
		// Context已取消，符合预期
	default:
		t.Error("Context should be cancelled after Close()")
	}
}

// TestLoadConfig 测试配置加载
func TestLoadConfig(t *testing.T) {
	// 保存原始配置
	originalConfig := config.GlobalConfig
	defer func() {
		config.GlobalConfig = originalConfig
	}()

	t.Run("nil global config", func(t *testing.T) {
		config.GlobalConfig = nil
		cfg, err := LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "global config not loaded")
	})

	t.Run("k8s not enabled", func(t *testing.T) {
		config.GlobalConfig = &config.Config{
			K8s: config.K8sConfig{
				Enabled: false,
			},
		}
		cfg, err := LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "kubernetes is not enabled")
	})

	t.Run("k8s enabled with in-cluster", func(t *testing.T) {
		config.GlobalConfig = &config.Config{
			K8s: config.K8sConfig{
				Enabled:   true,
				InCluster: true,
				Namespace: "test-ns",
			},
		}
		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "test-ns", cfg.Namespace)
		assert.True(t, cfg.InCluster)
	})

	t.Run("k8s enabled but validation fails", func(t *testing.T) {
		config.GlobalConfig = &config.Config{
			K8s: config.K8sConfig{
				Enabled:    true,
				InCluster:  false,
				KubeConfig: "/nonexistent/kubeconfig",
				Namespace:  "test-ns",
			},
		}
		cfg, err := LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "kubeconfig file not found")
	})
}

// TestGetClient 测试获取全局客户端
func TestGetClient(t *testing.T) {
	// 保存原始配置和客户端
	originalConfig := config.GlobalConfig
	originalClient := globalClient
	defer func() {
		config.GlobalConfig = originalConfig
		globalClient = originalClient
	}()

	t.Run("nil global config", func(t *testing.T) {
		config.GlobalConfig = nil
		globalClient = nil

		client, err := GetClient()
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("k8s not enabled", func(t *testing.T) {
		config.GlobalConfig = &config.Config{
			K8s: config.K8sConfig{
				Enabled: false,
			},
		}
		globalClient = nil

		client, err := GetClient()
		assert.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("return existing client", func(t *testing.T) {
		// 创建一个fake client
		fakeClientset := fake.NewSimpleClientset()
		existingClient := NewClientWithClientset(fakeClientset, "test")
		globalClient = existingClient

		client, err := GetClient()
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, existingClient, client)
	})
}

// TestInitClient 测试初始化全局客户端
func TestInitClient(t *testing.T) {
	// 保存原始客户端
	originalClient := globalClient
	defer func() {
		globalClient = originalClient
	}()

	t.Run("init with valid config", func(t *testing.T) {
		globalClient = nil

		cfg := &Config{
			InCluster: true,
			Namespace: "test-ns",
			Timeout:   30 * time.Second,
		}

		// 由于InitClient会调用NewClient，而NewClient需要真实的K8s环境
		// 这里只能测试错误路径
		err := InitClient(cfg)
		assert.Error(t, err) // 预期会失败，因为没有真实的K8s环境
	})

	t.Run("replace existing client", func(t *testing.T) {
		// 创建一个现有的fake client
		fakeClientset := fake.NewSimpleClientset()
		existingClient := NewClientWithClientset(fakeClientset, "old")
		globalClient = existingClient

		cfg := &Config{
			InCluster: true,
			Namespace: "new-ns",
			Timeout:   30 * time.Second,
		}

		// 尝试初始化新客户端（会失败，但会先关闭旧客户端）
		err := InitClient(cfg)
		assert.Error(t, err) // 预期会失败
	})
}
