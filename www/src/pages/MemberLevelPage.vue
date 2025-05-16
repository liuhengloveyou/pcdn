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
                    label="级别名称搜索"
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
                <div class="col-md-9 col-sm-12 col-xs-12 text-right">
                  <q-btn color="primary" icon="add" label="添加会员级别" @click="openAddDialog" />
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
                :rows="memberLevels"
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
                    <q-td key="commissionRate" :props="props">{{ props.row.commissionRate }}%</q-td>
                    <q-td key="minPoints" :props="props">{{ props.row.minPoints }}</q-td>
                    <q-td key="description" :props="props">{{ props.row.description }}</q-td>
                    <q-td key="actions" :props="props">
                      <q-btn
                        flat
                        round
                        dense
                        color="primary"
                        icon="edit"
                        @click="editLevel(props.row)"
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
      </div>
    </div>

    <!-- 添加/编辑会员级别对话框 -->
    <q-dialog v-model="addEditDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">{{ isEditing ? '编辑会员级别' : '添加会员级别' }}</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input
              v-model="form.name"
              label="级别名称 *"
              :rules="[(val) => !!val || '级别名称不能为空']"
              outlined
              dense
            />
            <q-input
              v-model.number="form.commissionRate"
              label="提成比例(%) *"
              type="number"
              :rules="[
                (val) => !!val || '提成比例不能为空',
                (val) => val >= 0 || '提成比例不能为负',
                (val) => val <= 100 || '提成比例不能超过100%',
              ]"
              outlined
              dense
            />
            <q-input
              v-model.number="form.minPoints"
              label="最低积分要求 *"
              type="number"
              :rules="[
                (val) => !!val || '最低积分要求不能为空',
                (val) => val >= 0 || '最低积分要求不能为负',
              ]"
              outlined
              dense
            />
            <q-input v-model="form.description" label="级别描述" type="textarea" outlined dense />
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
          <span class="q-ml-sm">确定要删除此会员级别吗？此操作不可撤销。</span>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="primary" v-close-popup />
          <q-btn flat label="删除" color="negative" @click="deleteLevel" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { QTableColumn, useQuasar } from 'quasar';

const $q = useQuasar();

// 表格列定义
const columns = [
  { name: 'id', align: 'left', label: 'ID', field: 'id', sortable: true },
  { name: 'name', align: 'left', label: '级别名称', field: 'name', sortable: true },
  {
    name: 'commissionRate',
    align: 'left',
    label: '提成比例',
    field: 'commissionRate',
    sortable: true,
  },
  { name: 'minPoints', align: 'left', label: '最低积分要求', field: 'minPoints', sortable: true },
  { name: 'description', align: 'left', label: '描述', field: 'description' },
  { name: 'actions', align: 'center', label: '操作', field: 'actions' },
];

// 模拟数据
const memberLevels = ref([
  {
    id: 1,
    name: '普通会员',
    commissionRate: 5,
    minPoints: 0,
    description: '基础会员级别',
  },
  {
    id: 2,
    name: '银牌会员',
    commissionRate: 8,
    minPoints: 1000,
    description: '中级会员待遇',
  },
  {
    id: 3,
    name: '金牌会员',
    commissionRate: 12,
    minPoints: 5000,
    description: '高级会员待遇',
  },
  {
    id: 4,
    name: 'VIP会员',
    commissionRate: 15,
    minPoints: 10000,
    description: '尊贵会员待遇',
  },
]);

// 状态变量
const searchText = ref('');
const loading = ref(false);
const pagination = ref({
  sortBy: 'id',
  descending: false,
  page: 1,
  rowsPerPage: 10,
});

// 对话框状态
const addEditDialog = ref(false);
const deleteDialog = ref(false);
const isEditing = ref(false);
const currentLevel = ref(null);

// 表单数据
const form = reactive({
  id: 0,
  name: '',
  commissionRate: 0,
  minPoints: 0,
  description: '',
});

// 搜索方法
function onSearch() {
  loading.value = true;
  // 这里应该是实际的API调用
  setTimeout(() => {
    loading.value = false;
  }, 500);
}

// 打开添加对话框
function openAddDialog() {
  isEditing.value = false;
  resetForm();
  addEditDialog.value = true;
}

// 编辑会员级别
function editLevel(level: null) {
  isEditing.value = true;
  currentLevel.value = level;
  Object.assign(form, level);
  addEditDialog.value = true;
}

// 确认删除
function confirmDelete(level: null) {
  currentLevel.value = level;
  deleteDialog.value = true;
}

// 删除会员级别
function deleteLevel() {
  if (!currentLevel.value) return;

  loading.value = true;
  // 这里应该是实际的API调用
  setTimeout(() => {
    if (!currentLevel.value) return;

    const index = memberLevels.value.findIndex((l) => l.id === currentLevel.value?.id);
    if (index !== -1) {
      memberLevels.value.splice(index, 1);
    }
    loading.value = false;
    $q.notify({
      color: 'positive',
      message: '会员级别已删除',
      icon: 'check',
    });
  }, 500);
}

// 提交表单
function onSubmit() {
  loading.value = true;
  // 这里应该是实际的API调用
  setTimeout(() => {
    if (isEditing.value) {
      // 更新现有级别
      const index = memberLevels.value.findIndex((l) => l.id === form.id);
      if (index !== -1) {
        memberLevels.value[index] = { ...form };
      }
      $q.notify({
        color: 'positive',
        message: '会员级别已更新',
        icon: 'check',
      });
    } else {
      // 添加新级别
      const newId = Math.max(0, ...memberLevels.value.map((l) => l.id)) + 1;
      memberLevels.value.push({
        ...form,
        id: newId,
      });
      $q.notify({
        color: 'positive',
        message: '会员级别已添加',
        icon: 'check',
      });
    }
    loading.value = false;
    addEditDialog.value = false;
  }, 500);
}

// 重置表单
function resetForm() {
  form.id = 0;
  form.name = '';
  form.commissionRate = 0;
  form.minPoints = 0;
  form.description = '';
}
</script>
