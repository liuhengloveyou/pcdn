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
import { BotModel, BotService, HighThrowBargainRobotModel } from "@/services/BotService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "HighThrowBargainRobotForm",
});
const emits = defineEmits(["commited"]);

let model = ref<HighThrowBargainRobotModel>({
  maxSell: "",
  sellPrice: "",
  minSellSize: 0,
  maxSellSize: 0,
  minSellInterval: 0,
  maxSellInterval: 0,
  maxBuy: "",
  buyPrice: "",
  minBuySize: 0,
  maxBuySize: 0,
  minBuyInterval: 0,
  maxBuyInterval: 0,
  launchMode: "immediateStart"
});
const loading = ref(false);
const martParamStore = useMartParamStore();
let { t } = useI18n({ useScope: "global" });
const formSchema = toTypedSchema(
  z.object({
    maxSell: z.string().min(1),
    sellPrice: z.string().min(1),
    minSellSize: z.number().min(0),
    maxSellSize: z.number().min(0),
    minSellInterval: z.number().min(0),
    maxSellInterval: z.number().min(0),
    maxBuy: z.string().min(1),
    buyPrice: z.string().min(1),
    minBuySize: z.number().min(0),
    maxBuySize: z.number().min(0),
    minBuyInterval: z.number().min(0),
    maxBuyInterval: z.number().min(0),
    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a notification type.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    maxSell: "1",
    sellPrice: "1",
    minSellSize: 1,
    maxSellSize: 2,
    minSellInterval: 1,
    maxSellInterval: 2,
    maxBuy: "1",
    buyPrice: "1",
    minBuySize: 1,
    maxBuySize: 2,
    minBuyInterval: 1,
    maxBuyInterval: 2,
    launchMode: "immediateStart",
  },
});


onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await BotService.LoadOne(
    111,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  console.log("bot: ", resp?.data);

  if (resp && resp.code == 0) {
    if (resp.data && resp.data.highThrowBargainBot) {
      model.value = resp.data.highThrowBargainBot;
      setFieldValue("maxSell", model.value.maxSell);
      setFieldValue("sellPrice", model.value.sellPrice);
      setFieldValue("minSellSize", model.value.minSellSize);
      setFieldValue("maxSellSize", model.value.maxSellSize);
      setFieldValue("minSellInterval", model.value.minSellInterval);
      setFieldValue("maxSellInterval", model.value.maxSellInterval);
      setFieldValue("maxBuy", model.value.maxBuy);
      setFieldValue("buyPrice", model.value.buyPrice);
      setFieldValue("minBuySize", model.value.minBuySize);
      setFieldValue("maxBuySize", model.value.maxBuySize);
      setFieldValue("minBuyInterval", model.value.minBuyInterval);
      setFieldValue("maxBuyInterval", model.value.maxBuyInterval);
      setFieldValue("launchMode", model.value.launchMode);
    }
  }
});


const onSubmit = handleSubmit((values) => {
  onRealSubmit(values);

  // toast({
  //   variant: "destructive",
  //   title: "You submitted the following values:",
  //   description: h(
  //     "pre",
  //     { class: "mt-2 w-[340px] rounded-md bg-slate-950 p-4" },
  //     h("code", { class: "text-white" }, JSON.stringify(values, null, 2))
  //   ),
  // });
});


async function onRealSubmit(m: GenericObject) {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  let data = { ...model.value };

  data.maxSell = m.maxSell;
  const maxSell = Number(data.maxSell);
  if (Number.isNaN(maxSell) || maxSell <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Max Sell>.",
    });
    return;
  }

  data.sellPrice = m.sellPrice;
  const sellPrice = Number(data.sellPrice);
  if (Number.isNaN(sellPrice) || sellPrice <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Sell Price>.",
    });
    return;
  }

  data.minSellSize = m.minSellSize;
  data.maxSellSize = m.maxSellSize;
  const minSellSize = Number(data.minSellSize);
  const maxSellSize = Number(data.maxSellSize);
  if (Number.isNaN(minSellSize) ||
    minSellSize <= 0 ||
    Number.isNaN(maxSellSize) ||
    maxSellSize <= 0 ||
    minSellSize > maxSellSize) {
    toast({
      variant: "destructive",
      title: "please input <Sell Size>.",
    });
    return;
  }

  data.minSellInterval = m.minSellInterval;
  data.maxSellInterval = m.maxSellInterval;
  const minSellInterval = Number(data.minSellInterval);
  const maxSellInterval = Number(data.maxSellInterval);
  if (Number.isNaN(minSellInterval) ||
    minSellInterval <= 0 ||
    Number.isNaN(maxSellInterval) ||
    maxSellInterval <= 0 ||
    minSellInterval > maxSellInterval) {
    toast({
      variant: "destructive",
      title: "please input <Sell Interval>.",
    });
    return;
  }

  data.maxBuy = m.maxBuy;
  const maxBuy = Number(data.maxBuy);
  if (Number.isNaN(maxBuy) || maxBuy <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Max Buy>.",
    });
    return;
  }

  data.buyPrice = m.buyPrice;
  const buyPrice = Number(data.buyPrice);
  if (Number.isNaN(buyPrice) || buyPrice <= 0) {
    toast({
      variant: "destructive",
      title: "please input <Buy Price>.",
    });
    return;
  }
  data.minBuySize = m.minBuySize;
  data.maxBuySize = m.maxBuySize;
  const minBuySize = Number(data.minBuySize);
  const maxBuySize = Number(data.maxBuySize);
  if (Number.isNaN(minBuySize) ||
    minBuySize <= 0 ||
    Number.isNaN(maxBuySize) ||
    maxBuySize <= 0 ||
    minBuySize > maxBuySize) {
    toast({
      variant: "destructive",
      title: "please input <Buy Size>.",
    });
    return;
  }

  data.minBuyInterval = m.minBuyInterval;
  data.maxBuyInterval = m.maxBuyInterval;
  const minBuyInterval = Number(data.minBuyInterval);
  const maxBuyInterval = Number(data.maxBuyInterval);
  if (Number.isNaN(minBuyInterval) ||
    minBuyInterval <= 0 ||
    Number.isNaN(maxBuyInterval) ||
    maxBuyInterval <= 0 ||
    minBuyInterval > maxBuyInterval) {
    toast({
      variant: "destructive",
      title: "please input <Buy Size>.",
    });
    return;
  }

  data.launchMode = m.launchMode;

  let botModel: BotModel = {
    id: 0,
    botType: 111,
    botName: "HighThrowBargain",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    highThrowBargainBot: data,
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
    <form class="space-y-4" @submit="onSubmit">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="sellPrice">
          <FormLabel>{{ $t("sellPrice") }}: </FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="maxSell">
          <FormLabel>{{ $t("maxSell") }}(Token): </FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <Label>{{ $t("sellOrderSize") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minSellSize">
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
          <FormField v-slot="{ componentField }" name="maxSellSize">
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
        <Label>{{ $t("sellInterval") }}(Min):</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minSellInterval">
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
          <FormField v-slot="{ componentField }" name="maxSellInterval">
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


      <!-- buy -->
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="buyPrice">
          <FormLabel>{{ $t("buyPrice") }}: </FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="maxBuy">
          <FormLabel>{{ $t("maxBuy") }}(U): </FormLabel>
          <FormItem>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <Label>{{ $t("buyOrderSize") }}:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minBuySize">
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
          <FormField v-slot="{ componentField }" name="maxBuySize">
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
        <Label>{{ $t("buyInterval") }}(Min):</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="minBuyInterval">
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
          <FormField v-slot="{ componentField }" name="maxBuyInterval">
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

      <!-- <FormField v-slot="{ componentField }" name="description">
        <FormItem>
          <FormLabel>{{ $t("note") }}</FormLabel>
          <FormControl>
            <Textarea class="w-full" v-bind="componentField"></Textarea>
          </FormControl>
        </FormItem>
      </FormField> -->

      <Button :disabled="loading" type="submit" class="w-full">
        <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
        {{ $t("submit") }}
      </Button>
    </form>
  </ScrollArea>
  <Toaster />
</template>
