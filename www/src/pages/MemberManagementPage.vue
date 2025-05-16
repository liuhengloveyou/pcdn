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
                    label="会员名称/ID搜索"
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
                    label="会员等级"
                    dense
                    outlined
                    emit-value
                    map-options
                    clearable
                  />
                </div>
                <div class="col-md-6 col-sm-12 col-xs-12 text-right">
                  <q-btn color="primary" icon="add" label="添加会员" @click="openAddDialog" />
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
                :rows="members"
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
                    <q-td key="phone" :props="props">{{ props.row.phone }}</q-td>
                    <q-td key="email" :props="props">{{ props.row.email }}</q-td>
                    <q-td key="level" :props="props">
                      <q-chip :color="getLevelColor(props.row.level)" text-color="white" dense>
                        {{ props.row.level }}
                      </q-chip>
                    </q-td>
                    <q-td key="balance" :props="props">{{ props.row.balance }}</q-td>
                    <q-td key="registerDate" :props="props">{{ props.row.registerDate }}</q-td>
                    <q-td key="actions" :props="props">
                      <q-btn
                        flat
                        round
                        dense
                        color="primary"
                        icon="edit"
                        @click="editMember(props.row)"
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
                        @click="viewMemberDetails(props.row)"
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

    <!-- 添加/编辑会员对话框 -->
    <q-dialog v-model="addDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">{{ isEditing ? '编辑会员' : '添加会员' }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input
              v-model="form.name"
              label="会员名称 *"
              :rules="[(val) => !!val || '会员名称不能为空']"
              outlined
              dense
            />
            <q-input
              v-model="form.phone"
              label="手机号码 *"
              :rules="[(val) => !!val || '手机号码不能为空']"
              outlined
              dense
            />
            <q-input v-model="form.email" label="电子邮箱" type="email" outlined dense />
            <q-select
              v-model="form.level"
              :options="levelOptions"
              label="会员等级 *"
              :rules="[(val) => !!val || '请选择会员等级']"
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
          <span class="q-ml-sm">确定要删除此会员吗？此操作不可撤销。</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn flat label="删除" color="negative" @click="deleteMember" v-close-popup />
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
  { name: 'id', align: 'left', label: '会员ID', field: 'id', sortable: true },
  { name: 'name', align: 'left', label: '会员名称', field: 'name', sortable: true },
  { name: 'phone', align: 'left', label: '手机号码', field: 'phone', sortable: true },
  { name: 'email', align: 'left', label: '电子邮箱', field: 'email', sortable: true },
  { name: 'level', align: 'left', label: '会员等级', field: 'level', sortable: true },
  { name: 'balance', align: 'left', label: '账户余额', field: 'balance', sortable: true },
  { name: 'registerDate', align: 'left', label: '注册日期', field: 'registerDate', sortable: true },
  { name: 'actions', align: 'center', label: '操作', field: 'actions' },
];

// 会员等级选项
const levelOptions = [
  { label: '普通会员', value: '普通会员' },
  { label: '银卡会员', value: '银卡会员' },
  { label: '金卡会员', value: '金卡会员' },
  { label: '钻石会员', value: '钻石会员' },
];

// 表格数据
const members = ref([
  {
    id: 1,
    name: '张三',
    phone: '13800138001',
    email: 'zhangsan@example.com',
    level: '普通会员',
    balance: 100.0,
    registerDate: '2023-01-15',
  },
  {
    id: 2,
    name: '李四',
    phone: '13800138002',
    email: 'lisi@example.com',
    level: '银卡会员',
    balance: 500.0,
    registerDate: '2023-02-20',
  },
  {
    id: 3,
    name: '王五',
    phone: '13800138003',
    email: 'wangwu@example.com',
    level: '金卡会员',
    balance: 1200.0,
    registerDate: '2023-03-10',
  },
  {
    id: 4,
    name: '赵六',
    phone: '13800138004',
    email: 'zhaoliu@example.com',
    level: '钻石会员',
    balance: 5000.0,
    registerDate: '2023-01-05',
  },
  {
    id: 5,
    name: '钱七',
    phone: '13800138005',
    email: 'qianqi@example.com',
    level: '普通会员',
    balance: 50.0,
    registerDate: '2023-04-25',
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
const currentMember = ref(null);

// 表单数据
const form = reactive({
  id: null,
  name: '',
  phone: '',
  email: '',
  level: '',
  balance: 0,
});

// 获取等级颜色
function getLevelColor(level: string) {
  switch (level) {
    case '普通会员':
      return 'blue';
    case '银卡会员':
      return 'grey';
    case '金卡会员':
      return 'amber';
    case '钻石会员':
      return 'purple';
    default:
      return 'blue';
  }
}

// 搜索会员
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

// 编辑会员
function editMember(member: any) {
  isEditing.value = true;
  currentMember.value = member;
  Object.assign(form, member);
  addDialog.value = true;
}

// 确认删除
function confirmDelete(member: any) {
  currentMember.value = member;
  deleteDialog.value = true;
}

// 删除会员
function deleteMember() {
  if (currentMember.value) {
    // 实际应用中这里会调用API进行删除
    const index = members.value.findIndex((m) => m.id === currentMember.value.id);
    if (index !== -1) {
      members.value.splice(index, 1);
    }
  }
}

// 查看会员详情
function viewMemberDetails(member: any) {
  // 实际应用中这里可能会跳转到详情页或打开详情对话框
  console.log('查看会员详情:', member);
}

// 提交表单
function onSubmit() {
  // 实际应用中这里会调用API进行保存
  if (isEditing.value && currentMember.value) {
    // 更新现有会员
    const index = members.value.findIndex((m) => m.id === currentMember.value.id);
    if (index !== -1) {
      members.value[index] = { ...members.value[index], ...form };
    }
  } else {
    // 添加新会员
    const newId = Math.max(...members.value.map((m) => m.id)) + 1;
    members.value.push({
      ...form,
      id: newId,
      registerDate: new Date().toISOString().split('T')[0],
    });
  }
  addDialog.value = false;
}

// 重置表单
function resetForm() {
  form.id = null;
  form.name = '';
  form.phone = '';
  form.email = '';
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
