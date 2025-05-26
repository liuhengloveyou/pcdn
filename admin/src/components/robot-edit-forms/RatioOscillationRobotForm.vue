<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription
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
import { BotModel, BotService, RatioOscillationBotModel } from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "RatioOscillationRobotComponent",
});
const emits = defineEmits(["commited"]);

let model = ref<RatioOscillationBotModel>({
  totalCost: "100",
  totalTime: 1,
  minOrderSize: 0,
  maxOrderSize: 0,
  minOscillate: 0,
  maxOscillate: 0,
  quantMode: "continuousDist",
  launchMode: "immediateStart",
  description: "",
});

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
const loading = ref(false);
const formSchema = toTypedSchema(
  z.object({
    totalCost: z.string().min(1),
    totalTime: z.number().min(1),
    minOrderSize: z.number().min(0),
    maxOrderSize: z.number().min(0),
    minOscillate: z.number().min(1).max(100),
    maxOscillate: z.number().min(1).max(100),
    quantMode: z.enum(["continuousDist", "randomDist"]),
    launchMode: z.enum(["immediateStart", "scheduledStart"]),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    totalCost: "100",
    totalTime: 30,
    minOrderSize: 100,
    maxOrderSize: 200,
    minOscillate: 5,
    maxOscillate: 10,
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
    112,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  console.log("bot: ", resp?.data);

  if (resp && resp.code == 0) {
    if (resp.data && resp.data.ratioOscillationBot) {
      model.value = resp.data.ratioOscillationBot;
      setFieldValue("totalCost", model.value.totalCost);
      setFieldValue("totalTime", model.value.totalTime);
      setFieldValue("minOrderSize", model.value.minOrderSize);
      setFieldValue("maxOrderSize", model.value.maxOrderSize);
      setFieldValue("minOscillate", model.value.minOscillate);
      setFieldValue("maxOscillate", model.value.maxOscillate);
      setFieldValue("launchMode", model.value.launchMode);
      setFieldValue("quantMode", model.value.quantMode);
    }
  }
});

const onSubmit = handleSubmit((values) => {
  onRealSubmit(values);
});

async function onRealSubmit(m: GenericObject) {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  let data = { ...model.value };

  data.totalCost = m.totalCost;
  const totalCost = Number(data.totalCost);
  if (Number.isNaN(totalCost) || totalCost <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Max Cost>.",
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

  data.minOscillate = Number(m.minOscillate);
  data.maxOscillate = Number(m.maxOscillate);
  if (
    data.minOscillate <= 0 ||
    data.maxOscillate <= 0 ||
    data.maxOscillate <= data.minOscillate ||
    data.minOscillate > 100 ||
    data.maxOscillate > 100
  ) {
    toast({
      variant: "destructive",
      title: "please input valid <Oscillate Range> (1-100%).",
    });
    return;
  }

  // data.oscillateRatio = Number(m.oscillateRatio);
  // if (Number.isNaN(data.oscillateRatio) || data.oscillateRatio <= 0 || data.oscillateRatio > 100) {
  //   toast({
  //     variant: "destructive",
  //     title: "please input valid <Oscillate Ratio> (1-100%).",
  //   });
  //   return;
  // }

  data.launchMode = m.launchMode;
  data.quantMode = m.quantMode;

  let botModel: BotModel = {
    id: 0,
    botType: 112,
    botName: "RatioOscillation",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    ratioOscillationBot: data,
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
    <FormField v-slot="{ componentField }" name="totalCost">
      <FormItem>
        <FormLabel>{{ $t("totalCost") }}(U):</FormLabel>
        <FormControl>
          <Input type="text" v-bind="componentField" />
        </FormControl>
        <FormDescription class="text-xs">{{ $t("totalCostDesc") }}</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="totalTime">
      <FormItem>
        <FormLabel>{{ $t("totalTime") }}(Min):</FormLabel>
        <FormControl>
          <Input type="number" v-bind="componentField" />
        </FormControl>
        <FormDescription class="text-xs">{{ $t("totalTimeDesc") }}</FormDescription>
      </FormItem>
    </FormField>

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
            <FormDescription class="text-xs">{{ $t("orderSizeDesc") }}</FormDescription>
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
      <Label>{{ $t("oscillateRange") }}(%):</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minOscillate">
          <FormItem>
            <NumberField class="gap-2" :min="1" :max="100" v-bind="componentField">
              <NumberFieldContent>
                <NumberFieldDecrement />
                <FormControl>
                  <NumberFieldInput />
                </FormControl>
                <NumberFieldIncrement />
              </NumberFieldContent>
            </NumberField>
            <FormDescription class="text-xs">{{ $t("oscillateRangeDesc") }}</FormDescription>
          </FormItem>
        </FormField>
        <span style="text-align: center; align-content: center"> ~ </span>
        <FormField v-slot="{ componentField }" name="maxOscillate">
          <FormItem>
            <NumberField class="gap-2" :min="1" :max="100" v-bind="componentField">
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

    <!-- <FormField v-slot="{ componentField }" name="note">
      <FormItem>
        <FormLabel>{{ $t("description") }}:</FormLabel>
        <FormControl>
          <Textarea v-model="model.description"  v-bind="componentField"/>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField> -->

    <Button type="submit" :disabled="loading" class="w-full">
      <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
      {{ $t("submit") }}
    </Button>
  </form>
  <Toaster />
</template>
