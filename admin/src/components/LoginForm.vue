<script setup lang="ts">
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AccountService } from "@/services/AccountService";
import { useUserSessionStore } from "@/stores/userSessionStore";
import { useRouter } from "vue-router";
import { onMounted, onUnmounted, ref } from "vue";
import { Toaster } from "@/components/ui/sonner";
import { toast } from "vue-sonner";
import { toTypedSchema } from "@vee-validate/zod";
import z from "zod";
import { useForm } from "vee-validate";

defineOptions({
  name: "LoginForm",
});

const cellphone = ref("");
const password = ref("");
const isLoading = ref(false);
const sessionStore = useUserSessionStore();
const router = useRouter();

onMounted(() => {
  // 重新登录
  document.cookie =
    "trade-sess=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;max-age=0";
});

onUnmounted(() => {
  // clearTimeout(smsTimer);
});


const formSchema = toTypedSchema(
  z.object({
    cellphone: z.string().length(11),
    password: z.string().min(6).max(50),
  })
);

const { handleSubmit } = useForm({
  validationSchema: formSchema,
});

const onSubmit = handleSubmit((values: { cellphone: string; password: string; }) => {
  // toast({
  //   title: "You submitted the following values:",
  //   description: h(
  //     "pre",
  //     { class: "mt-2 w-[340px] rounded-md bg-slate-950 p-4" },
  //     h("code", { class: "text-white" }, JSON.stringify(values, null, 2))
  //   ),
  // });

  onRealSubmit(values.cellphone, values.password);
});

async function onRealSubmit(cellphone: string, password: string) {
  // 表单验证
  if (cellphone.trim() === "") {
    toast("参数错误", {
      description: "请输入手机号码",
    });

    return;
  }
  if (password.trim() === "") {
    toast("参数错误", {
      description: "请输入密码",
    });
    return;
  }

  isLoading.value = true;

  try {
    // 调用登录API
    const resp = await AccountService.login(
      cellphone.trim(),
      password.trim(),
      ""
    );
    if (resp.code !== 0) {
      // 登录失败
      toast("错误", {
        description: resp.msg || "登录失败，请检查账号和密码",
      });
    } else {
      // 登录成功，保存用户信息到store
      console.log("登录成功:", resp.data);
      sessionStore.login(resp.data);
      // 跳转到首页
      router.replace({ path: "/" });
    }
  } catch (error) {
    console.error("登录出错:", error);
    toast("错误", {
        description: "登录过程中发生错误，请稍后再试",
      });
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <form :class="cn('flex flex-col gap-6')" @submit="onSubmit">
    <div class="flex flex-col items-center gap-2 text-center">
      <h1 class="text-2xl font-bold">欢迎登录智算PCDN系统</h1>
      <p class="text-balance text-sm text-muted-foreground">
        在下方输入您的手机号码和密码以登录您的账户
      </p>
    </div>
    <div class="grid gap-6">
      <div class="grid gap-2">
        <Label for="cellphone">手机号码:</Label>
        <Input
          id="cellphone"
          v-model="cellphone"
          type="text"
          placeholder="请输入您的手机号码"
          required
        />
      </div>
      <div class="grid gap-2">
        <div class="flex items-center">
          <Label for="password">密码</Label>
          <a
            href="#"
            class="ml-auto text-sm underline-offset-4 hover:underline"
          >
            忘记密码？
          </a>
        </div>
        <Input id="password" v-model="password" type="password" required />
      </div>
      <Button type="submit" class="w-full" :disabled="isLoading">
        {{ isLoading ? "登录中..." : "登录" }}
      </Button>
    </div>
    <div class="text-center text-sm">
      还没有账户？
      <a href="#" class="underline underline-offset-4"> 注册 </a>
    </div>
  </form>

  <Toaster />
</template>
