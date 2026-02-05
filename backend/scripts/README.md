# Scripts 目录说明

测试和运维脚本集合。

## 目录结构

```
scripts/
├── shell/     # Shell 脚本
├── python/    # Python 脚本
├── http/      # HTTP 请求文件 (VS Code REST Client)
└── go/        # Go 脚本
```

## 使用方法

### Shell 脚本

```bash
# 测试设备 API
./shell/test-machine-api.sh

# 添加测试设备
export API_TOKEN=your_token
./shell/add-test-machine.sh 192.168.1.100 test-gpu-01
```

### HTTP 请求

使用 VS Code REST Client 插件打开 `.http` 文件执行请求。
