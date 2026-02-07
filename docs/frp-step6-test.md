# frp方案 - 第六步:测试验证

## 目标

验证整个frp方案是否正常工作。

---

## SSH访问测试

### 从本地电脑测试

```bash
# 测试GPU1 SSH
ssh -p 10001 user@云服务器IP

# 测试GPU2 SSH
ssh -p 10002 user@云服务器IP
```

**预期结果**: 能够成功登录到对应的GPU机器。

---

## Web服务测试

### 测试Jupyter

```bash
# 测试GPU1 Jupyter
curl -I https://gpu1-jupyter.gpu.domain.com

# 测试GPU2 Jupyter
curl -I https://gpu2-jupyter.gpu.domain.com
```

**预期结果**: 返回HTTP 200或302状态码。

### 浏览器访问

直接在浏览器访问:
- `https://gpu1-jupyter.gpu.domain.com`
- `https://gpu1-tensorboard.gpu.domain.com`

**预期结果**: 能够看到Jupyter或TensorBoard界面。

---

## 检查SSL证书

```bash
# 检查证书有效性
openssl s_client -connect gpu1-jupyter.gpu.domain.com:443 -servername gpu1-jupyter.gpu.domain.com
```

**预期结果**: 显示证书信息,无错误。

---

## 性能测试

### 测试SSH传输速度

```bash
# 在GPU机器创建测试文件
dd if=/dev/zero of=/tmp/test.dat bs=1M count=100

# 从本地下载测试
scp -P 10001 user@云服务器IP:/tmp/test.dat /tmp/
```

### 测试Web服务响应时间

```bash
curl -o /dev/null -s -w "Time: %{time_total}s\n" https://gpu1-jupyter.gpu.domain.com
```

---

## 故障排查

### 问题1: SSH无法连接

**排查步骤**:
1. 检查frpc是否运行: `sudo systemctl status frpc`
2. 检查frps是否运行: `sudo systemctl status frps`
3. 查看frpc日志: `sudo journalctl -u frpc -n 50`
4. 测试端口: `telnet 云服务器IP 10001`

### 问题2: Web服务无法访问

**排查步骤**:
1. 检查nginx状态: `sudo systemctl status nginx`
2. 检查nginx配置: `sudo nginx -t`
3. 查看nginx日志: `sudo tail -f /var/log/nginx/error.log`
4. 测试本地端口: `curl http://127.0.0.1:11001`

### 问题3: SSL证书错误

**排查步骤**:
1. 检查证书文件: `sudo ls -l /etc/letsencrypt/live/gpu.domain.com/`
2. 检查nginx SSL配置
3. 重新申请证书: `sudo certbot renew --force-renewal`

---

## 完成!

如果所有测试通过,说明frp方案配置成功!

**批量配置**: 对于200台GPU机器,使用批量脚本简化配置,详见 `frp-batch-scripts.md`
