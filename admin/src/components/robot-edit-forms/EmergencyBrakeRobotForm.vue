<script setup lang="ts">
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { useToast } from "@/components/ui/toast/use-toast";
const { toast } = useToast();
import { toTypedSchema } from "@vee-validate/zod";
import { GenericObject, useForm } from "vee-validate";
import { onMounted, ref } from "vue";
import * as z from "zod";
import {
  EmergencyBrakeRobotModel,
  RiskManBotModel,
  RiskManService,
} from "@/services/RiskManService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "EmergencyBrakeRobotComponent",
});
const emits = defineEmits(["commited"]);

let { t } = useI18n({ useScope: "global" });
const martParamStore = useMartParamStore();
let loading = ref(false);

let model = ref<EmergencyBrakeRobotModel>({
  id: 0,
  mail1: "",
  mail2: "",
  mail3: "",
  balanceAlertThreshold: "1",
  tokenAlertThreshold: "1",
  launchMode: "immediateStart",
  isRunning: false,
  startTime: "",
});

const formSchema = toTypedSchema(
  z.object({
    mail1: z.string().email(),
    mail2: z.string().optional(),
    mail3: z.string().optional(),
    balanceAlertThreshold: z.string().min(1),
    tokenAlertThreshold: z.string().min(1),
    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a notification type.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    balanceAlertThreshold: "1",
    tokenAlertThreshold: "1",
    launchMode: "immediateStart",
  },
});

onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await RiskManService.LoadOne(
    303,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    console.log(">>>", resp);
    if (resp.code == 0) {
      if (resp.data && resp.data.emergencyBrakeBot) {
        model.value = resp.data.emergencyBrakeBot;
        setFieldValue("mail1", model.value.mail1);
        setFieldValue("mail2", model.value.mail2);
        setFieldValue("mail3", model.value.mail3);
        setFieldValue(
          "balanceAlertThreshold",
          model.value.balanceAlertThreshold
        );
        setFieldValue("tokenAlertThreshold", model.value.tokenAlertThreshold);
        setFieldValue("launchMode", model.value.launchMode);
      }
    }
  }
});

const onSubmit = handleSubmit((values) => {
  // console.log(">?>>", values);
  // toast({
  //   variant: "destructive",
  //   title: "You submitted the following values:",
  //   description: h(
  //     "pre",
  //     { class: "mt-2 w-[340px] rounded-md bg-slate-950 p-4" },
  //     h("code", { class: "text-white" }, JSON.stringify(values, null, 2))
  //   ),
  // });

  onRealSubmit(values);
});

async function onRealSubmit(m: GenericObject) {
  const activeMartParam = martParamStore.getActiveMartParam;

  if (activeMartParam == null) {
    return;
  }

  let data = { ...model.value };
  data.mail1 = m.mail1;
  data.mail2 = m.mail2;
  data.mail3 = m.mail3;
  data.balanceAlertThreshold = m.balanceAlertThreshold;
  data.tokenAlertThreshold = m.tokenAlertThreshold;
  data.launchMode = m.launchMode;

  let botModel: RiskManBotModel = {
    id: 0,
    botType: 303,
    botName: "EmergencyBrake",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    emergencyBrakeBot: data,
  };
  console.log(">>>", botModel);
  loading.value = true;
  const resp = await RiskManService.Set(botModel);
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
  <ScrollArea class="h-[90vh] w-full">
    <form class="space-y-4" @submit="onSubmit">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="mail1">
          <FormLabel>Mail 1:</FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormDescription> </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="mail2">
          <FormLabel>Mail 2:</FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormDescription> </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="mail3">
          <FormLabel>Mail 3:</FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormDescription> </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="balanceAlertThreshold">
          <FormLabel>Balance Alert Threshold:</FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormDescription> </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="tokenAlertThreshold">
          <FormLabel>Token Alert Threshold:</FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormDescription> </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

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

      <Button :disabled="loading" type="submit" class="w-full">
        <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
        {{ $t("submit") }}
      </Button>
    </form>
  </ScrollArea>
  <Toaster />
</template>
