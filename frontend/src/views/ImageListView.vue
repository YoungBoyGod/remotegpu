<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useRoleNavigation } from '@/composables/useRoleNavigation'
import PageHeader from '@/components/common/PageHeader.vue'
import FilterBar from '@/components/common/FilterBar.vue'

const router = useRouter()
const { navigateTo } = useRoleNavigation()

interface Image {
  id: string
  name: string
  tag: string
  type: string
  size: string
  createdAt: string
}

const images = ref<Image[]>([])
const loading = ref(false)
const searchText = ref('')

// 过滤后的镜像列表
const filteredImages = computed(() => {
  let result = images.value

  // 搜索过滤
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(img =>
      img.name.toLowerCase().includes(search) ||
      img.tag.toLowerCase().includes(search)
    )
  }

  return result
})

const loadImages = async () => {
  loading.value = true
  try {
    images.value = [
      {
        id: 'img-001',
        name: 'pytorch/pytorch',
        tag: '2.0-cuda11.8',
        type: 'official',
        size: '5.2 GB',
        createdAt: '2026-01-15',
      },
      {
        id: 'img-002',
        name: 'tensorflow/tensorflow',
        tag: '2.13-gpu',
        type: 'official',
        size: '4.8 GB',
        createdAt: '2026-01-16',
      },
    ]
  } catch (error) {
    ElMessage.error('加载镜像列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadImages()
})
</script>

<template>
  <div class="image-list">
    <PageHeader title="镜像管理">
      <template #actions>
        <el-button type="primary" :icon="Plus" @click="navigateTo('/images/build')">
          构建镜像
        </el-button>
      </template>
    </PageHeader>

    <FilterBar
      v-model:search-value="searchText"
      search-placeholder="搜索镜像名称"
    />

    <div class="image-grid">
      <el-card v-for="img in filteredImages" :key="img.id" class="image-card">
        <div class="image-info">
          <h3>{{ img.name }}</h3>
          <el-tag size="small">{{ img.tag }}</el-tag>
          <p class="size">{{ img.size }}</p>
          <p class="date">{{ img.createdAt }}</p>
        </div>
        <div class="image-actions">
          <el-button type="primary" size="small">使用</el-button>
          <el-button size="small">详情</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
.image-list {
  padding: 24px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.image-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.image-card:hover {
  transform: translateY(-4px);
}

.image-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
}

.image-info .size {
  color: #909399;
  margin: 8px 0;
}

.image-info .date {
  color: #c0c4cc;
  font-size: 12px;
  margin: 0;
}

.image-actions {
  margin-top: 16px;
  display: flex;
  gap: 8px;
}
</style>
