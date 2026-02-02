/**
 * 表格列配置
 * 统一管理所有表格的列定义
 */

export interface TableColumnConfig {
  prop: string
  label: string
  width?: string | number
  minWidth?: string | number
  sortable?: boolean
  fixed?: boolean | 'left' | 'right'
  slot?: string // 自定义插槽名称
  hidden?: boolean // 是否隐藏
}

/**
 * 环境列表表格列配置
 */
export const environmentColumns: TableColumnConfig[] = [
  {
    prop: 'name',
    label: '环境名称',
    minWidth: 180,
    sortable: true,
    slot: 'name',
  },
  {
    prop: 'status',
    label: '状态',
    width: 100,
    sortable: true,
    slot: 'status',
  },
  {
    prop: 'image',
    label: '镜像',
    minWidth: 200,
    sortable: true,
  },
  {
    prop: 'gpu',
    label: 'GPU',
    width: 150,
    sortable: true,
  },
  {
    prop: 'cpu',
    label: 'CPU/内存',
    width: 120,
    sortable: true,
    slot: 'cpu-memory',
  },
  {
    prop: 'runningTime',
    label: '运行时长',
    width: 120,
  },
  {
    prop: 'createdAt',
    label: '创建时间',
    width: 160,
    sortable: true,
  },
  {
    prop: 'actions',
    label: '操作',
    width: 200,
    fixed: 'right',
    slot: 'actions',
  },
]

/**
 * 数据集列表表格列配置
 */
export const datasetColumns: TableColumnConfig[] = [
  {
    prop: 'name',
    label: '数据集名称',
    minWidth: 200,
    sortable: true,
    slot: 'name',
  },
  {
    prop: 'version',
    label: '版本',
    width: 100,
    sortable: true,
  },
  {
    prop: 'size',
    label: '大小',
    width: 120,
    sortable: true,
  },
  {
    prop: 'fileCount',
    label: '文件数量',
    width: 120,
    sortable: true,
  },
  {
    prop: 'visibility',
    label: '可见性',
    width: 100,
    sortable: true,
    slot: 'visibility',
  },
  {
    prop: 'createdAt',
    label: '创建时间',
    width: 150,
    sortable: true,
  },
  {
    prop: 'actions',
    label: '操作',
    width: 150,
    fixed: 'right',
    slot: 'actions',
  },
]

