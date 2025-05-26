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
import { BotModel, BotService, DepthNearBotModel } from "@/services/BotService";
import { useI18n } from "vue-i18n";
import { useMartParamStore } from "@/stores/martParam-store";

defineOptions({
  name: "DepthNearRobotComponent",
});
const emits = defineEmits(["commited"]);

const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
let loading = ref(false);
let model = ref<DepthNearBotModel>({
  id: 0,
  minBuyOrderSize: 0,
  maxBuyOrderSize: 0,
  minBuySheet: 0,
  maxBuySheet: 0,
  minSellOrderSize: 0,
  maxSellOrderSize: 0,
  minSellSheet: 0,
  maxSellSheet: 0,
  minBuySpread: 0,
  maxBuySpread: 0,
  minSellSpread: 0,
  maxSellSpread: 0,
  minDuration: 0,
  maxDuration: 0,
  orderMode: "random",
  priceMode: "normalMode",
  launchMode: "immediateStart",
  description: "",
  isRunning: false,
  startTime: "",
});

const formSchema = toTypedSchema(
  z.object({
    minBuyOrderSize: z.number().min(1),
    maxBuyOrderSize: z.number().min(1),
    minBuySheet: z.number().min(1),
    maxBuySheet: z.number().min(1),
    minSellOrderSize: z.number().min(1),
    maxSellOrderSize: z.number().min(1),
    minSellSheet: z.number().min(1),
    maxSellSheet: z.number().min(1),
    minBuySpread: z.number().min(1),
    maxBuySpread: z.number().min(1),
    minSellSpread: z.number().min(1),
    maxSellSpread: z.number().min(1),
    minDuration: z.number().min(1),
    maxDuration: z.number().min(1),

    orderMode: z.enum(["random", "buyFirst", "sellFirst"], {
      required_error: "You need to select a order mode.",
    }),

    priceMode: z.enum(["normalMode", "safeMode"], {
      required_error: "You need to select a price mode.",
    }),

    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a lauch mode.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    minBuyOrderSize: 1,
    maxBuyOrderSize: 2,
    minBuySheet: 1,
    maxBuySheet: 2,
    minSellOrderSize: 1,
    maxSellOrderSize: 2,
    minSellSheet: 1,
    maxSellSheet: 2,
    minBuySpread: 1,
    maxBuySpread: 2,
    minSellSpread: 1,
    maxSellSpread: 2,
    minDuration: 1,
    maxDuration: 2,
    orderMode: "random",
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
    9,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    if (resp.code == 0) {
      if (resp.data && resp.data.depthNearBot) {
        model.value = resp.data.depthNearBot;

        setFieldValue("minBuyOrderSize", model.value.minBuyOrderSize);
        setFieldValue("maxBuyOrderSize", model.value.maxBuyOrderSize);
        setFieldValue("minBuySheet", model.value.minBuySheet);
        setFieldValue("maxBuySheet", model.value.maxBuySheet);
        setFieldValue("minSellOrderSize", model.value.minSellOrderSize);
        setFieldValue("maxSellOrderSize", model.value.maxSellOrderSize);
        setFieldValue("minSellSheet", model.value.minSellSheet);
        setFieldValue("maxSellSheet", model.value.maxSellSheet);
        setFieldValue("minBuySpread", model.value.minBuySpread);
        setFieldValue("maxBuySpread", model.value.maxBuySpread);
        setFieldValue("minSellSpread", model.value.minSellSpread);
        setFieldValue("maxSellSpread", model.value.maxSellSpread);
        setFieldValue("minDuration", model.value.minDuration);
        setFieldValue("maxDuration", model.value.maxDuration);
        setFieldValue("orderMode", model.value.orderMode);
        setFieldValue("priceMode", model.value.priceMode);
        setFieldValue("launchMode", model.value.launchMode);
      }
    }
  }

  console.log("bot: ", model.value);
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
  data.minBuyOrderSize = Number(m.minBuyOrderSize);
  data.maxBuyOrderSize = Number(m.maxBuyOrderSize);
  data.minBuySheet = Number(m.minBuySheet);
  data.maxBuySheet = Number(m.maxBuySheet);
  data.minSellOrderSize = Number(m.minSellOrderSize);
  data.maxSellOrderSize = Number(m.maxSellOrderSize);
  data.minSellSheet = Number(m.minSellSheet);
  data.maxSellSheet = Number(m.maxSellSheet);
  data.minBuySpread = Number(m.minBuySpread);
  data.maxBuySpread = Number(m.maxBuySpread);
  data.minSellSpread = Number(m.minSellSpread);
  data.maxSellSpread = Number(m.maxSellSpread);
  data.minDuration = Number(m.minDuration);
  data.maxDuration = Number(m.maxDuration);

  data.launchMode = m.launchMode;
  let botModel: BotModel = {
    id: 0,
    botType: 9,
    botName: "DepthNear",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    depthNearBot: data,
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
      <Label>{{ $t("buyOrderSize") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minBuyOrderSize">
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
        <FormField v-slot="{ componentField }" name="maxBuyOrderSize">
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
      <Label>{{ $t("buySheet") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minBuySheet">
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
        <FormField v-slot="{ componentField }" name="maxBuySheet">
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
      <Label>{{ $t("buySpread") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minBuySpread">
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
        <FormField v-slot="{ componentField }" name="maxBuySpread">
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
      <Label>{{ $t("sellOrderSize") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minSellOrderSize">
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
        <FormField v-slot="{ componentField }" name="maxSellOrderSize">
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
      <Label>{{ $t("sellSheet") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minSellSheet">
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
        <FormField v-slot="{ componentField }" name="maxSellSheet">
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
      <Label>{{ $t("sellSpread") }}:</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minSellSpread">
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
        <FormField v-slot="{ componentField }" name="maxSellSpread">
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
      <Label>{{ $t("duration") }}(s):</Label>
      <div class="flex row justify-between">
        <FormField v-slot="{ componentField }" name="minDuration">
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
        <FormField v-slot="{ componentField }" name="maxDuration">
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

    <FormField v-slot="{ componentField }" name="orderMode">
      <FormItem>
        <FormLabel>{{ $t("orderMode") }}:</FormLabel>
        <FormControl>
          <RadioGroup class="flex flex-col space-y-1" v-bind="componentField">
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="random" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("random") }}
              </FormLabel>
            </FormItem>
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="buyFirst" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("buyFirst") }}
              </FormLabel>
            </FormItem>
            <FormItem class="flex items-center space-y-0 gap-x-3">
              <FormControl>
                <RadioGroupItem value="sellFirst" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ $t("sellFirst") }}
              </FormLabel>
            </FormItem>
          </RadioGroup>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

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
      Submit
    </Button>
  </form>
  <Toaster />
</template>
