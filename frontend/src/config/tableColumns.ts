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

/**
 * 管理员任务列表表格列配置
 */
export const adminTaskColumns: TableColumnConfig[] = [
  {
    prop: 'id',
    label: '任务ID',
    width: 140,
    sortable: true,
  },
  {
    prop: 'name',
    label: '任务名称',
    minWidth: 180,
    sortable: true,
    slot: 'name',
  },
  {
    prop: 'status',
    label: '状态',
    width: 110,
    sortable: true,
    slot: 'status',
  },
  {
    prop: 'customer_id',
    label: '客户',
    minWidth: 140,
    slot: 'customer',
  },
  {
    prop: 'host_id',
    label: '机器',
    minWidth: 160,
    slot: 'host',
  },
  {
    prop: 'image',
    label: '镜像',
    minWidth: 160,
    slot: 'image',
    hidden: true,
  },
  {
    prop: 'command',
    label: '命令',
    minWidth: 220,
    slot: 'command',
    hidden: true,
  },
  {
    prop: 'created_at',
    label: '创建时间',
    width: 170,
    sortable: true,
    slot: 'created_at',
  },
  {
    prop: 'started_at',
    label: '开始时间',
    width: 170,
    sortable: true,
    slot: 'started_at',
    hidden: true,
  },
  {
    prop: 'finished_at',
    label: '结束时间',
    width: 170,
    sortable: true,
    slot: 'finished_at',
    hidden: true,
  },
  {
    prop: 'exit_code',
    label: '退出码',
    width: 90,
    hidden: true,
  },
  {
    prop: 'error_msg',
    label: '错误信息',
    minWidth: 200,
    hidden: true,
  },
  {
    prop: 'actions',
    label: '操作',
    width: 260,
    fixed: 'right',
    slot: 'actions',
  },
]

/**
 * 客户任务列表表格列配置
 */
export const customerTaskColumns: TableColumnConfig[] = [
  {
    prop: 'id',
    label: '任务ID',
    width: 140,
    sortable: true,
  },
  {
    prop: 'name',
    label: '任务名称',
    minWidth: 180,
    sortable: true,
    slot: 'name',
  },
  {
    prop: 'status',
    label: '状态',
    width: 110,
    sortable: true,
    slot: 'status',
  },
  {
    prop: 'host_id',
    label: '机器',
    minWidth: 160,
    slot: 'host',
  },
  {
    prop: 'image',
    label: '镜像',
    minWidth: 160,
    slot: 'image',
    hidden: true,
  },
  {
    prop: 'command',
    label: '命令',
    minWidth: 220,
    slot: 'command',
    hidden: true,
  },
  {
    prop: 'created_at',
    label: '创建时间',
    width: 170,
    sortable: true,
    slot: 'created_at',
  },
  {
    prop: 'started_at',
    label: '开始时间',
    width: 170,
    sortable: true,
    slot: 'started_at',
    hidden: true,
  },
  {
    prop: 'finished_at',
    label: '结束时间',
    width: 170,
    sortable: true,
    slot: 'finished_at',
    hidden: true,
  },
  {
    prop: 'exit_code',
    label: '退出码',
    width: 90,
    hidden: true,
  },
  {
    prop: 'error_msg',
    label: '错误信息',
    minWidth: 200,
    hidden: true,
  },
  {
    prop: 'actions',
    label: '操作',
    width: 240,
    fixed: 'right',
    slot: 'actions',
  },
]
