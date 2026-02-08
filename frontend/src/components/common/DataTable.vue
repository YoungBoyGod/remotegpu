<script setup lang="ts">
interface Props {
  data: any[]
  loading?: boolean
  emptyText?: string
  showPagination?: boolean
  total?: number
  currentPage?: number
  pageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: '暂无数据',
  showPagination: false,
  total: 0,
  currentPage: 1,
  pageSize: 10,
})

const emit = defineEmits<{
  'update:currentPage': [value: number]
  'update:pageSize': [value: number]
  'page-change': [page: number]
  'size-change': [size: number]
  'selection-change': [selection: any[]]
}>()

const handleCurrentChange = (page: number) => {
  emit('update:currentPage', page)
  emit('page-change', page)
}

const handleSizeChange = (size: number) => {
  emit('update:pageSize', size)
  emit('size-change', size)
}
</script>

<template>
  <div class="data-table">
    <el-table
      :data="data"
      v-loading="loading"
      stripe
      style="width: 100%"
      @selection-change="(val: any[]) => emit('selection-change', val)"
    >
      <slot />
      <template #empty>
        <el-empty :description="emptyText" />
      </template>
    </el-table>

    <el-pagination
      v-if="showPagination && total > 0"
      class="pagination"
      :current-page="currentPage"
      :page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<style scoped>
.data-table {
  background: white;
  border-radius: 4px;
  padding: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
