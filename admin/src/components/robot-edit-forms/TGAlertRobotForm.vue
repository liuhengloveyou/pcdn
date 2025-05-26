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
  RiskManBotModel,
  RiskManService,
  TGAlertRobotModel,
} from "@/services/RiskManService";
import { useMartParamStore } from "@/stores/martParam-store";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "TGAlertRobotFormComponent",
});
const emits = defineEmits(["commited"]);

let { t } = useI18n({ useScope: "global" });
const martParamStore = useMartParamStore();
let loading = ref(false);

let model = ref<TGAlertRobotModel>({
  id: 0,
  chatId: "",
  userName: "",
  availableMoney: "",
  availableToken: "",
  totalMoney: "",
  totalToken: "",
  priceRangeMin: "",
  priceRangeMax: "",
  remindInterval: 10,
  launchMode: "immediateStart",
  isRunning: false,
  startTime: "",
});

const formSchema = toTypedSchema(
  z.object({
    chatId: z.string(),
    userName: z.string(),
    availableMoney: z.string().optional(),
    availableToken: z.string().optional(),
    totalMoney: z.string().optional(),
    totalToken: z.string().optional(),
    priceRangeMin: z.string().optional(),
    priceRangeMax: z.string().optional(),
    remindInterval: z.number().min(10).max(1440),
    launchMode: z.enum(["immediateStart", "scheduledStart"], {
      required_error: "You need to select a notification type.",
    }),
  })
);

const { handleSubmit, setFieldValue } = useForm({
  validationSchema: formSchema,
  initialValues: {
    remindInterval: 10,
    launchMode: "immediateStart",
  },
});

onMounted(async () => {
  const activeMartParam = martParamStore.getActiveMartParam;
  if (activeMartParam == null) {
    return;
  }

  const resp = await RiskManService.LoadOne(
    305,
    activeMartParam.domain,
    activeMartParam.symbol
  );
  if (resp) {
    console.log(">>>", resp);
    if (resp.code == 0) {
      if (resp.data && resp.data.tgAlertBot) {
        model.value = resp.data.tgAlertBot;
        setFieldValue("chatId", model.value.chatId);
        setFieldValue("userName", model.value.userName);
        setFieldValue("availableMoney", model.value.availableMoney);
        setFieldValue("availableToken", model.value.availableToken);
        setFieldValue("totalMoney", model.value.totalMoney);
        setFieldValue("totalToken", model.value.totalToken);
        setFieldValue("priceRangeMin", model.value.priceRangeMin);
        setFieldValue("priceRangeMax", model.value.priceRangeMax);
        setFieldValue("remindInterval", model.value.remindInterval);
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
  data.chatId = m.chatId;
  data.userName = m.userName;
  data.launchMode = m.launchMode;

  let botModel: RiskManBotModel = {
    id: 0,
    botType: 305,
    botName: "TGAlert",
    isRunning: 0,
    martDomain: activeMartParam.domain,
    symbol: activeMartParam.symbol,
    startTime: 0,
    tgAlertBot: data,
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
    <p class="text-sm text-yellow-300">
      * If the balance data is lower than set value, Alarm will be issued.
    </p>
    <p class="text-sm text-yellow-300">
      * if the price data exceeds the set range, Alarm will be issued.
    </p>
    <form class="space-y-4" @submit="onSubmit">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="chatId">
          <FormLabel>Chat id:</FormLabel>
          <FormItem>
            <FormControl>
              <Input
                type="text"
                placeholder="chat_id_1,chat_id_2,chat_id_3,..."
                v-bind="componentField"
              />
            </FormControl>
            <FormDescription> </FormDescription>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>

      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="userName">
          <FormLabel>User name:</FormLabel>
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
        <FormField v-slot="{ componentField }" name="availableMoney">
          <FormLabel>Available Money:</FormLabel>
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
        <FormField v-slot="{ componentField }" name="availableToken">
          <FormLabel>Available Token:</FormLabel>
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
        <FormField v-slot="{ componentField }" name="totalMoney">
          <FormLabel>Total Money:</FormLabel>
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
        <FormField v-slot="{ componentField }" name="totalToken">
          <FormLabel>Total Token:</FormLabel>
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
        <Label> Price Range:</Label>
        <div class="flex row justify-between">
          <FormField v-slot="{ componentField }" name="priceRangeMin">
            <FormItem>
              <FormControl>
                <Input type="text" v-bind="componentField" />
              </FormControl>
              <FormDescription> </FormDescription>
              <FormMessage />
            </FormItem>
          </FormField>
          <span style="text-align: center; align-content: center"> ~ </span>
          <FormField v-slot="{ componentField }" name="priceRangeMax">
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
        <FormField v-slot="{ componentField }" name="remindInterval">
          <FormLabel>Remind Interval(Min):</FormLabel>
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
