package image

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"github.com/YoungBoyGod/remotegpu/pkg/harbor"
	"gorm.io/gorm"
)

type ImageService struct {
	imageDao     *dao.ImageDao
	harborClient *harbor.Client
	harborConfig *config.HarborConfig
}

func NewImageService(db *gorm.DB) *ImageService {
	svc := &ImageService{
		imageDao: dao.NewImageDao(db),
	}
	// 如果 Harbor 配置启用，初始化客户端
	if config.GlobalConfig != nil && config.GlobalConfig.Harbor.Enabled {
		cfg := &config.GlobalConfig.Harbor
		svc.harborConfig = cfg
		svc.harborClient = harbor.NewClient(cfg.Endpoint, cfg.Username, cfg.Password)
	}
	return svc
}

func (s *ImageService) List(ctx context.Context, params dao.ImageListParams) ([]entity.Image, int64, error) {
	return s.imageDao.List(ctx, params)
}

func (s *ImageService) Create(ctx context.Context, img *entity.Image) error {
	return s.imageDao.Create(ctx, img)
}

func (s *ImageService) GetByID(ctx context.Context, id uint) (*entity.Image, error) {
	return s.imageDao.FindByID(ctx, id)
}

func (s *ImageService) Update(ctx context.Context, img *entity.Image) error {
	return s.imageDao.Update(ctx, img)
}

func (s *ImageService) Delete(ctx context.Context, id uint) error {
	return s.imageDao.Delete(ctx, id)
}

// Sync 从 Harbor 镜像仓库同步镜像列表
// 同步逻辑：遍历 Harbor 项目下所有仓库和标签，去重后写入数据库
func (s *ImageService) Sync(ctx context.Context) error {
	if s.harborClient == nil || s.harborConfig == nil {
		return fmt.Errorf("Harbor 未配置或未启用")
	}

	project := s.harborConfig.Project
	if project == "" {
		return fmt.Errorf("Harbor 项目名未配置")
	}

	// 1. 从 Harbor 获取仓库列表
	repos, err := s.harborClient.ListRepositories(ctx, project)
	if err != nil {
		return fmt.Errorf("获取 Harbor 仓库列表失败: %w", err)
	}

	// 2. 收集所有镜像名称，用于批量去重
	var allImages []imageCandidate
	for _, repo := range repos {
		artifacts, err := s.harborClient.ListArtifacts(ctx, project, extractRepoName(repo.Name))
		if err != nil {
			log.Printf("[Harbor Sync] 获取仓库 %s 的制品失败: %v", repo.Name, err)
			continue
		}
		for _, artifact := range artifacts {
			for _, tag := range artifact.Tags {
				fullName := fmt.Sprintf("%s:%s", repo.Name, tag.Name)
				allImages = append(allImages, imageCandidate{
					name:        fullName,
					repoName:    repo.Name,
					tag:         tag.Name,
					description: repo.Description,
				})
			}
		}
	}

	if len(allImages) == 0 {
		return nil
	}

	// 3. 批量查询已存在的镜像名称
	names := make([]string, len(allImages))
	for i, img := range allImages {
		names[i] = img.name
	}
	existingNames, err := s.imageDao.ListExistingNames(ctx, names)
	if err != nil {
		return fmt.Errorf("查询已有镜像失败: %w", err)
	}

	// 4. 插入新镜像（跳过已存在的）
	var created int
	for _, img := range allImages {
		if existingNames[img.name] {
			continue
		}
		registryURL := fmt.Sprintf("%s/%s", s.harborConfig.Endpoint, img.name)
		newImage := &entity.Image{
			Name:        img.name,
			DisplayName: img.name,
			Description: img.description,
			RegistryURL: registryURL,
			IsOfficial:  true,
			Status:      "active",
		}
		// 尝试从仓库名称推断框架信息
		newImage.Framework, newImage.CUDAVersion = parseImageMeta(img.repoName)

		if err := s.imageDao.Create(ctx, newImage); err != nil {
			log.Printf("[Harbor Sync] 创建镜像 %s 失败: %v", img.name, err)
			continue
		}
		created++
	}

	log.Printf("[Harbor Sync] 同步完成: 扫描 %d 个镜像, 新增 %d 个", len(allImages), created)
	return nil
}

// imageCandidate 同步过程中的镜像候选项
type imageCandidate struct {
	name        string
	repoName    string
	tag         string
	description string
}

// extractRepoName 从 "project/repo" 格式中提取仓库名
// Harbor API 返回的仓库名格式为 "project/repo"，但查询制品时只需要 "repo" 部分
func extractRepoName(fullName string) string {
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return fullName
}

// parseImageMeta 从仓库名称推断框架和 CUDA 版本
// 例如 "pytorch-cuda11.8" → framework="pytorch", cuda="11.8"
func parseImageMeta(repoName string) (framework, cudaVersion string) {
	lower := strings.ToLower(repoName)

	// 推断框架
	frameworks := []string{"pytorch", "tensorflow", "paddlepaddle", "jax", "mxnet"}
	for _, fw := range frameworks {
		if strings.Contains(lower, fw) {
			framework = fw
			break
		}
	}

	// 推断 CUDA 版本
	if idx := strings.Index(lower, "cuda"); idx >= 0 {
		rest := lower[idx+4:]
		// 跳过可能的分隔符
		rest = strings.TrimLeft(rest, "-_")
		// 提取版本号（如 "11.8"）
		var ver strings.Builder
		for _, ch := range rest {
			if (ch >= '0' && ch <= '9') || ch == '.' {
				ver.WriteRune(ch)
			} else {
				break
			}
		}
		if ver.Len() > 0 {
			cudaVersion = ver.String()
		}
	}

	return
}
