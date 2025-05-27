<script lang="ts">
export const description = 'A two column login page with a cover image.'
</script>
<script setup lang="ts">
import { GalleryVerticalEnd } from 'lucide-vue-next'
import LoginForm from '@/components/LoginForm.vue'
import { onMounted, onUnmounted } from "vue";

defineOptions({
  name: "LoginPage",
});

onMounted(() => {
  // 重新登录
  document.cookie =
    "pcdn-sess=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;max-age=0";
});

onUnmounted(() => {
  // clearTimeout(smsTimer);
});


</script>

<template>
  <div class="grid min-h-svh lg:grid-cols-2">
    <div class="relative hidden bg-muted lg:block">
      <img
        src="/images/login-view.jpg"
        alt="Image"
        class="absolute inset-0 h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
      >
    </div>
    <div class="flex flex-col gap-4 p-6 md:p-10">
      <div class="flex justify-center gap-2 md:justify-start">
        <a href="#" class="flex items-center gap-2 font-medium">
          <div class="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
            <GalleryVerticalEnd class-name="size-4" />
          </div>
          智算末来PCDN系统
        </a>
      </div>
      <div class="flex flex-1 items-center justify-center">
        <div class="w-full max-w-xs">
          <LoginForm />
        </div>
      </div>
    </div>

  </div>
</template>
