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
import {
  NumberField,
  NumberFieldContent,
  NumberFieldDecrement,
  NumberFieldIncrement,
  NumberFieldInput,
} from "@/components/ui/number-field";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { useToast } from "@/components/ui/toast/use-toast";
const { toast } = useToast();
import { toTypedSchema } from "@vee-validate/zod";
import { GenericObject, useForm } from "vee-validate";
import { onMounted, ref } from "vue";
import * as z from "zod";
import {
  PriceMonitorRobotModel,
  RiskManBotModel,
  RiskManService,
} from "@/services/RiskManService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "PriceMonitorRobotComponent",
});
const emits = defineEmits(["commited"]);

let { t } = useI18n({ useScope: "global" });
const martParamStore = useMartParamStore();
let loading = ref(false);

let model = ref<PriceMonitorRobotModel>({
  id: 0,
  mail1: "",
  remindInterval: 10,
  minPrice: "1",
  maxPrice: "2",
  launchMode: "immediateStart",
  isRunning: false,
  startTime: "",
});

const formSchema = toTypedSchema(
  z.object({
    mail1: z.string().email(),
    remindInterval: z.number().min(1).max(1440),
    minPrice: z.string().min(1).max(100000),
    maxPrice: z.string().min(1).max(100000),
    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a notification type.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    remindInterval: 10,
    minPrice: "1",
    maxPrice: "2",
    launchMode: "immediateStart",
  },
});

onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await RiskManService.LoadOne(
    301,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    console.log(">>>", resp);
    if (resp.code == 0) {
      if (resp.data && resp.data.priceMonitorBot) {
        model.value = resp.data.priceMonitorBot;
        setFieldValue("mail1", model.value.mail1);

        setFieldValue("remindInterval", model.value.remindInterval);
        setFieldValue("minPrice", model.value.minPrice);
        setFieldValue("maxPrice", model.value.maxPrice);
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
  data.remindInterval = Number(m.remindInterval);
  data.minPrice = m.minPrice;
  data.maxPrice = m.maxPrice;
  data.mail1 = m.mail1;
  data.launchMode = m.launchMode;

  const minPrice = Number(data.minPrice);
  const maxPrice = Number(data.maxPrice);
  if (
    Number.isNaN(minPrice) ||
    Number.isNaN(maxPrice) ||
    minPrice >= maxPrice
  ) {
    toast({
      variant: "destructive",
      title: "please input <Price range>.",
    });
    return;
  }

  let botModel: RiskManBotModel = {
    id: 0,
    botType: 301,
    botName: "BalanceMonitor",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    priceMonitorBot: data,
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
          <FormLabel>{{ $t("mail") }}:</FormLabel>
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
        <FormField v-slot="{ componentField }" name="remindInterval">
          <FormLabel>{{ $t("interval") }}(Min):</FormLabel>
          <FormItem>
            <NumberField class="gap-2" :min="0" v-bind="componentField">
              <NumberFieldContent>
                <NumberFieldDecrement />
                <FormControl>
                  <NumberFieldInput />
                </FormControl>
                <NumberFieldIncrement />
              </NumberFieldContent>
            </NumberField>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <Label> {{ $t("priceRange") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minPrice">
            <FormItem>
              <FormControl>
                <Input type="text" v-bind="componentField" />
              </FormControl>
              <FormDescription> </FormDescription>
              <FormMessage />
            </FormItem>
          </FormField>
          <span style="text-align: center; align-content: center"> ~ </span>
          <FormField v-slot="{ componentField }" name="maxPrice">
            <FormItem>
              <FormControl>
                <Input type="text" v-bind="componentField" />
              </FormControl>
              <FormDescription> </FormDescription>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>
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
