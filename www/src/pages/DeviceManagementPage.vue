<template>
  <q-page class="q-pa-sm">
    <div class="row q-col-gutter-md">
      <!-- 搜索和添加按钮区域 -->
      <div class="col-12">
        <q-card flat bordered>
          <q-card-section style="padding: 0">
            <div class="row q-col-gutter-md items-center">
              <div class="col-md-3 col-sm-6 col-xs-12">
                <q-input
                  v-model="searchText"
                  label="设备名称/ID搜索"
                  dense
                  outlined
                  clearable
                  @keyup.enter="onLoad"
                >
                  <template v-slot:append>
                    <q-icon name="search" @click="onLoad" class="cursor-pointer" />
                  </template>
                </q-input>
              </div>
              <div class="col-md-3 col-sm-6 col-xs-12">
                <q-select
                  v-model="statusFilter"
                  :options="statusOptions"
                  label="设备状态"
                  dense
                  outlined
                  emit-value
                  map-options
                  clearable
                />
              </div>
              <div class="col-md-6 col-sm-12 col-xs-12 text-right">
                <q-btn color="primary" icon="add" label="添加设备" @click="openAddDialog" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- 数据表格区域 -->
      <q-card flat style="flex: 1; display: flex; flex-direction: column">
        <q-card-section style="padding: 0">
          <q-table
            flat
            hide-bottom
            :rows="devices"
            :columns="columns"
            row-key="id"
            :loading="loading"
            :filter="searchText"
            virtual-scroll
          >
            <template v-slot:body="props">
              <q-tr :props="props">
                <q-td key="id" :props="props">{{ props.row.id }}</q-td>
                <q-td key="sn" :props="props">{{ props.row.sn }}</q-td>

                <q-td key="status" :props="props">
                  <q-chip :color="getStatusColor(props.row.status)" text-color="white" dense>
                    {{ props.row.status }}
                  </q-chip>
                </q-td>
                <q-td key="ip" :props="props">{{ props.row.remote_addr }}</q-td>
                <q-td key="lastOnline" :props="props">{{ props.row.last_heartbear_str }}</q-td>
                <q-td key="version" :props="props">{{ props.row.version }}</q-td>
                <q-td key="actions" :props="props">
                  <q-btn
                    flat
                    round
                    dense
                    color="primary"
                    icon="edit"
                    @click="editDevice(props.row)"
                  />
                  <q-btn
                    flat
                    round
                    dense
                    color="negative"
                    icon="delete"
                    @click="confirmDelete(props.row)"
                  />
                </q-td>
              </q-tr>
            </template>
          </q-table>
        </q-card-section>
      </q-card>
    </div>

    <q-footer>
      <q-toolbar class="bg-grey-1">
        <q-pagination v-model="currentPage" :max="totalItems / pageSize + 1" input />
        <q-space></q-space>
        <q-toggle v-model="autoUpdateEnabled" label="自动刷新" @update:model-value="handleAutoUpdate" style="color: black;" />
      </q-toolbar>
    </q-footer>
    <!-- 添加/编辑设备对话框 -->
    <q-drawer v-model="addDrawer" side="right" :width="600" overlay elevated persistent>
      <q-card flat>
        <q-card-section>
          <div class="text-h6">{{ isEditing ? '编辑设备' : '添加设备' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input
              v-model="form.sn"
              label="设备SN *"
              :rules="[(val) => !!val || '设备SN不能为空']"
              outlined
            />
          </q-form>
        </q-card-section>
        <q-card-actions class="q-gutter-sm q-mt-md row full-width">
          <q-btn
            outline
            label="取消"
            color="primary"
            @click="addDrawer = false"
            class="col q-py-sm text-subtitle1"
          />
          <q-btn
            unelevated
            label="保存"
            color="primary"
            @click="onSubmit"
            class="col q-ml-sm q-py-sm text-subtitle1"
          />
        </q-card-actions>
      </q-card>
    </q-drawer>

    <!-- 删除确认对话框 -->
    <q-dialog v-model="deleteDialog" persistent>
      <q-card>
        <q-card-section class="row items-center">
          <q-avatar icon="warning" color="negative" text-color="white" />
          <span class="q-ml-sm">确定要删除此设备吗？此操作不可撤销。</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn flat label="删除" color="negative" @click="deleteDevice" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
// eslint-disable @typescript-eslint/no-explicit-any
import { useQuasar, type QTableColumn } from 'quasar';
import type { DeviceModel } from 'src/service/DeviceService';
import { addDevice, listDevice } from 'src/service/DeviceService';
import { ref, reactive, onMounted, onUnmounted } from 'vue';

const $q = useQuasar();

// 对话框控制
const addDrawer = ref(false);
const deleteDialog = ref(false);
const isEditing = ref(false);
const currentDevice = ref<DeviceModel | null>(null);
// 自动更新控制
const autoUpdateEnabled = ref(false);
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let autoUpdateTimer:any = null;

// 分页控制
const currentPage = ref(1);
const pageSize = 30;
const totalItems = ref(0);

// 表单数据
const form = reactive({
  id: 0,
  uid: 0,
  sn: '',
  createTime: 0,
  updateTime: 0,
  remote_addr: '',
  version: '',
  timestamp: 0,
  last_heartbear: 0,
});

// 表格数据
const devices = ref<DeviceModel[]>([]);
// 表格列定义
const columns = [
  { name: 'id', align: 'left', label: '设备ID', field: 'id', sortable: true },
  { name: 'sn', align: 'left', label: '设备SN', field: 'sn', sortable: true },
  { name: 'status', align: 'left', label: '状态', field: 'status', sortable: true },
  { name: 'ip', align: 'left', label: 'IP地址', field: 'ip', sortable: true },
  { name: 'lastOnline', align: 'left', label: '最后在线时间', field: 'lastOnline', sortable: true },
  { name: 'version', align: 'left', label: '版本', field: 'version' },
  { name: 'actions', align: 'center', label: '操作', field: 'actions' },
] as QTableColumn[];
// 状态选项
const statusOptions = [
  { label: '在线', value: '在线' },
  { label: '离线', value: '离线' },
  { label: '维护中', value: '维护中' },
  { label: '故障', value: '故障' },
];

// 状态过滤器
const statusFilter = ref(null);

// 搜索文本
const searchText = ref('');

// 加载状态
const loading = ref(false);

// 获取状态颜色
function getStatusColor(status: string) {
  switch (status) {
    case '在线':
      return 'positive';
    case '离线':
      return 'grey';
    case '维护中':
      return 'warning';
    case '故障':
      return 'negative';
    default:
      return 'grey';
  }
}

// 打开添加对话框
function openAddDialog() {
  isEditing.value = false;
  resetForm();
  addDrawer.value = true;
}

// 编辑设备
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function editDevice(device: any) {
  isEditing.value = true;
  currentDevice.value = device;
  Object.assign(form, device);
  addDrawer.value = true;
}

// 确认删除
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function confirmDelete(device: any) {
  currentDevice.value = device;
  deleteDialog.value = true;
}

// 删除设备
function deleteDevice() {
  if (currentDevice.value) {
    // 实际应用中这里会调用API进行删除
    const index = devices.value.findIndex((d) => d.id === currentDevice.value?.id);
    if (index !== -1) {
      devices.value.splice(index, 1);
    }
  }
}

// 提交表单
async function onSubmit() {
  if (isEditing.value && currentDevice.value) {
    // 更新现有设备
    // const index = devices.value.findIndex((d) => d.id === currentDevice.value.id);
    // if (index !== -1) {
    //   devices.value[index] = { ...devices.value[index], ...form };
    // }
  } else {
    // 添加新设备
    const resp = await addDevice({
      id: 0,
      sn: form.sn,
      uid: 0,
      createTime: 0,
      updateTime: 0,
      remote_addr: '',
      version: '',
      timestamp: 0,
      last_heartbear: 0,
    });
    console.log(resp);
  }

  resetForm();
  addDrawer.value = false;

  await onLoad();
}

// 重置表单
function resetForm() {
  form.id = 0;
  form.sn = '';
  form.uid = 0;
  form.createTime = 0;
  form.updateTime = 0;
  form.remote_addr = '';
  form.version = '';
  form.timestamp = 0;
  form.last_heartbear = 0;
}

function handleAutoUpdate() {
  if (autoUpdateEnabled.value) {
    // 启用自动更新
    autoUpdateTimer = setInterval(() => {
      onLoad().then(() => {
        console.log('Auto update completed.');
      }).catch((error) => {
        // 处理错误
        console.error('Auto update error:', error);
      });
    }, 10000);
  } else {
    // 禁用自动更新，清除定时器
    if (autoUpdateTimer) {
      clearInterval(autoUpdateTimer);
      autoUpdateTimer = null;
    }
  }
}

// 搜索设备
async function onLoad() {
  loading.value = true;
  const resp = await listDevice({ page: currentPage, pageSize: pageSize });
  loading.value = false;

  console.log(resp);
  if (!resp) {
    $q.notify({
      color: 'negative',
      textColor: 'white',
      icon: 'report_problem',
      message: '网络错误',
    });
    return;
  }

  if (resp.code != 0) {
    $q.notify(resp.msg);
    return;
  }

  if (resp.code === 0 && resp.data) {
    for (let i = 0; i < resp.data.length; i++) {
      resp.data[i]!.status = Date.now() - resp.data[i]!.last_heartbear <= 60000 ? '在线' : '离线';

      resp.data[i]!.last_heartbear_str = new Date(resp.data[i]!.last_heartbear).toLocaleString(); // 转换时间戳为可读格式
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
  if (autoUpdateTimer) {
    clearInterval(autoUpdateTimer);
    autoUpdateTimer = null;
  }
});
</script>
