<script setup lang="ts">
import { ref, onMounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'

// 注册 ECharts 组件
use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
])

const option = ref({
  title: {
    text: '资源使用趋势',
    left: 'left',
    textStyle: {
      fontSize: 16,
      fontWeight: 600,
    },
  },
  tooltip: {
    trigger: 'axis',
  },
  legend: {
    data: ['CPU', 'GPU', '内存'],
    top: 'bottom',
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '15%',
    containLabel: true,
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00'],
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      formatter: '{value}%',
    },
  },
  series: [
    {
      name: 'CPU',
      type: 'line',
      smooth: true,
      data: [30, 35, 45, 60, 55, 50, 40],
      itemStyle: {
        color: '#409EFF',
      },
    },
    {
      name: 'GPU',
      type: 'line',
      smooth: true,
      data: [20, 40, 60, 80, 75, 65, 50],
      itemStyle: {
        color: '#67C23A',
      },
    },
    {
      name: '内存',
      type: 'line',
      smooth: true,
      data: [40, 45, 50, 55, 60, 58, 52],
      itemStyle: {
        color: '#E6A23C',
      },
    },
  ],
})
</script>

<template>
  <div class="resource-chart">
    <VChart :option="option" :autoresize="true" class="chart" />
  </div>
</template>

<style scoped>
.resource-chart {
  background: white;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  height: 400px;
}

.chart {
  width: 100%;
  height: 100%;
}
</style>
