<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Input } from "@/components/ui/input";
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
import { GenericObject, useForm } from "vee-validate";
import { onMounted, ref } from "vue";
import {
  BotModel,
  BotService,
  RangeBrushBotModel,
} from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";
import * as z from "zod";
import { toTypedSchema } from "@vee-validate/zod";

defineOptions({
  name: "EditRangeBrushRobotComponent",
});
const emits = defineEmits(["commited"]);

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
let loading = ref(false);
let model = ref<RangeBrushBotModel>({
  id: 0,
  // 单笔交易量：
  minNumPeerOrder: 10,
  maxNumPeerOrder: 100,

  minPrice: "1",
  maxPrice: "2",

  // 时间间隔 s
  minTimeInterval: 100,
  maxTimeInterval: 200,

  // 启动模式：
  launchMode: "immediateStart", // 立即启动 / 定时启动
  // 数量分布类型：
  quantMode: "continuousDist", // 连续分布 / 随机分布

  // 机器人描述：
  description: "",
  isRunning: false,
  startTime: "",
});
const formSchema = toTypedSchema(
  z.object({
    // 单笔交易量：
    minNumPeerOrder: z.number().min(1),
    maxNumPeerOrder: z.number().min(1),

    // 价格区间
    minPrice: z.string(),
    maxPrice: z.string(),

    // 时间间隔 s
    minTimeInterval: z.number().min(1),
    maxTimeInterval: z.number().min(1),

    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a lauch mode.",
    }),

    quantMode: z.enum(["continuousDist", "randomDist"], {
      required_error: "You need to select a quant mode.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    minNumPeerOrder: 10,
    maxNumPeerOrder: 20,

    minPrice: "1",
    maxPrice: "2",

    // 时间间隔 s
    minTimeInterval: 10,
    maxTimeInterval: 200,

    launchMode: "immediateStart",
    quantMode: "continuousDist",
  },
});

onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await BotService.LoadOne(
    5,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    if (resp.code == 0) {
      if (resp.data && resp.data.rangeBrushBot) {
        model.value = resp.data.rangeBrushBot;
        setFieldValue("minNumPeerOrder", model.value.minNumPeerOrder);
        setFieldValue("maxNumPeerOrder", model.value.maxNumPeerOrder);
        setFieldValue("minTimeInterval", model.value.minTimeInterval);
        setFieldValue("maxTimeInterval", model.value.maxTimeInterval);
        setFieldValue("minPrice", model.value.minPrice);
        setFieldValue("maxPrice", model.value.maxPrice);

        setFieldValue("launchMode", model.value.launchMode);
        setFieldValue("quantMode", model.value.quantMode);
      }
    }
  }

  console.log("RangeBrushRobotComponent>>>bot: ", model);
});

const onSubmit = handleSubmit((values) => {
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
  data.minNumPeerOrder = Number(m.minNumPeerOrder);
  data.maxNumPeerOrder = Number(m.maxNumPeerOrder);
  if (
    data.minNumPeerOrder <= 0 ||
    data.maxNumPeerOrder <= 0 ||
    data.maxNumPeerOrder <= data.minNumPeerOrder
  ) {
    toast({
      variant: "destructive",
      title: "please input <Order Size>.",
    });
    return;
  }

  data.minPrice = m.minPrice;
  data.maxPrice = m.maxPrice;

  if (
    Number.isNaN(Number(data.minPrice)) ||
    Number.isNaN(Number(data.maxPrice)) ||
    Number(data.minPrice) <= 0 ||
    Number(data.maxPrice) <= 0 ||
    (Number(data.minPrice) >= Number(data.maxPrice))
  ) {
    toast({
      variant: "destructive",
      title: "please input Price range.",
    });
    return;
  }

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
  data.launchMode = m.launchMode;
  data.quantMode = m.quantMode;

  let botModel: BotModel = {
    id: 0,
    botType: 5,
    botName: "RangeBrush",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    rangeBrushBot: data,
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
    <div class="space-y-2">
      <Label>{{ $t("orderSize") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minNumPeerOrder">
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
        <FormField v-slot="{ componentField }" name="maxNumPeerOrder">
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
      <Label>{{ $t("priceRange") }}:</Label>
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

    <FormField v-slot="{ componentField }" name="quantMode">
      <FormItem>
        <FormLabel>{{ $t("amountDistributionType") }}:</FormLabel>
        <FormControl>
          <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="continuousDist" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("continuousDistribution") }}
              </FormLabel>
            </FormItem>
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="randomDist" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("randomDistribution") }}
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
  <Toaster  />
</template>
