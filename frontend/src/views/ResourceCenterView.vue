<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Monitor,
  Management,
  Platform,
  FolderOpened,
  Picture,
  Shop,
  DataBoard,
  TrendCharts
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'

const router = useRouter()
const loading = ref(false)

// 资源统计
const resourceStats = ref({
  totalServers: 0,
  activeServers: 0,
  totalGpus: 0,
  availableGpus: 0,
  totalStorage: 0,
  usedStorage: 0,
  totalDatasets: 0,
  totalImages: 0
})

// 快捷入口
const quickLinks = [
  {
    title: '资源平台',
    description: '查看资源使用情况和统计数据',
    icon: DataBoard,
    path: '/resource-platform',
    color: '#409EFF'
  },
  {
    title: '资源管理',
    description: '管理和配置资源信息',
    icon: Management,
    path: '/resource-list',
    color: '#67C23A'
  },
  {
    title: '算力市场',
    description: '浏览和租用算力资源',
    icon: Shop,
    path: '/computing-market',
    color: '#E6A23C'
  },
  {
    title: '云主机',
    description: '选择和管理云主机',
    icon: Platform,
    path: '/hosts',
    color: '#F56C6C'
  },
  {
    title: '数据集',
    description: '管理训练数据集',
    icon: FolderOpened,
    path: '/datasets',
    color: '#909399'
  },
  {
    title: '镜像',
    description: '管理容器镜像',
    icon: Picture,
    path: '/images',
    color: '#409EFF'
  }
]

// 加载资源统计
const loadResourceStats = async () => {
  loading.value = true
  try {
    // TODO: 调用API获取数据
    await new Promise(resolve => setTimeout(resolve, 500))
    resourceStats.value = {
      totalServers: 156,
      activeServers: 142,
      totalGpus: 1248,
      availableGpus: 856,
      totalStorage: 5120,
      usedStorage: 3240,
      totalDatasets: 89,
      totalImages: 45
    }
  } catch (error) {
    ElMessage.error('加载资源统计失败')
  } finally {
    loading.value = false
  }
}

// 导航到页面
const navigateTo = (path: string) => {
  router.push(path)
}

onMounted(() => {
  loadResourceStats()
})
</script>

<template>
  <div class="resource-center">
    <PageHeader title="资源中心" />

    <div v-loading="loading" class="center-content">
      <!-- 资源统计卡片 -->
      <div class="stats-section">
        <h3 class="section-title">资源概览</h3>
        <div class="stats-grid">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #409EFF">
                <el-icon :size="32"><Platform /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">服务器总数</div>
                <div class="stat-value">{{ resourceStats.totalServers }}</div>
                <div class="stat-sub">活跃: {{ resourceStats.activeServers }}</div>
              </div>
            </div>
          </el-card>

          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #67C23A">
                <el-icon :size="32"><Monitor /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">GPU总数</div>
                <div class="stat-value">{{ resourceStats.totalGpus }}</div>
                <div class="stat-sub">可用: {{ resourceStats.availableGpus }}</div>
              </div>
            </div>
          </el-card>

          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #E6A23C">
                <el-icon :size="32"><FolderOpened /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">存储空间 (TB)</div>
                <div class="stat-value">{{ resourceStats.totalStorage }}</div>
                <div class="stat-sub">已用: {{ resourceStats.usedStorage }} TB</div>
              </div>
            </div>
          </el-card>

          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #909399">
                <el-icon :size="32"><Picture /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-label">数据集/镜像</div>
                <div class="stat-value">{{ resourceStats.totalDatasets + resourceStats.totalImages }}</div>
                <div class="stat-sub">数据集: {{ resourceStats.totalDatasets }} | 镜像: {{ resourceStats.totalImages }}</div>
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 快捷入口 -->
      <div class="quick-links-section">
        <h3 class="section-title">快捷入口</h3>
        <div class="quick-links-grid">
          <el-card
            v-for="link in quickLinks"
            :key="link.path"
            class="quick-link-card"
            shadow="hover"
            @click="navigateTo(link.path)"
          >
            <div class="quick-link-content">
              <div class="quick-link-icon" :style="{ background: link.color }">
                <el-icon :size="40">
                  <component :is="link.icon" />
                </el-icon>
              </div>
              <div class="quick-link-info">
                <h4 class="quick-link-title">{{ link.title }}</h4>
                <p class="quick-link-description">{{ link.description }}</p>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.resource-center {
  padding: 24px;
}

.center-content {
  min-height: 600px;
}

.section-title {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

/* 资源统计样式 */
.stats-section {
  margin-bottom: 32px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.stat-sub {
  font-size: 12px;
  color: #909399;
}

/* 快捷入口样式 */
.quick-links-section {
  margin-bottom: 24px;
}

.quick-links-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.quick-link-card {
  cursor: pointer;
  transition: all 0.3s;
}

.quick-link-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.quick-link-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.quick-link-icon {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.quick-link-info {
  flex: 1;
  min-width: 0;
}

.quick-link-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.quick-link-description {
  margin: 0;
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .quick-links-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .quick-links-grid {
    grid-template-columns: 1fr;
  }
}
</style>
