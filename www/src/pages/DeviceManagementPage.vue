<template>
  <q-page padding>
    <div class="q-pa-md">
      <div class="row q-col-gutter-md q-mb-md">
        <!-- 搜索和添加按钮区域 -->
        <div class="col-12">
          <q-card flat bordered>
            <q-card-section>
              <div class="row q-col-gutter-md items-center">
                <div class="col-md-3 col-sm-6 col-xs-12">
                  <q-input
                    v-model="searchText"
                    label="设备名称/ID搜索"
                    dense
                    outlined
                    clearable
                    @keyup.enter="onSearch"
                  >
                    <template v-slot:append>
                      <q-icon name="search" @click="onSearch" class="cursor-pointer" />
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
        <div class="col-12">
          <q-card flat bordered>
            <q-card-section>
              <q-table
                :rows="devices"
                :columns="columns as QTableColumn[]"
                row-key="id"
                :loading="loading"
                :pagination.sync="pagination"
                :filter="searchText"
                binary-state-sort
              >
                <template v-slot:body="props">
                  <q-tr :props="props">
                    <q-td key="id" :props="props">{{ props.row.id }}</q-td>
                    <q-td key="name" :props="props">{{ props.row.name }}</q-td>
                    <q-td key="type" :props="props">{{ props.row.type }}</q-td>
                    <q-td key="status" :props="props">
                      <q-chip :color="getStatusColor(props.row.status)" text-color="white" dense>
                        {{ props.row.status }}
                      </q-chip>
                    </q-td>
                    <q-td key="ip" :props="props">{{ props.row.ip }}</q-td>
                    <q-td key="lastOnline" :props="props">{{ props.row.lastOnline }}</q-td>
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
                      <q-btn
                        flat
                        round
                        dense
                        color="info"
                        icon="visibility"
                        @click="viewDeviceDetails(props.row)"
                      />
                    </q-td>
                  </q-tr>
                </template>
              </q-table>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- 添加/编辑设备对话框 -->
    <q-dialog v-model="addDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">{{ isEditing ? '编辑设备' : '添加设备' }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input
              v-model="form.name"
              label="设备名称 *"
              :rules="[(val) => !!val || '设备名称不能为空']"
              outlined
              dense
            />
            <q-select
              v-model="form.type"
              :options="typeOptions"
              label="设备类型 *"
              :rules="[(val) => !!val || '请选择设备类型']"
              outlined
              dense
              emit-value
              map-options
            />
            <q-input v-model="form.ip" label="IP地址" outlined dense />
            <q-select
              v-model="form.status"
              :options="statusOptions"
              label="设备状态 *"
              :rules="[(val) => !!val || '请选择设备状态']"
              outlined
              dense
              emit-value
              map-options
            />
          </q-form>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn flat label="保存" color="primary" @click="onSubmit" />
        </q-card-actions>
      </q-card>
    </q-dialog>

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
import { QTableColumn } from 'quasar';
import { ref, reactive, onMounted } from 'vue';

// 表格列定义
const columns = [
  { name: 'id', align: 'left', label: '设备ID', field: 'id', sortable: true },
  { name: 'name', align: 'left', label: '设备名称', field: 'name', sortable: true },
  { name: 'type', align: 'left', label: '设备类型', field: 'type', sortable: true },
  { name: 'status', align: 'left', label: '状态', field: 'status', sortable: true },
  { name: 'ip', align: 'left', label: 'IP地址', field: 'ip', sortable: true },
  { name: 'lastOnline', align: 'left', label: '最后在线时间', field: 'lastOnline', sortable: true },
  { name: 'actions', align: 'center', label: '操作', field: 'actions' },
];

// 状态选项
const statusOptions = [
  { label: '在线', value: '在线' },
  { label: '离线', value: '离线' },
  { label: '维护中', value: '维护中' },
  { label: '故障', value: '故障' },
];

// 设备类型选项
const typeOptions = [
  { label: '路由器', value: '路由器' },
  { label: '交换机', value: '交换机' },
  { label: '服务器', value: '服务器' },
  { label: '其他', value: '其他' },
];

// 表格数据
const devices = ref([
  {
    id: 1,
    name: '设备-001',
    type: '路由器',
    status: '在线',
    ip: '192.168.1.1',
    lastOnline: '2023-05-15 14:30:22',
  },
  {
    id: 2,
    name: '设备-002',
    type: '交换机',
    status: '离线',
    ip: '192.168.1.2',
    lastOnline: '2023-05-14 09:15:43',
  },
  {
    id: 3,
    name: '设备-003',
    type: '服务器',
    status: '维护中',
    ip: '192.168.1.3',
    lastOnline: '2023-05-15 11:22:05',
  },
  {
    id: 4,
    name: '设备-004',
    type: '路由器',
    status: '故障',
    ip: '192.168.1.4',
    lastOnline: '2023-05-13 16:45:30',
  },
  {
    id: 5,
    name: '设备-005',
    type: '服务器',
    status: '在线',
    ip: '192.168.1.5',
    lastOnline: '2023-05-15 15:10:18',
  },
]);

// 状态过滤器
const statusFilter = ref(null);

// 搜索文本
const searchText = ref('');

// 加载状态
const loading = ref(false);

// 分页设置
const pagination = ref({
  sortBy: 'id',
  descending: false,
  page: 1,
  rowsPerPage: 10,
  rowsNumber: 10,
});

// 对话框控制
const addDialog = ref(false);
const deleteDialog = ref(false);
const isEditing = ref(false);
const currentDevice = ref(null);

// 表单数据
const form = reactive({
  id: null,
  name: '',
  type: '',
  status: '',
  ip: '',
});

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

// 搜索设备
function onSearch() {
  // 实际应用中这里会调用API进行搜索
  loading.value = true;
  setTimeout(() => {
    loading.value = false;
  }, 500);
}

// 打开添加对话框
function openAddDialog() {
  isEditing.value = false;
  resetForm();
  addDialog.value = true;
}

// 编辑设备
function editDevice(device: any) {
  isEditing.value = true;
  currentDevice.value = device;
  Object.assign(form, device);
  addDialog.value = true;
}

// 确认删除
function confirmDelete(device: any) {
  currentDevice.value = device;
  deleteDialog.value = true;
}

// 删除设备
function deleteDevice() {
  if (currentDevice.value) {
    // 实际应用中这里会调用API进行删除
    // const index = currentDevice.value?.id ? devices.value.findIndex((d) => d.id === currentDevice.value.id) : -1;
    // if (index !== -1) {
    //   devices.value.splice(index, 1);
    // }
  }
}

// 查看设备详情
function viewDeviceDetails(device: any) {
  // 实际应用中这里可能会跳转到详情页或打开详情对话框
  console.log('查看设备详情:', device);
}

// 提交表单
function onSubmit() {
  // // 实际应用中这里会调用API进行保存
  // if (isEditing.value && currentDevice.value) {
  //   // 更新现有设备
  //   const index = devices.value.findIndex((d) => d.id === currentDevice.value.id);
  //   if (index !== -1) {
  //     devices.value[index] = { ...devices.value[index], ...form };
  //   }
  // } else {
  //   // 添加新设备
  //   const newId = Math.max(...devices.value.map((d) => d.id)) + 1;
  //   devices.value.push({
  //     ...form,
  //     id: newId,
  //     lastOnline: new Date().toLocaleString(),
  //   });
  // }
  // addDialog.value = false;
}

// 重置表单
function resetForm() {
  form.id = null;
  form.name = '';
  form.type = '';
  form.status = '';
  form.ip = '';
}

// 页面加载时获取数据
onMounted(() => {
  // 实际应用中这里会调用API获取数据
  loading.value = true;
  setTimeout(() => {
    loading.value = false;
  }, 500);
});
</script>
