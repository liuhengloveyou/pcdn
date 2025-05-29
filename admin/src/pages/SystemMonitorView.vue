<script lang="ts" setup>
import Page from '@/components/basic-page.vue'
import { ref, onMounted, onUnmounted } from 'vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, GaugeChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  ToolboxComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import { computed } from 'vue'

// 注册必要的组件
use([
  CanvasRenderer,
  LineChart,
  GaugeChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  ToolboxComponent
])

// 模拟数据
const cpuUsage = ref(0)
const memoryUsage = ref(0)
const networkSpeed = ref([])
const diskUsage = ref([])

// 添加CPU和内存的历史数据数组
const cpuHistory = ref([])
const memoryHistory = ref([])

// 时间轴数据
const timeData = ref([])

// 更新间隔（毫秒）
const updateInterval = 2000
let timer: number | null | undefined = null

// CPU使用率图表选项 - 改为线图
const cpuOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    formatter: '{b}<br />CPU: {c0}%'
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: timeData.value
  },
  yAxis: {
    type: 'value',
    min: 0,
    max: 100,
    name: '%'
  },
  series: [
    {
      name: 'CPU',
      type: 'line',
      smooth: true,
      data: cpuHistory.value,
      itemStyle: {
        color: '#3b82f6'
      },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [{
            offset: 0, color: 'rgba(29, 78, 216, 0.3)'
          }, {
            offset: 1, color: 'rgba(59, 130, 246, 0.1)'
          }]
        }
      }
    }
  ]
}))

// 内存使用率图表选项 - 改为线图
const memoryOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    formatter: '{b}<br />内存: {c0}%'
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: timeData.value
  },
  yAxis: {
    type: 'value',
    min: 0,
    max: 100,
    name: '%'
  },
  series: [
    {
      name: '内存',
      type: 'line',
      smooth: true,
      data: memoryHistory.value,
      itemStyle: {
        color: '#10b981'
      },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [{
            offset: 0, color: 'rgba(5, 150, 105, 0.3)'
          }, {
            offset: 1, color: 'rgba(16, 185, 129, 0.1)'
          }]
        }
      }
    }
  ]
}))

// 网络速度图表选项
const networkOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    formatter: '{b}<br />上传: {c0} KB/s<br />下载: {c1} KB/s'
  },
  legend: {
    data: ['上传', '下载'],
    top: 10
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: timeData.value
  },
  yAxis: {
    type: 'value',
    name: 'KB/s'
  },
  series: [
    {
      name: '上传',
      type: 'line',
      smooth: true,
      data: networkSpeed.value.map(item => item.upload),
      itemStyle: {
        color: '#3b82f6'
      }
    },
    {
      name: '下载',
      type: 'line',
      smooth: true,
      data: networkSpeed.value.map(item => item.download),
      itemStyle: {
        color: '#10b981'
      }
    }
  ]
}))

// 硬盘使用率图表选项
const diskOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c}GB ({d}%)'
  },
  legend: {
    orient: 'vertical',
    left: 'left',
    top: 'center'
  },
  series: [
    {
      type: 'pie',
      radius: ['50%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 16,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: diskUsage.value
    }
  ]
}))

// 模拟数据更新
const updateData = () => {
  // 更新CPU使用率（0-100之间的随机值）
  cpuUsage.value = Math.floor(Math.random() * 100)
  
  // 更新内存使用率（20-90之间的随机值）
  memoryUsage.value = Math.floor(Math.random() * 70) + 20
  
  // 更新时间轴
  const now = new Date()
  const timeStr = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`
  
  // 保持最多10个数据点
  if (timeData.value.length >= 10) {
    timeData.value.shift()
    networkSpeed.value.shift()
    cpuHistory.value.shift()
    memoryHistory.value.shift()
  }
  
  timeData.value.push(timeStr)
  
  // 更新CPU和内存历史数据
  cpuHistory.value.push(cpuUsage.value)
  memoryHistory.value.push(memoryUsage.value)
  
  // 更新网络速度
  networkSpeed.value.push({
    upload: Math.floor(Math.random() * 100),
    download: Math.floor(Math.random() * 500)
  })
  
  // 初始化硬盘使用情况（只在第一次更新时设置）
  if (diskUsage.value.length === 0) {
    diskUsage.value = [
      { value: 120, name: '系统 (C:)' },
      { value: 234, name: '数据 (D:)' },
      { value: 335, name: '备份 (E:)' },
      { value: 148, name: '媒体 (F:)' },
      { value: 548, name: '可用空间' }
    ]
  }
}

onMounted(() => {
  // 初始更新一次数据
  updateData()
  
  // 设置定时更新
  timer = setInterval(updateData, updateInterval)
})

onUnmounted(() => {
  // 清除定时器
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<template>
  <Page
    title="系统监控"
    description="实时监控系统资源使用情况"
    sticky
  >
    <!-- CPU和内存使用率 -->
    <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle class="text-base font-semibold">CPU 使用率</CardTitle>
          <!-- 使用更现代的CPU图标 -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-blue-500">
            <rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect><rect x="9" y="9" width="6" height="6"></rect><line x1="9" y1="1" x2="9" y2="4"></line><line x1="15" y1="1" x2="15" y2="4"></line><line x1="9" y1="20" x2="9" y2="23"></line><line x1="15" y1="20" x2="15" y2="23"></line><line x1="1" y1="9" x2="4" y2="9"></line><line x1="1" y1="14" x2="4" y2="14"></line><line x1="20" y1="9" x2="23" y2="9"></line><line x1="20" y1="14" x2="23" y2="14"></line>
          </svg>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-blue-600 mb-2">{{ cpuUsage.toFixed(1) }}%</div>
          <v-chart class="chart" :option="cpuOption" autoresize />
        </CardContent>
      </Card>
      
      <Card>
        <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle class="text-base font-semibold">内存使用率</CardTitle>
          <!-- 使用更现代的内存图标 -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-emerald-500">
            <path d="M2 13h6l2-6 2 12 2-8 2 8h6"/>
          </svg>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-emerald-600 mb-2">{{ memoryUsage.toFixed(1) }}%</div>
          <v-chart class="chart" :option="memoryOption" autoresize />
        </CardContent>
      </Card>
      
      <Card class="sm:col-span-2">
        <CardHeader class="flex flex-row items-center justify-between pb-2 space-y-0">
          <CardTitle class="text-base font-semibold">系统信息</CardTitle>
          <!-- 使用信息图标 -->
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-gray-500">
            <circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line>
          </svg>
        </CardHeader>
        <CardContent>
          <div class="space-y-3 text-sm">
            <div class="flex justify-between items-center">
              <span class="text-muted-foreground">操作系统:</span>
              <span class="font-medium">Windows 10 Pro</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-muted-foreground">处理器:</span>
              <span class="font-medium truncate max-w-[180px]" title="Intel Core i7-10700K @ 3.80GHz">Intel Core i7-10700K @ 3.80GHz</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-muted-foreground">内存总量:</span>
              <span class="font-medium">32 GB DDR4</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-muted-foreground">系统运行时间:</span>
              <span class="font-medium">10天 14小时 23分钟</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-muted-foreground">IP地址:</span>
              <span class="font-medium">192.168.1.100</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
    
    <!-- 网络和硬盘使用情况 -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-7">
      <Card class="col-span-1 lg:col-span-4">
        <CardHeader>
          <CardTitle class="text-base font-semibold">网络流量</CardTitle>
          <CardDescription class="text-sm">实时网络上传和下载速度</CardDescription>
        </CardHeader>
        <CardContent>
          <v-chart class="chart-large" :option="networkOption" autoresize />
        </CardContent>
      </Card>
      
      <Card class="col-span-1 lg:col-span-3">
        <CardHeader>
          <CardTitle class="text-base font-semibold">硬盘使用情况</CardTitle>
          <CardDescription class="text-sm">各分区存储空间使用情况</CardDescription>
        </CardHeader>
        <CardContent>
          <v-chart class="chart-large" :option="diskOption" autoresize />
        </CardContent>
      </Card>
    </div>
  </Page>
</template>

<style scoped>
.chart {
  height: 180px; /* 稍微减小一点高度以适应新增的数值显示 */
}

.chart-large {
  height: 300px;
}

/* 可以添加一些全局的卡片悬浮效果 */
.card:hover {
  /* transform: translateY(-2px); */ /* 轻微上浮效果 */
  /* box-shadow: 0 4px 12px rgba(0,0,0,0.1); */ /* 增加阴影 */
  /* transition: all 0.2s ease-in-out; */
}
</style>