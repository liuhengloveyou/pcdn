<template>
  <q-page class="flex flex-center bg-grey-2">
    <q-card class="login-card q-pa-lg">
      <q-card-section class="text-center">
        <div class="text-h5 q-mb-md">系统登录</div>
      </q-card-section>

      <q-form @submit="onSubmit" class="q-gutter-md">
        <q-input
          v-model="username"
          label="用户名"
          outlined
          :rules="[val => !!val || '请输入用户名']"
        >
          <template v-slot:prepend>
            <q-icon name="person" />
          </template>
        </q-input>

        <q-input
          v-model="password"
          label="密码"
          outlined
          :type="isPwd ? 'password' : 'text'"
          :rules="[val => !!val || '请输入密码']"
        >
          <template v-slot:prepend>
            <q-icon name="lock" />
          </template>
          <template v-slot:append>
            <q-icon
              :name="isPwd ? 'visibility_off' : 'visibility'"
              class="cursor-pointer"
              @click="isPwd = !isPwd"
            />
          </template>
        </q-input>

        <div class="flex justify-between items-center q-mt-sm">
          <q-checkbox v-model="rememberMe" label="记住我" />
          <q-btn flat color="primary" label="忘记密码？" no-caps />
        </div>

        <div class="q-mt-lg">
          <q-btn
            type="submit"
            color="primary"
            class="full-width"
            label="登录"
            :loading="loading"
          />
        </div>
      </q-form>
    </q-card>
  </q-page>
</template>

<script>
import { defineComponent, ref } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'LoginPage',
  setup () {
    const $q = useQuasar()
    const router = useRouter()

    const username = ref('')
    const password = ref('')
    const isPwd = ref(true)
    const rememberMe = ref(false)
    const loading = ref(false)

    const onSubmit = async () => {
      loading.value = true

      try {
        // 这里添加实际的登录API调用
        // 例如: const response = await api.login(username.value, password.value)

        // 模拟登录延迟
        await new Promise(resolve => setTimeout(resolve, 1000))

        // 登录成功处理
        $q.notify({
          color: 'positive',
          message: '登录成功',
          icon: 'check_circle'
        })

        // 保存登录状态（如果需要）
        if (rememberMe.value) {
          localStorage.setItem('username', username.value)
          // 注意：实际项目中不要在localStorage中存储密码，这里只是示例
        }

        // 跳转到首页或者之前的页面
        router.push('/')
      } catch (error) {
        // 登录失败处理
        $q.notify({
          color: 'negative',
          message: '登录失败: ' + (error.message || '用户名或密码错误'),
          icon: 'error'
        })
      } finally {
        loading.value = false
      }
    }

    return {
      username,
      password,
      isPwd,
      rememberMe,
      loading,
      onSubmit
    }
  }
})
</script>

<style lang="scss" scoped>
.login-card {
  width: 100%;
  max-width: 400px;
  border-radius: 8px;
}
</style>
