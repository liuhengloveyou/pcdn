<script setup lang="ts">
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { useToast } from "@/components/ui/toast/use-toast";
const { toast } = useToast();
import { toTypedSchema } from "@vee-validate/zod";
import { GenericObject, useForm } from "vee-validate";
import { onMounted, ref } from "vue";
import * as z from "zod";
import {
  BotModel,
  BotService,
  MinimumBrushRobotModel,
} from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

// 最低维持刷量
defineOptions({
  name: "MinimumBrushRobotForm",
});
const emits = defineEmits(["commited"]);

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
let loading = ref(false);
let model = ref<MinimumBrushRobotModel>({
  id: 0,

  // 价格模式：
  priceMode: "normalMode", // 一般模式 / 安全模式
  // 启动模式：
  launchMode: "immediateStart", // 立即启动 / 定时启动

  // 定时启动时间
  launchTime: "",

  // 机器人描述：
  description: "",
  isRunning: false,
  startTime: "",
});

const formSchema = toTypedSchema(
  z.object({
    priceMode: z.enum(["normalMode", "safeMode"], {
      required_error: "You need to select a PriceMode.",
    }),
    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a LaunchMode.",
    }),
    launchTime: z.string().optional(),
  })
);
const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    priceMode: "normalMode",
    launchMode: "immediateStart",
  },
});

onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await BotService.LoadOne(
    1,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    if (resp.code == 0) {
      if (resp.data && resp.data.minimumBrushBot) {
        model.value = resp.data.minimumBrushBot;
        setFieldValue("priceMode", model.value.priceMode);
        setFieldValue("launchMode", model.value.launchMode);
        setFieldValue("launchTime", model.value.launchTime);
      }
    }
  }

  console.log("MinimumBrushRobotModel>>>bot: ", model);
});

const onSubmit = handleSubmit((values) => {
  onRealSubmit(values);
});

async function onRealSubmit(m: GenericObject) {
  const activeMartParam = martParamStore.getActiveMartParam;

  if (activeMartParam == null) {
    return;
  }

  let data = { ...m.value };
  data.priceMode = m.priceMode;
  data.launchMode = m.launchMode;
  data.launchTime = m.launchTime;

  let botModel: BotModel = {
    id: 0,
    botType: 1,
    botName: "MinimumMaintenanceBrush",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    minimumBrushBot: data,
  };
  console.log(">>>", botModel);
  loading.value = true;
  const resp = await BotService.Set(botModel);
  loading.value = false;
  console.log("resp: ", resp);
  if (!resp) {
    return;
  }

  if (resp?.code != 0) {
    toast({
      variant: "destructive",
      title: t("error." + resp.code),
    });
    return;
  }

  toast({
    title: "success.",
  });
  emits("commited");
}
</script>

<template>
  <form class="space-y-6" @submit="onSubmit">
    <FormField v-slot="{ componentField }" name="priceMode">
      <FormItem>
        <FormLabel>{{ $t("priceMode") }}:</FormLabel>
        <FormControl>
          <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="normalMode" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("normalMode") }}
              </FormLabel>
            </FormItem>
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="safeMode" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("safeMode") }}
              </FormLabel>
            </FormItem>
          </RadioGroup>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="launchMode">
      <FormItem>
        <FormLabel>{{ $t("launchMode") }}:</FormLabel>
        <FormControl>
          <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="immediateStart" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("immediateStart") }}
              </FormLabel>
            </FormItem>
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="scheduledStart" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("scheduledStart") }}
              </FormLabel>
            </FormItem>
          </RadioGroup>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="note">
      <FormItem>
        <FormLabel>{{ $t("note") }}</FormLabel>
        <FormControl>
          <Textarea class="w-full" v-bind="componentField"></Textarea>
        </FormControl>
      </FormItem>
    </FormField>

    <Button :disabled="loading" type="submit" class="w-full">
      <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
      {{ $t("submit") }}
    </Button>
  </form>
  <Toaster />
</template>
