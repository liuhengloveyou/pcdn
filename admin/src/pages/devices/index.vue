<script setup lang="ts">
import Page from '@/components/basic-page.vue'
import { columns } from './components/columns'
import DataTable from './components/data-table.vue'
import DeviceCreate from './components/device-create.vue'
import { onMounted, onUnmounted, ref } from 'vue'
import { listDevice, type DeviceModel } from '@/services/DeviceService'
import { toast } from 'vue-sonner'

// 表格数据
const devices = ref<DeviceModel[]>([]);

// 加载状态
const loading = ref(false);
// 分页控制
const currentPage = ref(1);
const pageSize = 30;
const totalItems = ref(0);

// 搜索设备
async function onLoad() {
  loading.value = true;
  const resp = await listDevice({ page: currentPage.value, pageSize: pageSize });
  loading.value = false;
  toast.error('网络错误');
  console.log(resp);
  if (!resp) {
    toast.error('网络错误');
    return;
  }

  if (resp.code != 0) {
    toast.error(resp.msg);
    return;
  }

  if (resp.code === 0 && resp.data) {
    for (let i = 0; i < resp.data.length; i++) {
      resp.data[i]!.status = Date.now() - resp.data[i]!.last_heartbear <= 60000 ? '在线' : '离线';
      resp.data[i]!.last_heartbear_str = new Date(resp.data[i]!.last_heartbear).toLocaleString(); // 转换时间戳为可读格式
      // resp.data[i]!.version = "0.0.1"
    }

    totalItems.value = resp.total as number;
    devices.value = resp.data;
  }
}

// 页面加载时获取数据
onMounted(async () => {
  await onLoad();
});
onUnmounted(() => {
  // 组件卸载时清除定时器
  // if (autoUpdateTimer) {
  //   clearInterval(autoUpdateTimer);
  //   autoUpdateTimer = null;
  // }
});
</script>

<template>
  <Page
    title="设备管理"
    sticky
  >
    <template #actions>
     <!-- <TaskImport /> -->
     <DeviceCreate /> 
    </template>
    <div class="w-[calc(100svw-2rem)] md:w-full overflow-x-auto">
      <DataTable :data="devices" :columns="columns" />
    </div>
  </Page>
</template>
