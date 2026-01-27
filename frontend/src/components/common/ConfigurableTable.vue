<script setup lang="ts">
import { computed } from 'vue'
import type { TableColumnConfig } from '@/config/tableColumns'

interface Props {
  columns: TableColumnConfig[]
  data: any[]
  loading?: boolean
  stripe?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  stripe: true,
})

// 过滤掉隐藏的列
const visibleColumns = computed(() => {
  return props.columns.filter(col => !col.hidden)
})
</script>

<template>
  <el-table :data="data" :loading="loading" :stripe="stripe" style="width: 100%">
    <el-table-column
      v-for="column in visibleColumns"
      :key="column.prop"
      :prop="column.slot ? undefined : column.prop"
      :label="column.label"
      :width="column.width"
      :min-width="column.minWidth"
      :sortable="column.sortable"
      :fixed="column.fixed"
    >
      <template v-if="column.slot" #default="scope">
        <slot :name="column.slot" :row="scope.row" :column="column" />
      </template>
    </el-table-column>
  </el-table>
</template>
