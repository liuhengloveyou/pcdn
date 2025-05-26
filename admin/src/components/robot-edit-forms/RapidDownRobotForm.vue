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
  RapidDownBotModel,
} from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "RapidDownRobotComponent",
});
const emits = defineEmits(["commited"]);

let model = ref<RapidDownBotModel>({
  targetPrice: "",
  totalCost: "",
  totalTime: 10,
  minOrderSize: 0,
  maxOrderSize: 0,
  minInterval: 0,
  maxInterval: 0,
  quantMode: "continuousDist",
  launchMode: "immediateStart",
  description: "",
  isRunning: false,
  startTime: ""
});

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
const loading = ref(false);
const formSchema = toTypedSchema(
  z.object({
    targetPrice: z.string().min(1),
    totalCost: z.string().min(1),
    totalTime: z.number().min(1),
    minOrderSize: z.number().min(0),
    maxOrderSize: z.number().min(0),
    minInterval: z.number().min(0),
    maxInterval: z.number().min(0),
    quantMode: z.enum(["continuousDist", "randomDist"]),
    launchMode: z.enum(["immediateStart", "scheduledStart"]),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    targetPrice: "1",
    totalCost: "1",
    totalTime: 10,
    minOrderSize: 1,
    maxOrderSize: 2,
    minInterval: 1, 
    maxInterval: 2,
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
    103,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  console.log("bot: ", resp?.data);

  if (resp && resp.code == 0) {
    if (resp.data && resp.data.rapidDownBot) {
      model.value = resp.data.rapidDownBot;
      setFieldValue("targetPrice", model.value.targetPrice);
      setFieldValue("totalCost", model.value.totalCost);
      setFieldValue("totalTime", model.value.totalTime);
      setFieldValue("minOrderSize", model.value.minOrderSize);
      setFieldValue("maxOrderSize", model.value.maxOrderSize);
      setFieldValue("minInterval", model.value.minInterval);
      setFieldValue("maxInterval", model.value.maxInterval);
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

  data.minInterval = Number(m.minInterval);
  data.maxInterval = Number(m.maxInterval);
  if (
    data.minInterval <= 0 ||
    data.maxInterval <= 0 ||
    data.maxInterval <= data.minInterval
  ) {
    toast({
      variant: "destructive",
      title: "please input <Interval>.",
    });
    return;
  }

  data.launchMode = m.launchMode;
  data.quantMode = m.quantMode;

  let botModel: BotModel = {
    id: 0,
    botType: 103,
    botName: "RapidDown",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    rapidDownBot: data,
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
      <Label>{{ $t("orderInterval") }}</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minInterval">
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
        <FormField v-slot="{ componentField }" name="maxInterval">
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
      {{ $t("submit") }}
    </Button>
  </form>
  <Toaster />
</template>
