<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
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
  BotModel,
  BotService,
  OrderbookFlashBotModel,
} from "@/services/BotService";
import { useI18n } from "vue-i18n";
import { useMartParamStore } from "@/stores/martParam-store";

defineOptions({
  name: "OrderBookFlashRobotComponent",
});
const emits = defineEmits(["commited"]);

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
let loading = ref(false);
let model = ref<OrderbookFlashBotModel>({
  id: 0,
  minTimeInterval: 1,
  maxTimeInterval: 1,
  minGear: 1,
  maxGear: 1,
  minOrderSize: 100,
  maxOrderSize: 200,
  minDivtions: 1,
  maxDivtions: 1,
  launchMode: "immediateStart",
  description: "",
  isRunning: false,
  startTime: "",
});

const formSchema = toTypedSchema(
  z.object({
    // 时间间隔 s
    minTimeInterval: z.number().min(1).max(3600),
    maxTimeInterval: z.number().min(1).max(3600),

    minGear: z.number().min(1).max(100),
    maxGear: z.number().min(1).max(100),
    minOrderSize: z.number().min(1),
    maxOrderSize: z.number().min(1),
    minDivtions: z.number().min(1),
    maxDivtions: z.number().min(1),

    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a lauch mode.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    // 时间间隔 s
    minTimeInterval: 5,
    maxTimeInterval: 10,

    // 檔位範圍
    minGear: 1,
    maxGear: 5,

    // 單筆數量
    minOrderSize: 100,
    maxOrderSize: 200,

    // 分筆範圍
    minDivtions: 1,
    maxDivtions: 10,

    launchMode: "immediateStart",
  },
});

onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await BotService.LoadOne(
    7,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    if (resp.code == 0) {
      if (resp.data && resp.data.orderbookFlashBot) {
        model.value = resp.data.orderbookFlashBot;

        setFieldValue("minGear", model.value.minGear);
        setFieldValue("maxGear", model.value.maxGear);
        setFieldValue("minTimeInterval", model.value.minTimeInterval);
        setFieldValue("maxTimeInterval", model.value.maxTimeInterval);
        setFieldValue("minOrderSize", model.value.minOrderSize);
        setFieldValue("maxOrderSize", model.value.maxOrderSize);
        setFieldValue("minDivtions", model.value.minDivtions);
        setFieldValue("maxDivtions", model.value.maxDivtions);
        setFieldValue("launchMode", model.value.launchMode);
      }
    }
  }

  console.log("OrderBookFlashRobotComponent>>>bot: ", model);
});

const onSubmit = handleSubmit((values) => {
  // console.log(">>>", values);
  // toast({
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

  data.minTimeInterval = Number(m.minTimeInterval);
  data.maxTimeInterval = Number(m.maxTimeInterval);
  if (
    data.minTimeInterval <= 0 ||
    data.maxTimeInterval <= 0 ||
    data.maxTimeInterval <= data.minTimeInterval
  ) {
    toast({
      variant: "destructive",
      title: "please input <interval>.",
    });
    return;
  }

  data.minOrderSize = Number(m.minOrderSize);
  data.maxOrderSize = Number(m.maxOrderSize);
  if (
    data.minOrderSize <= 0 ||
    data.maxOrderSize <= 0 ||
    data.maxOrderSize <= data.minOrderSize
  ) {
    toast({
      variant: "destructive",
      title: "please input <Order Size>.",
    });
    return;
  }

  data.minGear = Number(m.minGear);
  data.maxGear = Number(m.maxGear);
  if (
    data.minGear <= 0 ||
    data.maxGear <= 0 ||
    data.maxGear <= data.minGear
  ) {
    toast({
      variant: "destructive",
      title: "please input <Gear Range>.",
    });
    return;
  }

  data.minDivtions = Number(m.minDivtions);
  data.maxDivtions = Number(m.maxDivtions);
  if (
    data.minDivtions <= 0 ||
    data.maxDivtions <= 0 ||
    data.maxDivtions <= data.minDivtions
  ) {
    toast({
      variant: "destructive",
      title: "please input <Divtions>.",
    });
    return;
  }
  data.launchMode = m.launchMode;

  let botModel: BotModel = {
    id: 0,
    botType: 7,
    botName: "OrderbookFlash",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    orderbookFlashBot: data,
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
  <ScrollArea class="h-[90vh] w-full">
    <form class="space-y-6" @submit="onSubmit">
      <div class="space-y-2">
        <Label>{{ $t("orderInterval") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minTimeInterval">
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
          <span style="text-align: center; align-content: center"> ~ </span>
          <FormField v-slot="{ componentField }" name="maxTimeInterval">
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
      </div>

      <div class="space-y-2">
        <Label>{{ $t("gearRange") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minGear">
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
          <span style="text-align: center; align-content: center"> ~ </span>
          <FormField v-slot="{ componentField }" name="maxGear">
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
      </div>

      <div class="space-y-2">
        <Label>{{ $t("orderSize") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minOrderSize">
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
          <span style="text-align: center; align-content: center"> ~ </span>
          <FormField v-slot="{ componentField }" name="maxOrderSize">
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
      </div>

      <div class="space-y-2">
        <Label>{{ $t("divitions") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minDivtions">
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
          <span style="text-align: center; align-content: center"> ~ </span>
          <FormField v-slot="{ componentField }" name="maxDivtions">
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
        Submit
      </Button>
    </form>
  </ScrollArea>
  <Toaster />
</template>
