<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { Label } from "@/components/ui/label";
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
import { Input } from "@/components/ui/input";
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
  GradualDownBotModel,
} from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "GradualDownRobotComponent",
});
const emits = defineEmits(["commited"]);

let model = ref<GradualDownBotModel>({
  targetPrice: "",
  totalCost: "",
  totalTime: 100,
  minOrderSize: 10,
  maxOrderSize: 100,
  quantMode: "continuousDist",
  launchMode: "immediateStart",
  description: "",
  isRunning: false,
  startTime: "",
});

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
const loading = ref(false);
const formSchema = toTypedSchema(
  z.object({
    targetPrice: z.string().min(1),
    totalCost: z.string().min(1),
    totalTime: z.number().min(1),
    minOrderSize: z.number().min(0).max(1000000),
    maxOrderSize: z.number().min(0).max(1000000),
    quantMode: z.enum(["continuousDist", "randomDist"]),
    launchMode: z.enum(["immediateStart", "scheduledStart"]),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    targetPrice: "1",
    totalCost: "1",
    totalTime: 100,
    minOrderSize: 1,
    maxOrderSize: 2,
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
    101,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  console.log("bot: ", resp?.data);

  if (resp && resp.code == 0) {
    if (resp.data && resp.data.gradualDownBot) {
      model.value = resp.data.gradualDownBot;
      setFieldValue("targetPrice", model.value.targetPrice);
      setFieldValue("totalCost", model.value.totalCost);
      setFieldValue("totalTime", model.value.totalTime);
      setFieldValue("minOrderSize", model.value.minOrderSize);
      setFieldValue("maxOrderSize", model.value.maxOrderSize);
      setFieldValue("launchMode", model.value.launchMode);
      setFieldValue("quantMode", model.value.quantMode);
    }
  }
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

  data.targetPrice = m.targetPrice;
  const targetPrice = Number(data.targetPrice);
  if (Number.isNaN(targetPrice) || targetPrice <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Target Price>.",
    });
    return;
  }

  data.totalCost = m.totalCost;
  const totalCost = Number(data.totalCost);
  if (Number.isNaN(totalCost) || totalCost <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Total Cost>.",
    });
    return;
  }

  data.totalTime = m.totalTime;
  const totalTime = Number(data.totalTime);
  if (Number.isNaN(totalTime) || totalTime <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Total Time>.",
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

  // data.minTimeInterval = Number(m.minTimeInterval);
  // data.maxTimeInterval = Number(m.maxTimeInterval);
  // if (
  //   data.minTimeInterval <= 0 ||
  //   data.maxTimeInterval <= 0 ||
  //   data.maxTimeInterval <= data.minTimeInterval
  // ) {
  //   toast({
  //     variant: "destructive",
  //     title: "please input <interval>.",
  //   });
  //   return;
  // }

  data.launchMode = m.launchMode;
  data.quantMode = m.quantMode;

  let botModel: BotModel = {
    id: 0,
    botType: 101,
    botName: "GradualDown",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    gradualDownBot: data,
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
    <FormField v-slot="{ componentField }" name="targetPrice">
      <FormItem>
        <FormLabel>{{ $t("targetPrice") }}:</FormLabel>
        <FormControl>
          <Input type="text" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="totalCost">
      <FormItem>
        <FormLabel>{{ $t("totalCost") }}(U):</FormLabel>
        <FormControl>
          <Input type="text" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="totalTime">
      <FormItem>
        <FormLabel>{{ $t("totalTime") }}(Min):</FormLabel>
        <FormControl>
          <Input type="number" v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <div class="space-y-2">
      <Label>{{ $t("orderSize") }}</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minOrderSize">
          <FormItem>
            <FormControl>
              <Input type="number" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
        <span style="text-align: center; align-content: center"> ~ </span>
        <FormField v-slot="{ componentField }" name="maxOrderSize">
          <FormItem>
            <FormControl>
              <Input type="number" v-bind="componentField" />
            </FormControl>
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
      {{ $t("submit") }}
    </Button>
  </form>
  <Toaster />
</template>
