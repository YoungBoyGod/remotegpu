# Windows 机器远程访问与管理方案

> Windows GPU 服务器远程访问技术方案
>
> 创建日期：2026-01-26

---

## 目录

1. [整体架构](#1-整体架构)
2. [远程访问方案对比](#2-远程访问方案对比)
3. [RDP 远程桌面方案](#3-rdp-远程桌面方案)
4. [SSH 访问方案](#4-ssh-访问方案)
5. [容器化方案](#5-容器化方案)
6. [混合架构设计](#6-混合架构设计)
7. [实施方案](#7-实施方案)

---

## 1. 整体架构

### 1.1 Windows vs Linux 对比

| 维度 | Linux 机器 | Windows 机器 |
|------|-----------|-------------|
| **主要访问方式** | SSH (Port 22) | RDP (Port 3389) |
| **容器技术** | Docker (原生支持) | Docker Desktop / Windows Container |
| **GPU 支持** | NVIDIA Docker Runtime | NVIDIA Container Toolkit for Windows |
| **开发环境** | JupyterLab + SSH | Jupyter + RDP / VSCode Remote |
| **用户隔离** | Linux 用户 | Windows 用户 |
| **资源管理** | cgroups | Job Objects |

### 1.2 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      用户访问层                               │
├─────────────────────────────────────────────────────────────┤
│  Linux 用户:  ssh developer@host:30001                      │
│  Windows 用户: rdp://host:33001  或  ssh developer@host:30002│
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    负载均衡 / 网关                            │
│                  gateway.example.com                         │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Linux 主机1   │  │ Linux 主机2   │  │ Windows 主机1 │
│ GPU x4       │  │ GPU x4       │  │ GPU x4       │
│ Docker       │  │ Docker       │  │ Hyper-V      │
└──────────────┘  └──────────────┘  └──────────────┘
```

---

## 2. 远程访问方案对比

### 2.1 方案总览

| 方案 | 协议 | 端口 | 优点 | 缺点 | 适用场景 |
|------|------|------|------|------|---------|
| **RDP** | Remote Desktop | 3389 | 图形界面、原生支持 | 带宽占用大 | Windows 开发、GUI 应用 |
| **SSH** | OpenSSH | 22 | 轻量、跨平台 | 需要安装 | 命令行操作、脚本 |
| **VNC** | VNC Protocol | 5900 | 跨平台 | 性能差、不安全 | 不推荐 |
| **Web RDP** | HTML5 | 443 | 浏览器访问 | 性能一般 | 临时访问 |
| **容器化** | Docker | - | 隔离性好 | Windows 容器限制多 | 特定场景 |

### 2.2 推荐方案

**主推：RDP（远程桌面）+ SSH（可选）**

- ✅ RDP 作为主要访问方式（图形界面）
- ✅ SSH 作为辅助方式（命令行、脚本）
- ✅ Web RDP 作为补充（浏览器访问）

---

## 3. RDP 远程桌面方案

### 3.1 架构设计

```
用户 -> RDP 客户端 -> 网关 (Port 33001) -> Windows 主机 (Port 3389)
```

### 3.2 Windows 配置

#### A. 启用远程桌面

```powershell
# PowerShell 脚本：启用 RDP
Set-ItemProperty -Path 'HKLM:\System\CurrentControlSet\Control\Terminal Server' `
    -Name "fDenyTSConnections" -Value 0

# 启用防火墙规则
Enable-NetFirewallRule -DisplayGroup "Remote Desktop"

# 允许特定用户远程访问
Add-LocalGroupMember -Group "Remote Desktop Users" -Member "developer"
```

#### B. 创建开发用户

```powershell
# 创建用户脚本
$Username = "developer"
$Password = ConvertTo-SecureString "GeneratedPassword123!" -AsPlainText -Force

# 创建用户
New-LocalUser -Name $Username -Password $Password -FullName "Developer User" `
    -Description "Development Environment User"

# 添加到用户组
Add-LocalGroupMember -Group "Users" -Member $Username
Add-LocalGroupMember -Group "Remote Desktop Users" -Member $Username

# 设置用户权限（可选：管理员权限）
# Add-LocalGroupMember -Group "Administrators" -Member $Username
```

#### C. 配置 RDP 安全策略

```powershell
# 设置网络级别身份验证（NLA）
Set-ItemProperty -Path 'HKLM:\System\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp' `
    -Name "UserAuthentication" -Value 1

# 设置加密级别（高）
Set-ItemProperty -Path 'HKLM:\System\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp' `
    -Name "MinEncryptionLevel" -Value 3

# 设置最大连接数
Set-ItemProperty -Path 'HKLM:\System\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp' `
    -Name "MaxInstanceCount" -Value 2
```

### 3.3 端口映射

#### 方案 A：动态端口映射（推荐）

```powershell
# 使用 netsh 配置端口转发
netsh interface portproxy add v4tov4 `
    listenport=33001 `
    listenaddress=0.0.0.0 `
    connectport=3389 `
    connectaddress=192.168.1.10

# 查看端口映射
netsh interface portproxy show all

# 删除端口映射
netsh interface portproxy delete v4tov4 listenport=33001 listenaddress=0.0.0.0
```

#### 方案 B：修改 RDP 默认端口

```powershell
# 修改 RDP 监听端口
Set-ItemProperty -Path 'HKLM:\System\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp' `
    -Name "PortNumber" -Value 33001

# 重启 RDP 服务
Restart-Service TermService -Force

# 更新防火墙规则
New-NetFirewallRule -DisplayName "RDP Custom Port" `
    -Direction Inbound -Protocol TCP -LocalPort 33001 -Action Allow
```

### 3.4 Go 管理代码

```go
// windows/rdp_manager.go
package windows

import (
    "fmt"
    "os/exec"
)

type RDPManager struct {
    host string
}

// 创建 Windows 用户
func (m *RDPManager) CreateUser(username, password string) error {
    // 通过 WinRM 远程执行 PowerShell
    script := fmt.Sprintf(`
        $Password = ConvertTo-SecureString "%s" -AsPlainText -Force
        New-LocalUser -Name "%s" -Password $Password -FullName "Dev User"
        Add-LocalGroupMember -Group "Remote Desktop Users" -Member "%s"
    `, password, username, username)

    return m.executeRemotePowerShell(script)
}

// 删除用户
func (m *RDPManager) DeleteUser(username string) error {
    script := fmt.Sprintf(`Remove-LocalUser -Name "%s"`, username)
    return m.executeRemotePowerShell(script)
}

// 配置端口映射
func (m *RDPManager) ConfigurePortMapping(externalPort, internalPort int) error {
    script := fmt.Sprintf(`
        netsh interface portproxy add v4tov4 `
            listenport=%d listenaddress=0.0.0.0 `
            connectport=%d connectaddress=127.0.0.1
    `, externalPort, internalPort)

    return m.executeRemotePowerShell(script)
}

// 执行远程 PowerShell（通过 WinRM）
func (m *RDPManager) executeRemotePowerShell(script string) error {
    cmd := exec.Command(
        "winrm",
        "invoke",
        "Create",
        "wmicimv2/Win32_Process",
        fmt.Sprintf("-r:http://%s:5985", m.host),
        fmt.Sprintf(`-CommandLine:"powershell.exe -Command %s"`, script),
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("执行失败: %s, %v", string(output), err)
    }

    return nil
}

// 获取 RDP 连接信息
func (m *RDPManager) GetConnectionInfo(envID string) (*RDPConnectionInfo, error) {
    var env Environment
    if err := db.Where("id = ?", envID).First(&env).Error; err != nil {
        return nil, err
    }

    return &RDPConnectionInfo{
        Host:     m.host,
        Port:     env.RDPPort,
        Username: env.Username,
        Password: env.Password,
        RDPFile:  m.generateRDPFile(env),
    }, nil
}

// 生成 RDP 连接文件
func (m *RDPManager) generateRDPFile(env Environment) string {
    return fmt.Sprintf(`
screen mode id:i:2
use multimon:i:0
desktopwidth:i:1920
desktopheight:i:1080
session bpp:i:32
winposstr:s:0,3,0,0,800,600
compression:i:1
keyboardhook:i:2
audiocapturemode:i:0
videoplaybackmode:i:1
connection type:i:7
networkautodetect:i:1
bandwidthautodetect:i:1
displayconnectionbar:i:1
enableworkspacereconnect:i:0
disable wallpaper:i:0
allow font smoothing:i:0
allow desktop composition:i:0
disable full window drag:i:1
disable menu anims:i:1
disable themes:i:0
disable cursor setting:i:0
bitmapcachepersistenable:i:1
full address:s:%s:%d
audiomode:i:0
redirectprinters:i:0
redirectcomports:i:0
redirectsmartcards:i:0
redirectclipboard:i:1
redirectposdevices:i:0
autoreconnection enabled:i:1
authentication level:i:2
prompt for credentials:i:0
negotiate security layer:i:1
remoteapplicationmode:i:0
alternate shell:s:
shell working directory:s:
gatewayhostname:s:
gatewayusagemethod:i:4
gatewaycredentialssource:i:4
gatewayprofileusagemethod:i:0
promptcredentialonce:i:0
gatewaybrokeringtype:i:0
use redirection server name:i:0
rdgiskdcproxy:i:0
kdcproxyname:s:
username:s:%s
`, m.host, env.RDPPort, env.Username)
}
```

### 3.5 API 实现

```go
// API: 创建 Windows 开发环境
// POST /api/environments/windows
func CreateWindowsEnvironment(c *gin.Context) {
    var req struct {
        Name      string `json:"name"`
        Resources struct {
            CPU    int `json:"cpu"`
            Memory int `json:"memory"`
            GPU    int `json:"gpu"`
        } `json:"resources"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    customer := GetCurrentCustomer(c)

    // 1. 选择 Windows 主机
    host := scheduler.SelectWindowsHost(req.Resources)
    if host == nil {
        c.JSON(503, gin.H{"error": "暂无可用的 Windows 主机"})
        return
    }

    // 2. 分配端口
    rdpPort, _ := portPool.Allocate()

    // 3. 生成用户名和密码
    username := fmt.Sprintf("dev%d", customer.ID)
    password := GeneratePassword()

    // 4. 创建 Windows 用户
    rdpManager := NewRDPManager(host.IPAddress)
    if err := rdpManager.CreateUser(username, password); err != nil {
        c.JSON(500, gin.H{"error": "创建用户失败"})
        return
    }

    // 5. 配置端口映射
    if err := rdpManager.ConfigurePortMapping(rdpPort, 3389); err != nil {
        c.JSON(500, gin.H{"error": "配置端口失败"})
        return
    }

    // 6. 创建数据库记录
    env := &Environment{
        ID:         GenerateEnvID(),
        CustomerID: customer.ID,
        HostID:     host.ID,
        Type:       "windows",
        RDPPort:    rdpPort,
        Username:   username,
        Password:   password,
        Status:     "running",
    }
    db.Create(env)

    // 7. 返回连接信息
    c.JSON(200, gin.H{
        "env_id":   env.ID,
        "rdp_host": host.PublicIP,
        "rdp_port": rdpPort,
        "username": username,
        "password": password,
        "rdp_file": rdpManager.generateRDPFile(*env),
    })
}
```

### 3.6 前端展示

```typescript
// WindowsEnvironment.tsx
import React from 'react';

interface RDPConnectionProps {
  host: string;
  port: number;
  username: string;
  password: string;
  rdpFile: string;
}

export const RDPConnection: React.FC<RDPConnectionProps> = ({
  host, port, username, password, rdpFile
}) => {
  const downloadRDPFile = () => {
    const blob = new Blob([rdpFile], { type: 'application/rdp' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'connection.rdp';
    a.click();
  };

  return (
    <div className="rdp-connection">
      <h3>Windows 远程桌面连接</h3>

      <div className="connection-info">
        <p><strong>主机:</strong> {host}</p>
        <p><strong>端口:</strong> {port}</p>
        <p><strong>用户名:</strong> {username}</p>
        <p><strong>密码:</strong> {password} <button>复制</button></p>
      </div>

      <div className="connection-methods">
        <h4>连接方式：</h4>

        <div className="method">
          <h5>方式 1: 下载 RDP 文件（推荐）</h5>
          <button onClick={downloadRDPFile}>下载 RDP 文件</button>
          <p>双击下载的文件即可连接</p>
        </div>

        <div className="method">
          <h5>方式 2: 手动连接</h5>
          <ol>
            <li>打开"远程桌面连接"（mstsc.exe）</li>
            <li>输入地址: {host}:{port}</li>
            <li>输入用户名和密码</li>
          </ol>
        </div>

        <div className="method">
          <h5>方式 3: Web RDP（浏览器）</h5>
          <a href={`/web-rdp/${env.id}`} target="_blank">
            在浏览器中打开
          </a>
        </div>
      </div>
    </div>
  );
};
```

---

## 4. SSH 访问方案

### 4.1 安装 OpenSSH Server

```powershell
# 安装 OpenSSH Server
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0

# 启动 SSH 服务
Start-Service sshd

# 设置自动启动
Set-Service -Name sshd -StartupType 'Automatic'

# 配置防火墙
New-NetFirewallRule -Name sshd -DisplayName 'OpenSSH Server (sshd)' `
    -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22
```

### 4.2 SSH 配置

```powershell
# 编辑 SSH 配置文件
# C:\ProgramData\ssh\sshd_config

# 允许密码认证
PasswordAuthentication yes

# 允许公钥认证
PubkeyAuthentication yes

# 设置默认 Shell 为 PowerShell
New-ItemProperty -Path "HKLM:\SOFTWARE\OpenSSH" -Name DefaultShell `
    -Value "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" -PropertyType String -Force

# 重启 SSH 服务
Restart-Service sshd
```

### 4.3 SSH 密钥管理

```powershell
# 配置用户的 SSH 公钥
$Username = "developer"
$PublicKey = "ssh-rsa AAAAB3NzaC1yc2E..."

# 创建 .ssh 目录
$SSHDir = "C:\Users\$Username\.ssh"
New-Item -ItemType Directory -Path $SSHDir -Force

# 写入公钥
$PublicKey | Out-File -FilePath "$SSHDir\authorized_keys" -Encoding ascii

# 设置权限（只有用户可以访问）
icacls "$SSHDir\authorized_keys" /inheritance:r
icacls "$SSHDir\authorized_keys" /grant "${Username}:F"
```

### 4.4 SSH 使用示例

```bash
# 用户连接
ssh developer@windows-host:30022

# 执行命令
ssh developer@windows-host:30022 "Get-Process | Select-Object -First 10"

# 文件传输
scp local-file.txt developer@windows-host:30022:C:/Users/developer/
```

---

## 5. 容器化方案

### 5.1 Windows Container 限制

**Windows 容器的限制：**
- ❌ 不支持 GUI 应用（无法运行 RDP）
- ❌ GPU 支持有限
- ❌ 镜像体积大（几 GB）
- ❌ 生态不如 Linux 容器

**适用场景：**
- ✅ 命令行应用
- ✅ Web 服务
- ✅ 批处理任务

### 5.2 Hyper-V 虚拟机方案（推荐）

```powershell
# 创建虚拟机
New-VM -Name "dev-env-123" `
    -MemoryStartupBytes 16GB `
    -Generation 2 `
    -NewVHDPath "C:\VMs\dev-env-123.vhdx" `
    -NewVHDSizeBytes 500GB

# 配置 CPU
Set-VMProcessor -VMName "dev-env-123" -Count 8

# 配置 GPU（DDA - Discrete Device Assignment）
Dismount-VMHostAssignableDevice -LocationPath "PCIROOT(0)#PCI(0300)#PCI(0000)"
Add-VMAssignableDevice -VMName "dev-env-123" -LocationPath "PCIROOT(0)#PCI(0300)#PCI(0000)"

# 启动虚拟机
Start-VM -Name "dev-env-123"
```

### 5.3 GPU 直通配置

```powershell
# 1. 启用 SR-IOV（如果 GPU 支持）
Set-VM -VMName "dev-env-123" -GuestControlledCacheTypes $true
Set-VM -VMName "dev-env-123" -LowMemoryMappedIoSpace 3GB
Set-VM -VMName "dev-env-123" -HighMemoryMappedIoSpace 33GB

# 2. 配置 GPU 直通
$GPU = Get-PnpDevice | Where-Object {$_.FriendlyName -like "*NVIDIA*"}
Disable-PnpDevice -InstanceId $GPU.InstanceId -Confirm:$false
Dismount-VMHostAssignableDevice -LocationPath $GPU.LocationPath
Add-VMAssignableDevice -VMName "dev-env-123" -LocationPath $GPU.LocationPath
```

---

## 6. 混合架构设计（Linux + Windows）

### 6.1 统一管理架构

```
┌─────────────────────────────────────────────────────────────┐
│                      管理平台 API                             │
│                  (Go/Python 后端)                            │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Linux Agent  │  │ Linux Agent  │  │Windows Agent │
│ (Go)         │  │ (Go)         │  │ (Go/PS)      │
├──────────────┤  ├──────────────┤  ├──────────────┤
│ Docker       │  │ Docker       │  │ WinRM        │
│ SSH:22       │  │ SSH:22       │  │ RDP:3389     │
│              │  │              │  │ SSH:22       │
└──────────────┘  └──────────────┘  └──────────────┘
```

### 6.2 数据库设计

```sql
-- 主机表（支持 Linux 和 Windows）
CREATE TABLE hosts (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128),
    ip_address VARCHAR(64) NOT NULL,
    os_type VARCHAR(20) NOT NULL,        -- linux, windows
    status VARCHAR(20) DEFAULT 'online',

    -- 资源信息
    total_cpu INT,
    total_memory BIGINT,
    total_gpu INT,
    used_cpu INT DEFAULT 0,
    used_memory BIGINT DEFAULT 0,
    used_gpu INT DEFAULT 0,

    -- Windows 特有
    winrm_port INT,                      -- WinRM 端口（5985）
    hyperv_enabled BOOLEAN DEFAULT false,

    last_heartbeat TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),

    INDEX idx_os_type (os_type),
    INDEX idx_status (status)
);

-- 环境表（支持多种类型）
CREATE TABLE environments (
    id VARCHAR(64) PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    host_id VARCHAR(64) NOT NULL,

    type VARCHAR(20) NOT NULL,           -- linux-docker, windows-native, windows-vm
    os_type VARCHAR(20),                 -- linux, windows

    -- SSH 信息（Linux 和 Windows 都支持）
    ssh_port INT,
    ssh_username VARCHAR(128),
    ssh_password VARCHAR(128),

    -- RDP 信息（仅 Windows）
    rdp_port INT,
    rdp_username VARCHAR(128),
    rdp_password VARCHAR(128),

    -- 容器/虚拟机 ID
    container_id VARCHAR(128),           -- Docker 容器 ID
    vm_name VARCHAR(128),                -- Hyper-V 虚拟机名称

    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),

    FOREIGN KEY (host_id) REFERENCES hosts(id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_type (type)
);
```

### 6.3 统一调度器

```go
// scheduler/unified_scheduler.go
package scheduler

type UnifiedScheduler struct {
    linuxHosts   []*Host
    windowsHosts []*Host
}

// 选择主机（根据操作系统类型）
func (s *UnifiedScheduler) SelectHost(req CreateEnvRequest) (*Host, error) {
    var hosts []*Host

    // 根据用户需求选择主机池
    if req.OSType == "windows" {
        hosts = s.windowsHosts
    } else {
        hosts = s.linuxHosts
    }

    // 筛选可用主机
    var candidates []*Host
    for _, host := range hosts {
        if host.Status != "online" {
            continue
        }

        if host.AvailableCPU >= req.CPU &&
           host.AvailableMemory >= req.Memory &&
           host.AvailableGPU >= req.GPU {
            candidates = append(candidates, host)
        }
    }

    if len(candidates) == 0 {
        return nil, errors.New("no available host")
    }

    // 负载均衡策略
    sort.Slice(candidates, func(i, j int) bool {
        return candidates[i].UsedCPU < candidates[j].UsedCPU
    })

    return candidates[0], nil
}

// 创建环境（统一接口）
func (s *UnifiedScheduler) CreateEnvironment(req CreateEnvRequest) (*Environment, error) {
    // 1. 选择主机
    host, err := s.SelectHost(req)
    if err != nil {
        return nil, err
    }

    // 2. 根据主机类型调用不同的创建方法
    if host.OSType == "windows" {
        return s.createWindowsEnvironment(host, req)
    } else {
        return s.createLinuxEnvironment(host, req)
    }
}

// 创建 Linux 环境
func (s *UnifiedScheduler) createLinuxEnvironment(host *Host, req CreateEnvRequest) (*Environment, error) {
    // 使用 Docker 创建容器
    agent := NewLinuxAgent(host.IPAddress)

    sshPort, _ := portPool.Allocate()
    password := GeneratePassword()

    containerID, err := agent.CreateContainer(CreateContainerRequest{
        Image:    req.Image,
        SSHPort:  sshPort,
        Password: password,
        Resources: req.Resources,
    })

    if err != nil {
        return nil, err
    }

    return &Environment{
        ID:          GenerateEnvID(),
        Type:        "linux-docker",
        OSType:      "linux",
        HostID:      host.ID,
        ContainerID: containerID,
        SSHPort:     sshPort,
        SSHUsername: "developer",
        SSHPassword: password,
        Status:      "running",
    }, nil
}

// 创建 Windows 环境
func (s *UnifiedScheduler) createWindowsEnvironment(host *Host, req CreateEnvRequest) (*Environment, error) {
    agent := NewWindowsAgent(host.IPAddress)

    rdpPort, _ := portPool.Allocate()
    sshPort, _ := portPool.Allocate()
    username := fmt.Sprintf("dev%d", req.CustomerID)
    password := GeneratePassword()

    // 创建 Windows 用户
    if err := agent.CreateUser(username, password); err != nil {
        return nil, err
    }

    // 配置端口映射
    agent.ConfigurePortMapping(rdpPort, 3389)
    agent.ConfigurePortMapping(sshPort, 22)

    return &Environment{
        ID:          GenerateEnvID(),
        Type:        "windows-native",
        OSType:      "windows",
        HostID:      host.ID,
        RDPPort:     rdpPort,
        RDPUsername: username,
        RDPPassword: password,
        SSHPort:     sshPort,
        SSHUsername: username,
        SSHPassword: password,
        Status:      "running",
    }, nil
}
```

### 6.4 前端统一界面

```typescript
// EnvironmentSelector.tsx
import React, { useState } from 'react';

export const EnvironmentSelector: React.FC = () => {
  const [osType, setOSType] = useState<'linux' | 'windows'>('linux');

  return (
    <div className="environment-selector">
      <h3>创建开发环境</h3>

      <div className="os-selector">
        <label>操作系统：</label>
        <select value={osType} onChange={e => setOSType(e.target.value as any)}>
          <option value="linux">Linux (Ubuntu 20.04)</option>
          <option value="windows">Windows Server 2022</option>
        </select>
      </div>

      {osType === 'linux' && (
        <div className="linux-options">
          <h4>访问方式：</h4>
          <ul>
            <li>SSH 命令行</li>
            <li>JupyterLab（浏览器）</li>
            <li>VSCode Remote-SSH</li>
          </ul>
        </div>
      )}

      {osType === 'windows' && (
        <div className="windows-options">
          <h4>访问方式：</h4>
          <ul>
            <li>远程桌面（RDP）</li>
            <li>SSH 命令行（可选）</li>
            <li>Web RDP（浏览器）</li>
          </ul>
        </div>
      )}

      <button onClick={() => createEnvironment(osType)}>
        创建环境
      </button>
    </div>
  );
};
```

---

## 7. 实施方案

### 7.1 Windows 主机准备清单

```powershell
# Windows 主机初始化脚本
# setup-windows-host.ps1

# 1. 启用 Hyper-V（如果使用虚拟机方案）
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All

# 2. 安装 OpenSSH Server
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0
Start-Service sshd
Set-Service -Name sshd -StartupType 'Automatic'

# 3. 启用 WinRM（远程管理）
Enable-PSRemoting -Force
Set-Item WSMan:\localhost\Client\TrustedHosts -Value "*" -Force

# 4. 配置防火墙
New-NetFirewallRule -DisplayName "RDP" -Direction Inbound -Protocol TCP -LocalPort 3389 -Action Allow
New-NetFirewallRule -DisplayName "SSH" -Direction Inbound -Protocol TCP -LocalPort 22 -Action Allow
New-NetFirewallRule -DisplayName "WinRM" -Direction Inbound -Protocol TCP -LocalPort 5985 -Action Allow

# 5. 安装 NVIDIA 驱动和 CUDA
# 手动下载并安装：https://www.nvidia.com/Download/index.aspx

# 6. 安装 Python 和常用工具
choco install python git vscode -y

# 7. 创建工作目录
New-Item -ItemType Directory -Path "C:\Environments" -Force
```

### 7.2 部署步骤

#### 阶段 1：基础设施准备

**Linux 主机：**
```bash
# 1. 安装 Docker + NVIDIA Runtime
curl -fsSL https://get.docker.com | sh
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
    tee /etc/apt/sources.list.d/nvidia-docker.list
apt-get update && apt-get install -y nvidia-docker2
systemctl restart docker

# 2. 部署 Worker Agent
wget https://releases.example.com/linux-agent
chmod +x linux-agent
./linux-agent --master=http://master:8000 --host-id=linux-1
```

**Windows 主机：**
```powershell
# 1. 运行初始化脚本
.\setup-windows-host.ps1

# 2. 部署 Worker Agent
Invoke-WebRequest -Uri "https://releases.example.com/windows-agent.exe" -OutFile "C:\agent\windows-agent.exe"
Start-Process -FilePath "C:\agent\windows-agent.exe" -ArgumentList "--master=http://master:8000 --host-id=windows-1"
```

#### 阶段 2：测试验证

```bash
# 测试 Linux 环境创建
curl -X POST http://api.example.com/api/environments \
  -H "Content-Type: application/json" \
  -d '{
    "os_type": "linux",
    "resources": {"cpu": 4, "memory": 16, "gpu": 1}
  }'

# 测试 Windows 环境创建
curl -X POST http://api.example.com/api/environments \
  -H "Content-Type: application/json" \
  -d '{
    "os_type": "windows",
    "resources": {"cpu": 8, "memory": 32, "gpu": 1}
  }'
```

### 7.3 监控和维护

```go
// 监控 Windows 主机状态
func MonitorWindowsHost(hostID string) {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        // 1. 检查 WinRM 连接
        if !checkWinRMConnection(hostID) {
            alertManager.Send("Windows 主机 WinRM 连接失败", hostID)
        }

        // 2. 检查 RDP 服务
        if !checkRDPService(hostID) {
            alertManager.Send("Windows 主机 RDP 服务异常", hostID)
        }

        // 3. 检查资源使用
        resources := getWindowsHostResources(hostID)
        if resources.CPUUsage > 90 {
            alertManager.Send("Windows 主机 CPU 使用率过高", hostID)
        }
    }
}
```

---

## 8. 方案对比与选型

### 8.1 Windows 访问方案对比

| 方案 | 实现复杂度 | 用户体验 | 性能 | 成本 | 推荐度 |
|------|-----------|---------|------|------|--------|
| **RDP（原生）** | ⭐⭐ 简单 | ⭐⭐⭐⭐⭐ 优秀 | ⭐⭐⭐⭐⭐ 优秀 | 低 | ⭐⭐⭐⭐⭐ |
| **SSH** | ⭐⭐⭐ 中等 | ⭐⭐⭐ 良好 | ⭐⭐⭐⭐⭐ 优秀 | 低 | ⭐⭐⭐⭐ |
| **Web RDP** | ⭐⭐⭐⭐ 复杂 | ⭐⭐⭐ 良好 | ⭐⭐⭐ 一般 | 中 | ⭐⭐⭐ |
| **Hyper-V VM** | ⭐⭐⭐⭐⭐ 复杂 | ⭐⭐⭐⭐ 很好 | ⭐⭐⭐⭐ 很好 | 高 | ⭐⭐⭐ |

### 8.2 推荐方案

**MVP 阶段（快速启动）：**
```
✅ RDP 原生访问（主要）
✅ SSH 访问（辅助）
✅ 直接在 Windows 主机上创建用户
```

**生产阶段（规模化）：**
```
✅ RDP 原生访问
✅ SSH 访问
✅ Hyper-V 虚拟机隔离（可选）
✅ Web RDP（浏览器访问）
```

---

## 9. 常见问题

### Q1: Windows 和 Linux 能否共用端口池？
**A:** 可以，但建议分开管理：
- Linux SSH: 30000-31000
- Windows RDP: 33000-34000
- Windows SSH: 31000-32000

### Q2: Windows 用户隔离如何实现？
**A:** 三种方案：
1. **Windows 用户隔离**（推荐 MVP）：每个客户一个 Windows 用户
2. **Hyper-V 虚拟机**（推荐生产）：每个客户一个虚拟机
3. **Windows 容器**（不推荐）：限制太多

### Q3: RDP 带宽占用大怎么办？
**A:** 优化策略：
- 降低色彩深度（16 位）
- 禁用桌面背景和动画
- 使用 RemoteFX 压缩
- 限制分辨率（1920x1080）

### Q4: GPU 如何分配给 Windows 环境？
**A:** 两种方案：
1. **直接使用**：用户直接使用主机 GPU（多用户共享）
2. **GPU 直通**：Hyper-V DDA 技术（独占 GPU）

### Q5: Windows 环境如何挂载数据集？
**A:**
```powershell
# 方案 1：网络共享
New-SmbShare -Name "dataset-123" -Path "\\storage\datasets\customer-123\dataset-001"
net use Z: \\storage\dataset-123

# 方案 2：本地复制
robocopy \\storage\datasets\customer-123\dataset-001 C:\Data\dataset-001 /E
```

---

## 10. 总结

### 核心要点

1. **Windows 主要使用 RDP 访问**
   - 图形界面友好
   - 性能好
   - 原生支持

2. **SSH 作为辅助方式**
   - 命令行操作
   - 脚本执行
   - 文件传输

3. **统一管理 Linux 和 Windows**
   - 统一的 API 接口
   - 统一的调度器
   - 统一的前端界面

4. **渐进式实施**
   - MVP：直接创建 Windows 用户
   - 生产：Hyper-V 虚拟机隔离

### 技术栈总结

```
Windows 主机管理：
├── 远程管理：WinRM (PowerShell Remoting)
├── 用户管理：Windows 本地用户
├── 端口映射：netsh portproxy
├── 访问方式：RDP + SSH
└── 虚拟化：Hyper-V（可选）

统一管理平台：
├── 后端：Go (统一 API)
├── 调度器：支持 Linux + Windows
├── 数据库：PostgreSQL（统一元数据）
└── 前端：React（统一界面）
```

---

**文档结束**

这套方案支持 Linux 和 Windows 混合部署，提供统一的管理体验。
