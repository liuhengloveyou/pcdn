<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import Toaster from "@/components/ui/toast/Toaster.vue";
import { toast } from "@/components/ui/toast";
import { toTypedSchema } from "@vee-validate/zod";
import { h, onMounted, ref } from "vue";
import * as z from "zod";
import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Check, ChevronsUpDown } from "lucide-vue-next";
import { useForm } from "vee-validate";
import { cn } from "@/lib/utils";
import { MartParamModel, MartParamService } from "@/services/MartParamService";
import { useI18n } from "vue-i18n";

defineOptions({
  name: "MartParamComponent",
});
const props = defineProps({
  id: {
    type: Number,
    required: true,
  },
});

const emits = defineEmits(["commited"]);
const { t } = useI18n();
const open = ref(false);
const popoverOpen = ref(false);
const confirmDialog = ref(false);
const loading = ref(false);

const martOptions = ref([
  {
    value: "binance",
    label: "binance(https://www.binance.com)",
  },
  {
    value: "okx",
    label: "okx(https://www.okx.com)",
  },
  { value: "gate", label: "gate(https://www.gate.io)" },
  { value: "mexc", label: "mexc(https://www.mexc.com)" },
  { value: "xt", label: "xt(https://www.xt.com)" },
  {
    value: "bitmart",
    label: "bitmart(https://www.bitmart.com)",
  },
  {
    value: "biconomy",
    label: "biconomy(https://www.biconomy.com)",
  },
]);
let model = ref<MartParamModel>({
  id: 0,
  uid: 0,
  domain: "",
  symbol: "",
  memo: "",
  accessKey: "",
  secretKey: "",
  active: 0,
});

onMounted(async () => {
  if (props.id && props.id > 0) {
    const resp = await MartParamService.LoadOne(props.id);
    if (resp) {
      if (resp.code == 0 && resp.data) {
        model.value = resp.data;
      }
    }
    console.log(">>>mart: ", resp);
  }
});

const formSchema = toTypedSchema(
  z.object({
    domain: z.string({
      required_error: "Please select a mart.",
    }),
    symbol: z.string({
      required_error: "Please input symbol.",
    }),
    accessKey: z.string({
      required_error: "Please input accessKey.",
    }),
    secretKey: z.string({
      required_error: "Please input secretKey.",
    }),
    memo: z.string(),
  })
);

const { handleSubmit, setFieldValue, values } = useForm({
  validationSchema: formSchema,
});

const onSubmit = handleSubmit((values) => {
  if (props.id <= 0) {
    onRealSubmit({
      domain: values.domain,
      symbol: values.symbol,
      memo: values.memo,
      accessKey: values.accessKey,
      secretKey: values.secretKey,
      id: 0,
      uid: 0,
      active: 0,
    });
  }
});

async function onRealSubmit(model: MartParamModel) {
  let data = { ...model };

  if (
    !data.symbol ||
    data.symbol.length === 0 ||
    data.symbol.indexOf("_") <= 0
  ) {
    toast({
      variant: "destructive",
      title: "symbol wrong.",
    });
    return;
  }
  if (!data.accessKey || data.accessKey.length === 0) {
    toast({
      variant: "destructive",
      title: "accessKey wrong.",
    });
    return;
  }
  if (!data.secretKey || data.secretKey.length === 0) {
    toast({
      variant: "destructive",
      title: "secretKey wrong.",
    });
    return;
  }

  loading.value = true;
  const resp = await MartParamService.Set(data);
  loading.value = false;

  console.log("resp: ", resp);
  if (!resp) {
    return;
  }

  if (resp?.code !== 0) {
    toast({
      variant: "destructive",
      title: t("error." + resp.code),
    });
    return;
  }

  toast({
    title: "You submitted the following values:",
    description: h(
      "pre",
      { class: "mt-2 w-[340px] rounded-md bg-slate-950 p-4" },
      h("code", { class: "text-white" }, JSON.stringify(values, null, 2))
    ),
  });
  open.value = false;
  emits("commited", data);
}

async function onDelete() {
  if (props.id <= 0) {
    return;
  }

  confirmDialog.value = true;
}

async function onRealDelete() {
  if (props.id <= 0) {
    return;
  }

  loading.value = true;
  const resp = await MartParamService.Del(props.id);
  loading.value = false;

  console.log("resp: ", resp);
  if (!resp) {
    return;
  }

  if (resp?.code !== 0) {
    toast({
      variant: "destructive",
      title: t("error." + resp.code),
    });
    return;
  }

  open.value = false;
  emits("commited", undefined);
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogTrigger as-child>
      <Button
        class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground hover:bg-primary/90"
      >
        {{ $t(props.id > 0 ? "edit" : "userView.importAPIKEY") }}
      </Button>
    </DialogTrigger>
    <form id="dialogForm" @submit="onSubmit">
      <DialogContent class="space-y-6 p-8">
        <DialogHeader>
          <DialogTitle>{{ $t("bindApiKey") }}</DialogTitle>
        </DialogHeader>

        <FormField name="domain" :value="model.domain">
          <FormItem class="flex flex-col">
            <FormLabel>{{ $t("mart") }}:</FormLabel>
            <Popover v-model:open="popoverOpen">
              <PopoverTrigger as-child>
                <FormControl>
                  <Button
                    variant="outline"
                    role="combobox"
                    :class="
                      cn(
                        'w-full justify-between',
                        !values.domain && 'text-muted-foreground'
                      )
                    "
                  >
                    {{
                      values.domain
                        ? martOptions.find(
                            (mart) => mart.value === values.domain
                          )?.label
                        : "Select mart..."
                    }}
                    <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </FormControl>
              </PopoverTrigger>
              <PopoverContent class="w-full p-2">
                <Command>
                  <CommandList>
                    <CommandGroup>
                      <CommandItem
                        v-for="mart in martOptions"
                        :key="mart.value"
                        :value="mart.label"
                        @select="
                          () => {
                            setFieldValue('domain', mart.value);
                            popoverOpen = false;
                          }
                        "
                      >
                        <Check
                          :class="
                            cn(
                              'mr-2 h-4 w-4',
                              mart.value === values.domain
                                ? 'opacity-100'
                                : 'opacity-0'
                            )
                          "
                        />
                        {{ mart.label }}
                      </CommandItem>
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
          </FormItem>
        </FormField>

        <FormField
          v-slot="{ componentField }"
          name="symbol"
          :value="model.symbol"
        >
          <FormItem>
            <FormLabel> {{ $t("symbol") }}:</FormLabel>
            <FormControl>
              <Input
                type="text"
                placeholder="MOLI_USDT"
                v-bind="componentField"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField
          v-slot="{ componentField }"
          name="accessKey"
          :value="model.accessKey"
        >
          <FormItem>
            <FormLabel> AccessKey:</FormLabel>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField
          v-slot="{ componentField }"
          name="secretKey"
          :value="model.secretKey"
        >
          <FormItem>
            <FormLabel> SecretKey:</FormLabel>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="memo" :value="model.memo">
          <FormItem>
            <FormLabel> MEMO:</FormLabel>
            <FormControl>
              <Input type="text" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter>
          <Button
            v-if="props.id <= 0"
            :disabled="loading"
            type="submit"
            form="dialogForm"
          >
            <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
            {{ $t("submit") }}
          </Button>

          <Button
            v-if="props.id > 0"
            :disabled="loading"
            form="dialogForm"
            class="bg-red-600 dark:bg-red-600 text-white"
            @click="onDelete"
          >
            <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
            {{ $t("delete") }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </form>
  </Dialog>

  <AlertDialog v-model:open="confirmDialog">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Confirm</AlertDialogTitle>
        <AlertDialogDescription>
          Once deleted, it cannot be recovered. Are you sure you want to delete it?
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>Cancel</AlertDialogCancel>
        <AlertDialogAction @click="onRealDelete">Confirm</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>

  <Toaster />
</template>
