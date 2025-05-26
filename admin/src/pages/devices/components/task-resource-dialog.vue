<script lang="ts" setup>
import { computed } from "vue";
import { DialogHeader, DialogTitle } from "@/components/ui/dialog";
import {
  FormControl,
  FormDescription,
  FormField,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import * as z from "zod";
import FormItem from "@/components/ui/form/FormItem.vue";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { addDevice, type DeviceModel } from "@/services/DeviceService";
import type { Device } from "./columns";

const props = defineProps<{
  task: Device | null;
}>();
defineEmits(["close"]);

const task = computed(() => props.task);
const title = computed(() => (task.value?.id ? `编程设备` : "添加设备"));

const formSchema = toTypedSchema(
  z.object({
    sn: z
      .string()
      .min(2)
      .max(50)
      .default(props.task?.sn ?? ""),
  })
);

const { isFieldDirty, handleSubmit } = useForm({
  validationSchema: formSchema,
});
const onSubmit = handleSubmit((values) => {
  console.log(values);
  onRealSubmit({
    id: 0,
    uid: 0,
    sn: values.sn,
    createTime: 0,
    updateTime: 0,
    remote_addr: "",
    version: "",
    timestamp: 0,
    last_heartbear: 0,
  });
});

// 提交表单
async function onRealSubmit(val: DeviceModel) {
  // if (isEditing.value && currentDevice.value) {
  //   // 更新现有设备
  //   // const index = devices.value.findIndex((d) => d.id === currentDevice.value.id);
  //   // if (index !== -1) {
  //   //   devices.value[index] = { ...devices.value[index], ...form };
  //   // }
  // } else {
  // 添加新设备
  const resp = await addDevice(val);
  console.log(resp);
  // }

  // addDrawer.value = false;

  // await onLoad();
}
</script>

<template>
  <div>
    <DialogHeader>
      <DialogTitle>
        {{ title }}
      </DialogTitle>
    </DialogHeader>
    <form class="space-y-6 mt-4" @submit="onSubmit">
      <FormField
        v-slot="{ componentField }"
        name="sn"
        :validate-on-blur="!isFieldDirty"
      >
        <FormItem>
          <FormLabel>设备编码(SN): </FormLabel>
          <FormControl>
            <Input type="text" v-bind="componentField" />
          </FormControl>
          <FormDescription />
          <FormMessage />
        </FormItem>
      </FormField>

      <Button type="submit" class="w-full"> 提交 </Button>
    </form>
  </div>
</template>

<style scoped></style>
