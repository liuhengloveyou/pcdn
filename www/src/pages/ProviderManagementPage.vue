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
                    label="服务商名称/ID搜索"
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
                    v-model="levelFilter"
                    :options="levelOptions"
                    label="服务商等级"
                    dense
                    outlined
                    emit-value
                    map-options
                    clearable
                  />
                </div>
                <div class="col-md-6 col-sm-12 col-xs-12 text-right">
                  <q-btn color="primary" icon="add" label="添加服务商" @click="openAddDialog" />
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
                :rows="providers"
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
                    <q-td key="contact" :props="props">{{ props.row.contact }}</q-td>
                    <q-td key="phone" :props="props">{{ props.row.phone }}</q-td>
                    <q-td key="level" :props="props">
                      <q-chip :color="getLevelColor(props.row.level)" text-color="white" dense>
                        {{ props.row.level }}
                      </q-chip>
                    </q-td>
                    <q-td key="balance" :props="props">{{ props.row.balance }}</q-td>
                    <q-td key="joinDate" :props="props">{{ props.row.joinDate }}</q-td>
                    <q-td key="actions" :props="props">
                      <q-btn
                        flat
                        round
                        dense
                        color="primary"
                        icon="edit"
                        @click="editProvider(props.row)"
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
                        @click="viewProviderDetails(props.row)"
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

    <!-- 添加/编辑服务商对话框 -->
    <q-dialog v-model="addDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">{{ isEditing ? '编辑服务商' : '添加服务商' }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input
              v-model="form.name"
              label="服务商名称 *"
              :rules="[(val) => !!val || '服务商名称不能为空']"
              outlined
              dense
            />
            <q-input
              v-model="form.contact"
              label="联系人 *"
              :rules="[(val) => !!val || '联系人不能为空']"
              outlined
              dense
            />
            <q-input
              v-model="form.phone"
              label="联系电话 *"
              :rules="[(val) => !!val || '联系电话不能为空']"
              outlined
              dense
            />
            <q-input v-model="form.address" label="地址" outlined dense />
            <q-select
              v-model="form.level"
              :options="levelOptions"
              label="服务商等级 *"
              :rules="[(val) => !!val || '请选择服务商等级']"
              outlined
              dense
              emit-value
              map-options
            />
            <q-input v-model="form.balance" label="账户余额" type="number" outlined dense />
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
          <span class="q-ml-sm">确定要删除此服务商吗？此操作不可撤销。</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn flat label="删除" color="negative" @click="deleteProvider" v-close-popup />
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
  { name: 'id', align: 'left', label: '服务商ID', field: 'id', sortable: true },
  { name: 'name', align: 'left', label: '服务商名称', field: 'name', sortable: true },
  { name: 'contact', align: 'left', label: '联系人', field: 'contact', sortable: true },
  { name: 'phone', align: 'left', label: '联系电话', field: 'phone', sortable: true },
  { name: 'level', align: 'left', label: '服务商等级', field: 'level', sortable: true },
  { name: 'balance', align: 'left', label: '账户余额', field: 'balance', sortable: true },
  { name: 'joinDate', align: 'left', label: '加入日期', field: 'joinDate', sortable: true },
  { name: 'actions', align: 'center', label: '操作', field: 'actions' },
];

// 服务商等级选项
const levelOptions = [
  { label: '普通服务商', value: '普通服务商' },
  { label: '银牌服务商', value: '银牌服务商' },
  { label: '金牌服务商', value: '金牌服务商' },
  { label: '钻石服务商', value: '钻石服务商' },
];

// 表格数据
const providers = ref([
  {
    id: 1,
    name: '北京科技有限公司',
    contact: '张经理',
    phone: '13900139001',
    address: '北京市海淀区中关村',
    level: '普通服务商',
    balance: 5000.0,
    joinDate: '2023-01-10',
  },
  {
    id: 2,
    name: '上海网络科技有限公司',
    contact: '李总',
    phone: '13900139002',
    address: '上海市浦东新区张江高科技园区',
    level: '银牌服务商',
    balance: 15000.0,
    joinDate: '2023-02-15',
  },
  {
    id: 3,
    name: '广州数据服务有限公司',
    contact: '王经理',
    phone: '13900139003',
    address: '广州市天河区珠江新城',
    level: '金牌服务商',
    balance: 30000.0,
    joinDate: '2023-01-05',
  },
  {
    id: 4,
    name: '深圳云计算有限公司',
    contact: '赵总',
    phone: '13900139004',
    address: '深圳市南山区科技园',
    level: '钻石服务商',
    balance: 50000.0,
    joinDate: '2022-12-20',
  },
  {
    id: 5,
    name: '杭州互联网科技有限公司',
    contact: '钱经理',
    phone: '13900139005',
    address: '杭州市西湖区文三路',
    level: '普通服务商',
    balance: 8000.0,
    joinDate: '2023-03-01',
  },
]);

// 等级过滤器
const levelFilter = ref(null);

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
const currentProvider = ref(null);

// 表单数据
const form = reactive({
  id: null,
  name: '',
  contact: '',
  phone: '',
  address: '',
  level: '',
  balance: 0,
});

// 获取等级颜色
function getLevelColor(level: string) {
  switch (level) {
    case '普通服务商':
      return 'blue';
    case '银牌服务商':
      return 'grey';
    case '金牌服务商':
      return 'amber';
    case '钻石服务商':
      return 'purple';
    default:
      return 'blue';
  }
}

// 搜索服务商
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

// 编辑服务商
function editProvider(provider: any) {
  isEditing.value = true;
  currentProvider.value = provider;
  Object.assign(form, provider);
  addDialog.value = true;
}

// 确认删除
function confirmDelete(provider: any) {
  currentProvider.value = provider;
  deleteDialog.value = true;
}

// 删除服务商
function deleteProvider() {
  if (currentProvider.value) {
    // 实际应用中这里会调用API进行删除
    const index = providers.value.findIndex((p) => p.id === currentProvider.value.id);
    if (index !== -1) {
      providers.value.splice(index, 1);
    }
  }
}

// 查看服务商详情
function viewProviderDetails(provider: any) {
  // 实际应用中这里可能会跳转到详情页或打开详情对话框
  console.log('查看服务商详情:', provider);
}

// 提交表单
function onSubmit() {
  // 实际应用中这里会调用API进行保存
  if (isEditing.value && currentProvider.value) {
    // 更新现有服务商
    const index = providers.value.findIndex((p) => p.id === currentProvider.value.id);
    if (index !== -1) {
      providers.value[index] = { ...providers.value[index], ...form };
    }
  } else {
    // 添加新服务商
    const newId = Math.max(...providers.value.map((p) => p.id)) + 1;
    providers.value.push({
      ...form,
      id: newId,
      joinDate: new Date().toISOString().split('T')[0],
    });
  }
  addDialog.value = false;
}

// 重置表单
function resetForm() {
  form.id = null;
  form.name = '';
  form.contact = '';
  form.phone = '';
  form.address = '';
  form.level = '';
  form.balance = 0;
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
