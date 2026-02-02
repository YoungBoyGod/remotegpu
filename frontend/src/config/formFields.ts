/**
 * 表单字段配置
 * 统一管理所有表单的字段定义
 */

export interface FormFieldConfig {
  prop: string
  label: string
  type: 'input' | 'textarea' | 'select' | 'number' | 'slider' | 'switch' | 'radio' | 'checkbox' | 'date' | 'upload'
  required?: boolean
  placeholder?: string
  options?: Array<{ label: string; value: any }>
  min?: number
  max?: number
  step?: number
  unit?: string
  rules?: any[]
  disabled?: boolean
  hidden?: boolean
  span?: number // 栅格占位
}

/**
 * 环境创建表单配置
 */
export const environmentFormFields: FormFieldConfig[] = [
  {
    prop: 'name',
    label: '环境名称',
    type: 'input',
    required: true,
    placeholder: '请输入环境名称',
    span: 24,
  },
  {
    prop: 'image',
    label: '镜像',
    type: 'select',
    required: true,
    placeholder: '请选择镜像',
    options: [
      { label: 'pytorch/pytorch:2.0-cuda11.8', value: 'pytorch/pytorch:2.0-cuda11.8' },
      { label: 'tensorflow/tensorflow:2.13-gpu', value: 'tensorflow/tensorflow:2.13-gpu' },
    ],
    span: 24,
  },
  {
    prop: 'gpu',
    label: 'GPU 数量',
    type: 'slider',
    required: true,
    min: 1,
    max: 8,
    step: 1,
    unit: '个',
    span: 24,
  },
  {
    prop: 'cpu',
    label: 'CPU 核心',
    type: 'slider',
    required: true,
    min: 1,
    max: 32,
    step: 1,
    unit: '核',
    span: 24,
  },
  {
    prop: 'memory',
    label: '内存',
    type: 'slider',
    required: true,
    min: 4,
    max: 128,
    step: 4,
    unit: 'GB',
    span: 24,
  },
  {
    prop: 'description',
    label: '描述',
    type: 'textarea',
    placeholder: '请输入环境描述',
    span: 24,
  },
]

/**
 * 数据集上传表单配置
 */
export const datasetFormFields: FormFieldConfig[] = [
  {
    prop: 'name',
    label: '数据集名称',
    type: 'input',
    required: true,
    placeholder: '请输入数据集名称',
    span: 24,
  },
  {
    prop: 'version',
    label: '版本',
    type: 'input',
    required: true,
    placeholder: '例如: v1.0',
    span: 12,
  },
  {
    prop: 'visibility',
    label: '可见性',
    type: 'radio',
    required: true,
    options: [
      { label: '公开', value: 'public' },
      { label: '私有', value: 'private' },
    ],
    span: 12,
  },
  {
    prop: 'description',
    label: '描述',
    type: 'textarea',
    placeholder: '请输入数据集描述',
    span: 24,
  },
]

/**
 * 用户设置表单配置
 */
export const userSettingsFormFields: FormFieldConfig[] = [
  {
    prop: 'username',
    label: '用户名',
    type: 'input',
    required: true,
    disabled: true,
    span: 12,
  },
  {
    prop: 'email',
    label: '邮箱',
    type: 'input',
    required: true,
    placeholder: '请输入邮箱',
    span: 12,
  },
  {
    prop: 'fullName',
    label: '姓名',
    type: 'input',
    placeholder: '请输入姓名',
    span: 12,
  },
  {
    prop: 'phone',
    label: '手机号',
    type: 'input',
    placeholder: '请输入手机号',
    span: 12,
  },
]
