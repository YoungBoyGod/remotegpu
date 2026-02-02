package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/database"
	"gorm.io/gorm"
)

// PortPoolService 端口池管理服务
type PortPoolService struct {
	portMappingDao *dao.PortMappingDao
	mu             sync.Mutex // 保护端口分配的并发安全
}

// NewPortPoolService 创建端口池服务
func NewPortPoolService() *PortPoolService {
	return &PortPoolService{
		portMappingDao: dao.NewPortMappingDao(),
	}
}

// PortRange 端口范围配置
type PortRange struct {
	Start int // 起始端口
	End   int // 结束端口
}

// 默认端口范围配置
var (
	// SSHPortRange SSH 端口范围 (22000-22999)
	SSHPortRange = PortRange{Start: 22000, End: 22999}
	// RDPPortRange RDP 端口范围 (33890-34889)
	RDPPortRange = PortRange{Start: 33890, End: 34889}
	// JupyterPortRange Jupyter 端口范围 (8888-9887)
	JupyterPortRange = PortRange{Start: 8888, End: 9887}
	// VNCPortRange VNC 端口范围 (5900-6899)
	VNCPortRange = PortRange{Start: 5900, End: 6899}
	// NoVNCPortRange noVNC 端口范围 (6080-7079)
	NoVNCPortRange = PortRange{Start: 6080, End: 7079}
	// CustomPortRange 自定义服务端口范围 (10000-19999)
	CustomPortRange = PortRange{Start: 10000, End: 19999}
)

// GetPortRangeByServiceType 根据服务类型获取端口范围
func GetPortRangeByServiceType(serviceType string) PortRange {
	switch serviceType {
	case "ssh":
		return SSHPortRange
	case "rdp":
		return RDPPortRange
	case "jupyter":
		return JupyterPortRange
	case "vnc":
		return VNCPortRange
	case "novnc":
		return NoVNCPortRange
	case "tensorboard", "vscode", "custom":
		return CustomPortRange
	default:
		return CustomPortRange
	}
}

// AllocatePort 分配端口
func (s *PortPoolService) AllocatePort(envID string, serviceType string, internalPort int, protocol string, description string) (*entity.PortMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取端口范围
	portRange := GetPortRangeByServiceType(serviceType)

	// 查询已分配的端口
	db := database.GetDB()
	var usedPorts []int
	err := db.Model(&entity.PortMapping{}).
		Where("status = ?", "active").
		Pluck("external_port", &usedPorts).Error
	if err != nil {
		return nil, fmt.Errorf("查询已分配端口失败: %w", err)
	}

	// 创建已使用端口的 map
	usedPortMap := make(map[int]bool)
	for _, port := range usedPorts {
		usedPortMap[port] = true
	}

	// 查找可用端口
	var availablePort int
	for port := portRange.Start; port <= portRange.End; port++ {
		if !usedPortMap[port] {
			availablePort = port
			break
		}
	}

	if availablePort == 0 {
		return nil, fmt.Errorf("端口范围 %d-%d 已用尽", portRange.Start, portRange.End)
	}

	// 创建端口映射记录
	portMapping := &entity.PortMapping{
		EnvID:        envID,
		ServiceType:  serviceType,
		InternalPort: internalPort,
		ExternalPort: availablePort,
		Protocol:     protocol,
		Description:  description,
		Status:       "active",
	}

	if err := db.Create(portMapping).Error; err != nil {
		return nil, fmt.Errorf("创建端口映射失败: %w", err)
	}

	return portMapping, nil
}

// ReleasePort 释放端口
func (s *PortPoolService) ReleasePort(portMappingID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	db := database.GetDB()
	now := time.Now()

	result := db.Model(&entity.PortMapping{}).
		Where("id = ? AND status = ?", portMappingID, "active").
		Updates(map[string]interface{}{
			"status":      "released",
			"released_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("释放端口失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("端口映射不存在或已释放")
	}

	return nil
}

// ReleasePortsByEnvID 释放环境的所有端口
func (s *PortPoolService) ReleasePortsByEnvID(envID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	db := database.GetDB()
	now := time.Now()

	result := db.Model(&entity.PortMapping{}).
		Where("env_id = ? AND status = ?", envID, "active").
		Updates(map[string]interface{}{
			"status":      "released",
			"released_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("释放环境端口失败: %w", result.Error)
	}

	return nil
}

// AllocatePorts 批量分配端口
func (s *PortPoolService) AllocatePorts(envID string, portRequests []PortRequest) ([]*entity.PortMapping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var portMappings []*entity.PortMapping

	// 查询已分配的端口
	db := database.GetDB()
	var usedPorts []int
	err := db.Model(&entity.PortMapping{}).
		Where("status = ?", "active").
		Pluck("external_port", &usedPorts).Error
	if err != nil {
		return nil, fmt.Errorf("查询已分配端口失败: %w", err)
	}

	// 创建已使用端口的 map
	usedPortMap := make(map[int]bool)
	for _, port := range usedPorts {
		usedPortMap[port] = true
	}

	// 为每个请求分配端口
	for _, req := range portRequests {
		portRange := GetPortRangeByServiceType(req.ServiceType)

		// 查找可用端口
		var availablePort int
		for port := portRange.Start; port <= portRange.End; port++ {
			if !usedPortMap[port] {
				availablePort = port
				usedPortMap[port] = true // 标记为已使用,避免重复分配
				break
			}
		}

		if availablePort == 0 {
			return nil, fmt.Errorf("端口范围 %d-%d 已用尽", portRange.Start, portRange.End)
		}

		// 创建端口映射记录
		portMapping := &entity.PortMapping{
			EnvID:        envID,
			ServiceType:  req.ServiceType,
			InternalPort: req.InternalPort,
			ExternalPort: availablePort,
			Protocol:     req.Protocol,
			Description:  req.Description,
			Status:       "active",
		}

		portMappings = append(portMappings, portMapping)
	}

	// 批量插入数据库
	if err := db.Create(&portMappings).Error; err != nil {
		return nil, fmt.Errorf("批量创建端口映射失败: %w", err)
	}

	return portMappings, nil
}

// PortRequest 端口分配请求
type PortRequest struct {
	ServiceType  string // 服务类型
	InternalPort int    // 内部端口
	Protocol     string // 协议
	Description  string // 描述
}

// GetAllocatedPorts 获取环境的所有已分配端口
func (s *PortPoolService) GetAllocatedPorts(envID string) ([]*entity.PortMapping, error) {
	db := database.GetDB()
	var portMappings []*entity.PortMapping

	err := db.Where("env_id = ? AND status = ?", envID, "active").
		Order("service_type, internal_port").
		Find(&portMappings).Error

	if err != nil {
		return nil, fmt.Errorf("查询端口映射失败: %w", err)
	}

	return portMappings, nil
}

// GetPortByServiceType 获取环境指定服务类型的端口映射
func (s *PortPoolService) GetPortByServiceType(envID string, serviceType string) (*entity.PortMapping, error) {
	db := database.GetDB()
	var portMapping entity.PortMapping

	err := db.Where("env_id = ? AND service_type = ? AND status = ?", envID, serviceType, "active").
		First(&portMapping).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询端口映射失败: %w", err)
	}

	return &portMapping, nil
}
