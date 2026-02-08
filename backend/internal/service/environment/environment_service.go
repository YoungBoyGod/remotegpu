package environment

import (
	"context"
	"errors"
	"strconv"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

// EnvironmentService 环境管理业务逻辑层
type EnvironmentService struct {
	db     *gorm.DB
	envDao *dao.EnvironmentDao
}

func NewEnvironmentService(db *gorm.DB) *EnvironmentService {
	return &EnvironmentService{
		db:     db,
		envDao: dao.NewEnvironmentDao(db),
	}
}

// Create 创建环境
func (s *EnvironmentService) Create(ctx context.Context, env *entity.Environment) error {
	return s.envDao.Create(ctx, env)
}

// GetByID 获取环境详情
func (s *EnvironmentService) GetByID(ctx context.Context, id string) (*entity.Environment, error) {
	return s.envDao.FindByID(ctx, id)
}

// List 获取用户的环境列表
func (s *EnvironmentService) List(ctx context.Context, customerID uint, workspaceID *uint) ([]entity.Environment, error) {
	db := s.db.WithContext(ctx).
		Preload("Host").
		Where("user_id = ?", customerID)

	if workspaceID != nil {
		db = db.Where("workspace_id = ?", *workspaceID)
	}

	var envs []entity.Environment
	err := db.Order("created_at DESC").Find(&envs).Error
	return envs, err
}

// Start 启动环境
func (s *EnvironmentService) Start(ctx context.Context, id string, customerID uint) error {
	env, err := s.envDao.FindByID(ctx, id)
	if err != nil {
		return errors.New("环境不存在")
	}
	if env.UserID != customerID {
		return errors.New("无权操作该环境")
	}
	if env.Status != "stopped" {
		return errors.New("仅已停止的环境可启动")
	}
	return s.db.WithContext(ctx).Model(&entity.Environment{}).
		Where("id = ?", id).Update("status", "running").Error
}

// Stop 停止环境
func (s *EnvironmentService) Stop(ctx context.Context, id string, customerID uint) error {
	env, err := s.envDao.FindByID(ctx, id)
	if err != nil {
		return errors.New("环境不存在")
	}
	if env.UserID != customerID {
		return errors.New("无权操作该环境")
	}
	if env.Status != "running" {
		return errors.New("仅运行中的环境可停止")
	}
	return s.db.WithContext(ctx).Model(&entity.Environment{}).
		Where("id = ?", id).Update("status", "stopped").Error
}

// Delete 删除环境
func (s *EnvironmentService) Delete(ctx context.Context, id string, customerID uint) error {
	env, err := s.envDao.FindByID(ctx, id)
	if err != nil {
		return errors.New("环境不存在")
	}
	if env.UserID != customerID {
		return errors.New("无权操作该环境")
	}
	if env.Status != "stopped" && env.Status != "error" {
		return errors.New("仅已停止或异常的环境可删除")
	}
	return s.db.WithContext(ctx).Delete(&entity.Environment{}, "id = ?", id).Error
}

// AccessInfo 访问信息结构
type AccessInfo struct {
	SSH     *SSHAccess     `json:"ssh,omitempty"`
	Jupyter *JupyterAccess `json:"jupyter,omitempty"`
	VNC     *VNCAccess     `json:"vnc,omitempty"`
}

type SSHAccess struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}

type JupyterAccess struct {
	URL   string `json:"url"`
	Token string `json:"token,omitempty"`
}

type VNCAccess struct {
	URL string `json:"url"`
}

// GetAccessInfo 获取环境访问信息
func (s *EnvironmentService) GetAccessInfo(ctx context.Context, id string, customerID uint) (*AccessInfo, error) {
	env, err := s.envDao.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("环境不存在")
	}
	if env.UserID != customerID {
		return nil, errors.New("无权访问该环境")
	}
	if env.Status != "running" {
		return nil, errors.New("仅运行中的环境可获取访问信息")
	}

	info := &AccessInfo{}

	// 从 Host 获取外部地址（优先使用 ExternalIP，其次 PublicIP，最后 IPAddress）
	hostAddr := ""
	if env.Host != nil {
		if env.Host.ExternalIP != "" {
			hostAddr = env.Host.ExternalIP
		} else if env.Host.PublicIP != "" {
			hostAddr = env.Host.PublicIP
		} else {
			hostAddr = env.Host.IPAddress
		}
	}

	// 根据端口映射组装访问信息
	for _, pm := range env.PortMappings {
		if pm.Status != "active" {
			continue
		}
		switch pm.ServiceType {
		case "ssh":
			info.SSH = &SSHAccess{
				Host:     hostAddr,
				Port:     pm.ExternalPort,
				Username: "root",
			}
		case "jupyter":
			info.Jupyter = &JupyterAccess{
				URL: "http://" + hostAddr + ":" + intToStr(pm.ExternalPort),
			}
		case "vnc", "rdp":
			info.VNC = &VNCAccess{
				URL: "http://" + hostAddr + ":" + intToStr(pm.ExternalPort),
			}
		}
	}

	return info, nil
}

func intToStr(n int) string {
	return strconv.Itoa(n)
}
