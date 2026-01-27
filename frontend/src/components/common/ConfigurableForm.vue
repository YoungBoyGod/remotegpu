<script setup lang="ts">
import { computed } from 'vue'
import type { FormFieldConfig } from '@/config/formFields'
import ResourceSlider from './ResourceSlider.vue'

interface Props {
  fields: FormFieldConfig[]
  modelValue: Record<string, any>
  labelWidth?: string
}

const props = withDefaults(defineProps<Props>(), {
  labelWidth: '120px',
})

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
}>()

// 过滤掉隐藏的字段
const visibleFields = computed(() => {
  return props.fields.filter(field => !field.hidden)
})

// 更新字段值
const updateField = (prop: string, value: any) => {
  emit('update:modelValue', {
    ...props.modelValue,
    [prop]: value,
  })
}
</script>

<template>
  <el-form :model="modelValue" :label-width="labelWidth">
    <el-row :gutter="20">
      <el-col
        v-for="field in visibleFields"
        :key="field.prop"
        :span="field.span || 24"
      >
        <el-form-item
          :label="field.label"
          :prop="field.prop"
          :required="field.required"
          :rules="field.rules"
        >
          <!-- Input -->
          <el-input
            v-if="field.type === 'input'"
            :model-value="modelValue[field.prop]"
            :placeholder="field.placeholder"
            :disabled="field.disabled"
            @update:model-value="updateField(field.prop, $event)"
          />

          <!-- Textarea -->
          <el-input
            v-else-if="field.type === 'textarea'"
            :model-value="modelValue[field.prop]"
            type="textarea"
            :placeholder="field.placeholder"
            :disabled="field.disabled"
            :rows="4"
            @update:model-value="updateField(field.prop, $event)"
          />

          <!-- Select -->
          <el-select
            v-else-if="field.type === 'select'"
            :model-value="modelValue[field.prop]"
            :placeholder="field.placeholder"
            :disabled="field.disabled"
            style="width: 100%"
            @update:model-value="updateField(field.prop, $event)"
          >
            <el-option
              v-for="option in field.options"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>

          <!-- Number -->
          <el-input-number
            v-else-if="field.type === 'number'"
            :model-value="modelValue[field.prop]"
            :min="field.min"
            :max="field.max"
            :step="field.step"
            :disabled="field.disabled"
            @update:model-value="updateField(field.prop, $event)"
          />

          <!-- Slider -->
          <ResourceSlider
            v-else-if="field.type === 'slider'"
            :model-value="modelValue[field.prop]"
            :label="field.label"
            :min="field.min"
            :max="field.max"
            :step="field.step"
            :unit="field.unit"
            @update:model-value="updateField(field.prop, $event)"
          />

          <!-- Switch -->
          <el-switch
            v-else-if="field.type === 'switch'"
            :model-value="modelValue[field.prop]"
            :disabled="field.disabled"
            @update:model-value="updateField(field.prop, $event)"
          />

          <!-- Radio -->
          <el-radio-group
            v-else-if="field.type === 'radio'"
            :model-value="modelValue[field.prop]"
            :disabled="field.disabled"
            @update:model-value="updateField(field.prop, $event)"
          >
            <el-radio
              v-for="option in field.options"
              :key="option.value"
              :label="option.value"
            >
              {{ option.label }}
            </el-radio>
          </el-radio-group>

          <!-- Checkbox -->
          <el-checkbox-group
            v-else-if="field.type === 'checkbox'"
            :model-value="modelValue[field.prop]"
            :disabled="field.disabled"
            @update:model-value="updateField(field.prop, $event)"
          >
            <el-checkbox
              v-for="option in field.options"
              :key="option.value"
              :label="option.value"
            >
              {{ option.label }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>
